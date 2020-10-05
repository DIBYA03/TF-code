package bbva

import (
	"errors"
	"strings"
	"time"

	partnerbank "github.com/wiseco/core-platform/partner/bank"
)

type Currency string

const (
	CurrencyUSD = Currency("USD")
)

func (c Currency) String() string {
	return string(c)
}

func (c Currency) Upper() string {
	return strings.ToUpper(c.String())
}

var partnerCurrencyTo = map[Currency]partnerbank.Currency{
	CurrencyUSD: partnerbank.CurrencyUSD,
}

var partnerCurrencyFrom = map[partnerbank.Currency]Currency{
	partnerbank.CurrencyUSD: CurrencyUSD,
}

type AccountType string

const (
	AccountTypeChecking = AccountType("checking")
)

var partnerAccountTypeFrom = map[partnerbank.AccountType]AccountType{
	partnerbank.AccountTypeChecking: AccountTypeChecking,
}

var partnerAccountTypeTo = map[AccountType]partnerbank.AccountType{
	AccountTypeChecking: partnerbank.AccountTypeChecking,
}

type AccountStatus string

const (
	AccountStatusActive    = AccountStatus("active")
	AccountStatusBlocked   = AccountStatus("blocked")
	AccountStatusClosed    = AccountStatus("closed")
	AccountStatusDormant   = AccountStatus("dormant")
	AccountStatusAbandoned = AccountStatus("abandoned")
	AccountStatusEscheated = AccountStatus("escheated")
	AccountStatusChargeOff = AccountStatus("charge_off")
)

var partnerAccountStatusFrom = map[partnerbank.AccountStatus]AccountStatus{
	partnerbank.AccountStatusActive:    AccountStatusActive,
	partnerbank.AccountStatusBlocked:   AccountStatusBlocked,
	partnerbank.AccountStatusLocked:    AccountStatusBlocked,
	partnerbank.AccountStatusClosed:    AccountStatusClosed,
	partnerbank.AccountStatusDormant:   AccountStatusDormant,
	partnerbank.AccountStatusAbandoned: AccountStatusAbandoned,
	partnerbank.AccountStatusEscheated: AccountStatusEscheated,
	partnerbank.AccountStatusChargeOff: AccountStatusChargeOff,
}

var partnerAccountStatusTo = map[AccountStatus]partnerbank.AccountStatus{
	AccountStatusActive:    partnerbank.AccountStatusActive,
	AccountStatusBlocked:   partnerbank.AccountStatusBlocked,
	AccountStatusClosed:    partnerbank.AccountStatusClosed,
	AccountStatusDormant:   partnerbank.AccountStatusDormant,
	AccountStatusAbandoned: partnerbank.AccountStatusAbandoned,
	AccountStatusEscheated: partnerbank.AccountStatusEscheated,
	AccountStatusChargeOff: partnerbank.AccountStatusChargeOff,
}

/*
CreateAccountRequest
Can be used when creating accounts for consumers and non-consumers

Non-Consumers:
- participants must contain at least two participant_user_id /participant_role groups.
- multiple_participants must equal true.
- One participant_user_id must specify a valid consumer identifier.
- The second participant_user_id (and any other additional participants) can specify
    either consumer or business identifiers.
*/
type CreateAccountRequest struct {
	AccountType          AccountType          `json:"account_type"`
	MultipleParticipants bool                 `json:"multiple_participants"`
	Participants         []ParticipantRequest `json:"participants"`
	BusinessType         BusinessType         `json:"business_type"`
}

type BusinessType string

const (
	BusinessTypeCorporate = BusinessType("business_corporate")
	BusinessTypeForeign   = BusinessType("foreign_holder")
	BusinessTypeLLC       = BusinessType("llc")
	BusinessTypeNonProfit = BusinessType("non_profit_organizations")
)

var partnerBusinessTypeFrom = map[partnerbank.BusinessEntity]BusinessType{
	partnerbank.BusinessEntitySoleProprietor:              BusinessTypeCorporate,
	partnerbank.BusinessEntityAssociation:                 BusinessTypeCorporate,
	partnerbank.BusinessEntityProfessionalAssociation:     BusinessTypeCorporate,
	partnerbank.BusinessEntitySingleMemberLLC:             BusinessTypeLLC,
	partnerbank.BusinessEntityLimitedLiabilityCompany:     BusinessTypeLLC,
	partnerbank.BusinessEntityGeneralPartnership:          BusinessTypeCorporate,
	partnerbank.BusinessEntityLimitedPartnership:          BusinessTypeCorporate,
	partnerbank.BusinessEntityLimitedLiabilityPartnership: BusinessTypeCorporate,
	partnerbank.BusinessEntityProfessionalCorporation:     BusinessTypeCorporate,
	partnerbank.BusinessEntityUnlistedCorporation:         BusinessTypeCorporate,
}

