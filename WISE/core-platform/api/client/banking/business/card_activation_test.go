/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

// Package for handling contact api requests
package business

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/wiseco/core-platform/api"
	"github.com/wiseco/core-platform/test"
)

type testCardActivateRequests struct {
	Method     string
	StatusCode int
	Resource   string
}

var cardActivateTests = []testCardActivateRequests{
	testCardActivateRequests{http.MethodPost, http.StatusOK, ""},
}

func testActivatePostBody() string {
	testBody := ActivatePostBody{}
	testBody.PANLast6 = "06 4687"

	b, _ := json.Marshal(testBody)

	return string(b)
}

func runCardActivateJSONRequestTest(t *testing.T, resource, method string) api.APIResponse {

	request := test.TestRequest(method)

	request.PathParameters = map[string]string{"businessId": "17689405-8710-43d9-ad24-5ce7a65f188c", "cardId": "a81d3d09-535d-4e0c-95ca-133e5bedd3c5"}
	request.UserID = "604123ef-9090-4636-bb39-199197533096"

	switch method {
	case http.MethodPost:
		request.Body = testActivatePostBody()
	default:
		break
	}

	resp, err := HandleBankCardActivationAPIRequests(*request)
	if err != nil {
		t.Errorf("error handling request request for method %s details: %v", method, err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("error handling request for Method: %s expecting 200,received: %d", method, resp.StatusCode)
		return resp
	}

	return resp
}

func TestHandleCardActivateRequests(t *testing.T) {

	for _, test := range cardActivateTests {
		resp := runCardActivateJSONRequestTest(t, test.Resource, test.Method)
		t.Log(resp.Body)
		if resp.StatusCode != 200 {
			t.Errorf("error handling request expecting 200, received: %d", resp.StatusCode)
		}

	}
}
