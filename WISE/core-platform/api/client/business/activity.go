package business

import (
	"encoding/json"

	"github.com/wiseco/core-platform/api"
	activitysrv "github.com/wiseco/core-platform/services/activitystream/business"
	"github.com/wiseco/core-platform/shared"

	"net/http"
	"strings"
)

func activityList(r api.APIRequest, id shared.BusinessID, userID shared.UserID) (api.APIResponse, error) {
	offset, _ := r.GetQueryIntParamWithDefault("offset", 0)
	limit, _ := r.GetQueryIntParamWithDefault("limit", 20)

	list, err := activitysrv.New(r.SourceRequest()).List(offset, limit, id, userID)
	if err != nil {

		return api.InternalServerError(r, err)
	}

	listJSON, _ := json.Marshal(list)

	return api.Success(r, string(listJSON), false)
}

func activityByID(r api.APIRequest, id string) (api.APIResponse, error) {
	activity, err := activitysrv.New(r.SourceRequest()).GetByID(id)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	activityJSON, _ := json.Marshal(activity)

	return api.Success(r, string(activityJSON), false)
}

func HandleActivityRequest(r api.APIRequest) (api.APIResponse, error) {
	method := strings.ToUpper(r.HTTPMethod)
	businessID, err := shared.ParseBusinessID(r.GetPathParam("businessId"))
	if err != nil {
		return api.BadRequestError(r, err)
	}

	activityID := r.GetPathParam("activityId")

	if activityID != "" {
		switch method {
		case http.MethodGet:
			return activityByID(r, activityID)
		default:
			return api.NotSupported(r)
		}
	}

	switch method {
	case http.MethodGet:
		return activityList(r, businessID, r.UserID)
	}

	return api.NotSupported(r)
}
