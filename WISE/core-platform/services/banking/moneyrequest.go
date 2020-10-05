/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

// Package for all internal banking related items such as accounts and cards
package banking

import (
	"time"

	"github.com/wiseco/core-platform/shared"
)

// Wise money transfer
type MoneyRequest struct {
	// Transaction id
	ID shared.PaymentRequestID `json:"id" db:"id"`

	// Created user id
	CreatedUserID shared.UserID `json:"createdUserId" db:"created_user_id"`

	// Contact id
	ContactId string `json:"contactId" db:"contact_id"`

	// Transfer amount
	Amount float64 `json:"amount" db:"amount"`

	// Transfer currency
	Currency Currency `json:"currency" db:"currency"`

	// Transfer Notes
	Notes string `json:"notes" db:"notes"`

	// Transfer status
	Status *MoneyRequestStatus `json:"status" db:"request_status"`

	// Message id
	MessageId string `json:"messageId" db:"message_id"`

	// Payment Intent Secret
	PaymentIntentToken string `json:"paymentIntentToken"`

	// Request type
	RequestType *PaymentRequestType `json:"requestType" db:"request_type"`

	// POS ID
	CardReaderID *shared.CardReaderID `json:"cardReaderId" db:"pos_id"`

	// Created date
	Created time.Time `json:"created" db:"created"`

	// Modified date
	Modified time.Time `json:"modified" db:"modified"`
}

// Wise money transfer update
type MoneyRequestUpdate struct {
	// Transaction id
	ID shared.PaymentRequestID `json:"id" db:"id"`

	// Transfer status
	Status MoneyRequestStatus `json:"status" db:"request_status"`
}

// Wise money request message ID update
type MoneyRequestIDUpdate struct {
	// Transaction id
	ID shared.PaymentRequestID `json:"id" db:"id"`

	// Contact id
	ContactId string `json:"contactId" db:"contact_id"`

	// Message ID
	MessageId string `json:"messageId" db:"message_id"`
}

type ConsumerMoneyRequest struct {
	MoneyRequest
}

type MoneyRequestStatus string

const (
	MoneyRequestStatusPending  = MoneyRequestStatus("pending")
	MoneyRequestStatusFailed   = MoneyRequestStatus("failed")
	MoneyRequestStatusComplete = MoneyRequestStatus("complete")
)

type PaymentRequestType string

const (
	PaymentRequestTypePOS     = PaymentRequestType("pos")
	PaymentRequestTypeInvoice = PaymentRequestType("invoice")
)
