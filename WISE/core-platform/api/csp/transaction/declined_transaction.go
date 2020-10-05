package transaction

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/wiseco/core-platform/api"
	"github.com/wiseco/core-platform/services/csp/transaction"
	"github.com/wiseco/core-platform/shared"
)

func declinedTransactionList(r api.APIRequest) (api.APIResponse, error) {
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
	params["subtype"] = r.GetQueryParam("subtype")

	businessID := r.GetQueryParam("businessId")
	if len(businessID) > 0 {
		bID, err := shared.ParseBusinessID(businessID)
		if err != nil {
			return api.BadRequestError(r, err)
		}
		params["businessId"] = bID
	}

	list, err := transaction.NewDeclinedTransactionService(r.SourceRequest()).ListAllDeclinedTransaction(params)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	listJSON, _ := json.Marshal(list)

	return api.Success(r, string(listJSON), false)
}

func declinedTransactionByID(r api.APIRequest, txnID shared.PendingTransactionID) (api.APIResponse, error) {
	businessID := r.GetQueryParam("businessId")
	if len(businessID) == 0 {
		return api.BadRequestError(r, errors.New("businessId query param missing"))
	}

	bID, err := shared.ParseBusinessID(businessID)
	if err != nil {
		return api.BadRequestError(r, err)
	}

	transaction, err := transaction.NewDeclinedTransactionService(r.SourceRequest()).GetDeclinedTransactionByID(txnID, bID)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	transactionJSON, _ := json.Marshal(transaction)

	return api.Success(r, string(transactionJSON), false)
}

//HandleTransactionRequest handles the api request for transactions
func HandleDeclinedTransactionRequest(r api.APIRequest) (api.APIResponse, error) {
	method := strings.ToUpper(r.HTTPMethod)

	transactionID := r.GetPathParam("transactionId")
	if len(transactionID) > 0 {
		txnID, err := shared.ParsePendingTransactionID(transactionID)
		if err != nil {
			return api.BadRequestError(r, err)
		}

		switch method {
		case http.MethodGet:
			return declinedTransactionByID(r, txnID)
		default:
			return api.NotSupported(r)
		}
	}

	switch method {
	case http.MethodGet:
		return declinedTransactionList(r)
	}

	return api.NotSupported(r)
}

func exportDeclinedTransaction(r api.APIRequest) (api.APIResponse, error) {
	startDate := r.GetQueryParam("startDate")
	endDate := r.GetQueryParam("endDate")
	if len(startDate) == 0 || len(endDate) == 0 {
		return api.InternalServerError(r, errors.New("Missing start and end date query params"))
	}

	var busID *shared.BusinessID
	businessID := r.GetQueryParam("businessId")
	if len(businessID) > 0 {
		bID, err := shared.ParseBusinessID(businessID)
		if err != nil {
			return api.BadRequestError(r, err)
		}

		busID = &bID
	}

	t, err := transaction.NewDeclinedTransactionService(r.SourceRequest()).ExportDeclinedTransaction(busID, startDate, endDate)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	csvJSON, _ := json.Marshal(t)

	return api.Success(r, string(csvJSON), false)
}

//HandleDeclinedTransactionExportRequest handles the api request for transactions csv export
func HandleDeclinedTransactionExportRequest(r api.APIRequest) (api.APIResponse, error) {
	method := strings.ToUpper(r.HTTPMethod)

	switch method {
	case http.MethodGet:
		return exportDeclinedTransaction(r)
	default:
		return api.NotSupported(r)
	}
}
