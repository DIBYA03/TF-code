package main

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	csp "github.com/wiseco/core-platform/api/csp/wise_user"
	"github.com/wiseco/core-platform/api/gateway"
)

func handleConsumerRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	resp, err := csp.PhoneNumberChangeRequest(gateway.NewCSPAPIRequest(request))
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
	lambda.Start(handleConsumerRequest)
}
