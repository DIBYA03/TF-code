/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

package address

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"strings"
	"time"

	partnerbank "github.com/wiseco/core-platform/partner/bank"
	"github.com/wiseco/core-platform/shared"
)

type AddressType string

const AddressTypeNone = AddressType("")

const (
	// AddressTypeLegal is a legal business or residential address
	AddressTypeLegal = AddressType("legal")

	// AddressTypeMailing is a postal service compatible mailing address
	AddressTypeMailing = AddressType("mailing")

	// AddressTypeMailing is the work address for a user
	AddressTypeWork = AddressType("work")

	// AddressTypeMailing is the businesses headquarter address
	AddressTypeHeadquarter = AddressType("headquarter")

	// AddressTypeBilling is a postal service compatible billing address
	AddressTypeBilling = AddressType("billing")

	// Other address type
	AddressTypeOther = AddressType("other")
)

type AddressState string

const (
	AddressStateInactive = "inactive"
	AddressStateActive   = "active"
)

// Address refers to a physical or mailing address
type Address struct {
	// Address id
	ID shared.AddressID `json:"id" db:"id"`

	// Consumer ID reference
	ConsumerID shared.ConsumerID `json:"consumerID" db:"consumer_id"`

	// Contact ID reference
	ContactID shared.ContactID `json:"contactID" db:"contact_id"`

	// Business ID reference
	BusinessID shared.BusinessID `json:"businessID" db:"business_id"`

	// StreetAddress
	StreetAddress string `json:"streetAddress" db:"street"`

	// AddressLine2 (unit/suite #)
	AddressLine2 string `json:"addressLine2,omitempty" db:"line2"`

	// locality or town
	Locality string `json:"locality" db:"locality"`

	// State or province
	AdminArea string `json:"adminArea" db:"admin_area"`

	// Country ISO 3166 2-Alpha
	Country string `json:"country" db:"country"`

	// PostalCode postal code of address
	PostalCode string `json:"postalCode" db:"postal_code"`

	// Latitude of address
	Latitude float64 `json:"latitude,omitempty" db:"latitude"`

	// Longitude of address
	Longitude float64 `json:"longitude,omitempty" db:"longitude"`

	// AddressType refers to legal or mailing
	Type AddressType `json:"addressType,omitempty" db:"address_type"`

	// The state this record is in active/inactive
	AddressState AddressState `json:"addressState,omitempty" db:"address_state"`

	// Created timestamp
	Created time.Time `json:"created" db:"created"`

	// Modified timestamp
	Modified time.Time `json:"modified" db:"modified"`
}

type AddressCreate struct {
	// Address id
	ID shared.AddressID `json:"id" db:"id"`

	// Consumer ID reference
	ConsumerID shared.ConsumerID `json:"consumerID" db:"consumer_id"`

	// Contact ID reference
	ContactID shared.ContactID `json:"contactID" db:"contact_id"`

	// Business ID reference
	BusinessID shared.BusinessID `json:"businessID" db:"business_id"`

	// StreetAddress
	StreetAddress string `json:"streetAddress" db:"street"`

	// AddressLine2 (unit/suite #)
	AddressLine2 string `json:"addressLine2" db:"line2"`

	// locality or town
	Locality string `json:"locality" db:"locality"`

	// State or province
	AdminArea string `json:"adminArea" db:"admin_area"`

	// Country ISO 3166 2-Alpha
	Country string `json:"country" db:"country"`

	// PostalCode postal code of address
	PostalCode string `json:"postalCode" db:"postal_code"`

	// Latitude of address
	Latitude float64 `json:"latitude" db:"latitude"`

	// Longitude of address
	Longitude float64 `json:"longitude" db:"longitude"`

	// AddressType refers to legal or mailing
	Type AddressType `json:"addressType" db:"address_type"`

	// The state this record is in active/inactive
	AddressState AddressState `json:"addressState" db:"address_state"`
}

type AddressUpdate struct {
	// Address id
	ID shared.AddressID `json:"id" db:"id"`

	// StreetAddress
	StreetAddress string `json:"streetAddress" db:"street"`

	// AddressLine2 (unit/suite #)
	AddressLine2 string `json:"addressLine2" db:"line2"`

	// locality or town
	Locality string `json:"locality" db:"locality"`

	// State or province
	AdminArea string `json:"adminArea" db:"admin_area"`

	// Country ISO 3166 2-Alpha
	Country string `json:"country" db:"country"`

	// PostalCode postal code of address
	PostalCode string `json:"postalCode" db:"postal_code"`

	// Latitude of address
	Latitude float64 `json:"latitude db:"latitude`

	// Longitude of address
	Longitude float64 `json:"longitude" db:"longitude"`
}

// SQL value marshaller
func (a Address) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a Address) ToPartnerBankAddress(at AddressType) partnerbank.AddressRequest {
	addrType := a.Type
	if at != AddressTypeNone {
		addrType = at
	}

	return partnerbank.AddressRequest{
		Type:    partnerbank.AddressRequestType(addrType),
		Line1:   a.StreetAddress,
		Line2:   a.AddressLine2,
		City:    a.Locality,
		State:   a.AdminArea,
		ZipCode: a.PostalCode,
		Country: partnerbank.Country(a.Country),
	}
}

// Validate formats and returns new value
func ValidateAddress(a *Address, addrType AddressType) (*Address, error) {
	// Ignore nil address values
	if a == nil {
		return nil, nil
	}

	addr := &Address{}
	if addrType != AddressTypeNone {
		addr.Type = addrType
	} else {
		addr.Type = a.Type
	}

	addr.StreetAddress = strings.TrimSpace(a.StreetAddress)
	addr.AddressLine2 = strings.TrimSpace(a.AddressLine2)
	addr.Locality = strings.TrimSpace(a.Locality)
	addr.AdminArea = strings.TrimSpace(a.AdminArea)
	addr.PostalCode = strings.TrimSpace(a.PostalCode)
	addr.Country = strings.TrimSpace(a.Country)
	addr.Latitude = a.Latitude
	addr.Longitude = a.Longitude

	// Replace with call to address validator
	if len(addr.StreetAddress) < 2 {
		return addr, errors.New("invalid street address")
	}

	if len(addr.Locality) < 2 {
		return addr, errors.New("invalid locality")
	}

	if len(addr.AdminArea) < 2 {
		return addr, errors.New("invalid state or provence")
	}

	if len(addr.PostalCode) < 2 {
		return addr, errors.New("invalid postal code")
	}

	return addr, nil
}
