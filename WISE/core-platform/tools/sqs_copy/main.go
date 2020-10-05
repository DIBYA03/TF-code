package main

import (
	"context"
	"log"

	"github.com/wiseco/core-platform/shared"
)

func main() {
	urlIn := "https://sqs.us-west-2.amazonaws.com/058450407364/prd-client-api-bbva-notifications-dead-letter"
	regionIn := "us-west-2"
	inq, err := shared.NewSQSMessageQueueFromURL(urlIn, regionIn)
	if err != nil {
		panic(err)
	}

	urlOut := "https://sqs.us-west-2.amazonaws.com/058450407364/prd-client-api-bbva-notifications"
	regionOut := "us-west-2"
	outq, err := shared.NewSQSMessageQueueFromURL(urlOut, regionOut)
	if err != nil {
		panic(err)
	}

	h := handlerInfo{
		outq: outq,
	}

	err = inq.ReceiveMessages(context.Background(), &h)
	if err != nil {
		panic(err)
	}
}

type handlerInfo struct {
	outq shared.MessageQueue
}

func (h *handlerInfo) HandleMessage(_ context.Context, m shared.Message) error {
	_, err := h.outq.SendMessages([]shared.Message{m})
	if err != nil {
		log.Println(err)
	}

	return err
}
