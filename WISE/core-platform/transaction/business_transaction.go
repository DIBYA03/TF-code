package transaction

import (
	"time"

	grpcTypes "github.com/golang/protobuf/ptypes"
	"github.com/jmoiron/sqlx/types"
	partnerbank "github.com/wiseco/core-platform/partner/bank"
	"github.com/wiseco/core-platform/services"
	bsrv "github.com/wiseco/core-platform/services/business"
	"github.com/wiseco/core-platform/services/payment"
	"github.com/wiseco/core-platform/shared"
	"github.com/wiseco/go-lib/id"
	"github.com/wiseco/go-lib/num"
	grpcShopify "github.com/wiseco/protobuf/golang/shopping/shopify"
	grpcTxn "github.com/wiseco/protobuf/golang/transaction"
	grpcBankTxn "github.com/wiseco/protobuf/golang/transaction/bank"
)

//Currency ..
type Currency string

const (
	CurrencyUSD = Currency("USD")
)

//SourceType the source type of the transaction e.g `card`
type SourceType string

const (
	//SourceCard source of type `card`
	SourceCard = SourceType("card")

	//SourceAccount source of type `account`
	SourceAccount = SourceType("account")
)

type TransactionFilterType string

const (
	//CreditTransactionFilter
	CreditTransactionFilterType = TransactionFilterType("creditPosted")

	//DebitTransactionFilter
	DebitTransactionFilterType = TransactionFilterType("debitPosted")
)

// BusinessCardTransaction ..
type BusinessCardTransaction struct {
	ID string `json:"id" db:"id"`

	// Related transaction
	TransactionID string `json:"transactionId" db:"transaction_id"`

	// Card holder ID
	CardHolderID shared.UserID `json:"cardholderID" db:"cardholder_id"`

	// Card specific transaction id
	CardTransactionID string `json:"cardTransactionId" db:"card_transaction_id"`

	// Network used for card transaction
	TransactionNetwork string `json:"transactionNetwork" db:"transaction_network"`

	// Authorization amount
	AuthAmount    num.Decimal `json:"authAmount" db:"auth_amount"`
	AuthAmountDep float64     `json:"authAmountDep" db:"auth_amount_dep"` // Deprecated

	// Authorization date
	AuthDate time.Time `json:"authDate" db:"auth_date"`

	// Authorization response code
	AuthResponseCode string `json:"authResponseCode" db:"auth_response_code"`

	// Authorization number
	AuthNumber string `json:"authNumber" db:"auth_number"`

	// Card transaction code
	TransactionType CardTransactionType `json:"transactionType" db:"transaction_type"`

	// Local currency amount
	LocalAmount    num.Decimal `json:"localAmount" db:"local_amount"`
	LocalAmountDep float64     `json:"localAmountDep" db:"local_amount_dep"` // Deprecated

	// Local currency
	LocalCurrency string `json:"localCurrency" db:"local_currency"`

	// Local date of card transaction
	LocalDate time.Time `json:"localDate" db:"local_date"`

	// Billing currency
	BillingCurrency string `json:"billingCurrency" db:"billing_currency"`

	// Point of sale entry mode describes how the card was entered or used - e.g. swipe or chip
	POSEntryMode string `json:"posEntryMode" db:"pos_entry_mode"`

	// Point of sale condition
	POSConditionCode string `json:"posConditionCode" db:"pos_condition_code"`

	// Acquiring bank identification number (BIN)
	AcquirerBIN string `json:"acquirerBIN" db:"acquirer_bin"`

	// Merchant (acceptor) id
	MerchantID string `json:"merchantId" db:"merchant_id"`

	// Merchant category code describes the merchants transaction category - e.g. restaurants
	MerchantCategoryCode string `json:"merchantCategoryCode" db:"merchant_category_code"`

	// Merchant (acceptor) terminal
	MerchantTerminal string `json:"merchantTerminal" db:"merchant_terminal"`

	// Merchant (acceptor) name
	MerchantName string `json:"merchantName" db:"merchant_name"`

	// Merchant (acceptor) address
	MerchantStreetAddress string `json:"merchantStreetAddress" db:"merchant_street_address"`

	// Merchant (acceptor) city
	MerchantCity string `json:"merchantCity" db:"merchant_city"`

	// Merchant (acceptor) state
	MerchantState string `json:"merchantState" db:"merchant_state"`

	// Merchant (acceptor) country
	MerchantCountry string `json:"merchantCountry" db:"merchant_country"`

	// Created
	Created time.Time `json:"created" db:"created"`

	// Modified
	Modified *time.Time `json:"modified" db:"modified"`
}

