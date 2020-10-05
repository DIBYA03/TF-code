package business

import (
	"encoding/json"

	"github.com/wiseco/core-platform/api"
	"github.com/wiseco/core-platform/services/csp/review"
	"github.com/wiseco/core-platform/shared"
)

func handleBusinessStatus(businessID string, r api.APIRequest) (api.APIResponse, error) {
	bID, err := shared.ParseBusinessID(businessID)
	if err != nil {
		return api.BadRequest(r, err)
	}

	status, err := review.New(r.SourceRequest()).GetStatus(bID)
	if err != nil {
		return api.InternalServerError(r, err)
	}
	resp, err := json.Marshal(status)
	if err != nil {
		return api.InternalServerError(r, err)
	}
	return api.Success(r, string(resp), false)
}
