package business

import (
	"time"

	"github.com/wiseco/core-platform/services"
	"github.com/wiseco/core-platform/services/banking"
	"github.com/wiseco/core-platform/shared"
)

// Business transfer request
type MoneyRequest struct {
	BusinessID shared.BusinessID `json:"businessId" db:"business_id"` // Related business account for this card
	banking.MoneyRequest
}

// Wise money request message ID update
type MoneyRequestIDUpdate struct {
	BusinessID shared.BusinessID `json:"businessId" db:"business_id"` // Related business account for this card
	banking.MoneyRequestIDUpdate
}

type RequestInitiate struct {
	// User Id
	CreatedUserID shared.UserID `json:"createdUserId" db:"created_user_id"`

	// Business Id
	BusinessID shared.BusinessID `json:"businessId" db:"business_id"`

	// Name on the account
	ContactId string `json:"contactId" db:"contact_id"`

	// Transfer amount
	Amount float64 `json:"amount" db:"amount"`

	// Denominated currency
	Currency banking.Currency `json:"currency" db:"currency"`

	// Transfer Notes
	Notes string `json:"notes" db:"notes"`

	// Request type
	RequestType *banking.PaymentRequestType `json:"requestType" db:"request_type"`

	// POS ID
	POSID *string `json:"posId" db:"pos_id"`
}

type RequestCancel struct {
	// Transaction id
	Id string `json:"id"`

	// Message
	Message string `json:"message"`
}

type PaymentResponse struct {
	// Request id
	RequestID shared.PaymentRequestID `json:"requestId" db:"request_id"`

	// Source payment id
	SourcePaymentID string `json:"sourcePaymentId" db:"source_payment_id"`

	// Status
	Status string `json:"status" db:"status"`

	// Expiration Date
	ExpirationDate *time.Time `json:"expirationDate" db:"expiration_date"`

	// Payment Intent Client Secret
	ClientSecret string `json:"clientSecret"`
}

type Event struct {
	Id     string `json:"id"`
	Type   string `json:"type"`
	Object Object `json:"data"`
}

type Object struct {
	Payment Payment `json:"object"`
}

type Payment struct {
	Id string `json:"id" db:"id"`

	RequestId string `json:"requestId" db:"request_id"`

	SourcePaymentID *string `json:"sourcePaymentId" db:"source_payment_id"`

	Status string `json:"status" db:"status"`

	// Token
	PaymentToken *string `json:"paymentToken" db:"payment_token"`

	// Expiration Date
	ExpirationDate *time.Time `json:"expirationDate" db:"expiration_date"`

	// Card Brand
	CardBrand *string `json:"cardBrand" db:"card_brand"`

	// Receipt Number
	CardLast4 *string `json:"cardLast4" db:"card_number"`

	// Paid Date
	PaymentDate *time.Time `json:"paidDate" db:"payment_date"`

	ReceiptID *string `json:"receiptId" db:"receipt_id"`

	ReceiptMode *CardReaderReceiptMode `json:"receiptMode" db:"receipt_mode"`

	CustomerContact *string `json:"customerContact" db:"customer_contact"`

	// Receipt token
	ReceiptToken *string `json:"receiptToken" db:"receipt_token"`

	// Purchase address
	PurchaseAddress *services.Address `json:"purchaseAddress" db:"purchase_address"`

	// Linked account ID
	LinkedBankAccountID *string `json:"linkedBankAccountId" db:"linked_bank_account_id"`

	// Created timestamp
	Created time.Time `json:"created" db:"created"`

	// Modified timestamp
	Modified time.Time `json:"modified" db:"modified"`
}

type PaymentStatus string

const (
	PaymentStatusRequiresPaymentMethod = PaymentStatus("requiresPaymentMethod")
	PaymentStatusRequiresConfirmation  = PaymentStatus("requiresConfirmation")
	PaymentStatusRequiresAction        = PaymentStatus("requiresAction")
	PaymentStatusRequiresProcessing    = PaymentStatus("processing")
	PaymentStatusRequiresCapture       = PaymentStatus("requiresCapture")
	PaymentStatusRequiresCanceled      = PaymentStatus("canceled")
	PaymentStatusSucceeded             = PaymentStatus("succeeded")
)

type Invoice struct {
	// Id
	ID string `json:"id" db:"id"`

	// Request Id
	RequestID shared.PaymentRequestID `json:"requestId" db:"request_id"`

	// User Id
	CreatedUserID shared.UserID `json:"createdUserId" db:"created_user_id"`

	// Business Id
	BusinessID shared.BusinessID `json:"businessId" db:"business_id"`

	// Name on the account
	ContactId string `json:"contactId" db:"contact_id"`

	// Invoice number
	InvoiceNumber string `json:"invoiceNumber" db:"invoice_number"`

	// Storage key
	StorageKey string `json:"storageKey" db:"storage_key"`

	// Created timestamp
	Created time.Time `json:"created" db:"created"`

	// Modified timestamp
	Modified time.Time `json:"modified" db:"modified"`
}

type InvoiceCreate struct {

	// Request Id
	RequestID shared.PaymentRequestID `json:"requestId" db:"request_id"`

	// User Id
	CreatedUserID shared.UserID `json:"createdUserId" db:"created_user_id"`

	// Business Id
	BusinessID shared.BusinessID `json:"businessId" db:"business_id"`

	// Name on the account
	ContactId string `json:"contactId" db:"contact_id"`

	// Invoice number
	InvoiceNumber string `json:"invoiceNumber" db:"invoice_number"`

	// Storage key
	StorageKey string `json:"storageKey" db:"storage_key"`
}

type Receipt struct {
	// Id
	ID string `json:"id" db:"id"`

	// Request Id
	RequestID shared.PaymentRequestID `json:"requestId" db:"request_id"`

	// Invoice Id
	InvoiceId string `json:"invoiceId" db:"invoice_id"`

	// User Id
	CreatedUserID shared.UserID `json:"createdUserId" db:"created_user_id"`

	// Business Id
	BusinessID shared.BusinessID `json:"businessId" db:"business_id"`

	// Name on the account
	ContactId string `json:"contactId" db:"contact_id"`

	// Invoice number
	ReceiptNumber string `json:"receiptNumber" db:"receipt_number"`

	// Storage key
	StorageKey string `json:"storageKey" db:"storage_key"`

	// Created timestamp
	Created time.Time `json:"created" db:"created"`

	// Modified timestamp
	Modified time.Time `json:"modified" db:"modified"`
}

type ReceiptCreate struct {
	// Request Id
	RequestID shared.PaymentRequestID `json:"requestId" db:"request_id"`

	// Invoice Id
	InvoiceId string `json:"invoiceId" db:"invoice_id"`

	// User Id
	CreatedUserID shared.UserID `json:"createdUserId" db:"created_user_id"`

	// Business Id
	BusinessID shared.BusinessID `json:"businessId" db:"business_id"`

	// Name on the account
	ContactId string `json:"contactId" db:"contact_id"`

	// Invoice number
	ReceiptNumber string `json:"receiptNumber" db:"receipt_number"`

	// Storage key
	StorageKey string `json:"storageKey" db:"storage_key"`
}

type CardReaderReceiptMode string

const (
	CardReaderReceiptModeEmail = CardReaderReceiptMode("email")
	CardReaderReceiptModeSMS   = CardReaderReceiptMode("sms")
	CardReaderReceiptModeNone  = CardReaderReceiptMode("none")
)
