package device

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/wiseco/core-platform/api"
	usersrv "github.com/wiseco/core-platform/services/user"
)

func logout(r api.APIRequest) (api.APIResponse, error) {
	var body usersrv.DeviceLogout
	err := json.Unmarshal([]byte(r.Body), &body)
	if err != nil {
		return api.BadRequest(r, err)
	}

	if err := usersrv.NewDeviceService(r.SourceRequest()).Logout(&body); err != nil {
		return api.InternalServerError(r, err)
	}

	return api.Success(r, "", false)
}

//HandleDeviceLogoutAPIRequest handles device logout
func HandleDeviceLogoutAPIRequest(r api.APIRequest) (api.APIResponse, error) {
	method := strings.ToUpper(r.HTTPMethod)
	switch method {
	case http.MethodPost:
		return logout(r)
	}

	return api.NotSupported(r)
}