/*
CreateBankAccountResponse
Your application must store the value of account_id, account_number and routing_number so
that you can reference the business account record when accessing other Open Platform APIs.
{
  "account_id": "AC-7cb2a96c-9c23-4dd2-992f-edcd11d34be0",
  "account_number": "5953572185",
  "routing_number": "062001186"
}

Know Your Customer status and business accounts
The JSON response body for a new consumer record contains a Know Your Customer (KYC) status.
When creating a business account, all consumer identifiers supplied in the request body
must have a KYC status of APPROVED.

In sandbox applications, the KYC status by default equals APPROVED. This status enables your
sandbox application to immediately use the information in the consumer response to access
other API endpoints.

In production applications, the KYC status by default equals REVIEW. This status requires
the customer to provide supplementary identification such as a driver’s license or passport
as part of the approval process for new customers.
*/
type CreateBankAccountResponse struct {
	AccountID     string `json:"account_id"`
	AccountNumber string `json:"account_number"`
	RoutingNumber string `json:"routing_number"`
}

type GetBankAccountResponse struct {
	AccountID        string        `json:"account_id"`
	RoutingNumber    string        `json:"routing_number"`
	AccountNumber    string        `json:"account_number"`
	AccountType      AccountType   `json:"account_type"`
	Alias            string        `json:"account_alias"`
	Status           AccountStatus `json:"account_status"`
	OpeningDate      string        `json:"opening_date"`
	AvailableBalance float64       `json:"available_balance"`
	PostedBalance    float64       `json:"posted_balance"`
	Currency         Currency      `json:"currency"`
}

func (resp *GetBankAccountResponse) toPartnerGetBankAccountResponse(req partnerbank.APIRequest) (*partnerbank.GetBankAccountResponse, error) {
	accountType, ok := partnerAccountTypeTo[resp.AccountType]
	if !ok {
		return nil, errors.New("invalid account type")
	}

	status, ok := partnerAccountStatusTo[resp.Status]
	if !ok {
		return nil, errors.New("invalid account status")
	}

	t, err := time.Parse("2006-01-02", resp.OpeningDate)
	if err != nil {
		return nil, err
	}

	currency, ok := partnerCurrencyTo[resp.Currency]
	if !ok {
		iso := strings.ToUpper(strings.TrimSpace(resp.Currency.String())[:3])
		currency, ok = partnerCurrencyTo[Currency(iso)]
		if !ok {
			return nil, errors.New("invalid currency type")
		}
	}

	// Return full response object
	return &partnerbank.GetBankAccountResponse{
		AccountID:        partnerbank.AccountBankID(resp.AccountID),
		AccountType:      accountType,
		AccountNumber:    resp.AccountNumber,
		RoutingNumber:    resp.RoutingNumber,
		Alias:            resp.Alias,
		Status:           status,
		Opened:           t,
		AvailableBalance: resp.AvailableBalance,
		PostedBalance:    resp.PostedBalance,
		Currency:         currency,
	}, nil
}

// business account PATCH requests
// After a business account has been created, the account endpoint’s
// PATCH method can modify only a subset of the information submitted in the original request.
type PatchAccountRequest struct {
	Alias        string             `json:"account_alias,omitempty"`
	Status       AccountStatus      `json:"account_status,omitempty"`
	StatusReason AccountCloseReason `json:"account_status_reason,omitempty"`
}

type AccountCloseReason string

const (
	AccountCloseReasonDeath              = AccountCloseReason("death")
	AccountCloseReasonBankError          = AccountCloseReason("bank_error")
	AccountCloseReasonNonSufficientFunds = AccountCloseReason("nsf_funds")
	AccountCloseReasonLegal              = AccountCloseReason("legal")
	AccountCloseReasonConsolidating      = AccountCloseReason("consolidating_accounts")
	AccountCloseReasonBank               = AccountCloseReason("bank_request")
	AccountCloseReasonFraud              = AccountCloseReason("fraud")
	AccountCloseReasonCustomer           = AccountCloseReason("customer_request")
)

var partnerAccountCloseReasonFrom = map[partnerbank.AccountCloseReason]AccountCloseReason{
	partnerbank.AccountCloseReasonDeath:              AccountCloseReasonDeath,
	partnerbank.AccountCloseReasonBankError:          AccountCloseReasonBankError,
	partnerbank.AccountCloseReasonNonSufficientFunds: AccountCloseReasonNonSufficientFunds,
	partnerbank.AccountCloseReasonLegal:              AccountCloseReasonLegal,
	partnerbank.AccountCloseReasonConsolidating:      AccountCloseReasonConsolidating,
	partnerbank.AccountCloseReasonBank:               AccountCloseReasonBank,
	partnerbank.AccountCloseReasonFraud:              AccountCloseReasonFraud,
	partnerbank.AccountCloseReasonCustomer:           AccountCloseReasonCustomer,
}

