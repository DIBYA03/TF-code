/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

package payment

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/wiseco/core-platform/api"
	payment "github.com/wiseco/core-platform/services/payment"
	"github.com/wiseco/core-platform/shared"
)

type RequestPostBody = payment.RequestInitiate
type RequestPatchBody = payment.RequestStatusUpdate

func getRequest(r api.APIRequest, id shared.PaymentRequestID, businessID shared.BusinessID) (api.APIResponse, error) {
	b, err := payment.NewRequestService(r.SourceRequest()).GetByID(id, businessID)
	if err != nil {
		return api.BadRequest(r, err)
	}

	requestJSON, err := json.Marshal(b)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	return api.Success(r, string(requestJSON), false)
}

func updateRequest(r api.APIRequest, ID shared.PaymentRequestID, businessID shared.BusinessID) (api.APIResponse, error) {
	var requestBody RequestPatchBody

	err := json.Unmarshal([]byte(r.Body), &requestBody)
	if err != nil {
		return api.BadRequest(r, err)
	}

	requestBody.BusinessID = businessID
	requestBody.ID = ID

	resp, err := payment.NewRequestService(r.SourceRequest()).UpdateRequestStatus(&requestBody)
	if err != nil {
		return api.BadRequest(r, err)
	}

	respJSON, err := json.Marshal(resp)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	return api.Success(r, string(respJSON), false)
}

func getRequests(r api.APIRequest, businessID shared.BusinessID) (api.APIResponse, error) {
	offset, _ := r.GetQueryIntParamWithDefault("offset", 0)
	limit, _ := r.GetQueryIntParamWithDefault("limit", 20)
	status := r.GetQueryParam("status")

	list, err := payment.NewRequestService(r.SourceRequest()).List(businessID, offset, limit, status)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	requestListJSON, err := json.Marshal(list)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	return api.Success(r, string(requestListJSON), false)
}

func getContactRequests(r api.APIRequest, contactID string, businessID shared.BusinessID) (api.APIResponse, error) {
	offset, _ := r.GetQueryIntParamWithDefault("offset", 0)
	limit, _ := r.GetQueryIntParamWithDefault("limit", 20)
	status := r.GetQueryParam("status")

	list, err := payment.NewRequestService(r.SourceRequest()).ListByContactID(businessID, contactID, offset, limit, status)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	requestListJSON, err := json.Marshal(list)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	return api.Success(r, string(requestListJSON), false)
}

func initiateRequest(r api.APIRequest, businessID shared.BusinessID) (api.APIResponse, error) {
	var requestBody RequestPostBody

	err := json.Unmarshal([]byte(r.Body), &requestBody)
	if err != nil {
		return api.BadRequest(r, err)
	}

	requestBody.CreatedUserID = r.UserID
	requestBody.BusinessID = businessID
	requestBody.IPAddress = &r.SourceIP

	connection, err := payment.NewRequestService(r.SourceRequest()).Request(&requestBody)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	connectionJSON, err := json.Marshal(connection)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	return api.Success(r, string(connectionJSON), false)
}

func HandleRequestAPIRequests(request api.APIRequest) (api.APIResponse, error) {

	var method = strings.ToUpper(request.HTTPMethod)

	if request.BusinessID == nil {
		return api.BadRequestError(request, errors.New("missing header X-Wise-Business-ID"))
	}

	requestID := request.GetPathParam("requestId")
	contactID := request.GetQueryParam("contactId")

	if len(requestID) > 0 {
		rID, err := shared.ParsePaymentRequestID(requestID)
		if err != nil {
			api.BadRequestError(request, err)
		}

		switch method {
		case http.MethodGet:
			return getRequest(request, rID, *request.BusinessID)
		case http.MethodPatch:
			return updateRequest(request, rID, *request.BusinessID)
		default:
			return api.NotSupported(request)
		}
	}

	if len(contactID) > 0 {
		switch method {
		case http.MethodGet:
			return getContactRequests(request, contactID, *request.BusinessID)
		default:
			return api.NotSupported(request)
		}
	}

	switch method {
	case http.MethodGet:
		return getRequests(request, *request.BusinessID)
	case http.MethodPost:
		return initiateRequest(request, *request.BusinessID)
	default:
		return api.NotSupported(request)
	}
}
