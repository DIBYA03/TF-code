package presignup

import (
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestHandleCognitoPreSignupRequestOk(t *testing.T) {
	var event = events.CognitoEventUserPoolsPreSignup{
		CognitoEventUserPoolsHeader: events.CognitoEventUserPoolsHeader{
			Version:       "2",
			TriggerSource: "PreSignUp_SignUp",
			Region:        "us-west-2",
			UserPoolID:    "us-west-2_uzJhQpMkO",
		},
		Request: events.CognitoEventUserPoolsPreSignupRequest{
			UserAttributes: map[string]string{
				"given_name":   "Joe",
				"middle_name":  "Smithers",
				"family_name":  "Sample",
				"email":        "user@example.com",
				"phone_number": "+14084659283",
				"address":      "255 Constitution Dr, Menlo Park, CA 94025",
			},
			ValidationData: map[string]string{},
		},
		Response: events.CognitoEventUserPoolsPreSignupResponse{},
	}

	_, err := HandleCognitoPreSignupRequest(event)
	if err != nil {
		t.Errorf("TestHandleCognitoPreSignupRequestOk Failed")
	}
}

func TestHandleCognitoPreSignupRequestEmailMissing(t *testing.T) {
	var event = events.CognitoEventUserPoolsPreSignup{
		CognitoEventUserPoolsHeader: events.CognitoEventUserPoolsHeader{
			Version:       "2",
			TriggerSource: "PreSignUp_SignUp",
			Region:        "us-west-2",
			UserPoolID:    "us-west-2_uzJhQpMkO",
		},
		Request: events.CognitoEventUserPoolsPreSignupRequest{
			UserAttributes: map[string]string{
				"given_name":   "Joe",
				"middle_name":  "Smithers",
				"family_name":  "Sample",
				"phone_number": "+14084659283",
				"address":      "255 Constitution Dr, Menlo Park, CA 94025",
			},
			ValidationData: map[string]string{},
		},
		Response: events.CognitoEventUserPoolsPreSignupResponse{},
	}

	_, err := HandleCognitoPreSignupRequest(event)
	if err == nil {
		t.Errorf("TestHandleCognitoPreSignupRequestEmailMissing Failed")
	}
}

func TestHandleCognitoPreSignupRequestEmailInvalid(t *testing.T) {
	var event = events.CognitoEventUserPoolsPreSignup{
		CognitoEventUserPoolsHeader: events.CognitoEventUserPoolsHeader{
			Version:       "2",
			TriggerSource: "PreSignUp_SignUp",
			Region:        "us-west-2",
			UserPoolID:    "us-west-2_uzJhQpMkO",
		},
		Request: events.CognitoEventUserPoolsPreSignupRequest{
			UserAttributes: map[string]string{
				"given_name":   "Joe",
				"middle_name":  "Smithers",
				"family_name":  "Sample",
				"email":        "@example.com",
				"phone_number": "+14084659283",
				"address":      "255 Constitution Dr, Menlo Park, CA 94025",
			},
			ValidationData: map[string]string{},
		},
		Response: events.CognitoEventUserPoolsPreSignupResponse{},
	}

	_, err := HandleCognitoPreSignupRequest(event)
	if err == nil {
		t.Errorf("TestHandleCognitoPreSignupRequestEmailInvalid Failed")
	}
}
