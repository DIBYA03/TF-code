/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

// Package for handling business businessRequest.requests
package business

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/wiseco/core-platform/api"
	bussrv "github.com/wiseco/core-platform/services/business"
	"github.com/wiseco/core-platform/shared"
)

type BusinessPostBody = bussrv.BusinessCreate
type BusinessPatchBody = bussrv.BusinessUpdate

func getBusiness(r api.APIRequest, id shared.BusinessID) (api.APIResponse, error) {
	b, err := bussrv.NewBusinessService(r.SourceRequest()).GetById(id)
	if err != nil {
		return api.BadRequest(r, err)
	}

	businessJSON, err := json.Marshal(b)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	return api.Success(r, string(businessJSON), false)
}

func getBusinesses(r api.APIRequest) (api.APIResponse, error) {
	offset, _ := r.GetQueryIntParamWithDefault("offset", 0)
	limit, _ := r.GetQueryIntParamWithDefault("limit", 10)

	list, err := bussrv.NewBusinessService(r.SourceRequest()).List(offset, limit, r.UserID)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	businessListJSON, err := json.Marshal(list)
	if err != nil {
		return api.InternalServerError(r, err)
	}
	return api.Success(r, string(businessListJSON), false)
}

func createBusiness(r api.APIRequest) (api.APIResponse, error) {
	var requestBody BusinessPostBody
	userID := r.UserID
	err := json.Unmarshal([]byte(r.Body), &requestBody)

	if err != nil {
		return api.BadRequest(r, err)
	}

	requestBody.OwnerID = userID
	id, err := bussrv.NewBusinessService(r.SourceRequest()).Create(&requestBody)

	if err != nil {
		return api.InternalServerError(r, err)
	}

	businessJSON, _ := json.Marshal(id)
	return api.Success(r, string(businessJSON), false)
}

func updateBusiness(r api.APIRequest, id shared.BusinessID) (api.APIResponse, error) {
	var requestBody BusinessPatchBody
	err := json.Unmarshal([]byte(r.Body), &requestBody)
	if err != nil {
		return api.BadRequest(r, err)
	}

	requestBody.ID = id
	b, err := bussrv.NewBusinessService(r.SourceRequest()).Update(requestBody)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	businessJSON, _ := json.Marshal(b)
	return api.Success(r, string(businessJSON), false)
}

//HandleBusinessAPIRequests handles the api request
func HandleBusinessAPIRequests(request api.APIRequest) (api.APIResponse, error) {
	var method = strings.ToUpper(request.HTTPMethod)

	if request.GetPathParam("businessId") != "" {
		businessID, err := shared.ParseBusinessID(request.GetPathParam("businessId"))
		if err != nil {
			return api.BadRequest(request, err)
		}

		switch strings.ToUpper(request.HTTPMethod) {
		case http.MethodGet:
			return getBusiness(request, businessID)
		case http.MethodPatch:
			return updateBusiness(request, businessID)
		default:
			return api.NotSupported(request)
		}

	}

	switch method {
	case http.MethodGet:
		return getBusinesses(request)
	case http.MethodPost:
		return createBusiness(request)
	default:
		return api.NotSupported(request)
	}
}
