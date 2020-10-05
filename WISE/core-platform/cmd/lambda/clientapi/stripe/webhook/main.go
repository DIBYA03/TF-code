package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/wiseco/core-platform/api"
	stripe "github.com/wiseco/core-platform/api/client/webhook"
)

func HandleRequest(ctx context.Context, sqsEvent events.SQSEvent) error {

	for _, message := range sqsEvent.Records {
		signature := message.MessageAttributes["Stripe-Signature"].StringValue

		headers := map[string]string{
			"Stripe-Signature": *signature,
		}

		request := api.APIRequest{
			UserID:     "", //User id will be empty
			StartedAt:  time.Now(),
			RequestID:  message.MessageId,
			UserAgent:  *message.MessageAttributes["User-Agent"].StringValue,
			PoolID:     os.Getenv("CLIENT_API_POOL_ID"),
			Body:       message.Body,
			Headers:    headers,
			HTTPMethod: http.MethodPost,
			SourceIP:   *message.MessageAttributes["Source-IP"].StringValue,
		}

		err := stripe.HandleStripeRequest(request)
		if err != nil {
			return err
		}

	}

	return nil
}

func main() {
	lambda.Start(HandleRequest)
}
