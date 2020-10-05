package device

import (
	"encoding/json"
	"testing"

	usersrv "github.com/wiseco/core-platform/services/user"
	"github.com/wiseco/core-platform/test"
)

var postRequest = test.TestRequest("post")

func bodyPostRequest() string {
	body := usersrv.PushRegistrationCreate{
		Token:      "test-token", // add test token
		TokenType:  "fcm",
		DeviceType: "ios",
		Language:   "en-US",
	}

	b, _ := json.Marshal(body)
	return string(b)
}

func TestRegister(t *testing.T) {
	postRequest.Body = bodyPostRequest()
	postRequest.UserID = "aa41fcc7-acbb-480d-b285-dded8bcdf645"
	_, err := HandlePushNotification(*postRequest)
	if err != nil {
		t.Errorf("registering the token fail details: %v", err)
	}
}

func TestGetByID(t *testing.T) {
}

func TestUnregister(t *testing.T) {
	//TODO: to be implemented with device logout
}
