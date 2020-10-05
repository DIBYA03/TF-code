package main

import (
	"log"

	"github.com/wiseco/core-platform/partner/bank"
	"github.com/wiseco/core-platform/partner/bank/bbva/data"
)

func main() {
	err := data.GetSubscriptionService(bank.NewAPIRequest()).UnsubscribeAll()
	if err != nil {
		log.Println(err)
		panic(err)
	}

	log.Println("Completed unsubscribing to BBVA notifications")
}
