package createchallenge

import (
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/wiseco/core-platform/api/auth/customevents"
)

func TestHandleCognitoCreateChallengeRequestOk(t *testing.T) {
	var event = customevents.CognitoEventUserPoolsCreateChallenge{
		CognitoEventUserPoolsHeader: events.CognitoEventUserPoolsHeader{
			Version:       "2",
			TriggerSource: "CustomMessage_SignUp",
			Region:        "us-west-2",
			UserPoolID:    "us-west-2_uzJhQpMkO",
			UserName:      "8bbaabe8-09d0-49c0-b4ee-b0dc8dd62eff",
		},
		Request: customevents.CognitoEventUserPoolsCreateChallengeRequest{
			UserAttributes: map[string]string{
				"username": "burim@wise.us",
				"password": "12345678",
			},
			ChallengeName: "CUSTOM_CHALLENGE",

			Session: map[string]string{
				"ChallengeResult": "false",
			},
		},
		Response: customevents.CognitoEventUserPoolsCreateChallengeResponse{
			PublicChallengeParameters: map[string]string{
				"captchaUrl": "chaptcha/Image",
			},
			PrivateChallengeParameters: map[string]string{
				"answer": "5",
			},
		},
	}

	_, err := HandleCognitoCreateChallengeRequest(event)
	if err != nil {
		t.Errorf("TestHandleCognitoCreateChallengeRequestOk Failed")
	}

}
