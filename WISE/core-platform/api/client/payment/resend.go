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

type RequestResendPostBody struct {
	Requests []payment.PaymentRequestResend `json:"requests"`
}

func resendPaymentRequest(r api.APIRequest, businessID shared.BusinessID) (api.APIResponse, error) {
	var postBody RequestResendPostBody

	err := json.Unmarshal([]byte(r.Body), &postBody)
	if err != nil {
		return api.BadRequest(r, err)
	}

	err = payment.NewRequestService(r.SourceRequest()).Resend(postBody.Requests, businessID)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	return api.Success(r, "", false)
}

func HandleRequestAPIResendRequest(request api.APIRequest) (api.APIResponse, error) {

	var method = strings.ToUpper(request.HTTPMethod)

	if request.BusinessID == nil {
		return api.BadRequestError(request, errors.New("missing header X-Wise-Business-ID"))
	}

	switch method {
	case http.MethodPost:
		return resendPaymentRequest(request, *request.BusinessID)
	default:
		return api.NotSupported(request)
	}
}
