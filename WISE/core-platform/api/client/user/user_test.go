package user

import (
	"net/http"
	"testing"

	"github.com/wiseco/core-platform/api"
	"github.com/wiseco/core-platform/test"
)

var testRequest = []test.TestResource{
	test.TestResource{http.MethodDelete,
		http.StatusOK,
		"/user"},
}

func runRequestTest(t *testing.T, method, resouce string) api.APIResponse {
	request := test.TestRequest(method)
	request.UserId = "eca9f55a-186e-40c0-b99e-d17c78a5e980"
	request.ResourcePath = resouce
	request.PathParameters = map[string]string{"userId": "eca9f55a-186e-40c0-b99e-d17c78a5e980"}
	resp, err := HandleUserDeleteRequest(*request)
	if err != nil {
		t.Errorf("request for deleting user failed, details: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("request for deteling user failed, expected 200 got:%d", resp.StatusCode)
	}
	return resp
}

func TestRequest(t *testing.T) {
	for _, test := range testRequest {
		resp := runRequestTest(t, test.Method, test.Resource)
		t.Log(resp.Body)
	}
}
