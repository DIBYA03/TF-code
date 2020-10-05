package business

import (
	"time"

	"github.com/jmoiron/sqlx/types"
	"github.com/wiseco/core-platform/services"
	"github.com/wiseco/core-platform/shared"
)

type ExternalBankAccount struct {
	ID                  string            `json:"id" db:"id"`
	BusinessID          shared.BusinessID `json:"businessId" db:"business_id"`
	LinkedAccountID     *string           `json:"linkedAccountId" db:"linked_account_id"`
	PartnerAccountID    string            `json:"partnerAccountId" db:"partner_account_id"`
	PartnerName         string            `json:"partnerName" db:"partner_name"`
	AccountName         string            `json:"accountName" db:"account_name"`
	OfficialAccountName string            `json:"officialAccountName" db:"official_account_name"`
	AccountType         string            `json:"accountType" db:"account_type"`
	AccountSubtype      string            `json:"accountSubtype" db:"account_subtype"`
	AccountNumber       string            `json:"accountNumber" db:"account_number"`
	RoutingNumber       string            `json:"routingNumber" db:"routing_number"`
	WireRouting         string            `json:"wireRouting" db:"wire_routing"`
	AvailableBalance    float64           `json:"availableBalance" db:"available_balance"`
	PostedBalance       float64           `json:"postedBalance" db:"posted_balance"`
	Currency            string            `json:"currency" db:"currency"`
	LastLogin           time.Time         `json:"lastLogin" db:"last_login"`
	Created             time.Time         `json:"created" db:"created"`
	Modified            time.Time         `json:"modified" db:"modified"`
}

type ExternalBankAccountUpdate struct {
	BusinessID          shared.BusinessID `db:"business_id"`
	LinkedAccountID     *string           `db:"linked_account_id"`
	PartnerAccountID    *string           `db:"partner_account_id"`
	PartnerName         string            `db:"partner_name"`
	AccountName         string            `db:"account_name"`
	OfficialAccountName string            `db:"official_account_name"`
	AccountType         string            `db:"account_type"`
	AccountSubtype      string            `db:"account_subtype"`
	AccountNumber       string            `db:"account_number"`
	RoutingNumber       string            `db:"routing_number"`
	WireRouting         string            `db:"wire_routing"`
	AvailableBalance    *float64          `db:"available_balance"`
	PostedBalance       *float64          `db:"posted_balance"`
	Currency            *string           `db:"currency"`
	Owner               []ExternalBankAccountOwnerCreate
	LastLogin           *time.Time `db:"last_login"`
}

type ExternalBankAccountOwner struct {
	ID                string               `json:"id" db:"id"`
	ExternalAccountID string               `json:"externalAccountId" db:"external_bank_account_id"`
	AccountHolderName services.StringArray `json:"accountHolderName" db:"account_holder_name"`
	Phone             types.JSONText       `json:"phone" db:"phone"`
	Email             types.JSONText       `json:"email" db:"email"`
	OwnerAddress      types.JSONText       `json:"ownerAddress" db:"owner_address"`
	Created           time.Time            `json:"created" db:"created"`
	Modified          time.Time            `json:"modified" db:"modified"`
}

type ExternalBankAccountOwnerCreate struct {
	ExternalAccountID string               `db:"external_bank_account_id"`
	AccountHolderName services.StringArray `db:"account_holder_name"`
	Phone             types.JSONText       `db:"phone"`
	Email             types.JSONText       `db:"email"`
	OwnerAddress      types.JSONText       `db:"owner_address"`
}

type VerificationStatus string

const (
	VerificationStatusSucceeded  = VerificationStatus("succeeded")
	VerificationStatusUnverified = VerificationStatus("unverified")
)

const (
	NameMismatch                      = "NAME_MISMATCH"
	EmailMismatch                     = "EMAIL_MISMATCH"
	PhoneMismatch                     = "PHONE_MISMATCH"
	AddressMismatch                   = "ADDRESS_MISMATCH"
	BusinessMemberNameMatch           = "BUSINESS_MEMBER_NAME_MATCH"
	BusinessNameMatch                 = "BUSINESS_NAME_MATCH"
	BusinessMemberPhoneMatch          = "BUSINESS_MEMBER_PHONE_MATCH"
	BusinessPhoneMatch                = "BUSINESS_PHONE_MATCH"
	BusinessMemberEmailMatch          = "BUSINESS_MEMBER_EMAIL_MATCH"
	BusinessEmailMatch                = "BUSINESS_EMAIL_MATCH"
	BusinessMemberMailingAddressMatch = "BUSINESS_MEMBER_MAILING_ADDRESS_MATCH"
	BusinessMemberLegalAddressMatch   = "BUSINESS_MEMBER_LEGAL_ADDRESS_MATCH"
	BusinessMemberWorkAddressMatch    = "BUSINESS_MEMBER_WORK_ADDRESS_MATCH"
	BusinessLegalAddressMatch         = "BUSINESS_LEGAL_ADDRESS_MATCH"
	BusinessMailingAddressMatch       = "BUSINESS_MAILING_ADDRESS_MATCH"
	BusinessHeadquarterAddressMatch   = "BUSINESS_HEADQUARTER_ADDRESS_MATCH"
)

const (
	ExternalVerificationErrorNameMismatch    = "Account name mismatch"
	ExternalVerificationErrorAddressMismatch = "Account address mismatch"
	ExternalVerificationErrorEmailMismatch   = "Account email mismatch"
	ExternalVerificationErrorPhoneMismatch   = "Account phone mismatch"
	ExternalVerificationErrorGeneric         = "Unable to verify account owner"
)

type ExternalAccountVerificationRequest struct {
	BusinessID       shared.BusinessID
	AccessToken      string
	PartnerItemID    string
	PartnerAccountID string
}

type ExternalAccountVerificationCreate struct {
	BusinessID         shared.BusinessID    `db:"business_id"`
	ExternalAccountID  string               `db:"external_bank_account_id"`
	SourceIPAddress    string               `db:"source_ip_address"`
	AccessToken        string               `db:"access_token"`
	PartnerItemID      string               `db:"partner_item_id"`
	VerificationStatus VerificationStatus   `db:"verification_status"`
	VerificationResult services.StringArray `db:"verification_result"`
}

type ExternalAccountVerificationResult struct {
	ID                 string               `db:"id"`
	BusinessID         shared.BusinessID    `db:"business_id"`
	ExternalAccountID  string               `db:"external_bank_account_id"`
	SourceIPAddress    string               `db:"source_ip_address"`
	AccessToken        string               `db:"access_token"`
	PartnerItemID      string               `db:"partner_item_id"`
	VerificationStatus VerificationStatus   `db:"verification_status"`
	VerificationResult services.StringArray `db:"verification_result"`
	Created            time.Time            `db:"created"`
}
