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

type CardReaderPostBody = payment.CardReaderCreate
type CardReaderPatchBody = payment.CardReaderUpdate

func deactivateCardReader(r api.APIRequest, cardReaderID shared.CardReaderID, businessID shared.BusinessID) (api.APIResponse, error) {
	err := payment.NewCardReaderService(r.SourceRequest()).Deactivate(cardReaderID, r.UserID, businessID)
	if err != nil {
		return api.BadRequest(r, err)
	}

	return api.Success(r, "", false)
}

func getCardReader(r api.APIRequest, cardReaderID shared.CardReaderID, businessID shared.BusinessID) (api.APIResponse, error) {
	cardReader, err := payment.NewCardReaderService(r.SourceRequest()).GetById(cardReaderID, r.UserID, businessID)
	if err != nil {
		return api.BadRequest(r, err)
	}

	cardReaderJSON, err := json.Marshal(cardReader)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	return api.Success(r, string(cardReaderJSON), false)
}

func getCardReaderList(r api.APIRequest, businessID shared.BusinessID) (api.APIResponse, error) {

	cardReader, err := payment.NewCardReaderService(r.SourceRequest()).List(r.UserID, businessID)
	if err != nil {
		return api.BadRequest(r, err)
	}

	cardReaderJSON, err := json.Marshal(cardReader)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	return api.Success(r, string(cardReaderJSON), false)
}

func createCardReader(r api.APIRequest, businessID shared.BusinessID) (api.APIResponse, error) {
	var requestBody CardReaderPostBody

	err := json.Unmarshal([]byte(r.Body), &requestBody)
	if err != nil {
		return api.BadRequest(r, err)
	}

	requestBody.BusinessID = businessID

	cardReader, err := payment.NewCardReaderService(r.SourceRequest()).Create(requestBody)
	if err != nil {
		return api.BadRequest(r, err)
	}

	cardReaderJSON, err := json.Marshal(cardReader)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	return api.Success(r, string(cardReaderJSON), false)
}

func updateCardReader(r api.APIRequest, cardReaderID shared.CardReaderID, businessID shared.BusinessID) (api.APIResponse, error) {
	var requestBody CardReaderPatchBody

	err := json.Unmarshal([]byte(r.Body), &requestBody)
	if err != nil {
		return api.BadRequest(r, err)
	}

	requestBody.ID = cardReaderID
	requestBody.BusinessID = businessID

	cardReader, err := payment.NewCardReaderService(r.SourceRequest()).Update(requestBody)
	if err != nil {
		return api.BadRequest(r, err)
	}

	cardReaderJSON, err := json.Marshal(cardReader)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	return api.Success(r, string(cardReaderJSON), false)
}

func HandleCardReaderAPIRequests(request api.APIRequest) (api.APIResponse, error) {
	var method = strings.ToUpper(request.HTTPMethod)

	if request.BusinessID == nil {
		return api.InternalServerError(request, errors.New("Missing header X-Wise-Business-ID"))
	}

	cardReaderID := request.GetPathParam("cardReaderId")

	switch method {
	case http.MethodPatch:
		if len(cardReaderID) > 0 {
			cardReaderID, err := shared.ParseCardReaderID(cardReaderID)
			if err != nil {
				return api.BadRequestError(request, err)
			}

			return updateCardReader(request, cardReaderID, *request.BusinessID)
		}
	case http.MethodGet:
		if len(cardReaderID) > 0 {
			cardReaderID, err := shared.ParseCardReaderID(cardReaderID)
			if err != nil {
				return api.BadRequestError(request, err)
			}

			return getCardReader(request, cardReaderID, *request.BusinessID)
		} else {
			return getCardReaderList(request, *request.BusinessID)
		}
	case http.MethodDelete:
		if len(cardReaderID) > 0 {
			cardReaderID, err := shared.ParseCardReaderID(cardReaderID)
			if err != nil {
				return api.BadRequestError(request, err)
			}

			return deactivateCardReader(request, cardReaderID, *request.BusinessID)
		}
	default:
		return api.NotSupported(request)
	}

	return api.NotSupported(request)
}
