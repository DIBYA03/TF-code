package consumer

import (
	"encoding/json"

	"github.com/wiseco/core-platform/api"
	"github.com/wiseco/core-platform/services"
	csp "github.com/wiseco/core-platform/services/csp/consumer"
	"github.com/wiseco/core-platform/shared"
)

func handleConsumerList(r api.APIRequest) (api.APIResponse, error) {
	list, err := csp.New().UserList()
	if err != nil {
		return api.InternalServerError(r, err)
	}

	resp, err := json.Marshal(list)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	return api.Success(r, string(resp), false)
}

func handleConsumerID(consumerID string, r api.APIRequest) (api.APIResponse, error) {
	conID, err := shared.ParseConsumerID(consumerID)
	if err != nil {
		return api.BadRequest(r, err)
	}

	user, err := csp.New().ByID(conID)
	if err != nil {
		_, isENF := err.(*services.ErrorNotFound)
		if isENF {
			return api.NotFound(r)
		}

		return api.InternalServerError(r, err)
	}

	resp, err := json.Marshal(user)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	return api.Success(r, string(resp), false)

}

func handleConsumerCSPConsumer(consumerID string, r api.APIRequest) (api.APIResponse, error) {

	conID, err := shared.ParseConsumerID(consumerID)
	if err != nil {
		return api.BadRequest(r, err)
	}
	review, err := csp.NewCSPService().CSPConsumerByConsumerID(conID)
	if err != nil {
		_, isENF := err.(*services.ErrorNotFound)
		if isENF {
			return api.NotFound(r)
		}

		return api.InternalServerError(r, err)
	}

	resp, err := json.Marshal(review)
	if err != nil {
		return api.InternalServerError(r, err)
	}
	return api.Success(r, string(resp), false)
}

func handleConsumerUpdate(consumerID string, r api.APIRequest) (api.APIResponse, error) {
	conID, err := shared.ParseConsumerID(consumerID)
	if err != nil {
		return api.BadRequest(r, err)
	}

	var body csp.Update
	if err := json.Unmarshal([]byte(r.Body), &body); err != nil {
		return api.BadRequest(r, err)
	}
	user, err := csp.NewWithSource(r.SourceRequest()).UpdateID(conID, &body)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	resp, err := json.Marshal(user)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	return api.Success(r, string(resp), false)
}

func handleState(CSPConsumerID string, r api.APIRequest) (api.APIResponse, error) {
	list, err := csp.NewStateService().List(CSPConsumerID)
	if err != nil {
		return api.InternalServerError(r, err)
	}
	resp, err := json.Marshal(list)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	return api.Success(r, string(resp), false)
}
