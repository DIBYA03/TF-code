/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

package business

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/wiseco/core-platform/api"
	"github.com/wiseco/core-platform/services/signature"
	"github.com/wiseco/core-platform/shared"
)

type SignaturePOSTBody struct {
	TemplateType string
}

func createSignatureRequest(r api.APIRequest, id shared.BusinessID) (api.APIResponse, error) {
	var requestBody SignaturePOSTBody
	err := json.Unmarshal([]byte(r.Body), &requestBody)
	if err != nil {
		return api.BadRequest(r, err)
	}

	b, err := signature.NewSignatureService(r.SourceRequest()).Create(id, signature.SignatureRequestTemplate(requestBody.TemplateType))
	if err != nil {
		return api.InternalServerError(r, err)
	}

	if b == nil {
		return api.InternalServerError(r, err)
	}

	bJSON, _ := json.Marshal(b)
	return api.Success(r, string(bJSON), false)
}

func getSignatureRequest(r api.APIRequest, id shared.BusinessID) (api.APIResponse, error) {
	templateType := r.GetQueryParam("templateType")
	if len(templateType) == 0 {
		return api.BadRequestError(r, errors.New("query param missing"))
	}

	b, err := signature.NewSignatureService(r.SourceRequest()).GetByBusinessID(id, signature.SignatureRequestTemplate(templateType))
	if err != nil {
		return api.InternalServerError(r, err)
	}

	if b == nil {
		return api.Success(r, "", false)
	}

	bJSON, _ := json.Marshal(b)
	return api.Success(r, string(bJSON), false)
}

// HandleBusinessSignatureRequest handle user kyc requests
func HandleBusinessSignatureRequest(r api.APIRequest) (api.APIResponse, error) {
	var method = strings.ToUpper(r.HTTPMethod)

	businessID, err := shared.ParseBusinessID(r.GetPathParam("businessId"))
	if err != nil {
		return api.BadRequestError(r, err)
	}

	switch method {
	case http.MethodPost:
		return createSignatureRequest(r, businessID)
	case http.MethodGet:
		return getSignatureRequest(r, businessID)
	default:
		return api.NotSupported(r)
	}
}
