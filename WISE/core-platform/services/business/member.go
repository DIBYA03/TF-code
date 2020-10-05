/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

// Package for business services
package business

import (
	"strings"
	"time"

	partnerbank "github.com/wiseco/core-platform/partner/bank"
	"github.com/wiseco/core-platform/services"
	"github.com/wiseco/core-platform/shared"
)

type TitleType string

const (
	TitleTypeCEO       = TitleType("chiefExecutiveOfficer")
	TitleTypeCFO       = TitleType("chiefFinancialOfficer")
	TitleTypeCOO       = TitleType("chiefOperatingOfficer")
	TitleTypePresident = TitleType("president")
	TitleTypeVP        = TitleType("vicePresident")
	TitleTypeSVP       = TitleType("seniorVicePresident")
	TitleTypeTreasurer = TitleType("treasurer")
	TitleSecretary     = TitleType("secretary")
	TitleTypeGP        = TitleType("generalPartner")
	TitleManager       = TitleType("manager")
	TitleMember        = TitleType("member")
	TitleOwner         = TitleType("owner")
	TitleTypeOther     = TitleType("other")
)

var PartnerTitleTypeFrom = map[TitleType]partnerbank.BusinessMemberTitle{
	TitleTypeCEO:       partnerbank.BusinessMemberTitleCEO,
	TitleTypeCFO:       partnerbank.BusinessMemberTitleCFO,
	TitleTypeCOO:       partnerbank.BusinessMemberTitleCOO,
	TitleTypePresident: partnerbank.BusinessMemberTitlePresident,
	TitleTypeVP:        partnerbank.BusinessMemberTitleVP,
	TitleTypeSVP:       partnerbank.BusinessMemberTitleSVP,
	TitleTypeTreasurer: partnerbank.BusinessMemberTitleTreasurer,
	TitleSecretary:     partnerbank.BusinessMemberTitleSecretary,
	TitleTypeGP:        partnerbank.BusinessMemberTitleGP,
	TitleManager:       partnerbank.BusinessMemberTitleManager,
	TitleMember:        partnerbank.BusinessMemberTitleMember,
	TitleOwner:         partnerbank.BusinessMemberTitleOwner,
	TitleTypeOther:     partnerbank.BusinessMemberTitleOther,
}

type BusinessMember struct {
	// Business member id
	ID shared.BusinessMemberID `json:"id" db:"id"`

	// Consumer Id
	ConsumerID shared.ConsumerID `json:"consumerId" db:"consumer_id"`

	// BankID bank id (CO)
	BankID *partnerbank.ConsumerBankID `json:"bankId"`

	// Business id
	BusinessID shared.BusinessID `json:"businessId" db:"business_id"`

	// Title type
	TitleType TitleType `json:"titleType" db:"title_type"`

	// Title other
	TitleOther *string `json:"titleOther" db:"title_other"`

	// Ownership percentage
	Ownership int `json:"ownership" db:"ownership"`

	// Is controlling manager of business
	IsControllingManager bool `json:"isControllingManager" db:"is_controlling_manager"`

	// Deactivated
	Deactivated *time.Time `json:"deactivated" db:"deactivated"`

	// Created time
	Created time.Time `json:"created" db:"created"`

	// Modified time
	Modified time.Time `json:"modified" db:"modified"`

	/*
	 * Consumer Properties
	 */

	// First name
	FirstName string `json:"firstName" db:"first_name"`

	// Middle name
	MiddleName string `json:"middleName" db:"middle_name"`

	// Last name
	LastName string `json:"lastName" db:"last_name"`

	// Email
	Email string `json:"email" db:"email"`

	// Phone
	Phone string `json:"phone" db:"phone"`

	// Date of birth
	DateOfBirth *shared.Date `json:"dateOfBirth" db:"date_of_birth"`

	// Tax id number
	TaxID *services.TaxID `json:"taxIdMasked" db:"tax_id"`

	// Tax id type e.g. SSN or ITIN
	TaxIDType *services.TaxIDType `json:"taxIdType" db:"tax_id_type"`

	// KYC Status
	KYCStatus services.KYCStatus `json:"kycStatus" db:"kyc_status"`

	// Legal address
	LegalAddress *services.Address `json:"legalAddress" db:"legal_address"`

	// Mailing address
	MailingAddress *services.Address `json:"mailingAddress" db:"mailing_address"`

	// Work address
	WorkAddress *services.Address `json:"workAddress" db:"work_address"`

	// Residency
	Residency *services.Residency `json:"residency" db:"residency"`

	// List of citizenships
	CitizenshipCountries services.StringArray `json:"citizenshipCountries" db:"citizenship_countries"`

	// User's occupation
	Occupation *string `json:"occupation" db:"occupation"`

	// User income source e.g. salary or inheritance
	IncomeType services.StringArray `json:"incomeType" db:"income_type"`

	// User activity type
	ActivityType services.StringArray `json:"activityType" db:"activity_type"`

	// Is this a restricted person
	IsRestricted bool `json:"isRestricted" db:"is_restricted"`

	// Email ID
	EmailID shared.EmailID `json:"emailID" db:"email_id"`

	/* User ID */
	UserID *shared.UserID `json:"userId" db:"user_id"`

	/* User Phone */
	UserPhone *string `json:"userPhone" db:"user_phone"`
}

