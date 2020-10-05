/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

// Package for handling contact's linked account api requests
package contact

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/wiseco/core-platform/services"

	"github.com/wiseco/core-platform/api"
	"github.com/wiseco/core-platform/test"
)

type testCardRequests struct {
	Method     string
	StatusCode int
	Resource   string
}

var cardTests = []testCardRequests{
	//testRequests{http.MethodGet, http.StatusOK, "/contact"},
	//testCardRequests{http.MethodGet, http.StatusOK, "/contact/3560b088-7be1-4ca3-a514-782a9d982761"},
	//testCardRequests{http.MethodPost, http.StatusOK, "/contact"},
	testCardRequests{http.MethodDelete, http.StatusOK, ""},
}

func testCardPostBody() string {

	testBody := CardPostBody{}

	testBody.CardHolderName = "John Doe"
	testBody.CardNumber = "1234123412341234"
	testBody.CVVCode = "123"
	//testBody.ExpirationDate = time.Now()
	address := services.Address{
		Type:          services.AddressTypeBilling,
		StreetAddress: "255 constitution drive",
		City:          "Menlo Park",
		State:         "CA",
		Country:       "US",
		PostalCode:    "94025",
	}
	testBody.BillingAddress = &address

	b, _ := json.Marshal(testBody)

	return string(b)
}

func runCardJSONRequestTest(t *testing.T, resource, method string) api.APIResponse {

	request := test.TestRequest(method)

	switch method {
	case http.MethodPost:
		request.PathParameters = map[string]string{"businessId": "17689405-8710-43d9-ad24-5ce7a65f188c", "contactId": "3feb9743-81cc-4d1b-b253-477ce1b4c074"}
		request.Body = testCardPostBody()
	case http.MethodGet:
		request.PathParameters = map[string]string{"businessId": "17689405-8710-43d9-ad24-5ce7a65f188c", "contactId": "3feb9743-81cc-4d1b-b253-477ce1b4c074"}
	case http.MethodDelete:
		request.PathParameters = map[string]string{"businessId": "17689405-8710-43d9-ad24-5ce7a65f188c",
			"contactId": "3feb9743-81cc-4d1b-b253-477ce1b4c074", "cardId": "17689405-8710-43d9-ad24-5ce7a65f188c"}
	default:
		break
	}

	request.UserID = "604123ef-9090-4636-bb39-199197533096"

	resp, err := HandleLinkedCardAPIRequests(*request)
	if err != nil {
		t.Errorf("error handling request request for method %s details: %v", method, err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("error handling request for Method: %s expecting 200,received: %d", method, resp.StatusCode)
		return resp
	}

	return resp
}

func TestHandleLinkedCardRequests(t *testing.T) {

	for _, test := range cardTests {
		resp := runCardJSONRequestTest(t, test.Resource, test.Method)
		t.Log(resp.Body)
		if resp.StatusCode != 200 {
			t.Errorf("error handling request expecting 200, received: %d", resp.StatusCode)
		}

	}
}
