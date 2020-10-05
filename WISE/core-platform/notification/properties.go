package notification

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx/types"
	"github.com/wiseco/core-platform/notification/activity"
	"github.com/wiseco/core-platform/services"
	"github.com/wiseco/core-platform/services/banking"
	"github.com/wiseco/core-platform/services/payment"
	"github.com/wiseco/core-platform/shared"
	"github.com/wiseco/core-platform/transaction"
	"github.com/wiseco/go-lib/num"
)

// Type notification type
type Type string

// EntityType notification entity type
type EntityType string

const (
	EntityTypeConsumer = EntityType("consumer")
	EntityTypeBusiness = EntityType("business")
	EntityTypeMember   = EntityType("member")
)

// Attribute notification type attribute
type Attribute string

// Action notitication action e.g `create`
type Action string

// Status the status
type Status string

// TransactionType ..
type TransactionType string

const (
	// TypeConsumer notification of consumer
	TypeConsumer = Type("consumer")

	// TypeBusiness notification of business
	TypeBusiness = Type("business")

	// TypeAccount notification of account
	TypeAccount = Type("account")

	// TypeCard notification of card
	TypeCard = Type("card")

	// TypeMoneyTransfer notification of money transfer
	TypeMoneyTransfer = Type("moneyTransfer")

	// TypeTransaction notification of transaction
	TypeTransaction = Type("transaction")

	// TypePendingTransfer notification of pending transfer
	TypePendingTransfer = Type("pendingTransfer")
)

const (
	// ActionCreate notification action create
	ActionCreate = Action("create")

	// ActionOpen
	ActionOpen = Action("open")

	// ActionActivate
	ActionActivate = Action("activate")

	// ActionReissue
	ActionReissue = Action("reissue")

	// ActionUpdate notification action update
	ActionUpdate = Action("update")

	// ActionAdd
	ActionAdd = Action("add")

	// ActionRemove
	ActionRemove = Action("remove")

	// ActionDelete notification action delete
	ActionDelete = Action("delete")

	// ActionCancel
	ActionCancel = Action("cancel")

	// ActionCorrected action corrected
	ActionCorrected = Action("corrected")

	// ActionAuthorize notification action authorize
	ActionAuthorize = Action("authorize")

	// ActionHold notification action hold
	ActionHold = Action("hold")

	// ActionPosted notification action posted
	ActionPosted = Action("posted")
)

const (
	AttributeEmail = Attribute("email")

	AttributeKYC = Attribute("kyc")

	AttributePhone = Attribute("phone")

	AttributeAddress = Attribute("address")

	AttributeBlock = Attribute("block")
)

const (
	//TransactionTypeACH ACH
	TransactionTypeACH = TransactionType("ach")

	// TransactionTypeAdjustment adjustment
	TransactionTypeAdjustment = TransactionType("adjustment")

	// TransactionTypeATM ATM
	TransactionTypeATM = TransactionType("atm")

	// TransactionTypeCheck check
	TransactionTypeCheck = TransactionType("check")

	// TransactionTypeDeposit deposit
	TransactionTypeDeposit = TransactionType("deposit")

	// TransactionTypeFee fee
	TransactionTypeFee = TransactionType("fee")

	// TransactionTypeOtherCredit other credit
	TransactionTypeOtherCredit = TransactionType("otherCredit")

	// TransactionTypeOtherDebit other debit
	TransactionTypeOtherDebit = TransactionType("otherDebit")

	//TransactionTypePurchase Purchase transaction
	TransactionTypePurchase = TransactionType("purchase")

	//TransactionTypeRefund Refund
	TransactionTypeRefund = TransactionType("refund")

	// TransactionTypeReturn return
	TransactionTypeReturn = TransactionType("return")

	// TransactionTypeReversal reversal
	TransactionTypeReversal = TransactionType("reversal")

	// TransactionTypeTransfer transfer
	TransactionTypeTransfer = TransactionType("transfer")

	// TransactionTypeVisaCredit visa credit
	TransactionTypeVisaCredit = TransactionType("visaCredit")

	// TransactionTypeWithdrawal withdrawal
	TransactionTypeWithdrawal = TransactionType("withdrawal")

	//TransactionTypeOther Miscellaneous
	TransactionTypeOther = TransactionType("other")
)

