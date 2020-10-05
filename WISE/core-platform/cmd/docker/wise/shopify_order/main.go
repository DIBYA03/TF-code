package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/wiseco/core-platform/services"
	"github.com/wiseco/core-platform/services/order"
	"github.com/wiseco/core-platform/shared"
)

type shopifyOrderService struct{}

func main() {
	log.Println("starting with main")

	region := os.Getenv("SQS_REGION")
	queueURL := os.Getenv("SHOPIFY_ORDER_SQS_URL")
	if queueURL == "" {
		panic("shopify order sqs url missing")
	}

	queue, err := shared.NewSQSMessageQueueFromURL(queueURL, region)
	if err != nil {
		panic(err)
	}

	ss := &shopifyOrderService{}
	err = queue.ReceiveMessages(context.Background(), ss)
	if err != nil {
		log.Println(err)
	}
}

func (ss *shopifyOrderService) HandleMessage(_ context.Context, m shared.Message) error {
	body := string(m.Body[:])
	msg := order.ShopifyOrderMessage{}

	err := json.Unmarshal([]byte(body), &msg)
	if err != nil {
		return err
	}

	return order.NewShopifyOrderService(services.NewSourceRequest()).HandleShopifyOrder(msg)
}
