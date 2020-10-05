/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

// Package for business banking apis
package business

import (
	"database/sql"
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

type UnblockPostBody = block.BankCardBlockDelete

func unblockCard(r api.APIRequest, cardId string, businessID shared.BusinessID) (api.APIResponse, error) {

	requestBody := UnblockPostBody{}

	s := r.SourceRequest()

	// Only temporary blocks can be removed
	block, err := business.NewCardService(s).GetByPartnerBlockID(banking.CardBlockIDLocked, cardId)
	if err != nil && err != sql.ErrNoRows {
		return api.InternalServerError(r, err)
	}

	if err != nil {
		return api.InternalServerError(r, errors.New("Card cannot be unblocked"))
	}

	requestBody.BusinessID = businessID
	requestBody.CardID = cardId
	requestBody.CardholderID = r.UserID
	requestBody.ID = block.ID

	id, err := business.NewCardService(s).UnBlockBankCard(&requestBody)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	respJSON, err := json.Marshal(id)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	return api.Success(r, string(respJSON), false)
}

func HandleBankCardUnblockAPIRequests(request api.APIRequest) (api.APIResponse, error) {
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
		return unblockCard(request, cardId, businessID)
	default:
		return api.NotSupported(request)
	}

}
