/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

// Package for handling member memberRequest.requests
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

func submitMember(r api.APIRequest, memberID shared.BusinessMemberID, businessID shared.BusinessID) (api.APIResponse, error) {
	member, err := bsrv.NewMemberService(r.SourceRequest()).Submit(memberID, businessID)
	if err != nil {
		return api.BadRequest(r, err)
	} else if member == nil {
		return api.NotFound(r)
	}

	memberJSON, _ := json.Marshal(member)
	return api.Success(r, string(memberJSON), false)
}

//HandleMemberVerificationRequest handle member kyc requests
func HandleMemberSubmissionRequest(r api.APIRequest) (api.APIResponse, error) {
	var method = strings.ToUpper(r.HTTPMethod)
	if r.BusinessID == nil {
		return api.BadRequestError(r, errors.New("Missing header X-Wise-Business-ID"))
	}

	memberID, err := shared.ParseBusinessMemberID(r.GetPathParam("memberId"))
	if err != nil {
		return api.BadRequestError(r, err)
	}

	switch method {
	case http.MethodPost:
		return submitMember(r, memberID, *r.BusinessID)
	default:
		return api.NotSupported(r)
	}
}
