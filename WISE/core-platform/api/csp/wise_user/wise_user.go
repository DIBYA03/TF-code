package wise_user

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/wiseco/core-platform/api"
	usr "github.com/wiseco/core-platform/services/csp/wise_user"
	"github.com/wiseco/core-platform/shared"
)

//PhoneNumberChangeRequest ..
func PhoneNumberChangeRequest(r api.APIRequest) (api.APIResponse, error) {

	userID := r.GetPathParam("userId")

	if userID == "" {
		return api.NotSupported(r)
	}

	switch strings.ToUpper(r.HTTPMethod) {
	case http.MethodPost:
		return handlePhoneNumberChange(userID, r)
	case http.MethodGet:
		return handleListPhoneNumberChange(userID, r)
	}

	return api.NotSupported(r)
}

func handlePhoneNumberChange(userID string, r api.APIRequest) (api.APIResponse, error) {
	usrID, err := shared.ParseUserID(userID)
	if err != nil {
		return api.BadRequest(r, err)
	}

	var body usr.PhoneChangeRequestCreate
	if err := json.Unmarshal([]byte(r.Body), &body); err != nil {
		return api.BadRequest(r, err)
	}

	body.UserID = usrID

	err = usr.NewWithCognitoSource(r.SourceRequest(), &r.CognitoID).ChangePhone(body)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	return api.Success(r, string(""), false)
}

func handleListPhoneNumberChange(userID string, r api.APIRequest) (api.APIResponse, error) {
	usrID, err := shared.ParseUserID(userID)
	if err != nil {
		return api.BadRequest(r, err)
	}

	list, err := usr.NewWithCognitoSource(r.SourceRequest(), &r.CognitoID).ListPhoneChangeRequest(usrID)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	resp, err := json.Marshal(list)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	return api.Success(r, string(resp), false)
}
