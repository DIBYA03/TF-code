/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

package user

import (
	"encoding/json" //"net/http" -- Not being used right now.
	"net/http"
	"strings"

	"github.com/wiseco/core-platform/api"
	usersrv "github.com/wiseco/core-platform/services/user"
	"github.com/wiseco/core-platform/shared"
)

type UserPatchBody = usersrv.UserUpdate

func getUser(r api.APIRequest, id shared.UserID) (api.APIResponse, error) {
	user, err := usersrv.NewUserService(r.SourceRequest()).GetById(id)
	if err != nil || user == nil {
		return api.NotFoundError(r, err)
	}

	userJSON, _ := json.Marshal(user)
	return api.Success(r, string(userJSON), false)
}

func updateUser(r api.APIRequest, id shared.UserID) (api.APIResponse, error) {
	var requestBody UserPatchBody
	if err := json.Unmarshal([]byte(r.Body), &requestBody); err != nil {
		return api.BadRequest(r, err)
	}

	requestBody.ID = id
	user, err := usersrv.NewUserService(r.SourceRequest()).Update(requestBody)
	if err != nil {
		return api.InternalServerError(r, err)
	} else if user == nil {
		return api.NotFound(r)
	}

	userJSON, _ := json.Marshal(user)
	return api.Success(r, string(userJSON), false)
}

func HandleUserAPIRequest(r api.APIRequest) (api.APIResponse, error) {

	var method = strings.ToUpper(r.HTTPMethod)

	userId, err := shared.ParseUserID(r.GetPathParam("userId"))
	if err == nil {
		switch method {
		case http.MethodGet:
			return getUser(r, userId)
		case http.MethodPatch:
			return updateUser(r, userId)
		default:
			return api.NotSupported(r)
		}
	}

	return api.NotSupported(r)
}
