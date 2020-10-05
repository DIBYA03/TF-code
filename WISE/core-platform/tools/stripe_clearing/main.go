package main

import (
	"log"
	"os"

	"github.com/wiseco/core-platform/services"
	"github.com/wiseco/core-platform/services/banking"
	"github.com/wiseco/core-platform/services/banking/business"
	"github.com/wiseco/core-platform/shared"
)

func main() {
	uid := os.Getenv("WISE_CLEARING_USER_ID")
	if uid == "" {
		log.Fatal("user id required")
	}

	bid := os.Getenv("WISE_CLEARING_BUSINESS_ID")
	if bid == "" {
		log.Fatal("business id required")
	}

	ac := business.BankAccountCreate{
		shared.BusinessID(bid),
		business.UsageTypeClearing,
		banking.BankAccountCreate{
			BankName:    banking.BankNameBBVA,
			AccountType: "checking",
		},
	}

	acc, err := business.NewBankAccountService(services.NewSourceRequest()).Create(&ac)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Clearing Account: ", acc.Id)

	la, err := business.NewLinkedAccountService(services.NewSourceRequest()).GetByAccountIDInternal(shared.BusinessID(bid), acc.Id)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Clearing Linked Account: ", la.Id)
}
