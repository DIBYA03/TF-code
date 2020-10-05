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

type ConnectionPostBody = payment.PaymentConnectionRequest

func getConnectionToken(r api.APIRequest, businessID shared.BusinessID) (api.APIResponse, error) {
	requestBody := ConnectionPostBody{}

	requestBody.UserID = r.UserID
	requestBody.BusinessID = businessID

	connection, err := payment.NewRequestService(r.SourceRequest()).GetConnectionToken(requestBody)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	connectionJSON, err := json.Marshal(connection)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	return api.Success(r, string(connectionJSON), false)
}

func HandlePaymentAPIRequests(request api.APIRequest) (api.APIResponse, error) {
	var method = strings.ToUpper(request.HTTPMethod)

	if request.BusinessID == nil {
		return api.BadRequestError(request, errors.New("Missing header X-Wise-Business-ID"))
	}

	switch method {
	case http.MethodPost:
		// Get connection token
		return getConnectionToken(request, *request.BusinessID)
	default:
		return api.NotSupported(request)
	}
}