func (m *BusinessMember) Name() string {
	return strings.TrimSpace(m.FirstName + " " + m.LastName)
}

type BusinessMemberCreate struct {
	// Business id
	BusinessID shared.BusinessID `json:"businessId" db:"business_id"`

	// Title type
	TitleType TitleType `json:"titleType" db:"title_type"`

	// Title other
	TitleOther *string `json:"titleOther,omitempty" db:"title_other"`

	// Ownership percentage
	Ownership int `json:"ownership,omitempty" db:"ownership"`

	// Is controlling manager
	IsControllingManager bool `json:"isControllingManager" db:"is_controlling_manager"`

	/* User ID */
	UserID *shared.UserID `json:"userId"`

	/*
	 * Consumer Properties
	 */

	// First name
	FirstName string `json:"firstName" db:"first_name"`

	// Middle name
	MiddleName string `json:"middleName,omitempty" db:"middle_name"`

	// Last name
	LastName string `json:"lastName" db:"last_name"`

	// Email
	Email string `json:"email" db:"email"`

	// Phone
	Phone string `json:"phone" db:"phone"`

	// Date of birth
	DateOfBirth *shared.Date `json:"dateOfBirth,omitempty" db:"date_of_birth"`

	// Tax id number
	TaxID *services.TaxID `json:"taxId" db:"tax_id"`

	// Tax id type e.g. SSN or ITIN
	TaxIDType *services.TaxIDType `json:"taxIdType"  db:"tax_id_type"`

	// Legal address
	LegalAddress *services.Address `json:"legalAddress,omitempty" db:"legal_address"`

	// Residency
	Residency *services.Residency `json:"residency,omitempty" db:"residency"`

	// List of citizenships
	CitizenshipCountries services.StringArray `json:"citizenshipCountries" db:"citizenship_countries"`

	// Consumer's occupation
	Occupation string `json:"occupation" db:"occupation"`

	// Consumer income source e.g. salary or inheritance
	IncomeType *services.StringArray `json:"incomeType" db:"income_type"`

	// Consumer activity type
	ActivityType *services.StringArray `json:"activityType" db:"activity_type"`
}

type BusinessMemberUpdate struct {
	// Member id
	ID shared.BusinessMemberID `json:"id" db:"id"`

	// Title type
	TitleType *TitleType `json:"titleType,omitempty" db:"title_type"`

	// Title other
	TitleOther *string `json:"titleOther,omitempty" db:"title_other"`

	// Ownership percentage
	Ownership *int `json:"ownership,omitempty" db:"ownership"`

	// Is controlling manager
	IsControllingManager *bool `json:"isControllingManager" db:"is_controlling_manager"`

	/*
	 * Consumer Properties
	 */

	// First name
	FirstName *string `json:"firstName,omitempty" db:"first_name"`

	// Middle name
	MiddleName *string `json:"middleName,omitempty" db:"middle_name"`

	// Last name
	LastName *string `json:"lastName,omitempty" db:"last_name"`

	// Email
	Email *string `json:"email,omitempty" db:"email"`

	// Phone
	Phone string `json:"phone,omitempty" db:"phone"`

	// Date of birth
	DateOfBirth *shared.Date `json:"dateOfBirth,omitempty" db:"date_of_birth"`

	// Tax id number
	TaxID *services.TaxID `json:"taxId,omitempty" db:"tax_id"`

	// Tax id type e.g. SSN or ITIN
	TaxIDType *services.TaxIDType `json:"taxIdType,omitempty"  db:"tax_id_type"`

	// Legal address
	LegalAddress *services.Address `json:"legalAddress,omitempty" db:"legal_address"`

	MailingAddress *services.Address `json:"mailingAddress" db:"mailing_address"`

	// Residency
	Residency *services.Residency `json:"residency,omitempty" db:"residency"`

	// List of citizenships
	CitizenshipCountries *services.StringArray `json:"citizenshipCountries,omitempty" db:"citizenship_countries"`

	// Consumer's occupation
	Occupation *string `json:"occupation,omitempty" db:"occupation"`

	// Consumer income source e.g. salary or inheritance
	IncomeType *services.StringArray `json:"incomeType,omitempty" db:"income_type"`

	// Consumer activity type
	ActivityType *services.StringArray `json:"activityType,omitempty" db:"activity_type"`
}

type MemberKYCResponse struct {
	Status      services.KYCStatus `json:"status"`
	ReviewItems []string           `json:"reviewItems"`
	Member      BusinessMember     `json:"member"`
}

type MemberKYCError struct {
	RawError  error                    `json:"-"`
	ErrorType string                   `json:"errorType"`
	Values    []string                 `json:"values"`
	MemberID  *shared.BusinessMemberID `json:"memberId"`
}

func (e MemberKYCError) Error() string {
	switch e.ErrorType {
	case services.KYCErrorTypeOther:
		return e.RawError.Error()
	case services.KYCErrorTypeInProgress:
		return "Consumer verification in progress"
	case services.KYCErrorTypeParam:
		return "Consumer verification invalid parameter(s)"
	case services.KYCErrorTypeReview:
		return "Consumer verification user in review"
	case services.KYCErrorTypeDeactivated:
		return "Consumer deactivated"
	case services.KYCErrorTypeRestricted:
		return "Consumer restricted"
	default:
		return "Error with unknown KYC error type"
	}
}
