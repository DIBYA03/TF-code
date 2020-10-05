package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	grpcTypes "github.com/golang/protobuf/ptypes"
	"github.com/wiseco/core-platform/services/banking/business"
	"github.com/wiseco/core-platform/shared"
	"github.com/wiseco/core-platform/transaction"
	txndata "github.com/wiseco/core-platform/transaction"
	"github.com/wiseco/go-lib/grpc"
	"github.com/wiseco/go-lib/id"
	"github.com/wiseco/go-lib/num"
	grpcBankTxn "github.com/wiseco/protobuf/golang/transaction/bank"
)

var dayStart, dayEnd, dayStartLocal, dayEndLocal time.Time
var txnClient grpcBankTxn.TransactionServiceClient

func main() {
	// Get time zone and determine start/end
	tz := os.Getenv("BATCH_TZ")
	if tz == "" {
		panic(errors.New("Local timezone missing"))
	}

	loc, err := time.LoadLocation(tz)
	if err != nil {
		panic(err)
	}

	utcLoc, err := time.LoadLocation("UTC")
	if err != nil {
		panic(err)
	}

	nowUTC := time.Now().UTC()
	nowLocal := nowUTC.In(loc)

	dayEndLocal = time.Date(nowLocal.Year(), nowLocal.Month(), nowLocal.Day(), 0, 0, 0, 0, loc)
	dayStartLocal = dayEndLocal.AddDate(0, 0, -1)

	dayStart = dayStartLocal.In(utcLoc)
	dayEnd = dayEndLocal.In(utcLoc)
	if dayEnd.After(nowUTC) {
		panic(fmt.Errorf("Error: day end (%v) is after current time (%v)", dayEnd, nowUTC))
	}

	dayStart = dayEnd.AddDate(0, 0, -1)
	dayStartLocal = dayStart.In(loc)

	if os.Getenv("USE_TRANSACTION_SERVICE") == "true" {
		sn, err := grpc.GetConnectionStringForService(grpc.ServiceNameTransaction)
		if err != nil {
			panic(err)
		}

		cl, err := grpc.NewInsecureClient(sn)
		if err != nil {
			panic(err)
		}

		defer cl.CloseAndCancel()
		txnClient = grpcBankTxn.NewTransactionServiceClient(cl.GetConn())
	}

	// Fetch bank accounts in groups of 10
	limit := 10
	offset := 0
	for {
		accounts, err := business.NewAccountService().ListInternal(limit, offset)
		if err != nil {
			panic(err)
		}

		if len(accounts) == 0 {
			log.Println("No more accounts")
			break
		}

		offset += limit

		// TODO: Does SQS or transactioning per action work better?
		var wait sync.WaitGroup
		wait.Add(len(accounts))
		for _, account := range accounts {
			go func(account *business.BankAccount) {
				defer wait.Done()

				// Balance Update
				err := updateDailyAccountBalance(account)
				if err != nil {
					log.Println("updateDailyAccountBalance:", err)
				}
			}(account)
		}

		wait.Wait()
	}
}

func updateDailyAccountBalance(account *business.BankAccount) error {
	zb, err := shared.NewDecimalFin(0)
	if err != nil {
		return err
	}

	zero := shared.Decimal{V: zb}
	bal, err := business.NewAccountService().GetBalanceByIDInternal(account.Id)
	if err != nil {
		log.Println("GetBalanceByIDInternal:", err)
		return err
	}

	postedBal, err := num.NewFromFloat(bal.PostedBalance)
	if err != nil {
		log.Println("NewFromFloat:", err)
		return err

	}

	// Get total credits and debits
	var debitCreditSummary struct {
		AmountCredited shared.Decimal `db:"amount_credited"`
		AmountDebited  shared.Decimal `db:"amount_debited"`
	}

	s := account.Id
	if !strings.HasPrefix(account.Id, id.IDPrefixBankAccount.String()) {
		s = fmt.Sprintf("%s%s", id.IDPrefixBankAccount, account.Id)
	}

	baID, err := id.ParseBankAccountID(s)
	if err != nil {
		return err
	}

	if os.Getenv("USE_TRANSACTION_SERVICE") == "true" {
		start, err := grpcTypes.TimestampProto(dayStart)
		if err != nil {
			return err
		}

		end, err := grpcTypes.TimestampProto(dayEnd)
		if err != nil {
			return err
		}

		req := &grpcBankTxn.StatsRequest{
			AccountId: baID.String(),
			DateRange: &grpcBankTxn.DateRange{
				Filter: grpcBankTxn.DateRangeFilter_DRF_START_END,
				Start:  start,
				End:    end,
			},
		}

		resp, err := txnClient.GetStats(context.Background(), req)
		if err != nil {
			return err
		}

		debitCreditSummary.AmountCredited, err = shared.ParseDecimal(resp.AmountCredited)
		if err != nil {
			return err
		}

		debitCreditSummary.AmountDebited, err = shared.ParseDecimal(resp.AmountDebited)
		if err != nil {
			return err
		}
	} else {
		err = txndata.DBRead.Get(
			&debitCreditSummary,
			`
		SELECT
			sum(case when code_type = $1 then amount else 0 end) as amount_credited,
			sum(case when code_type = $2 then amount else 0 end) as amount_debited
		FROM business_transaction
		WHERE
			account_id = $3 AND
			transaction_date >= $4 AND
			transaction_date < $5
		GROUP BY account_id`,
			transaction.TransactionCodeTypeCreditPosted,
			transaction.TransactionCodeTypeDebitPosted,
			account.Id,
			dayStart,
			dayEnd,
		)
		if err != nil && err != sql.ErrNoRows {
			return err
		}
	}

	if debitCreditSummary.AmountCredited.IsNil() {
		debitCreditSummary.AmountCredited = zero
	}

	if debitCreditSummary.AmountDebited.IsNil() {
		debitCreditSummary.AmountDebited = zero
	}

	s = string(account.BusinessID)
	if !strings.HasPrefix(s, id.IDPrefixBusiness.String()) {
		s = fmt.Sprintf("%s%s", id.IDPrefixBusiness, s)
	}

	bID, err := id.ParseBusinessID(s)
	if err != nil {
		return err
	}

	// Create entry in daily balance table
	c := transaction.BusinessAccountDailyBalanceCreate{
		AccountID:      transaction.AccountID(baID.UUIDString()),
		BusinessID:     shared.BusinessID(bID.UUIDString()),
		PostedBalance:  postedBal,
		Currency:       transaction.Currency(bal.Currency),
		AmountCredited: debitCreditSummary.AmountCredited.NumDecimal(),
		AmountDebited:  debitCreditSummary.AmountDebited.NumDecimal(),
		APR:            99,
		RecordedDate:   shared.Date(dayStartLocal),
	}
	_, err = transaction.CreateDailyBalance(&c)
	if err != nil {
		log.Println("CreateDailyBalance:", bal.AccountID)
	} else {
		log.Println("Daily Balannce Complete:", bal.AccountID)
	}

	return err
}
