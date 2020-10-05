package delegate

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/jmoiron/sqlx"
	partnerbank "github.com/wiseco/core-platform/partner/bank"
	_ "github.com/wiseco/core-platform/partner/bank/bbva" //Should be fix later
	"github.com/wiseco/core-platform/services"
	"github.com/wiseco/core-platform/services/data"
	"github.com/wiseco/core-platform/services/user"
	"github.com/wiseco/core-platform/shared"
)

type Proxy interface {
	Execute(Resource) ProxyResponse
}

func getBaseAPIURL(r partnerbank.APIRequest) string {
	bank, err := partnerbank.GetProxyBank(partnerbank.ProviderNameBBVA)
	if err != nil {
		panic(errors.New("unable to get bank"))
	}

	return bank.ProxyService(r).GetBaseAPIURL()
}

var netClient = &http.Client{}

type header = map[string]string
type proxy struct {
	*sqlx.DB
}

// Resource of the incoming request
type Resource struct {
	Path          string      `json:"path"`
	HTTPMethod    string      `json:"method"`
	NC            bool        `json:"nc"`
	Content       interface{} `json:"content"`
	SourceRequest services.SourceRequest
}

func NewProxyService() Proxy {

	return &proxy{data.DBWrite}
}

type ProxyResponse struct {
	Body       []byte
	StatusCode int
	Error      error
}

func (p *proxy) Execute(r Resource) ProxyResponse {
	businessID := p.getBusinessID(r.SourceRequest.UserID)

	req, err := prepare(r, businessID)
	if err != nil {
		return errorResponse(err, 0)
	}

	resp, err := netClient.Do(req)
	if err != nil {
		return errorResponse(err, resp.StatusCode)
	}

	b, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	if err != nil {
		return errorResponse(err, 0)
	}

	return sucessResponse(nil, resp.StatusCode, b)
}

func getOP(r Resource, businessID shared.BusinessID, useNC bool) (*string, partnerbank.ProxyServiceProvider, error) {
	bank, err := partnerbank.GetProxyBank(partnerbank.ProviderNameBBVA)
	if err != nil {
		return nil, nil, errors.New("unable to get bank")
	}

	if useNC == true {
		nc, err := bank.ProxyService(r.SourceRequest.PartnerBankRequest()).GetBusinessBankID(partnerbank.BusinessID(businessID))
		if err != nil {
			return nil, bank, err
		}

		id := string(*nc)
		return &id, bank, nil
	}

	usr, err := user.NewUserService(r.SourceRequest).GetById(r.SourceRequest.UserID)
	if err != nil {
		return nil, bank, err
	}

	co, err := bank.ProxyService(r.SourceRequest.PartnerBankRequest()).GetConsumerBankID(partnerbank.ConsumerID(usr.ConsumerID))
	if err != nil {
		return nil, bank, err
	}

	id := string(*co)
	return &id, bank, nil
}

func prepare(r Resource, businessID shared.BusinessID) (*http.Request, error) {

	op, bank, err := getOP(r, businessID, r.NC)
	if err != nil {
		log.Println("Get OP error: ", err)
		return nil, err
	}

	jwt, err := bank.ProxyService(r.SourceRequest.PartnerBankRequest()).GetAccessToken()

	if err != nil {
		log.Println("unable to generate jwt token", err)
		return nil, err
	}

	url := getBaseAPIURL(r.SourceRequest.PartnerBankRequest()) + "/" + r.Path
	req := request(r, url)
	log.Println("Delegate url request: ", url)

	req.Header = buildHeaders(r.SourceRequest.RequestID, r.SourceRequest.SourceIP, string(*op), *jwt)

	return req, err
}

func request(r Resource, url string) *http.Request {
	var rq *http.Request

	switch strings.ToUpper(r.HTTPMethod) {
	case http.MethodPost:
		rq, _ = http.NewRequest("POST", url, nil)
		b, _ := json.Marshal(r.Content)
		log.Println(string(b))
		rq.Body = ioutil.NopCloser(bytes.NewReader(b))
	case http.MethodPatch:
		rq, _ = http.NewRequest("PATCH", url, nil)
		b, _ := json.Marshal(r.Content)
		rq.Body = ioutil.NopCloser(bytes.NewReader(b))
		rq.Method = http.MethodPatch
	default:
		rq, _ = http.NewRequest("GET", url, nil)
	}

	return rq
}

func buildHeaders(reqID, IP, OP, jwt string) http.Header {

	h := http.Header{}
	h.Set("X-Unique-Transaction-Id", reqID)
	h.Set("Content-Type", "application/json")
	h.Set("X-Customer-IP", IP)
	h.Set("OP-User-Id", OP)
	h.Set("Authorization", "jwt "+jwt)

	return h
}

func (p *proxy) getBusinessID(userID shared.UserID) shared.BusinessID {
	var id shared.BusinessID
	err := p.Get(&id, "SELECT id FROM business WHERE owner_id = $1", userID)
	if err != nil {
		log.Println("error getting business id", err)
	}

	return id
}

//errorResponse is the proxy Error Response
func errorResponse(err error, status int) ProxyResponse {
	return ProxyResponse{
		Error:      err,
		StatusCode: status,
	}
}

//sucessResponse the proxy success response
func sucessResponse(err error, status int, body []byte) ProxyResponse {
	return ProxyResponse{
		Error:      err,
		StatusCode: status,
		Body:       body,
	}
}
