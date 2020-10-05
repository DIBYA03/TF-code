package main

import (
	"fmt"
	"log"
	"os"

	"github.com/wiseco/core-platform/partner/bank"
	"github.com/wiseco/core-platform/partner/bank/bbva/data"
)

func main() {
	url := os.Getenv("SQS_BBVA_URL")
	if url == "" {
		err := fmt.Errorf("data notification sqs url missing")
		log.Println(err)
		panic(err)
	}

	srv := data.GetSubscriptionService(bank.NewAPIRequest())
	for _, configTypes := range data.AllEventConfigTypes {
		_, err := srv.Subscribe(data.SubscriptionChannelTypeSQS, data.ChannelURL(url), configTypes)
		if err != nil {
			log.Println(err.Error())
			panic(err.Error())
		}
	}

	log.Println("Completed subscribing to notifications")
}
