package api

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"os"
	"time"

	retryhttp "github.com/hashicorp/go-retryablehttp"
	"github.com/wiseco/go-lib/log"
	"github.com/wiseco/go-lib/url"
)

const (
	// RetryWaitMin is the minimum time before retrying a request
	RetryWaitMin = 10 * time.Second
	// RetryMax is the maximum number of retries
	RetryMax = 3
)

//HTTPClient an interface describing the httper
type HTTPClient interface {
	Do(*retryhttp.Request) (*http.Response, error)
}

//RawBody is an interface used to help readability and assert the correct named interfaces are passed as the correct parameters
type RawBody interface{}

//Response is an interface used to help readability and assert the correct named interfaces are passed as the correct parameters
type Response interface{}

//Client an interface describing this packages client struct
type Client interface {
	Post(string, RawBody, Response, ...ResponseFormatter) error
	Patch(string, RawBody, Response, ...ResponseFormatter) error
	Get(string, url.Params, Response, ...ResponseFormatter) error
	Delete(string, url.Params) error
}

//Type BasicAuthConfig when you want to use basic auth but not with oauth
type BasicAuthConfig struct {
	ClientID     string
	ClientSecret string
}

type BearerTokenConfig struct {
	AccessToken string
}

// TLSConfig - Struct to store certificate file path and key file path
type TLSConfig struct {
	CertFilePath string
	KeyFilePath  string
}

//ClientConfig dictates how this api package will interact with the outside world
type ClientConfig struct {
	BaseURL           string
	ContentType       string
	AcceptType        string
	Headers           Headers
	BasicAuth         *BasicAuthConfig
	OAUTHClientConfig *OAUTHClientConfig
	BearerToken       *BearerTokenConfig
	CustomHeader      http.Header
	RequestLogging    bool

	// Allow shared OAuth configuration to avoid refetching the access token
	OAuthClient OAUTHClient

	TLSClientConfig *TLSConfig
}

//Header represents a http header name and value
type Header struct {
	Name  string
	Value string
}

//Headers are a list of headers to be used in the request
type Headers []Header

type client struct {
	clientConfig ClientConfig
	oAuthClient  OAUTHClient
	httpClient   HTTPClient
	log          log.Logger
}

//TestingHTTPClient used for testing only, overrides the http client functionality
var TestingHTTPClient HTTPClient

type ClientBuilder func(ClientConfig, log.Logger) *client

var DefaultClientBuilder ClientBuilder = func(cc ClientConfig, l log.Logger) *client {
	return &client{
		clientConfig: cc,
		httpClient:   getHTTPClient(l, cc, nil),
		log:          l,
	}
}

var TLSClientBuilder ClientBuilder = func(cc ClientConfig, l log.Logger) *client {
	if cc.TLSClientConfig == nil {
		panic("tlsConfig must not be empty")
	}
	// Load client cert
	cert, err := tls.LoadX509KeyPair(cc.TLSClientConfig.CertFilePath, cc.TLSClientConfig.KeyFilePath)
	if err != nil {
		panic(err)
	}

	// Setup HTTPS client
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}
	tlsConfig.BuildNameToCertificate()
	transport := &http.Transport{TLSClientConfig: tlsConfig}
	return &client{
		clientConfig: cc,
		httpClient:   getHTTPClient(l, cc, transport),
		log:          l,
	}
}

//NewClient returns a client to interact with api's
func NewClient(cc ClientConfig, l log.Logger, cb ...ClientBuilder) Client {
	var clientBuilder ClientBuilder
	if len(cb) == 0 {
		clientBuilder = DefaultClientBuilder
	} else {
		clientBuilder = cb[0]
	}
	c := clientBuilder(cc, l)

	if cc.OAuthClient != nil {
		c.oAuthClient = cc.OAuthClient
	} else if cc.OAUTHClientConfig != nil {
		oac := NewOAUTHClient(*cc.OAUTHClientConfig, l)
		c.oAuthClient = oac
	}

	return c
}

type ResponseFormatter func([]byte, interface{}) error

var ResponseFormatterJSON ResponseFormatter = func(resp []byte, respFormatted interface{}) error {
	return json.Unmarshal(resp, &respFormatted)
}

var ResponseFormatterEmpty ResponseFormatter = func(resp []byte, respFormatted interface{}) error {
	return nil
}

var ResponseFormatterString ResponseFormatter = func(resp []byte, respFormatted interface{}) error {
	r := string(resp)
	var dptr, ok = respFormatted.(**string)
	var err error = nil
	if ok && dptr != nil {
		*dptr = &r
	} else {
		err = errors.New("Expected a non nil double pointer string for returing response")
	}
	return err
}

func (c *client) Post(path string, rawBody RawBody, resp Response, rf ...ResponseFormatter) error {
	u, err := url.BuildAbsolute(path, c.clientConfig.BaseURL, nil)
	if err != nil {
		return err
	}

	req, err := c.getRequest(http.MethodPost, u, rawBody)
	if err != nil {
		return err
	}
	respFormatter := ResponseFormatterJSON
	if len(rf) != 0 {
		respFormatter = rf[0]
	}
	return c.do(req, resp, respFormatter)
}

