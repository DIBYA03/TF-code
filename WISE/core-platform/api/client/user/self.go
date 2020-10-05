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
)

func getSelf(r api.APIRequest) (api.APIResponse, error) {

	user, err := usersrv.NewUserService(r.SourceRequest()).GetById(r.UserID)
	if err != nil {
		return api.NotFoundError(r, err)
	} else if user == nil {
		return api.NotFound(r)
	}

	userJSON, _ := json.Marshal(user)
	return api.Success(r, string(userJSON), false)
}

func HandleUserSelfAPIRequest(r api.APIRequest) (api.APIResponse, error) {

	var method = strings.ToUpper(r.HTTPMethod)

	switch method {
	case http.MethodGet:
		return getSelf(r)
	default:
		return api.NotSupported(r)
	}
}
