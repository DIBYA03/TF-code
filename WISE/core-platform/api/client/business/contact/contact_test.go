/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

// Package for handling contact api requests
package contact

import (
	"encoding/json"
	"net/http"
	"testing"

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
	testRequests{http.MethodPost, http.StatusOK, "/contact"},
	//testRequests{http.MethodPatch, http.StatusOK, "/contact/d98805a7-a2e3-4499-9b8d-39ad0bc839c9"},
}

func testPostBody() string {
	var engagement = "part-time"
	var title = "CEO"

	testBody := ContactPostBody{}

	testBody.Type = "business"
	testBody.Email = "john.doe@gmail.com"
	testBody.PhoneNumber = "+11231234567"
	testBody.Engagement = &engagement
	testBody.JobTitle = &title

	busName := "Joe Inc"
	testBody.BusinessName = &busName

	b, _ := json.Marshal(testBody)

	return string(b)
}

func testPatchBody() string {
	var engagement = "full-time"
	var title = "CFO"
	var email = "contractor"
	var phone = "+11231231231"
	var firstName = "Johniee"
	var lastName = "Walkers"

	testBody := ContactPatchBody{}

	testBody.Email = &email
	testBody.PhoneNumber = &phone
	testBody.FirstName = &firstName
	testBody.LastName = &lastName
	testBody.Engagement = &engagement
	testBody.JobTitle = &title

	b, _ := json.Marshal(testBody)

	return string(b)
}

func runJSONRequestTest(t *testing.T, resource, method string) api.APIResponse {

	println("running json request")

	request := test.TestRequest(method)

	request.PathParameters = map[string]string{"businessId": "7b59ebb0-824c-4019-8d8c-2c1a9ed6192b"}

	switch method {
	case http.MethodPost:
		request.Body = testPostBody()
	case http.MethodPatch:
		request.PathParameters = map[string]string{"businessId": "7b59ebb0-824c-4019-8d8c-2c1a9ed6192b", "contactId": "6618f8c6-4f13-43de-adb5-6584af8aee0e"}
		request.Body = testPatchBody()
	case http.MethodDelete:
		request.PathParameters = map[string]string{"businessId": "7b59ebb0-824c-4019-8d8c-2c1a9ed6192b", "contactId": "6618f8c6-4f13-43de-adb5-6584af8aee0e"}
	default:
		break
	}

	request.UserID = "74246caf-a317-4f96-aff1-8ad72c16e484"

	resp, err := HandleContactAPIRequests(*request)
	if err != nil {
		t.Errorf("error handling request request for method %s details: %v", method, err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("error handling request for Method: %s expecting 200,received: %d", method, resp.StatusCode)
		return resp
	}

	return resp
}

func TestHandleContactRequests(t *testing.T) {

	for _, test := range tests {
		resp := runJSONRequestTest(t, test.Resource, test.Method)
		t.Log(resp.Body)
		if resp.StatusCode != 200 {
			t.Errorf("error handling request expecting 200, received: %d", resp.StatusCode)
		}

	}
}
