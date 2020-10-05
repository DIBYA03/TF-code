package main

import (
	"context"
	"log"
	"os"

	analytics "github.com/wiseco/core-platform/analytics"
	"github.com/wiseco/core-platform/shared"
)

type analyticsService struct{}

func main() {
	region := os.Getenv("SQS_REGION")
	queueURL := os.Getenv("SEGMENT_SQS_URL")
	if queueURL == "" {
		panic("Analytics sqs url missing")
	}

	queue, err := shared.NewSQSMessageQueueFromURL(queueURL, region)
	if err != nil {
		panic(err)
	}

	as := &analyticsService{}

	err = queue.ReceiveMessages(context.Background(), as)
	if err != nil {
		log.Println(err)
	}

}

func (a *analyticsService) HandleMessage(_ context.Context, m shared.Message) error {
	body := string(m.Body[:])
	err := analytics.HandleMessages(&body)
	return err
}
