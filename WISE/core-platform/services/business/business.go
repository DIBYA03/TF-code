/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

// Package for business services
package business

import (
	"strings"
	"time"

	"github.com/wiseco/core-platform/services"
	"github.com/wiseco/core-platform/shared"
)

// Business entity object
type Business struct {
	// Business id
	ID shared.BusinessID `json:"id" db:"id"`

	//Owner business id
	OwnerID shared.UserID `json:"ownerId" db:"owner_id"`

	// Generated employer number (8-Digit Code)
	EmployerNumber string `json:"employerNumber" db:"employer_number"`

	// Company legal name
	LegalName *string `json:"legalName" db:"legal_name"`

	// DBA name
	DBA services.StringArray `json:"dba" db:"dba"`

	// Is phone verified?
	PhoneVerified bool `json:"phoneVerified" db:"phone_verified"`

	// Business activity type e.g. cash, check, domesticAch
	ActivityType services.StringArray `json:"activityType" db:"activity_type"`

	//Business Handles cash
	HandlesCash bool `json:"handlesCash" db:"handles_cash"`

	// Is email verified?
	EmailVerified bool `json:"emailVerified" db:"email_verified"`

	// Business members >= 25% owners
	Members []BusinessMember `json:"members" db:"members"`

	// Know your customer status e.g. approved, declined, etc.
	KYCStatus services.KYCStatus `json:"kycStatus" db:"kyc_status"`

	// Entity type e.g. llc, etc
	EntityType *string `json:"entityType" db:"entity_type"`

	// Is business restricted
	IsRestricted bool `json:"isRestricted" db:"is_restricted"`

	// Industry type e.g. hotels, etc
	IndustryType *string `json:"industryType" db:"industry_type"`

	// Tax id number e.g. ein
	TaxID *services.TaxID `json:"taxIdMasked" db:"tax_id"`

	// Tax id type e.g. ssn, ein, etc
	TaxIDType *services.TaxIDType `json:"taxIdType" db:"tax_id_type"`

	// Origin or country of incorporation
	OriginCountry *string `json:"originCountry" db:"origin_country"`

	// Origin or state of incorporation
	OriginState *string `json:"originState" db:"origin_state"`

	// Origin or date of incorporation
	OriginDate *shared.Date `json:"originDate" db:"origin_date"`

	// Freeform business purpose
	Purpose *string `json:"purpose" db:"purpose"`

	// Business operations e.g. foreign, domestic, etc.
	OperationType *string `json:"operationType" db:"operation_type"`

	// Business email
	Email *string `json:"email" db:"email"`

	// Email ID
	EmailID shared.EmailID `json:"emailID" db:"email_id"`

	// Business phone number
	Phone *string `json:"phone" db:"phone"`

	// Legal address
	LegalAddress *services.Address `json:"legalAddress" db:"legal_address"`

	// Headquarter address
	HeadquarterAddress *services.Address `json:"headquarterAddress" db:"headquarter_address"`

	// Mailing address
	MailingAddress *services.Address `json:"mailingAddress" db:"mailing_address"`

	// Formation Document Id
	FormationDocumentID *shared.BusinessDocumentID `json:"formationDocumentId" db:"formation_document_id"`

	//Deactivated
	Deactivated *time.Time `json:"deactivated" db:"deactivated"`

	// Created timestamp
	Created time.Time `json:"created" db:"created"`

	// Modified timestamp
	Modified time.Time `json:"modified" db:"modified"`

	// ACH Pull Enabled
	ACHPullEnabled bool `json:"achPullEnabled"`

	// Business website
	Website *services.Website `json:"website" db:"website"`

	// Business online info
	OnlineInfo *string `json:"onlineInfo" db:"online_info"`

	//Subscription decision date
	SubscriptionDecisionDate *time.Time `json:"subscriptionDecisionDate" db:"subscription_decision_date"`

	// Subscription status
	SubscriptionStatus *services.SubscriptionStatus `json:"subscriptionStatus" db:"subscription_status"`

	// Subscription start date
	SubscriptionStartDate *shared.Date `json:"subscriptionStartDate" db:"subscription_start_date"`
}

