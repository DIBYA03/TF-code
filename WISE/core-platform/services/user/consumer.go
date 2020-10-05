/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

// Package for all user related services
package user

import (
	"strings"
	"time"

	"github.com/wiseco/core-platform/services"
	"github.com/wiseco/core-platform/shared"
)

// Consumer entity object provides access to accounts, businesses, and cards
type Consumer struct {

	// Consumer id (uuid)
	ID shared.ConsumerID `json:"id" db:"id"`

	// First name
	FirstName string `json:"firstName" db:"first_name"`

	// Middle name
	MiddleName string `json:"middleName" db:"middle_name"`

	// Last name
	LastName string `json:"lastName" db:"last_name"`

	// Email
	Email *string `json:"email" db:"email"`

	// Email ID
	EmailID shared.EmailID `json:"emailID" db:"email_id"`

	// Phone
	Phone *string `json:"phone" db:"phone"`

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

	// Consumer's occupation
	Occupation *string `json:"occupation" db:"occupation"`

	// Consumer income source e.g. salary or inheritance
	IncomeType services.StringArray `json:"incomeType" db:"income_type"`

	// Consumer activity type
	ActivityType services.StringArray `json:"activityType" db:"activity_type"`

	// Is this a restricted person
	IsRestricted bool `json:"isRestricted" db:"is_restricted"`

	// Deactivated
	Deactivated *time.Time `json:"deactivated" db:"deactivated"`

	// Created time
	Created time.Time `json:"created" db:"created"`

	// Last modified time
	Modified time.Time `json:"modified" db:"modified"`
}

func (c *Consumer) FullName() string {
	if c.MiddleName != "" {
		return strings.TrimSpace(c.FirstName + " " + c.MiddleName + " " + c.LastName)
	}

	return strings.TrimSpace(c.FirstName + " " + c.LastName)
}

// ConsumerAuthCreate is a convenniece struct to pass on to `user_service`
// so we can create a minimal user on `DB`
type ConsumerAuthCreate struct {
	// Phone number
	Phone *string `db:"phone"`

	// Email
	Email *string `json:"email" db:"email"`

	// Email ID
	EmailID shared.EmailID `json:"emailID" db:"email_id"`
}

type ConsumerCreate struct {
	// First name
	FirstName string `json:"firstName" db:"first_name"`

	// Middle name
	MiddleName string `json:"middleName" db:"middle_name"`

	// Last name
	LastName string `json:"lastName" db:"last_name"`

	// Email
	Email *string `json:"email" db:"email"`

	// Email ID
	EmailID shared.EmailID `json:"emailID" db:"email_id"`

	// Phone
	Phone *string `json:"phone" db:"phone"`

	// Date of birth
	DateOfBirth *shared.Date `json:"dateOfBirth" db:"date_of_birth"`

	// Tax id number
	TaxID *services.TaxID `json:"taxId" db:"tax_id"`

	// Tax id type e.g. SSN or ITIN
	TaxIDType *services.TaxIDType `json:"taxIdType" db:"tax_id_type"`

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

	// Consumer's occupation
	Occupation *string `json:"occupation" db:"occupation"`

	// Consumer income source e.g. salary or inheritance
	IncomeType *services.StringArray `json:"incomeType" db:"income_type"`

	// Consumer activity type
	ActivityType *services.StringArray `json:"activityType" db:"activity_type"`

	// Is this a business member?
	IsBusinessMember bool
}

type ConsumerUpdate struct {
	// Consumer id
	ID shared.ConsumerID `json:"userId" db:"id"`

	// First name
	FirstName *string `json:"firstName" db:"first_name"`

	// Middle name
	MiddleName *string `json:"middleName" db:"middle_name"`

	// Last name
	LastName *string `json:"lastName" db:"last_name"`

	// Email
	Email *string `json:"email" db:"email"`

	// Email ID
	EmailID shared.EmailID `json:"emailID" db:"email_id"`

	// Phone
	Phone *string `json:"phone" db:"phone"`

	// Date of birth
	DateOfBirth *shared.Date `json:"dateOfBirth" db:"date_of_birth"`

	// Tax id number
	TaxID *services.TaxID `json:"taxId" db:"tax_id"`

	// Tax id type e.g. SSN or ITIN
	TaxIDType *services.TaxIDType `json:"taxIdType" db:"tax_id_type"`

	// Legal address
	LegalAddress *services.Address `json:"legalAddress" db:"legal_address"`

	// Mailing address
	MailingAddress *services.Address `json:"mailingAddress" db:"mailing_address"`

	// Work address
	WorkAddress *services.Address `json:"workAddress" db:"work_address"`

	// Residency
	Residency *services.Residency `json:"residency" db:"residency"`

	// List of citizenships
	CitizenshipCountries *services.StringArray `json:"citizenshipCountries" db:"citizenship_countries"`

	// Consumer's occupation
	Occupation *string `json:"occupation" db:"occupation"`

	// Consumer income source e.g. salary or inheritance
	IncomeType *services.StringArray `json:"incomeType" db:"income_type"`

	// Consumer activity type
	ActivityType *services.StringArray `json:"activityType" db:"activity_type"`
}

type ConsumerVerificationUpdate struct {
	// Consumer id
	ID shared.ConsumerID `json:"id" db:"id"`

	// KYC Status
	KYCStatus services.KYCStatus `json:"kycStatus" db:"kyc_status"`
}

type ConsumerKYCResponse struct {
	Status      services.KYCStatus `json:"status"`
	ReviewItems []string           `json:"reviewItems"`
	ConsumerID  *shared.ConsumerID `json:"consumerId"`
}

type ConsumerKYCError struct {
	RawError   error              `json:"-"`
	ErrorType  string             `json:"errorType"`
	Values     []string           `json:"values"`
	ConsumerID *shared.ConsumerID `json:"consumerId"`
}

func (e ConsumerKYCError) Error() string {
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

func (c ConsumerUpdate) GetFirstName() string {
	if c.FirstName != nil {
		return *c.FirstName
	}

	return ""
}

func (c ConsumerUpdate) GetFullName() string {
	if c.FirstName != nil && c.LastName != nil {
		return *c.FirstName + " " + *c.LastName
	}

	return ""
}
