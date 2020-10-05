package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	netURL "net/url"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	retryhttp "github.com/hashicorp/go-retryablehttp"
	"github.com/wiseco/go-lib/log"
)

const (
	//TypeBasic oauth basic type
	TypeBasic = oauthType("basic")
	//TypeAuthCode oauth auth code type
	TypeAuthCode = oauthType("auth_code")
)

type oauthType string

//OAUTHToken our normalized view of oauth token responses
type OAUTHToken struct {
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

//OAUTHClientConfig config passed in to a new client
type OAUTHClientConfig struct {
	URL               string
	ClientID          string
	ClientSecret      string
	GrantType         string
	Code              string
	RedirectURI       string
	Type              oauthType
	HeaderValuePrefix *string
}

type authCodeBody struct {
	Code         string `json:"code"`
	RedirectURI  string `json:"redirect_uri"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	GrantType    string `json:"grant_type"`
}

//OAUTHClient describes an oauth client
type OAUTHClient interface {
	Token() (*OAUTHToken, error)
}

type oAuthClient struct {
	config     OAUTHClientConfig
	oAuthToken *OAUTHToken
	mux        *sync.Mutex
	log        log.Logger
}

//NewOAUTHClient returns a Client interface
func NewOAUTHClient(c OAUTHClientConfig, l log.Logger) OAUTHClient {
	return &oAuthClient{
		config: c,
		mux:    &sync.Mutex{},
		log:    l,
	}
}

//Token attempt to request a new oauth token
func (o *oAuthClient) Token() (*OAUTHToken, error) {
	var err error
	o.mux.Lock()
	defer o.mux.Unlock()
	if o.oAuthToken == nil || o.oAuthToken.Expired() {
		o.oAuthToken, err = o.getToken()
	}

	return o.oAuthToken, err
}

//Expired returns whether or not a token has expired, gives 60 second window
func (t *OAUTHToken) Expired() bool {
	// Check 60 seconds before expiry
	if t.ExpiresAt.Unix() <= (time.Now().Unix() - 60) {
		return true
	}
	return false
}

func (o *oAuthClient) getToken() (*OAUTHToken, error) {
	var resp OAUTHToken

	rBody, contentType, err := o.getRequestBodyForType()
	if err != nil {
		return &resp, err
	}

	req, err := retryhttp.NewRequest(http.MethodPost, o.config.URL, rBody)
	if err != nil {
		return &resp, err
	}

	if contentType != "" {
		req.Header.Add("Content-Type", contentType)
	}

	o.setHeadersForType(req)

	rcl := retryhttp.NewClient()
	rcl.Logger = retryhttp.LeveledLogger(o.log)

	cl := client{
		httpClient: rcl,
		log:        o.log,
	}

	err = cl.do(req, &resp, ResponseFormatterJSON)
	if err != nil {
		return &resp, err
	}

	resp.ExpiresAt = time.Now().Add(time.Duration(resp.ExpiresIn) * time.Second)

	return &resp, nil
}

func (o *oAuthClient) setHeadersForType(req *retryablehttp.Request) {
	switch o.config.Type {
	case TypeBasic:
		req.SetBasicAuth(o.config.ClientID, o.config.ClientSecret)
	case TypeAuthCode:
		// handle auth code when required
	}
}

func (o *oAuthClient) getRequestBodyForType() (io.Reader, string, error) {
	switch o.config.Type {
	case TypeBasic:
		v := netURL.Values{}
		v.Add("grant_type", o.config.GrantType)
		return strings.NewReader(v.Encode()), "application/x-www-form-urlencoded", nil
	case TypeAuthCode:
		buf := &bytes.Buffer{}
		body := authCodeBody{
			Code:         o.config.Code,
			RedirectURI:  o.config.RedirectURI,
			ClientID:     o.config.ClientID,
			ClientSecret: o.config.ClientSecret,
			GrantType:    o.config.GrantType,
		}

		err := json.NewEncoder(buf).Encode(body)
		return buf, "application/json", err
	default:
		return strings.NewReader(""), "", nil
	}
}
