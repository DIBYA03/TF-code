/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

package business

import (
	"encoding/json" //"net/http" -- Not being used right now.
	"net/http"
	"strings"

	"github.com/wiseco/core-platform/api"
	subsrv "github.com/wiseco/core-platform/services/subscription"
	"github.com/wiseco/core-platform/shared"
)

type RequestBody = subsrv.SubscriptionUpdate

func subscribe(r api.APIRequest, businessID shared.BusinessID) (api.APIResponse, error) {
	var requestBody RequestBody
	if err := json.Unmarshal([]byte(r.Body), &requestBody); err != nil {
		return api.BadRequest(r, err)
	}

	requestBody.BusinessID = businessID
	requestBody.UserID = r.UserID

	b, err := subsrv.NewSubscriptionService(r.SourceRequest()).Update(requestBody)
	if err != nil {
		return api.InternalServerError(r, err)
	} else if b == nil {
		return api.NotFound(r)
	}

	userJSON, _ := json.Marshal(b)
	return api.Success(r, string(userJSON), false)
}

func HandleSubscriptionAPIRequest(r api.APIRequest) (api.APIResponse, error) {

	var method = strings.ToUpper(r.HTTPMethod)

	businessId, err := shared.ParseBusinessID(r.GetPathParam("businessId"))
	if err == nil {
		switch method {
		case http.MethodPatch:
			return subscribe(r, businessId)
		default:
			return api.NotSupported(r)
		}
	}

	return api.NotSupported(r)
}
