/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

// Package for handling contact's request money api requests
package contact

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/wiseco/core-platform/api"
	"github.com/wiseco/core-platform/services/banking"
	"github.com/wiseco/core-platform/test"
)

type testMoneyRequests struct {
	Method     string
	StatusCode int
	Resource   string
}

var requestTests = []testMoneyRequests{
	//testMoneyRequests{http.MethodGet, http.StatusOK, "/contact"},
	//testMoneyRequests{http.MethodGet, http.StatusOK, "/contact/3560b088-7be1-4ca3-a514-782a9d982761"},
	testMoneyRequests{http.MethodPost, http.StatusOK, "/contact"},
	//testRequests{http.MethodPatch, http.StatusOK, "/contact/d98805a7-a2e3-4499-9b8d-39ad0bc839c9"},
}

func testRequestPostBody() string {
	notes := "<h1>Example</h1>"

	testBody := RequestPostBody{
		Currency: banking.CurrencyUSD,
		Amount:   60,
		Notes:    notes,
	}

	b, _ := json.Marshal(testBody)

	return string(b)
}

func runRequestJSONTest(t *testing.T, resource, method string) api.APIResponse {

	request := test.TestRequest(method)

	switch method {
	case http.MethodPost:
		println("post call here")
		request.PathParameters = map[string]string{"businessId": "7b59ebb0-824c-4019-8d8c-2c1a9ed6192b", "contactId": "9ae5cce8-1587-45a9-99ea-cd2d885ee856"}
		request.Body = testRequestPostBody()
	case http.MethodGet:
		request.PathParameters = map[string]string{"businessId": "7b59ebb0-824c-4019-8d8c-2c1a9ed6192b", "contactId": "9ae5cce8-1587-45a9-99ea-cd2d885ee856"}
	default:
		break
	}

	request.UserID = "74246caf-a317-4f96-aff1-8ad72c16e484"

	resp, err := HandleMoneyRequestAPIRequests(*request)
	if err != nil {
		t.Errorf("error handling request request for method %s details: %v", method, err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("error handling request for Method: %s expecting 200,received: %d", method, resp.StatusCode)
		return resp
	}

	return resp
}

func TestHandleRequestMoneyRequests(t *testing.T) {

	for _, test := range requestTests {
		resp := runRequestJSONTest(t, test.Resource, test.Method)
		t.Log(resp.Body)
		if resp.StatusCode != 200 {
			t.Errorf("error handling request expecting 200, received: %d", resp.StatusCode)
		}

	}
}