type BusinessCreate struct {
	//Owner business id
	OwnerID shared.UserID `db:"owner_id"`

	// Company legal name
	LegalName *string `json:"legalName" db:"legal_name"`

	// Generated employer number (8-Digit Code)
	EmployerNumber string `db:"employer_number"`

	// DBA name
	DBA services.StringArray `json:"dba" db:"dba"`

	// Business activity type e.g. cash, check, domesticAch
	ActivityType services.StringArray `json:"activityType" db:"activity_type"`

	// Is this a restricted business
	IsRestrictedBusiness bool `db:"is_restricted"`

	// Business usage purpose
	Purpose *string `json:"purpose" db:"purpose"`

	// Entity type e.g. llc, etc
	EntityType *string `json:"entityType" db:"entity_type"`

	// Industry type e.g. hotels, etc
	IndustryType *string `json:"industryType" db:"industry_type"`

	// Tax id number e.g. ein
	TaxID *services.TaxID `json:"taxId,omitempty" db:"tax_id"`

	// Tax id type e.g. ssn, ein, etc
	TaxIDType *services.TaxIDType `json:"taxIdType,omitempty" db:"tax_id_type"`

	// Origin or country of incorporation
	OriginCountry *string `json:"originCountry,omitempty" db:"origin_country"`

	// Origin or state of incorporation
	OriginState *string `json:"originState,omitempty" db:"origin_state"`

	// Origin or date of incorporation
	OriginDate *shared.Date `json:"originDate,omitempty" db:"origin_date"`

	// Business operations e.g. foreign, domestic, etc.
	OperationType *string `json:"operationType" db:"operation_type"`

	// Business email
	Email *string `json:"email,omitempty" db:"email"`

	// Email ID
	EmailID shared.EmailID `json:"emailID" db:"email_id"`

	// Business phone number
	Phone *string `json:"phone,omitempty" db:"phone"`

	// Legal address
	LegalAddress *services.Address `json:"legalAddress,omitempty" db:"legal_address"`

	// Headquarter address
	HeadquarterAddress *services.Address `json:"headquarterAddress" db:"headquarter_address"`

	// Mailing address
	MailingAddress *services.Address `json:"mailingAddress,omitempty" db:"mailing_address"`

	// Business website
	Website *services.Website `json:"website" db:"website"`

	// Business online info
	OnlineInfo *string `json:"onlineInfo" db:"online_info"`

	//Subscription decision date
	SubscriptionDecisionDate *time.Time `json:"subscriptionDecisionDate" db:"subscription_decision_date"`

	// Subscription status
	SubscriptionStatus *services.SubscriptionStatus `json:"subscriptionStatus" db:"subscription_status"`

	// Subscription start date
	SubscriptionStartDate *shared.Date `json:"subscriptionStartDate" db:"subscription_start_date"`
}

