package main

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	csp "github.com/wiseco/core-platform/api/csp/transaction"
	"github.com/wiseco/core-platform/api/gateway"
)

func handleCSPTransactionRequests(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	resp, err := csp.HandlePostedTransactionExportRequest(gateway.NewCSPAPIRequest(request))
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
	lambda.Start(handleCSPTransactionRequests)
}
