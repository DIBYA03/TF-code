/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

// Package for handling contact's linked accounts
package contact

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/wiseco/core-platform/api"
	request "github.com/wiseco/core-platform/services/banking/business"
	contact "github.com/wiseco/core-platform/services/banking/business/contact"
	"github.com/wiseco/core-platform/shared"
)

type RequestPostBody = request.RequestInitiate

func getRequest(r api.APIRequest, id shared.PaymentRequestID, contactId string, businessID shared.BusinessID) (api.APIResponse, error) {
	b, err := contact.NewMoneyRequestService(r.SourceRequest()).GetById(id, contactId, businessID)
	if err != nil {
		return api.BadRequest(r, err)
	}

	requestJSON, err := json.Marshal(b)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	return api.Success(r, string(requestJSON), false)
}

func getContactRequests(r api.APIRequest, contactId string, businessID shared.BusinessID) (api.APIResponse, error) {
	offset, _ := r.GetQueryIntParamWithDefault("offset", 0)
	limit, _ := r.GetQueryIntParamWithDefault("limit", 10)

	list, err := contact.NewMoneyRequestService(r.SourceRequest()).GetByContactId(offset, limit, contactId, businessID)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	requestListJSON, err := json.Marshal(list)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	return api.Success(r, string(requestListJSON), false)
}

func initiateRequest(r api.APIRequest, contactId string, businessID shared.BusinessID) (api.APIResponse, error) {
	var requestBody RequestPostBody

	err := json.Unmarshal([]byte(r.Body), &requestBody)
	if err != nil {
		return api.BadRequest(r, err)
	}

	requestBody.BusinessID = businessID
	requestBody.ContactId = contactId
	requestBody.CreatedUserID = r.UserID

	id, err := contact.NewMoneyRequestService(r.SourceRequest()).Request(&requestBody)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	transferJSON, _ := json.Marshal(id)
	return api.Success(r, string(transferJSON), false)
}

//HandleMoneyRequestAPIRequests handles the api request
func HandleMoneyRequestAPIRequests(request api.APIRequest) (api.APIResponse, error) {
	var method = strings.ToUpper(request.HTTPMethod)

	contactId := request.GetPathParam("contactId")
	if contactId == "" {
		return api.NotFoundError(request, errors.New("missing contact id"))
	}

	businessID, err := shared.ParseBusinessID(request.GetPathParam("businessId"))
	if err != nil {
		return api.BadRequestError(request, err)
	}

	requestId := request.GetPathParam("requestId")
	if requestId != "" {
		rID, err := shared.ParsePaymentRequestID(requestId)
		if err != nil {
			return api.BadRequestError(request, err)
		}

		switch method {
		case http.MethodGet:
			return getRequest(request, rID, contactId, businessID)
		default:
			return api.NotSupported(request)
		}

	}

	switch method {
	case http.MethodGet:
		return getContactRequests(request, contactId, businessID)
	case http.MethodPost:
		return initiateRequest(request, contactId, businessID)
	default:
		return api.NotSupported(request)
	}

}
