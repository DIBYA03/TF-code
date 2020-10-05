/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

// Package for handling contact's linked account api requests
package contact

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/wiseco/core-platform/services/banking"

	"github.com/wiseco/core-platform/api"
	"github.com/wiseco/core-platform/test"
)

type testRequests struct {
	Method     string
	StatusCode int
	Resource   string
}

var tests = []testRequests{
	//testRequests{http.MethodGet, http.StatusOK, "/contact"},
	//testRequests{http.MethodGet, http.StatusOK, "/contact/3560b088-7be1-4ca3-a514-782a9d982761"},
	//testRequests{http.MethodPost, http.StatusOK, "/contact"},
	testRequests{http.MethodDelete, http.StatusOK, ""},
}

func testPostBody() string {

	testBody := AccountPostBody{}

	currency := banking.CurrencyUSD
	permission := "sendOnly"

	testBody.AccountNumber = "123123123"
	testBody.AccountType = "checking"
	testBody.RoutingNumber = "898989"
	testBody.Currency = currency
	testBody.Permission = banking.LinkedAccountPermission(permission)

	b, _ := json.Marshal(testBody)

	return string(b)
}

func runJSONRequestTest(t *testing.T, resource, method string) api.APIResponse {

	request := test.TestRequest(method)

	switch method {
	case http.MethodPost:
		request.PathParameters = map[string]string{"businessId": "17689405-8710-43d9-ad24-5ce7a65f188c", "contactId": "3feb9743-81cc-4d1b-b253-477ce1b4c074"}
		request.Body = testPostBody()
	case http.MethodDelete:
		request.PathParameters = map[string]string{"businessId": "17689405-8710-43d9-ad24-5ce7a65f188c", "contactId": "3feb9743-81cc-4d1b-b253-477ce1b4c074",
			"accountId": "3feb9743-81cc-4d1b-b253-477ce1b4c074"}
	case http.MethodGet:
		request.PathParameters = map[string]string{"businessId": "17689405-8710-43d9-ad24-5ce7a65f188c", "contactId": "3feb9743-81cc-4d1b-b253-477ce1b4c074"}
	default:
		break
	}

	request.UserID = "604123ef-9090-4636-bb39-199197533096"

	resp, err := HandleLinkedAccountAPIRequests(*request)
	if err != nil {
		t.Errorf("error handling request request for method %s details: %v", method, err)
	} else if resp.StatusCode != 200 {
		t.Errorf("error handling request for Method: %s expecting 200,received: %d", method, resp.StatusCode)
		return resp
	}

	return resp
}

func TestHandleLinkedAccountRequests(t *testing.T) {

	for _, test := range tests {
		resp := runJSONRequestTest(t, test.Resource, test.Method)
		t.Log(resp.Body)
		if resp.StatusCode != 200 {
			t.Errorf("error handling request expecting 200, received: %d", resp.StatusCode)
		}

	}
}
