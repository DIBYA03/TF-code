/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

// Package for handling register account

package business

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/wiseco/core-platform/api"
	bankingsrv "github.com/wiseco/core-platform/services/banking/business"
	"github.com/wiseco/core-platform/shared"
)

type LinkAccountCreateBody = bankingsrv.LinkedExternalAccountCreate

// HandleLinkAccountRequest handle register account requests
func HandleLinkAccountRequest(r api.APIRequest) (api.APIResponse, error) {
	var method = strings.ToUpper(r.HTTPMethod)

	bID, err := shared.ParseBusinessID(r.GetPathParam("businessId"))
	if err != nil {
		return api.BadRequestError(r, err)
	}

	aID := r.GetPathParam("accountId")
	if aID != "" {
		switch method {
		case http.MethodGet:
			return getLinkedAccount(r, aID, bID)
		case http.MethodDelete:
			return unlinkBankAccount(r, aID, bID)
		default:
			return api.NotSupported(r)
		}
	}

	switch method {
	case http.MethodGet:
		return getLinkedAccounts(r, bID)
	case http.MethodPost:
		return linkBankAccount(r, r.UserID, bID)
	default:
		return api.NotSupported(r)
	}
}

func linkBankAccount(r api.APIRequest, uID shared.UserID, bID shared.BusinessID) (api.APIResponse, error) {
	var requestBody LinkAccountCreateBody
	err := json.Unmarshal([]byte(r.Body), &requestBody)
	if err != nil {
		return api.BadRequest(r, err)
	}

	requestBody.BusinessID = bID
	requestBody.UserID = uID

	b, err := bankingsrv.NewLinkedAccountService(r.SourceRequest()).LinkExternalBankAccount(&requestBody)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	linkedAccountJSON, _ := json.Marshal(b)
	return api.Success(r, string(linkedAccountJSON), false)
}

func unlinkBankAccount(r api.APIRequest, id string, businessId shared.BusinessID) (api.APIResponse, error) {
	b, err := bankingsrv.NewLinkedAccountService(r.SourceRequest()).UnlinkBankAccount(id, businessId)
	if err != nil {
		return api.BadRequest(r, err)
	}

	Accountjson, err := json.Marshal(b)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	return api.Success(r, string(Accountjson), false)
}

func getLinkedAccount(r api.APIRequest, id string, businessId shared.BusinessID) (api.APIResponse, error) {
	b, err := bankingsrv.NewLinkedAccountService(r.SourceRequest()).GetById(id, businessId)
	if err != nil {
		return api.BadRequest(r, err)
	}

	Accountjson, err := json.Marshal(b)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	return api.Success(r, string(Accountjson), false)
}

func getLinkedAccounts(r api.APIRequest, businessId shared.BusinessID) (api.APIResponse, error) {
	offset, _ := r.GetQueryIntParamWithDefault("offset", 0)
	limit, _ := r.GetQueryIntParamWithDefault("limit", 10)

	b, err := bankingsrv.NewLinkedAccountService(r.SourceRequest()).List(offset, limit, businessId)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	AccountListjson, err := json.Marshal(b)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	return api.Success(r, string(AccountListjson), false)
}
