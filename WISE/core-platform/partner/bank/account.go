package bank

import (
	"time"

	"github.com/jmoiron/sqlx/types"
)

type AccountType string

func (t AccountType) String() string {
	return string(t)
}

const (
	AccountTypeChecking = AccountType("checking")
	AccountTypeSavings  = AccountType("savings")
)

type AccountRole string

func (r AccountRole) String() string {
	return string(r)
}

const (
	AccountRoleHolder     = AccountRole("holder")
	AccountRoleAuthorized = AccountRole("authorized")
	AccountRoleAttorney   = AccountRole("attorney")
	AccountRoleSpouse     = AccountRole("spouse")
	AccountRoleMinor      = AccountRole("minor")
	AccountRoleCustodian  = AccountRole("custodian")
)

type AccountParticipantRequest struct {
	ConsumerID ConsumerID  `json:"participantId"`
	Role       AccountRole `json:"role"`
}

type CreateConsumerBankAccountRequest struct {
	AccountHolderID   ConsumerID                  `json:"accountHolderId"`
	ExtraParticipants []AccountParticipantRequest `json:"extraParticipants"`
	AccountType       AccountType                 `json:"accountType"`
	Alias             *string                     `json:"alias"`
}

type CreateBusinessBankAccountRequest struct {
	BusinessID        BusinessID                  `json:"businessId"`
	ExtraParticipants []AccountParticipantRequest `json:"extraParticipants"`
	AccountType       AccountType                 `json:"accountType"`
	Alias             *string                     `json:"alias"`
	BusinessType      BusinessEntity              `json:"businessType"`
	IsForeign         bool                        `json:"isForeign"`
}

type AccountBankID string

func (id AccountBankID) String() string {
	return string(id)
}

type AccountStatus string

func (s AccountStatus) String() string {
	return string(s)
}

const (
	AccountStatusActive    = AccountStatus("active")
	AccountStatusBlocked   = AccountStatus("blocked")
	AccountStatusLocked    = AccountStatus("locked")
	AccountStatusClosed    = AccountStatus("closed")
	AccountStatusDormant   = AccountStatus("dormant")
	AccountStatusAbandoned = AccountStatus("abandoned")
	AccountStatusEscheated = AccountStatus("escheated")
	AccountStatusChargeOff = AccountStatus("chargeOff")
)

type Currency string

func (c Currency) String() string {
	return string(c)
}

const (
	CurrencyUSD = Currency("usd")
)

type AccountParticipantResponse struct {
	ConsumerID ConsumerID  `json:"participantId"`
	Role       AccountRole `json:"role"`
}

type CreateConsumerBankAccountResponse struct {
	AccountID        AccountBankID                `json:"bankAccountId"`
	BankName         ProviderName                 `json:"bankName"`
	BankExtra        types.JSONText               `json:"bankExtra"`
	AccountHolderID  ConsumerID                   `json:"accountHolderId"`
	Participants     []AccountParticipantResponse `json:"participants"`
	AccountType      AccountType                  `json:"accountType"`
	AccountNumber    string                       `json:"accountNumber"`
	RoutingNumber    string                       `json:"routingNumber"`
	WireRouting      string                       `json:"wireRouting"`
	Alias            string                       `json:"alias"`
	Status           AccountStatus                `json:"status"`
	Opened           time.Time                    `json:"opened"`
	AvailableBalance float64                      `json:"availableBalance"`
	PostedBalance    float64                      `json:"postedBalance"`
	Currency         Currency                     `json:"currency"`
}

type CreateBusinessBankAccountResponse struct {
	AccountID        AccountBankID                `json:"bankAccountId"`
	BankName         ProviderName                 `json:"bankName"`
	BankExtra        types.JSONText               `json:"bankExtra"`
	BusinessID       BusinessID                   `json:"businessId"`
	Participants     []AccountParticipantResponse `json:"participants"`
	AccountType      AccountType                  `json:"accountType"`
	AccountNumber    string                       `json:"accountNumber"`
	RoutingNumber    string                       `json:"routingNumber"`
	WireRouting      string                       `json:"wireRouting"`
	Alias            string                       `json:"alias"`
	Status           AccountStatus                `json:"status"`
	Opened           time.Time                    `json:"opened"`
	AvailableBalance float64                      `json:"availableBalance"`
	PostedBalance    float64                      `json:"postedBalance"`
	Currency         Currency                     `json:"currency"`
}

type PatchBankAccountRequest struct {
	Alias string `json:"alias"`
}

type AccountCloseReason string

func (r AccountCloseReason) String() string {
	return string(r)
}

const (
	AccountCloseReasonDeath              = AccountCloseReason("death")
	AccountCloseReasonBankError          = AccountCloseReason("error")
	AccountCloseReasonNonSufficientFunds = AccountCloseReason("nsf")
	AccountCloseReasonLegal              = AccountCloseReason("legal")
	AccountCloseReasonConsolidating      = AccountCloseReason("consolidating")
	AccountCloseReasonBank               = AccountCloseReason("bank")
	AccountCloseReasonFraud              = AccountCloseReason("fraud")
	AccountCloseReasonCustomer           = AccountCloseReason("customer")
)

type CloseBankAccountRequest struct {
	Reason AccountCloseReason `json:"reason"`
}

type GetBankAccountResponse struct {
	AccountID        AccountBankID `json:"bankAccountId"`
	AccountType      AccountType   `json:"accountType"`
	AccountNumber    string        `json:"accountNumber"`
	RoutingNumber    string        `json:"routingNumber"`
	WireRouting      string        `json:"wireRouting"`
	Alias            string        `json:"alias"`
	Status           AccountStatus `json:"status"`
	Opened           time.Time     `json:"opened"`
	AvailableBalance float64       `json:"availableBalance"`
	PostedBalance    float64       `json:"postedBalance"`
	Currency         Currency      `json:"currency"`
}

type AccountBlockType string

const (
	AccountBlockTypeDebit  = AccountBlockType("debits")
	AccountBlockTypeCredit = AccountBlockType("credits")
	AccountBlockTypeCheck  = AccountBlockType("checks")
	AccountBlockTypeAll    = AccountBlockType("all")
)

type AccountBlockRequest struct {
	AccountID AccountBankID    `json:"bankAccountId"`
	Type      AccountBlockType `json:"blockType"`
	Reason    string           `json:"reason"`
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

type AccountBlockResponse struct {
	BlockID      AccountBlockBankID `json:"BlockId"`
	Type         AccountBlockType   `json:"blockType"`
	Status       AccountBlockStatus `json:"status"`
	Reason       string             `json:"reason"`
	BlockDate    time.Time          `json:"blockDate"`
	CanceledDate *time.Time         `json:"canceledDate"`
}

type AccountUnblockRequest struct {
	AccountID AccountBankID      `json:"bankAccountId"`
	BlockID   AccountBlockBankID `json:"BlockId"`
}

type AccountStatementBankID string

type AccountStatementResponse struct {
	StatementID AccountStatementBankID `json:"statementId"`
	Description string                 `json:"description"`
	Created     time.Time              `json:"created"`
	PageCount   int                    `json:"pageCount"`
}

type GetAccountStatementDocument struct {
	ContentType string `json:"contentType"`
	Content     []byte `json:"content"`
}
