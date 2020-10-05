package business

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/wiseco/core-platform/api"
	partnerbank "github.com/wiseco/core-platform/partner/bank"
	"github.com/wiseco/core-platform/services/csp/business"
	"github.com/wiseco/core-platform/shared"
)

func getCSPBusinessList(r api.APIRequest) (api.APIResponse, error) {

	params := make(map[string]interface{})
	params["offset"], _ = r.GetQueryIntParamWithDefault("offset", 0)
	params["limit"], _ = r.GetQueryIntParamWithDefault("limit", 20)

	businessID := r.GetQueryParam("businessId")
	if len(businessID) > 0 {
		bID, err := shared.ParseBusinessID(businessID)
		if err != nil {
			return api.BadRequestError(r, err)
		}
		params["businessId"] = bID
	}

	bankID := r.GetQueryParam("bankId")
	if len(bankID) > 0 {
		params["bankId"] = partnerbank.BusinessBankID(bankID)
	}

	userID := r.GetQueryParam("userId")
	if len(userID) > 0 {
		uID, err := shared.ParseUserID(userID)
		if err != nil {
			return api.BadRequestError(r, err)
		}
		params["userId"] = uID
	}

	consumerBankID := r.GetQueryParam("consumerBankId")
	if len(consumerBankID) > 0 {
		params["consumerBankId"] = partnerbank.ConsumerBankID(consumerBankID)
	}

	ownerFirstName := r.GetQueryParam("ownerFirstName")
	if len(ownerFirstName) > 0 {
		params["ownerFirstName"] = string(ownerFirstName)
	}

	ownerPhoneNumber := r.GetQueryParam("ownerPhoneNumber")
	if len(ownerPhoneNumber) > 0 {
		params["ownerPhoneNumber"] = string(ownerPhoneNumber)
	}

	ownerEmailID := r.GetQueryParam("ownerEmailId")
	if len(ownerEmailID) > 0 {
		params["ownerEmailId"] = string(ownerEmailID)
	}

	params["name"] = r.GetQueryParam("name")

	list, err := business.NewCSPService().CSPBusinessList(params)
	if err != nil {
		return api.InternalServerError(r, err)
	}
	jsonList, _ := json.Marshal(list)
	return api.Success(r, string(jsonList), false)
}

// HandleCSPBusinessAPIRequests handles the api request
func HandleCSPBusinessAPIRequests(request api.APIRequest) (api.APIResponse, error) {
	var method = strings.ToUpper(request.HTTPMethod)

	switch method {
	case http.MethodGet:
		return getCSPBusinessList(request)
	default:
		return api.NotSupported(request)
	}
}
