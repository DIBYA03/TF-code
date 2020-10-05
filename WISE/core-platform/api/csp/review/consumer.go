package review

import (
	"encoding/json"

	"github.com/wiseco/core-platform/api"
	"github.com/wiseco/core-platform/services/csp/consumer"
)

func handleConsumerReviewList(status string, limit, offset int, r api.APIRequest) (api.APIResponse, error) {
	list, err := consumer.NewCSPService().List(status, limit, offset)
	if err != nil {
		return api.InternalServerError(r, err)
	}
	resp, err := json.Marshal(list)
	if err != nil {
		return api.InternalServerError(r, err)
	}
	return api.Success(r, string(resp), false)
}

func handleConsumerReviewID(reviewID string, r api.APIRequest) (api.APIResponse, error) {
	review, err := consumer.NewCSPService().CSPConsumerByID(reviewID)
	if err != nil {
		return api.InternalServerError(r, err)
	}
	resp, err := json.Marshal(review)
	if err != nil {
		return api.InternalServerError(r, err)
	}
	return api.Success(r, string(resp), false)
}
