package document

import (
	"encoding/json"
	"errors"

	"github.com/wiseco/core-platform/api"
	"github.com/wiseco/core-platform/services"
	"github.com/wiseco/core-platform/services/document"
	"github.com/wiseco/core-platform/shared"
)

func handleCreateDocument(r api.APIRequest) (api.APIResponse, error) {
	consumerID := r.GetPathParam("consumerId")
	if consumerID == "" {
		return api.BadRequest(r, errors.New("missing or invalid params"))
	}

	conID, err := shared.ParseConsumerID(consumerID)
	if err != nil {
		return api.BadRequest(r, err)
	}

	var body document.ConsumerDocumentCreate
	if err := json.Unmarshal([]byte(r.Body), &body); err != nil {
		return api.BadRequest(r, err)
	}

	doc, err := document.NewConsumerDocumentService(r.SourceRequest()).Create(conID, body)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	resp, err := json.Marshal(doc)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	return api.Success(r, string(resp), false)
}

func handleUpdateDocument(consumerID, documentID string, r api.APIRequest) (api.APIResponse, error) {
	conID, err := shared.ParseConsumerID(consumerID)
	if err != nil {
		return api.BadRequest(r, err)
	}

	docID, err := shared.ParseConsumerDocumentID(documentID)
	if err != nil {
		return api.BadRequest(r, err)
	}

	body := document.ConsumerDocumentUpdate{ID: docID}
	if err := json.Unmarshal([]byte(r.Body), &body); err != nil {
		return api.BadRequest(r, err)
	}

	doc, err := document.NewConsumerDocumentService(r.SourceRequest()).Update(conID, body)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	resp, err := json.Marshal(doc)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	return api.Success(r, string(resp), false)
}

func handleDeleteDocument(consumerID, documentID string, r api.APIRequest) (api.APIResponse, error) {
	conID, err := shared.ParseConsumerID(consumerID)
	if err != nil {
		return api.BadRequest(r, err)
	}

	docID, err := shared.ParseConsumerDocumentID(documentID)
	if err != nil {
		return api.BadRequest(r, err)
	}

	err = document.NewConsumerDocumentService(r.SourceRequest()).Delete(conID, docID)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	return api.Success(r, "", false)
}

func handleDocumentURL(consumerID, documentID string, r api.APIRequest) (api.APIResponse, error) {
	conID, err := shared.ParseConsumerID(consumerID)
	if err != nil {
		return api.BadRequest(r, err)
	}

	docID, err := shared.ParseConsumerDocumentID(documentID)
	if err != nil {
		return api.BadRequest(r, err)
	}

	url, err := document.NewConsumerDocumentService(r.SourceRequest()).SignedURL(conID, docID)
	if err != nil {
		_, isENF := err.(*services.ErrorNotFound)
		if isENF {
			return api.NotFound(r)
		}

		return api.InternalServerError(r, err)
	}

	if url == nil {
		return api.InternalServerError(r, errors.New("Unable to get signed url"))
	}

	URL := struct {
		URL *string `json:"url"`
	}{URL: url}

	resp, _ := json.Marshal(URL)
	return api.Success(r, string(resp), false)

}

func handleDocumentByID(consumerID, documentID string, r api.APIRequest) (api.APIResponse, error) {
	conID, err := shared.ParseConsumerID(consumerID)
	if err != nil {
		return api.BadRequest(r, err)
	}

	docID, err := shared.ParseConsumerDocumentID(documentID)
	if err != nil {
		return api.BadRequest(r, err)
	}

	doc, err := document.NewConsumerDocumentService(r.SourceRequest()).GetByID(conID, docID)
	if err != nil && err == err.(*services.ErrorNotFound) {
		return api.NotFound(r)
	}

	if err != nil {
		return api.InternalServerError(r, err)
	}

	resp, err := json.Marshal(doc)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	return api.Success(r, string(resp), false)
}

func handleDocumentList(consumerID string, limit, offset int, r api.APIRequest) (api.APIResponse, error) {
	conID, err := shared.ParseConsumerID(consumerID)
	if err != nil {
		return api.BadRequest(r, err)
	}

	list, err := document.NewConsumerDocumentService(r.SourceRequest()).List(conID, offset, limit)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	resp, err := json.Marshal(list)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	return api.Success(r, string(resp), false)
}
