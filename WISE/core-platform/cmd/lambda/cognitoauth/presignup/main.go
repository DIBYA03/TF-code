package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/wiseco/core-platform/api/auth/presignup"
)

func HandleRequest(ctx context.Context, event events.CognitoEventUserPoolsPreSignup) (events.CognitoEventUserPoolsPreSignup, error) {
	if _, ok := event.Request.UserAttributes["lambda_warmer"]; ok {
		log.Println("lambda warming...")
		return events.CognitoEventUserPoolsPreSignup{}, nil
	}

	return presignup.HandleCognitoPreSignupRequest(event)
}

func main() {
	lambda.Start(HandleRequest)
}
