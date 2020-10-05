package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/wiseco/core-platform/api/csp/auth/preauthentication"
)

// HandleRequest Handle presignup
func HandleRequest(ctx context.Context, event events.CognitoEventUserPoolsPreAuthentication) (events.CognitoEventUserPoolsPreAuthentication, error) {
	return preauthentication.HandleCognitoPreAuthenticationRequest(event)
}

func main() {
	lambda.Start(HandleRequest)
}