type BusinessCardTransactionCreate struct {
	// Related transaction
	TransactionID string `json:"transactionId,omitempty" db:"transaction_id"`

	// Card holder ID
	CardHolderID shared.UserID `json:"cardholderID,omitempty" db:"cardholder_id"`

	// Card specific transaction id
	CardTransactionID string `json:"cardTransactionId" db:"card_transaction_id"`

	// Network used for card transaction
	TransactionNetwork string `json:"transactionNetwork" db:"transaction_network"`

	// Authorization amount
	AuthAmount num.Decimal `json:"authAmount" db:"auth_amount"`

	// Authorization date
	AuthDate time.Time `json:"authDate" db:"auth_date"`

	// Authorization response code
	AuthResponseCode string `json:"authResponseCode" db:"auth_response_code"`

	// Authorization number
	AuthNumber string `json:"authNumber" db:"auth_number"`

	// Card transaction code
	TransactionType CardTransactionType `json:"transactionType" db:"transaction_type"`

	// Local currency amount
	LocalAmount num.Decimal `json:"localAmount" db:"local_amount"`

	// Local currency
	LocalCurrency string `json:"localCurrency" db:"local_currency"`

	// Local date of card transaction
	LocalDate time.Time `json:"localDate" db:"local_date"`

	// Billing currency
	BillingCurrency string `json:"billingCurrency" db:"billing_currency"`

	// Point of sale entry mode describes how the card was entered or used - e.g. swipe or chip
	POSEntryMode string `json:"posEntryMode" db:"pos_entry_mode"`

	// Point of sale condition
	POSConditionCode string `json:"posConditionCode" db:"pos_condition_code"`

	// Acquiring bank identification number (BIN)
	AcquirerBIN string `json:"acquirerBIN" db:"acquirer_bin"`

	// Merchant (acceptor) id
	MerchantID string `json:"merchantId" db:"merchant_id"`

	// Merchant category code describes the merchants transaction category - e.g. restaurants
	MerchantCategoryCode string `json:"merchantCategoryCode" db:"merchant_category_code"`

	// Merchant (acceptor) terminal
	MerchantTerminal string `json:"merchantTerminal" db:"merchant_terminal"`

	// Merchant (acceptor) name
	MerchantName string `json:"merchantName" db:"merchant_name"`

	// Merchant (acceptor) address
	MerchantStreetAddress string `json:"merchantStreetAddress" db:"merchant_street_address"`

	// Merchant (acceptor) city
	MerchantCity string `json:"merchantCity" db:"merchant_city"`

	// Merchant (acceptor) state
	MerchantState string `json:"merchantState" db:"merchant_state"`

	// Merchant (acceptor) country
	MerchantCountry string `json:"merchantCountry" db:"merchant_country"`
}

type BusinessCardTransactionMini struct {
	ID *string `json:"id,omitempty" db:"id"`

	// Authorization amount
	AuthAmount *num.Decimal `json:"authAmount,omitempty" db:"auth_amount"`

	// Card transaction code
	TransactionType *CardTransactionType `json:"transactionType,omitempty" db:"transaction_type"`

	// Local currency amount
	LocalAmount *num.Decimal `json:"localAmount,omitempty" db:"local_amount"`

	// Local currency
	LocalCurrency *string `json:"localCurrency,omitempty" db:"local_currency"`

	// Local date of card transaction
	LocalDate *time.Time `json:"localDate,omitempty" db:"local_date"`

	// Billing currency
	BillingCurrency *string `json:"billingCurrency,omitempty" db:"billing_currency"`

	// Merchant category code describes the merchants transaction category - e.g. restaurants
	MerchantCategoryCode *string `json:"merchantCategoryCode,omitempty" db:"merchant_category_code"`

	// Merchant (acceptor) name
	MerchantName *string `json:"merchantName,omitempty" db:"merchant_name"`

	// Merchant (acceptor) address
	MerchantStreetAddress *string `json:"merchantStreetAddress,omitempty" db:"merchant_street_address"`

	// Merchant (acceptor) city
	MerchantCity *string `json:"merchantCity,omitempty" db:"merchant_city"`

	// Merchant (acceptor) state
	MerchantState *string `json:"merchantState,omitempty" db:"merchant_state"`

	// Merchant (acceptor) country
	MerchantCountry *string `json:"merchantCountry,omitempty" db:"merchant_country"`
}

// BusinessHoldTransaction ..
type BusinessHoldTransaction struct {
	ID string ` json:"id" db:"id"`

	TransactionID *string `json:"transactionId" db:"transaction_id"`

	// Hold number
	Number string `json:"number" db:"hold_number"`

	// Transaction amount
	Amount    num.Decimal `json:"amount" db:"amount"`
	AmountDep float64     `json:"amountDep" db:"amount_dep"`

	// Date of hold
	Date time.Time `json:"date" db:"transaction_date"`

	// Date of hold expiry
	ExpiryDate time.Time `json:"expiryDate" db:"expiry_date"`

	// Created
	Created time.Time `json:"created" db:"created"`

	// Modified
	Modified *time.Time `json:"modified" db:"modified"`
}

type BusinessHoldTransactionMini struct {
	ID *string ` json:"id,omitempty" db:"id"`

	// Hold number
	Number *string `json:"number,omitempty" db:"hold_number"`

	// Transaction amount
	Amount *num.Decimal `json:"amount,omitempty" db:"amount"`

	// Date of hold
	Date *time.Time `json:"date,omitempty" db:"transaction_date"`

	// Date of hold expiry
	ExpiryDate *time.Time `json:"expiryDate,omitempty" db:"expiry_date"`
}