type AccountBlockType string

const (
	AccountBlockTypeDebits  = AccountBlockType("debits")
	AccountBlockTypeCredits = AccountBlockType("credits")
	AccountBlockTypeChecks  = AccountBlockType("checks")
	AccountBlockTypeAll     = AccountBlockType("all")
)

var partnerAccountBlockTypeFrom = map[partnerbank.AccountBlockType]AccountBlockType{
	partnerbank.AccountBlockTypeDebit:  AccountBlockTypeDebits,
	partnerbank.AccountBlockTypeCredit: AccountBlockTypeCredits,
	partnerbank.AccountBlockTypeCheck:  AccountBlockTypeChecks,
	partnerbank.AccountBlockTypeAll:    AccountBlockTypeAll,
}

var partnerAccountBlockTypeTo = map[AccountBlockType]partnerbank.AccountBlockType{
	AccountBlockTypeDebits:  partnerbank.AccountBlockTypeDebit,
	AccountBlockTypeCredits: partnerbank.AccountBlockTypeCredit,
	AccountBlockTypeChecks:  partnerbank.AccountBlockTypeCheck,
	AccountBlockTypeAll:     partnerbank.AccountBlockTypeAll,
}

type AccountBlockStatus string

const (
	AccountBlockStatusActive   = AccountBlockStatus("active")
	AccountBlockStatusCanceled = AccountBlockStatus("canceled")
)

var partnerAccountBlockStatusFrom = map[partnerbank.AccountBlockStatus]AccountBlockStatus{
	partnerbank.AccountBlockStatusActive:   AccountBlockStatusActive,
	partnerbank.AccountBlockStatusCanceled: AccountBlockStatusCanceled,
}

var partnerAccountBlockStatusTo = map[AccountBlockStatus]partnerbank.AccountBlockStatus{
	AccountBlockStatusActive:   partnerbank.AccountBlockStatusActive,
	AccountBlockStatusCanceled: partnerbank.AccountBlockStatusCanceled,
}

type CreateAccountBlockRequest struct {
	BlockType AccountBlockType `json:"block_type"`
	Reason    string           `json:"block_reason"`
}

type CreateAccountBlockResponse struct {
	BlockID string `json:"block_id"`
}

type GetAccountBlockResponse struct {
	BlockID      string             `json:"block_id"`
	BlockType    AccountBlockType   `json:"block_type"`
	Status       AccountBlockStatus `json:"block_status"`
	Reason       string             `json:"block_reason"`
	CreationDate string             `json:"block_creation_date"`
	CancelDate   *time.Time         `json:"cancellation_date"`
}

func (r GetAccountBlockResponse) toPartnerAccountBlockResponse() (*partnerbank.AccountBlockResponse, error) {
	date, err := time.Parse("2006-01-02", r.CreationDate)
	if err != nil {
		date, err = time.Parse("2006-01-02T15:04:05Z", r.CreationDate)
		if err != nil {
			return nil, err
		}
	}

	blockType, ok := partnerAccountBlockTypeTo[r.BlockType]
	if !ok {
		return nil, errors.New("invalid account block type")
	}

	status, ok := partnerAccountBlockStatusTo[r.Status]
	if !ok {
		return nil, errors.New("invalid account block status")
	}

	return &partnerbank.AccountBlockResponse{
		BlockID:      partnerbank.AccountBlockBankID(r.BlockID),
		Type:         blockType,
		Status:       status,
		Reason:       r.Reason,
		BlockDate:    date,
		CanceledDate: r.CancelDate,
	}, nil
}

type GetAccountBlocks struct {
	Blocks []GetAccountBlockResponse `json:"blocks"`
}

type AccountStatement struct {
	StatementID  string `json:"statement_id"`
	Description  string `json:"description"`
	CreationDate string `json:"creation_date"`
	TotalPages   int    `json:"total_pages"`
}

func (s *AccountStatement) toPartnerAccountStatement() (*partnerbank.AccountStatementResponse, error) {
	date, err := time.Parse("2006-01-02", s.CreationDate)
	if err != nil {
		return nil, err
	}

	return &partnerbank.AccountStatementResponse{
		StatementID: partnerbank.AccountStatementBankID(s.StatementID),
		Description: s.Description,
		Created:     date,
		PageCount:   s.TotalPages,
	}, nil
}

type GetAccountStatements struct {
	Statements []AccountStatement `json:"statements"`
}
