package analytics

import (
	"time"
)

const (
	ConsumerID            = "consumer_id"
	ConsumerPhone         = "phone"
	ConsumerPhoneVerified = "phone_verified"
	ConsumerFirstName     = "first_name"
	ConsumerMiddleName    = "middle_name"
	ConsumerLastName      = "last_name"
	ConsumerEmail         = "email"
	ConsumerDateOfBirth   = "date_of_birth"
	ConsumerKYCStatus     = "consumer_kyc_status"
	ConsumerJobTitle      = "job_title"

	SubscriptionStatus = "subscription_status"

	BusinessId             = "business_id"
	BusinessLegalName      = "business_legal_name"
	BusinessDBA            = "dba"
	BusinessEntityType     = "entity_type"
	BusinessIndustryType   = "industry_type"
	BusinessKYCStatus      = "business_kyc_status"
	BusinessMailingAddress = "mailing_address"
	BusinessOriginDate     = "date_of_establishment"

	BusinessCardId           = "business_card_id"
	BusinessCardStatus       = "business_card_status"
	BusinessTransactionCount = "business_transaction_count"

	BusinessAccountId               = "business_account_id"
	BusinessAccountType             = "business_account_type"
	BusinessAccountStatus           = "business_account_status"
	BusinessAccountAlias            = "business_account_alias"
	BusinessAccountBankName         = "business_bank_name"
	BusinessAccountOpened           = "business_account_opened"
	BusinessNewPromoFunding         = "new_business_promo_funding"
	BusinessAccountAvailableBalance = "account_available_balance"

	LinkedAccountCount = "linked_accounts"
	LinkedCardCount    = "linked_cards"
	CardReaderCount    = "card_readers"
	ContactCount       = "account_contacts"
	InvoiceCount       = "invoices"

	PostedTransactionCount       = "accounts_posted_transactions"
	PushToDebitTransactionCount  = "card_push_to_debit_count"
	DebitCardATMTransactionCount = "debit_card_atm_count"
	CardReaderTransactionCount   = "card_reader_transaction_count"
	ACHDebitTransactionCount     = "ach_transfer_debit_count"
	ACHCreditTransactionCount    = "ach_transfer_credit_count"
	WireCreditTransactionCount   = "wire_transfer_credit_count"
	DebitCardTransactionCount    = "debit_card_transaction_count"
)

// BusinessPromotionUpdate ...
type CSPBusinessUpdate struct {
	KYCBStatus  *string    `json:"kybStatus"`
	PromoFunded *time.Time `json:"promoFunded"`
	Amount      *float64   `json:"amount"`
}
