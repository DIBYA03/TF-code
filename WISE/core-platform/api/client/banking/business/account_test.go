package business

import (
	"encoding/json"
	"log"
	"net/http"
	"testing"

	"github.com/wiseco/core-platform/api"
	bankaccountsrv "github.com/wiseco/core-platform/services/banking/business"
	"github.com/wiseco/core-platform/test"
)

type testRequests struct {
	method     string
	statusCode int
	resource   string
}

var tests = []testRequests{
	testRequests{http.MethodGet, http.StatusOK, ""},
	//testRequests{http.MethodPost, http.StatusOK, ""},
	//testRequests{http.MethodPatch, http.StatusOK, ""},
}

func postRequestBody() string {
	testBody := bankaccountsrv.BankAccountCreate{}

	b, _ := json.Marshal(testBody)

	return string(b)
}

func patchRequestBody() string {
	testBody := bankaccountsrv.BankAccountCreate{}

	b, _ := json.Marshal(testBody)

	return string(b)
}

func runJSONRequestTest(t *testing.T, method string) api.APIResponse {
	request := test.TestRequest(method)
	switch method {
	case http.MethodGet, http.MethodPost:
		request.PathParameters = map[string]string{"businessId": "fe197a11-6757-4f62-a191-77a4abfbd610"}
		request.Body = postRequestBody()
	case http.MethodPatch:
		request.Body = patchRequestBody()
		request.PathParameters = map[string]string{"accountId": "6d0c19a6-0074-4ffc-aa51-9122f0457ac0",
			"businessId": "fe197a11-6757-4f62-a191-77a4abfbd610",
		}
	}

	request.UserID = "1bf0faf1-cb57-4053-a2ae-e693ec99d7a3"
	resp, err := HandleAccountAPIRequest(*request)
	if err != nil {
		t.Errorf("error handling request for method: %s. details %v", method, err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("error handling request for Method: %s expecting 200,received: %d", method, resp.StatusCode)
		return resp
	}

	return resp
}

func TestHandleAccountAPIRequests(t *testing.T) {

	for _, ts := range tests {
		resp := runJSONRequestTest(t, ts.method)
		log.Println(resp.Body)
		if resp.StatusCode != 200 {
			t.Errorf("error handling request: %s. expecting 200, received: %d", ts.resource, resp.StatusCode)
		}

	}
}
