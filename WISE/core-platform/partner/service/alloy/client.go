/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

package alloy

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

	retryhttp "github.com/hashicorp/go-retryablehttp"
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
}

func newHTTPClient() *retryhttp.Client {
	return retryhttp.NewClient()
}

var apiSubdomain = map[string]string{
	"sandbox": "sandbox",
	"prod":    "api",
}

func newClient() *client {
	return &client{alloyConfig, oAuthSrv}
}

func getBaseAPIURL() string {
	// Alloy App Env
	appEnv := os.Getenv("ALLOY_APP_ENV")
	if len(appEnv) == 0 {
		panic(errors.New("ALLOY_APP_ENV is missing"))
	}

	// Default to sandbox if none available or incorrect
	subDomain, ok := apiSubdomain[appEnv]
	if !ok {
		panic(errors.New("invalid app environment"))
	}

	return "https://" + subDomain + ".alloy.co"
}

var alloyConfig *config = func() *config {
	// Alloy App Name
	appName := os.Getenv("ALLOY_APP_NAME")
	if len(appName) == 0 {
		panic(errors.New("ALLOY_APP_NAME is missing"))
	}

	// Alloy App Id
	appID := os.Getenv("ALLOY_APP_ID")
	if len(appID) == 0 {
		panic(errors.New("ALLOY_APP_ID is missing"))
	}

	// Alloy OAuth Secret
	appSecret := os.Getenv("ALLOY_APP_SECRET")
	if len(appSecret) == 0 {
		panic(errors.New("ALLOY_APP_SECRET is missing"))
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

func (s *client) newAuthenticatedRequest(method, u string, buf *bytes.Buffer) (*retryhttp.Request, error) {
	token, err := s.oAuth.Token()
	if err != nil {
		return nil, err
	}

	req, err := s.newRequest(method, u, buf)
	if err != nil {
		return nil, err
	}

	// Set custom headers
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))

	return req, nil
}

func (s *client) do(req *retryhttp.Request, response interface{}) error {
	// Disallow for production
	if os.Getenv("ALLOY_APP_ENV") != "prod" {
		dump, _ := httputil.DumpRequest(toHttpRequest(req), true)
		log.Println("Request: ", string(dump))
	}

	start := time.Now()
	httpResp, err := s.httpClient.Do(req)
	log.Printf("Time elapsed (%s %s): %d ms", req.Method, req.URL, int64(time.Now().Sub(start)/time.Millisecond))
	if err != nil {
		return fmt.Errorf("http client: %v", err)
	}

	return handleAPIResponse(&clientResponse{body: response}, httpResp)
}

func (s *client) doClientResp(req *retryhttp.Request, resp *clientResponse) error {
	// Disallow for production
	if os.Getenv("ALLOY_APP_ENV") != "prod" {
		dump, _ := httputil.DumpRequest(toHttpRequest(req), true)
		log.Println("Request: ", string(dump))
	}

	start := time.Now()
	httpResp, err := s.httpClient.Do(req)
	log.Printf("Time elapsed (%s %s): %s ms", req.Method, req.URL, int64(time.Now().Sub(start)/time.Millisecond))
	if err != nil {
		return fmt.Errorf("http client: %v", err)
	}

	return handleAPIResponse(resp, httpResp)
}

func (s *client) get(path string) (*retryhttp.Request, error) {
	u := fmt.Sprintf("%s/%s", s.baseURL, path)
	req, err := s.newAuthenticatedRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, fmt.Errorf("newAuthenticatedRequest: %v", err)
	}
	return req, nil
}

func (s *client) post(path string, payload interface{}) (*retryhttp.Request, error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(payload); err != nil {
		return nil, err
	}

	// Disallow for production
	if os.Getenv("ALLOY_APP_ENV") != "prod" {
		log.Println("Post: ", string(buf.Bytes()))
	}

	u := fmt.Sprintf("%s/%s", s.baseURL, path)
	req, err := s.newAuthenticatedRequest(http.MethodPost, u, &buf)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (s *client) patch(path string, payload interface{}) (*retryhttp.Request, error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(payload); err != nil {
		return nil, err
	}

	// Disallow for production
	if os.Getenv("ALLOY_APP_ENV") != "prod" {
		log.Println("Patch: ", string(buf.Bytes()))
	}

	u := fmt.Sprintf("%s/%s", s.baseURL, path)
	req, err := s.newAuthenticatedRequest(http.MethodPatch, u, &buf)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (s *client) delete(path string) (*retryhttp.Request, error) {
	u := fmt.Sprintf("%s/%s", s.baseURL, path)
	req, err := s.newAuthenticatedRequest(http.MethodDelete, u, nil)
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

var Debug bool

func handleAPIResponse(resp *clientResponse, httpResp *http.Response) error {
	// Disallow for production
	if os.Getenv("ALLOY_APP_ENV") != "prod" {
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

	// Handle error cases
	if httpResp.StatusCode >= 400 {
		if os.Getenv("ALLOY_APP_ENV") == "prod" && len(respBody) > 0 {
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
