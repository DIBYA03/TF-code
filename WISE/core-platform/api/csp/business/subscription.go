package business

import (
	"encoding/json"

	"github.com/wiseco/core-platform/api"
	"github.com/wiseco/core-platform/services/csp/business"
	csp "github.com/wiseco/core-platform/services/csp/services"
	"github.com/wiseco/core-platform/shared"
)

func handleSubscriptionByBusinessID(businessID string, r api.APIRequest) (api.APIResponse, error) {
	bID, err := shared.ParseBusinessID(businessID)
	if err != nil {
		return api.BadRequestError(r, err)
	}

	sub, err := business.NewSubscriptionService(csp.NewSRRequest(r.CognitoID)).GetByBusinessID(bID)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	subscriptionJSON, err := json.Marshal(sub)
	if err != nil {
		return api.InternalServerError(r, err)
	}
	return api.Success(r, string(subscriptionJSON), false)
}

func handleSubscriptionUpdate(businessID string, r api.APIRequest) (api.APIResponse, error) {
	var requestBody business.SubscriptionUpdate
	bID, err := shared.ParseBusinessID(businessID)
	if err != nil {
		return api.BadRequestError(r, err)
	}

	err = json.Unmarshal([]byte(r.Body), &requestBody)
	if err != nil {
		return api.BadRequest(r, err)
	}

	requestBody.BusinessID = bID

	sub, err := business.NewSubscriptionService(csp.NewSRRequest(r.CognitoID)).Update(requestBody)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	subscriptionJSON, err := json.Marshal(sub)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	return api.Success(r, string(subscriptionJSON), false)
}
