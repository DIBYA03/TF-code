/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

// Package for handling linked account

package business

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/wiseco/core-platform/api"
	bankingsrv "github.com/wiseco/core-platform/services/banking/business"
	"github.com/wiseco/core-platform/shared"
)

type LinkAccountBody = bankingsrv.LinkedBankAccountRequest

// HandleConnectAccountRequest handle connect account requests
func HandleConnectAccountRequest(r api.APIRequest) (api.APIResponse, error) {
	var method = strings.ToUpper(r.HTTPMethod)

	bID, err := shared.ParseBusinessID(r.GetPathParam("businessId"))
	if err != nil {
		return api.BadRequestError(r, err)
	}

	switch method {
	case http.MethodPost:
		return connectBankAccount(r, bID)
	default:
		return api.NotSupported(r)
	}

}

func connectBankAccount(r api.APIRequest, id shared.BusinessID) (api.APIResponse, error) {
	var requestBody LinkAccountBody
	err := json.Unmarshal([]byte(r.Body), &requestBody)
	if err != nil {
		return api.BadRequest(r, err)
	}

	b, err := bankingsrv.NewLinkedAccountService(r.SourceRequest()).ConnectBankAccount(requestBody, id)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	linkedAccountJSON, _ := json.Marshal(b)
	return api.Success(r, string(linkedAccountJSON), false)
}
