/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

// Package for handling contact's linked accounts
package contact

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/wiseco/core-platform/api"
	transfer "github.com/wiseco/core-platform/services/banking/business"
	contact "github.com/wiseco/core-platform/services/banking/business/contact"
	"github.com/wiseco/core-platform/shared"
)

type TransferPostBody = transfer.TransferInitiate

func getTransfer(r api.APIRequest, id string, contactId string, businessID shared.BusinessID) (api.APIResponse, error) {
	s := r.SourceRequest()

	b, err := contact.NewMoneyTransferService(s).GetById(id, contactId, businessID)
	if err != nil {
		return api.BadRequest(r, err)
	}

	transferJSON, err := json.Marshal(b)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	return api.Success(r, string(transferJSON), false)
}

func getContactTransfers(r api.APIRequest, contactId string, businessID shared.BusinessID) (api.APIResponse, error) {
	s := r.SourceRequest()

	offset, _ := r.GetQueryIntParamWithDefault("offset", 0)
	limit, _ := r.GetQueryIntParamWithDefault("limit", 10)

	list, err := contact.NewMoneyTransferService(s).GetByContactId(offset, limit, contactId, businessID)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	transferListJSON, err := json.Marshal(list)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	return api.Success(r, string(transferListJSON), false)
}

func initiateTransfer(r api.APIRequest, contactId string, businessID shared.BusinessID) (api.APIResponse, error) {
	s := r.SourceRequest()

	var requestBody TransferPostBody

	err := json.Unmarshal([]byte(r.Body), &requestBody)
	if err != nil {
		return api.BadRequest(r, err)
	}

	requestBody.BusinessID = businessID
	requestBody.ContactId = &contactId
	requestBody.CreatedUserID = r.UserID

	id, err := contact.NewMoneyTransferService(s).Transfer(&requestBody)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	transferJSON, _ := json.Marshal(id)
	return api.Success(r, string(transferJSON), false)
}

func Cancel(r api.APIRequest, id string, contactId string, businessID shared.BusinessID) (api.APIResponse, error) {
	s := r.SourceRequest()

	b, err := contact.NewMoneyTransferService(s).Cancel(id, contactId, businessID)
	if err != nil {
		return api.BadRequest(r, err)
	}

	transferJSON, err := json.Marshal(b)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	return api.Success(r, string(transferJSON), false)
}

//HandleLinkedCardAPIRequests handles the api request
func HandleMoneyTransferAPIRequests(request api.APIRequest) (api.APIResponse, error) {
	var method = strings.ToUpper(request.HTTPMethod)

	contactId := request.GetPathParam("contactId")
	if contactId == "" {
		return api.NotFoundError(request, errors.New("missing contact id"))
	}

	businessID, err := shared.ParseBusinessID(request.GetPathParam("businessId"))
	if err != nil {
		return api.BadRequestError(request, err)
	}

	transferId := request.GetPathParam("transferId")
	if transferId != "" {
		switch method {
		case http.MethodGet:
			return getTransfer(request, transferId, contactId, businessID)
		case http.MethodDelete:
			return Cancel(request, transferId, contactId, businessID)
		default:
			return api.NotSupported(request)
		}

	}

	switch method {
	case http.MethodGet:
		return getContactTransfers(request, contactId, businessID)
	case http.MethodPost:
		return initiateTransfer(request, contactId, businessID)
	default:
		return api.NotSupported(request)
	}

}
