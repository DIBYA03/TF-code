package definechallenge

import (
	"fmt"

	"github.com/wiseco/core-platform/api/auth/customevents"
)

// HandleCognitoDefineChallengeRequest is used on custom messages
func HandleCognitoDefineChallengeRequest(event customevents.CognitoEventUserPoolsDefineChallenge) (customevents.CognitoEventUserPoolsDefineChallenge, error) {
	//event.Response.challengeName = "CUSTOM_CHALLENGE", "PASSWORD_VERIFIER", "SMS_MFA", "DEVICE_SRP_AUTH", "DEVICE_PASSWORD_VERIFIER", or "ADMIN_NO_SRP_AUTH"
	event.Response.IssueTokens = true // Issue access tokens for now, on future this can be conditional
	fmt.Println(event)
	return event, nil
}
