/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

package payment

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/wiseco/go-lib/id"

	"github.com/wiseco/core-platform/services/invoice"

	"github.com/wiseco/core-platform/api"
	payment "github.com/wiseco/core-platform/services/payment"
	"github.com/wiseco/core-platform/shared"
)

func receiptInvoiceUrlByID(r api.APIRequest, requestID string, businessID shared.BusinessID) (api.APIResponse, error) {
	var respURL string
	var err error
	// check if the request type is POS
	rID, err := shared.ParsePaymentRequestID(requestID)
	if err != nil {
		return api.BadRequestError(r, err)
	}
	request, err := payment.NewRequestService(r.SourceRequest()).GetByID(rID, *r.BusinessID)
	if err != nil {
		return api.BadRequestError(r, err)
	}
	if os.Getenv("USE_INVOICE_SERVICE") == "true" && *request.RequestType != payment.PaymentRequestTypePOS {
		invSvc, err := invoice.NewInvoiceService()
		if err != nil {
			return api.InternalServerError(r, err)
		}
		invoiceID, err := id.ParseInvoiceID(fmt.Sprintf("%s%s", id.IDPrefixInvoice, requestID))
		if err != nil {
			return api.InternalServerError(r, err)
		}
		inv, err := invSvc.GetInvoiceByID(invoiceID)
		if err != nil {
			return api.InternalServerError(r, err)
		}
		respURL = inv.InvoiceViewLink
	} else {
		url, err := payment.NewInvoiceService(r.SourceRequest()).GetSignedURL(requestID, businessID)
		if err != nil {
			return api.InternalServerError(r, err)
		}
		if url == nil {
			return api.InternalServerError(r, errors.New("Unable to get signed url"))
		}
		respURL = *url
	}

	if err != nil && err.Error() == "not found" {
		return api.NotFound(r)
	} else if err != nil {
		return api.InternalServerError(r, err)
	}

	URL := struct {
		URL *string `json:"url"`
	}{URL: &respURL}

	b, _ := json.Marshal(URL)
	return api.Success(r, string(b), false)
}

func HandleInvoiceAPIRequests(request api.APIRequest) (api.APIResponse, error) {
	var method = strings.ToUpper(request.HTTPMethod)

	if request.BusinessID == nil {
		return api.InternalServerError(request, errors.New("Missing header X-Wise-Business-ID"))
	}

	receiptID := request.GetPathParam("invoiceId")
	if len(receiptID) == 0 {
		return api.InternalServerError(request, errors.New("Missing invoice id"))
	}

	switch method {
	case http.MethodGet:
		return receiptInvoiceUrlByID(request, receiptID, *request.BusinessID)
	default:
		return api.NotSupported(request)
	}
}