type BusinessHoldTransactionCreate struct {
	// Related transaction
	TransactionID string `json:"transactionId,omitempty" db:"transaction_id"`

	// Hold number
	Number string `json:"number" db:"hold_number"`

	// Transaction amount
	Amount num.Decimal `json:"amount" db:"amount"`

	// Date of hold
	Date time.Time `json:"date" db:"transaction_date"`

	// Date of hold expiry
	ExpiryDate time.Time `json:"expiryDate" db:"expiry_date"`
}

type BusinessPostedTransaction struct {
	//Transaction ID
	ID shared.PostedTransactionID `json:"id" db:"id"`

	// Business Id
	BusinessID shared.BusinessID `json:"businessId" db:"business_id"`

	// Bank name - e.g. 'bbva'
	BankName string `json:"bankName" db:"bank_name"`

	// Bank account id
	BankTransactionID *string `json:"bankTransactionId" db:"bank_transaction_id"`

	// Partner bank account extra data
	BankExtra types.JSONText `json:"-" db:"bank_extra"`

	// Transaction type
	TransactionType TransactionType `json:"transactionType" db:"transaction_type"`

	// Bank account id if the trasanction is of type account
	AccountID *string `json:"accountId" db:"account_id"`

	// Card id
	CardID *string `json:"cardId" db:"card_id"`

	// Transaction code
	CodeType TransactionCodeType `json:"codeType" db:"code_type"`

	// Amount
	Amount    num.Decimal `json:"amount" db:"amount"`
	AmountDep float64     `json:"amountDep" db:"amount_dep"`

	// Currency
	Currency Currency `json:"currency" db:"currency"`

	// Card transaction details
	CardTransaction *BusinessCardTransactionMini `json:"cardTransaction,omitempty" db:"business_card_transaction"`

	// Card hold details
	HoldTransaction *BusinessHoldTransactionMini `json:"holdTransaction,omitempty" db:"business_hold_transaction"`

	// Money transfer id
	MoneyTransferID *string `json:"moneyTransferId" db:"money_transfer_id"`

	// Money request id
	MoneyRequestID *shared.PaymentRequestID `json:"moneyRequestId" db:"money_request_id"`

	// Notes
	SourceNotes *string `json:"sourceNotes" db:"source_notes"`
	Notes       *string `json:"notes" db:"notes"`

	// Contact id
	ContactID *string `json:"contactId" db:"contact_id"`

	// Origin account
	OriginAccount *string `json:"originAccount"`

	// Destination account
	DestinationAccount *string `json:"destinationAccount"`

	// Money transfer description
	BankTransactionDesc *string `json:"bankTransactionDesc" db:"bank_transaction_desc"`
	MoneyTransferDesc   *string `json:"moneyTransferDesc" db:"money_transfer_desc"`

	// Readable transaction description
	TransactionDesc string `json:"transactionDesc" db:"transaction_desc"`

	PartnerName string `json:"partnerName"`

	// Contact Detail
	Contact *Contact `json:"contact,omitempty" db:"business_contact"`

	// Dispute Detail
	Dispute *DisputeMini `json:"dispute" db:"business_transaction_dispute"`

	// MoneyTransfer Detail
	MoneyTransfer *MoneyTransfer `json:"moneyTransfer" db:"business_money_transfer"`

	// Transaction Date Created
	TransactionDate time.Time `json:"transactionDate" db:"transaction_date"`

	// Invoice Detail
	Invoice *payment.InvoiceMini `json:"invoice" db:"business_invoice"`

	// Receipt Detail
	Receipt *payment.ReceiptMini `json:"receipt" db:"business_receipt"`

	// Money Request Details
	MoneyRequest *payment.RequestMini `json:"-" db:"business_money_request"`

	// Payment Detail
	Payment *payment.PaymentMini `json:"-" db:"business_money_request_payment"`

	// Payment Summary
	PaymentSummary *PaymentSummary `json:"paymentSummary"`

	Source *TransactionSource `json:"source"`

	Destination *TransactionDestination `json:"destination"`

	AttachmentID *string `json:"attachmentId" db:"business_transaction_attachment.id"`

	AttachmentDeleted *time.Time `json:"-" db:"business_transaction_attachment.deleted"`

	// Transaction Notes
	TransactionNotes *string `json:"transactionNotes" db:"business_transaction_annotation.transaction_notes"`

	// Transaction title
	TransactionTitle *string `json:"transactionTitle" db:"transaction_title"`

	// Transaction sub type
	TransactionSubtype *TransactionSubtype `json:"transactionSubtype" db:"transaction_subtype"`

	ShopifyPayout *ShopifyPayout `json:"shopifyPayout,omitempty"`

	// Pass through to GRPC
	Status       TransactionStatus `json:"transactionStatus" db:"-"`
	Counterparty string            `json:"counterparty" db:"-"`

	NotificationID *string `json:"-" db:"notification_id"`

	//Created date
	Created time.Time `json:"created" db:"created"`

	// Modified date
	Modified *time.Time `json:"modified" db:"modified"`
}

type PaymentSummary struct {
	InvoiceAmount float64     `json:"invoiceAmount"`
	FeeAmount     num.Decimal `json:"feeAmount"`
	NetAmount     num.Decimal `json:"netAmount"`
}

