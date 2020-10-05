/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

// Package for all internal banking related items such as accounts and cards
package transaction

import (
	"time"

	"github.com/wiseco/core-platform/services/banking"
)

// Mini object for transactions
type MoneyTransfer struct {
	// Used internally
	SourceAccountID string               `json:"-" db:"source_account_id"`
	SourceType      banking.TransferType `json:"-" db:"source_type"`
	DestAccountID   string               `json:"-" db:"dest_account_id"`
	DestType        banking.TransferType `json:"-" db:"dest_type"`

	// Transfer amount
	Amount *float64 `json:"amount,omitempty" db:"amount"`

	// Transfer currency
	Currency *Currency `json:"currency,omitempty" db:"currency"`

	// Transfer Notes
	Notes *string `json:"notes,omitempty" db:"notes"`

	// Transfer status
	Status *string `json:"status,omitempty" db:"status"`

	// Created timestamp
	Created *time.Time `json:"created,omitempty" db:"created"`
}
