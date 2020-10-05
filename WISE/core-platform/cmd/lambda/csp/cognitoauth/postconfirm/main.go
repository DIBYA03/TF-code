package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/wiseco/core-platform/api/csp/auth/postconfirm"
)

func HandleRequest(ctx context.Context, event events.CognitoEventUserPoolsPostConfirmation) (events.CognitoEventUserPoolsPostConfirmation, error) {
	return postconfirm.HandleCognitoPostConfirmRequest(event)
}

func main() {
	lambda.Start(HandleRequest)
}
