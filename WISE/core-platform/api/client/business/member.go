/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

// Package for handling business members businessRequest.businessRequest

package business

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/wiseco/core-platform/api"
	bsrv "github.com/wiseco/core-platform/services/business"
	"github.com/wiseco/core-platform/shared"
)

type businessMemberPostBody = bsrv.BusinessMemberCreate
type businessMemberPatchBody = bsrv.BusinessMemberUpdate

func deactivateMember(r api.APIRequest, memberId shared.BusinessMemberID, businessID shared.BusinessID) (api.APIResponse, error) {
	err := bsrv.NewMemberService(r.SourceRequest()).Deactivate(memberId, businessID)
	if err != nil {
		return api.BadRequest(r, err)
	}

	return api.Success(r, "", false)
}

func getMember(r api.APIRequest, memberId shared.BusinessMemberID, businessID shared.BusinessID) (api.APIResponse, error) {
	b, err := bsrv.NewMemberService(r.SourceRequest()).GetById(memberId, businessID)
	if err != nil {
		return api.BadRequest(r, err)
	}

	memberJSON, err := json.Marshal(b)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	return api.Success(r, string(memberJSON), false)
}

func getMembers(r api.APIRequest, id shared.BusinessID) (api.APIResponse, error) {
	offset, _ := r.GetQueryIntParamWithDefault("offset", 0)
	limit, _ := r.GetQueryIntParamWithDefault("limit", 10)

	members, err := bsrv.NewMemberService(r.SourceRequest()).List(offset, limit, id)
	if err != nil {
		return api.NotFound(r)
	}

	memberListJSON, err := json.Marshal(members)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	return api.Success(r, string(memberListJSON), false)
}

func createMember(r api.APIRequest, id shared.BusinessID) (api.APIResponse, error) {
	var requestBody businessMemberPostBody
	err := json.Unmarshal([]byte(r.Body), &requestBody)
	if err != nil {
		return api.BadRequest(r, err)
	}

	requestBody.BusinessID = id

	member, err := bsrv.NewMemberService(r.SourceRequest()).Create(&requestBody)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	memberJSON, err := json.Marshal(member)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	return api.Success(r, string(memberJSON), false)
}

func updateMember(r api.APIRequest, id shared.BusinessMemberID, businessID shared.BusinessID) (api.APIResponse, error) {
	var requestBody businessMemberPatchBody
	err := json.Unmarshal([]byte(r.Body), &requestBody)
	if err != nil {
		return api.BadRequest(r, err)
	}

	requestBody.ID = id

	member, err := bsrv.NewMemberService(r.SourceRequest()).Update(id, businessID, &requestBody)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	memberJSON, err := json.Marshal(member)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	return api.Success(r, string(memberJSON), false)
}

//HandleMemberAPIRequests handles the api request
func HandleMemberAPIRequests(request api.APIRequest) (api.APIResponse, error) {
	var method = strings.ToUpper(request.HTTPMethod)

	businessID, err := shared.ParseBusinessID(request.GetPathParam("businessId"))
	if err != nil {
		return api.BadRequestError(request, err)
	}

	mID := request.GetPathParam("memberId")
	if len(mID) > 0 {
		memberID, err := shared.ParseBusinessMemberID(mID)
		if err != nil {
			return api.BadRequestError(request, err)
		}

		switch method {
		case http.MethodGet:
			return getMember(request, memberID, businessID)
		case http.MethodPatch:
			return updateMember(request, memberID, businessID)
		case http.MethodDelete:
			return deactivateMember(request, memberID, businessID)
		default:
			return api.NotSupported(request)
		}
	}

	switch method {
	case http.MethodGet:
		return getMembers(request, businessID)
	case http.MethodPost:
		return createMember(request, businessID)
	default:
		return api.NotSupported(request)
	}
}
