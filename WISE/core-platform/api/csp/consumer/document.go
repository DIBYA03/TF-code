package consumer

import (
	"encoding/json"
	"errors"

	"github.com/wiseco/core-platform/api"
	"github.com/wiseco/core-platform/services"
	"github.com/wiseco/core-platform/services/csp/document"
	"github.com/wiseco/core-platform/services/csp/review"
	"github.com/wiseco/core-platform/shared"
)

func handleDocumentList(consumerID string, limit, offset int, r api.APIRequest) (api.APIResponse, error) {
	conID, err := shared.ParseConsumerID(consumerID)
	if err != nil {
		return api.BadRequest(r, err)
	}

	list, err := document.NewConsumerDocumentService().List(conID, offset, limit)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	resp, err := json.Marshal(list)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	return api.Success(r, string(resp), false)

}

func handleDocumentID(consumerID, documentID string, r api.APIRequest) (api.APIResponse, error) {
	_, err := shared.ParseConsumerID(consumerID)
	if err != nil {
		return api.BadRequest(r, err)
	}

	docID, err := shared.ParseConsumerDocumentID(documentID)
	if err != nil {
		return api.BadRequest(r, err)
	}

	doc, err := document.NewConsumerDocumentService().GetByID(docID)
	if err != nil {
		_, isENF := err.(*services.ErrorNotFound)
		if isENF {
			return api.NotFound(r)
		}

		return api.InternalServerError(r, err)
	}

	resp, err := json.Marshal(doc)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	return api.Success(r, string(resp), false)

}

func handleDocumentCreate(consumerID string, r api.APIRequest) (api.APIResponse, error) {
	conID, err := shared.ParseConsumerID(consumerID)
	if err != nil {
		return api.BadRequest(r, err)
	}

	var body document.ConsumerDocumentCreate
	if err := json.Unmarshal([]byte(r.Body), &body); err != nil {
		return api.BadRequest(r, err)
	}

	doc, err := document.NewConsumerDocumentService().Create(conID, body)
	if err != nil {
		return api.InternalServerError(r, err)
	}
	resp, err := json.Marshal(doc)
	if err != nil {
		return api.InternalServerError(r, err)
	}
	return api.Success(r, string(resp), false)
}

func handleDocumentUpdate(consumerID, documentID string, r api.APIRequest) (api.APIResponse, error) {
	_, err := shared.ParseConsumerID(consumerID)
	if err != nil {
		return api.BadRequest(r, err)
	}

	docID, err := shared.ParseConsumerDocumentID(documentID)
	if err != nil {
		return api.BadRequest(r, err)
	}

	var updates document.ConsumerDocumentUpdate
	json.Unmarshal([]byte(r.Body), &updates)
	doc, err := document.NewConsumerDocumentService().Update(docID, updates)
	if err != nil {
		_, isENF := err.(*services.ErrorNotFound)
		if isENF {
			return api.NotFound(r)
		}

		return api.InternalServerError(r, err)
	}

	resp, err := json.Marshal(doc)
	if err != nil {
		return api.InternalServerError(r, err)
	}
	return api.Success(r, string(resp), false)
}

func handleDocumentDelete(consumerID, documentID string, r api.APIRequest) (api.APIResponse, error) {
	_, err := shared.ParseConsumerID(consumerID)
	if err != nil {
		return api.BadRequest(r, err)
	}

	_, err = shared.ParseConsumerDocumentID(documentID)
	if err != nil {
		_, isENF := err.(*services.ErrorNotFound)
		if isENF {
			return api.NotFound(r)
		}

		return api.InternalServerError(r, err)
	}

	status := struct {
		Message string `json:"message"`
	}{Message: "Success"}

	resp, _ := json.Marshal(status)

	return api.Success(r, string(resp), false)
}

func handleDocumentURL(consumerID, documentID string, r api.APIRequest) (api.APIResponse, error) {
	_, err := shared.ParseConsumerID(consumerID)
	if err != nil {
		return api.BadRequest(r, err)
	}

	docID, err := shared.ParseConsumerDocumentID(documentID)
	if err != nil {
		return api.BadRequest(r, err)
	}

	url, err := document.NewConsumerDocumentService().SignedURL(docID)
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

func handleDocumentStatus(consumerID, documentID string, r api.APIRequest) (api.APIResponse, error) {
	conID, err := shared.ParseConsumerID(consumerID)
	if err != nil {
		return api.BadRequest(r, err)
	}

	docID, err := shared.ParseConsumerDocumentID(documentID)
	if err != nil {
		return api.BadRequest(r, err)
	}

	status, err := document.NewConsumerDocumentService().Status(conID, docID)
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

func handleSendDocument(consumerID, documentID string, r api.APIRequest) (api.APIResponse, error) {
	id, err := shared.ParseConsumerID(consumerID)
	if err != nil {
		return api.BadRequest(r, err)
	}

	docID, err := shared.ParseConsumerDocumentID(documentID)
	if err != nil {
		return api.BadRequest(r, err)
	}
	if err := review.NewUploader(r.SourceRequest()).ConsumerSingleUpload(id, docID); err != nil {
		return api.InternalServerError(r, err)
	}

	status := struct {
		Message string `json:"message"`
	}{Message: "Success"}

	resp, _ := json.Marshal(status)
	return api.Success(r, string(resp), false)

}

func handleReUploadDocument(consumerID string, r api.APIRequest) (api.APIResponse, error) {
	id, err := shared.ParseConsumerID(consumerID)
	if err != nil {
		return api.BadRequest(r, err)
	}

	if err := review.NewUploader(r.SourceRequest()).ConsumerSingleReUpload(id); err != nil {
		return api.InternalServerError(r, err)
	}

	status := struct {
		Message string `json:"message"`
	}{Message: "Success"}

	resp, _ := json.Marshal(status)
	return api.Success(r, string(resp), false)

}
