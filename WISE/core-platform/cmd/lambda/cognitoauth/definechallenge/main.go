package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/wiseco/core-platform/api/auth/customevents"
	"github.com/wiseco/core-platform/api/auth/definechallenge"
)

// HandleDefineChallengeRequest handle define challenge req
func HandleDefineChallengeRequest(ctx context.Context, event customevents.CognitoEventUserPoolsDefineChallenge) (customevents.CognitoEventUserPoolsDefineChallenge, error) {
	if _, ok := event.Request.UserAttributes["lambda_warmer"]; ok {
		log.Println("lambda warming...")
		return customevents.CognitoEventUserPoolsDefineChallenge{}, nil
	}

	return definechallenge.HandleCognitoDefineChallengeRequest(event)
}

func main() {
	lambda.Start(HandleDefineChallengeRequest)
}
