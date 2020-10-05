/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

// Package for all internal banking related items such as accounts and cards
package banking

import (
	"encoding/json"
	"time"

	"github.com/jmoiron/sqlx/types"
	"github.com/wiseco/core-platform/services"
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

//BankAccount  Wise bank account object
type BankAccount struct {
	// Bank account id
	Id string `json:"id" db:"id"`

	// Account holder user id
	AccountHolderID shared.UserID `json:"accountHolderId" db:"account_holder_id"`

	// Partner bank name - e.g. 'bbva'
	BankName BankName `json:"bankName" db:"bank_name"`

	// Partner bank account id
	BankAccountId string `json:"bankAccountId" db:"bank_account_id"`

	// Partner bank account extra data
	BankExtra types.JSONText `json:"-" db:"bank_extra"`

	// Account type e.g. checking or savings
	AccountType string `json:"accountType" db:"account_type"`

	// Account status e.g. active or closed
	AccountStatus string `json:"accountStatus" db:"account_status"`

	// Primary account number
	AccountNumber string `json:"accountNumber" db:"account_number"`

	// Primary account routing number
	RoutingNumber string `json:"routingNumber" db:"routing_number"`

	// Wire Routing Number
	WireRouting *string `json:"wireRouting" db:"wire_routing"`

	// Account alias
	Alias *string `json:"alias" db:"alias"`

	// Available balance for widthdrawals
	AvailableBalance float64 `json:"availableBalance" db:"available_balance"`

	// Current posted balance
	PostedBalance float64 `json:"postedBalance" db:"posted_balance"`

	// Current balance pending in review
	PendingDebitBalance float64 `json:"pendingDebitBalance" db:"-"`

	// Current AvailableBalance - PendingDebitBalance
	ActualBalance float64 `json:"actualBalance" db:"-"`

	// Denominated currency
	Currency Currency `json:"currency" db:"currency"`

	// Interest earned YTD
	InterestYTD float64 `json:"interestYTD" db:"interest_ytd"`

	// Monthly Cycle Start Day (1-28) - used for interest calculation, etc.
	MonthlyCycleStartDay int `json:"monthlyCycleStartDay" db:"monthly_cycle_start_day"`

	// RemainingFundingBalance
	RemainingFundingBalance float64 `json:"remainingFundingBalance" db"-"`

	// FundingLimit
	FundingLimit float64 `json:"fundingLimit" db"-"`

	// Timestamp opened (UTC)
	Opened time.Time `json:"opened" db:"opened"`

	// Timestamp created (UTC)
	Created time.Time `json:"created" db:"created"`

	// Timestamp modified (UTC)
	Modified time.Time `json:"modified" db:"modified"`
}

// BankAccountCreate a bank account creation struct with mandatory fields
type BankAccountCreate struct {
	// Account alias
	Alias *string `json:"alias" db:"alias"`

	// Account type e.g. checking or savings
	AccountType string `json:"accountType" db:"account_type"`

	// Partner bank name - e.g. 'bbva'
	BankName BankName `json:"bankName" db:"bank_name"`
}

//BankAccountUpdate updating struct with minimal update fields
type BankAccountUpdate struct {
	// Bank account id
	Id string `json:"id" db:"id"`

	// Account alias
	Alias string `json:"alias" db:"alias"`

	Status string `json:"status" db:"-"`
}

type BankAccountStatement struct {
	AccountHolderID shared.UserID `json:"accountHolderId"`
	StatementID     string        `json:"statementId"`
	Description     string        `json:"description"`
	Created         time.Time     `json:"created"`
	PageCount       int           `json:"pageCount"`
}

type BankAccountStatementDocument struct {
	ContentType string `json:"contentType"`
	Content     []byte `json:"content"`
}

type AccountBlockStatus string

const (
	AccountBlockStatusActive   = AccountBlockStatus("active")
	AccountBlockStatusCanceled = AccountBlockStatus("canceled")
)

type AccountBlockBankID string

func (id AccountBlockBankID) String() string {
	return string(id)
}

type AccountBlockType string

const (
	AccountBlockTypeDebit  = AccountBlockType("debits")
	AccountBlockTypeCredit = AccountBlockType("credits")
	AccountBlockTypeCheck  = AccountBlockType("checks")
	AccountBlockTypeAll    = AccountBlockType("all")
)

type BankAccountBlockResponse struct {
	BlockID AccountBlockBankID `json:"BlockId"`
	Type    AccountBlockType   `json:"blockType"`
	Status  AccountBlockStatus `json:"status"`
}

type AccountBlock struct {
	ID             string           `json:"id" db:"id"`
	AccountID      string           `json:"accountId" db:"account_id"`
	BlockID        string           `json:"blockId" db:"block_id"`
	BlockType      AccountBlockType `json:"blockType" db:"block_type"`
	Reason         string           `json:"reason" db:"reason"`
	Created        time.Time        `json:"created" db:"created"`
	Deactivated    *time.Time       `json:"deactivated" db:"deactivated"`
	OriginatedFrom string           `json:"originatedFrom" db:"originated_from"`
}

type AccountBlockCreate struct {
	AccountID      string           `json:"accountId" db:"account_id"`
	BlockID        string           `json:"blockId" db:"block_id"`
	BlockType      AccountBlockType `json:"blockType" db:"block_type"`
	Reason         string           `json:"reason" db:"reason"`
	OriginatedFrom string           `json:"originatedFrom" db:"originated_from"`
}
