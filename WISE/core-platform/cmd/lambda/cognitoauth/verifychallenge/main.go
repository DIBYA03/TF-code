package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/wiseco/core-platform/api/auth/customevents"
	"github.com/wiseco/core-platform/api/auth/verifychallenge"
)

// HandleVerifyChallengeRequest handle define challenge req
func HandleVerifyChallengeRequest(ctx context.Context, event customevents.CognitoEventUserPoolsVerifyChallenge) (customevents.CognitoEventUserPoolsVerifyChallenge, error) {
	if _, ok := event.Request.UserAttributes["lambda_warmer"]; ok {
		log.Println("lambda warming...")
		return customevents.CognitoEventUserPoolsVerifyChallenge{}, nil
	}

	return verifychallenge.HandleCognitoVerifyChallengeRequest(event)
}

func main() {
	lambda.Start(HandleVerifyChallengeRequest)
}
