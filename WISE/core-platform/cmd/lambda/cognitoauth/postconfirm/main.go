package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/wiseco/core-platform/api/auth/postconfirm"
)

func HandleRequest(ctx context.Context, event events.CognitoEventUserPoolsPostConfirmation) (events.CognitoEventUserPoolsPostConfirmation, error) {
	if _, ok := event.Request.UserAttributes["lambda_warmer"]; ok {
		log.Println("lambda warming...")
		return events.CognitoEventUserPoolsPostConfirmation{}, nil
	}

	return postconfirm.HandleCognitoPostConfirmRequest(event)
}

func main() {
	lambda.Start(HandleRequest)
}
