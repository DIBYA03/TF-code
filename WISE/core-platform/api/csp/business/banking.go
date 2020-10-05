package business

import (
	"encoding/json"

	"github.com/wiseco/core-platform/api"
	"github.com/wiseco/core-platform/services"
	"github.com/wiseco/core-platform/services/banking/business"
	csp "github.com/wiseco/core-platform/services/csp/business"
	"github.com/wiseco/core-platform/shared"
)

func handleAccountCreation(id shared.BusinessID, r api.APIRequest) (api.APIResponse, error) {
	err := csp.NewBanking(r.SourceRequest()).CreateBankAccount(id)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	status := struct {
		Message string `json:"message"`
	}{Message: "Success"}
	resp, _ := json.Marshal(status)
	return api.Success(r, string(resp), false)
}

func handleCardCreation(id shared.BusinessID, r api.APIRequest) (api.APIResponse, error) {
	if err := csp.NewBanking(r.SourceRequest()).CreateCard(id); err != nil {
		return api.InternalServerError(r, err)
	}
	status := struct {
		Message string `json:"message"`
	}{Message: "Success"}
	resp, _ := json.Marshal(status)
	return api.Success(r, string(resp), false)
}

func handleCardList(businessID shared.BusinessID, r api.APIRequest) (api.APIResponse, error) {
	list, err := csp.NewBanking(r.SourceRequest()).Cards(businessID)
	if err != nil {
		return api.InternalServerError(r, err)
	}
	resp, err := json.Marshal(list)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	return api.Success(r, string(resp), false)
}

func handleCardByID(businessID shared.BusinessID, cardID string, r api.APIRequest) (api.APIResponse, error) {
	card, err := csp.NewBanking(r.SourceRequest()).CardByID(businessID, cardID)
	if err != nil {
		_, isENF := err.(*services.ErrorNotFound)
		if isENF {
			return api.NotFound(r)
		}

		return api.InternalServerError(r, err)
	}

	resp, err := json.Marshal(card)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	return api.Success(r, string(resp), false)
}

func handleBusinessAccount(businessID shared.BusinessID, r api.APIRequest) (api.APIResponse, error) {
	list, err := csp.NewBanking(r.SourceRequest()).AccountByBusinessID(businessID)
	if err != nil {
		_, isENF := err.(*services.ErrorNotFound)
		if isENF {
			return api.NotFound(r)
		}

		return api.InternalServerError(r, err)
	}

	resp, _ := json.Marshal(list)
	return api.Success(r, string(resp), false)
}

func handleExternalBusinessAccount(businessID shared.BusinessID, r api.APIRequest) (api.APIResponse, error) {
	list, err := csp.NewBanking(r.SourceRequest()).ExternalAccountByBusinessID(businessID)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	resp, _ := json.Marshal(list)
	return api.Success(r, string(resp), false)
}

func handlePromotionalFund(id shared.BusinessID, r api.APIRequest) (api.APIResponse, error) {
	if err := csp.NewBanking(r.SourceRequest()).FundPromotion(id); err != nil {
		return api.InternalServerError(r, err)
	}
	status := struct {
		Message string `json:"message"`
	}{Message: "Success"}

	resp, _ := json.Marshal(status)

	return api.Success(r, string(resp), false)
}

func handleCardReissue(bID shared.BusinessID, cID string, r api.APIRequest) (api.APIResponse, error) {
	var body business.CardReissueRequest
	if err := json.Unmarshal([]byte(r.Body), &body); err != nil {
		return api.BadRequest(r, err)
	}

	body.BusinessID = bID
	body.CardID = cID

	card, err := csp.NewBanking(r.SourceRequest()).ReissueCard(body)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	resp, _ := json.Marshal(card)
	return api.Success(r, string(resp), false)
}
