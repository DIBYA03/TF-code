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

type DailyTransactionID string

type BusinessDailyTransaction struct {
	ID                 DailyTransactionID `json:"id" db:"id"`
	BusinessID         shared.BusinessID  `json:"businessId" db:"business_id"`
	AmountRequested    num.Decimal        `json:"amountRequested" db:"amount_requested"`
	AmountPaid         num.Decimal        `json:"amountPaid" db:"amount_paid"`
	AmountSent         num.Decimal        `json:"amountSent" db:"amount_sent"`
	AmountCredited     num.Decimal        `json:"amountCredited" db:"amount_credited"`
	AmountDebited      num.Decimal        `json:"amountDebited" db:"amount_debited"`
	AmountRequestedDep float64            `json:"amountRequestedDep" db:"amount_requested_dep"` // Deprecated
	AmountPaidDep      float64            `json:"amountPaidDep" db:"amount_paid_dep"`           // Deprecated
	AmountSentDep      float64            `json:"amountSentDep" db:"amount_sent_dep"`           // Deprecated
	AmountCreditedDep  float64            `json:"amountCreditedDep" db:"amount_credited_dep"`   // Deprecated
	AmountDebitedDep   float64            `json:"amountDebitedDep" db:"amount_debited_dep"`     // Deprecated
	Currency           Currency           `json:"currency" db:"currency"`
	RecordedDate       shared.Date        `json:"recordedDate" db:"recorded_date"`
	Created            time.Time          `json:"created" db:"created"`
}

type BusinessDailyTransactionCreate struct {
	BusinessID      shared.BusinessID `json:"businessId" db:"business_id"`
	AmountRequested num.Decimal       `json:"amountRequested" db:"amount_requested"`
	AmountPaid      num.Decimal       `json:"amountPaid" db:"amount_paid"`
	AmountSent      num.Decimal       `json:"mneySent" db:"amount_sent"`
	AmountCredited  num.Decimal       `json:"amountCredited" db:"amount_credited"`
	AmountDebited   num.Decimal       `json:"amountDebited" db:"amount_debited"`
	Currency        Currency          `json:"currency" db:"currency"`
	RecordedDate    shared.Date       `json:"recordedDate" db:"recorded_date"`
}

func CreateDailyTransactionStats(c *BusinessDailyTransactionCreate) (*BusinessDailyTransaction, error) {
	q := `
	    INSERT INTO business_daily_transaction_stats(
			business_id, amount_requested, amount_paid, amount_sent, amount_credited, amount_debited,
			currency, recorded_date
		)
		VALUES(
			:business_id, :amount_requested, :amount_paid, :amount_sent, :amount_credited, :amount_debited,
			:currency, :recorded_date
		)
		RETURNING *`

	rows, err := DBWrite.NamedQuery(q, c)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var b BusinessDailyTransaction
		err = rows.StructScan(&b)
		rows.Close()

		if err == nil {
			return &b, nil
		}

		return nil, err
	}

	return nil, errors.New("create daily transaction internal error")
}
