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

func getAccountBalance(r api.APIRequest, accountID string, businessID shared.BusinessID) (api.APIResponse, error) {

	account, err := accountsrv.NewBankAccountService(r.SourceRequest()).GetBalanceByID(accountID, businessID)
	if err != nil {
		return api.NotFound(r)
	}

	accountJSON, err := json.Marshal(account)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	return api.Success(r, string(accountJSON), false)

}

func HandleAccountBalanceAPIRequest(request api.APIRequest) (api.APIResponse, error) {
	method := strings.ToUpper(request.HTTPMethod)
	businessID, err := shared.ParseBusinessID(request.GetPathParam("businessId"))
	if err != nil {
		return api.BadRequestError(request, err)
	}

	accountID := request.GetPathParam("accountId")
	if accountID != "" {
		switch method {
		case http.MethodGet:
			return getAccountBalance(request, accountID, businessID)
		}
	}

	return api.NotSupported(request)
}
