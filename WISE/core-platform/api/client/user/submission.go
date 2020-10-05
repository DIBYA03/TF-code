/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

// Package for handling user userRequest.requests
package user

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/wiseco/core-platform/api"
	"github.com/wiseco/core-platform/services/user"
)

func submitUser(r api.APIRequest) (api.APIResponse, error) {
	user, err := user.NewUserService(r.SourceRequest()).Submit(r.UserID)
	if err != nil {
		return api.BadRequest(r, err)
	} else if user == nil {
		return api.NotFound(r)
	}

	userJSON, _ := json.Marshal(user)
	return api.Success(r, string(userJSON), false)
}

//HandleUserVerificationRequest handle user kyc requests
func HandleUserSubmissionRequest(r api.APIRequest) (api.APIResponse, error) {
	var method = strings.ToUpper(r.HTTPMethod)
	switch method {
	case http.MethodPost:
		return submitUser(r)
	default:
		return api.NotSupported(r)
	}
}
