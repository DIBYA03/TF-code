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

type CapturePostBody = payment.PaymentCaptureRequest

func capturePayment(r api.APIRequest, businessID shared.BusinessID) (api.APIResponse, error) {
	var requestBody CapturePostBody

	err := json.Unmarshal([]byte(r.Body), &requestBody)
	if err != nil {
		return api.BadRequest(r, err)
	}

	requestBody.BusinessID = businessID
	requestBody.CreatedUserID = r.UserID

	err = payment.NewRequestService(r.SourceRequest()).CapturePayment(requestBody)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	return api.Success(r, "", false)
}

func HandlePaymentCaptureAPIRequests(request api.APIRequest) (api.APIResponse, error) {
	var method = strings.ToUpper(request.HTTPMethod)

	if request.BusinessID == nil {
		return api.BadRequestError(request, errors.New("Missing header X-Wise-Business-ID"))
	}

	switch method {
	case http.MethodPost:
		return capturePayment(request, *request.BusinessID)
	default:
		return api.NotSupported(request)
	}
}
