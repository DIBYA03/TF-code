package review

import (
	"errors"
	"net/http"
	"strings"

	"github.com/wiseco/core-platform/api"
)

// BusinessStatusListRequest list business review status
func BusinessStatusListRequest(r api.APIRequest) (api.APIResponse, error) {
	if strings.ToUpper(r.HTTPMethod) != http.MethodGet {
		return api.NotSupported(r)
	}
	limit, _ := r.GetQueryIntParamWithDefault("limit", 30)
	offset, _ := r.GetQueryIntParamWithDefault("offset", 0)
	reviewSubstatus := r.GetQueryParam("reviewSubstatus")
	if status := r.GetPathParam("status"); status != "" {
		return handleReview(status, reviewSubstatus, limit, offset, r)
	}
	return api.NotSupported(r)
}

// BusinessReviewItemRequest gets a business review item by id
func BusinessReviewItemRequest(r api.APIRequest) (api.APIResponse, error) {
	if strings.ToUpper(r.HTTPMethod) != http.MethodGet {
		return api.NotSupported(r)
	}

	id := r.GetPathParam("id")

	return handleReviewID(id, r)
}

// BusinessStatusChangeRequest .. business review changes
func BusinessStatusChangeRequest(r api.APIRequest) (api.APIResponse, error) {
	if strings.ToUpper(r.HTTPMethod) != http.MethodPost {
		return api.NotSupported(r)
	}
	status := r.GetPathParam("status")
	id := r.GetPathParam("businessId")

	if status != "" {
		return handleReviewUpdate(status, id, r)
	}

	return api.NotSupported(r)
}

// ConsumerReviewItemRequest get consumer review item by id
func ConsumerReviewItemRequest(r api.APIRequest) (api.APIResponse, error) {
	if strings.ToUpper(r.HTTPMethod) != http.MethodGet {
		return api.NotSupported(r)
	}
	reviewID := r.GetPathParam("id")
	if reviewID != "" {
		return handleConsumerReviewID(reviewID, r)
	}
	return api.NotSupported(r)
}

// ConsumerStatusListRequest Consumer review list by status
func ConsumerStatusListRequest(r api.APIRequest) (api.APIResponse, error) {
	if strings.ToUpper(r.HTTPMethod) != http.MethodGet {
		return api.NotSupported(r)
	}
	limit, _ := r.GetQueryIntParamWithDefault("limit", 30)
	offset, _ := r.GetQueryIntParamWithDefault("offset", 0)
	if status := r.GetPathParam("status"); status != "" {
		return handleConsumerReviewList(status, limit, offset, r)
	}
	return api.NotSupported(r)
}

// ApproveRequest approve a business from  memberReview - docReview - RiskReview
func ApproveRequest(r api.APIRequest) (api.APIResponse, error) {
	if strings.ToUpper(r.HTTPMethod) != http.MethodPost {
		return api.NotSupported(r)
	}
	businessID := r.GetPathParam("businessId")
	if businessID == "" {
		return api.BadRequest(r, errors.New("Missing or invalid params"))
	}

	return handleApprove(businessID, r)
}

// DeclineRequest declines a business from  memberReview - docReview - RiskReview
func DeclineRequest(r api.APIRequest) (api.APIResponse, error) {
	if strings.ToUpper(r.HTTPMethod) != http.MethodPost {
		return api.NotSupported(r)
	}
	businessID := r.GetPathParam("businessId")
	if businessID == "" {
		return api.BadRequest(r, errors.New("Missing or invalid params"))
	}
	return handleDecline(businessID, r)
}
