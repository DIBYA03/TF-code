package business

import (
	"encoding/json"
	"errors"

	"github.com/wiseco/core-platform/api"
	"github.com/wiseco/core-platform/services"
	csp "github.com/wiseco/core-platform/services/csp/document"
	"github.com/wiseco/core-platform/services/csp/review"
	"github.com/wiseco/core-platform/shared"
)

var ds = csp.NewDocumentService()

func handleDocumentList(businessID string, limit, offset int, r api.APIRequest) (api.APIResponse, error) {
	bID, err := shared.ParseBusinessID(businessID)
	if err != nil {
		return api.BadRequest(r, err)
	}

	list, err := ds.List(bID, limit, offset)
	if err != nil {
		return api.InternalServerError(r, err)
	}
	resp, err := json.Marshal(list)
	if err != nil {
		return api.InternalServerError(r, err)
	}
	return api.Success(r, string(resp), false)
}

func handleDocumentID(businessID, documentID string, r api.APIRequest) (api.APIResponse, error) {
	bID, err := shared.ParseBusinessID(businessID)
	if err != nil {
		return api.BadRequest(r, err)
	}

	docID, err := shared.ParseBusinessDocumentID(documentID)
	if err != nil {
		return api.BadRequest(r, err)
	}

	doc, err := ds.GetByID(bID, docID)
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

func handleDocumentCreate(businessID string, r api.APIRequest) (api.APIResponse, error) {
	bID, err := shared.ParseBusinessID(businessID)
	if err != nil {
		return api.BadRequest(r, err)
	}

	var body csp.BusinessDocumentCreate
	if err := json.Unmarshal([]byte(r.Body), &body); err != nil {
		return api.BadRequest(r, err)
	}
	body.BusinessID = &bID
	doc, err := ds.Create(bID, body)
	if err != nil {
		return api.InternalServerError(r, err)
	}
	resp, err := json.Marshal(doc)
	if err != nil {
		return api.InternalServerError(r, err)
	}
	return api.Success(r, string(resp), false)
}

func handleDocumentDelete(businessID, documentID string, r api.APIRequest) (api.APIResponse, error) {
	bID, err := shared.ParseBusinessID(businessID)
	if err != nil {
		return api.BadRequest(r, err)
	}

	docID, err := shared.ParseBusinessDocumentID(documentID)
	if err != nil {
		return api.BadRequest(r, err)
	}

	err = ds.Delete(bID, docID)
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

func handleDocumentUpdate(businessID, documentID string, r api.APIRequest) (api.APIResponse, error) {
	bID, err := shared.ParseBusinessID(businessID)
	if err != nil {
		return api.BadRequest(r, err)
	}

	docID, err := shared.ParseBusinessDocumentID(documentID)
	if err != nil {
		return api.BadRequest(r, err)
	}

	var updates csp.BusinessDocumentUpdate

	json.Unmarshal([]byte(r.Body), &updates)

	doc, err := ds.Update(bID, docID, updates)
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

func handleDocumentURL(businessID, documentID string, r api.APIRequest) (api.APIResponse, error) {
	bID, err := shared.ParseBusinessID(businessID)
	if err != nil {
		return api.BadRequest(r, err)
	}

	docID, err := shared.ParseBusinessDocumentID(documentID)
	if err != nil {
		return api.BadRequest(r, err)
	}

	url, err := ds.URL(bID, docID)
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

func handleSetAsFormation(businessID, documentID string, r api.APIRequest) (api.APIResponse, error) {
	bID, err := shared.ParseBusinessID(businessID)
	if err != nil {
		return api.BadRequest(r, err)
	}

	docID, err := shared.ParseBusinessDocumentID(documentID)
	if err != nil {
		return api.BadRequest(r, err)
	}

	if err := ds.SetAsFormation(bID, docID, false); err != nil {
		return api.InternalServerError(r, err)
	}

	status := struct {
		Message string `json:"message"`
	}{Message: "Success"}
	resp, _ := json.Marshal(status)
	return api.Success(r, string(resp), false)

}

func handleRemoveFormation(businessID, documentID string, r api.APIRequest) (api.APIResponse, error) {
	bID, err := shared.ParseBusinessID(businessID)
	if err != nil {
		return api.BadRequest(r, err)
	}

	docID, err := shared.ParseBusinessDocumentID(documentID)
	if err != nil {
		return api.BadRequest(r, err)
	}

	if err := ds.SetAsFormation(bID, docID, true); err != nil {
		return api.InternalServerError(r, err)
	}
	status := struct {
		Message string `json:"message"`
	}{Message: "Success"}
	resp, _ := json.Marshal(status)
	return api.Success(r, string(resp), false)

}

func handleDocumentStatus(businessID, documentID string, r api.APIRequest) (api.APIResponse, error) {

	docID, err := shared.ParseBusinessDocumentID(documentID)
	if err != nil {
		return api.BadRequest(r, err)
	}

	status, err := csp.NewCSPBusinessService().CSPDocumentByDocumentID(docID)
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

func handleSubmitDocument(businessID, documentID string, r api.APIRequest) (api.APIResponse, error) {
	bID, err := shared.ParseBusinessID(businessID)
	if err != nil {
		return api.BadRequest(r, err)
	}

	docID, err := shared.ParseBusinessDocumentID(documentID)
	if err != nil {
		return api.BadRequest(r, err)
	}

	if err := review.NewUploader(r.SourceRequest()).BusinessSingleUpload(bID, docID); err != nil {
		return api.InternalServerError(r, err)
	}

	status := struct {
		Message string `json:"message"`
	}{Message: "Success"}
	resp, _ := json.Marshal(status)
	return api.Success(r, string(resp), false)
}

func handleReUploadDocument(businessID string, r api.APIRequest) (api.APIResponse, error) {
	bID, err := shared.ParseBusinessID(businessID)
	if err != nil {
		return api.BadRequest(r, err)
	}

	if err := review.NewUploader(r.SourceRequest()).BusinessSingleReUpload(bID); err != nil {
		return api.InternalServerError(r, err)
	}

	status := struct {
		Message string `json:"message"`
	}{Message: "Success"}

	resp, _ := json.Marshal(status)
	return api.Success(r, string(resp), false)

}
