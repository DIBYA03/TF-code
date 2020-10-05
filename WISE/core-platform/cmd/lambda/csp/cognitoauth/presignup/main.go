package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/wiseco/core-platform/api/csp/auth/presignup"
)

func HandleRequest(ctx context.Context, event events.CognitoEventUserPoolsPreSignup) (events.CognitoEventUserPoolsPreSignup, error) {
	return presignup.HandleCognitoPreSignupRequest(event)
}

func main() {
	lambda.Start(HandleRequest)
}
