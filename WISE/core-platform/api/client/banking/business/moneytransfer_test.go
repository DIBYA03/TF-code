/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

// Package for handling linked account and its related api requests
package business

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
	testTransferRequests{http.MethodGet, http.StatusOK, "/contact"},
	//testTransferRequests{http.MethodPost, http.StatusOK, "/contact/3560b088-7be1-4ca3-a514-782a9d982761"},
}

func testTransferPostBody() string {

	testBody := TransferPostBody{}

	testBody.Amount = 1.00
	testBody.SourceAccountId = "6291096d-27ff-46e7-83d5-fdbe5fd05283"
	testBody.SourceType = "account"
	testBody.DestAccountId = "f4e46263-d27d-4754-ad0a-5c2343ecb47c"
	testBody.DestType = "account"
	testBody.Currency = "usd"

	b, _ := json.Marshal(testBody)

	return string(b)
}

func runTransferJSONRequestTest(t *testing.T, resource, method string) api.APIResponse {

	request := test.TestRequest(method)

	switch method {
	case http.MethodPost:
		request.PathParameters = map[string]string{"businessId": "e21cdcb3-895b-433a-8bc0-e41e624e110e"}
		request.Body = testTransferPostBody()
	case http.MethodGet:
		request.PathParameters = map[string]string{"businessId": "e21cdcb3-895b-433a-8bc0-e41e624e110e"}
	default:
		break
	}

	request.UserID = "fc66f20b-4c7f-4a7e-b1a2-e57e88efca5a"

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
