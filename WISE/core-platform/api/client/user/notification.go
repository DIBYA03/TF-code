package user

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/pkg/errors"
	"github.com/wiseco/core-platform/api"
	usersrv "github.com/wiseco/core-platform/services/user"
)

func updateNotification(r api.APIRequest) (api.APIResponse, error) {
	var notification usersrv.UserNotificationUpdate
	if err := json.Unmarshal([]byte(r.Body), &notification); err != nil {
		return api.BadRequest(r, err)
	}

	n, err := usersrv.NewUserService(r.SourceRequest()).UpdateNotification(&notification)
	if err != nil {
		return api.InternalServerError(r, errors.Wrap(err, "error updating notifications"))
	}

	notif, err := json.Marshal(n)
	if err != nil {
		return api.InternalServerError(r, errors.Wrap(err, "error encoding"))
	}

	return api.Success(r, string(notif), false)
}

//HandleUserNotification handles user notification updates
func HandleUserNotification(r api.APIRequest) (api.APIResponse, error) {
	if strings.ToUpper(r.HTTPMethod) == http.MethodPatch {
		return updateNotification(r)
	}

	return api.NotSupported(r)
}
