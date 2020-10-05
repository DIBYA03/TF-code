/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

package transaction

import (
	"net/http"
	"testing"

	"github.com/wiseco/core-platform/api"
	"github.com/wiseco/core-platform/test"
)

type testTransactionRequests struct {
	Method     string
	StatusCode int
	Resource   string
}

var transactionTests = []testTransactionRequests{
	testTransactionRequests{http.MethodGet, http.StatusOK, ""},
}

func runTransactionJSONRequestTest(t *testing.T, resource, method string) api.APIResponse {

	request := test.TestRequest(method)

	request.PathParameters = map[string]string{"businessId": "69c9fb64-a084-4621-91ab-e84edbe31292"}
	request.UserID = "604123ef-9090-4636-bb39-199197533096"
	request.QueryStringParameters = map[string]string{"limit": "20", "offset": "0",
		"startDate": "2019-07-01T13:30:00Z", "endDate": "2019-07-03T16:08:00Z", "type": "creditPosted", "maxAmount": "1", "text": "SADA"}

	switch method {
	case http.MethodPost:
		request.Body = testPostBody()
	default:
		break
	}

	println("running transactions..")

	resp, err := HandleTransactionRequest(*request)
	if err != nil {
		t.Errorf("error handling request request for method %s details: %v", method, err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("error handling request for Method: %s expecting 200,received: %d", method, resp.StatusCode)
		return resp
	}

	return resp
}

func runTransactionExportJSONRequestTest(t *testing.T, resource, method string) api.APIResponse {

	request := test.TestRequest(method)

	request.PathParameters = map[string]string{"businessId": "69c9fb64-a084-4621-91ab-e84edbe31292"}
	request.UserID = "604123ef-9090-4636-bb39-199197533096"
	request.QueryStringParameters = map[string]string{
		"startDate": "2019-07-01T13:30:00Z", "endDate": "2019-07-03T16:30:00Z"}

	resp, err := HandleTransactionExportRequest(*request)
	if err != nil {
		t.Errorf("error handling request request for method %s details: %v", method, err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("error handling request for Method: %s expecting 200,received: %d", method, resp.StatusCode)
		return resp
	}

	return resp
}

func TestHandleTransactionRequests(t *testing.T) {

	for _, test := range transactionTests {
		resp := runTransactionExportJSONRequestTest(t, test.Resource, test.Method)
		t.Log(resp.Body)
		if resp.StatusCode != 200 {
			t.Errorf("error handling request expecting 200, received: %d", resp.StatusCode)
		}

	}
}
