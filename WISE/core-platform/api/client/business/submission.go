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
	bsrv "github.com/wiseco/core-platform/services/business"
	"github.com/wiseco/core-platform/shared"
)

func submitBusiness(r api.APIRequest, id shared.BusinessID) (api.APIResponse, error) {
	b, err := bsrv.NewBusinessService(r.SourceRequest()).Submit(id)
	if err != nil {
		return api.InternalServerError(r, err)
	}
	if b == nil {
		return api.InternalServerError(r, err)
	}

	bJSON, _ := json.Marshal(b)
	return api.Success(r, string(bJSON), false)
}

// HandleBusinessSubmissionRequest handle user kyc requests
func HandleBusinessSubmissionRequest(r api.APIRequest) (api.APIResponse, error) {
	var method = strings.ToUpper(r.HTTPMethod)
	if r.BusinessID == nil {
		return api.BadRequestError(r, errors.New("Missing header X-Wise-Business-ID"))
	}

	switch method {
	case http.MethodPost:
		return submitBusiness(r, *r.BusinessID)
	default:
		return api.NotSupported(r)
	}
}
