package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	csp "github.com/wiseco/core-platform/services/csp/document"
)

// HandleRequest is the lambda entrypoint
func HandleRequest(ctx context.Context, s3Event events.S3Event) error {

	// Download and scan the files that triggered the lambda
	for _, record := range s3Event.Records {
		csp.HandleS3Event(record.S3.Object.Key)
	}
	return nil
}
func main() {
	lambda.Start(HandleRequest)
}