type POSEntryMode string

func (m POSEntryMode) isOnlinePayment() bool {

	str := m.String()
	m = POSEntryMode(string(str[0:2]))

	switch m {
	case POSEntryModeKeyEntry, POSEntryModeManualECommerce, POSEntryModeStoredCheckout:
		return true
	default:
		return false
	}
}

func (m POSEntryMode) String() string {
	return string(m)
}

const (
	POSEntryModeKeyEntry        = POSEntryMode("06")
	POSEntryModeManualECommerce = POSEntryMode("81")
	POSEntryModeStoredCheckout  = POSEntryMode("96")
)

//Notification coming from SQS
type Notification struct {
	ID         string         `json:"id" `
	EntityID   string         `json:"entityId"`
	EntityType EntityType     `json:"entityType"`
	BankName   string         `json:"bankName"`
	SourceID   string         `json:"sourceId"`
	Type       Type           `json:"type" `
	Attribute  *Attribute     `json:"attribute"`
	Action     Action         `json:"action"`
	Version    string         `json:"version"`
	Created    time.Time      `json:"created"`
	Data       types.JSONText `json:"data"`
}

type AddressResponse AddressRequest

type AddressRequest struct {
	Type string `json:"type"`

	Line1 string `json:"line1"`

	Line2 string `json:"line2,omitempty"`

	City string `json:"city"`

	State string `json:"state"`

	ZipCode string `json:"zipCode"`

	Country string `json:"country,omitempty"`
}

//ConsumerNotification ..
type ConsumerNotification struct {
	BankEntityID       string            `json:"bankEntityId"`
	TaxID              string            `json:"taxId"`
	TaxIDType          string            `json:"taxIdType"`
	FirstName          string            `json:"firstName"`
	MiddleName         *string           `json:"middleName"`
	LastName           string            `json:"lastName"`
	Residency          ConsumerResidency `json:"residency"`
	CitizenshipCountry string            `json:"citizenshipCountry"`
	Document           string            `json:"document"`
}

//ConsumerResidency ..
type ConsumerResidency struct {
	Country string `json:"country"`
	Status  string `json:"status"`
}

//- Handles consumer email/phone number update
type ConsumerContactUpdateNotification struct {
	BankEntityID string           `json:"bankEntityId"`
	AttributeID  string           `json:"attributeId"`
	Phone        *string          `json:"phone,omitempty"`
	Email        *string          `json:"email,omitempty"`
	Address      *AddressResponse `json:"address,omitempty"`
}

type ConsumerKYCNotesNotification struct {
	Source string `json:"source"`
	Desc   string `json:"desc"`
}

type ConsumerKYCUpdateNotification struct {
	BankEntityId string                         `json:"bankId"`
	Risk         string                         `json:"risk"`
	Result       string                         `json:"result"`
	KYCStatus    string                         `json:"kycStatus"`
	KYCNotes     []ConsumerKYCNotesNotification `json:"kycNotes"`
}

//NonConsumerNotification  ..
type NonConsumerNotification struct {
	BankEntityID    string     `json:"bankEntityId"`
	TaxID           string     `json:"taxId"`
	TaxIDType       string     `json:"taxIdType"`
	LegalName       string     `json:"legalName"`
	EntityType      EntityType `json:"entityType"`
	IndustryType    string     `json:"industryType"`
	EntityFormation *string    `json:"entityFormation"`
}

//- Handles business email/phone number update
type BusinessContactUpdateNotification struct {
	BankEntityID string           `json:"bankEntityId"`
	AttributeID  string           `json:"attributeId"`
	Phone        *string          `json:"phone,omitempty"`
	Email        *string          `json:"email,omitempty"`
	Address      *AddressResponse `json:"address,omitempty"`
}

type BusinessKYCNotesNotification struct {
	Source string `json:"source"`
	Desc   string `json:"desc"`
}

type BusinessKYCNotification struct {
	BankEntityID string                         `json:"bankEntityId"`
	Risk         string                         `json:"risk"`
	Result       string                         `json:"result"`
	KYCStatus    string                         `json:"kycStatus"`
	KYCNotes     []BusinessKYCNotesNotification `json:"kycNotes"`
}

