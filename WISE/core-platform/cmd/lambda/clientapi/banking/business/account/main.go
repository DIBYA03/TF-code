package main

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	busbanking "github.com/wiseco/core-platform/api/client/banking/business"
	"github.com/wiseco/core-platform/api/gateway"
)

func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	resp, err := busbanking.HandleAccountAPIRequest(gateway.NewAPIRequest(request))
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
	lambda.Start(HandleRequest)
}
