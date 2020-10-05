package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/wiseco/core-platform/api/auth/pretoken"
)

// HandleRequest Handle presignup
func HandleRequest(ctx context.Context, event events.CognitoEventUserPoolsPreTokenGen) (events.CognitoEventUserPoolsPreTokenGen, error) {
	if _, ok := event.Request.UserAttributes["lambda_warmer"]; ok {
		log.Println("lambda warming...")
		return events.CognitoEventUserPoolsPreTokenGen{}, nil
	}

	return pretoken.HandleCognitoPreTokenRequest(event)
}

func main() {
	lambda.Start(HandleRequest)
}
