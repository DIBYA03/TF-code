package main

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	csp "github.com/wiseco/core-platform/api/csp/account_closure"
	"github.com/wiseco/core-platform/api/gateway"
)

func handleAccountClosureRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	resp, err := csp.HandleCSPAccountClosureStateAPIRequests(gateway.NewCSPAPIRequest(request))
	if err != nil {
		return events.APIGatewayProxyResponse{
				StatusCode: http.StatusInternalServerError,
				Headers: map[string]string{
					"Content-Type": "application/json",
				},
			},
			nil
	}

	return gateway.ProxyResponse(resp), nil
}

func main() {
	lambda.Start(handleAccountClosureRequest)
}
