package main

import (
	"context"
	"log"
	"os"

	"github.com/wiseco/core-platform/partner/bank/bbva"
	_ "github.com/wiseco/core-platform/partner/bank/bbva"
	"github.com/wiseco/core-platform/shared"
)

func main() {
	// BBVA SQS
	regionIn := os.Getenv("SQS_BBVA_ENV_REGION")
	urlIn := os.Getenv("SQS_BBVA_ENV_URL")
	if urlIn == "" {
		panic("BBVA notification sqs url missing")
	}

	// Acvitity SQS
	regionOut := os.Getenv("SQS_BANKING_REGION")
	urlOut := os.Getenv("SQS_BANKING_URL")
	if urlOut == "" {
		panic("Activity sqs url missing")
	}

	// Stream
	streamRegion := os.Getenv("KINESIS_BANK_NOTIF_REGION")
	if streamRegion == "" {
		panic("Notification stream region missing")
	}

	streamName := os.Getenv("KINESIS_BANK_NOTIF_NAME")
	if streamName == "" {
		panic("Notification stream name missing")
	}

	inQueue, err := shared.NewSQSMessageQueueFromURL(urlIn, regionIn)
	if err != nil {
		panic(err)
	}

	outQueue, err := shared.NewSQSMessageQueueFromURL(urlOut, regionOut)
	if err != nil {
		panic(err)
	}

	info := bbva.NotificationHandlerInfo{
		InQueue:        inQueue,
		OutQueue:       outQueue,
		StreamProvider: shared.NewKinesisStreamProvider(shared.StreamProviderRegion(streamRegion), streamName),
	}

	err = bbva.NewNotificationService(info).HandleNotifications(context.Background())
	if err != nil {
		log.Println("Exiting handle notifications: ", err.Error())
	} else {
		log.Println("Exiting handle notifications")
	}
}
