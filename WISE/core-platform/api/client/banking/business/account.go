/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

// Package for business  bankingbanking apis
package business

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/wiseco/core-platform/api"
	accountsrv "github.com/wiseco/core-platform/services/banking/business"
	"github.com/wiseco/core-platform/shared"
)

type accountPostBody = accountsrv.BankAccountCreate
type accountPatchBody = accountsrv.BankAccountUpdate

func getAccount(r api.APIRequest, accountId string, businessID shared.BusinessID) (api.APIResponse, error) {
	account, err := accountsrv.NewBankAccountService(r.SourceRequest()).GetByID(accountId, businessID)
	if err != nil {
		return api.NotFound(r)
	}

	accountJSON, err := json.Marshal(account)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	return api.Success(r, string(accountJSON), false)
}

func getAccounts(r api.APIRequest, businessID shared.BusinessID) (api.APIResponse, error) {
	offset, _ := r.GetQueryIntParamWithDefault("offset", 0)
	limit, _ := r.GetQueryIntParamWithDefault("limit", 10)

	accounts, err := accountsrv.NewBankAccountService(r.SourceRequest()).List(businessID, limit, offset)
	if err != nil {
		return api.NotFound(r)
	}

	accountsJSON, err := json.Marshal(accounts)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	return api.Success(r, string(accountsJSON), false)
}

func updateAccount(r api.APIRequest, accountId string, businessID shared.BusinessID) (api.APIResponse, error) {

	var requestBody accountPatchBody
	if err := json.Unmarshal([]byte(r.Body), &requestBody); err != nil {
		return api.BadRequest(r, err)
	}

	requestBody.Id = accountId
	requestBody.BusinessID = businessID
	account, err := accountsrv.NewBankAccountService(r.SourceRequest()).Update(&requestBody)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	accountJSON, err := json.Marshal(account)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	return api.Success(r, string(accountJSON), false)
}

func HandleAccountAPIRequest(request api.APIRequest) (api.APIResponse, error) {
	method := strings.ToUpper(request.HTTPMethod)
	businessID, err := shared.ParseBusinessID(request.GetPathParam("businessId"))
	if err != nil {
		return api.BadRequest(request, err)
	}

	accountId := request.GetPathParam("accountId")
	//return getAccounts(request, businessID)

	if accountId != "" && businessID != "" {
		switch method {
		case http.MethodGet:
			return getAccount(request, accountId, businessID)
		case http.MethodPatch:
			return updateAccount(request, accountId, businessID)
		default:
			return api.NotSupported(request)
		}
	}

	if businessID != "" {
		switch method {
		case http.MethodGet:
			return getAccounts(request, businessID)
		default:
			return api.NotSupported(request)
		}
	}

	return api.NotSupported(request)
}
