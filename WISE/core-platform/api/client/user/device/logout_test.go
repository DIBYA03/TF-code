package device

import (
	"encoding/json"
	"testing"

	usersrv "github.com/wiseco/core-platform/services/user"
	"github.com/wiseco/core-platform/test"
)

var request = test.TestRequest("post")

func requestBody() string {
	b := usersrv.DeviceLogout{
		DeviceKey: "device_key",
	}
	json, _ := json.Marshal(b)
	return string(json)
}
func TestLogout(t *testing.T) {
	request.Body = requestBody()
	HandleDeviceLogoutAPIRequest(*request)
}
