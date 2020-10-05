package transaction

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/wiseco/core-platform/api"
	"github.com/wiseco/core-platform/services/banking/business"
	"github.com/wiseco/core-platform/shared"
	"github.com/wiseco/core-platform/transaction"
)

type TransactionPatchBody = transaction.BusinessPostedTransactionUpdate

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

	accounts, err := business.NewBankAccountService(r.SourceRequest()).GetByUsageType(r.UserID, businessID, business.UsageTypePrimary)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	// Use first primary account for now - API will be deprecated - List must be by account
	list, err := transaction.NewBusinessService().ListAll(params, r.UserID, businessID, accounts[0].Id)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	listJSON, _ := json.Marshal(list)

	return api.Success(r, string(listJSON), false)
}

func transactionByID(r api.APIRequest, txnID shared.PostedTransactionID, businessID shared.BusinessID) (api.APIResponse, error) {
	transaction, err := transaction.NewBusinessService().GetByID(txnID, r.UserID, businessID)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	transactionJSON, _ := json.Marshal(transaction)

	return api.Success(r, string(transactionJSON), false)
}

func transactionUpdate(r api.APIRequest, txnID shared.PostedTransactionID, businessID shared.BusinessID) (api.APIResponse, error) {
	var requestBody TransactionPatchBody
	err := json.Unmarshal([]byte(r.Body), &requestBody)
	if err != nil {
		return api.BadRequest(r, err)
	}

	requestBody.TransactionID = txnID
	requestBody.BusinessID = businessID
	transaction, err := transaction.NewBusinessService().Update(requestBody, r.UserID)

	if err != nil {
		return api.InternalServerError(r, err)
	}

	transactionJSON, _ := json.Marshal(transaction)
	return api.Success(r, string(transactionJSON), false)
}

//HandleTransactionRequest handles the api request for transactions
func HandleTransactionRequest(r api.APIRequest) (api.APIResponse, error) {
	method := strings.ToUpper(r.HTTPMethod)
	businessID, err := shared.ParseBusinessID(r.GetPathParam("businessId"))
	if err != nil {
		return api.BadRequestError(r, err)
	}

	transactionID := r.GetPathParam("transactionId")
	if transactionID != "" {
		txnID, err := shared.ParsePostedTransactionID(transactionID)
		if err != nil {
			return api.BadRequestError(r, err)
		}

		switch method {
		case http.MethodGet:
			return transactionByID(r, txnID, businessID)
		case http.MethodPatch:
			return transactionUpdate(r, txnID, businessID)
		default:
			return api.NotSupported(r)
		}
	}

	if businessID != "" {
		switch method {
		case http.MethodGet:
			return transactionList(r, businessID)
		}
	}

	return api.NotSupported(r)
}

func exportTransaction(r api.APIRequest, id shared.BusinessID) (api.APIResponse, error) {
	startDate := r.GetQueryParam("startDate")
	endDate := r.GetQueryParam("endDate")
	if len(startDate) == 0 || len(endDate) == 0 {
		return api.InternalServerError(r, errors.New("Missing start and end date query params"))
	}

	t, err := transaction.NewBusinessService().Export(r.UserID, id, startDate, endDate)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	csvJSON, _ := json.Marshal(t)

	return api.Success(r, string(csvJSON), false)
}

//HandleTransactionExportRequest handles the api request for transactions csv export
func HandleTransactionExportRequest(r api.APIRequest) (api.APIResponse, error) {
	method := strings.ToUpper(r.HTTPMethod)
	businessID, err := shared.ParseBusinessID(r.GetPathParam("businessId"))
	if err != nil {
		return api.BadRequestError(r, err)
	}

	if businessID != "" {
		switch method {
		case http.MethodGet:
			return exportTransaction(r, businessID)
		}
	}

	return api.NotSupported(r)
}
