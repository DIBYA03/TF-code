package support

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/wiseco/core-platform/api"
	business "github.com/wiseco/core-platform/services/banking/business"
	csp "github.com/wiseco/core-platform/services/csp/support"
)

//HandleUnblockCard ..
func HandleUnblockCard(r api.APIRequest) (api.APIResponse, error) {
	cardID := r.GetPathParam("cardId")
	if cardID != "" {
		switch strings.ToUpper(r.HTTPMethod) {
		case http.MethodPost:
			return unblockCard(cardID, r)
		default:
			return api.NotSupported(r)
		}
	}

	return api.NotSupported(r)
}

//HandleAccountBlocks ..
func HandleCardBlocks(r api.APIRequest) (api.APIResponse, error) {
	cardID := r.GetPathParam("cardId")
	if cardID != "" {
		switch strings.ToUpper(r.HTTPMethod) {
		case http.MethodPost:
			return blockCard(cardID, r)
		case http.MethodGet:
			return listBlocks(cardID, r)
		default:
			return api.NotSupported(r)
		}
	}

	return api.NotSupported(r)
}

//HandleCardStatus ..
func HandleCardStatus(r api.APIRequest) (api.APIResponse, error) {
	cardID := r.GetPathParam("cardId")
	if cardID != "" {
		switch strings.ToUpper(r.HTTPMethod) {
		case http.MethodGet:
			return cardStatus(cardID, r)
		default:
			return api.NotSupported(r)
		}
	}

	return api.NotSupported(r)
}

func blockCard(ID string, r api.APIRequest) (api.APIResponse, error) {

	var body business.BankCardBlockCreate
	err := json.Unmarshal([]byte(r.Body), &body)
	if err != nil {
		return api.BadRequest(r, err)
	}

	if len(body.OriginatedFrom) == 0 || len(body.BlockID) == 0 {
		return api.BadRequest(r, fmt.Errorf("Missing or invalid params. body:%v", body))
	}

	body.CardID = ID
	card, err := csp.NewCardSupportService(r.SourceRequest()).Block(&body)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	resp, err := json.Marshal(card)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	return api.Success(r, string(resp), false)
}

func listBlocks(ID string, r api.APIRequest) (api.APIResponse, error) {

	blocks, err := csp.NewCardSupportService(r.SourceRequest()).ListBlocks(ID)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	resp, err := json.Marshal(blocks)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	return api.Success(r, string(resp), false)
}

func unblockCard(ID string, r api.APIRequest) (api.APIResponse, error) {
	card, err := csp.NewCardSupportService(r.SourceRequest()).UnBlock(ID)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	resp, err := json.Marshal(card)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	return api.Success(r, string(resp), false)
}

func cardStatus(ID string, r api.APIRequest) (api.APIResponse, error) {
	card, err := csp.NewCardSupportService(r.SourceRequest()).CheckBlockStatus(ID)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	resp, err := json.Marshal(card)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	return api.Success(r, string(resp), false)
}
