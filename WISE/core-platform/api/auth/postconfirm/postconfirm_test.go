package postconfirm

import (
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestHandleCognitoPostConfirmRequestOk(t *testing.T) {
	var event = events.CognitoEventUserPoolsPostConfirmation{
		CognitoEventUserPoolsHeader: events.CognitoEventUserPoolsHeader{
			Version:       "2",
			TriggerSource: "PostConfirmation_ConfirmSignUp",
			Region:        "us-west-2",
			UserPoolID:    "us-west-2_uzJhQpMkO",
		},
		Request: events.CognitoEventUserPoolsPostConfirmationRequest{
			UserAttributes: map[string]string{
				"given_name":   "Joe",
				"middle_name":  "Smithers",
				"family_name":  "Sample",
				"email":        "burim@wise.us",
				"phone_number": "+16503954859",
				"address":      "my address",
			},
		},
		Response: events.CognitoEventUserPoolsPostConfirmationResponse{},
	}

	_, err := HandleCognitoPostConfirmRequest(event)
	if err != nil {
		t.Errorf("TestHandleCognitoPostConfirmRequestOk Failed")
	}

}