type NotificationCardReissueReason string

// New card creation
type NewCardNotification struct {
	BankCardID    string  `json:"bankCardId"`
	BankAccountID string  `json:"bankAccountId"`
	Reason        *string `json:"reason"`
}

// Block card
type CardBlockNotification struct {
	BankCardID    string `json:"cardId"`
	BankAccountID string `json:"bankAccountId"`
	Status        string `json:"status"`
	Reason        string `json:"reason"`
}

type CardStatusNotification struct {
	BankCardID    string     `json:"cardId"`
	BankAccountID string     `json:"bankAccountId"`
	Status        CardStatus `json:"status"`
}

type BankCardDetails struct {
	BusinessLegalName *string              `db:"legal_name"`
	DBA               services.StringArray `db:"dba"`
	ID                string               `db:"business_bank_card.id"`
	CardNumberMasked  string               `db:"card_number_masked"`
}

func (b BankCardDetails) GetCardNumberLastFour() (string, error) {
	if len(b.CardNumberMasked) < 4 {
		return "", fmt.Errorf("Unable to get last four digits of card, CardNumberMasked not long enough. length: %d", len(b.CardNumberMasked))
	}

	return b.CardNumberMasked[len(b.CardNumberMasked)-4:], nil
}

// Account specific notification
type AccountOpenedNotification struct {
	BankAccountID string     `json:"bankAccountId"`
	Opened        *time.Time `json:"accountType"`
}

type AccountBlockNotification struct {
	BankAccountID string `json:"bankAccountId"`
	Block         string `json:"block"`
	Reason        string `json:"reason"`
}

//AccountStatusNotification ..
type AccountStatusNotification struct {
	BankAccountID     string    `json:"bankAccountID"`
	Status            Status    `json:"status"`
	StatusDate        time.Time `json:"opened"`
	Description       *string   `json:"description"`
	BlockID           *string   `json:"blockId"`
	BlockType         *string   `json:"blockType"`
	BlockStatus       *string   `json:"blockStatus"`
	BlockReason       *string   `json:"blockReason"`
	ChargeOffDaysLeft *int      `json:"chargeOffDaysLeft"`
}

//AccountNotification ..
type AccountNotification struct {
	BankAccountID     string    `json:"bankAccountID"`
	AccountType       string    `json:"accountType"`
	AccountNumber     string    `json:"accountNumber"`
	RoutingNumber     string    `json:"routingNumber"`
	FullAccountNumber string    `json:"fullAccountNumber"`
	Opened            time.Time `json:"opened"`
	PendingBalance    float64   `json:"pendingBalance"`
	AvailableBalance  float64   `json:"availableBalance"`
	PostedBalance     float64   `json:"postedBalance"`
	Status            Status    `json:"status"`
}

//MoneyTransferStatusNotification ..
type MoneyTransferStatusNotification struct {
	MoneyTransferID   string    `json:"moneyTransferId"`
	OriginAccountID   string    `json:"originAccountId"`
	OriginAccountType string    `json:"originAccountType"`
	DestAccountID     string    `json:"destAccountId"`
	DestAccountType   string    `json:"destAccountType"`
	Amount            float64   `json:"amount"`
	Currency          string    `json:"currency,omitempty"`
	Status            string    `json:"status"`
	BankStatus        string    `json:"bank_status"`
	StatusUpdated     time.Time `json:"statusUpdated"`
	ReasonCode        *string   `json:"reasonCode"`
	ReasonDescription *string   `json:"reasonDescription"`
}

//MoneyTransferCorrectedNotification ..
type MoneyTransferCorrectedNotification struct {
	MoneyTransferID   string    `json:"moneyTransferId"`
	SourceAccountID   string    `json:"sourceAccountId"`
	DestAccountID     string    `json:"destAccountId"`
	Amount            float64   `json:"amount"`
	Currency          string    `json:"currency,omitempty"`
	Status            Status    `json:"status"`
	StatusUpdated     time.Time `json:"statusUpdated"`
	ChangeCode        string    `json:"changeCode"`
	ChangeDescription string    `json:"changeDescription"`
	ChangeAccountType string    `json:"changeAccountType"`
}

