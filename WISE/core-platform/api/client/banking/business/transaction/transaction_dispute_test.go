/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

package transaction

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/wiseco/core-platform/api"
	"github.com/wiseco/core-platform/test"
)

type testTransactionDisputeRequests struct {
	Method     string
	StatusCode int
	Resource   string
}

var transactionDisputeTests = []testTransactionDisputeRequests{
	testTransactionDisputeRequests{http.MethodPost, http.StatusOK, ""},
}

func testDisputePostBody() string {

	summary := "My card was incorrectly charged"
	receiptID := "ab6806ab-385c-4411-8972-1463f76b25b1"

	testBody := DisputePostBody{}

	testBody.Category = "stillBeingCharged"
	testBody.Summary = &summary
	testBody.ReceiptID = &receiptID

	b, _ := json.Marshal(testBody)

	return string(b)
}

func runTransactionDisputeJSONRequestTest(t *testing.T, resource, method string) api.APIResponse {

	request := test.TestRequest(method)

	request.PathParameters = map[string]string{"businessId": "69c9fb64-a084-4621-91ab-e84edbe31292", "transactionId": "ab6806ab-385c-4411-8972-1463f76b25b1"}
	request.UserID = "08124497-9261-4ca0-a470-b570a7badd58"

	switch method {
	case http.MethodPost:
		request.Body = testDisputePostBody()
	default:
		break
	}

	resp, err := HandleTransactionDisputeRequest(*request)
	if err != nil {
		t.Errorf("error handling request request for method %s details: %v", method, err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("error handling request for Method: %s expecting 200,received: %d", method, resp.StatusCode)
		return resp
	}

	return resp
}

func TestHandleTransactionDisputeRequests(t *testing.T) {

	for _, test := range transactionDisputeTests {
		resp := runTransactionDisputeJSONRequestTest(t, test.Resource, test.Method)
		t.Log(resp.Body)
		if resp.StatusCode != 200 {
			t.Errorf("error handling request expecting 200, received: %d", resp.StatusCode)
		}

	}
}
