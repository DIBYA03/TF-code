package main

import (
	"context"
	"log"
	"os"

	"github.com/wiseco/core-platform/services/csp/messages"
	"github.com/wiseco/core-platform/shared"
)

type handler struct{}

func main() {
	region := os.Getenv("CSP_REVIEW_SQS_REGION")
	queueURL := os.Getenv("CSP_REVIEW_SQS_URL")
	if queueURL == "" {
		panic("csp  review sqs url missing")
	}

	queue, err := shared.NewSQSMessageQueueFromURL(queueURL, region)
	if err != nil {
		panic(err)
	}

	hd := &handler{}
	err = queue.ReceiveMessages(context.Background(), hd)
	if err != nil {
		log.Println(err)
	}
}

func (h *handler) HandleMessage(ctx context.Context, m shared.Message) error {
	body := string(m.Body)
	return messages.HandleSQS(body)
}
