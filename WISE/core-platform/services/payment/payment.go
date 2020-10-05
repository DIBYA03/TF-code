/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

package payment

import (
	"time"

	"github.com/wiseco/go-lib/id"

	"github.com/wiseco/core-platform/services"
	"github.com/wiseco/core-platform/shared"
	"github.com/wiseco/go-lib/num"
)

type Currency string

const (
	// ISO currency code
	CurrencyUSD = Currency("usd")
)

type PaymentRequestType string

const (
	PaymentRequestTypePOS                = PaymentRequestType("pos")
	PaymentRequestTypeInvoiceCard        = PaymentRequestType("invoiceCard")
	PaymentRequestTypeInvoiceBank        = PaymentRequestType("invoiceBank")
	PaymentRequestTypeInvoiceCardAndBank = PaymentRequestType("invoiceCardAndBank")
	PaymentRequestTypeInvoiceNone        = PaymentRequestType("invoiceNone")
)

type PaymentRequestStatus string

const (
	PaymentRequestStatusPending   = PaymentRequestStatus("pending")
	PaymentRequestStatusFailed    = PaymentRequestStatus("failed")
	PaymentRequestStatusComplete  = PaymentRequestStatus("complete")
	PaymentRequestStatusInProcess = PaymentRequestStatus("inProcess")
	PaymentRequestStatusCanceled  = PaymentRequestStatus("canceled")
)

var MoveMoneyStatusToRequestStatus = map[string]PaymentRequestStatus{
	"posted":    PaymentRequestStatusComplete,
	"inProcess": PaymentRequestStatusInProcess,
}

type PaymentStatus string

const (
	PaymentStatusRequiresPaymentMethod = PaymentStatus("requiresPaymentMethod")
	PaymentStatusRequiresConfirmation  = PaymentStatus("requiresConfirmation")
	PaymentStatusRequiresAction        = PaymentStatus("requiresAction")
	PaymentStatusRequiresProcessing    = PaymentStatus("processing")
	PaymentStatusRequiresCapture       = PaymentStatus("requiresCapture")
	PaymentStatusPending               = PaymentStatus("pending")
	PaymentStatusCanceled              = PaymentStatus("canceled")
	PaymentStatusSucceeded             = PaymentStatus("succeeded")
)

var PartnerTransferStatusToPaymentStatus = map[string]PaymentStatus{
	"posted":                PaymentStatusSucceeded,
	"inProcess":             PaymentStatusRequiresProcessing,
	"requiresPaymentMethod": PaymentStatusRequiresPaymentMethod,
	"requiresConfirmation":  PaymentStatusRequiresConfirmation,
	"requiresAction":        PaymentStatusRequiresAction,
	"processing":            PaymentStatusRequiresProcessing,
	"requiresCapture":       PaymentStatusRequiresCapture,
	"canceled":              PaymentStatusCanceled,
	"succeeded":             PaymentStatusSucceeded,
}

type RequestSource string

const (
	RequestSourceShopify = RequestSource("shopify")
)

type WalletType string

const (
	GooglePayWallet = WalletType("google_pay")
	ApplePayWallet  = WalletType("apple_pay")
)

var PaymentWalletTypeMap = map[WalletType]string{
	GooglePayWallet: "Google",
	ApplePayWallet:  "Apple",
}

// Wise payment request
type Request struct {
	// Transaction id
	ID shared.PaymentRequestID `json:"id" db:"id"`

	BusinessID shared.BusinessID `json:"businessId" db:"business_id"`

	// Created user id
	CreatedUserID shared.UserID `json:"createdUserId" db:"created_user_id"`

	// Contact id
	ContactID *string `json:"contactId" db:"contact_id"`

	// Transfer amount
	Amount float64 `json:"amount" db:"amount"`

	// Transfer currency
	Currency Currency `json:"currency" db:"currency"`

	// Transfer Notes
	Notes *string `json:"notes" db:"notes"`

	// Transfer status
	Status *PaymentRequestStatus `json:"status" db:"request_status"`

	// Message id
	MessageID string `json:"messageId" db:"message_id"`

	// Payment Intent Secret - used only for 'POS' request type
	PaymentIntentToken string `json:"paymentIntentToken"`

	// Request type
	RequestType *PaymentRequestType `json:"requestType" db:"request_type"`

	// POS ID
	CardReaderID *shared.CardReaderID `json:"cardReaderId" db:"pos_id"`

	// IP address
	IPAddress *string `json:"-" db:"request_ip_address"`

	// Request Source - ex: shopify
	RequestSource *RequestSource `json:"requestSource" db:"request_source"`

	// Request Source ID - ex: shopify order Id
	RequestSourceID *string `json:"requestSourceId" db:"request_source_id"`

	// Created date
	Created time.Time `json:"created" db:"created"`

	// Modified date
	Modified time.Time `json:"modified" db:"modified"`
}