func (c client) Patch(path string, rawBody RawBody, resp Response, rf ...ResponseFormatter) error {
	u, err := url.BuildAbsolute(path, c.clientConfig.BaseURL, nil)
	if err != nil {
		return err
	}

	req, err := c.getRequest(http.MethodPatch, u, rawBody)
	if err != nil {
		return err
	}
	respFormatter := ResponseFormatterJSON
	if len(rf) != 0 {
		respFormatter = rf[0]
	}
	return c.do(req, resp, respFormatter)
}

func (c client) Get(path string, up url.Params, resp Response, rf ...ResponseFormatter) error {
	var buf bytes.Buffer

	u, err := url.BuildAbsolute(path, c.clientConfig.BaseURL, up)
	if err != nil {
		return err
	}

	req, err := c.getRequest(http.MethodGet, u, &buf)
	if err != nil {
		return err
	}
	respFormatter := ResponseFormatterJSON
	if len(rf) != 0 {
		respFormatter = rf[0]
	}
	return c.do(req, resp, respFormatter)
}

func (c client) Delete(path string, up url.Params) error {
	var buf bytes.Buffer

	u, err := url.BuildAbsolute(path, c.clientConfig.BaseURL, up)
	if err != nil {
		return err
	}

	req, err := c.getRequest(http.MethodDelete, u, &buf)
	if err != nil {
		return err
	}
	respFormatter := ResponseFormatterEmpty
	return c.do(req, nil, respFormatter)
}

func (c client) do(req *retryhttp.Request, resp interface{}, respFormatter ResponseFormatter) error {
	start := time.Now()
	httpResp, err := c.httpClient.Do(req)
	defer httpResp.Body.Close()

	duration := time.Now().Sub(start)

	if err != nil {
		c.log.ErrorD("ApiReq", log.Fields{"time_elapsed": duration, "method": req.Method, "url": req.URL, "error": err.Error()})

		return fmt.Errorf("http client: %v", err)
	}

	c.log.InfoD("ApiReq", log.Fields{"time_elapsed": duration, "method": req.Method, "url": req.URL})

	respBody, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return err
	}

	if httpResp.StatusCode >= 400 {
		return errors.New(string(respBody))
	}

	if respBody == nil || len(respBody) == 0 || httpResp.StatusCode == 204 {
		return nil
	}
	return respFormatter(respBody, resp)
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

func (c *client) getRequest(m, u string, rawBody interface{}) (*retryhttp.Request, error) {
	req, err := retryhttp.NewRequest(m, u, rawBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Cache-Control", "no-cache")

	// No content type for HTTP Get
	if m != http.MethodGet {
		if c.clientConfig.ContentType != "" {
			req.Header.Set("Content-Type", c.clientConfig.ContentType)
		} else {
			req.Header.Set("Content-Type", "application/json")
		}
	}
	if c.clientConfig.AcceptType != "" {
		req.Header.Set("Accept", c.clientConfig.AcceptType)
	}

	for i := range c.clientConfig.CustomHeader {
		req.Header.Set(i, c.clientConfig.CustomHeader[i][0])
	}

	if c.oAuthClient != nil {
		t, err := c.oAuthClient.Token()
		if err != nil {
			return nil, err
		}

		headerValuePrefix := "Bearer"
		if c.clientConfig.OAUTHClientConfig != nil && c.clientConfig.OAUTHClientConfig.HeaderValuePrefix != nil {
			headerValuePrefix = *c.clientConfig.OAUTHClientConfig.HeaderValuePrefix
		}

		req.Header.Set("Authorization", fmt.Sprintf("%s %s", headerValuePrefix, t.AccessToken))
	}

	if c.clientConfig.BasicAuth != nil {
		req.SetBasicAuth(c.clientConfig.BasicAuth.ClientID, c.clientConfig.BasicAuth.ClientSecret)
	}

	if c.clientConfig.BearerToken != nil {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.clientConfig.BearerToken.AccessToken))
	}

	if len(c.clientConfig.Headers) > 0 {
		for _, h := range c.clientConfig.Headers {
			req.Header.Set(h.Name, h.Value)
		}
	}
	if os.Getenv("API_ENV") != "prod" {
		dump, _ := httputil.DumpRequest(toHttpRequest(req), true)
		fmt.Println(string(dump))
	}

	return req, nil
}

func getHTTPClient(l log.Logger, cc ClientConfig, transport *http.Transport) HTTPClient {
	if TestingHTTPClient != nil {
		return TestingHTTPClient
	}

	c := retryhttp.NewClient()
	if transport != nil {
		c.HTTPClient.Transport = transport
	}
	c.RetryWaitMin = RetryWaitMin
	c.RetryMax = RetryMax
	c.ErrorHandler = retryhttp.PassthroughErrorHandler

	c.Logger = retryhttp.LeveledLogger(l)

	if cc.RequestLogging {
		c.RequestLogHook = requestLogHook
	}

	return c
}

func requestLogHook(logger retryhttp.Logger, req *http.Request, numRetries int) {
	bb, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		logger.Printf("REQUEST LOGGING ERROR READING BYTES: %s", err.Error())
		return
	}

	logger.Printf("REQUEST LOGGING: %s", string(bb))

	return
}
