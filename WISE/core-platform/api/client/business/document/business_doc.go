/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

// Package for handling business document
package document

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/wiseco/core-platform/api"
	docsrv "github.com/wiseco/core-platform/services/document"
	"github.com/wiseco/core-platform/shared"
)

type businessDocumentPostBody = docsrv.BusinessDocumentCreate
type businessDocumentPatchBody = docsrv.BusinessDocumentUpdate

func getDocument(r api.APIRequest, documentID shared.BusinessDocumentID, businessID shared.BusinessID) (api.APIResponse, error) {
	doc, err := docsrv.NewBusinessDocumentService(r.SourceRequest()).GetByID(documentID, businessID)
	if err != nil {
		return api.NotFound(r)
	}

	documentJSON, err := json.Marshal(doc)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	return api.Success(r, string(documentJSON), false)
}

func getDocuments(r api.APIRequest, id shared.BusinessID) (api.APIResponse, error) {
	offset, _ := r.GetQueryIntParamWithDefault("offset", 0)
	limit, _ := r.GetQueryIntParamWithDefault("limit", 10)
	docs, err := docsrv.NewBusinessDocumentService(r.SourceRequest()).List(id, "", limit, offset)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	documentListJSON, err := json.Marshal(docs)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	return api.Success(r, string(documentListJSON), false)
}

func createDocument(r api.APIRequest, id shared.BusinessID) (api.APIResponse, error) {
	var requestBody businessDocumentPostBody
	err := json.Unmarshal([]byte(r.Body), &requestBody)
	if err != nil {
		return api.BadRequest(r, err)
	}

	requestBody.CreatedUserID = r.UserID
	requestBody.BusinessID = id
	doc, err := docsrv.NewBusinessDocumentService(r.SourceRequest()).Create(requestBody)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	documentJSON, err := json.Marshal(doc)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	return api.Success(r, string(documentJSON), false)
}

func updateDocument(r api.APIRequest, documentID shared.BusinessDocumentID, businessID shared.BusinessID) (api.APIResponse, error) {
	var requestBody businessDocumentPatchBody
	err := json.Unmarshal([]byte(r.Body), &requestBody)
	if err != nil {
		return api.BadRequest(r, err)
	}

	requestBody.ID = documentID
	requestBody.BusinessID = businessID
	doc, err := docsrv.NewBusinessDocumentService(r.SourceRequest()).Update(requestBody)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	documentJSON, err := json.Marshal(doc)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	return api.Success(r, string(documentJSON), false)
}

func deleteDocument(r api.APIRequest, id shared.BusinessDocumentID) (api.APIResponse, error) {
	// TODO: Hook up to platform methods
	return api.Success(r, "", false)
}

//HandleDocumentAPIRequests handle methods
func HandleDocumentAPIRequests(request api.APIRequest) (api.APIResponse, error) {
	var method = strings.ToUpper(request.HTTPMethod)
	businessID, err := shared.ParseBusinessID(request.GetPathParam("businessId"))
	if err != nil {
		return api.BadRequestError(request, err)
	}

	// documents level path
	documentID := request.GetPathParam("documentId")
	if documentID != "" {
		docID, err := shared.ParseBusinessDocumentID(documentID)
		if err != nil {
			return api.BadRequestError(request, err)
		}

		switch method {
		case http.MethodGet:
			return getDocument(request, docID, businessID)
		case http.MethodPatch:
			return updateDocument(request, docID, businessID)
		case http.MethodDelete:
			return deleteDocument(request, docID)
		default:
			return api.NotSupported(request)
		}
	}
	//business level path
	switch method {
	case http.MethodGet:
		return getDocuments(request, businessID)
	case http.MethodPost:
		return createDocument(request, businessID)
	default:
		return api.NotSupported(request)
	}
}
