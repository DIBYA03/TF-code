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
	account "github.com/wiseco/core-platform/services/banking/business"
	contact "github.com/wiseco/core-platform/services/banking/business/contact"
	"github.com/wiseco/core-platform/shared"
)

type AccountPostBody = account.ContactLinkedAccountCreate

func getAccount(r api.APIRequest, id string, contactId string, businessID shared.BusinessID) (api.APIResponse, error) {
	b, err := contact.NewLinkedAccountService(r.SourceRequest()).GetById(id, contactId, businessID)
	if err != nil {
		return api.BadRequest(r, err)
	}

	contactJSON, err := json.Marshal(b)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	return api.Success(r, string(contactJSON), false)
}

func getAccounts(r api.APIRequest, contactId string, businessID shared.BusinessID) (api.APIResponse, error) {
	offset, _ := r.GetQueryIntParamWithDefault("offset", 0)
	limit, _ := r.GetQueryIntParamWithDefault("limit", 10)

	list, err := contact.NewLinkedAccountService(r.SourceRequest()).List(offset, limit, contactId, businessID)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	contactListJSON, err := json.Marshal(list)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	return api.Success(r, string(contactListJSON), false)
}

func createAccount(r api.APIRequest, contactId string, businessID shared.BusinessID) (api.APIResponse, error) {
	var requestBody AccountPostBody
	err := json.Unmarshal([]byte(r.Body), &requestBody)

	if err != nil {
		return api.BadRequest(r, err)
	}

	requestBody.BusinessID = businessID
	requestBody.ContactId = contactId
	requestBody.UserID = r.UserID

	id, err := contact.NewLinkedAccountService(r.SourceRequest()).Create(&requestBody)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	contactJSON, err := json.Marshal(id)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	return api.Success(r, string(contactJSON), false)
}

func deactivateAccount(r api.APIRequest, ID string, contactID string, businessID shared.BusinessID) (api.APIResponse, error) {

	err := contact.NewLinkedAccountService(r.SourceRequest()).Deactivate(ID, contactID, businessID)
	if err != nil {
		return api.BadRequest(r, err)
	}

	return api.Success(r, "", false)
}

//HandleLinkedAccountAPIRequests handles the api request
func HandleLinkedAccountAPIRequests(request api.APIRequest) (api.APIResponse, error) {
	var method = strings.ToUpper(request.HTTPMethod)

	contactId := request.GetPathParam("contactId")
	if contactId == "" {
		return api.NotFoundError(request, errors.New("missing contact id"))
	}

	businessID, err := shared.ParseBusinessID(request.GetPathParam("businessId"))
	if err != nil {
		return api.BadRequestError(request, err)
	}

	accountId := request.GetPathParam("accountId")
	if accountId != "" {
		switch method {
		case http.MethodGet:
			return getAccount(request, accountId, contactId, businessID)
		case http.MethodDelete:
			return deactivateAccount(request, accountId, contactId, businessID)
		default:
			return api.NotSupported(request)
		}

	}

	switch method {
	case http.MethodGet:
		return getAccounts(request, contactId, businessID)
	case http.MethodPost:
		return createAccount(request, contactId, businessID)
	default:
		return api.NotSupported(request)
	}

}
