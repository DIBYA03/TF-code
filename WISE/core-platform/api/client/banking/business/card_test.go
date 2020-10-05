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

type testCardRequests struct {
	Method     string
	StatusCode int
	Resource   string
}

var cardTests = []testCardRequests{
	testCardRequests{http.MethodGet, http.StatusOK, ""},
	//testCardRequests{http.MethodPost, http.StatusOK, ""},
}

func testPostBody() string {
	testBody := CardPostBody{}
	testBody.CardType = "debit"

	b, _ := json.Marshal(testBody)

	return string(b)
}

func runCardJSONRequestTest(t *testing.T, resource, method string) api.APIResponse {

	request := test.TestRequest(method)

	request.PathParameters = map[string]string{"businessId": "17689405-8710-43d9-ad24-5ce7a65f188c", "accountId": "1bbaf747-c330-45a2-9e21-1fecef8d9205"}
	request.UserID = "604123ef-9090-4636-bb39-199197533096"

	switch method {
	case http.MethodPost:
		request.Body = testPostBody()
	default:
		break
	}

	resp, err := HandleBusinessBankCardAPIRequests(*request)
	if err != nil {
		t.Errorf("error handling request request for method %s details: %v", method, err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("error handling request for Method: %s expecting 200,received: %d", method, resp.StatusCode)
		return resp
	}

	return resp
}

func TestHandleCardRequests(t *testing.T) {

	for _, test := range cardTests {
		resp := runCardJSONRequestTest(t, test.Resource, test.Method)
		t.Log(resp.Body)
		if resp.StatusCode != 200 {
			t.Errorf("error handling request expecting 200, received: %d", resp.StatusCode)
		}

	}
}
