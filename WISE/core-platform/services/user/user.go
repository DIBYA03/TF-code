/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

// Package for all user related services
package user

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/wiseco/core-platform/services"
	"github.com/wiseco/core-platform/shared"
)

//UserNotification settings object controls how and when push notifications are sent
type UserNotification struct {
	// Send notification on transactions
	Transactions *bool `json:"transactions"`

	// Send notification on transfers
	Transfers *bool `json:"transfers"`

	// Send notification on contact changes
	Contacts *bool `json:"contacts"`
}

// SQL value marshaller
func (n UserNotification) Value() (driver.Value, error) {
	return json.Marshal(n)
}

// SQL scan unmarshaller
func (n *UserNotification) Scan(src interface{}) error {
	source, ok := src.([]byte)
	if !ok {
		return errors.New("Type convertible to []byte")
	}

	var out UserNotification
	err := json.Unmarshal(source, &out)
	if err != nil {
		return err
	}

	*n = out
	return nil
}

type UserNotificationUpdate UserNotification

func (u UserNotificationUpdate) Value() (driver.Value, error) {
	return json.Marshal(u)
}

// User entity object provides access to accounts, businesses, and cards
type User struct {

	// User id (uuid)
	ID shared.UserID `json:"id" db:"id"`

	// ConsumerID is shared by users and members
	ConsumerID shared.ConsumerID `json:"consumerId" db:"consumer_id"`

	// AuthID from auth server
	IdentityID shared.IdentityID `json:"-" db:"identity_id"`

	// Partner id
	PartnerID *shared.PartnerID `json:"partnerId" db:"partner_id"`

	// Email
	Email *string `json:"email" db:"email"`

	// Email verified
	EmailVerified bool `json:"emailVerified" db:"email_verified"`

	// Phone
	Phone string `json:"phone" db:"phone"`

	// Phone verified
	PhoneVerified bool `json:"phoneVerified" db:"phone_verified"`

	// Notification settings object
	Notification UserNotification `json:"notification" db:"notification"`

	// Deactivated
	Deactivated *time.Time `json:"deactivated" db:"deactivated"`

	// Created time
	Created time.Time `json:"created" db:"created"`

	// Last modified time
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

	// Date of birth
	DateOfBirth *shared.Date `json:"dateOfBirth" db:"date_of_birth"`

	// Tax id number
	TaxID *services.TaxID `json:"taxIdMasked" db:"tax_id"`

	// Tax id type e.g. SSN or ITIN
	TaxIDType *services.TaxIDType `json:"taxIdType" db:"tax_id_type"`

	// KYC Status
	KYCStatus string `json:"kycStatus" db:"kyc_status"`

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

	// Subscription
	SubscriptionStatus *services.SubscriptionStatus `json:"subscriptionStatus" db:"subscription_status"`
}

func (u *User) Name() string {
	return strings.TrimSpace(u.FirstName + " " + u.LastName)
}

// UserAuthCreate is a convenniece struct to pass on to `user_service`
// so we can create a minimal user on `DB`
type UserAuthCreate struct {
	// AuthID from auth server
	IdentityID shared.IdentityID `json:"-" db:"identity_id"`

	// Phone number
	Phone string `db:"phone"`

	// phone number verified
	PhoneVerified bool `db:"phone_verified"`

	// Email
	Email *string `json:"email" db:"email"`

	// Email verified
	EmailVerified bool `json:"emailVerified" db:"email_verified"`
}

type UserCreate struct {
	UserAuthCreate

	/*
	 * ConsumerCreate
	 */

	// First name
	FirstName string `json:"firstName" db:"first_name"`

	// Middle name
	MiddleName string `json:"middleName" db:"middle_name"`

	// Last name
	LastName string `json:"lastName" db:"last_name"`

	// Date of birth
	DateOfBirth *shared.Date `json:"dateOfBirth" db:"date_of_birth"`

	// Tax id number
	TaxID *services.TaxID `json:"taxId" db:"tax_id"`

	// Tax id type e.g. SSN or ITIN
	TaxIDType *services.TaxIDType `json:"taxIdType" db:"tax_id_type"`

	// Legal address
	LegalAddress *services.Address `json:"legalAddress" db:"legal_address"`

	// Mailing address
	MailingAddress *services.Address `json:"mailingAddress db:"mailing_address"`

	// Work address
	WorkAddress *services.Address `json:"workAddress" db:"work_address"`

	// Residency
	Residency *services.Residency `json:"residency" db:"residency"`

	// List of citizenships
	CitizenshipCountries services.StringArray `json:"citizenshipCountries" db:"citizenship_countries"`

	// User's occupation
	Occupation *string `json:"occupation" db:"occupation"`

	// User income source e.g. salary or inheritance
	IncomeType *services.StringArray `json:"incomeType" db:"income_type"`

	// User activity type
	ActivityType *services.StringArray `json:"activityType" db:"activity_type"`
}

type UserUpdate struct {
	// User id
	ID shared.UserID `json:"userId" db:"id"`

	// First name
	FirstName *string `json:"firstName" db:"first_name"`

	// Middle name
	MiddleName *string `json:"middleName" db:"middle_name"`

	// Last name
	LastName *string `json:"lastName" db:"last_name"`

	// Email
	Email *string `json:"email" db:"email"`

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

	// User's occupation
	Occupation *string `json:"occupation" db:"occupation"`

	// User income source e.g. salary or inheritance
	IncomeType *services.StringArray `json:"incomeType" db:"income_type"`

	// User activity type
	ActivityType *services.StringArray `json:"activityType" db:"activity_type"`
}

type UserVerificationUpdate struct {
	// User id
	ID shared.UserID `json:"id" db:"id"`

	// KYC Status
	KYCStatus string `json:"kycStatus" db:"kyc_status"`
}

type UserPartnerUpdate struct {
	// User id
	ID shared.UserID `json:"id" db:"id"`

	// Partner id
	PartnerID shared.PartnerID `json:"partnerId" db:"partner_id"`
}

type UserKYCResponse struct {
	Status      services.KYCStatus `json:"status"`
	ReviewItems []string           `json:"reviewItems"`
	User        User               `json:"user"`
}

type UserKYCError struct {
	RawError  error          `json:"-"`
	ErrorType string         `json:"errorType"`
	Values    []string       `json:"values"`
	UserID    *shared.UserID `json:"userId"`
}

func (e UserKYCError) Error() string {
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
