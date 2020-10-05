package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/wiseco/core-platform/services"
	"github.com/wiseco/core-platform/services/banking"
	"github.com/wiseco/core-platform/services/banking/business"
	bsrv "github.com/wiseco/core-platform/services/business"
	"github.com/wiseco/core-platform/shared"
	"github.com/wiseco/core-platform/transaction"
	txndata "github.com/wiseco/core-platform/transaction"
	"github.com/wiseco/go-lib/id"
)

var dayStart, dayEnd, dayStartLocal, dayEndLocal, recordedMonth time.Time
var clearingUserID, clearingBusinessID, clearingAccountID, clearingLinkedAccountID string

func main() {
	// Send money transfer from wise clearing account
	clearingUserID, err := shared.ParseUserID(os.Getenv("WISE_CLEARING_USER_ID"))
	if err != nil {
		panic(err)
	}

	clearingBusinessID, err := shared.ParseBusinessID(os.Getenv("WISE_CLEARING_BUSINESS_ID"))
	if err != nil {
		panic(err)
	}

	clearingAccountID = os.Getenv("WISE_CLEARING_ACCOUNT_ID")
	if clearingAccountID == "" {
		panic(errors.New("clearing account missing"))
	}

	clearingLinkedAccountID = os.Getenv("WISE_CLEARING_LINKED_ACCOUNT_ID")
	if clearingLinkedAccountID == "" {
		panic(errors.New("clearing linked account missing"))
	}

	if os.Getenv("USE_BANKING_SERVICE") == "true" {
		clearingLinkedAccountID = id.IDPrefixLinkedBankAccount.String() + clearingLinkedAccountID
	}

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

	dayEndLocal = time.Date(nowLocal.Year(), nowLocal.Month(), 1, 0, 0, 0, 0, loc)
	dayStartLocal = dayEndLocal.AddDate(0, -1, 0)
	recordedMonth = dayStartLocal

	dayStart = dayStartLocal.In(utcLoc)
	dayEnd = dayEndLocal.In(utcLoc)
	if dayEnd.After(nowUTC) {
		//panic(fmt.Errorf("Error: day end (%v) is after current time (%v)", dayEnd, nowUTC))
	}

	// Calculate Days in Year
	yearStartLocal := time.Date(dayStartLocal.Year(), 1, 1, 0, 0, 0, 0, loc)
	yearEndLocal := yearStartLocal.AddDate(1, 0, 0)
	daysInCurrentYear := int64(yearEndLocal.Sub(yearStartLocal).Hours() / 24)

	// 1% APY
	apr := 99
	var summaries []transaction.BusinessAccountMonthlyInterestCreate
	err = txndata.DBRead.Select(
		&summaries, `
		SELECT
			account_id,
			business_id,
			ROUND(AVG(posted_balance), 2) avg_posted_balance,
			COUNT(*) days,
			currency,
			ROUND(AVG(posted_balance) * (CAST($1 AS DECIMAL(19,4)) / CAST(10000 AS DECIMAL(19,4))) * COUNT(*) / CAST($2 AS DECIMAL(19,4)), 4) interest_amount,
			ROUND(FLOOR(AVG(posted_balance) * (CAST($1 AS DECIMAL(19,4)) / CAST(10000 AS DECIMAL(19,4))) * COUNT(*) / CAST(365 AS DECIMAL(19,4)) * 100) / 100, 2) interest_payout,
			$1 apr,
			MIN(recorded_date) start_date,
			MAX(recorded_date) end_date,
			to_date($4, 'YYYY-MM-DD') recorded_date
			FROM business_account_daily_balance
			WHERE
	            recorded_date >= to_date($3, 'YYYY-MM-DD') AND
    	        recorded_date < to_date($4, 'YYYY-MM-DD') AND
				account_id != $5
			GROUP BY account_id, business_id, currency
			ORDER BY avg_posted_balance DESC`,
		apr,
		daysInCurrentYear,
		shared.Date(dayStartLocal),
		shared.Date(dayEndLocal),
		clearingAccountID,
	)
	if err != nil {
		panic(err)
	}

	for _, s := range summaries {
		if string(s.AccountID) == clearingAccountID {
			log.Printf("%s: ignore clearning account", s.AccountID)
			continue
		}

		mi, err := transaction.CreateMonthlyInterest(&s)
		if err != nil {
			log.Printf("%s: %v", s.AccountID, err)
			continue
		}

		interest, ok := mi.InterestPayout.Float64()
		if !ok {
			log.Printf("%s: invalid interest value", mi.AccountID)
			continue
		} else if interest <= 0 {
			log.Printf(
				"%s: zero or negative interest value of %s %s ",
				mi.AccountID,
				strconv.FormatFloat(interest, 'f', 2, 64),
				strings.ToUpper(string(mi.Currency)),
			)
			continue
		}

		// Get Business
		srcReq := services.NewSourceRequest()
		b, err := bsrv.NewBusinessServiceWithout().GetByIdInternal(mi.BusinessID)
		if err != nil {
			log.Printf("%s: %v", s.BusinessID, err)
			continue
		}

		// Get linked account
		srcReq = services.NewSourceRequest()
		srcReq.UserID = clearingUserID
		acc, err := business.NewBankAccountService(srcReq).GetByIDInternal(string(mi.AccountID))
		if err != nil {
			log.Printf("%s: %v", s.AccountID, err)
			continue
		}

		destAccountID := acc.Id
		if os.Getenv("USE_BANKING_SERVICE") != "true" {
			srcReq = services.NewSourceRequest()
			srcReq.UserID = clearingUserID
			linkedAccount, err := business.NewLinkedAccountService(srcReq).GetByAccountNumber(
				clearingBusinessID,
				business.AccountNumber(acc.AccountNumber),
				acc.RoutingNumber,
			)
			if err != nil {
				srcReq = services.NewSourceRequest()
				srcReq.UserID = clearingUserID
				linkedAccount, err = business.NewLinkedAccountService(srcReq).LinkMerchantBankAccount(
					&business.MerchantLinkedAccountCreate{
						UserID:            clearingUserID,
						BusinessID:        clearingBusinessID,
						AccountHolderName: b.Name(),
						AccountNumber:     business.AccountNumber(acc.AccountNumber),
						AccountType:       banking.AccountType(acc.AccountType),
						RoutingNumber:     acc.RoutingNumber,
						Currency:          acc.Currency,
						Permission:        banking.LinkedAccountPermissionSendAndRecieve,
					},
				)
				if err != nil {
					log.Printf("%s: %v", s.AccountID, err)
					continue
				}
			}

			destAccountID = linkedAccount.Id
		}
		notes := fmt.Sprintf("Interest earned in %s", recordedMonth.Format("Jan 2006"))
		log.Printf(notes)

		// Move money
		mid := string(mi.ID)
		ti := business.TransferInitiate{
			CreatedUserID:     clearingUserID,
			BusinessID:        clearingBusinessID,
			SourceAccountId:   clearingLinkedAccountID,
			DestAccountId:     destAccountID,
			Amount:            interest,
			SourceType:        banking.TransferTypeAccount,
			DestType:          banking.TransferTypeAccount,
			Currency:          banking.Currency(mi.Currency),
			MonthlyInterestID: &mid,
			Notes:             &notes,
		}

		srcReq = services.NewSourceRequest()
		srcReq.UserID = clearingUserID
		mt, err := business.NewMoneyTransferService(srcReq).Transfer(&ti)
		if err != nil {
			log.Printf("%s: %v", mi.AccountID, err)
			continue
		}

		log.Printf(
			"Interest transfer (id: %s) sent to account (id: %s): %s %s",
			mt.Id,
			mi.AccountID,
			strconv.FormatFloat(interest, 'f', 2, 64),
			strings.ToUpper(string(mi.Currency)),
		)
	}
}