type ShopifyPayout struct {
	PayoutID                  int64                           `json:"payoutId"`
	PayoutStatus              grpcShopify.ShopifyPayoutStatus `json:"payoutStatus"`
	PayoutDate                time.Time                       `json:"payoutDate"`
	Currency                  string                          `json:"currency"`
	Amount                    num.Decimal                     `json:"amount"`
	AdjustmentsFeeAmount      num.Decimal                     `json:"adjustmentsFeeAmount"`
	AdjustmentsGrossAmount    num.Decimal                     `json:"adjustmentsGrossAmount"`
	ChargesFeeAmount          num.Decimal                     `json:"chargesFeeAmount"`
	ChargesGrossAmount        num.Decimal                     `json:"chargesGrossAmount"`
	RefundsFeeAmount          num.Decimal                     `json:"refundsFeeAmount"`
	RefundsGrossAmount        num.Decimal                     `json:"refundsGrossAmount"`
	ReservedFundsFeeAmount    num.Decimal                     `json:"reservedFundsFeeAmount"`
	ReservedFundsGrossAmount  num.Decimal                     `json:"reservedFundsGrossAmount"`
	RetriedPayoutsFeeAmount   num.Decimal                     `json:"retriedPayoutsFeeAmount"`
	RetriedPayoutsGrossAmount num.Decimal                     `json:"retriedPayoutsGrossAmount"`
}

type BusinessPendingTransaction struct {
	//Transaction ID
	ID shared.PendingTransactionID `json:"id" db:"id"`

	// Business Id
	BusinessID shared.BusinessID `json:"businessId" db:"business_id"`

	// Bank name - e.g. 'bbva'
	BankName string `json:"bankName" db:"bank_name"`

	// Bank account id
	BankTransactionID *string `json:"bankTransactionId" db:"bank_transaction_id"`

	// Partner bank account extra data
	BankExtra types.JSONText `json:"-" db:"bank_extra"`

	// Transaction type
	TransactionType TransactionType `json:"transactionType" db:"transaction_type"`

	// Bank account id if the trasanction is of type account
	AccountID *string `json:"accountId" db:"account_id"`

	// Card id
	CardID *string `json:"cardId" db:"card_id"`

	// Transaction code
	CodeType TransactionCodeType `json:"codeType" db:"code_type"`

	// Amount
	Amount    num.Decimal `json:"amount" db:"amount"`
	AmountDep float64     `json:"amountDep" db:"amount_dep"`

	// Currency
	Currency Currency `json:"currency" db:"currency"`

	// Card transaction details
	CardTransaction *BusinessCardTransactionMini `json:"cardTransaction,omitempty" db:"business_card_pending_transaction"`

	// Card hold details
	HoldTransaction *BusinessHoldTransactionMini `json:"holdTransaction,omitempty" db:"business_hold_pending_transaction"`

	// Money transfer id
	MoneyTransferID *string `json:"moneyTransferId" db:"money_transfer_id"`

	// Money request id
	MoneyRequestID *shared.PaymentRequestID `json:"moneyRequestId" db:"money_request_id"`

	// Notes
	SourceNotes *string `json:"sourceNotes" db:"source_notes"`
	Notes       *string `json:"notes" db:"notes"`

	// Contact id
	ContactID *string `json:"contactId" db:"contact_id"`

	// Origin account
	OriginAccount *string `json:"originAccount"`

	// Destination account
	DestinationAccount *string `json:"destinationAccount"`

	// Money transfer description
	BankTransactionDesc *string `json:"bankTransactionDesc" db:"bank_transaction_desc"`
	MoneyTransferDesc   *string `json:"moneyTransferDesc" db:"money_transfer_desc"`

	// Readable transaction description
	TransactionDesc string `json:"transactionDesc" db:"transaction_desc"`

	// Contact Detail
	Contact *Contact `json:"contact,omitempty" db:"business_contact"`

	// Dispute Detail
	Dispute *DisputeMini `json:"dispute" db:"business_transaction_dispute"`

	// MoneyTransfer Detail
	MoneyTransfer *MoneyTransfer `json:"moneyTransfer" db:"business_money_transfer"`

	// Transaction Date Created
	TransactionDate time.Time `json:"transactionDate" db:"transaction_date"`

	TransactionStatus *string `json:"transactionStatus" db:"transaction_status"`

	PartnerName string `json:"partnerName" db:"partner_name"`

	Payment *payment.PaymentMini `json:"-" db:"business_money_request_payment"`

	Source *TransactionSource `json:"source"`

	Destination *TransactionDestination `json:"destination"`

	// Transaction title
	TransactionTitle *string `json:"transactionTitle" db:"transaction_title"`

	// Transaction sub type
	TransactionSubtype *TransactionSubtype `json:"transactionSubtype" db:"transaction_subtype"`

	// Pass through to GRPC
	Status       TransactionStatus `json:"transactionStatus" db:"-"`
	Counterparty string            `json:"counterparty" db:"-"`

	NotificationID *string `json:"-" db:"notification_id"`

	//Created date
	Created *time.Time `json:"created" db:"created"`

	// Modified date
	Modified *time.Time `json:"modified" db:"modified"`
}

