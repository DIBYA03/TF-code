package user

import (
	"encoding/json"
	"testing"

	usersrv "github.com/wiseco/core-platform/services/user"
	"github.com/wiseco/core-platform/test"
)

func TestUpdateNotification(t *testing.T) {
	request := test.TestRequest("patch")
	request.PathParameters["userId"] = "aa41fcc7-acbb-480d-b285-dded8bcdf645"
	request.Body = userNotificationpatchBody(t)
	request.UserId = "aa41fcc7-acbb-480d-b285-dded8bcdf645"
	HandleUserNotification(*request)
}

func userNotificationpatchBody(t *testing.T) string {
	n := usersrv.UserNotificationUpdate{}
	JSON, err := json.Marshal(n)
	if err != nil {
		t.Errorf("Error marshalling %v into json", n)
		return ""
	}
	return string(JSON)
}
