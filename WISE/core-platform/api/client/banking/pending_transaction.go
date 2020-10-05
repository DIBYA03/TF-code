package banking

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/wiseco/core-platform/api"
	"github.com/wiseco/core-platform/shared"
	"github.com/wiseco/core-platform/transaction"
)

func transactionList(r api.APIRequest, businessID shared.BusinessID) (api.APIResponse, error) {
	params := make(map[string]interface{})
	params["offset"], _ = r.GetQueryIntParamWithDefault("offset", 0)
	params["limit"], _ = r.GetQueryIntParamWithDefault("limit", 20)
	params["startDate"] = r.GetQueryParam("startDate")
	params["endDate"] = r.GetQueryParam("endDate")
	params["type"] = r.GetQueryParam("type")
	params["minAmount"] = r.GetQueryParam("minAmount")
	params["maxAmount"] = r.GetQueryParam("maxAmount")
	params["text"] = r.GetQueryParam("text")
	params["contactId"] = r.GetQueryParam("contactId")

	list, err := transaction.NewPendingTransactionService().ListAll(params, r.UserID, businessID)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	listJSON, _ := json.Marshal(list)

	return api.Success(r, string(listJSON), false)
}

func transactionByID(r api.APIRequest, txnID shared.PendingTransactionID, businessID shared.BusinessID) (api.APIResponse, error) {
	transaction, err := transaction.NewPendingTransactionService().GetByID(txnID, r.UserID, businessID)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	transactionJSON, _ := json.Marshal(transaction)

	return api.Success(r, string(transactionJSON), false)
}

//HandlePendingTransactionAPIRequest handles the api request for pending transactions
func HandlePendingTransactionAPIRequest(r api.APIRequest) (api.APIResponse, error) {
	if r.BusinessID == nil {
		return api.InternalServerError(r, errors.New("Missing header X-Wise-Business-ID"))
	}

	transactionID := r.GetPathParam("transactionId")

	method := strings.ToUpper(r.HTTPMethod)
	switch method {
	case http.MethodGet:
		if transactionID != "" {
			transactionID, err := shared.ParsePendingTransactionID(transactionID)
			if err != nil {
				return api.BadRequestError(r, err)
			}

			return transactionByID(r, transactionID, *r.BusinessID)
		} else {
			return transactionList(r, *r.BusinessID)
		}
	default:
		return api.NotSupported(r)
	}
}

func exportTransaction(r api.APIRequest, id shared.BusinessID) (api.APIResponse, error) {
	startDate := r.GetQueryParam("startDate")
	endDate := r.GetQueryParam("endDate")
	if len(startDate) == 0 || len(endDate) == 0 {
		return api.InternalServerError(r, errors.New("Missing start and end date query params"))
	}

	t, err := transaction.NewPendingTransactionService().Export(r.UserID, id, startDate, endDate)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	csvJSON, _ := json.Marshal(t)

	return api.Success(r, string(csvJSON), false)
}

//HandlePendingTransactionExportRequest handles the api request for transactions csv export
func HandlePendingTransactionExportRequest(r api.APIRequest) (api.APIResponse, error) {
	method := strings.ToUpper(r.HTTPMethod)
	businessID, err := shared.ParseBusinessID(r.GetPathParam("businessId"))
	if err != nil {
		return api.BadRequestError(r, err)
	}

	switch method {
	case http.MethodGet:
		return exportTransaction(r, businessID)
	default:
		return api.NotSupported(r)
	}
}