type BusinessPostedTransactionCreate struct {
	ID shared.PostedTransactionID `json:"id" db:"id"`

	// Business Id
	BusinessID shared.BusinessID `json:"businessId" db:"business_id"`

	// Bank name - e.g. 'bbva'
	BankName string `json:"bankName" db:"bank_name"`

	// Bank account id
	BankTransactionID *string `json:"bankTransactionId" db:"bank_transaction_id"`

	// Partner bank account extra data
	BankExtra types.JSONText `json:"-" db:"bank_extra"`

	// Transaction type
	TransactionType TransactionType `json:"transactionType" db:"transaction_type"`

	// Bank account id if the trasanction is of type account
	AccountID *string `json:"accountId" db:"account_id"`

	// Card id
	CardID *string `json:"cardId" db:"card_id"`

	// Transaction code
	CodeType TransactionCodeType `json:"codeType" db:"code_type"`

	// Amount
	Amount num.Decimal `json:"amount" db:"amount"`

	// Currency
	Currency Currency `json:"currency" db:"currency"`

	// Money transfer id
	MoneyTransferID *string `json:"moneyTransferId" db:"money_transfer_id"`

	// Money request id
	MoneyRequestID *shared.PaymentRequestID `json:"moneyRequestId" db:"money_request_id"`

	// Notes
	SourceNotes *string `json:"sourceNotes" db:"source_notes"`

	// Contact id
	ContactID *string `json:"contactId" db:"contact_id"`

	// Money transfer description
	BankTransactionDesc *string `json:"bankTransactionDesc" db:"bank_transaction_desc"`

	// Readable transaction description
	TransactionDesc string `json:"transactionDesc" db:"transaction_desc"`

	// Transaction title
	TransactionTitle string `json:"transactionTitle" db:"transaction_title"`

	// Transaction sub type
	TransactionSubtype TransactionSubtype `json:"transactionSubtype" db:"transaction_subtype"`

	// Pass through to GRPC
	Status       TransactionStatus `json:"transactionStatus" db:"-"`
	Counterparty string            `json:"counterparty" db:"-"`

	// Transaction Date Created
	TransactionDate time.Time `json:"transactionDate" db:"transaction_date"`

	NotificationID *string `json:"-" db:"notification_id"`
}

type BusinessPendingTransactionCreate struct {
	// Transaction ID
	ID shared.PendingTransactionID `json:"id" db:"id"`

	// Business Id
	BusinessID shared.BusinessID `json:"businessId" db:"business_id"`

	// Bank name - e.g. 'bbva'
	BankName string `json:"bankName" db:"bank_name"`

	// Bank account id
	BankTransactionID *string `json:"bankTransactionId" db:"bank_transaction_id"`

	// Partner bank account extra data
	BankExtra types.JSONText `json:"-" db:"bank_extra"`

	// Transaction type
	TransactionType TransactionType `json:"transactionType" db:"transaction_type"`

	// Bank account id if the trasanction is of type account
	AccountID *string `json:"accountId" db:"account_id"`

	// Card id
	CardID *string `json:"cardId" db:"card_id"`

	// Transaction code
	CodeType TransactionCodeType `json:"codeType" db:"code_type"`

	// Amount
	Amount num.Decimal `json:"amount" db:"amount"`

	// Currency
	Currency Currency `json:"currency" db:"currency"`

	// Money transfer id
	MoneyTransferID *string `json:"moneyTransferId" db:"money_transfer_id"`

	// Money request id
	MoneyRequestID *shared.PaymentRequestID `json:"moneyRequestId" db:"money_request_id"`

	// Notes
	SourceNotes *string `json:"sourceNotes" db:"source_notes"`

	// Contact id
	ContactID *string `json:"contactId" db:"contact_id"`

	// Money transfer description
	BankTransactionDesc *string `json:"bankTransactionDesc" db:"bank_transaction_desc"`

	// Readable transaction description
	TransactionDesc string `json:"transactionDesc" db:"transaction_desc"`

	// Transaction Date Created
	TransactionDate time.Time `json:"transactionDate" db:"transaction_date"`

	TransactionStatus *string `json:"transactionStatus" db:"transaction_status"`

	PartnerName string `json:"partnerName" db:"partner_name"`

	// Transaction title
	TransactionTitle string `json:"transactionTitle" db:"transaction_title"`

	// Transaction sub type
	TransactionSubtype TransactionSubtype `json:"transactionSubtype" db:"transaction_subtype"`

	// Pass through to GRPC
	Status       TransactionStatus `json:"transactionStatus" db:"-"`
	Counterparty string            `json:"counterparty" db:"-"`

	NotificationID *string `json:"-" db:"notification_id"`
}

type BusinessPendingTransactionUpdate struct {
	BusinessID        shared.BusinessID
	MoneyTransferID   string
	BankTransactionID string
	TransactionDate   time.Time
	Status            string
}

type CSVTransaction struct {
	// Data in CSV format
	Data string `json:"data"`
}

