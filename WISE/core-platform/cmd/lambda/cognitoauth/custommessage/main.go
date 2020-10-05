package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/wiseco/core-platform/api/auth/customevents"
	"github.com/wiseco/core-platform/api/auth/custommessage"
)

// HandleCustomMessageRequest handle Custom message request
func HandleCustomMessageRequest(ctx context.Context, event customevents.CognitoEventUserPoolsCustomMessage) (customevents.CognitoEventUserPoolsCustomMessage, error) {
	if _, ok := event.Request.UserAttributes["lambda_warmer"]; ok {
		log.Println("lambda warming...")
		return customevents.CognitoEventUserPoolsCustomMessage{}, nil
	}

	return custommessage.HandleCognitoCustomMessageRequest(event)
}

func main() {
	lambda.Start(HandleCustomMessageRequest)
}
