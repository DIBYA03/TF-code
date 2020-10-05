/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

package payment

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/wiseco/core-platform/api"
	payment "github.com/wiseco/core-platform/services/payment"
	"github.com/wiseco/core-platform/shared"
)

type ReceiptPostBody = payment.CardReaderReceiptCreate

func createPOSReceipt(r api.APIRequest, businessID shared.BusinessID) (api.APIResponse, error) {
	var requestBody ReceiptPostBody

	err := json.Unmarshal([]byte(r.Body), &requestBody)
	if err != nil {
		return api.BadRequest(r, err)
	}

	requestBody.BusinessID = businessID

	err = payment.NewPaymentService(r.SourceRequest()).SendCardReaderReceipt(requestBody)
	if err != nil {
		return api.BadRequest(r, err)
	}

	return api.Success(r, "", false)
}

func receiptReceiptUrlByID(r api.APIRequest, id string, businessID shared.BusinessID) (api.APIResponse, error) {
	rcpSvc := payment.NewReceiptService(r.SourceRequest())
	isPOSRequest := rcpSvc.IsPOSInvoice(id, businessID)
	var respURL string
	if os.Getenv("USE_INVOICE_SERVICE") == "true" && !isPOSRequest {
		url, err := rcpSvc.GetReceiptURLForInvoice(id, businessID)
		if err != nil {
			return api.InternalServerError(r, err)
		}

		if url == nil {
			return api.InternalServerError(r, errors.New("Unable to get signed url"))
		}

		respURL = *url
	} else {

		url, err := payment.NewReceiptService(r.SourceRequest()).GetSignedURL(id, businessID)
		if err != nil {
			return api.InternalServerError(r, err)
		}

		if err != nil && err.Error() == "not found" {
			return api.NotFound(r)
		} else if err != nil {
			return api.InternalServerError(r, err)
		}

		if url == nil {
			return api.InternalServerError(r, errors.New("Unable to get signed url"))
		}
		respURL = *url
	}

	URL := struct {
		URL *string `json:"url"`
	}{URL: &respURL}

	b, _ := json.Marshal(URL)
	return api.Success(r, string(b), false)
}

func HandleReceiptAPIRequests(request api.APIRequest) (api.APIResponse, error) {
	var method = strings.ToUpper(request.HTTPMethod)

	if request.BusinessID == nil {
		return api.BadRequestError(request, errors.New("missing header X-Wise-Business-ID"))
	}

	receiptID := request.GetPathParam("receiptId")
	if len(receiptID) > 0 {
		switch method {
		case http.MethodGet:
			return receiptReceiptUrlByID(request, receiptID, *request.BusinessID)
		default:
			return api.NotSupported(request)
		}
	}

	switch method {
	case http.MethodPost:
		return createPOSReceipt(request, *request.BusinessID)
	default:
		return api.NotSupported(request)
	}
}
