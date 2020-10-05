/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

package user

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/wiseco/core-platform/api"
	usersrv "github.com/wiseco/core-platform/services/user"
)

type PartnerVerificationBody = usersrv.PartnerVerification

func verifyPartnerCode(r api.APIRequest) (api.APIResponse, error) {
	var requestBody PartnerVerificationBody
	userId := r.UserID
	err := json.Unmarshal([]byte(r.Body), &requestBody)

	if err != nil {
		return api.BadRequest(r, err)
	}

	requestBody.UserID = userId

	id, err := usersrv.NewPartnerService(r.SourceRequest()).Verify(&requestBody)

	if err != nil {
		return api.InternalServerError(r, err)
	}

	userJSON, _ := json.Marshal(id)
	return api.Success(r, string(userJSON), false)
}

func HandlePartnerCodeVerificationRequest(r api.APIRequest) (api.APIResponse, error) {

	var method = strings.ToUpper(r.HTTPMethod)

	switch method {
	case http.MethodPost:
		return verifyPartnerCode(r)
	default:
		return api.NotSupported(r)
	}
}
