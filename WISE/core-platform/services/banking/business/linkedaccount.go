/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

// Package for all internal banking related items such as accounts and cards
package business

import (
	"encoding/json"
	"time"

	"github.com/wiseco/core-platform/services"
	"github.com/wiseco/core-platform/services/banking"
	"github.com/wiseco/core-platform/shared"
)

type AccountNumber string

func (n *AccountNumber) String() string {
	return string(*n)
}

// Marshal and transform fields as needed
func (n *AccountNumber) MarshalJSON() ([]byte, error) {
	return json.Marshal(services.MaskLeft(n.String(), 4))
}

// Linked accounts can be used for paying contractors, vendors, etc. via ACH push
// ACH pull requires full verification
type LinkedBankAccount struct {
	// Bank account id
	Id string `json:"id" db:"id"`

	// Account holder user id - Cognito ID
	UserID shared.UserID `json:"userId" db:"user_id"`

	// Business Id
	BusinessID shared.BusinessID `json:"businessId" db:"business_id"`

	// Wise Business Bank account id
	BusinessBankAccountId *string `json:"businessBankAccountId" db:"business_bank_account_id"`

	// Contact id
	ContactId *string `json:"contactId" db:"contact_id"`

	// Account Reference Id - This Id is from registered bank(bbva for example)
	RegisteredAccountId string `json:"registeredAccountId" db:"registered_account_id"`

	// Registered bank name - Ex: bbva
	RegisteredBankName string `json:"registeredBankName"  db:"registered_bank_name"`

	// Ex: Plaid Account Id
	SourceAccountId *string `json:"sourceAccountId" db:"source_account_id"`

	// Account holder
	AccountHolderName string `json:"accountHolderName" db:"account_holder_name"`

	// Account name - Ex: Plaid Checking
	AccountName *string `json:"accountName" db:"account_name"`

	// Account number masked
	AccountNumber AccountNumber `json:"accountNumberMasked" db:"account_number"`

	// Bank name - Ex: Wells Fargo
	BankName *string `json:"bankName" db:"bank_name"`

	// Denominated currency
	Currency banking.Currency `json:"currency" db:"currency"`

	// Type - Ex: checking or savings
	AccountType banking.AccountType `json:"accountType" db:"account_type"`

	// Account usage type
	UsageType *UsageType `json:"usageType" db:"usage_type"`

	// Primary account routing number
	RoutingNumber string `json:"routingNumber" db:"routing_number"`

	// Wire Routing Number
	WireRouting *string `json:"wireRouting" db:"wire_routing"`

	// Plaid Request Id
	SourceId *string `json:"sourceId" db:"source_id"`

	// Plaid Request Name
	SourceName *string `json:"sourceName" db:"source_name"`

	// Account Usage - Ex: send_and_receive, send_only or receive_only
	Permission banking.LinkedAccountPermission `json:"permission" db:"account_permission"`

	// Account alias
	Alias *string `json:"alias" db:"alias"`

	// External accounts need to be verified
	Verified bool `json:"verified" db:"verified"`

	// Deactivated
	Deactivated *time.Time `json:"deactivated" db:"deactivated"`

	// Created timestamp
	Created time.Time `json:"created" db:"created"`

	// Modified timestamp
	Modified time.Time `json:"modified" db:"modified"`
}

/*LinkedBankAccountBase has minimal bank account account thats sent back to client and user can choose which account
needs to be registered */
type LinkedBankAccountBase struct {
	// Plaid Account Id
	SourceAccountId *string `json:"sourceAccountId" db:"source_account_id"`

	// Account name - Ex: Plaid Checking
	AccountName string `json:"accountName" db:"account_name"`

	// Account number masked
	AccountNumber AccountNumber `json:"accountNumberMasked" db:"account_number"`
}

type LinkedExternalAccountCreate struct {
	// Bank account id - Plaid Account Id
	SourceAccountId string `json:"sourceAccountId" db:"source_account_id"`

	// Public token
	PublicToken string `json:"publicToken"`

	// Account holder user id - Cognito ID
	UserID shared.UserID `json:"userId" db:"user_id"`

	// Business Id
	BusinessID shared.BusinessID `json:"businessId" db:"business_id"`

	// Contact ID
	ContactID *string `json:"contactId" db:"contact_id"`
}

type LinkedBankAccountRequest struct {
	// Public token
	PublicToken string `json:"publicToken"`
}

