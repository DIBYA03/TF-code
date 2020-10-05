package business

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/wiseco/core-platform/api"
	bsrv "github.com/wiseco/core-platform/services/business"
	"github.com/wiseco/core-platform/shared"
)

type accountClosurePostBody = bsrv.AccountClosureCreate

//HandleAccountClosurAPIRequests handles the api request
func HandleAccountClosurAPIRequests(request api.APIRequest) (api.APIResponse, error) {
	var method = strings.ToUpper(request.HTTPMethod)
	businessID, err := shared.ParseBusinessID(request.GetPathParam("businessId"))
	if err != nil {
		return api.BadRequestError(request, err)
	}

	switch method {
	case http.MethodGet:
		return getAccountClosureDetails(request, businessID)
	case http.MethodPost:
		return createAccountClosure(request, businessID)
	default:
		return api.NotSupported(request)
	}
}

func getAccountClosureDetails(r api.APIRequest, businessID shared.BusinessID) (api.APIResponse, error) {
	req, err := bsrv.NewAccountClosureService(r.SourceRequest()).GetByBusinessID(businessID)
	if err != nil {
		if err == sql.ErrNoRows {
			resp := map[string]string{"status": "not_requested"}
			reqJSON, _ := json.Marshal(resp)
			return api.Success(r, string(reqJSON), false)
		}
		return api.BadRequest(r, err)
	}

	reqJSON, err := json.Marshal(req)
	return api.Success(r, string(reqJSON), false)
}

func createAccountClosure(r api.APIRequest, businessID shared.BusinessID) (api.APIResponse, error) {
	var body accountClosurePostBody
	err := json.Unmarshal([]byte(r.Body), &body)
	if err != nil {
		return api.BadRequest(r, err)
	}
	body.BusinessID = businessID
	resp, err := bsrv.NewAccountClosureService(r.SourceRequest()).Create(body)
	if err != nil {
		return api.BadRequest(r, err)
	}

	reqJSON, err := json.Marshal(resp)
	return api.Success(r, string(reqJSON), false)
}
