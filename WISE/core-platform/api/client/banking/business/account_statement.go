/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

// Package for business  bankingbanking apis
package business

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/wiseco/core-platform/api"
	accountsrv "github.com/wiseco/core-platform/services/banking/business"
	"github.com/wiseco/core-platform/shared"
)

func getStatementByID(r api.APIRequest, stmtID string, accountID string, businessID shared.BusinessID) (api.APIResponse, error) {
	s, err := accountsrv.NewAccountStmtService(r.SourceRequest()).GetByID(stmtID, accountID, r.UserID, businessID)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	// PLATFORM-281: Pass PDF as base64 encoded binary file to API Gateway
	resp, err := api.Success(r, base64.StdEncoding.EncodeToString(s.Content), true)
	resp.Headers[http.CanonicalHeaderKey("Content-Type")] = s.ContentType
	return resp, err
}

func getStatements(r api.APIRequest, accountID string, businessID shared.BusinessID) (api.APIResponse, error) {
	statements, err := accountsrv.NewAccountStmtService(r.SourceRequest()).List(accountID, r.UserID, businessID)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	statementsJSON, err := json.Marshal(statements)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	resp, err := api.Success(r, string(statementsJSON), false)
	resp.Headers[http.CanonicalHeaderKey("Content-Type")] = "application/json"

	return resp, err
}

func HandleStatementAPIRequest(request api.APIRequest) (api.APIResponse, error) {
	method := strings.ToUpper(request.HTTPMethod)
	businessID, err := shared.ParseBusinessID(request.GetPathParam("businessId"))
	if err != nil {
		return api.BadRequest(request, err)
	}

	accountID := request.GetPathParam("accountId")
	statementID := request.GetPathParam("statementId")

	if accountID != "" && businessID != "" && statementID != "" {
		switch method {
		case http.MethodGet:
			return getStatementByID(request, statementID, accountID, businessID)
		default:
			return api.NotSupported(request)
		}
	}

	if accountID != "" && businessID != "" {
		switch method {
		case http.MethodGet:
			return getStatements(request, accountID, businessID)
		default:
			return api.NotSupported(request)
		}
	}

	return api.NotSupported(request)
}
