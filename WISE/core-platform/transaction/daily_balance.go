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

type DailyBalanceID string
type AccountID string

type BusinessAccountDailyBalance struct {
	ID                DailyBalanceID    `json:"id" db:"id"`
	AccountID         AccountID         `json:"accountId" db:"account_id"`
	BusinessID        shared.BusinessID `json:"businessId" db:"business_id"`
	PostedBalance     num.Decimal       `json:"postedBalance" db:"posted_balance"`
	AmountCredited    num.Decimal       `json:"amountCredited" db:"amount_credited"`
	AmountDebited     num.Decimal       `json:"amountDebited" db:"amount_debited"`
	PostedBalanceDep  float64           `json:"postedBalanceDep" db:"posted_balance_dep"`   // Deprecated
	AmountCreditedDep float64           `json:"amountCreditedDep" db:"amount_credited_dep"` // Deprecated
	AmountDebitedDep  float64           `json:"amountDebitedDep" db:"amount_debited_dep"`   // Deprecated
	Currency          Currency          `json:"currency" db:"currency"`
	APR               int               `json:"apr" db:"apr"`
	RecordedDate      shared.Date       `json:"recordedDate" db:"recorded_date"`
	Created           time.Time         `json:"created" db:"created"`
}

type BusinessAccountDailyBalanceCreate struct {
	AccountID      AccountID         `json:"accountId" db:"account_id"`
	BusinessID     shared.BusinessID `json:"businessId" db:"business_id"`
	PostedBalance  num.Decimal       `json:"postedBalance" db:"posted_balance"`
	AmountCredited num.Decimal       `json:"amountCredited" db:"amount_credited"`
	AmountDebited  num.Decimal       `json:"amountDebited" db:"amount_debited"`
	Currency       Currency          `json:"currency" db:"currency"`
	APR            int               `json:"apr" db:"apr"`
	RecordedDate   shared.Date       `json:"recordedDate" db:"recorded_date"`
}

func CreateDailyBalance(c *BusinessAccountDailyBalanceCreate) (*BusinessAccountDailyBalance, error) {

	q := `
	    INSERT INTO business_account_daily_balance(
			account_id, business_id, posted_balance, amount_credited, amount_debited, currency, apr,
			recorded_date
		)
		VALUES(
			:account_id, :business_id, :posted_balance, :amount_credited, :amount_debited, :currency, :apr,
			:recorded_date
		)
		RETURNING *`

	rows, err := DBWrite.NamedQuery(q, c)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var b BusinessAccountDailyBalance
		err = rows.StructScan(&b)
		rows.Close()

		if err == nil {
			return &b, nil
		}

		return nil, err
	}

	return nil, errors.New("create daily balance internal error")
}
