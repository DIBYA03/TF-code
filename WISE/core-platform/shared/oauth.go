package shared

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
	"sync"
	"time"

	retryhttp "github.com/hashicorp/go-retryablehttp"
)

type OAuthConfig struct {
	BaseURL   string
	AppID     string
	AppSecret string
}

type clientResponse struct {
	httpResp *http.Response
	body     interface{}
	bytes    []byte
}

type oAuthService struct {
	config     OAuthConfig
	httpClient *retryhttp.Client
	oAuthToken *OAuthTokenResponse
	mux        *sync.Mutex
}

type OAuthService interface {
	Token() (*OAuthTokenResponse, error)
}

func NewOAuthService(c OAuthConfig) OAuthService {
	return &oAuthService{
		config:     c,
		httpClient: retryhttp.NewClient(),
		mux:        &sync.Mutex{},
	}
}

type OAuthTokenResponse struct {
	// AccessToken
	// Access token to call API services.
	// STRING | REQUIRED
	AccessToken string `json:"access_token"`

	// TokenType
	// Security token type.
	// STRING | REQUIRED
	TokenType string `json:"token_type"`

	// ExpiresIn
	// TTL for this access token in seconds.
	// NUMBER | REQUIRED
	ExpiresIn int64 `json:"expires_in"`

	// Scope
	// STRING | REQUIRED
	Scope string `json:"scope"`

	ExpiresAt time.Time `json:"-"`
}

func (t *OAuthTokenResponse) Expired() bool {
	// Check 60 seconds before expiry
	if t.ExpiresAt.Unix() <= (time.Now().Unix() - 60) {
		return true
	}
	return false
}

func (s *oAuthService) Token() (*OAuthTokenResponse, error) {
	var err error
	s.mux.Lock()
	if s.oAuthToken == nil || s.oAuthToken.Expired() {
		s.oAuthToken, err = s.token()
	}

	s.mux.Unlock()
	return s.oAuthToken, err
}

func (s *oAuthService) token() (*OAuthTokenResponse, error) {
	var body = map[string]interface{}{}
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(body); err != nil {
		return nil, err
	}

	url := s.config.BaseURL + "?grant_type=client_credentials"
	req, err := retryhttp.NewRequest("POST", url, &buf)
	if err != nil {
		return nil, err
	}

	// add some headers
	req.SetBasicAuth(s.config.AppID, s.config.AppSecret)
	start := time.Now()
	httpResp, err := s.httpClient.Do(req)
	log.Printf("Time elapsed (%s %s): %d ms", req.Method, req.URL, int64(time.Now().Sub(start)/time.Millisecond))
	if err != nil {
		return nil, err
	}

	var response OAuthTokenResponse
	if err := handleAPIResponse(&clientResponse{body: &response}, httpResp); err != nil {
		return nil, err
	}

	response.ExpiresAt = time.Now().Add(time.Duration(response.ExpiresIn) * time.Second)
	return &response, nil
}

func handleAPIResponse(resp *clientResponse, httpResp *http.Response) error {
	// Disallow for production
	if os.Getenv("BBVA_APP_ENV") != "prod" {
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
		if os.Getenv("BBVA_APP_ENV") == "prod" && len(respBody) > 0 {
			dump, _ := httputil.DumpResponse(httpResp, true)
			log.Println("Response: ", string(dump))
		}

		return errors.New("http error")
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