//CardTransactionNotification ..
type CardTransactionNotification struct {
	// Card specific transaction id
	CardTransactionID string `json:"cardTransactionId"`

	// Network used for card transaction
	TransactionNetwork string `json:"transactionNetwork"`

	// Authorization amount
	AuthAmount float64 `json:"authAmount"`

	// Authorization date
	AuthDate time.Time `json:"authDate"`

	// Authorization response code
	AuthResponseCode AuthResponseCode `json:"authResponseCode"`

	// Authorization number
	AuthNumber string `json:"authNumber"`

	// Card transaction code
	TransactionType transaction.CardTransactionType `json:"transactionType"`

	// Local currency amount
	LocalAmount float64 `json:"localAmount"`

	// Local currency
	LocalCurrency string `json:"localCurrency"`

	// Local date of card transaction
	LocalDate time.Time `json:"localDate"`

	// Billing Currency
	BillingCurrency string `json:"billingCurrency"`

	// Point of sale entry mode describes how the card was entered or used - e.g. swipe or chip
	POSEntryMode POSEntryMode `json:"posEntryMode"`

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

//HoldTransactionNotification ..
type HoldTransactionNotification struct {
	// Hold number
	Number string `json:"number"`

	// Transaction amount
	Amount float64 `json:"amount"`

	// Date of hold
	Date time.Time `json:"date"`

	// Date of hold expiry
	ExpiryDate time.Time `json:"expiryDate"`
}

//TransactionNotification ..
type TransactionNotification struct {
	// Bank name - e.g. 'bbva'
	BankName string `json:"bankName"`

	// Bank account id
	BankTransactionID string `json:"bankTransactionId"`

	// Partner bank account extra data
	BankExtra types.JSONText `json:"-"`

	// Transaction type
	TransactionType TransactionType `json:"type"`

	// Bank account id
	BankAccountID *string `json:"accountId"`

	// Card id
	BankCardID *string `json:"cardId"`

	// Transaction code
	CodeType transaction.TransactionCodeType `json:"codeType"`

	// Amount
	Amount num.Decimal `json:"amount"`

	// Currency
	Currency string `json:"currency"`

	// Card transaction details
	CardTransaction *CardTransactionNotification `json:"cardTransaction"`

	// Card hold details
	HoldTransaction *HoldTransactionNotification `json:"holdTransaction"`

	// Money transfer id
	BankMoneyTransferID *string `json:"bankMoneyTransferId"`

	// Money transfer description
	MoneyTransferDesc   *string `json:"moneyTransferDesc"`
	BankTransactionDesc *string `json:"bankTransactionDesc"`

	// Transaction Date Created
	TransactionDate time.Time `json:"transactionDate"`
}

type CardStatus string

const (
	CardStatusShipped      = CardStatus("shipped")
	CardStatusActivated    = CardStatus("activated")
	CardStatusCanceled     = CardStatus("canceled")
	CardStatusBlocked      = CardStatus("blocked")
	CardStatusUnblocked    = CardStatus("unblocked")
	CardStatusReissued     = CardStatus("reissued")
	CardStatusLimitChanged = CardStatus("limitChanged")
)

var bankCardStatusTo = map[CardStatus]banking.CardStatus{
	CardStatusShipped:      banking.CardStatusShipped,
	CardStatusActivated:    banking.CardStatusActive,
	CardStatusCanceled:     banking.CardStatusCanceled,
	CardStatusBlocked:      banking.CardStatusBlocked,
	CardStatusUnblocked:    banking.CardStatusUnblocked,
	CardStatusReissued:     banking.CardStatusReissued,
	CardStatusLimitChanged: banking.CardStatusLimitChanged,
}

type PendingTransferNotification struct {
	// Bank name - e.g. 'bbva'
	BankName string `json:"bankName"`

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
	PartnerName string `json:"partnerName"`

	// Currency
	Currency string `json:"currency"`

	// Money transfer id
	MoneyTransferID *string `json:"moneyTransferId"`

	// Request id
	MoneyRequestID *shared.PaymentRequestID `json:"moneyRequestId"`

	// Contact id
	ContactID *string `json:"contactId"`

	CodeType transaction.TransactionCodeType `json:"codeType"`

	// Transaction Date Created
	TransactionDate time.Time `json:"transactionDate"`
}

type TransactionMessage struct {
	SenderReceiverName     *string
	ActivityType           activity.Type
	NotificationHeader     string
	NotificationBody       string
	TransactionTitle       string
	TransactionDescription string
	ActivtiyID             *string
	BusinessName           string
	InterestEarnedMonth    *string
	MerchantName           *string
	Counterparty           string
}

type TransferDetails struct {
	MoneyTransferID     *string                  `db:"business_money_transfer.id"`
	TransferContactID   *string                  `db:"business_money_transfer.contact_id"`
	RequestID           *shared.PaymentRequestID `db:"business_money_transfer.money_request_id"`
	PaymentID           *string                  `db:"business_money_request_payment.id"`
	Notes               *string                  `db:"business_money_transfer.notes"`
	RequestContactID    *string                  `db:"business_money_request.contact_id"`
	RequestType         *string                  `db:"business_money_request.request_type"`
	RequestSource       *payment.RequestSource   `db:"business_money_request.request_source"`
	MonthlyInterestID   *string                  `db:"business_money_transfer.account_monthly_interest_id"`
	ContactFirstName    *string                  `db:"business_contact.first_name"`
	ContactLastName     *string                  `db:"business_contact.last_name"`
	ContactBusinessName *string                  `db:"business_contact.business_name"`
	ContactEmail        *string                  `db:"business_contact.email"`
	ContactPhone        *string                  `db:"business_contact.phone_number"`
	UserID              shared.UserID            `db:"business.owner_id"`
	BusinessLegalName   *string                  `db:"business.legal_name"`
	BusinessDBA         services.StringArray     `db:"business.dba"`
	BusinessPhone       string                   `db:"business.phone"`
	BusinessEmail       string                   `db:"business.email"`
	BankName            *string                  `db:"business_linked_bank_account.bank_name"`
	AccountNumber       *string                  `db:"business_linked_bank_account.account_number"`
	LinkedBankAccountID *string                  `db:"business_money_request_payment.linked_bank_account_id"`
	PaymentDate         *time.Time               `db:"business_money_request_payment.payment_date"`
	InvoiceID           *string                  `db:"business_invoice.id"`
	InvoiceNumber       *string                  `db:"business_invoice.invoice_number"`
}

const (
	PayoutStatusPaid = "paid"
)

type AuthResponseCode string

const (
	AuthResponseCodeInvalidMerchant         = "03"
	AuthResponseCodeRequestInProgress       = "09"
	AuthResponseCodeInvalidTransaction      = "12"
	AuthResponseCodeInvalidAmount           = "13"
	AuthResponseCodeInvalidCardNumber       = "14"
	AuthResponseCodeNoSuchIssuer            = "15"
	AuthResponseCodeInvalidResponse         = "20"
	AuthResponseCodeUnacceptableTxnFee      = "23"
	AuthResponseCodeFormatError             = "30"
	AuthResponseCodeBankNotSupported        = "31"
	AuthResponseCodeExpiredCardPickUp       = "33"
	AuthResponseCodeSuspectedFraudPickUp    = "34"
	AuthResponseCodeRestrictedCardPickUp    = "36"
	AuthResponseCodePINTriesExceededPickUp  = "38"
	AuthResponseCodeLostCard                = "41"
	AuthResponseCodeStolenCard              = "43"
	AuthResponseCodeInsufficientFunds       = "51"
	AuthResponseCodeExpiredCard             = "54"
	AuthResponseCodeIncorrectPIN            = "55"
	AuthResponseCodeSuspectedFraud          = "59"
	AuthResponseCodeAmountLimitExceeded     = "61"
	AuthResponseCodeRestrictedCard          = "62"
	AuthResponseCodeSecurityViolation       = "63"
	AuthResponseCodeFrequencyLimitExceeded  = "65"
	AuthResponseCodePINTriesExceeded        = "75"
	AuthResponseCodeKeySynchronizationError = "76"
	AuthResponseCodeCVVVerificationFailed   = "89"
	AuthResponseCodeLawViolation            = "93"
	AuthResponseCodeDuplicateTransaction    = "94"
	AuthResponseCodeSystemMalfunction       = "96"
)

const (
	DeclinedInvalidMerchant         = "invalid merchant"
	DeclinedRequestInProgress       = "request in progress"
	DeclinedInvalidTransaction      = "invalid transaction"
	DeclinedInvalidAmount           = "invalid amount"
	DeclinedInvalidCardNumber       = "invalid card number"
	DeclinedNoSuchIssuer            = "no such issuer"
	DeclinedInvalidResponse         = "invalid response"
	DeclinedUnacceptableTxnFee      = "unacceptable fee"
	DeclinedFormatError             = "format error"
	DeclinedBankNotSupported        = "bank not supported"
	DeclinedExpiredCard             = "expired card"
	DeclinedSuspectedFraud          = "suspected fraud"
	DeclinedRestrictedCard          = "restricted card"
	DeclinedPINTriesExceeded        = "exceeded PIN tries"
	DeclinedLostCard                = "lost card"
	DeclinedStolenCard              = "stolen card"
	DeclinedInsufficientFunds       = "insufficient funds"
	DeclinedIncorrectPIN            = "incorrect PIN"
	DeclinedAmountLimitExceeded     = "amount limit exceeded"
	DeclinedSecurityViolation       = "security violation"
	DeclinedFrequencyLimitExceeded  = "frequency limit exceeded"
	DeclinedKeySynchronizationError = "synchronization error"
	DeclinedCVVVerificationFailed   = "CVV verification failed"
	DeclinedLawViolation            = "violation of law"
	DeclinedDuplicateTransaction    = "duplicate transaction"
	DeclinedSystemMalfunction       = "system malfunction"
)

var AuthRespCodeToMessage = map[AuthResponseCode]string{
	AuthResponseCodeInvalidMerchant:         DeclinedInvalidMerchant,
	AuthResponseCodeRequestInProgress:       DeclinedRequestInProgress,
	AuthResponseCodeInvalidTransaction:      DeclinedInvalidTransaction,
	AuthResponseCodeInvalidAmount:           DeclinedInvalidAmount,
	AuthResponseCodeInvalidCardNumber:       DeclinedInvalidCardNumber,
	AuthResponseCodeNoSuchIssuer:            DeclinedNoSuchIssuer,
	AuthResponseCodeInvalidResponse:         DeclinedInvalidResponse,
	AuthResponseCodeUnacceptableTxnFee:      DeclinedUnacceptableTxnFee,
	AuthResponseCodeFormatError:             DeclinedFormatError,
	AuthResponseCodeBankNotSupported:        DeclinedBankNotSupported,
	AuthResponseCodeExpiredCardPickUp:       DeclinedExpiredCard,
	AuthResponseCodeSuspectedFraudPickUp:    DeclinedSuspectedFraud,
	AuthResponseCodeRestrictedCardPickUp:    DeclinedRestrictedCard,
	AuthResponseCodePINTriesExceededPickUp:  DeclinedPINTriesExceeded,
	AuthResponseCodeLostCard:                DeclinedLostCard,
	AuthResponseCodeStolenCard:              DeclinedStolenCard,
	AuthResponseCodeInsufficientFunds:       DeclinedInsufficientFunds,
	AuthResponseCodeExpiredCard:             DeclinedExpiredCard,
	AuthResponseCodeIncorrectPIN:            DeclinedIncorrectPIN,
	AuthResponseCodeSuspectedFraud:          DeclinedSuspectedFraud,
	AuthResponseCodeAmountLimitExceeded:     DeclinedAmountLimitExceeded,
	AuthResponseCodeRestrictedCard:          DeclinedRestrictedCard,
	AuthResponseCodeSecurityViolation:       DeclinedLawViolation,
	AuthResponseCodeFrequencyLimitExceeded:  DeclinedFrequencyLimitExceeded,
	AuthResponseCodePINTriesExceeded:        DeclinedPINTriesExceeded,
	AuthResponseCodeKeySynchronizationError: DeclinedKeySynchronizationError,
	AuthResponseCodeCVVVerificationFailed:   DeclinedCVVVerificationFailed,
	AuthResponseCodeLawViolation:            DeclinedLawViolation,
	AuthResponseCodeDuplicateTransaction:    DeclinedDuplicateTransaction,
	AuthResponseCodeSystemMalfunction:       DeclinedSystemMalfunction,
}