const (
	// Depository account
	PlaidAccountTypeDepository = "depository"
)

type ClearingLinkedAccountCreate struct {
	// User Id
	UserID shared.UserID `json:"userId" db:"user_id"`

	// Business Id
	BusinessID shared.BusinessID `json:"businessId" db:"business_id"`

	// Wise Business Bank account id
	BusinessBankAccountId *string `json:"businessBankAccountId" db:"business_bank_account_id"`

	// Account Reference Id - This Id is from registered bank(bbva for example)
	RegisteredAccountId string `json:"registeredAccountId" db:"registered_account_id"`

	// Registered bank name - Ex: bbva
	RegisteredBankName string `json:"registeredBankName"  db:"registered_bank_name"`

	// Account number masked
	AccountNumber AccountNumber `json:"accountNumber" db:"account_number"`

	// Type - Ex: checking or savings
	AccountType banking.AccountType `json:"accountType" db:"account_type"`

	// Account usage type
	UsageType *UsageType `json:"usageType" db:"usage_type"`

	// Primary account routing number
	RoutingNumber string `json:"routingNumber" db:"routing_number"`

	// Denominated currency
	Currency banking.Currency `json:"currency" db:"currency"`

	// Account Usage - Ex: send_and_receive, send_only or receive_only
	Permission banking.LinkedAccountPermission `json:"permission" db:"account_permission"`
}

type ContactLinkedAccountCreate struct {
	// User Id
	UserID shared.UserID `json:"userId" db:"user_id"`

	// Business Id
	BusinessID shared.BusinessID `json:"businessId" db:"business_id"`

	// Name on the account
	ContactId string `json:"contactId" db:"contact_id"`

	// Account Reference Id - This Id is from registered bank(bbva for example)
	RegisteredAccountId string `json:"registeredAccountId" db:"registered_account_id"`

	// Registered bank name - Ex: bbva
	RegisteredBankName string `json:"registeredBankName"  db:"registered_bank_name"`

	// Account number masked
	AccountNumber AccountNumber `json:"accountNumber" db:"account_number"`

	// Type - Ex: checking or savings
	AccountType banking.AccountType `json:"accountType" db:"account_type"`

	// Account usage type
	UsageType *UsageType `json:"usageType" db:"usage_type"`

	// Primary account routing number
	RoutingNumber string `json:"routingNumber" db:"routing_number"`

	// Denominated currency
	Currency banking.Currency `json:"currency" db:"currency"`

	// Account Usage - Ex: send_and_receive, send_only or receive_only
	Permission banking.LinkedAccountPermission `json:"permission" db:"account_permission"`
}

type MerchantLinkedAccountCreate struct {
	// User Id
	UserID shared.UserID `json:"userId" db:"user_id"`

	// Business Id
	BusinessID shared.BusinessID `json:"businessId" db:"business_id"`

	// Account Reference Id - This Id is from registered bank(bbva for example)
	RegisteredAccountId string `json:"registeredAccountId" db:"registered_account_id"`

	// Registered bank name - Ex: bbva
	RegisteredBankName string `json:"registeredBankName"  db:"registered_bank_name"`

	// Account holder
	AccountHolderName string `json:"accountHolderName" db:"account_holder_name"`

	// Account number masked
	AccountNumber AccountNumber `json:"accountNumber" db:"account_number"`

	// Type - Ex: checking or savings
	AccountType banking.AccountType `json:"accountType" db:"account_type"`

	// Account usage type
	UsageType UsageType `json:"usageType" db:"usage_type"`

	// Primary account routing number
	RoutingNumber string `json:"routingNumber" db:"routing_number"`

	// Denominated currency
	Currency banking.Currency `json:"currency" db:"currency"`

	// Account Usage - Ex: send_and_receive, send_only or receive_only
	Permission banking.LinkedAccountPermission `json:"permission" db:"account_permission"`
}

type LinkedAccountUpdate struct {
	// Bank account id
	ID string `json:"id" db:"id"`

	// Account usage type
	UsageType *UsageType `json:"usageType" db:"usage_type"`
}

func (l *LinkedBankAccount) Source() LinkedAccountSource {
	if l.SourceName == nil {
		return ""
	}

	return LinkedAccountSource(*l.SourceName)
}

func (a LinkedBankAccount) IntraBankAccount() bool {
	if a.RoutingNumber == "062001186" {
		return true
	}

	return false
}
