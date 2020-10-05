/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

// Package for handling external card registeration

package business

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/wiseco/core-platform/api"
	bankingsrv "github.com/wiseco/core-platform/services/banking/business"
	"github.com/wiseco/core-platform/shared"
)

type LinkCardCreateBody = bankingsrv.LinkedCardCreate

// HandleLinkCardRequest handle register account requests
func HandleLinkCardRequest(r api.APIRequest) (api.APIResponse, error) {
	var method = strings.ToUpper(r.HTTPMethod)

	bID, err := shared.ParseBusinessID(r.GetPathParam("businessId"))
	if err != nil {
		return api.BadRequestError(r, err)
	}

	cID := r.GetPathParam("cardId")
	if cID != "" {
		switch method {
		case http.MethodGet:
			return getLinkedCard(r, cID, bID)
		case http.MethodDelete:
			return unlinkCard(r, cID, bID)
		default:
			return api.NotSupported(r)
		}
	}

	switch method {
	case http.MethodGet:
		return getLinkedCards(r, bID)
	case http.MethodPost:
		return linkCard(r, r.UserID, bID)
	default:
		return api.NotSupported(r)
	}
}

func linkCard(r api.APIRequest, uID shared.UserID, bID shared.BusinessID) (api.APIResponse, error) {
	var requestBody LinkCardCreateBody
	err := json.Unmarshal([]byte(r.Body), &requestBody)
	if err != nil {
		return api.BadRequest(r, err)
	}

	requestBody.BusinessID = bID
	requestBody.UserID = uID

	b, err := bankingsrv.NewLinkedCardService(r.SourceRequest()).Create(&requestBody)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	linkedAccountJSON, _ := json.Marshal(b)
	return api.Success(r, string(linkedAccountJSON), false)
}

func getLinkedCard(r api.APIRequest, id string, businessId shared.BusinessID) (api.APIResponse, error) {
	b, err := bankingsrv.NewLinkedCardService(r.SourceRequest()).GetByID(id, businessId)
	if err != nil {
		return api.BadRequest(r, err)
	}

	Accountjson, err := json.Marshal(b)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	return api.Success(r, string(Accountjson), false)
}

func unlinkCard(r api.APIRequest, id string, businessId shared.BusinessID) (api.APIResponse, error) {
	err := bankingsrv.NewLinkedCardService(r.SourceRequest()).Delete(id, businessId)
	if err != nil {
		return api.BadRequest(r, err)
	}

	return api.Success(r, "success", false)
}

func getLinkedCards(r api.APIRequest, businessId shared.BusinessID) (api.APIResponse, error) {
	offset, _ := r.GetQueryIntParamWithDefault("offset", 0)
	limit, _ := r.GetQueryIntParamWithDefault("limit", 10)

	b, err := bankingsrv.NewLinkedCardService(r.SourceRequest()).List(offset, limit, businessId)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	AccountListjson, err := json.Marshal(b)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	return api.Success(r, string(AccountListjson), false)
}
