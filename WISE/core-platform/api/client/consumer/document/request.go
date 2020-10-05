package document

import (
	"errors"
	"net/http"
	"strings"

	"github.com/wiseco/core-platform/api"
)

//DocumentRequest user document create request
func DocumentRequest(r api.APIRequest) (api.APIResponse, error) {
	switch strings.ToUpper(r.HTTPMethod) {
	case http.MethodPost:
		return handleCreateDocument(r)
	case http.MethodGet:
		return getDocumentRequest(r)
	case http.MethodPatch:
		return updateDocumentRequest(r)
	case http.MethodDelete:
		return deleteDocumentRequest(r)
	}
	return api.NotSupported(r)
}

func updateDocumentRequest(r api.APIRequest) (api.APIResponse, error) {
	consumerID := r.GetPathParam("consumerId")
	documentID := r.GetPathParam("documentId")
	if consumerID == "" || documentID == "" {
		return api.BadRequest(r, errors.New("missing or invalid params"))
	}
	return handleUpdateDocument(consumerID, documentID, r)
}

func deleteDocumentRequest(r api.APIRequest) (api.APIResponse, error) {
	consumerID := r.GetPathParam("consumerId")
	documentID := r.GetPathParam("documentId")
	if consumerID == "" || documentID == "" {
		return api.BadRequest(r, errors.New("missing or invalid params"))
	}
	return handleDeleteDocument(consumerID, documentID, r)
}

func getDocumentRequest(r api.APIRequest) (api.APIResponse, error) {
	consumerID := r.GetPathParam("consumerId")
	documentID := r.GetPathParam("documentId")
	if consumerID == "" {
		return api.BadRequest(r, errors.New("missing or invalid params"))
	}
	if documentID == "" {
		limit, _ := r.GetQueryIntParamWithDefault("limit", 30)
		offset, _ := r.GetQueryIntParamWithDefault("offset", 0)
		return handleDocumentList(consumerID, limit, offset, r)
	}
	return handleDocumentByID(consumerID, documentID, r)
}

//DocumentURLRequest user document URL request
func DocumentURLRequest(r api.APIRequest) (api.APIResponse, error) {
	consumerID := r.GetPathParam("consumerId")
	documentID := r.GetPathParam("documentId")
	if consumerID == "" || documentID == "" {
		return api.BadRequest(r, errors.New("missing or invalid params"))
	}
	if strings.ToUpper(r.HTTPMethod) != http.MethodGet {
		return api.NotSupported(r)
	}
	return handleDocumentURL(consumerID, documentID, r)
}
