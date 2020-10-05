package main

import (
	"context"
	"log"
	"sync"
	"time"

	grpcTypes "github.com/golang/protobuf/ptypes"
	"github.com/google/uuid"
	"github.com/wiseco/core-platform/services/banking/business"
	"github.com/wiseco/core-platform/services/data"
	"github.com/wiseco/core-platform/services/user"
	"github.com/wiseco/go-lib/id"
	grpcBanking "github.com/wiseco/protobuf/golang/banking"
	grpcMonitor "github.com/wiseco/protobuf/golang/transaction/monitor"
)

func processDebitCard(monitorClient grpcMonitor.BankTransactionMonitorServiceClient, cardID id.DebitCardID) error {
	c, err := business.NewCardServiceWithout().GetByIDInternal(uuid.UUID(cardID).String())
	if err != nil {
		log.Println(err, cardID)
		return err
	}

	created, err := grpcTypes.TimestampProto(c.Created)
	if err != nil {
		log.Println(err, cardID)
		return err
	}

	modified, err := grpcTypes.TimestampProto(c.Modified)
	if err != nil {
		log.Println(err, cardID)
		return err
	}

	status := grpcBanking.DebitCardStatus_DCS_ACTIVE
	bUID, err := uuid.Parse(string(c.BusinessID))
	if err != nil {
		log.Println(err, cardID)
		return err
	}

	bID := id.BusinessID(bUID)
	if bID.IsZero() {
		log.Println("Business ID zero:", cardID)
		return nil
	}

	accUUID, err := uuid.Parse(string(c.BankAccountId))
	if err != nil {
		log.Println(err, cardID)
		return err
	}

	accID := id.BankAccountID(accUUID)

	u, err := user.NewUserServiceWithout().GetByIdInternal(c.CardholderID)
	if err != nil {
		log.Println(err, cardID)
		return err
	}

	conUUID, err := uuid.Parse(string(u.ConsumerID))
	if err != nil {
		log.Println(err, cardID)
		return err
	}

	conID := id.ConsumerID(conUUID)

	last4, err := c.GetCardNumberLastFour()
	if err != nil {
		log.Println(err, cardID)
		return err
	}

	creq := &grpcMonitor.DebitCardRequest{
		Id:         cardID.String(),
		ConsumerId: conID.String(),
		AccountId:  accID.String(),
		CardLast_4: last4,
		Status:     status,
		Created:    created,
		Modified:   modified,
	}

	resp, err := monitorClient.AddUpdateDebitCard(context.Background() /* client.GetContext() */, creq)
	if err != nil {
		log.Println(err, cardID)
		return err
	}

	log.Println("Success: ", resp.Id)
	return nil
}

func sendDebitCardUpdates(monitorClient grpcMonitor.BankTransactionMonitorServiceClient, dayStart, dayEnd time.Time) {
	// Process in groups of 5
	offset := 0
	limit := 5
	for {
		var cardIDs []id.DebitCardID
		err := data.DBWrite.Select(
			&cardIDs,
			`
			SELECT id from business_bank_card
			WHERE
				(created >= $1 AND created < $2) OR
				(modified >= $1 AND modified < $2)
			ORDER BY created ASC OFFSET $3 LIMIT $4`,
			dayStart,
			dayEnd,
			offset,
			limit,
		)
		if err != nil {
			panic(err)
		} else if len(cardIDs) == 0 {
			log.Println("No more debit cards", dayStart, dayEnd)
			break
		}

		wg := sync.WaitGroup{}
		wg.Add(len(cardIDs))
		for _, cardID := range cardIDs {
			go func(id id.DebitCardID) {
				defer wg.Done()
				_ = processDebitCard(monitorClient, id)
			}(cardID)
		}

		wg.Wait()
		offset += 5
	}
}
