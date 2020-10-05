package kyc

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/wiseco/core-platform/api"
	"github.com/wiseco/core-platform/services/csp/business"
)

type KYCPost struct {
	SSN string `json:"ssn"`
} 

// HandleKYCAPIRequests handles the api request
func HandleKYCAPIRequests(request api.APIRequest) (api.APIResponse, error) {
	var method = strings.ToUpper(request.HTTPMethod)

	switch method {
	case http.MethodPost:
		return runKYC(request)
	default:
		return api.NotSupported(request)
	}
}

func runKYC(r api.APIRequest) (api.APIResponse, error) {
	var body KYCPost
	if err := json.Unmarshal([]byte(r.Body), &body); err != nil {
		return api.BadRequest(r, err)
	}

	kycResponse, err := business.NewMemberService(r.SourceRequest()).RunKYC(body.SSN)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	resp, err := json.Marshal(kycResponse)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	return api.Success(r, string(resp), false)
}
