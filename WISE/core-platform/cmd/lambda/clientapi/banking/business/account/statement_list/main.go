package main

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	banking "github.com/wiseco/core-platform/api/client/banking/business"
	"github.com/wiseco/core-platform/api/gateway"
)

//HandleRequest handles a transaction request
func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	resp, err := banking.HandleStatementAPIRequest(gateway.NewAPIRequest(request))
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
