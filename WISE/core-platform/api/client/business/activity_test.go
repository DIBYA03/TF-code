package business

import (
	"fmt"
	"testing"

	"github.com/wiseco/core-platform/test"
)

func TestActivityList(t *testing.T) {
	request := test.TestRequest("GET")

	request.UserID = "08124497-9261-4ca0-a470-b570a7badd58"

	request.PathParameters["businessId"] = "69c9fb64-a084-4621-91ab-e84edbe31292"
	resp, err := HandleActivityRequest(*request)
	if err != nil {
		t.Errorf("getting list failed details:%v", err)
	}
	fmt.Println(resp)
}
func TestGetActivityByID(t *testing.T) {
	request := test.TestRequest("GET")

	request.PathParameters["businessId"] = "1bf0faf1-cb57-4053-a2ae-e693ec99d7a3"
	request.PathParameters["activitId"] = "1bf0faf1-cb57-4053-a2ae-e693ec99d7a3"
	resp, err := HandleActivityRequest(*request)
	if err != nil {
		t.Errorf("getting list failed details:%v", err)
	}
	fmt.Println(resp)
}
