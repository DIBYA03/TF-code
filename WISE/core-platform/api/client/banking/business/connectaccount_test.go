/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

// Package for handling linked account api requests
package business

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/wiseco/core-platform/api"
	"github.com/wiseco/core-platform/test"
)

type linkedAccountTestRequests struct {
	Method     string
	StatusCode int
	Resource   string
}

var linkedAccountTests = []linkedAccountTestRequests{
	linkedAccountTestRequests{http.MethodPost, http.StatusOK, ""},
	//linkedAccountTestRequests{http.MethodGet, http.StatusOK, ""},
}

func testConnectBankAccount() string {
	testBody := LinkAccountBody{
		PublicToken: "public-sandbox-6cb08832-d6dc-463c-b81a-e3989c941120",
	}

	b, _ := json.Marshal(testBody)

	return string(b)
}

func runConnectJSONRequestTest(t *testing.T, resource, method string) api.APIResponse {

	request := test.TestRequest(method)

	request.PathParameters = map[string]string{"businessId": "17689405-8710-43d9-ad24-5ce7a65f188c"}

	switch method {
	case http.MethodPost:
		request.Body = testConnectBankAccount()
	default:
		break
	}

	request.UserID = "604123ef-9090-4636-bb39-199197533096"

	resp, err := HandleConnectAccountRequest(*request)
	if err != nil {
		t.Errorf("error handling request request for method %s details: %v", method, err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("error handling request for Method: %s expecting 200,received: %d", method, resp.StatusCode)
		return resp
	}

	return resp
}

func TestHandleConnectAccountRequests(t *testing.T) {

	for _, test := range linkedAccountTests {
		resp := runConnectJSONRequestTest(t, test.Resource, test.Method)
		t.Log(resp.Body)
		if resp.StatusCode != 200 {
			t.Errorf("error handling request expecting 200, received: %d", resp.StatusCode)
		}

	}
}
