package hellosign

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	netURL "net/url"
	"os"
	"strings"

	"github.com/jmoiron/sqlx/types"
	"github.com/wiseco/go-lib/api"
	"github.com/wiseco/go-lib/log"
	"github.com/wiseco/go-lib/url"
)

type SignatureRequest struct {
	TemplateName string
	CustomField  []CustomField
	Email        string
	Name         string
}

type SignatureResponse struct {
	SignatureRequest SignatureRequestDetail `json:"signature_request"`
}

type SignatureRequestDetail struct {
	EmailAddress          types.JSONText `json:"cc_email_addresses"`
	CustomFields          types.JSONText `json:"custom_fields"`
	DetailsURL            string         `json:"details_url"`
	HasError              bool           `json:"has_error"`
	IsComplete            bool           `json:"is_complete"`
	Message               string         `json:"message"`
	RequesterEmailAddress string         `json:"requester_email_address"`
	ResponseData          types.JSONText `json:"response_data"`
	SignatureRequestID    string         `json:"signature_request_id"`
	Signatures            []Signature    `json:"signatures"`
	SigningURL            string         `json:"signing_url"`
	SigningRedirectURL    string         `json:"signing_redirect_url"`
	Subject               string         `json:"subject"`
	Title                 string         `json:"title"`
}

type Signature struct {
	SignatureID string `json:"signature_id"`
}

type HellosignRequest struct {
	ClientID   string `json:"client_id"`
	TemplateID string `json:"template_id"`
}

type hellosignService struct {
	log log.Logger
}

type CustomField struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type SignURLResponse struct {
	Embedded Embedded `json:"embedded"`
}

type Embedded struct {
	SignURL   string `json:"sign_url"`
	ExpiresAt int64  `json:"expires_at"`
}

type FileResponse struct {
	FileURL   string `json:"file_url"`
	ExpiresAt int64  `json:"expires_at"`
}

type HellosignService interface {
	CreateSignatureRequest(SignatureRequest) (*SignatureResponse, error)
	CancelSignatureRequest(string) error
	GetEmbeddedSignURL(string) (*SignURLResponse, error)
	DownloadDocument(string) (*string, error)
}

func NewHellosignService(l log.Logger) HellosignService {
	return &hellosignService{log: l}
}

func (s *hellosignService) CreateSignatureRequest(req SignatureRequest) (*SignatureResponse, error) {

	// Hellosign API Key
	apiKey := os.Getenv("HELLOSIGN_API_KEY")
	if len(apiKey) == 0 {
		return nil, errors.New("HELLOSIGN_API_KEY is missing")
	}

	// URL
	u, err := url.BuildAbsoluteForHellosign("", url.Params{})
	if err != nil {
		return nil, err
	}

	client := api.NewClient(
		api.ClientConfig{
			BaseURL:     u,
			ContentType: "application/x-www-form-urlencoded",
			BasicAuth:   &api.BasicAuthConfig{ClientID: apiKey},
		},
		s.log)

	val, err := req.URLValues()
	if err != nil {
		s.log.Error(err.Error())
		return nil, err
	}

	// Send request
	sr := &SignatureResponse{}
	err = client.Post("v3/signature_request/create_embedded_with_template", strings.NewReader(val.Encode()), sr)
	if err != nil {
		s.log.Error(err.Error())
		return nil, err
	}

	if len(sr.SignatureRequest.Signatures) == 0 {
		return nil, errors.New("There are no signatures")
	}

	r, err := s.GetEmbeddedSignURL(sr.SignatureRequest.Signatures[0].SignatureID)
	if err != nil {
		s.log.Error(err.Error())
		return nil, err
	}

	url := r.Embedded.SignURL
	sr.SignatureRequest.SigningURL = url

	// Return response
	return sr, nil
}

