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

type TransferRequestPostBody = payment.TransferRequestCreate

func sendTransferRequest(r api.APIRequest, businessID shared.BusinessID) (api.APIResponse, error) {
	var requestBody TransferRequestPostBody

	err := json.Unmarshal([]byte(r.Body), &requestBody)
	if err != nil {
		return api.BadRequest(r, err)
	}

	requestBody.CreatedUserID = r.UserID
	requestBody.BusinessID = businessID

	connection, err := payment.NewTransferService(r.SourceRequest()).SendTransferRequest(requestBody)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	requestJSON, err := json.Marshal(connection)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	return api.Success(r, string(requestJSON), false)
}

func HandleTransferAPIRequests(request api.APIRequest) (api.APIResponse, error) {
	var method = strings.ToUpper(request.HTTPMethod)

	if request.BusinessID == nil {
		return api.BadRequestError(request, errors.New("Missing header X-Wise-Business-ID"))
	}

	switch method {
	case http.MethodPost:
		return sendTransferRequest(request, *request.BusinessID)
	default:
		return api.NotSupported(request)
	}
}