type BusinessPostedTransactionAnnotation struct {
	ID string `json:"id" db:"id"`

	TransactionID shared.PostedTransactionID `json:"transactionId" db:"transaction_id"`

	TransactionNotes *string `json:"transactionNotes" db:"transaction_notes"`

	//Created date
	Created *time.Time `json:"created" db:"created"`

	// Modified date
	Modified *time.Time `json:"modified" db:"modified"`
}

type BusinessPostedTransactionUpdate struct {
	BusinessID shared.BusinessID `json:"businessId" db:"business_id"`

	TransactionID shared.PostedTransactionID `json:"transactionId" db:"transaction_id"`

	TransactionNotes *string `json:"transactionNotes" db:"transaction_notes"`
}

type TransactionSource struct {
	AccountNumber     *string           `json:"accountNumber"`
	AccountHolderName *string           `json:"accountHolderName"`
	BankName          *string           `json:"bankName"`
	CardLast4         *string           `json:"cardLast4"`
	CardHolderName    *string           `json:"cardHolderName"`
	CardBrand         *string           `json:"cardBrand"`
	PurchaseAddress   *services.Address `json:"purchaseAddress"`
	WalletType        *string           `json:"walletType"`
}

type TransactionDestination struct {
	AccountNumber     *string `json:"accountNumber"`
	AccountHolderName *string `json:"accountHolderName"`
	BankName          *string `json:"bankName"`
}

func BusinessPostedTransactionFromProto(gtxn *grpcBankTxn.Transaction, bus *bsrv.Business) (*BusinessPostedTransaction, error) {
	t := &BusinessPostedTransaction{}

	// Bank Transaction ID
	btxID, _ := id.ParseBankTransactionID(gtxn.Id)
	t.ID = shared.PostedTransactionID(btxID.UUIDString())

	// Business ID
	busID, _ := id.ParseBusinessID(gtxn.BusinessId)
	t.BusinessID = shared.BusinessID(busID.UUIDString())

	// Account ID
	accID, _ := id.ParseBankAccountID(gtxn.AccountId)
	if !accID.IsZero() {
		strID := accID.UUIDString()
		t.AccountID = &strID
	}

	// Debit Card ID
	dbcID, _ := id.ParseDebitCardID(gtxn.DebitCardId)
	if !dbcID.IsZero() {
		strID := dbcID.UUIDString()
		t.CardID = &strID
	}

	// Money tranfer ID
	btID, _ := id.ParseBankTransferID(gtxn.BankTransferId)
	if !btID.IsZero() {
		strID := btID.UUIDString()
		t.MoneyTransferID = &strID
	}

	// Payment Request ID
	prID, _ := id.ParsePaymentRequestID(gtxn.PaymentRequestId)
	if !prID.IsZero() {
		payReqID := shared.PaymentRequestID(prID.UUIDString())
		t.MoneyRequestID = &payReqID
	}

	// Contact ID
	contactID, _ := id.ParseContactID(gtxn.ContactId)
	if !contactID.IsZero() {
		strID := contactID.UUIDString()
		t.ContactID = &strID
	}

	var err error
	t.Amount, err = num.ParseDecimal(gtxn.Amount)
	if err != nil {
		return t, err
	}

	t.Currency = Currency(gtxn.Currency)

	// Status
	t.Status, _ = TransactionStatusFromProto[gtxn.Status]

	// Category/Type
	t.TransactionType, _ = TransactionTypeFromCategoryProto[gtxn.Category]

	// Type/Subtype
	subtype, ok := TransactionSubtypeFromTypeProto[gtxn.Type]
	if gtxn.CounterpartyType == grpcTxn.BankTransactionCounterpartyType_BTCT_SHOPPING_SHOPIFY_PAYOUT {
		subtype = TransactionSubtypeACHTransferShopifyCredit
	} else if gtxn.CounterpartyType == grpcTxn.BankTransactionCounterpartyType_BTCT_SHOPPING_SHOPIFY_REFUND {
		subtype = TransactionSubtypeACHTransferShopifyDebit
	} else if !ok {
		subtype = TransactionSubtype(gtxn.LegacySubtype)
	}

	t.TransactionSubtype = &subtype

	switch gtxn.Status {
	case grpcTxn.BankTransactionStatus_BTS_CARD_POSTED, grpcTxn.BankTransactionStatus_BTS_NONCARD_POSTED:
		if t.Amount.IsNegative() {
			t.CodeType = TransactionCodeTypeDebitPosted
		} else {
			t.CodeType = TransactionCodeTypeCreditPosted
		}
	}

	t.BankName = string(partnerbank.ProviderNameBBVA)
	t.BankTransactionID = &gtxn.PartnerTransactionId
	t.BankTransactionDesc = &gtxn.PartnerTransactionDesc
	t.MoneyTransferDesc = &gtxn.PartnerTransactionDesc

	// Transaction Date
	t.TransactionDate, _ = grpcTypes.Timestamp(gtxn.TransactionDate)

	t.Counterparty = gtxn.Counterparty

	// Title
	title := BankTransactionDisplayTitle(gtxn, bus)
	t.TransactionTitle = &title
	t.TransactionDesc = gtxn.LegacyDescription

	t.Notes = &gtxn.LegacyNotes
	t.SourceNotes = &gtxn.LegacyNotes

	t.Created, _ = grpcTypes.Timestamp(gtxn.Created)
	modified, _ := grpcTypes.Timestamp(gtxn.Modified)
	t.Modified = &modified

	return t, nil
}

