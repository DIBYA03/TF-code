package business

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/wiseco/core-platform/api"
	"github.com/wiseco/core-platform/services/csp/intercom"
)

// HandleIntercomAPIRequests handles the api request
func HandleIntercomAPIRequests(request api.APIRequest) (api.APIResponse, error) {
	var method = strings.ToUpper(request.HTTPMethod)

	emailID := request.GetQueryParam("emailId")
	userID := request.GetQueryParam("userId")

	if emailID == "" && userID == "" {
		return api.BadRequest(request, errors.New("emailId or userId is required"))
	}

	params := make(map[string]string)
	params["emailId"] = emailID
	params["userId"] = userID

	switch method {
	case http.MethodGet:
		return getUser(request, params)
	default:
		return api.NotSupported(request)
	}
}

func getUser(r api.APIRequest, params map[string]string) (api.APIResponse, error) {
	emailID := params["emailId"]
	userID := params["userId"]

	if emailID != "" {
		user, err := intercom.New(r.SourceRequest()).GetByEmailID(emailID)
		if err != nil {
			return api.InternalServerError(r, err)
		}
		return api.Success(r, *user, false)

	} else if userID != "" {
		user, err := intercom.New(r.SourceRequest()).GetByUserID(userID)
		if err != nil {
			return api.InternalServerError(r, err)
		}
		return api.Success(r, *user, false)

	}
	return api.BadRequestError(r, errors.New("userId/emailId query param missing"))

}

// HandleIntercomTagAPIRequests handles the tag api request
func HandleIntercomTagAPIRequests(request api.APIRequest) (api.APIResponse, error) {
	var method = strings.ToUpper(request.HTTPMethod)

	switch method {
	case http.MethodGet:
		return getTags(request)
	case http.MethodPost:
		return setTags(request)

	default:
		return api.NotSupported(request)
	}
}

func getTags(r api.APIRequest) (api.APIResponse, error) {
	tl, err := intercom.New(r.SourceRequest()).GetTags()
	if err != nil {
		return api.InternalServerError(r, err)
	}
	resp, err := json.Marshal(tl)
	return api.Success(r, string(resp), false)
}

func setTags(r api.APIRequest) (api.APIResponse, error) {
	var tagBodyItems []intercom.TagBodyItem
	err := json.Unmarshal([]byte(r.Body), &tagBodyItems)

	tl, err := intercom.New(r.SourceRequest()).SetTags(tagBodyItems)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	resp, err := json.Marshal(tl)
	return api.Success(r, string(resp), false)
}
