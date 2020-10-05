/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

// Package for handling business api requests
package document

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"testing"

	"github.com/wiseco/core-platform/api"
)

type documentTestRequests struct {
	filename   string
	statusCode int
}

var documentTests = []documentTestRequests{
	documentTestRequests{"testdata/document_get_request_200.json", 200},
	documentTestRequests{"testdata/document_post_request_200.json", 200},
	documentTestRequests{"testdata/document_patch_request_200.json", 200},
}

func runDocumentJSONRequestTest(t *testing.T, jsonFile string) (api.APIResponse, error) {
	requestJSON, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		t.Errorf("could not open test file. details: %v", err)
	}

	var request api.APIRequest
	if err := json.Unmarshal(requestJSON, &request); err != nil {
		t.Errorf("could not unmarshal request. details: %v", err)
	}

	resp, err := HandleDocumentAPIRequests(request)

	if resp.StatusCode != 200 {
		t.Errorf("error handling request file: %s. details %v", jsonFile, err)
	}

	return resp, err
}

func TestHandleDocumentAPIRequests(t *testing.T) {

	for _, test := range documentTests {
		resp, _ := runDocumentJSONRequestTest(t, test.filename)
		log.Println(resp.Body)
		if resp.StatusCode != 200 {
			t.Errorf("error handling request file: %s. expecting 200, received: %d", "testdata/business/document/document_get_request_200.json", resp.StatusCode)
		}

	}
}
