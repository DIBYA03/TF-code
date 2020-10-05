package document

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/wiseco/core-platform/api"
	docsrv "github.com/wiseco/core-platform/services/document"
	"github.com/wiseco/core-platform/shared"
)

func getURL(r api.APIRequest, documentID shared.BusinessDocumentID, businessID shared.BusinessID) (api.APIResponse, error) {
	url, err := docsrv.NewBusinessDocumentService(r.SourceRequest()).SignedURL(businessID, documentID)

	if err != nil && err.Error() == "not found" {
		return api.NotFound(r)
	} else if err != nil {
		return api.InternalServerError(r, err)
	}

	if url == nil {
		return api.InternalServerError(r, errors.New("Unable to get signed url"))
	}

	URL := struct {
		URL *string `json:"url"`
	}{URL: url}

	b, _ := json.Marshal(URL)
	return api.Success(r, string(b), false)

}

//HandleDocumentSignedAPIRequest handles the document signed url request
func HandleDocumentSignedAPIRequest(request api.APIRequest) (api.APIResponse, error) {
	var method = strings.ToUpper(request.HTTPMethod)
	businessID, err := shared.ParseBusinessID(request.GetPathParam("businessId"))
	if err != nil {
		return api.BadRequestError(request, err)
	}

	documentID, err := shared.ParseBusinessDocumentID(request.GetPathParam("documentId"))
	if err != nil {
		return api.BadRequestError(request, err)
	}

	switch method {
	case http.MethodGet:
		return getURL(request, documentID, businessID)
	default:
		return api.NotSupported(request)
	}
}
