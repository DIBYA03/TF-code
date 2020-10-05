/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

// Package for handling business member api requests
package business

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/wiseco/core-platform/api"
	"github.com/wiseco/core-platform/services"
	b "github.com/wiseco/core-platform/services/business"
	"github.com/wiseco/core-platform/test"
)

type memberTestRequests struct {
	Method     string
	StatusCode int
	Resource   string
}

var memberTests = []memberTestRequests{
	//memberTestRequests{http.MethodGet, http.StatusOK, "/business/fe197a11-6757-4f62-a191-77a4abfbd610/member/a4e4ebfc-caba-4e30-bde0-4d258f137f6a"},
	//memberTestRequests{http.MethodPost, http.StatusOK, "/business/17689405-8710-43d9-ad24-5ce7a65f188c/member"},
	memberTestRequests{http.MethodPatch, http.StatusOK, "/business/fe197a11-6757-4f62-a191-77a4abfbd610/member/a4e4ebfc-caba-4e30-bde0-4d258f137f6a"},
}

func testMemberPostBody() string {
	email := "henry.wilson@RILB.com"
	dateOfBirth := services.Date(time.Now())
	taxId := services.TaxID("123123123")
	taxIdType := "ssn"

	incomeType := services.StringArray{"Salary"}
	occupation := "Management"

	residencyStatus := services.Residency{
		Country: "US",
		Status:  "citizen",
	}

	testBody := businessMemberPostBody{
		//UserId:          "604123ef-9090-4636-bb39-199197533096",
		BusinessId:           "17689405-8710-43d9-ad24-5ce7a65f188c",
		FirstName:            "Henry",
		LastName:             "Wilson",
		Email:                email,
		Phone:                "+16501234567",
		TitleType:            "Manager",
		DateOfBirth:          &dateOfBirth,
		Ownership:            100,
		IsControllingManager: true,
		TaxId:                taxId,
		TaxIdType:            taxIdType,
		MemberResidency:      &residencyStatus,
		Occupation:           occupation,
	}

	testBody.IncomeType = incomeType

	CitizenshipCountries := services.StringArray{}
	CitizenshipCountries = append(CitizenshipCountries, "US")
	testBody.CitizenshipCountries = CitizenshipCountries

	address := services.Address{
		StreetAddress: "2175 cooley",
		State:         "CA",
		City:          "Palo Alto",
		PostalCode:    "12345",
	}
	testBody.LegalAddress = &address

	b, _ := json.Marshal(testBody)

	return string(b)
}

func testMemberPatchBody() string {
	o := 26
	titleType := b.TitleTypePresident

	testBody := businessMemberPatchBody{
		Ownership: &o,
		TitleType: &titleType,
	}

	b, _ := json.Marshal(testBody)

	return string(b)
}

func runMemberJSONRequestTest(t *testing.T, resource, method string) api.APIResponse {

	request := test.TestRequest(method)
	switch method {
	case http.MethodPost:
		request.PathParameters = map[string]string{"businessId": "17689405-8710-43d9-ad24-5ce7a65f188c"}
		request.Body = testMemberPostBody()
	case http.MethodPatch:
		request.PathParameters = map[string]string{"businessId": "17689405-8710-43d9-ad24-5ce7a65f188c",
			"memberId": "d012db7c-df0b-49c2-a44e-0d3e1be43b8a"}
		request.Body = testMemberPatchBody()
	case http.MethodGet:
		request.PathParameters = map[string]string{"businessId": "17689405-8710-43d9-ad24-5ce7a65f188c"}
	default:
		break
	}

	request.UserId = "604123ef-9090-4636-bb39-199197533096"

	resp, err := HandleMemberAPIRequests(*request)
	if err != nil {
		t.Errorf("error handling request request for method %s details: %v", method, err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("error handling request for Method: %s expecting 200,received: %d", method, resp.StatusCode)
		return resp
	}

	return resp
}

func TestHandleBusinessMemberAPIRequests(t *testing.T) {

	for _, test := range memberTests {
		resp := runMemberJSONRequestTest(t, test.Resource, test.Method)
		t.Log(resp.Body)
		if resp.StatusCode != 200 {
			t.Errorf("error handling request expecting 200, received: %d", resp.StatusCode)
		}

	}
}