func BusinessPostedTransactionFromFullProto(gtxn *grpcBankTxn.FullTransaction, bus *bsrv.Business) (*BusinessPostedTransaction, error) {
	t, err := BusinessPostedTransactionFromProto(gtxn.Transaction, bus)
	if err != nil {
		return t, err
	}

	// Card Transaction
	if gtxn.CardTransaction != nil {
		cardTxn := &BusinessCardTransactionMini{
			LocalCurrency:         &gtxn.CardTransaction.LocalCurrency,
			BillingCurrency:       &gtxn.CardTransaction.BillingCurrency,
			MerchantCategoryCode:  &gtxn.CardTransaction.MerchantCategoryCode,
			MerchantName:          &gtxn.CardTransaction.MerchantName,
			MerchantStreetAddress: &gtxn.CardTransaction.AcceptorStreetAddress,
			MerchantCity:          &gtxn.CardTransaction.AcceptorCity,
			MerchantState:         &gtxn.CardTransaction.AcceptorState,
			MerchantCountry:       &gtxn.CardTransaction.AcceptorCountry,
		}

		amount, err := num.ParseDecimal(gtxn.CardTransaction.AuthAmount)
		if err != nil {
			return t, err
		}

		cardTxn.AuthAmount = &amount

		cardTxnType := CardTransactionType(gtxn.CardTransaction.CardTransactionType)
		cardTxn.TransactionType = &cardTxnType

		amount, err = num.ParseDecimal(gtxn.CardTransaction.LocalAmount)
		if err != nil {
			return t, err
		}

		cardTxn.LocalAmount = &amount

		date, _ := grpcTypes.Timestamp(gtxn.CardTransaction.LocalDate)
		cardTxn.LocalDate = &date
		t.CardTransaction = cardTxn
	}

	if gtxn.HoldTransaction != nil {
		holdTxn := &BusinessHoldTransactionMini{
			Number: &gtxn.HoldTransaction.HoldNumber,
		}

		amount, err := num.ParseDecimal(gtxn.HoldTransaction.Amount)
		if err != nil {
			return t, err
		}

		holdTxn.Amount = &amount

		date, _ := grpcTypes.Timestamp(gtxn.HoldTransaction.TransactionDate)
		holdTxn.Date = &date

		date, _ = grpcTypes.Timestamp(gtxn.HoldTransaction.ExpiryDate)
		holdTxn.ExpiryDate = &date

		t.HoldTransaction = holdTxn
	}

	return t, nil
}