type BusinessUpdate struct {
	// Business id
	ID shared.BusinessID `json:"id" db:"id"`

	// Company legal name
	LegalName *string `json:"legalName,omitempty" db:"legal_name"`

	// Know your customer status e.g. approved, declined, etc.
	KYCStatus *services.KYCStatus `json:"kycStatus,omitempty" db:"kyc_status"`

	// Entity type e.g. llc, etc
	EntityType *string `json:"entityType,omitempty" db:"entity_type"`

	// Industry type e.g. hotels, etc
	IndustryType *string `json:"industryType,omitempty" db:"industry_type"`

	// DBA name
	DBA *services.StringArray `json:"dba,omitempty" db:"dba"`

	// Tax id number e.g. ein
	TaxID *services.TaxID `json:"taxId,omitempty" db:"tax_id"`

	// Tax id type e.g. ssn, ein, etc
	TaxIDType *services.TaxIDType `json:"taxIdType,omitempty" db:"tax_id_type"`

	// Origin or country of incorporation
	OriginCountry *string `json:"originCountry,omitempty" db:"origin_country"`

	// Origin or state of incorporation
	OriginState *string `json:"originState,omitempty" db:"origin_state"`

	// Origin or date of incorporation
	OriginDate *shared.Date `json:"originDate,omitempty" db:"origin_date"`

	// Freeform business purpose
	Purpose *string `json:"purpose,omitempty" db:"purpose"`

	// Business operations e.g. foreign, domestic, etc.
	OperationType *string `json:"operationType" db:"operation_type"`

	// Business email
	Email *string `json:"email,omitempty" db:"email"`

	// Email ID
	EmailID shared.EmailID `json:"emailID" db:"email_id"`

	// Business phone number
	Phone *string `json:"phone,omitempty" db:"phone"`

	// Business activity type e.g. cash, check, domesticAch
	ActivityType *services.StringArray `json:"activityType,omitempty" db:"activity_type"`

	// Is this a restricted business
	IsRestricted *bool `json:"isRestricted" db:"is_restricted"`

	// Legal address
	LegalAddress *services.Address `json:"legalAddress,omitempty" db:"legal_address"`

	// Headquarter address
	HeadquarterAddress *services.Address `json:"headquarterAddress" db:"headquarter_address"`

	// Mailing address
	MailingAddress *services.Address `json:"mailingAddress,omitempty" db:"mailing_address"`

	// Business website
	Website *services.Website `json:"website" db:"website"`

	// Business online info
	OnlineInfo *string `json:"onlineInfo" db:"online_info"`
}

type BusinessSubscriptionUpdate struct {
	ID shared.BusinessID `json:"-" db:"id"`

	//Subscription decision date
	SubscriptionDecisionDate *time.Time `json:"-" db:"subscription_decision_date"`

	// Subscription status
	SubscriptionStatus *services.SubscriptionStatus `json:"-" db:"subscription_status"`

	// Subscription start date
	SubscriptionStartDate *shared.Date `json:"-" db:"subscription_start_date"`
}

type BusinessVerificationUpdate struct {
	// User id
	ID shared.BusinessID `json:"id" db:"id"`

	// KYC Status
	KYCStatus services.KYCStatus `json:"kycStatus" db:"kyc_status"`
}

type BusinessKYCResponse struct {
	Status      string   `json:"status"`
	ReviewItems []string `json:"reviewItems"`
	Business    Business `json:"business"`
}

type BusinessKYCError struct {
	RawError  error     `json:"-"`
	ErrorType string    `json:"errorType"`
	Values    []string  `json:"values"`
	Business  *Business `json:"user"`
}

func (e BusinessKYCError) Error() string {
	switch e.ErrorType {
	case services.KYCErrorTypeOther:
		return e.RawError.Error()
	case services.KYCErrorTypeInProgress:
		return "Business verification in progress"
	case services.KYCErrorTypeParam:
		return "Business verification invalid parameter(s)"
	case services.KYCErrorTypeReview:
		return "Business verification user in review"
	case services.KYCErrorTypeDeactivated:
		return "Business deactivated"
	case services.KYCErrorTypeRestricted:
		return "Business restricted"
	default:
		return "Error with unknown KYC error type"
	}
}

func (b *Business) HasDBA() bool {
	if b == nil {
		return false
	}

	if len(b.DBA) == 0 {
		return false
	}

	return len(strings.TrimSpace(b.DBA[0])) > 0
}

// Return DBA else LegalName else empty string
func (b *Business) Name() string {
	if b == nil {
		return ""
	}

	if b.HasDBA() {
		// Use first DBA
		return b.DBA[0]
	}

	if b.LegalName != nil {
		return *b.LegalName
	}

	return ""
}
