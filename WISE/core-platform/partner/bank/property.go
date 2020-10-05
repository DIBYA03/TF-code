package bank

type PropertyBankID string
type ConsumerPropertyType string

const (
	ConsumerPropertyTypeContactEmail   = ConsumerPropertyType("contact.email")
	ConsumerPropertyTypeContactPhone   = ConsumerPropertyType("contact.phone")
	ConsumerPropertyTypeAddressLegal   = ConsumerPropertyType("address.legal")
	ConsumerPropertyTypeAddressMailing = ConsumerPropertyType("address.mailing")
	ConsumerPropertyTypeAddressWork    = ConsumerPropertyType("address.work")
	ConsumerPropertyTypeAddressBilling = ConsumerPropertyType("address.billing")
)

type BusinessPropertyType string

const (
	BusinessPropertyTypeContactEmail       = BusinessPropertyType("contact.email")
	BusinessPropertyTypeContactPhone       = BusinessPropertyType("contact.phone")
	BusinessPropertyTypeAddressLegal       = BusinessPropertyType("address.legal")
	BusinessPropertyTypeAddressMailing     = BusinessPropertyType("address.mailing")
	BusinessPropertyTypeAddressHeadquarter = BusinessPropertyType("address.headquarter")
	BusinessPropertyTypeAddressBilling     = BusinessPropertyType("address.billing")
)

type AddressRequestType string

const (
	AddressRequestTypeLegal       = AddressRequestType("legal")
	AddressRequestTypeMailing     = AddressRequestType("mailing")
	AddressRequestTypeWork        = AddressRequestType("work")
	AddressRequestTypeHeadquarter = AddressRequestType("headquarter")
	AddressRequestTypeBilling     = AddressRequestType("billing")
	AddressRequestTypeRemittance  = AddressRequestType("remittance")
	AddressRequestTypeOther       = AddressRequestType("other")
)

// AddressRequest should follow USPS normalization practices (Example: "St" instead of "street",
// common unit designator "APT" instead of "apartment"). Periods "." are not allowed. Zip+4 is not required.
type AddressRequest struct {
	// One of a standard set of values that indicate the customer's address type.
	// POSSIBLE VALUES:
	// legal: Legal address
	// mailing: Mailing address
	Type AddressRequestType `json:"type"`

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
	ZipCode string `json:"zipCode"`

	// Country
	Country Country `json:"country,omitempty"`
}

type AddressResponse AddressRequest
