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
	"github.com/wiseco/core-platform/services/banking"
	"github.com/wiseco/core-platform/services/banking/business"
	block "github.com/wiseco/core-platform/services/banking/business"
	"github.com/wiseco/core-platform/shared"
)

type BlockPostBody = block.BankCardBlockCreate

func blockCard(r api.APIRequest, cardId string, businessID shared.BusinessID) (api.APIResponse, error) {

	requestBody := BlockPostBody{}

	requestBody.BusinessID = businessID
	requestBody.CardID = cardId
	requestBody.CardholderID = r.UserID
	requestBody.OriginatedFrom = banking.OriginatedFromClientApplication

	// Client can only add temporary block
	requestBody.BlockID = banking.CardBlockIDLocked

	reason := "User blocked"
	requestBody.Reason = &reason

	s := r.SourceRequest()

	id, err := business.NewCardService(s).BlockBankCard(&requestBody)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	contactJSON, err := json.Marshal(id)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	return api.Success(r, string(contactJSON), false)
}

func HandleBankCardBlockAPIRequests(request api.APIRequest) (api.APIResponse, error) {
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
		return blockCard(request, cardId, businessID)
	default:
		return api.NotSupported(request)
	}

}
