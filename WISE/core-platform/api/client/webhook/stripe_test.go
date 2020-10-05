/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

// Package for handling contact api requests
package webhook

import (
	"encoding/json"
	"net/http"
	"testing"

	b "github.com/wiseco/core-platform/services/banking/business"
	"github.com/wiseco/core-platform/services/banking/business/contact"
	"github.com/wiseco/core-platform/test"
)

type testStripeRequests struct {
	Method     string
	StatusCode int
	Resource   string
}

var tests = []testStripeRequests{
	testStripeRequests{http.MethodPost, http.StatusOK, "/contact"},
}

func testPostBody() string {

	testBody := b.Payment{
		Id:     "pi_1Eae2mDu1MErS7u8h95WPcz4",
		Status: "succeeded",
	}

	b, _ := json.Marshal(testBody)

	return string(b)
}

func runJSONStripeTest(t *testing.T, resource, method string) {

	request := test.TestRequest(method)

	request.PathParameters = map[string]string{"businessId": "ef1249ae-6b60-4802-a2ae-8213b5c58bc4"}

	switch method {
	case http.MethodPost:
		request.Body = testPostBody()
	default:
		break
	}

	request.UserID = "604123ef-9090-4636-bb39-199197533096"

	testBody := b.Payment{
		Id:     "pi_1Ec4pLDu1MErS7u8S7Dboq3n",
		Status: "succeeded",
	}

	println("before calling handlewebhook")

	go contact.NewMoneyRequestService(request.SourceRequest()).HandleWebhook(&testBody)

	println("after calling handlewebhook")

}

func TestHandleStripeRequests(t *testing.T) {

	for _, test := range tests {
		runJSONStripeTest(t, test.Resource, test.Method)

	}
}
