/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

// Package for all internal banking related items such as accounts and cards
package banking

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx/types"
	"github.com/wiseco/core-platform/shared"
)

type BlockStatus string

const (
	BlockStatusInactive = BlockStatus("inactive")
	BlockStatusActive   = BlockStatus("active")
)

type BankCard struct {
	// Global business card identifier
	Id string `json:"id" db:"id"`

	// Owner id of this card
	CardholderID shared.UserID `json:"cardholderId" db:"cardholder_id"`

	// Bank account id
	BankAccountId string `json:"bankAccountId" db:"bank_account_id"`

	// Card type (debit or credit)
	CardType CardType `json:"cardType" db:"card_type"`

	// Card holder name
	CardholderName string `json:"cardholderName" db:"cardholder_name"`

	// Is virtual card
	IsVirtual bool `json:"isVirtual" db:"is_virtual"`

	// Partner bank name - e.g. 'bbva'
	BankName BankName `json:"bankName" db:"bank_name"`

	// Partner bank card id
	BankCardId string `json:"-" db:"bank_card_id"`

	// Partner bank account extra data
	BankExtra types.JSONText `json:"-" db:"bank_extra" db:"bank_extra"`

	// Card number
	CardNumberMasked string `json:"cardNumberMasked" db:"card_number_masked"`

	// Card brand
	CardBrand string `json:"cardBrand" db:"card_brand"`

	// Currency
	Currency Currency `json:"currency" db:"currency"`

	// Card status (active, blocked, etc.)
	CardStatus CardStatus `json:"cardStatus" db:"card_status"`

	// Alias or friendly name
	Alias *string `json:"alias" db:"alias"`

	// Daily ATM limit
	DailyATMLimit *float64 `json:"dailyATMLimit" db:"daily_withdrawal_limit"`

	// Daily POS limit
	DailyPOSLimit *float64 `json:"dailyPOSLimit" db:"daily_pos_limit"`

	// Daily transaction limit
	DailyTransactionLimit *int `json:"dailyTransactionLimit" db:"daily_transaction_limit"`

	// Card blocks - Soon to be deprecated
	Block *CardBlockArray `json:"block" db:"card_block"`

	// Date created
	Created time.Time `json:"created" db:"created"`

	// Date last modified
	Modified time.Time `json:"modified" db:"modified"`
}

func (b BankCard) GetCardNumberLastFour() (string, error) {
	if len(b.CardNumberMasked) < 4 {
		return "", fmt.Errorf("Unable to get last four digits of card, CardNumberMasked not long enough. length: %d", len(b.CardNumberMasked))
	}

	return b.CardNumberMasked[len(b.CardNumberMasked)-4:], nil
}

// Consumer debit or credit card
type ConsumerBankCard struct {
	BankCard
}

type BankCardCreate struct {
	// Owner id of this card
	CardholderID shared.UserID `json:"cardholderId"`

	// Bank account id
	BankAccountId string `json:"bankAccountId"`

	// Card type (debit or credit)
	CardType CardType `json:"cardType"`
}

type BankCardActivate struct {
	// Global business card identifier
	Id string `json:"id"`

	// Owner id of this card
	CardholderID shared.UserID `json:"cardholderId"`

	// Card last 6 digit
	PANLast6 string `json:"panLast6"`
}

type BankCardBlockCreate struct {
	// Global business card identifier
	CardID string `json:"cardId" db:"card_id"`

	// Cardholder Id
	CardholderID shared.UserID `json:"cardholderId"`

	// Card block code
	BlockID CardBlockID `json:"blockId" db:"block_id"`

	// Card block reason
	Reason *string `json:"reason" db:"reason"`

	// Originated from - CSP, app etc.
	OriginatedFrom OriginatedFrom `json:"originatedFrom" db:"originated_from"`

	// Block status
	BlockStatus BlockStatus `json:"blockStatus" db:"block_status"`

	// Block date
	BlockDate time.Time `json:"blockDate" db:"block_date"`
}

type BankCardBlockDelete struct {
	ID shared.BankCardBlockID `json:"id" db:"id"`

	// Cardholder Id
	CardholderID shared.UserID `json:"cardholderId"`

	// Global business card identifier
	CardID string `json:"cardId" db:"card_id"`
}

type BankCardBlock struct {
	ID shared.BankCardBlockID `json:"id" db:"id"`

	// Global business card identifier
	CardID string `json:"cardId" db:"card_id"`

	// Card block code
	BlockID CardBlockID `json:"blockId" db:"block_id"`

	// Card block reason
	Reason string `json:"reason" db:"reason"`

	// Originated from - CSP, app etc.
	OriginatedFrom OriginatedFrom `json:"originatedFrom" db:"originated_from"`

	// Block status
	BlockStatus BlockStatus `json:"blockStatus" db:"block_status"`

	// Block date
	BlockDate time.Time `json:"blockDate" db:"block_date"`

	// Created time
	Created time.Time `json:"created" db:"created"`

	// Date last modified
	Modified time.Time `json:"modified" db:"modified"`
}

type BankCardUpdate struct {
	Alias string `json:"alias"`
}

// BankCardMini is for segment integration
type BankCardMini struct {
	BankCardID string `json:"bankCardID"`

	CardStatus CardStatus `json:"cardStatus"`

	DailyTransactionLimit *int `json:"dailyTransactionLimit"`
}

type CardReissueReason string

const (
	CardReissueReasonNotReceived = CardReissueReason("notReceived")
	CardReissueReasonResendPin   = CardReissueReason("resendPin")
)

type CardReissueRequest struct {
	CardID string            `json:"cardId" db:"card_id"`
	Reason CardReissueReason `json:"reason" db:"reason"`
}

type CardReissueHistory struct {
	ID       string            `json:"id" db:"id"`
	CardID   string            `json:"cardId" db:"card_id"`
	Reason   CardReissueReason `json:"reason" db:"reason"`
	Created  time.Time         `json:"created" db:"created"`
	Modified time.Time         `json:"modified" db:"modified"`
}
type OriginatedFrom string

const (
	OriginatedFromCSP               = OriginatedFrom("Customer Success Portal")
	OriginatedFromClientApplication = OriginatedFrom("Client Application")
	OriginatedFromBank              = OriginatedFrom("Bank")
)
