/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

package transaction

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/wiseco/core-platform/api"
	"github.com/wiseco/core-platform/shared"
	"github.com/wiseco/core-platform/transaction"
)

type ReceiptPostBody = transaction.ReceiptCreate

func deleteReceiptByID(r api.APIRequest, ID string, txnID shared.PostedTransactionID, businessID shared.BusinessID) (api.APIResponse, error) {
	err := transaction.NewReceiptService(r.SourceRequest()).Delete(ID, txnID, businessID)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	return api.Success(r, "", false)
}

func receiptByID(r api.APIRequest, id string, txnID shared.PostedTransactionID, businessID shared.BusinessID) (api.APIResponse, error) {
	transaction, err := transaction.NewReceiptService(r.SourceRequest()).GetById(id, txnID, businessID)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	transactionJSON, _ := json.Marshal(transaction)

	return api.Success(r, string(transactionJSON), false)
}

func createReceipt(r api.APIRequest, txnID shared.PostedTransactionID, businessID shared.BusinessID) (api.APIResponse, error) {
	var requestBody ReceiptPostBody
	err := json.Unmarshal([]byte(r.Body), &requestBody)

	if err != nil {
		return api.BadRequest(r, err)
	}

	requestBody.CreatedUserID = r.UserID
	requestBody.BusinessID = businessID
	requestBody.TransactionID = txnID

	receipt, err := transaction.NewReceiptService(r.SourceRequest()).Create(&requestBody)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	rJSON, _ := json.Marshal(receipt)
	return api.Success(r, string(rJSON), false)
}

//HandleTransactionReceiptRequest handles the api request for transactions receipts
func HandleTransactionReceiptRequest(r api.APIRequest) (api.APIResponse, error) {
	method := strings.ToUpper(r.HTTPMethod)
	businessID, err := shared.ParseBusinessID(r.GetPathParam("businessId"))
	if err != nil {
		return api.BadRequestError(r, err)
	}

	transactionID, err := shared.ParsePostedTransactionID(r.GetPathParam("transactionId"))
	if err != nil {
		return api.BadRequestError(r, err)
	}

	receiptID := r.GetPathParam("receiptId")
	if receiptID != "" {
		switch method {
		case http.MethodGet:
			return receiptByID(r, receiptID, transactionID, businessID)
		case http.MethodDelete:
			return deleteReceiptByID(r, receiptID, transactionID, businessID)
		default:
			return api.NotSupported(r)
		}
	}

	switch method {
	case http.MethodPost:
		return createReceipt(r, transactionID, businessID)
	default:
		return api.NotSupported(r)
	}
}

func receiptUrlByID(r api.APIRequest, id string, txnID shared.PostedTransactionID, businessID shared.BusinessID) (api.APIResponse, error) {
	url, err := transaction.NewReceiptService(r.SourceRequest()).GetSignedURL(id, txnID, businessID)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	if err != nil && err.Error() == "not found" {
		return api.NotFound(r)
	} else if err != nil {
		return api.InternalServerError(r, err)
	}

	if url == nil {
		return api.InternalServerError(r, errors.New("Unable to get signed url"))
	}

	URL := struct {
		URL *string `json:"url"`
	}{URL: url}

	b, _ := json.Marshal(URL)
	return api.Success(r, string(b), false)
}

//HandleTxnReceiptUrlRequest handles the api request for transactions receipt url
func HandleTxnReceiptUrlRequest(r api.APIRequest) (api.APIResponse, error) {
	method := strings.ToUpper(r.HTTPMethod)
	businessID, err := shared.ParseBusinessID(r.GetPathParam("businessId"))
	if err != nil {
		return api.BadRequestError(r, err)
	}

	transactionID, err := shared.ParsePostedTransactionID(r.GetPathParam("transactionId"))
	if err != nil {
		return api.BadRequestError(r, err)
	}

	receiptID := r.GetPathParam("receiptId")
	if receiptID != "" {
		switch method {
		case http.MethodGet:
			return receiptUrlByID(r, receiptID, transactionID, businessID)
		default:
			return api.NotSupported(r)
		}
	}

	return api.NotSupported(r)
}
