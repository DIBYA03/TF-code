/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

// Package for transaction services
package transaction

import (
	"time"

	"github.com/wiseco/core-platform/shared"
)

//Category is the dispute category - for ex: still being charged, incorrect charge etc..
type Category string

const (
	CategoryStillBeingCharged = Category("stillBeingCharged")

	CategoryIncorrectCharge = Category("incorrectCharge")

	CategoryFraudulentCharge = Category("fraudulentCharge")

	CategoryNotSatisfied = Category("notSatisfied")
)

var categoryTo = map[Category]string{
	CategoryStillBeingCharged: "still being charged",
	CategoryIncorrectCharge:   "incorrectly charged",
	CategoryFraudulentCharge:  "fradulent charge",
	CategoryNotSatisfied:      "product or services were not as described",
}

//DisputeStatus is the dispute status - for ex: disputed, dispute credited, dispute not credited, dispute cancelled
type DisputeStatus string

const (
	DisputeStatusDisputed = DisputeStatus("disputed")

	DisputeStatusDisputedNotCredited = DisputeStatus("disputedNotCredited")

	DisputeStatusDisputedCredited = DisputeStatus("disputedCredited")

	DisputeStatusDisputedCancelled = DisputeStatus("disputedCancelled")
)

type DisputeCreate struct {
	DisputeNumber string `json:"disputeNumber" db:"dispute_number"`

	TransactionID shared.PostedTransactionID `json:"transactionId" db:"transaction_id"`

	ReceiptID *string `json:"receiptId" db:"receipt_id"`

	CreatedUserID shared.UserID `json:"createdUserId" db:"created_user_id"`

	BusinessID shared.BusinessID `json:"businessId" db:"business_id"`

	Category Category `json:"category" db:"category"`

	Summary *string `json:"summary" db:"summary"`

	DisputeStatus *DisputeStatus `json:"disputeStatus" db:"dispute_status"`
}

type Dispute struct {
	Id string `json:"id" db:"id"`

	DisputeNumber string `json:"disputeNumber" db:"dispute_number"`

	TransactionID shared.PostedTransactionID `json:"transactionId" db:"transaction_id"`

	ReceiptID *string `json:"receiptId" db:"receipt_id"`

	CreatedUserID shared.UserID `json:"createdUserId" db:"created_user_id"`

	BusinessID shared.BusinessID `json:"businessId" db:"business_id"`

	Category Category `json:"category" db:"category"`

	Summary *string `json:"summary" db:"summary"`

	DisputeStatus *DisputeStatus `json:"disputeStatus" db:"dispute_status"`

	Created *time.Time `json:"created" db:"created"`

	Modified *time.Time `json:"modified" db:"modified"`
}

type DisputeCancel struct {
	Id string `json:"id" db:"id"`

	TransactionID shared.PostedTransactionID `json:"transactionId" db:"transaction_id"`

	CreatedUserID shared.UserID `json:"createdUserId" db:"created_user_id"`

	BusinessID shared.BusinessID `json:"businessId" db:"business_id"`
}

type DisputeMini struct {
	Id *string `json:"id,omitempty" db:"id"`

	ReceiptID *string `json:"receiptId,omitempty" db:"receipt_id"`

	Category *Category `json:"category,omitempty" db:"category"`

	Summary *string `json:"summary,omitempty" db:"summary"`

	DisputeStatus *DisputeStatus `json:"disputeStatus,omitempty" db:"dispute_status"`

	DisputeNumber *string `json:"disputeNumber,omitempty" db:"dispute_number"`

	Created *time.Time `json:"created,omitempty" db:"created"`
}
