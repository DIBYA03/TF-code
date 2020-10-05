/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

// Package for handling contact's linked account api requests
package contact

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/wiseco/core-platform/api"
	"github.com/wiseco/core-platform/test"
)

type testTransferRequests struct {
	Method     string
	StatusCode int
	Resource   string
}

var transferTests = []testTransferRequests{
	//testRequests{http.MethodGet, http.StatusOK, "/contact"},
	testTransferRequests{http.MethodGet, http.StatusOK, "/contact/3560b088-7be1-4ca3-a514-782a9d982761"},
	//testCardRequests{http.MethodPost, http.StatusOK, "/contact"},
	//testRequests{http.MethodPatch, http.StatusOK, "/contact/d98805a7-a2e3-4499-9b8d-39ad0bc839c9"},
}

func testTransferPostBody() string {

	testBody := TransferPostBody{}

	testBody.Amount = 100.00

	b, _ := json.Marshal(testBody)

	return string(b)
}

func runTransferJSONRequestTest(t *testing.T, resource, method string) api.APIResponse {

	request := test.TestRequest(method)

	switch method {
	case http.MethodPost:
		request.PathParameters = map[string]string{"businessId": "17689405-8710-43d9-ad24-5ce7a65f188c", "contactId": "dc5467ed-e536-4798-8405-580fd5e81b4a"}
		request.Body = testTransferPostBody()
	case http.MethodGet:
		request.PathParameters = map[string]string{"businessId": "17689405-8710-43d9-ad24-5ce7a65f188c", "contactId": "dc5467ed-e536-4798-8405-580fd5e81b4a"}
	default:
		break
	}

	request.UserID = "604123ef-9090-4636-bb39-199197533096"

	resp, err := HandleMoneyTransferAPIRequests(*request)
	if err != nil {
		t.Errorf("error handling request request for method %s details: %v", method, err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("error handling request for Method: %s expecting 200,received: %d", method, resp.StatusCode)
		return resp
	}

	return resp
}

func TestHandleMoneyTransferRequests(t *testing.T) {

	for _, test := range transferTests {
		resp := runTransferJSONRequestTest(t, test.Resource, test.Method)
		t.Log(resp.Body)
		if resp.StatusCode != 200 {
			t.Errorf("error handling request expecting 200, received: %d", resp.StatusCode)
		}

	}
}
