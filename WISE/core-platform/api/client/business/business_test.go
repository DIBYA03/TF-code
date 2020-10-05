/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

// Package for handling business api requests
package business

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/wiseco/core-platform/api"
	"github.com/wiseco/core-platform/services"
	"github.com/wiseco/core-platform/test"
)

type testRequests struct {
	Method     string
	StatusCode int
	Resource   string
}

var tests = []testRequests{
	//testRequests{http.MethodGet, http.StatusOK, "/business"},
	//testRequests{http.MethodGet, http.StatusOK, "/business/3560b088-7be1-4ca3-a514-782a9d982761"},
	testRequests{http.MethodPost, http.StatusOK, "/business"},
	testRequests{http.MethodPatch, http.StatusOK, "/business/fe197a11-6757-4f62-a191-77a4abfbd610"},
}

func testPostBody() string {
	testBody := BusinessPostBody{
		EmployerNumber: "654321",
		//LegalName:      "testing body",
		//Purpose:        "payroll",
	}

	et := "LLC"
	dba := services.StringArray{}
	dba = append(dba, "cafe")
	dba = append(dba, "other cafe")
	address := services.Address{
		StreetAddress: "2175 cooley",
		State:         "CA",
		City:          "Palo Alto",
	}

	industry := "Cafe"
	testBody.EntityType = &et
	testBody.IndustryType = &industry
	testBody.LegalAddress = &address
	testBody.DBA = dba
	b, _ := json.Marshal(testBody)

	return string(b)
}

func testPatchBody() string {
	testBody := BusinessPostBody{
		EmployerNumber: "654321",
		//LegalName:      "testing body",
		//Purpose:        "payroll",
	}

	et := "LLC"
	dba := services.StringArray{}
	dba = append(dba, "cafe")
	dba = append(dba, "other cafe")

	address := services.Address{
		StreetAddress: "2175 cooley",
		State:         "CA",
		City:          "Palo Alto",
	}

	industry := "Cafe"
	testBody.EntityType = &et
	testBody.IndustryType = &industry
	testBody.LegalAddress = &address
	testBody.DBA = dba
	b, _ := json.Marshal(testBody)

	return string(b)
}

func runJSONRequestTest(t *testing.T, resource, method string) api.APIResponse {

	request := test.TestRequest(method)
	switch method {
	case http.MethodPost:
		request.Body = testPostBody()
	case http.MethodPatch:
		request.PathParameters = map[string]string{"businessId": "fe197a11-6757-4f62-a191-77a4abfbd610"}
		request.Body = testPatchBody()
	default:
		break
	}

	request.UserID = "aa41fcc7-acbb-480d-b285-dded8bcdf645"

	resp, err := HandleBusinessAPIRequests(*request)
	if err != nil {
		t.Errorf("error handling request request for method %s details: %v", method, err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("error handling request for Method: %s expecting 200,received: %d", method, resp.StatusCode)
		return resp
	}

	return resp
}

func TestHandleBusinessAPIRequests(t *testing.T) {

	for _, test := range tests {
		resp := runJSONRequestTest(t, test.Resource, test.Method)
		t.Log(resp.Body)
		if resp.StatusCode != 200 {
			t.Errorf("error handling request expecting 200, received: %d", resp.StatusCode)
		}

	}
}
