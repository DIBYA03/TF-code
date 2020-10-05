/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/
package transaction

import (
	"errors"
	"time"

	"github.com/wiseco/core-platform/shared"
	"github.com/wiseco/go-lib/num"
)

type MonthlyInterestID string

type BusinessAccountMonthlyInterest struct {
	ID               MonthlyInterestID `json:"id" db:"id"`
	AccountID        AccountID         `json:"accountId" db:"account_id"`
	BusinessID       shared.BusinessID `json:"businessId" db:"business_id"`
	AvgPostedBalance num.Decimal       `json:"avgBalance" db:"avg_posted_balance"`
	Days             int               `json:"days" db:"days"`
	InterestAmount   num.Decimal       `json:"interestAmount" db:"interest_amount"`
	InterestPayout   num.Decimal       `json:"interestPayout" db:"interest_payout"`
	Currency         Currency          `json:"currency" db:"currency"`
	APR              int               `json:"apr" db:"apr"`
	StartDate        shared.Date       `json:"startDate" db:"start_date"`
	EndDate          shared.Date       `json:"endDate" db:"end_date"`
	RecordedDate     shared.Date       `json:"recordedDate" db:"recorded_date"`
	Created          time.Time         `json:"created" db:"created"`
}

type BusinessAccountMonthlyInterestCreate struct {
	AccountID        AccountID         `json:"accountId" db:"account_id"`
	BusinessID       shared.BusinessID `json:"businessId" db:"business_id"`
	AvgPostedBalance num.Decimal       `json:"avgBalance" db:"avg_posted_balance"`
	Days             int               `json:"days" db:"days"`
	InterestAmount   num.Decimal       `json:"interestAmount" db:"interest_amount"`
	InterestPayout   num.Decimal       `json:"interestPayout" db:"interest_payout"`
	Currency         Currency          `json:"currency" db:"currency"`
	APR              int               `json:"apr" db:"apr"`
	StartDate        shared.Date       `json:"startDate" db:"start_date"`
	EndDate          shared.Date       `json:"endDate" db:"end_date"`
	RecordedDate     shared.Date       `json:"recordedDate" db:"recorded_date"`
}

func CreateMonthlyInterest(c *BusinessAccountMonthlyInterestCreate) (*BusinessAccountMonthlyInterest, error) {

	q := `
	    INSERT INTO business_account_monthly_interest(
			account_id, business_id, avg_posted_balance, days, interest_amount, interest_payout, currency,
			apr, start_date, end_date, recorded_date
		)
		VALUES(
			:account_id, :business_id, :avg_posted_balance, :days, :interest_amount, :interest_payout,
			:currency, :apr, :start_date, :end_date, :recorded_date
		)
		RETURNING *`

	rows, err := DBWrite.NamedQuery(q, c)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var b BusinessAccountMonthlyInterest
		err = rows.StructScan(&b)
		rows.Close()

		if err == nil {
			return &b, nil
		}

		return nil, err
	}

	return nil, errors.New("monthly interest internal error")
}
