/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

// Package for all internal banking related items such as accounts and cards
package business

import (
	"time"

	"github.com/wiseco/core-platform/partner/bank"
	"github.com/wiseco/core-platform/shared"
)

type PayeeStatus string

const (
	PayeeStatusInactive = PayeeStatus("inactive")
	PayeeStatusActive   = PayeeStatus("active")
)

// Linked Payees can be used to send checks
type LinkedPayee struct {
	// Bank account id
	ID string `json:"id" db:"id"`

	// Business Id
	BusinessID shared.BusinessID `json:"businessID" db:"business_id"`

	// Contact id
	ContactID shared.ContactID `json:"contactID" db:"contact_id"`

	// Address id
	AddressID shared.AddressID `json:"addressID" db:"address_id"`

	// Bank Payee ID - id sent to use from the bank
	BankPayeeID string `json:"BankPayeeID" db:"bank_payee_id"`

	// Bank Name associated with payee
	BankName bank.ProviderName `json:"BankName" db:"bank_name"`

	// Account holder name
	AccountHolderName string `json:"accountHolderName" db:"account_holder_name"`

	// Payee name aka business name
	PayeeName string `json:"payeeName" db:"payee_name"`

	//State
	Status PayeeStatus `json:"status" db:"status"`

	// Created timestamp
	Created time.Time `json:"created" db:"created"`

	// Modified timestamp
	Modified time.Time `json:"modified" db:"modified"`
}

type LinkedPayeeCreate struct {
	// Business Id
	BusinessID shared.BusinessID `json:"businessID" db:"business_id"`

	// Contact id
	ContactID string `json:"contactID" db:"contact_id"`

	// Address id
	AddressID shared.AddressID `json:"addressID" db:"address_id"`

	// Bank Payee ID - id sent to use from the bank
	BankPayeeID string `json:"bankPayeeID" db:"bank_payee_id"`

	// Bank Name associated with payee
	BankName bank.ProviderName `json:"BankName" db:"bank_name"`

	// Account holder name
	AccountHolderName string `json:"accountHolderName" db:"account_holder_name"`

	// Payee name aka business name
	PayeeName string `json:"payeeName" db:"payee_name"`

	//State
	Status PayeeStatus `json:"status" db:"status"`
}
