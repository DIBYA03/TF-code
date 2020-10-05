package verifychallenge

import (
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/wiseco/core-platform/api/auth/customevents"
)

func TestHandleCognitoCustomMessageRequestOk(t *testing.T) {
	var event = customevents.CognitoEventUserPoolsVerifyChallenge{
		CognitoEventUserPoolsHeader: events.CognitoEventUserPoolsHeader{
			Version:       "2",
			TriggerSource: "CustomMessage_SignUp",
			Region:        "us-west-2",
			UserPoolID:    "us-west-2_uzJhQpMkO",
			UserName:      "8bbaabe8-09d0-49c0-b4ee-b0dc8dd62eff",
		},
		Request: customevents.CognitoEventUserPoolsVerifyChallengeRequest{
			UserAttributes: map[string]string{
				"username": "burim@wise.us",
				"password": "12345678",
			},
			PrivateChallengeParameters: map[string]string{
				"answer": "5",
			},
			ChallengeAnswer: map[string]string{
				"answer": "5",
			},
		},
		Response: customevents.CognitoEventUserPoolsVerifyChallengeResponse{
			AnswerCorrect: true,
		},
	}

	_, err := HandleCognitoVerifyChallengeRequest(event)
	if err != nil {
		t.Errorf("TestHandleCognitoCustomMessageRequestOk Failed")
	}

}
