package consumer

import (
	"net/http"
	"strings"

	"github.com/wiseco/core-platform/api"
)

//Request consumer request
func Request(r api.APIRequest) (api.APIResponse, error) {
	consumerID := r.GetPathParam("consumerId")
	if consumerID != "" {
		switch strings.ToUpper(r.HTTPMethod) {
		case http.MethodGet:
			return handleConsumerID(consumerID, r)
		case http.MethodPatch:
			return handleConsumerUpdate(consumerID, r)
		default:
			return api.NotSupported(r)
		}

	}
	return handleConsumerList(r)
}

//DocumentRequest ..
func DocumentRequest(r api.APIRequest) (api.APIResponse, error) {

	consumerID := r.GetPathParam("consumerId")
	documentID := r.GetPathParam("documentId")

	limit, _ := r.GetQueryIntParamWithDefault("limit", 10)
	offet, _ := r.GetQueryIntParamWithDefault("offset", 0)

	if consumerID != "" && documentID != "" {
		switch strings.ToUpper(r.HTTPMethod) {
		case http.MethodGet:
			return handleDocumentID(consumerID, documentID, r)
		case http.MethodPatch:
			return handleDocumentUpdate(consumerID, documentID, r)
		case http.MethodDelete:
			return handleDocumentDelete(consumerID, documentID, r)
		}
	}

	if consumerID != "" {
		switch strings.ToUpper(r.HTTPMethod) {
		case http.MethodPost:
			return handleDocumentCreate(consumerID, r)
		case http.MethodGet:
			return handleDocumentList(consumerID, limit, offet, r)
		}
	}

	return api.NotSupported(r)
}

// DocumentURLRequest document url
func DocumentURLRequest(r api.APIRequest) (api.APIResponse, error) {
	if strings.ToUpper(r.HTTPMethod) != http.MethodGet {
		return api.NotSupported(r)
	}

	consumerID := r.GetPathParam("consumerId")
	documentID := r.GetPathParam("documentId")

	if consumerID != "" && documentID != "" {
		return handleDocumentURL(consumerID, documentID, r)
	}

	return api.NotSupported(r)
}

// DocumentStatusRequest document status
func DocumentStatusRequest(r api.APIRequest) (api.APIResponse, error) {
	if strings.ToUpper(r.HTTPMethod) != http.MethodGet {
		return api.NotSupported(r)
	}

	consumerID := r.GetPathParam("consumerId")
	documentID := r.GetPathParam("documentId")
	if consumerID != "" && documentID != "" {
		return handleDocumentStatus(consumerID, documentID, r)
	}

	return api.NotSupported(r)
}

// SendDocumentRequest sends a single document to BBVA, this requests post a sqs message to the document queue
func SendDocumentRequest(r api.APIRequest) (api.APIResponse, error) {
	if strings.ToUpper(r.HTTPMethod) != http.MethodPost {
		return api.NotSupported(r)
	}

	consumerID := r.GetPathParam("consumerId")
	documentID := r.GetPathParam("documentId")

	if consumerID != "" && documentID != "" {
		return handleSendDocument(consumerID, documentID, r)
	}

	return api.NotSupported(r)
}

// CheckStatusRequest checks kyc status on bbva side
func CheckStatusRequest(r api.APIRequest) (api.APIResponse, error) {
	if strings.ToUpper(r.HTTPMethod) != http.MethodGet {
		return api.NotSupported(r)
	}
	consumerID := r.GetPathParam("consumerId")
	if consumerID != "" {
		return handleCheckStatus(consumerID, r)
	}

	return api.NotSupported(r)
}

// ReviewItemRequest ...
func ReviewItemRequest(r api.APIRequest) (api.APIResponse, error) {
	if strings.ToUpper(r.HTTPMethod) != http.MethodGet {
		return api.NotSupported(r)
	}
	consumerID := r.GetPathParam("consumerId")
	if consumerID != "" {
		return handleConsumerCSPConsumer(consumerID, r)
	}

	return api.NotSupported(r)
}

// VerificationRequest verification to BBVA for any consumer
func VerificationRequest(r api.APIRequest) (api.APIResponse, error) {
	consumerID := r.GetPathParam("consumerId")
	if consumerID == "" {
		return api.NotSupported(r)
	}

	switch strings.ToUpper(r.HTTPMethod) {
	case http.MethodGet:
		return handleGetVerification(consumerID, r)
	case http.MethodPost:
		return handleVerification(consumerID, r)
	default:
		return api.NotSupported(r)
	}
}

func ListStatesRequest(r api.APIRequest) (api.APIResponse, error) {
	consumerID := r.GetPathParam("consumerId")
	if consumerID == "" {
		return api.NotSupported(r)
	}
	switch strings.ToUpper(r.HTTPMethod) {
	case http.MethodGet:
		return handleState(consumerID, r)
	}
	return api.NotSupported(r)
}

// ReUploadDocumentRequest ..
func ReUploadDocumentRequest(r api.APIRequest) (api.APIResponse, error) {
	consumerID := r.GetPathParam("consumerId")
	if consumerID == "" {
		return api.NotSupported(r)
	}
	switch strings.ToUpper(r.HTTPMethod) {
	case http.MethodPost:
		return handleReUploadDocument(consumerID, r)
	default:
		return api.NotSupported(r)
	}
}
