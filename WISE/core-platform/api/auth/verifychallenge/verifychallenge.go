package verifychallenge

import (
	"fmt"

	"github.com/wiseco/core-platform/api/auth/customevents"
)

// HandleCognitoVerifyChallengeRequest is used on custom messages
func HandleCognitoVerifyChallengeRequest(event customevents.CognitoEventUserPoolsVerifyChallenge) (customevents.CognitoEventUserPoolsVerifyChallenge, error) {

	// if event.Request.PrivateChallengeParameters["answer"] == event.Request.ChallengeAnswer["answer"] {
	// 	event.Response.AnswerCorrect = true
	// } else {
	// 	event.Response.AnswerCorrect = false
	// }

	fmt.Println(event)
	return event, nil
}
