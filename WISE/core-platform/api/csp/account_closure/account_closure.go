package accountclosure

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/wiseco/core-platform/api"
	accountclosure "github.com/wiseco/core-platform/services/csp/account_closure"
	csp "github.com/wiseco/core-platform/services/csp/services"
	"github.com/wiseco/core-platform/shared"
)

// HandleCSPAccountClosureAPIRequests ...
func HandleCSPAccountClosureAPIRequests(r api.APIRequest) (api.APIResponse, error) {
	method := strings.ToUpper(r.HTTPMethod)
	requestID := r.GetPathParam("requestId")

	if requestID == "" {
		if method == http.MethodPost {
			return createClosureRequest(r)
		} else if method == http.MethodGet {
			return getClosureRequestList(r)
		} else {
			return api.NotSupported(r)
		}
	} else {
		if method == http.MethodGet {
			return getClosureRequestItem(r)
		} else if method == http.MethodPatch {
			return updateClosureRequest(r)
		} else {
			return api.NotSupported(r)
		}
	}

}

// HandleCSPAccountClosureStateAPIRequests ...
func HandleCSPAccountClosureStateAPIRequests(r api.APIRequest) (api.APIResponse, error) {
	method := strings.ToUpper(r.HTTPMethod)
	requestID := r.GetPathParam("requestId")

	if requestID != "" {
		if method == http.MethodGet {
			return getClosureStateList(r)
		} else {
			return api.NotSupported(r)
		}
	} else {
		return api.NotSupported(r)
	}

}

func createClosureRequest(r api.APIRequest) (api.APIResponse, error) {
	var create accountclosure.CSPAccountClosureCreate
	err := json.Unmarshal([]byte(r.Body), &create)

	_, err = shared.ParseBusinessID(string(create.BusinessID))
	if err != nil {
		return api.BadRequestError(r, err)
	}

	list, err := accountclosure.NewCSPService(csp.NewSRRequest(r.CognitoID)).CSPClosureRequestCreate(create)
	if err != nil {
		return api.InternalServerError(r, err)
	}
	jsonList, _ := json.Marshal(list)
	return api.Success(r, string(jsonList), false)
}

func getClosureRequestList(r api.APIRequest) (api.APIResponse, error) {
	params := accountclosure.CSPAccountClosureQueryParams{}

	params.Status = r.GetQueryParam("status")
	params.StartDate = r.GetQueryParam("startDate")
	params.EndDate = r.GetQueryParam("endDate")

	params.BusinessID = r.GetQueryParam("businessId")
	params.BusinessName = r.GetQueryParam("businessName")
	params.OwnerName = r.GetQueryParam("ownerName")
	params.AvailableBalanceMin = r.GetQueryParam("availableBalanceMin")
	params.AvailableBalanceMax = r.GetQueryParam("availableBalanceMax")
	params.PostedBalanceMin = r.GetQueryParam("postedBalanceMin")
	params.PostedBalanceMax = r.GetQueryParam("postedBalanceMax")

	params.Offset = r.GetQueryParam("offset")
	params.Limit = r.GetQueryParam("limit")

	params.SortField = r.GetQueryParam("sortField")
	params.SortDirection = r.GetQueryParam("sortDirection")

	list, err := accountclosure.NewCSPService(csp.NewSRRequest(r.CognitoID)).CSPClosureRequestList(params)
	if err != nil {
		return api.InternalServerError(r, err)
	}
	jsonList, _ := json.Marshal(list)
	return api.Success(r, string(jsonList), false)
}

func getClosureRequestItem(r api.APIRequest) (api.APIResponse, error) {
	requestID := r.GetPathParam("requestId")

	requestItem, err := accountclosure.NewCSPService(csp.NewSRRequest(r.CognitoID)).CSPClosureRequestDetails(requestID)
	if err != nil {
		return api.InternalServerError(r, err)
	}
	reqJSON, _ := json.Marshal(requestItem)
	return api.Success(r, string(reqJSON), false)
}

func updateClosureRequest(r api.APIRequest) (api.APIResponse, error) {
	var body accountclosure.CSPClosureRequestPatch
	requestID := r.GetPathParam("requestId")

	if err := json.Unmarshal([]byte(r.Body), &body); err != nil {
		return api.BadRequest(r, err)
	}

	if body.Status != "" {
		status := body.Status
		isValid := accountclosure.CheckAccountClosureStatus(body.Status)
		if status != "" && !isValid {
			return api.BadRequest(r, errors.New("Invalid status"))
		}

		requestItem, err := accountclosure.NewCSPService(csp.NewSRRequest(r.CognitoID)).CSPClosureRequestUpdate(requestID, status)
		if err != nil {
			return api.InternalServerError(r, err)
		}

		requestJSON, _ := json.Marshal(requestItem)
		return api.Success(r, string(requestJSON), false)
	}

	return api.BadRequest(r, errors.New("Invalid patch body"))
}

func getClosureStateList(r api.APIRequest) (api.APIResponse, error) {
	requestID := r.GetPathParam("requestId")

	list, err := accountclosure.NewCSPService(csp.NewSRRequest(r.CognitoID)).CSPClosureStateList(requestID)
	if err != nil {
		return api.InternalServerError(r, err)
	}
	jsonList, _ := json.Marshal(list)
	return api.Success(r, string(jsonList), false)
}
