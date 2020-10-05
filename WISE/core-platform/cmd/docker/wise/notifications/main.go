package main

import (
	"context"
	"log"
	"os"

	"github.com/wiseco/core-platform/notification"
	"github.com/wiseco/core-platform/shared"
)

type notificationService struct{}

func main() {
	region := os.Getenv("SQS_REGION")
	queueURL := os.Getenv("SQS_URL")
	if queueURL == "" {
		panic("Notification sqs url missing")
	}

	queue, err := shared.NewSQSMessageQueueFromURL(queueURL, region)
	if err != nil {
		panic(err)
	}

	// Stream
	streamRegion := os.Getenv("KINESIS_TRX_REGION")
	if streamRegion == "" {
		panic("Notification stream region missing")
	}

	streamName := os.Getenv("KINESIS_TRX_NAME")
	if streamName == "" {
		panic("Notification stream name missing")
	}

	err = queue.ReceiveMessages(context.Background(), &notificationService{})
	if err != nil {
		log.Println(err)
	}
}

func (s *notificationService) HandleMessage(_ context.Context, m shared.Message) error {
	body := string(m.Body[:])
	err := notification.HandleNotification(&body)
	return err
}
