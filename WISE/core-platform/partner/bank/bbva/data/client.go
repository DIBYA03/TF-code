package data

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"time"

	"github.com/google/uuid"
	retryhttp "github.com/hashicorp/go-retryablehttp"
	partnerbank "github.com/wiseco/core-platform/partner/bank"
	"github.com/wiseco/core-platform/shared"
)

type config struct {
	baseURL        string
	httpClient     *retryhttp.Client
	appName        string
	appID          string
	secretOAuthKey string
}

type client struct {
	*config
	oAuth shared.OAuthService
}

func toHttpRequest(r *retryhttp.Request) *http.Request {
	return &http.Request{
		Method:           r.Method,
		URL:              r.URL,
		Proto:            r.Proto,
		ProtoMajor:       r.ProtoMajor,
		ProtoMinor:       r.ProtoMinor,
		Header:           r.Header,
		Body:             r.Body,
		ContentLength:    r.ContentLength,
		TransferEncoding: r.TransferEncoding,
		Close:            r.Close,
		Host:             r.Host,
		Form:             r.Form,
		PostForm:         r.PostForm,
		MultipartForm:    r.MultipartForm,
		Trailer:          r.Trailer,
		RemoteAddr:       r.RemoteAddr,
		RequestURI:       r.RequestURI,
		TLS:              r.TLS,
		Cancel:           r.Cancel,
		Response:         r.Response,
	}
}

type clientResponse struct {
	httpResp *http.Response
	body     interface{}
	err      apiError
	bytes    []byte
	start    time.Time
	end      time.Time
}

func newHTTPClient() *retryhttp.Client {
	client := retryhttp.NewClient()
	client.Logger = nil
	return client
}

var apiSubdomain = map[string]string{
	"sandbox": "sandbox-apis",
	"preprod": "preprod-apis",
	"prod":    "apis",
}

func newClient() *client {
	return &client{bbvaConfig, oAuthSrv}
}

func getBaseAPIURL() string {
	// BBVA App Env
	appEnv := os.Getenv("BBVA_APP_ENV")
	if len(appEnv) == 0 {
		panic(errors.New("BBVA_APP_ENV is missing"))
	}

	// Default to sandbox if none available or incorrect
	subDomain, ok := apiSubdomain[appEnv]
	if !ok {
		panic(errors.New("invalid app environment"))
	}

	return "https://" + subDomain + ".bbvaopenplatform.com"
}

var bbvaConfig *config = func() *config {
	// BBVA App Name
	appName := os.Getenv("BBVA_APP_NAME")
	if len(appName) == 0 {
		panic(errors.New("BBVA_APP_NAME is missing"))
	}

	// BBVA App Id
	appID := os.Getenv("BBVA_APP_ID")
	if len(appID) == 0 {
		panic(errors.New("BBVA_APP_ID is missing"))
	}

	// BBVA OAuth Secret
	appSecret := os.Getenv("BBVA_APP_SECRET")
	if len(appSecret) == 0 {
		panic(errors.New("BBVA_APP_SECRET is missing"))
	}

	return &config{
		baseURL:        getBaseAPIURL(),
		appName:        appName,
		appID:          appID,
		secretOAuthKey: appSecret,
		httpClient:     newHTTPClient(),
	}
}()