type RequestInitiate struct {
	// User Id
	CreatedUserID shared.UserID `json:"createdUserId" db:"created_user_id"`

	// Business Id
	BusinessID shared.BusinessID `json:"businessId" db:"business_id"`

	// Contact Id
	ContactID *string `json:"contactId" db:"contact_id"`

	// Transfer amount
	Amount float64 `json:"amount" db:"amount"`

	// Denominated currency
	Currency Currency `json:"currency" db:"currency"`

	// Transfer Notes
	Notes *string `json:"notes" db:"notes"`

	// Request type
	RequestType PaymentRequestType `json:"requestType" db:"request_type"`

	// POS ID
	CardReaderID *shared.CardReaderID `json:"cardReaderId" db:"pos_id"`

	// IP address
	IPAddress *string `json:"-" db:"request_ip_address"`

	// Request Source - ex: shopify
	RequestSource *RequestSource `json:"requestSource" db:"request_source"`

	// Request Source ID - ex: shopify order Id
	RequestSourceID *string `json:"requestSourceId" db:"request_source_id"`
}

func (r Request) isShopifyRequest() bool {
	if r.RequestSource != nil && *r.RequestSource == RequestSourceShopify {
		return true
	}

	return false
}

// Wise payment request status update
type RequestStatusUpdate struct {
	// Request id
	ID shared.PaymentRequestID `json:"id" db:"id"`

	// Business ID
	BusinessID shared.BusinessID `json:"businessId" db:"business_id"`

	// Transfer status
	Status PaymentRequestStatus `json:"status" db:"request_status"`
}

type RequestMetadata struct {
	BusinessName     string
	FirstName        string  `db:"first_name"`
	MiddleName       *string `db:"middle_name"`
	LastName         string  `db:"last_name"`
	AvailableBalance float64 `db:"available_balance"`
	IPAddress        *string `db:"request_ip_address"`
}

// Wise payment request update
type RequestUpdate struct {
	// Request id
	ID shared.PaymentRequestID `json:"id" db:"id"`

	// Transfer status
	Status PaymentRequestStatus `json:"status" db:"request_status"`

	// Request type
	RequestType *PaymentRequestType `json:"requestType" db:"request_type"`
}

// Wise payment request message ID update
type RequestIDUpdate struct {
	// Request id
	ID shared.PaymentRequestID `json:"id" db:"id"`

	BusinessID shared.BusinessID `json:"businessId" db:"business_id"`

	// Contact id
	ContactID *string `json:"contactId" db:"contact_id"`

	// Message ID
	MessageID string `json:"messageId" db:"message_id"`
}

type RequestCancel struct {
	// Transaction id
	ID string `json:"id"`

	// Message
	Message string `json:"message"`
}

type Payment struct {
	ID string `json:"id" db:"id"`

	RequestID *shared.PaymentRequestID `json:"requestId" db:"request_id"`

	SourcePaymentID *string `json:"sourcePaymentId" db:"source_payment_id"`

	Status PaymentStatus `json:"status" db:"status"`

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

	// Fee amount
	FeeAmount *num.Decimal `json:"feeAmount" db:"fee_amount"`

	// Invoice ID
	InvoiceID *id.InvoiceID `json:"invoiceId" db:"invoice_id"`

	// Wallet type
	WalletType *string `json:"walletType" db:"wallet_type"`

	// Created timestamp
	Created time.Time `json:"created" db:"created"`

	// Modified timestamp
	Modified time.Time `json:"modified" db:"modified"`
}

// Response from stripe
type PartnerPaymentResponse struct {
	// Request id
	RequestID shared.PaymentRequestID `json:"requestId" db:"request_id"`

	// Invoice id
	InvoiceID string `json:"invoiceId" db:"invoice_id"`

	// Source payment id
	SourcePaymentID string `json:"sourcePaymentId" db:"source_payment_id"`

	// Status
	Status string `json:"status" db:"status"`

	// Expiration Date
	ExpirationDate *time.Time `json:"expirationDate" db:"expiration_date"`

	// Payment Intent Client Secret
	ClientSecret string `json:"clientSecret"`
}

