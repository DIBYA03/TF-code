package user

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/wiseco/core-platform/api"
	activitysrv "github.com/wiseco/core-platform/services/activitystream/user"
)

func activityList(r api.APIRequest, id string) (api.APIResponse, error) {
	list, err := activitysrv.New(r.SourceRequest()).List(id)
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
	userID := r.GetPathParam("userId")
	activityID := r.GetPathParam("activityId")

	if userID != "" {
		switch method {
		case http.MethodGet:
			return activityList(r, userID)
		}
	}

	if userID != "" && activityID != "" {
		switch method {
		case http.MethodGet:
			return activityByID(r, activityID)
		default:
			return api.NotSupported(r)
		}
	}

	return api.NotSupported(r)
}