func BusinessPendingTransactionFromProto(gtxn *grpcBankTxn.Transaction, bus *bsrv.Business) (*BusinessPendingTransaction, error) {
	t := &BusinessPendingTransaction{}

	// Bank Transaction ID
	btxID, _ := id.ParseBankTransactionID(gtxn.Id)
	t.ID = shared.PendingTransactionID(btxID.UUIDString())

	// Business ID
	busID, _ := id.ParseBusinessID(gtxn.BusinessId)
	t.BusinessID = shared.BusinessID(busID.UUIDString())

	// Account ID
	accID, _ := id.ParseBankAccountID(gtxn.AccountId)
	if !accID.IsZero() {
		strID := accID.UUIDString()
		t.AccountID = &strID
	}

	// Debit Card ID
	dbcID, _ := id.ParseDebitCardID(gtxn.DebitCardId)
	if !dbcID.IsZero() {
		strID := dbcID.UUIDString()
		t.CardID = &strID
	}

	// Money tranfer ID
	btID, _ := id.ParseBankTransferID(gtxn.BankTransferId)
	if !btID.IsZero() {
		strID := btID.UUIDString()
		t.MoneyTransferID = &strID
	}

	// Payment Request ID
	prID, _ := id.ParsePaymentRequestID(gtxn.PaymentRequestId)
	if !prID.IsZero() {
		payReqID := shared.PaymentRequestID(prID.UUIDString())
		t.MoneyRequestID = &payReqID
	}

	// Contact ID
	contactID, _ := id.ParseContactID(gtxn.ContactId)
	if !contactID.IsZero() {
		strID := contactID.UUIDString()
		t.ContactID = &strID
	}

	var err error
	t.Amount, err = num.ParseDecimal(gtxn.Amount)
	if err != nil {
		return t, err
	}

	t.Currency = Currency(gtxn.Currency)

	// Status
	t.Status, _ = TransactionStatusFromProto[gtxn.Status]

	// Category/Type
	t.TransactionType, _ = TransactionTypeFromCategoryProto[gtxn.Category]

	// Type/Subtype
	subtype, ok := TransactionSubtypeFromTypeProto[gtxn.Type]
	if gtxn.CounterpartyType == grpcTxn.BankTransactionCounterpartyType_BTCT_SHOPPING_SHOPIFY_PAYOUT {
		subtype = TransactionSubtypeACHTransferShopifyCredit
	} else if gtxn.CounterpartyType == grpcTxn.BankTransactionCounterpartyType_BTCT_SHOPPING_SHOPIFY_REFUND {
		subtype = TransactionSubtypeACHTransferShopifyDebit
	} else if !ok {
		subtype = TransactionSubtype(gtxn.LegacySubtype)
	}

	t.TransactionSubtype = &subtype

	switch gtxn.Status {
	case grpcTxn.BankTransactionStatus_BTS_CARD_AUTHORIZED:
		t.CodeType = TransactionCodeTypeAuthApproved
	case grpcTxn.BankTransactionStatus_BTS_CARD_AUTH_REVERSED:
		t.CodeType = TransactionCodeTypeAuthReversed
	case grpcTxn.BankTransactionStatus_BTS_HOLD_SET:
		t.CodeType = TransactionCodeTypeHoldApproved
	case grpcTxn.BankTransactionStatus_BTS_HOLD_RELEASED:
		t.CodeType = TransactionCodeTypeHoldReleased
	case grpcTxn.BankTransactionStatus_BTS_VALIDATION:
		t.CodeType = TransactionCodeTypeValidation
	case grpcTxn.BankTransactionStatus_BTS_REVIEW:
		t.CodeType = TransactionCodeTypeReview
	case grpcTxn.BankTransactionStatus_BTS_CARD_AUTH_DECLINED:
		t.CodeType = TransactionCodeTypeAuthDeclined
	default:
		if t.Amount.IsNegative() {
			t.CodeType = TransactionCodeTypeDebitInProcess
		} else {
			t.CodeType = TransactionCodeTypeCreditInProcess
		}
	}

	t.BankName = string(partnerbank.ProviderNameBBVA)
	t.BankTransactionID = &gtxn.PartnerTransactionId
	t.BankTransactionDesc = &gtxn.PartnerTransactionDesc
	t.MoneyTransferDesc = &gtxn.PartnerTransactionDesc

	// Transaction Date
	t.TransactionDate, _ = grpcTypes.Timestamp(gtxn.TransactionDate)

	t.Counterparty = gtxn.Counterparty

	// Title
	title := BankTransactionDisplayTitle(gtxn, bus)
	t.TransactionTitle = &title
	t.TransactionDesc = gtxn.LegacyDescription

	t.Notes = &gtxn.LegacyNotes
	t.SourceNotes = &gtxn.LegacyNotes

	created, _ := grpcTypes.Timestamp(gtxn.Created)
	t.Created = &created

	modified, _ := grpcTypes.Timestamp(gtxn.Modified)
	t.Modified = &modified

	return t, nil
}

func BusinessPendingTransactionFromFullProto(gtxn *grpcBankTxn.FullTransaction, bus *bsrv.Business) (*BusinessPendingTransaction, error) {
	t, err := BusinessPendingTransactionFromProto(gtxn.Transaction, bus)
	if err != nil {
		return t, err
	}

	// Card Transaction
	if gtxn.CardTransaction != nil {
		cardTxn := &BusinessCardTransactionMini{
			LocalCurrency:         &gtxn.CardTransaction.LocalCurrency,
			BillingCurrency:       &gtxn.CardTransaction.BillingCurrency,
			MerchantCategoryCode:  &gtxn.CardTransaction.MerchantCategoryCode,
			MerchantName:          &gtxn.CardTransaction.MerchantName,
			MerchantStreetAddress: &gtxn.CardTransaction.AcceptorStreetAddress,
			MerchantCity:          &gtxn.CardTransaction.AcceptorCity,
			MerchantState:         &gtxn.CardTransaction.AcceptorState,
			MerchantCountry:       &gtxn.CardTransaction.AcceptorCountry,
		}

		amount, err := num.ParseDecimal(gtxn.CardTransaction.AuthAmount)
		if err != nil {
			return t, err
		}

		cardTxn.AuthAmount = &amount

		cardTxnType := CardTransactionType(gtxn.CardTransaction.CardTransactionType)
		cardTxn.TransactionType = &cardTxnType

		amount, err = num.ParseDecimal(gtxn.CardTransaction.LocalAmount)
		if err != nil {
			return t, err
		}

		cardTxn.LocalAmount = &amount

		date, _ := grpcTypes.Timestamp(gtxn.CardTransaction.LocalDate)
		cardTxn.LocalDate = &date
		t.CardTransaction = cardTxn
	}

	if gtxn.HoldTransaction != nil {
		holdTxn := &BusinessHoldTransactionMini{
			Number: &gtxn.HoldTransaction.HoldNumber,
		}

		amount, err := num.ParseDecimal(gtxn.HoldTransaction.Amount)
		if err != nil {
			return t, err
		}

		holdTxn.Amount = &amount

		date, _ := grpcTypes.Timestamp(gtxn.HoldTransaction.TransactionDate)
		holdTxn.Date = &date

		date, _ = grpcTypes.Timestamp(gtxn.HoldTransaction.ExpiryDate)
		holdTxn.ExpiryDate = &date

		t.HoldTransaction = holdTxn
	}

	return t, nil
}
