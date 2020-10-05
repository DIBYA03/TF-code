package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	csp "github.com/wiseco/core-platform/api/csp/business"
	"github.com/wiseco/core-platform/api/gateway"
)

func handleCSPBusinessRequests(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Println("In handler")
	resp, err := csp.BankAccountRequest(gateway.NewCSPAPIRequest(request))
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
	lambda.Start(handleCSPBusinessRequests)
}
