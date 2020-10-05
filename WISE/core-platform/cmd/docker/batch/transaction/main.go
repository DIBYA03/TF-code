package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/wiseco/core-platform/services/invoice"

	grpcTypes "github.com/golang/protobuf/ptypes"
	"github.com/wiseco/core-platform/services/banking"
	coredata "github.com/wiseco/core-platform/services/data"
	"github.com/wiseco/core-platform/shared"
	"github.com/wiseco/core-platform/transaction"
	txndata "github.com/wiseco/core-platform/transaction"
	"github.com/wiseco/go-lib/grpc"
	"github.com/wiseco/go-lib/id"
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

	// Update daily business data
	var count int
	row := coredata.DBRead.QueryRow("SELECT COUNT(*) FROM business")
	err = row.Scan(&count)
	if err != nil {
		panic(err)
	}

	// Fetch businesses in groups of 10
	limit := 10
	offset := 0
	for {
		var bIDs []shared.BusinessID
		err := coredata.DBRead.Select(
			&bIDs,
			`SELECT id FROM business ORDER BY id ASC OFFSET $1 LIMIT $2`,
			offset,
			limit,
		)
		if err != nil {
			panic(err)
		}

		if len(bIDs) == 0 {
			break
		}

		offset += limit

		var wait sync.WaitGroup
		wait.Add(len(bIDs))
		for _, bID := range bIDs {
			go func(bID shared.BusinessID) {
				defer wait.Done()

				err := createDailyTransactionStats(bID)
				if err != nil {
					log.Println(err)
				}
			}(bID)
		}

		wait.Wait()
	}

}

