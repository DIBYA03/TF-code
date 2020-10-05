package custommessage

import (
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/wiseco/core-platform/api/auth/customevents"
)

func TestHandleCognitoCustomMessageRequestOk(t *testing.T) {
	var event = customevents.CognitoEventUserPoolsCustomMessage{
		CognitoEventUserPoolsHeader: events.CognitoEventUserPoolsHeader{
			Version:       "2",
			TriggerSource: "CustomMessage_SignUp",
			Region:        "us-west-2",
			UserPoolID:    "us-west-2_uzJhQpMkO",
			UserName:      "8bbaabe8-09d0-49c0-b4ee-b0dc8dd62eff",
		},
		Request: customevents.CognitoEventUserPoolsCustomMessageRequest{
			UserAttributes: map[string]string{
				"given_name":            "Joe",
				"middle_name":           "Smithers",
				"family_name":           "Sample",
				"email":                 "burim@wise.us",
				"phone_number":          "+16503954859",
				"sub":                   "85539d6d-4aca-4cc6-8d79-a4c0e954313d",
				"email_verified":        "false",
				"phone_number_verified": "false",
				"cognito:user_status":   "UNCONFIRMED",
			},
			CodeParameter:     "345454",
			UsernameParameter: "8bbaabe8-09d0-49c0-b4ee-b0dc8dd62eff",
		},
		Response: customevents.CognitoEventUserPoolsCustomMessageResponse{},
	}

	_, err := HandleCognitoCustomMessageRequest(event)
	if err != nil {
		t.Errorf("TestHandleCognitoCustomMessageRequestOk Failed")
	}

}
