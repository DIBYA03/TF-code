package transaction

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/wiseco/core-platform/api"
	"github.com/wiseco/core-platform/services/csp/business"
)

//HandleTransactionApproveRequest handles the api request for transactions
func HandleTransactionApproveRequest(r api.APIRequest) (api.APIResponse, error) {
	method := strings.ToUpper(r.HTTPMethod)

	switch method {
	case http.MethodPost:
		return approveTransaction(r)
	}

	return api.NotSupported(r)
}

func approveTransaction(r api.APIRequest) (api.APIResponse, error) {
	transactionID := r.GetPathParam("transactionId")
	if len(transactionID) == 0 {
		return api.BadRequestError(r, errors.New("transactionId query param missing"))
	}

	err := business.NewBanking(r.SourceRequest()).ApproveTransferInReview(transactionID, r.CognitoID, r.SourceIP)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	return api.Success(r, "", false)
}

//HandleTransactionDeclineRequest handles the api request for transactions
func HandleTransactionDeclineRequest(r api.APIRequest) (api.APIResponse, error) {
	method := strings.ToUpper(r.HTTPMethod)

	switch method {
	case http.MethodPost:
		return declineTransaction(r)
	}

	return api.NotSupported(r)
}

func declineTransaction(r api.APIRequest) (api.APIResponse, error) {
	transactionID := r.GetPathParam("transactionId")
	if len(transactionID) == 0 {
		return api.BadRequestError(r, errors.New("transactionId query param missing"))
	}

	err := business.NewBanking(r.SourceRequest()).DeclineTransferInReview(transactionID, r.CognitoID, r.SourceIP)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	return api.Success(r, "", false)
}

//HandleTransactionDeclineRequest handles the api request for transactions
func HandleTransactionTransferInfo(r api.APIRequest) (api.APIResponse, error) {
	method := strings.ToUpper(r.HTTPMethod)

	switch method {
	case http.MethodGet:
		transactionID := r.GetPathParam("transactionId")

		trp, err := business.NewBanking(r.SourceRequest()).GetTransferForTransactionID(transactionID)
		if err != nil {
			return api.InternalServerError(r, err)
		}

		resp, err := json.Marshal(trp)
		if err != nil {
			return api.InternalServerError(r, err)
		}

		return api.Success(r, string(resp), false)
	}

	return api.NotSupported(r)
}
