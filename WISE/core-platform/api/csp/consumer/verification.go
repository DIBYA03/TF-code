package consumer

import (
	"encoding/json"

	"github.com/wiseco/core-platform/api"
	"github.com/wiseco/core-platform/shared"

	"github.com/wiseco/core-platform/services/csp/review"
)

func handleCheckStatus(consumerID string, r api.APIRequest) (api.APIResponse, error) {
	id, err := shared.ParseConsumerID(consumerID)
	if err != nil {
		return api.BadRequest(r, err)
	}

	status, err := review.NewConsumerVerfication(r.SourceRequest()).GetVerification(id)
	if err != nil {
		return api.InternalServerError(r, err)
	}
	resp, err := json.Marshal(status)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	return api.Success(r, string(resp), false)
}

func handleVerification(consumerID string, r api.APIRequest) (api.APIResponse, error) {
	id, err := shared.ParseConsumerID(consumerID)
	if err != nil {
		return api.BadRequest(r, err)
	}

	verification, err := review.NewConsumerVerfication(r.SourceRequest()).StartVerification(id)
	if err != nil {
		return api.InternalServerError(r, err)
	}
	resp, err := json.Marshal(verification)
	if err != nil {
		return api.InternalServerError(r, err)
	}
	return api.Success(r, string(resp), false)
}

func handleGetVerification(consumerID string, r api.APIRequest) (api.APIResponse, error) {
	id, err := shared.ParseConsumerID(consumerID)
	if err != nil {
		return api.BadRequest(r, err)
	}

	verification, err := review.NewConsumerVerfication(r.SourceRequest()).GetVerification(id)
	if err != nil {
		return api.InternalServerError(r, err)
	}
	resp, err := json.Marshal(verification)
	if err != nil {
		return api.InternalServerError(r, err)
	}
	return api.Success(r, string(resp), false)
}
