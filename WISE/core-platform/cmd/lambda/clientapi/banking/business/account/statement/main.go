package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	banking "github.com/wiseco/core-platform/api/client/banking/business"
	"github.com/wiseco/core-platform/api/gateway"
)

//HandleRequest handles a transaction request
func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (string, error) {
	resp, err := banking.HandleStatementAPIRequest(gateway.NewAPIRequest(request))
	if err != nil {
		return "", err
	}

	return resp.Body, nil
}

func main() {
	lambda.Start(HandleRequest)
}