func (r SignatureRequest) URLValues() (netURL.Values, error) {
	clientID := os.Getenv("HELLOSIGN_CLIENT_ID")
	if len(clientID) == 0 {
		return nil, errors.New("HELLOSIGN_CLIENT_ID is missing")
	}

	templateID := os.Getenv("HELLOSIGN_TEMPLATE_ID")
	if len(templateID) == 0 {
		return nil, errors.New("HELLOSIGN_TEMPLATE_ID is missing")
	}

	v := netURL.Values{}
	v.Add("client_id", clientID)
	v.Add("template_id", templateID)

	v.Add(fmt.Sprintf("signers[%s][name]", r.TemplateName), r.Name)
	v.Add(fmt.Sprintf("signers[%s][email_address]", r.TemplateName), r.Email)

	if r.CustomField != nil && len(r.CustomField) > 0 {
		buf, err := json.Marshal(r.CustomField)
		if err != nil {
			return nil, err
		}

		v.Add("custom_fields", string(buf))
	}

	if url.IsDevEnv() {
		v.Add("test_mode", "1")
	}

	return v, nil
}

func (s *hellosignService) CancelSignatureRequest(signatureRequestID string) error {
	// Hellosign API Key
	apiKey := os.Getenv("HELLOSIGN_API_KEY")
	if len(apiKey) == 0 {
		return errors.New("HELLOSIGN_API_KEY is missing")
	}

	// URL
	u, err := url.BuildAbsoluteForHellosign("", url.Params{})
	if err != nil {
		return err
	}

	client := api.NewClient(
		api.ClientConfig{
			BaseURL:     u,
			ContentType: "application/x-www-form-urlencoded",
			BasicAuth:   &api.BasicAuthConfig{ClientID: apiKey},
		},
		s.log)

	v := netURL.Values{}
	if url.IsDevEnv() {
		v.Add("test_mode", "1")
	}

	// Send request
	var resp interface{}
	err = client.Post(fmt.Sprintf("v3/signature_request/cancel/%s", signatureRequestID), strings.NewReader(v.Encode()), &resp)
	if err != nil {
		s.log.Error(err.Error())
		return err
	}

	return nil
}

func (s *hellosignService) GetEmbeddedSignURL(signatureID string) (*SignURLResponse, error) {
	// Hellosign API Key
	apiKey := os.Getenv("HELLOSIGN_API_KEY")
	if len(apiKey) == 0 {
		return nil, errors.New("HELLOSIGN_API_KEY is missing")
	}

	// URL
	u, err := url.BuildAbsoluteForHellosign("", url.Params{})
	if err != nil {
		return nil, err
	}

	client := api.NewClient(
		api.ClientConfig{
			BaseURL:   u,
			BasicAuth: &api.BasicAuthConfig{ClientID: apiKey},
		},
		s.log)

	// Send request
	sr := &SignURLResponse{}

	params := url.Params{}
	if url.IsDevEnv() {
		params["test_mode"] = "1"
	}

	err = client.Get(fmt.Sprintf("v3/embedded/sign_url/%s", signatureID), params, sr)
	if err != nil {
		s.log.Error(err.Error())
		return nil, err
	}

	return sr, nil
}

func (s *hellosignService) DownloadDocument(signatureRequestID string) (*string, error) {
	// Hellosign API Key
	apiKey := os.Getenv("HELLOSIGN_API_KEY")
	if len(apiKey) == 0 {
		return nil, errors.New("HELLOSIGN_API_KEY is missing")
	}

	// URL
	u, err := url.BuildAbsoluteForHellosign("", url.Params{})
	if err != nil {
		return nil, err
	}

	client := api.NewClient(
		api.ClientConfig{
			BaseURL:   u,
			BasicAuth: &api.BasicAuthConfig{ClientID: apiKey},
		},
		s.log)

	// Send request
	fr := &FileResponse{}

	params := url.Params{}
	if url.IsDevEnv() {
		params["test_mode"] = "1"
	}

	// Url to download
	params["get_url"] = "1"

	err = client.Get(fmt.Sprintf("v3/signature_request/files/%s", signatureRequestID), params, fr)
	if err != nil {
		s.log.Error(err.Error())
		return nil, err
	}

	return getFileContent(fr.FileURL)
}

func getFileContent(fileURL string) (*string, error) {
	resp, err := http.Get(fileURL)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	body := base64.StdEncoding.EncodeToString(respBody)

	return &body, nil
}
