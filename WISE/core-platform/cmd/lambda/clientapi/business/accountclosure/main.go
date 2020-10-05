package main

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/wiseco/core-platform/api/client/business"
	"github.com/wiseco/core-platform/api/gateway"
)

//HandleAccountClosureRequest handles business activity request
func HandleAccountClosureRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	resp, err := business.HandleAccountClosurAPIRequests(gateway.NewAPIRequest(request))
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
	lambda.Start(HandleAccountClosureRequest)
}
