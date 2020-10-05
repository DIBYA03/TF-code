package createchallenge

import (
	"fmt"

	"github.com/wiseco/core-platform/api/auth/customevents"
)

// HandleCognitoCreateChallengeRequest is used on custom messages
func HandleCognitoCreateChallengeRequest(event customevents.CognitoEventUserPoolsCreateChallenge) (customevents.CognitoEventUserPoolsCreateChallenge, error) {
	//event.Response.challengeName = "CUSTOM_CHALLENGE", "PASSWORD_VERIFIER", "SMS_MFA", "DEVICE_SRP_AUTH", "DEVICE_PASSWORD_VERIFIER", or "ADMIN_NO_SRP_AUTH"

	if event.Request.ChallengeName == "CUSTOM_CHALLENGE" { // example but commented out

		//event.Response.PublicChallengeParameters = make(map[string]string)
		//event.Response.PublicChallengeParameters.captchaUrl = "'url/123.jpg"
		//event.Response.PrivateChallengeParameters = map[string]string
		//event.Response.PrivateChallengeParameters["answer"] = "5"
		//event.Response.ChallengeMetadata = "CAPTCHA_CHALLENGE"

	}
	fmt.Println(event)
	return event, nil
}
