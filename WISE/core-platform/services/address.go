/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

package services

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"strings"

	partnerbank "github.com/wiseco/core-platform/partner/bank"
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

// Address refers to a physical or mailing address
type Address struct {
	// AddressType refers to legal or mailing
	Type AddressType `json:"addressType,omitempty"`

	// StreetAddress
	StreetAddress string `json:"streetAddress"`

	// AddressLine2 (unit/suite #)
	AddressLine2 string `json:"addressLine2,omitempty"`

	// City or town
	City string `json:"city"`

	// State or province
	State string `json:"state"`

	// Country ISO 3166 2-Alpha
	Country string `json:"country"`

	// PostalCode postal code of address
	PostalCode string `json:"postalCode"`

	// Latitude of address
	Latitude float64 `json:"latitude,omitempty"`

	// Longitude of address
	Longitude float64 `json:"longitude,omitempty"`
}

// SQL value marshaller
func (a Address) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// SQL scan unmarshaller
func (addr *Address) Scan(src interface{}) error {
	source, ok := src.([]byte)
	if !ok {
		return errors.New("Type convertible to []byte")
	}

	var out Address
	err := json.Unmarshal(source, &out)
	if err != nil {
		return err
	}

	*addr = out
	return nil
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
		City:    a.City,
		State:   a.State,
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
	addr.City = strings.TrimSpace(a.City)
	addr.State = strings.TrimSpace(a.State)
	addr.PostalCode = strings.TrimSpace(a.PostalCode)
	addr.Country = strings.TrimSpace(a.Country)
	addr.Latitude = a.Latitude
	addr.Longitude = a.Longitude

	// Replace with call to address validator
	if len(addr.StreetAddress) < 2 {
		return addr, errors.New("invalid street address")
	}

	if len(addr.City) < 2 {
		return addr, errors.New("invalid city")
	}

	if len(addr.State) < 2 {
		return addr, errors.New("invalid state or provence")
	}

	if len(addr.PostalCode) < 2 {
		return addr, errors.New("invalid postal code")
	}

	return addr, nil
}
