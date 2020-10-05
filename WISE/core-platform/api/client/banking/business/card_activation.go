/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

// Package for business banking apis
package business

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/wiseco/core-platform/api"
	"github.com/wiseco/core-platform/services/banking/business"
	activate "github.com/wiseco/core-platform/services/banking/business"
	"github.com/wiseco/core-platform/shared"
)

type ActivatePostBody = activate.BankCardActivate

func activateCard(r api.APIRequest, cardId string, businessID shared.BusinessID) (api.APIResponse, error) {
	var requestBody ActivatePostBody

	err := json.Unmarshal([]byte(r.Body), &requestBody)
	if err != nil {
		return api.BadRequest(r, err)
	}

	requestBody.BusinessID = businessID
	requestBody.Id = cardId
	requestBody.CardholderID = r.UserID

	s := r.SourceRequest()

	id, err := business.NewCardService(s).ActivateBankCard(&requestBody)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	contactJSON, err := json.Marshal(id)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	return api.Success(r, string(contactJSON), false)
}

func HandleBankCardActivationAPIRequests(request api.APIRequest) (api.APIResponse, error) {
	var method = strings.ToUpper(request.HTTPMethod)

	businessID, err := shared.ParseBusinessID(request.GetPathParam("businessId"))
	if err != nil {
		return api.BadRequest(request, err)
	}

	cardId := request.GetPathParam("cardId")
	if cardId == "" {
		return api.InternalServerError(request, errors.New("not found"))
	}

	switch method {
	case http.MethodPost:
		return activateCard(request, cardId, businessID)
	default:
		return api.NotSupported(request)
	}

}