func (s *client) newRequest(method, u string, buf *bytes.Buffer) (req *retryhttp.Request, err error) {
	if buf == nil {
		req, err = retryhttp.NewRequest(method, u, nil)
	} else {
		req, err = retryhttp.NewRequest(method, u, buf)
	}
	if err != nil {
		return nil, err
	}
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

func (s *client) newAuthenticatedRequest(method, u string, request partnerbank.APIRequest, buf *bytes.Buffer) (*retryhttp.Request, error) {
	token, err := s.oAuth.Token()
	if err != nil {
		return nil, err
	}

	req, err := s.newRequest(method, u, buf)
	if err != nil {
		return nil, err
	}

	// Source IP Required
	if len(request.SourceIP) == 0 {
		return nil, errors.New("Required: Request Context IP Address")
	}

	// Set custom headers
	req.Header.Set("Authorization", fmt.Sprintf("jwt %s", token.AccessToken))
	req.Header.Set("X-Customer-IP", request.SourceIP)
	req.Header.Set("X-Unique-Transaction-ID", uuid.New().String())

	return req, nil
}

func (s *client) do(req *retryhttp.Request, response interface{}) error {
	// Dump when debug is turned on
	if os.Getenv("DEBUG") != "" {
		dump, _ := httputil.DumpRequest(toHttpRequest(req), true)
		log.Println("Request: ", string(dump))
	}

	resp := &clientResponse{body: response, start: time.Now()}
	httpResp, err := s.httpClient.Do(req)
	resp.end = time.Now()
	if err != nil {
		return fmt.Errorf("http client: %v", err)
	}

	return handleAPIResponse(resp, httpResp)
}

func (s *client) doClientResp(req *retryhttp.Request, resp *clientResponse) error {
	// Dump when debug is turned on
	if os.Getenv("DEBUG") != "" {
		dump, _ := httputil.DumpRequest(toHttpRequest(req), true)
		log.Println("Request: ", string(dump))
	}

	resp.start = time.Now()
	httpResp, err := s.httpClient.Do(req)
	resp.end = time.Now()
	if err != nil {
		return fmt.Errorf("http client: %v", err)
	}

	return handleAPIResponse(resp, httpResp)
}

func (s *client) get(path string, request partnerbank.APIRequest) (*retryhttp.Request, error) {
	u := fmt.Sprintf("%s/%s", s.baseURL, path)
	req, err := s.newAuthenticatedRequest(http.MethodGet, u, request, nil)
	if err != nil {
		return nil, fmt.Errorf("newAuthenticatedRequest: %v", err)
	}
	return req, nil
}

func (s *client) post(path string, request partnerbank.APIRequest, payload interface{}) (*retryhttp.Request, error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(payload); err != nil {
		return nil, err
	}

	// Dump when debug is turned on
	if os.Getenv("DEBUG") != "" {
		log.Println("Post: ", string(buf.Bytes()))
	}

	u := fmt.Sprintf("%s/%s", s.baseURL, path)
	req, err := s.newAuthenticatedRequest(http.MethodPost, u, request, &buf)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (s *client) patch(path string, request partnerbank.APIRequest, payload interface{}) (*retryhttp.Request, error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(payload); err != nil {
		return nil, err
	}

	// Dump when debug is turned on
	if os.Getenv("DEBUG") != "" {
		log.Println("Patch: ", string(buf.Bytes()))
	}

	u := fmt.Sprintf("%s/%s", s.baseURL, path)
	req, err := s.newAuthenticatedRequest(http.MethodPatch, u, request, &buf)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (s *client) delete(path string, request partnerbank.APIRequest) (*retryhttp.Request, error) {
	u := fmt.Sprintf("%s/%s", s.baseURL, path)
	req, err := s.newAuthenticatedRequest(http.MethodDelete, u, request, nil)
	if err != nil {
		return nil, fmt.Errorf("newAuthenticatedRequest: %v", err)
	}
	return req, nil
}

type apiError struct {
	Result APIErrorResult `json:"result"`
	Errors []APIError     `json:"errors"`
}

func (a apiError) Error() error {
	if len(a.Errors) > 0 {
		return fmt.Errorf("%v", a.Errors)
	}
	return fmt.Errorf("%v", a.Result)
}

type APIErrorResult struct {
	Code         int    `json:"code"`
	Info         string `json:"info"`
	InternalCode string `json:"internal_code"`
}

type APIError struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}

func handleAPIResponse(resp *clientResponse, httpResp *http.Response) error {
	// Dump when debug is turned on
	if os.Getenv("DEBUG") != "" {
		dump, _ := httputil.DumpResponse(httpResp, true)
		log.Println("Response: ", string(dump))
	}

	defer httpResp.Body.Close()
	resp.httpResp = httpResp
	respBody, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return err
	}

	// Save response body as bytes
	resp.bytes = respBody

	// Dump out method/url, time elapsed, source ip, and transaction id for all BBVA calls
	log.Println("BBVA request: ", httpResp.Request.Method, httpResp.Request.URL.String())
	log.Println("Status code: ", httpResp.StatusCode, http.StatusText(httpResp.StatusCode))
	log.Printf("Time elapsed: %d ms", int64(resp.end.Sub(resp.start)/time.Millisecond))
	log.Println("OP-User-Id: ", httpResp.Request.Header.Get("OP-User-Id"))
	log.Println("X-Customer-IP: ", httpResp.Request.Header.Get("X-Customer-IP"))
	log.Println("X-Unique-Transaction-ID: ", httpResp.Request.Header.Get("X-Unique-Transaction-ID"))

	// Handle error cases
	if httpResp.StatusCode >= 400 {
		// Dump when not in debug
		if os.Getenv("DEBUG") == "" {
			dump, _ := httputil.DumpResponse(httpResp, true)
			log.Println("Response: ", string(dump))
		}

		if err := json.Unmarshal(resp.bytes, &resp.err); err != nil {
			return err
		}

		return resp.err.Error()
	}

	// Return bytes if no body or no content
	if resp.body == nil || httpResp.StatusCode == 204 {
		return nil
	}

	if err := json.Unmarshal(resp.bytes, &resp.body); err != nil {
		return fmt.Errorf("json.Unmarshal: %v", err)
	}

	return nil
}
