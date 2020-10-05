package bbva

import (
	"github.com/wiseco/core-platform/partner/bank"
)

type ContactRequest struct {
	Type  ContactType `json:"type"`
	Value string      `json:"value"`
}

type ContactUpdateRequestValue struct {
	Value string `json:"value"`
}

type ContactUpdateRequest struct {
	Contact ContactUpdateRequestValue `json:"contact"`
}

type ContactEntityResponse struct {
	ID      string      `json:"id"`
	Type    ContactType `json:"type"`
	Contact string      `json:"contact"`
	Value   string      `json:"value"`
}

type ContactIDResponse struct {
	Contact ContactEntityResponse `json:"contact"`
}

type ContactType string

const (
	ContactTypeEmail        = ContactType("email")
	ContactTypePhone        = ContactType("phone")
	ContactTypeMobileNumber = ContactType("mobile_number")
)

type AddressType string

const (
	// Translation helpers
	AddressTypeEmpty = AddressType("")

	// Normal values
	AddressTypeLegal       = AddressType("legal")
	AddressTypePostal      = AddressType("postal")
	AddressTypeWork        = AddressType("work")
	AddressTypeMailing     = AddressType("mailing")
	AddressTypeHeadquarter = AddressType("headquarter")
	AddressTypeBilling     = AddressType("billing")
)

// AddressRequest should follow USPS normalization practices (Example: "St" instead of "street",
// common unit designator "APT" instead of "apartment"). Periods "." are not allowed. Zip+4 is not required.
type AddressRequest struct {
	// One of a standard set of values that indicate the customer's address type.
	// POSSIBLE VALUES:
	// legal: Legal address
	// postal: Postal address
	// work: Work address

	Type AddressType `json:"type,omitempty"`

	// Customer's Address line 1. The valid character set consists of UPPERCASE and lowercase letters,
	// numbers (i.e., integers), apostrophes, hyphens and the SPACE characters ONLY.
	Line1 string `json:"line1"`

	// Customer's Address line 2. The valid character set consists of UPPERCASE and lowercase letters,
	// numbers (i.e., integers), apostrophes, hyphens and the SPACE characters ONLY.
	Line2 string `json:"line2,omitempty"`

	// Customer's Residential City. The valid character set consists of UPPERCASE and lowercase letters,
	// numbers (i.e., integers), apostrophes, hyphens and the SPACE characters ONLY
	City string `json:"city"`

	// USA two-letter State abbreviation as defined by the USPS format.
	State string `json:"state"`

	// Customer's Address Postal Code. Format is 5 numbers for US zip codes, or 12345-1234 for zip codes with extension.
	ZipCode string `json:"zip_code"`

	// Country code 3-alpha
	CountryCode Country3Alpha `json:"country_code,omitempty"`
}

func (a *AddressRequest) addressToPartner(at *AddressType) bank.AddressRequest {
	ar := bank.AddressRequest{
		Line1:   a.Line1,
		Line2:   a.Line2,
		City:    a.City,
		State:   a.State,
		ZipCode: a.ZipCode,
		Country: partnerCountryTo[a.CountryCode],
	}

	if at != nil {
		ar.Type = bank.AddressRequestType(*at)
	} else {
		ar.Type = bank.AddressRequestType(a.Type)
	}

	return ar
}

func toPartnerAddresses(addrs []AddressRequest) []bank.AddressRequest {
	var paddrs []bank.AddressRequest
	for _, a := range addrs {
		paddrs = append(paddrs, a.addressToPartner(nil))
	}

	return paddrs
}

func addressFromPartner(a bank.AddressRequest, at *AddressType) AddressRequest {
	ar := AddressRequest{
		Type:        AddressType(a.Type),
		Line1:       stripAddressPart(a.Line1),
		Line2:       stripAddressPart(a.Line2),
		City:        stripAddressPart(a.City),
		State:       stripAddressPart(a.State),
		ZipCode:     stripAddressPart(a.ZipCode),
		CountryCode: partnerCountryFrom[a.Country],
	}

	if at != nil {
		if *at == AddressTypeBilling {
			ar.Line1 = truncateTo(ar.Line1, 20)
			ar.Line2 = truncateTo(ar.Line2, 20)
			ar.City = truncateTo(ar.City, 20)
			ar.State = truncateTo(ar.State, 20)
			ar.ZipCode = truncateTo(ar.ZipCode, 20)
		} else {
			ar.Type = *at
		}
	}

	// Pass country for billing only
	if ar.Type != AddressTypeBilling {
		ar.CountryCode = CountryEmpty
	}

	return ar
}

func fromPartnerAddresses(paddrs []bank.AddressRequest) []AddressRequest {
	var addrs []AddressRequest
	for _, pa := range paddrs {
		addrs = append(addrs, addressFromPartner(pa, nil))
	}

	return addrs
}

type AddressCreateRequest struct {
	Address AddressRequest `json:"address"`
}

type AddressResponse AddressRequest

func (a *AddressResponse) addressToPartner(at *AddressType) bank.AddressResponse {
	ar := bank.AddressResponse{
		Line1:   a.Line1,
		Line2:   a.Line2,
		City:    a.City,
		State:   a.State,
		ZipCode: a.ZipCode,
		Country: partnerCountryTo[Country3Alpha(a.CountryCode)],
	}

	if at != nil {
		ar.Type = bank.AddressRequestType(*at)
	} else {
		ar.Type = bank.AddressRequestType(a.Type)
	}

	return ar
}

type AddressEntityResponse struct {
	ID string `json:"id"`
	AddressResponse
}

type AddressIDResponse struct {
	Address AddressEntityResponse `json:"address"`
}

type AddressCreateResponse struct {
	AddressID string `json:"address_id"`
}
