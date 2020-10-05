/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

package bank

import (
	"time"

	"github.com/jmoiron/sqlx/types"
)

type TransactionType string

const (
	// ACH
	TransactionTypeACH = TransactionType("ach")

	// Adjustment
	TransactionTypeAdjustment = TransactionType("adjustment")

	// ATM
	TransactionTypeATM = TransactionType("atm")

	// Check
	TransactionTypeCheck = TransactionType("check")

	// Deposit
	TransactionTypeDeposit = TransactionType("deposit")

	// Fee
	TransactionTypeFee = TransactionType("fee")

	// Other Credit
	TransactionTypeOtherCredit = TransactionType("otherCredit")

	// Other Debit
	TransactionTypeOtherDebit = TransactionType("otherDebit")

	// Purchase transaction
	TransactionTypePurchase = TransactionType("purchase")

	// Refund
	TransactionTypeRefund = TransactionType("refund")

	// Return
	TransactionTypeReturn = TransactionType("return")

	// Reversal
	TransactionTypeReversal = TransactionType("reversal")

	// Transfer
	TransactionTypeTransfer = TransactionType("transfer")

	// Visa credit
	TransactionTypeVisaCredit = TransactionType("visaCredit")

	// Withdrawal
	TransactionTypeWithdrawal = TransactionType("withdrawal")

	// Other
	TransactionTypeOther = TransactionType("other")
)

type TransactionCode string

const (
	TransactionCodeStatusChange   = TransactionCode("statusChange")
	TransactionCodeAuthApproved   = TransactionCode("authApproved")
	TransactionCodeAuthDeclined   = TransactionCode("authDeclined")
	TransactionCodeAuthReversed   = TransactionCode("authReversed")
	TransactionCodeHoldApproved   = TransactionCode("holdApproved")
	TransactionCodeHoldReleased   = TransactionCode("holdReleased")
	TransactionCodeDebitPosted    = TransactionCode("debitPosted")
	TransactionCodeCreditPosted   = TransactionCode("creditPosted")
	TransactionCodeTransferChange = TransactionCode("transferChange")
)

type HoldTransaction struct {
	// Hold number
	Number string `json:"number"`

	// Transaction amount
	Amount float64 `json:"amount"`

	// Date of hold
	Date time.Time `json:"date"`

	// Date of hold expiry
	ExpiryDate time.Time `json:"expiryDate"`
}

type CardTransactionNetwork string

const (
	CardTransactionNetworkVisa = CardTransactionNetwork("visa")
)

type CardTransaction struct {
	// Card specific transaction id
	CardTransactionID string `json:"cardTransactionId"`

	// Network used for card transaction
	TransactionNetwork CardTransactionNetwork `json:"transactionNetwork"`

	// Authorization amount
	AuthAmount float64 `json:"authAmount"`

	// Authorization date
	AuthDate time.Time `json:"authDate"`

	// Authorization response code
	AuthResponseCode string `json:"authResponseCode"`

	// Authorization number
	AuthNumber string `json:"authNumber"`

	// Card transaction code
	TransactionType string `json:"transactionType"`

	// Local currency amount
	LocalAmount float64 `json:"localAmount"`

	// Local currency
	LocalCurrency Currency `json:"localCurrency"`

	// Local date of card transaction
	LocalDate time.Time `json:"localDate"`

	// Billing Currency
	BillingCurrency Currency `json:"billingCurrency"`

	// Point of sale entry mode describes how the card was entered or used - e.g. swipe or chip
	POSEntryMode string `json:"posEntryMode"`

	// Point of sale condition
	POSConditionCode string `json:"posConditionCode"`

	// Acquiring bank identification number (BIN)
	AcquirerBIN string `json:"acquirerBIN"`

	// Merchant (acceptor) id
	MerchantID string `json:"merchantId"`

	// Merchant category code describes the merchants transaction category - e.g. restaurants
	MerchantCategoryCode string `json:"merchantCategoryCode"`

	// Merchant (acceptor) terminal
	MerchantTerminal string `json:"merchantTerminal"`

	// Merchant (acceptor) name
	MerchantName string `json:"merchantName"`

	// Merchant (acceptor) address
	MerchantStreetAddress string `json:"merchantStreetAddress"`

	// Merchant (acceptor) city
	MerchantCity string `json:"merchantCity"`

	// Merchant (acceptor) state
	MerchantState string `json:"merchantState"`

	// Merchant (acceptor) country
	MerchantCountry string `json:"merchantCountry"`
}

type TransactionID string

type Transaction struct {
	// Bank name - e.g. 'bbva'
	BankName ProviderName `json:"bankName"`

	// Bank account id
	BankTransactionID string `json:"bankTransactionId"`

	// Partner bank account extra data
	BankExtra types.JSONText `json:"-"`

	// Transaction type
	TransactionType TransactionType `json:"type"`

	// Bank account id
	AccountID *string `json:"accountId"`

	// Card id
	CardID *string `json:"cardId"`

	// Transaction code
	CodeType TransactionCode `json:"codeType"`

	// Amount
	Amount float64 `json:"amount"`

	// Posted Balance
	PostedBalance *float64 `json:"postedBalance"`

	// Currency
	Currency Currency `json:"currency"`

	// Card transaction details
	CardTransaction *CardTransaction `json:"cardTransaction"`

	// Card hold details
	HoldTransaction *HoldTransaction `json:"holdTransaction"`

	// Money transfer id
	BankMoneyTransferID *string `json:"bankMoneyTransferId"`

	// Bank transaction desc
	BankTransactionDesc *string `json:"bankTransactionDesc"`

	// Transaction Date Created
	TransactionDate time.Time `json:"transactionDate"`
}
