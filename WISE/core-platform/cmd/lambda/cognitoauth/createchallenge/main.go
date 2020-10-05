package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/wiseco/core-platform/api/auth/createchallenge"
	"github.com/wiseco/core-platform/api/auth/customevents"
)

// HandleCreateChallengeRequest handle define challenge req
func HandleCreateChallengeRequest(ctx context.Context, event customevents.CognitoEventUserPoolsCreateChallenge) (customevents.CognitoEventUserPoolsCreateChallenge, error) {
	if _, ok := event.Request.UserAttributes["lambda_warmer"]; ok {
		log.Println("lambda warming...")
		return customevents.CognitoEventUserPoolsCreateChallenge{}, nil
	}

	return createchallenge.HandleCognitoCreateChallengeRequest(event)
}

func main() {
	lambda.Start(HandleCreateChallengeRequest)
}
