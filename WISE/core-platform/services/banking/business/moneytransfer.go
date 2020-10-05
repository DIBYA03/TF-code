package business

import (
	"time"

	"github.com/jmoiron/sqlx/types"
	"github.com/wiseco/core-platform/services/banking"
	"github.com/wiseco/core-platform/shared"
)

// Business debit or credit card
type MoneyTransfer struct {
	BusinessID shared.BusinessID `json:"businessId" db:"business_id"` // Related business account for this card

	PostedDebitTransactionID *string `json:"postedDebitTransactionId" db:"posted_debit_transaction_id"`

	PostedCreditTransactionID *string `json:"-" db:"posted_credit_transaction_id"`

	MoneyRequestID *shared.PaymentRequestID `json:"-" db:"money_request_id"`

	// Interest ID if interest payment
	MonthlyInterestID *string `json:"-" db:"account_monthly_interest_id"`

	ErrorCause string `db:"-"`

	banking.MoneyTransfer
}

type TransferInitiate struct {
	// User Id
	CreatedUserID shared.UserID `json:"createdUserId" db:"created_user_id"`

	// Business Id
	BusinessID shared.BusinessID `json:"businessId" db:"business_id"`

	// Name on the account
	ContactId *string `json:"contactId" db:"contact_id"`

	// Address ID -- used with Pay by check Desttype == TransferTypeCheck
	AddressID shared.AddressID `json:"addressId" db:"address_id"`

	// Source bank account id
	SourceAccountId string `json:"sourceAccountId" db:"source_account_id"`

	// Ex - card or account
	SourceType banking.TransferType `json:"sourceType" db:"source_type"`

	// Destination bank account id
	DestAccountId string `json:"destAccountId" db:"dest_account_id"`

	// Ex - card or account
	DestType banking.TransferType `json:"destType" db:"dest_type"`

	// Transfer amount
	Amount float64 `json:"amount" db:"amount"`

	// Denominated currency
	Currency banking.Currency `json:"currency" db:"currency"`

	// Transfer Notes
	Notes *string `json:"notes" db:"notes"`

	// Send email
	SendEmail bool `json:"sendEmail" db:"send_email"`

	// Request ID
	MoneyRequestID *shared.PaymentRequestID `json:"moneyRequestId" db:"money_request_id"`

	// Interest ID if interest payment
	MonthlyInterestID *string `json:"-" db:"account_monthly_interest_id"`

	// CVC code
	CVVCode *string `json:"cvvCode"`
}

type TransferCancel struct {
	// Transaction id
	Id string `json:"id"`

	// Message
	Message string `json:"message"`
}

// Pending money transfer notification
type Notification struct {
	ID         string           `json:"id" `
	EntityID   string           `json:"entityId"`
	EntityType EntityType       `json:"entityType"`
	BankName   banking.BankName `json:"bankName"`
	SourceID   string           `json:"sourceId"`
	Type       Type             `json:"type" `
	Attribute  *Attribute       `json:"attribute"`
	Action     Action           `json:"action"`
	Version    string           `json:"version"`
	Created    time.Time        `json:"created"`
	Data       types.JSONText   `json:"data"`
}

type Type string
type EntityType string
type Attribute string
type Action string
type Status string

const (
	EntityTypeConsumer = EntityType("consumer")
	EntityTypeBusiness = EntityType("business")
	EntityTypeMember   = EntityType("member")
)

const PendingMoneyTransfer = "pendingTransfer"
const PendingMoneyTransferType = "transfer"
const NotificationVersion = "1.0.0"

type PendingTransferNotification struct {
	// Bank name - e.g. 'bbva'
	BankName banking.BankName `json:"bankName"`

	// Transaction type
	TransactionType string `json:"type"`

	// Bank account id
	BankAccountID string `json:"accountId"`

	// Transaction status
	Status string `json:"status"`

	// Amount
	Amount float64 `json:"amount"`

	// Notes
	Notes *string `json:"notes"`

	// Partner name - eg: plaid, stripe, bbva, etc.
	ParterName banking.PartnerName `json:"partnerName"`

	// Currency
	Currency banking.Currency `json:"currency"`

	// Money transfer id
	MoneyTransferID *string `json:"moneyTransferId"`

	// Money request id
	MoneyRequestID *shared.PaymentRequestID `json:"moneyRequestId"`

	// Contact id
	ContactID *string `json:"contactId"`

	CodeType string `json:"codeType"`

	// Transaction Date Created
	TransactionDate time.Time `json:"transactionDate"`
}

const (
	DebitInProcess  = "debitInProcess"
	CreditInProcess = "creditInProcess"
)
