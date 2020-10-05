/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

package transaction

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/wiseco/core-platform/api"
	"github.com/wiseco/core-platform/shared"
	"github.com/wiseco/core-platform/transaction"
)

type DisputePostBody = transaction.DisputeCreate
type DisputeCancelBody = transaction.DisputeCancel

func disputeByID(r api.APIRequest, ID string, txnID shared.PostedTransactionID, businessID shared.BusinessID) (api.APIResponse, error) {
	transaction, err := transaction.NewDisputeService(r.SourceRequest()).GetById(ID, txnID, businessID)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	transactionJSON, _ := json.Marshal(transaction)

	return api.Success(r, string(transactionJSON), false)
}

func createDispute(r api.APIRequest, txnID shared.PostedTransactionID, businessID shared.BusinessID) (api.APIResponse, error) {
	var requestBody DisputePostBody
	err := json.Unmarshal([]byte(r.Body), &requestBody)
	if err != nil {
		return api.BadRequest(r, err)
	}

	requestBody.CreatedUserID = r.UserID
	requestBody.BusinessID = businessID
	requestBody.TransactionID = txnID

	d, err := transaction.NewDisputeService(r.SourceRequest()).Create(&requestBody)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	dJSON, _ := json.Marshal(d)
	return api.Success(r, string(dJSON), false)
}

func cancelDispute(r api.APIRequest, id string, txnID shared.PostedTransactionID, businessID shared.BusinessID) (api.APIResponse, error) {
	requestBody := DisputeCancelBody{}
	requestBody.CreatedUserID = r.UserID
	requestBody.BusinessID = businessID
	requestBody.TransactionID = txnID
	requestBody.Id = id

	d, err := transaction.NewDisputeService(r.SourceRequest()).Cancel(&requestBody)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	dJSON, _ := json.Marshal(d)
	return api.Success(r, string(dJSON), false)
}

//HandleTransactionDisputeRequest handles the api request for transactions dispute
func HandleTransactionDisputeRequest(r api.APIRequest) (api.APIResponse, error) {
	method := strings.ToUpper(r.HTTPMethod)
	businessID, err := shared.ParseBusinessID(r.GetPathParam("businessId"))
	if err != nil {
		return api.BadRequestError(r, err)
	}

	transactionID, err := shared.ParsePostedTransactionID(r.GetPathParam("transactionId"))
	if err != nil {
		return api.BadRequestError(r, err)
	}

	disputeID := r.GetPathParam("disputeId")
	if disputeID != "" {
		switch method {
		case http.MethodGet:
			return disputeByID(r, disputeID, transactionID, businessID)
		case http.MethodDelete:
			return cancelDispute(r, disputeID, transactionID, businessID)
		default:
			return api.NotSupported(r)
		}
	}

	switch method {
	case http.MethodPost:
		return createDispute(r, transactionID, businessID)
	default:
		return api.NotSupported(r)
	}
}
