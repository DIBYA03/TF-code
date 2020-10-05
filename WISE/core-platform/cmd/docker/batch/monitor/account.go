package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	grpcTypes "github.com/golang/protobuf/ptypes"
	"github.com/google/uuid"
	"github.com/wiseco/core-platform/services"
	"github.com/wiseco/core-platform/services/banking/business"
	bsrv "github.com/wiseco/core-platform/services/business"
	"github.com/wiseco/core-platform/services/data"
	usrv "github.com/wiseco/core-platform/services/user"
	"github.com/wiseco/core-platform/shared"
	"github.com/wiseco/go-lib/id"
	grpcBanking "github.com/wiseco/protobuf/golang/banking"
	grpcMonitor "github.com/wiseco/protobuf/golang/transaction/monitor"
)

func processAccount(monitorClient grpcMonitor.BankTransactionMonitorServiceClient, aID id.BankAccountID) error {
	a, err := business.NewAccountService().GetByIDInternal(aID.UUIDString())
	if err != nil {
		log.Println(err, aID.String())
		return err
	}

	b, err := bsrv.NewBusinessServiceWithout().GetByIdInternal(a.BusinessID)
	if err != nil {
		log.Println(err, aID.String())
		return err
	}

	if b.KYCStatus != services.KYCStatusApproved {
		err := fmt.Errorf("kyc status not approved: %s", b.KYCStatus)
		log.Println(err, aID.String())
		return err
	}

	created, err := grpcTypes.TimestampProto(a.Created)
	if err != nil {
		log.Println(err, aID.String())
		return err
	}

	modified, err := grpcTypes.TimestampProto(a.Modified)
	if err != nil {
		log.Println(err, aID.String())
		return err
	}

	status := grpcBanking.AccountStatus_AS_ACTIVE
	bUID, err := uuid.Parse(string(a.BusinessID))
	if err != nil {
		log.Println(err, aID.String())
		return err
	}

	bID := id.BusinessID(bUID)

	usr, err := usrv.NewUserServiceWithout().GetByIdInternal(a.AccountHolderID)
	if err != nil {
		log.Println(err, aID.String())
		return err
	}

	conUID, err := uuid.Parse(string(usr.ConsumerID))
	if err != nil {
		log.Println(err, aID.String())
		return err
	}

	conID := id.ConsumerID(conUID)
	creq := &grpcMonitor.BankAccountRequest{
		Id:                  aID.String(),
		BusinessId:          bID.String(),
		AdditionalConsumers: []string{conID.String()},
		Type:                grpcBanking.AccountType_AT_BUSINESS,
		PartnerType:         grpcBanking.PartnerAccountType_PAT_DEPOSITORY_CHECKING,
		Subtype:             grpcBanking.AccountSubtype_AST_PRIMARY,
		AccountAlias:        shared.StringValue(a.Alias),
		AvailableBalance:    strconv.FormatFloat(a.AvailableBalance, 'f', 2, 64),
		PostedBalance:       strconv.FormatFloat(a.PostedBalance, 'f', 2, 64),
		Currency:            string(a.Currency),
		Status:              status,
		Created:             created,
		Modified:            modified,
	}

	resp, err := monitorClient.AddUpdateAccount(context.Background() /* client.GetContext() */, creq)
	if err != nil {
		log.Println(err, aID.String())
	} else {
		log.Println("Success: ", resp.Id)
	}

	return err
}

func sendAccountUpdates(service grpcMonitor.BankTransactionMonitorServiceClient, dayStart, dayEnd time.Time) {
	// Process in groups of 5
	offset := 0
	limit := 5
	for {
		var accountIDs []id.BankAccountID
		err := data.DBWrite.Select(
			&accountIDs,
			`
			SELECT id from business_bank_account
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
		} else if len(accountIDs) == 0 {
			log.Println("No more accounts", dayStart, dayEnd)
			break
		}

		wg := sync.WaitGroup{}
		wg.Add(len(accountIDs))
		for _, aID := range accountIDs {
			go func(id id.BankAccountID) {
				defer wg.Done()
				_ = processAccount(service, id)
			}(aID)
		}

		wg.Wait()
		offset += 5
	}
}
