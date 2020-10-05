/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

// Package for handling linked account and its related api requests
package business

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/wiseco/core-platform/api"
	transfer "github.com/wiseco/core-platform/services/banking/business"
	"github.com/wiseco/core-platform/shared"
)

type TransferPostBody = transfer.TransferInitiate

func getBusinessTransfers(r api.APIRequest, businessID shared.BusinessID) (api.APIResponse, error) {
	offset, _ := r.GetQueryIntParamWithDefault("offset", 0)
	limit, _ := r.GetQueryIntParamWithDefault("limit", 10)

	s := r.SourceRequest()

	list, err := transfer.NewMoneyTransferService(s).GetByBusinessID(offset, limit, businessID, r.UserID)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	transferListJSON, err := json.Marshal(list)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	return api.Success(r, string(transferListJSON), false)
}

func initiateTransfer(r api.APIRequest, businessID shared.BusinessID) (api.APIResponse, error) {
	var requestBody TransferPostBody

	err := json.Unmarshal([]byte(r.Body), &requestBody)
	if err != nil {
		return api.BadRequest(r, err)
	}

	requestBody.BusinessID = businessID
	requestBody.CreatedUserID = r.UserID

	s := r.SourceRequest()

	id, err := transfer.NewMoneyTransferService(s).Transfer(&requestBody)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	transferJSON, _ := json.Marshal(id)
	return api.Success(r, string(transferJSON), false)
}

//HandleLinkedCardAPIRequests handles the api request
func HandleMoneyTransferAPIRequests(request api.APIRequest) (api.APIResponse, error) {
	var method = strings.ToUpper(request.HTTPMethod)

	businessId, err := shared.ParseBusinessID(request.GetPathParam("businessId"))
	if err != nil {
		return api.BadRequestError(request, err)
	}

	switch method {
	case http.MethodGet:
		return getBusinessTransfers(request, businessId)
	case http.MethodPost:
		return initiateTransfer(request, businessId)
	default:
		return api.NotSupported(request)
	}

}
