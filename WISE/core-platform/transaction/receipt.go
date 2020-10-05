/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

// Package for transaction services
package transaction

import (
	"time"

	"github.com/wiseco/core-platform/shared"
)

type ReceiptCreate struct {
	TransactionID shared.PostedTransactionID `json:"transactionId" db:"transaction_id"`

	CreatedUserID shared.UserID `json:"createdUserId" db:"created_user_id"`

	BusinessID shared.BusinessID `json:"businessId" db:"business_id"`

	ContentType *string `json:"contentType" db:"content_type"`

	Content *string `json:"content"`

	StorageKey *string `json:"storageKey,omitempty" db:"storage_key"`
}

type Receipt struct {
	Id string `json:"id" db:"id"`

	TransactionID shared.PostedTransactionID `json:"transactionId" db:"transaction_id"`

	CreatedUserID shared.UserID `json:"createdUserId" db:"created_user_id"`

	BusinessID shared.BusinessID `json:"businessId" db:"business_id"`

	ContentType *string `json:"contentType" db:"content_type"`

	StorageKey *string `json:"-" db:"storage_key"`

	SignedURL *string `json:"uploadURL"`

	Deleted *time.Time `json:"deleted" db:"deleted"`

	Created time.Time `json:"created" db:"created"`

	Modified time.Time `json:"modified" db:"modified"`
}
