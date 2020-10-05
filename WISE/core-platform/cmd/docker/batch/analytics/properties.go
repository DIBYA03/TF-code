package main

import (
	"time"

	"github.com/wiseco/core-platform/services"
	"github.com/wiseco/core-platform/shared"
)

type Business struct {
	ID         shared.BusinessID `db:"business.id"`
	UserID     shared.UserID     `db:"business.owner_id"`
	ConsumerID string            `db:"wise_user.consumer_id"`

	IndustryType       *string              `db:"industry_type"`
	EntityType         string               `db:"entity_type"`
	LegalName          *string              `db:"legal_name"`
	DBA                services.StringArray `db:"dba"`
	BusinessKYCStatus  string               `db:"business.kyc_status"`
	MailingAddress     *string              `db:"business.mailing_address"`
	BusinessOriginDate *time.Time           `db:"origin_date"`

	ConsumerFirstName  string    `db:"consumer.first_name"`
	ConsumerMiddleName string    `db:"consumer.middle_name"`
	ConsumerLastName   string    `db:"consumer.last_name"`
	ConsumerJobTitle   string    `db:"title_type"`
	ConsumerDOB        time.Time `db:"date_of_birth"`
	ConsumerKYCStatus  string    `db:"consumer.kyc_status"`
	PhoneVerified      bool      `db:"phone_verified"`

	DebitCardStatus       string `db:"card_status"`
	DebitCardID           string `db:"business_bank_card.id"`
	DailyTransactionLimit int    `db:"daily_transaction_limit"`

	BankName         string
	AccountType      string    `db:"account_type"`
	AccountStatus    string    `db:"account_status"`
	AccountOpened    time.Time `db:"opened"`
	AvailableBalance float64   `db:"available_balance"`
	PostedBalance    float64   `db:"posted_balance"`
	AccountID        string    `db:"business_bank_account.id"`

	LinkedAccountCount *int `db:"linked_account_count"`
	LinkedCardCount    *int `db:"linked_card_count"`
	CardReaderCount    *int `db:"card_reader_count"`
	ContactCount       *int `db:"contact_count"`
	InvoiceCount       *int `db:"invoice_count"`

	// Transaction details
	PostedTransactionCount       int `db:"transaction_count"`
	PushToDebitTransactionCount  int `db:"push_to_debit_transaction_count"`
	DebitCardATMTransactionCount int `db:"debit_card_atm_transaction_count"`
	CardReaderTransactionCount   int `db:"card_reader_transaction_count"`
	ACHDebitTransactionCount     int `db:"ach_debit_transaction_count"`
	ACHCreditTransactionCount    int `db:"ach_credit_transaction_count"`
	WireCreditTransactionCount   int `db:"wire_credit_transaction_count"`
	DebitCardTransactionCount    int `db:"debit_card_transaction_count"`
}
