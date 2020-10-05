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
	card "github.com/wiseco/core-platform/services/banking/business"
	contact "github.com/wiseco/core-platform/services/banking/business/contact"
	"github.com/wiseco/core-platform/shared"
)

type CardPostBody = card.LinkedCardCreate

func getCard(r api.APIRequest, id string, contactId string, businessID shared.BusinessID) (api.APIResponse, error) {
	b, err := contact.NewLinkedCardService(r.SourceRequest()).GetById(id, contactId, businessID)
	if err != nil {
		return api.BadRequest(r, err)
	}

	cardJSON, err := json.Marshal(b)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	return api.Success(r, string(cardJSON), false)
}

func getCards(r api.APIRequest, contactId string, businessID shared.BusinessID) (api.APIResponse, error) {
	offset, _ := r.GetQueryIntParamWithDefault("offset", 0)
	limit, _ := r.GetQueryIntParamWithDefault("limit", 10)

	list, err := contact.NewLinkedCardService(r.SourceRequest()).List(offset, limit, contactId, businessID)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	cardListJSON, err := json.Marshal(list)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	return api.Success(r, string(cardListJSON), false)
}

func createCard(r api.APIRequest, contactId string, businessID shared.BusinessID) (api.APIResponse, error) {
	var requestBody CardPostBody

	err := json.Unmarshal([]byte(r.Body), &requestBody)
	if err != nil {
		return api.BadRequest(r, err)
	}

	requestBody.BusinessID = businessID
	requestBody.ContactId = &contactId
	requestBody.UserID = r.UserID

	id, err := contact.NewLinkedCardService(r.SourceRequest()).Create(&requestBody)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	contactJSON, _ := json.Marshal(id)
	return api.Success(r, string(contactJSON), false)
}

func deactivateCard(r api.APIRequest, ID string, contactID string, businessID shared.BusinessID) (api.APIResponse, error) {

	err := contact.NewLinkedCardService(r.SourceRequest()).Deactivate(ID, contactID, businessID)
	if err != nil {
		return api.BadRequest(r, err)
	}

	return api.Success(r, "", false)
}

//HandleLinkedCardAPIRequests handles the api request
func HandleLinkedCardAPIRequests(request api.APIRequest) (api.APIResponse, error) {
	var method = strings.ToUpper(request.HTTPMethod)

	contactId := request.GetPathParam("contactId")
	if contactId == "" {
		return api.NotFoundError(request, errors.New("missing contact id"))
	}

	businessID, err := shared.ParseBusinessID(request.GetPathParam("businessId"))
	if err != nil {
		return api.BadRequestError(request, err)
	}

	cardId := request.GetPathParam("cardId")
	if cardId != "" {
		switch method {
		case http.MethodGet:
			return getCard(request, cardId, contactId, businessID)
		case http.MethodDelete:
			return deactivateCard(request, cardId, contactId, businessID)
		default:
			return api.NotSupported(request)
		}

	}

	switch method {
	case http.MethodGet:
		return getCards(request, contactId, businessID)
	case http.MethodPost:
		return createCard(request, contactId, businessID)
	default:
		return api.NotSupported(request)
	}
}
