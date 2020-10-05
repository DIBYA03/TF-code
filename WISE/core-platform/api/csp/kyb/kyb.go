package kyb

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/wiseco/core-platform/api"
	"github.com/wiseco/core-platform/services/csp/business"
)

type KYBPost struct {
	TIN string `json:"tin"`
}

// HandleKYBAPIRequests handles the api request
func HandleKYBAPIRequests(request api.APIRequest) (api.APIResponse, error) {
	var method = strings.ToUpper(request.HTTPMethod)

	switch method {
	case http.MethodPost:
		return runKYB(request)
	default:
		return api.NotSupported(request)
	}
}

func runKYB(r api.APIRequest) (api.APIResponse, error) {
	var body KYBPost
	if err := json.Unmarshal([]byte(r.Body), &body); err != nil {
		return api.BadRequest(r, err)
	}

	kybResponse, err := business.New(r.SourceRequest()).RunKYB(body.TIN)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	resp, err := json.Marshal(kybResponse)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	return api.Success(r, string(resp), false)
}
