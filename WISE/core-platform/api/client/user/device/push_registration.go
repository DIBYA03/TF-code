package device

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/wiseco/core-platform/api"
	"github.com/wiseco/core-platform/services"
	usersrv "github.com/wiseco/core-platform/services/user"
	"github.com/wiseco/core-platform/shared"
)

func register(r api.APIRequest) (api.APIResponse, error) {
	var body usersrv.PushRegistrationCreate
	if err := json.Unmarshal([]byte(r.Body), &body); err != nil {
		return api.BadRequest(r, err)
	}

	push, err := usersrv.NewDeviceService(r.SourceRequest()).RegisterPush(&body)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	JSON, _ := json.Marshal(push)
	return api.Success(r, string(JSON), false)
}

func get(r api.APIRequest, id shared.UserDeviceID) (api.APIResponse, error) {
	push, err := usersrv.NewDeviceService(r.SourceRequest()).GetByID(id)
	if err != nil {
		_, isENF := err.(*services.ErrorNotFound)
		if isENF {
			return api.NotFound(r)
		}

		return api.InternalServerError(r, err)
	}

	pushJSON, err := json.Marshal(push)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	return api.Success(r, string(pushJSON), false)
}

func unregister(r api.APIRequest, token string) (api.APIResponse, error) {
	if ok, err := usersrv.NewDeviceService(r.SourceRequest()).UnregisterPush(token); !ok {
		return api.InternalServerError(r, err)
	}

	return api.Success(r, "", false)
}

//HandlePushNotification handles device registration for push notifications
func HandlePushNotification(r api.APIRequest) (api.APIResponse, error) {
	id := r.GetPathParam("registrationId")
	method := strings.ToUpper(r.HTTPMethod)
	if id != "" {
		switch method {
		case http.MethodDelete:
			return unregister(r, id)
		case http.MethodGet:
			return get(r, shared.UserDeviceID(id))
		}
		return api.NotSupported(r)
	}

	switch method {
	case http.MethodPost:
		return register(r)
	}

	return api.NotSupported(r)
}
