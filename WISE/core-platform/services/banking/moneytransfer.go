/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

// Package for all internal banking related items such as accounts and cards
package banking

import (
	"time"

	"github.com/jmoiron/sqlx/types"
	"github.com/wiseco/core-platform/shared"
)

const (
	// Business transfer
	MoneyTransferTypeBusiness = "business"

	// Business payroll transfer
	MoneyTransferTypePayroll = "payroll"

	// Tax transfer
	MoneyTransferTypeTax = "tax"
)

const (
	// Wise intrabank transfer (wise payment rail)
	MoneyTransferRailTypeWise = "wise"

	// Push to debit
	MoneyTransferRailTypeDebit = "debit"

	// ACH transfer
	MoneyTransferRailTypeACH = "ach"

	// Wire transfer
	MoneyTransferRailTypeWire = "wire"

	// Remittance transfer
	MoneyTransferRailTypeRemit = "remit"
)

const (
	// Transfer account type checking
	MoneyTransferAccountTypeChecking = "checking"

	// Transfer account type savings
	MoneyTransferAccountTypeSavings = "savings"

	// Transfer account type debit
	MoneyTransferAccountTypeDebit = "debit"
)

const (
	// Transfer is pending - can be cancelled
	MoneyTransferStatusPending = "pending"

	// Transfer submitted and is in process
	MoneyTransferStatusInProcess = "inProcess"

	// Transfer has been canceled
	MoneyTransferStatusCanceled = "canceled"

	// Transfer Posted (intrabank)
	MoneyTransferStatusPosted = "posted"

	// Transfer settled (extenbal or push to debit)
	MoneyTransferStatusSettled = "settled"

	// Debit sent to origin bank
	MoneyTransferStatusDebitSent = "debitSent"

	// Credit sent to destination bank
	MoneyTransferStatusCreditSent = "creditSent"

	// Review issue resolved
	MoneyTransferStatusReviewResolved = "reviewResolved"

	// Pull failed
	MoneyTransferStatusPullFailed = "pullFailed"

	// Pull refunded
	MoneyTransferStatusPullRefunded = "pullRefunded"

	// Pull transfer under review
	MoneyTransferStatusPullReview = "pullReview"

	// Push refunded
	MoneyTransferStatusPushRefunded = "pushRefunded"

	// Push under review
	MoneyTransferStatusPushReview = "pushReview"

	// Bank error
	MoneyTransferStatusBankError = "bankError"
)

type TransferType string

const (
	TransferTypeAccount = TransferType("account")
	TransferTypeCard    = TransferType("card")
	TransferTypeCheck   = TransferType("check")
)

// Wise money transfer
type MoneyTransfer struct {
	// Transaction id
	Id string `json:"id" db:"id"`

	// Created user id
	CreatedUserID shared.UserID `json:"createdUserId" db:"created_user_id"`

	// Contact id
	ContactId *string `json:"contactId" db:"contact_id"`

	//  bank name - e.g. 'bbva'
	BankName BankName `json:"bankName" db:"bank_name"`

	//  bank account id
	BankTransferId string `json:"bankTransferId" db:"bank_transfer_id"`

	// Partner bank account extra data
	BankExtra types.JSONText `json:"-" db:"bank_extra"`

	// Source bank account id
	SourceAccountId string `json:"sourceAccountId" db:"source_account_id"`

	// Ex - card or account
	SourceType TransferType `json:"sourceType" db:"source_type"`

	// Destination bank account id
	DestAccountId string `json:"destAccountId" db:"dest_account_id"`

	// Ex - card or account
	DestType TransferType `json:"destType" db:"dest_type"`

	// Transfer amount
	Amount float64 `json:"amount" db:"amount"`

	// Transfer currency
	Currency Currency `json:"currency" db:"currency"`

	// Transfer Notes
	Notes *string `json:"notes" db:"notes"`

	// Transfer status
	Status string `json:"status" db:"status"`

	// Send email
	SendEmail bool `json:"sendEmail" db:"send_email"`

	/* Related transaction
	Transactions *[]Transaction `json:"transactions"` */

	// Created date
	Created time.Time `json:"created" db:"created"`

	// Modified date
	Modified time.Time `json:"modified" db:"modified"`
}

type ConsumerMoneyTransfer struct {
	MoneyTransfer
}

// Contact mini object for notifications
type NotificationMoneyTransfer struct {
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
