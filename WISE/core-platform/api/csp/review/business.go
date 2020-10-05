package review

import (
	"encoding/json"
	"errors"

	"github.com/wiseco/core-platform/api"
	"github.com/wiseco/core-platform/services"
	"github.com/wiseco/core-platform/services/csp"
	"github.com/wiseco/core-platform/services/csp/business"
	"github.com/wiseco/core-platform/services/csp/review"
	"github.com/wiseco/core-platform/shared"
)

func handleReview(status string, reviewSubstatus string, limit, offset int, r api.APIRequest) (api.APIResponse, error) {

	st, ok := csp.ReviewStatus.NewStatus(status)
	if !ok {
		return api.BadRequest(r, errors.New("Invalid params"))
	}

	isValid := csp.CheckReviewSubstatus(reviewSubstatus)
	if reviewSubstatus != "" && !isValid {
		return api.BadRequest(r, errors.New("Invalid params"))
	}
	rst := csp.ReviewSubstatus(reviewSubstatus)

	list, err := business.NewCSPService().CSPBusinessListStatus(limit, offset, st, rst)
	if err != nil {
		return api.InternalServerError(r, err)
	}
	resp, err := json.Marshal(list)
	if err != nil {
		return api.InternalServerError(r, err)
	}
	return api.Success(r, string(resp), false)
}

func handleReviewUpdate(status, businessID string, r api.APIRequest) (api.APIResponse, error) {
	st, ok := csp.ReviewStatus.NewStatus(status)
	if !ok {
		return api.BadRequest(r, errors.New("Invalid params"))
	}

	bID, err := shared.ParseBusinessID(businessID)
	if err != nil {
		return api.BadRequest(r, err)
	}

	review, err := review.NewReviewService(r.SourceRequest()).UpdateReview(bID, st)
	if err != nil {
		return api.InternalServerError(r, err)
	}
	resp, err := json.Marshal(review)
	if err != nil {
		return api.InternalServerError(r, err)
	}
	return api.Success(r, string(resp), false)
}

func handleReviewID(id string, r api.APIRequest) (api.APIResponse, error) {

	review, err := business.NewCSPService().CSPBusinessByID(id)
	if err != nil {
		return api.InternalServerError(r, err)
	}
	resp, err := json.Marshal(review)
	if err != nil {
		return api.InternalServerError(r, err)
	}
	return api.Success(r, string(resp), false)
}

func handleDecline(id string, r api.APIRequest) (api.APIResponse, error) {
	bID, err := shared.ParseBusinessID(id)
	if err != nil {
		return api.BadRequest(r, err)
	}
	status, err := review.NewReviewService(r.SourceRequest()).Decline(bID)
	if err != nil {
		_, isENF := err.(*services.ErrorNotFound)
		if isENF {
			return api.NotFound(r)
		}

		return api.InternalServerError(r, err)
	}

	resp, err := json.Marshal(status)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	return api.Success(r, string(resp), false)
}

func handleApprove(id string, r api.APIRequest) (api.APIResponse, error) {
	bID, err := shared.ParseBusinessID(id)
	if err != nil {
		return api.BadRequest(r, err)
	}

	status, err := review.NewReviewService(r.SourceRequest()).Approve(bID)
	if err != nil {
		_, isENF := err.(*services.ErrorNotFound)
		if isENF {
			return api.NotFound(r)
		}

		return api.InternalServerError(r, err)
	}

	resp, err := json.Marshal(status)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	return api.Success(r, string(resp), false)
}
