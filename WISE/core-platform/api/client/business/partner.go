package business

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/wiseco/core-platform/api"
	"github.com/wiseco/core-platform/services/business"
	"github.com/wiseco/core-platform/shared"
)

func activateBusinessPartner(r api.APIRequest) (api.APIResponse, error) {
	var requestBody business.Partner

	err := json.Unmarshal([]byte(r.Body), &requestBody)
	if err != nil {
		return api.BadRequest(r, err)
	}

	businessID := r.GetPathParam("businessId")
	if len(businessID) > 0 {
		bID, err := shared.ParseBusinessID(businessID)
		if err != nil {
			return api.BadRequestError(r, err)
		}
		requestBody.BusinessID = bID
	}

	err = business.NewPartnerService(r.SourceRequest()).ActivatePartnerBusiness(&requestBody)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	return api.Success(r, "", false)
}

// HandleBusinessPartnerAPIRequests handles the api request
func HandleBusinessPartnerAPIRequests(request api.APIRequest) (api.APIResponse, error) {
	var method = strings.ToUpper(request.HTTPMethod)

	switch method {
	case http.MethodPost:
		return activateBusinessPartner(request)
	default:
		return api.NotSupported(request)
	}
}