// Response from payments table
type PaymentResponse struct {
	BusinessID            shared.BusinessID        `db:"business.id"`
	UserID                shared.UserID            `db:"owner_id"`
	ContactID             *string                  `db:"contact_id"`
	PaymentID             string                   `db:"business_money_request_payment.id"`
	MoneyRequestID        *shared.PaymentRequestID `db:"business_money_request_payment.request_id"`
	BusinessBankAccountID *string                  `db:"business_bank_account_id"`
	RegisteredAccountID   string                   `db:"business_linked_bank_account.id"`
	LegalName             *string                  `db:"legal_name"`
	DBA                   services.StringArray     `db:"dba"`
	BusinessName          string
	RequestType           PaymentRequestType `db:"request_type"`
	Amount                float64            `db:"amount"`
	Notes                 string             `db:"notes"`
	ClientSecret          *string
	StripeKey             string
	Status                PaymentRequestStatus `db:"request_status"`
	PaymentDate           *time.Time           `db:"business_money_request_payment.payment_date"`
	InvoiceID             *id.InvoiceID        `db:"invoice_id"`
	InvoiceIPAddress      *string
	InvoiceTitle          *string
}

type Event struct {
	ID     string `json:"id"`
	Type   string `json:"type"`
	Object Object `json:"data"`
}

type Object struct {
	Payment Payment `json:"object"`
}

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
	ContactID string `json:"contactId" db:"contact_id"`

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
	RequestID *shared.PaymentRequestID `json:"requestId" db:"request_id"`

	// User Id
	CreatedUserID shared.UserID `json:"createdUserId" db:"created_user_id"`

	// Business Id
	BusinessID shared.BusinessID `json:"businessId" db:"business_id"`

	// Name on the account
	ContactID string `json:"contactId" db:"contact_id"`

	// Invoice number
	InvoiceNumber string `json:"invoiceNumber" db:"invoice_number"`

	// Storage key
	StorageKey string `json:"storageKey" db:"storage_key"`

	InvoiceIdV2 *id.InvoiceID `json:"invoiceIdV2" db:"invoice_id_v2"`
}

type Receipt struct {
	// Id
	ID string `json:"id" db:"id"`

	// Request Id
	RequestID *shared.PaymentRequestID `json:"requestId" db:"request_id"`

	// Invoice Id
	InvoiceID *string `json:"invoiceId" db:"invoice_id"`

	// User Id
	CreatedUserID shared.UserID `json:"createdUserId" db:"created_user_id"`

	// Business Id
	BusinessID shared.BusinessID `json:"businessId" db:"business_id"`

	// Name on the account
	ContactID *string `json:"contactId" db:"contact_id"`

	// Invoice number
	ReceiptNumber string `json:"receiptNumber" db:"receipt_number"`

	// Storage key
	StorageKey string `json:"storageKey" db:"storage_key"`

	// Invoice id for invoice service
	InvoiceIdV2 *id.InvoiceID `json:"invoiceIdV2" db:"invoice_id_v2"`

	// Created timestamp
	Created time.Time `json:"created" db:"created"`

	// Modified timestamp
	Modified time.Time `json:"modified" db:"modified"`
}

type ReceiptCreate struct {
	// Request Id
	RequestID *shared.PaymentRequestID `json:"requestId" db:"request_id"`

	// Invoice Id
	InvoiceID *string `json:"invoiceId" db:"invoice_id"`

	// User Id
	CreatedUserID shared.UserID `json:"createdUserId" db:"created_user_id"`

	// Business Id
	BusinessID shared.BusinessID `json:"businessId" db:"business_id"`

	// Name on the account
	ContactID *string `json:"contactId" db:"contact_id"`

	// Invoice number
	ReceiptNumber string `json:"receiptNumber" db:"receipt_number"`

	// Storage key
	StorageKey string `json:"storageKey" db:"storage_key"`

	// Invoice id for invoice service
	InvoiceIdV2 *id.InvoiceID `json:"invoiceIdV2" db:"invoice_id_v2"`
}

type PaymentConnectionRequest struct {
	BusinessID shared.BusinessID `json:"businessId"`
	UserID     shared.UserID     `json:"userId"`
}

type PaymentConnectionResponse struct {
	ConnectionToken string `json:"connectionToken"`
}

type PaymentCaptureRequest struct {
	BusinessID      shared.BusinessID       `json:"businessId"`
	RequestID       shared.PaymentRequestID `json:"requestId"`
	CreatedUserID   shared.UserID           `json:"userId"`
	PurchaseAddress *services.Address       `json:"purchaseAddress" db:"purchase_address"`
}

type CardReaderReceiptMode string

const (
	CardReaderReceiptModeEmail = CardReaderReceiptMode("email")
	CardReaderReceiptModeSMS   = CardReaderReceiptMode("sms")
	CardReaderReceiptModeNone  = CardReaderReceiptMode("none")
)

