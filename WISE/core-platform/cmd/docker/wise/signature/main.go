package main

import (
	"context"
	"log"
	"os"

	"github.com/wiseco/core-platform/services/signature"
	"github.com/wiseco/core-platform/shared"
)

type signatureService struct{}

func main() {
	log.Println("starting with main")

	region := os.Getenv("SQS_REGION")
	queueURL := os.Getenv("SIGNATURE_SQS_URL")
	if queueURL == "" {
		panic("Signature sqs url missing")
	}

	queue, err := shared.NewSQSMessageQueueFromURL(queueURL, region)
	if err != nil {
		panic(err)
	}

	ss := &signatureService{}
	err = queue.ReceiveMessages(context.Background(), ss)
	if err != nil {
		log.Println(err)
	}
}

func (ss *signatureService) HandleMessage(_ context.Context, m shared.Message) error {
	body := string(m.Body[:])
	err := signature.HandleMessage(body)
	return err
}
