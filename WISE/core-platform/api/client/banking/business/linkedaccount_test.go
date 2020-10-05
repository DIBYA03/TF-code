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

type registerAccountTestRequests struct {
	Method     string
	StatusCode int
	Resource   string
}

var registerAccountTests = []registerAccountTestRequests{
	//registerAccountTestRequests{http.MethodPost, http.StatusOK, ""},
	//registerAccountTestRequests{http.MethodDelete, http.StatusOK, ""},
	registerAccountTestRequests{http.MethodGet, http.StatusOK, ""},
}

func testRegisterBankAccount() string {
	testBody := LinkAccountCreateBody{
		PublicToken:     "public-sandbox-6cb08832-d6dc-463c-b81a-e3989c941120",
		SourceAccountId: "Wmw1D8RwoAtQQgy4GvvgtMQzg8x6XdFlRJarr",
	}

	b, _ := json.Marshal(testBody)

	return string(b)
}

func runRegisterJSONRequestTest(t *testing.T, resource, method string) api.APIResponse {

	request := test.TestRequest(method)

	switch method {
	case http.MethodPost:
		request.PathParameters = map[string]string{"businessId": "17689405-8710-43d9-ad24-5ce7a65f188c"}
		request.Body = testRegisterBankAccount()
	case http.MethodDelete:
		request.PathParameters = map[string]string{"businessId": "e21cdcb3-895b-433a-8bc0-e41e624e110e",
			"accountId": "6291096d-27ff-46e7-83d5-fdbe5fd05283"}
	case http.MethodGet:
		request.PathParameters = map[string]string{"businessId": "e21cdcb3-895b-433a-8bc0-e41e624e110e"}
	default:
		break
	}

	request.UserID = "fc66f20b-4c7f-4a7e-b1a2-e57e88efca5a"

	resp, err := HandleLinkAccountRequest(*request)
	if err != nil {
		t.Errorf("error handling request request for method %s details: %v", method, err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("error handling request for Method: %s expecting 200,received: %d", method, resp.StatusCode)
		return resp
	}

	return resp
}

func TestHandleLinkAccountRequests(t *testing.T) {

	for _, test := range registerAccountTests {
		resp := runRegisterJSONRequestTest(t, test.Resource, test.Method)
		t.Log(resp.Body)
		if resp.StatusCode != 200 {
			t.Errorf("error handling request expecting 200, received: %d", resp.StatusCode)
		}

	}
}
