/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

package user

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/wiseco/core-platform/api"
	"github.com/wiseco/core-platform/test"
)

type partnerTestRequests struct {
	Method     string
	StatusCode int
	Resource   string
}

var tests = []partnerTestRequests{
	partnerTestRequests{http.MethodPost, http.StatusOK, ""},
}

func partnerVerificationBody() string {
	testBody := PartnerVerificationBody{}

	testBody.Code = "HOTELKEY"

	b, _ := json.Marshal(testBody)

	return string(b)
}

func runJSONRequestTest(t *testing.T, resource, method string) api.APIResponse {

	request := test.TestRequest(method)

	switch method {
	case http.MethodPost:
		request.Body = partnerVerificationBody()
	default:
		break
	}

	request.UserId = "604123ef-9090-4636-bb39-199197533096"

	resp, err := HandlePartnerCodeVerificationRequest(*request)
	if err != nil {
		t.Errorf("error handling request request for method %s details: %v", method, err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("error handling request for Method: %s expecting 200,received: %d", method, resp.StatusCode)
		return resp
	}

	return resp
}

func TestHandlePartnerVerificationRequests(t *testing.T) {

	for _, test := range tests {
		resp := runJSONRequestTest(t, test.Resource, test.Method)
		t.Log(resp.Body)
		if resp.StatusCode != 200 {
			t.Errorf("error handling request expecting 200, received: %d", resp.StatusCode)
		}

	}
}