type CardReaderReceiptCreate struct {
	BusinessID      shared.BusinessID       `json:"businessId" db:"business_id"`
	RequestID       shared.PaymentRequestID `json:"requestId" db:"request_id"`
	ReceiptMode     *CardReaderReceiptMode  `json:"receiptMode" db:"receipt_mode"`
	CustomerContact *string                 `json:"customerContact" db:"customer_contact"`
	ReceiptID       *string                 `json:"receiptId" db:"receipt_id"`
	ReceiptToken    *string                 `json:"receiptToken" db:"receipt_token"`
}

type CardReaderReceiptResponse struct {
	BusinessID      shared.BusinessID    `db:"business_id"`
	ReceiptID       string               `db:"receipt_id"`
	LegalName       *string              `db:"legal_name"`
	DBA             services.StringArray `db:"dba"`
	BusinessName    string
	Amount          float64           `db:"amount"`
	PaymentDate     time.Time         `db:"payment_date"`
	CardBrand       *string           `db:"card_brand"`
	CardLast4       *string           `db:"card_number"`
	PurchaseAddress *services.Address `db:"purchase_address"`
	ClientSecret    string
	StripeKey       string
}

type InvoiceMini struct {
	Id            *string `json:"id,omitempty" db:"id"`
	InvoiceNumber *string `json:"invoiceNumber,omitempty" db:"invoice_number"`
	ViewLink      *string `json:"viewLink,omitempty"`
}

type ReceiptMini struct {
	Id            *string `json:"id,omitempty" db:"id"`
	ReceiptNumber *string `json:"receiptNumber,omitempty" db:"receipt_number"`
	ViewLink      *string `json:"viewLink,omitempty"`
}

type RequestMini struct {
	Id     shared.PaymentRequestID `db:"id"`
	Amount float64                 `db:"amount"`
}

type PaymentMini struct {
	Id              *string           `json:"id,omitempty" db:"id"`
	RequestType     *string           `json:"requestType" db:"business_money_request.request_type"`
	ReceiptId       *string           `json:"receiptId,omitempty" db:"receipt_id"`
	PaymentDate     *time.Time        `json:"paymentDate,omitempty" db:"payment_date"`
	CardBrand       *string           `json:"cardBrand,omitempty" db:"card_brand"`
	CardLast4       *string           `json:"cardLast4,omitempty" db:"card_number"`
	PurchaseAddress *services.Address `json:"purchaseAddress,omitempty" db:"purchase_address"`
	FeeAmount       num.Decimal       `json:"feeAmount,omitempty" db:"fee_amount"`
	WalletType      *string           `json:"walletType,omitempty" db:"wallet_type"`
}

type RequestPayment struct {
	Request       Request                     `json:"request" db:"business_money_request"`
	PaymentDate   *time.Time                  `json:"paymentDate" db:"payment_date"`
	TransactionID *shared.PostedTransactionID `json:"transactionId" db:"posted_credit_transaction_id"`
}

type ReceiptGenerate struct {
	ContactFirstName    *string
	ContactLastName     *string
	ContactBusinessName *string
	ContactEmail        *string
	PaymentDate         *time.Time
	PaymentBrand        *string // VISA, Bank name, etc..
	PaymentNumber       *string // Account number or card number
	ReceiptNumber       string
	InvoiceNumber       *string
	BusinessName        string
	BusinessPhone       string
	Amount              float64
	Notes               *string
	UserID              shared.UserID
	BusinessID          shared.BusinessID
	ContactID           *string
	InvoiceID           *string
	RequestID           *shared.PaymentRequestID
	InvoiceIdV2         *id.InvoiceID
}

type ReceiptRequest struct {
	RequestSource       *RequestSource
	ContactFirstName    *string
	ContactLastName     *string
	ContactBusinessName *string
	ContactEmail        *string
	ContactPhone        *string
	BusinessName        string
	BusinessEmail       string
	Notes               *string
	Amount              float64
	ReceiptNumber       string
	PaymentDate         string
	Content             *string
}

type InvoicePaymentReceipt struct {
	CardBrand     *string       `db:"card_brand"`
	CardLast4     *string       `db:"card_number"`
	WalletType    *WalletType   `db:"wallet_type"`
	ReceiptNumber *string       `db:"receipt_number"`
	Status        PaymentStatus `db:"status"`
	PaymentDate   *time.Time    `json:"paidDate" db:"payment_date"`
}

func (r ReceiptRequest) isShopifyRequest() bool {
	if r.RequestSource != nil && *r.RequestSource == RequestSourceShopify {
		return true
	}

	return false
}
