package business

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/wiseco/core-platform/api"
	"github.com/wiseco/core-platform/services"
	csp "github.com/wiseco/core-platform/services/csp/business"
	"github.com/wiseco/core-platform/shared"
)

func deactivateCardReader(r api.APIRequest, cardReaderID shared.CardReaderID, businessID shared.BusinessID) (api.APIResponse, error) {
	err := csp.NewCardReaderService(r.SourceRequest()).Deactivate(cardReaderID, businessID)
	if err != nil {
		return api.BadRequest(r, err)
	}

	return api.Success(r, "", false)
}

func getCardReader(r api.APIRequest, cardReaderID shared.CardReaderID, businessID shared.BusinessID) (api.APIResponse, error) {
	cardReader, err := csp.NewCardReaderService(r.SourceRequest()).ByID(cardReaderID, businessID)
	if err != nil {
		_, isENF := err.(*services.ErrorNotFound)
		if isENF {
			return api.NotFound(r)
		}

		return api.InternalServerError(r, err)
	}

	cardReaderJSON, err := json.Marshal(cardReader)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	return api.Success(r, string(cardReaderJSON), false)
}

func getCardReaderList(r api.APIRequest, businessID shared.BusinessID) (api.APIResponse, error) {

	cardReader, err := csp.NewCardReaderService(r.SourceRequest()).List(businessID)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	cardReaderJSON, err := json.Marshal(cardReader)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	return api.Success(r, string(cardReaderJSON), false)
}

func createCardReader(r api.APIRequest, businessID shared.BusinessID) (api.APIResponse, error) {
	var requestBody csp.CardReaderCreate

	err := json.Unmarshal([]byte(r.Body), &requestBody)
	if err != nil {
		return api.BadRequest(r, err)
	}

	requestBody.BusinessID = businessID

	cardReader, err := csp.NewCardReaderService(r.SourceRequest()).Create(requestBody)
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
	var requestBody csp.CardReaderUpdate

	err := json.Unmarshal([]byte(r.Body), &requestBody)
	if err != nil {
		return api.BadRequest(r, err)
	}

	cardReader, err := csp.NewCardReaderService(r.SourceRequest()).Update(cardReaderID, requestBody)
	if err != nil {
		return api.BadRequest(r, err)
	}

	cardReaderJSON, err := json.Marshal(cardReader)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	return api.Success(r, string(cardReaderJSON), false)
}

func HandleCardReaderAPIRequests(r api.APIRequest) (api.APIResponse, error) {
	var method = strings.ToUpper(r.HTTPMethod)

	businessID, err := shared.ParseBusinessID(r.GetPathParam("businessId"))
	if err != nil {
		return api.BadRequestError(r, err)
	}

	cardReaderID := r.GetPathParam("readerId")
	if cardReaderID != "" {
		cID, err := shared.ParseCardReaderID(cardReaderID)
		if err != nil {
			return api.BadRequestError(r, err)
		}

		switch method {
		case http.MethodPatch:
			return updateCardReader(r, cID, businessID)
		case http.MethodGet:
			return getCardReader(r, cID, businessID)
		case http.MethodDelete:
			return deactivateCardReader(r, cID, businessID)
		default:
			return api.NotSupported(r)
		}
	}

	switch method {
	case http.MethodPost:
		return createCardReader(r, businessID)
	case http.MethodGet:
		return getCardReaderList(r, businessID)
	default:
		return api.NotSupported(r)
	}
}
