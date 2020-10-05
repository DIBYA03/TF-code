/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

// Package for business banking apis
package business

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/wiseco/core-platform/api"
	"github.com/wiseco/core-platform/services/banking/business"
	cardsrv "github.com/wiseco/core-platform/services/banking/business"
	"github.com/wiseco/core-platform/shared"
)

type CardPostBody = cardsrv.BankCardCreate
type CardPatchBody = cardsrv.BankCardUpdate

func getCard(r api.APIRequest, id string, businessID shared.BusinessID, userID shared.UserID) (api.APIResponse, error) {
	s := r.SourceRequest()

	b, err := business.NewCardService(s).GetById(id, businessID, userID)
	if err != nil {

		return api.BadRequest(r, err)
	}

	contactJSON, err := json.Marshal(b)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	return api.Success(r, string(contactJSON), false)
}

func getBusinessCards(r api.APIRequest, businessID shared.BusinessID, userID shared.UserID) (api.APIResponse, error) {
	offset, _ := r.GetQueryIntParamWithDefault("offset", 0)
	limit, _ := r.GetQueryIntParamWithDefault("limit", 10)

	s := r.SourceRequest()

	list, err := business.NewCardService(s).GetByBusinessID(offset, limit, businessID, userID)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	contactListJSON, err := json.Marshal(list)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	return api.Success(r, string(contactListJSON), false)
}

func getAccountCards(r api.APIRequest, accountId string, businessID shared.BusinessID, userID shared.UserID) (api.APIResponse, error) {
	offset, _ := r.GetQueryIntParamWithDefault("offset", 0)
	limit, _ := r.GetQueryIntParamWithDefault("limit", 10)

	s := r.SourceRequest()

	list, err := business.NewCardService(s).GetByAccountId(offset, limit, accountId, businessID, userID)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	contactListJSON, err := json.Marshal(list)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	return api.Success(r, string(contactListJSON), false)
}

func createCard(r api.APIRequest, accountId string, businessID shared.BusinessID) (api.APIResponse, error) {
	var requestBody CardPostBody
	err := json.Unmarshal([]byte(r.Body), &requestBody)

	if err != nil {
		return api.BadRequest(r, err)
	}

	requestBody.BusinessID = businessID
	requestBody.CardholderID = r.UserID
	requestBody.BankAccountId = accountId

	s := r.SourceRequest()

	id, err := business.NewCardService(s).CreateBankCard(&requestBody)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	contactJSON, err := json.Marshal(id)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	return api.Success(r, string(contactJSON), false)
}

func HandleBusinessBankCardAPIRequests(request api.APIRequest) (api.APIResponse, error) {
	var method = strings.ToUpper(request.HTTPMethod)

	businessID, err := shared.ParseBusinessID(request.GetPathParam("businessId"))
	if err != nil {
		return api.BadRequest(request, err)
	}

	accountId := request.GetPathParam("accountId")
	if accountId != "" {
		switch method {
		case http.MethodGet:
			return getAccountCards(request, accountId, businessID, request.UserID)
		case http.MethodPost:
			return createCard(request, accountId, businessID)
		default:
			return api.NotSupported(request)
		}

	}

	cardId := request.GetPathParam("cardId")
	if cardId != "" {
		switch method {
		case http.MethodGet:
			return getCard(request, cardId, businessID, request.UserID)
		default:
			return api.NotSupported(request)
		}

	}

	switch method {
	case http.MethodGet:
		return getBusinessCards(request, businessID, request.UserID)
	default:
		return api.NotSupported(request)
	}

}
