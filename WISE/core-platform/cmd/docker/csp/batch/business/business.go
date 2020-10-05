package main

import (
	"log"
	"sync"

	"github.com/wiseco/core-platform/services"
	"github.com/wiseco/core-platform/services/csp"
	"github.com/wiseco/core-platform/services/csp/bot"
	cspDB "github.com/wiseco/core-platform/services/csp/data"
	"github.com/wiseco/core-platform/services/csp/mail"
	"github.com/wiseco/core-platform/services/csp/review"
	"github.com/wiseco/core-platform/shared"
)

type business struct {
	BusinessID shared.BusinessID `db:"business_id"`
	Name       string            `db:"business_name"`
	Status     csp.Status        `db:"review_status"`
}

func fetchBusinesses() error {
	rows, err := cspDB.DBRead.Queryx("SELECT business_id, business_name, review_status FROM business WHERE review_status = 'bankReview'")
	if err != nil {
		log.Printf("error fetching businesess to check for kyc status %v", err)
		return err
	}
	defer rows.Close()
	var wg sync.WaitGroup
	for rows.Next() {
		var bus business
		if err := rows.StructScan(&bus); err != nil {
			log.Printf("error scanning row %v", err)
			return err
		}
		// add 1 to wait group so we can run multiple business at once
		wg.Add(1)
		go checkStatus(bus, &wg)
	}
	wg.Wait()
	return nil
}

func checkStatus(b business, wg *sync.WaitGroup) error {
	defer wg.Done()
	log.Printf("Checking status on busines %s", b.Name)
	resp, err := review.New(services.NoToCSPServiceRequest(shared.UserIDEmpty)).GetStatus(b.BusinessID)
	if err != nil {
		log.Printf("erro checking status %v", err)
		return err
	}
	if resp == nil {
		log.Print("no response")
		return nil
	}
	status := csp.KYCStatus(resp.Status)
	notify(b.Name, status)
	// SQS message to create account and card
	// temporarly disabled
	/*
		if status == csp.KYCStatusApproved {
			csp.SendReviewMessage(csp.Message{
				EntityID: string(b.BusinessID),
				Category: csp.CategoryAccount,
				Action:   csp.ActionCreate,
			})
		}
	*/
	return err
}

func notify(business string, status csp.KYCStatus) {
	team := []mail.WiseTeam{
		mail.WiseTeam{
			Name:  "Arjun",
			Email: "arjun@wise.us",
		},
		mail.WiseTeam{
			Name:  "Josh",
			Email: "josh@wise.us",
		},
	}
	switch status {
	case csp.KYCStatusApproved:
		mail.EmailWiseTeam(business, csp.KYCStatusApproved, team)
		bot.SendNotification(bot.StatusApproved, business)
	case csp.KYCStatusDeclined:
		mail.EmailWiseTeam(business, csp.KYCStatusDeclined, team)
		bot.SendNotification(bot.StatusDeclined, business)
	}

}