func createDailyTransactionStats(bID shared.BusinessID) error {
	zb, err := shared.NewDecimalFin(0)
	if err != nil {
		return err
	}

	zero := shared.Decimal{V: zb}

	// Get amount sent, amount requested and amount paid
	var amountSentReceived struct {
		AmountSent      *shared.Decimal `db:"amount_sent"`
		AmountRequested *shared.Decimal `db:"amount_requested"`
		AmountPaid      *shared.Decimal `db:"amount_paid"`
	}

	// Calculate using fixed precision values - doubles are approxomations
	if os.Getenv("USE_INVOICE_SERVICE") == "true" {
		// first get the POS data
		err = coredata.DBRead.Get(
			&amountSentReceived,
			`
		SELECT
			CAST(SUM(amount) AS DECIMAL(19,2)) AS amount_requested,
			CAST(SUM(case when request_status = $1 then amount else 0 end) AS DECIMAL(19,2)) AS amount_paid
		FROM business_money_request
		WHERE
			business_id = $2 AND
			created >= $3 AND
			created < $4 AND
            request_type = 'pos'
		GROUP BY business_id`,
			banking.MoneyRequestStatusComplete,
			bID,
			dayStart,
			dayEnd,
		)
		if err != nil && err != sql.ErrNoRows {
			return err
		}

		// then get the invoice data from invoice service
		invSvc, err := invoice.NewInvoiceService()
		if err != nil {
			log.Println(err)
			return err
		}
		amountResp, err := invSvc.GetInvoiceAmountsWithFilter(bID, dayStart, dayEnd)
		if err != nil {
			log.Println(err)
			return err
		}

		if amountResp.TotalRequest.IsNil() {
			amountResp.TotalRequest = shared.NewZero()
		}

		// lastly, make a sum of the number
		if !amountSentReceived.AmountRequested.IsNil() {
			totalRequested := amountResp.TotalRequest.Add(*amountSentReceived.AmountRequested)
			amountSentReceived.AmountRequested = &totalRequested
		} else {
			amountSentReceived.AmountRequested = &amountResp.TotalRequest
		}

		if amountResp.TotalPaid.IsNil() {
			amountResp.TotalPaid = shared.NewZero()
		}

		if !amountSentReceived.AmountPaid.IsNil() {
			totalPaid := amountResp.TotalPaid.Add(*amountSentReceived.AmountPaid)
			amountSentReceived.AmountPaid = &totalPaid
		} else {
			amountSentReceived.AmountPaid = &amountResp.TotalPaid
		}

	} else {
		err = coredata.DBRead.Get(
			&amountSentReceived,
			`
		SELECT
			CAST(SUM(amount) AS DECIMAL(19,2)) AS amount_requested,
			CAST(SUM(case when request_status = $1 then amount else 0 end) AS DECIMAL(19,2)) AS amount_paid
		FROM business_money_request
		WHERE
			business_id = $2 AND
			created >= $3 AND
			created < $4
		GROUP BY business_id`,
			banking.MoneyRequestStatusComplete,
			bID,
			dayStart,
			dayEnd,
		)
		if err != nil && err != sql.ErrNoRows {
			return err
		}
	}

	// TODO: Get aggregate data from banking service
	err = coredata.DBRead.Get(
		&amountSentReceived,
		`
    SELECT
        CAST(SUM(amount) AS DECIMAL(19,2)) AS amount_sent
    FROM business_money_transfer
    WHERE
        business_id = $1 AND
        created >= $2 AND
        created < $3
    GROUP BY business_id`,
		bID,
		dayStart,
		dayEnd,
	)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if amountSentReceived.AmountSent.IsNil() {
		amountSentReceived.AmountSent = &zero
	}

	if amountSentReceived.AmountRequested.IsNil() {
		amountSentReceived.AmountRequested = &zero
	}

	if amountSentReceived.AmountPaid.IsNil() {
		amountSentReceived.AmountPaid = &zero
	}

	// Get total credits and debits
	var amountDebitedCredited struct {
		AmountCredited shared.Decimal `db:"amount_credited"`
		AmountDebited  shared.Decimal `db:"amount_debited"`
	}

	// Calculate using fixed precision values - doubles are approxomations
	if os.Getenv("USE_TRANSACTION_SERVICE") == "true" {
		busUUID, err := id.ParseUUID(string(bID))
		if err != nil {
			return err
		}

		busID := id.BusinessID(busUUID)

		start, err := grpcTypes.TimestampProto(dayStart)
		if err != nil {
			return err
		}

		end, err := grpcTypes.TimestampProto(dayEnd)
		if err != nil {
			return err
		}

		req := &grpcBankTxn.StatsRequest{
			BusinessId: busID.String(),
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

		amountDebitedCredited.AmountCredited, err = shared.ParseDecimal(resp.AmountCredited)
		if err != nil {
			return err
		}

		amountDebitedCredited.AmountDebited, err = shared.ParseDecimal(resp.AmountDebited)
		if err != nil {
			return err
		}
	} else {
		err = txndata.DBRead.Get(
			&amountDebitedCredited,
			`
		SELECT
			CAST(SUM(case when code_type = $1 then amount else 0 end) AS DECIMAL(19,2)) AS amount_credited,
			CAST(SUM(case when code_type = $2 then amount else 0 end) AS DECIMAL(19,2)) AS amount_debited
		FROM business_transaction
		WHERE
			business_id = $3 AND
			transaction_date >= $4 AND
			transaction_date < $5
		GROUP BY business_id`,
			transaction.TransactionCodeTypeCreditPosted,
			transaction.TransactionCodeTypeDebitPosted,
			bID,
			dayStart,
			dayEnd,
		)
		if err != nil && err != sql.ErrNoRows {
			return err
		}
	}

	if amountDebitedCredited.AmountCredited.IsNil() {
		amountDebitedCredited.AmountCredited = zero
	}

	if amountDebitedCredited.AmountDebited.IsNil() {
		amountDebitedCredited.AmountDebited = zero
	}

	// Create entry in daily transaction stats table
	c := transaction.BusinessDailyTransactionCreate{
		BusinessID:      bID,
		AmountSent:      amountSentReceived.AmountSent.NumDecimal(),
		AmountRequested: amountSentReceived.AmountRequested.NumDecimal(),
		AmountPaid:      amountSentReceived.AmountPaid.NumDecimal(),
		AmountCredited:  amountDebitedCredited.AmountCredited.NumDecimal(),
		AmountDebited:   amountDebitedCredited.AmountDebited.NumDecimal(),
		Currency:        transaction.Currency(banking.CurrencyUSD),
		RecordedDate:    shared.Date(dayStartLocal),
	}
	_, err = transaction.CreateDailyTransactionStats(&c)
	if err != nil {
		log.Println("CreateDailyTransactionStats:", err, bID)
	} else {
		log.Println("Stats Created:", bID)
	}

	return err
}
