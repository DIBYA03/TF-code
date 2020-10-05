package bank

type BankPayeeID string

//TODO nothing about this will scale to many banks
//rethink this when creating a banking service
type PayeeAddressRequest struct {
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
}

type LinkedPayeeRequest struct {
	PayeeName         string               `json:"payee_name"`
	ACIID             string               `json:"payee_id,omitempty"` //refers to an optional id provided by https://www.aciworldwide.com/
	AccountNumber     string               `json:"account_number"`
	NameOnAccount     string               `json:"name_on_account"`
	RemittanceAddress *PayeeAddressRequest `json:"remittance_address"`
}

type LinkedPayeeResponse struct {
	AccountLast4 string      `json:"account_last4"`
	BankPayeeID  BankPayeeID `json:"account_reference_id"`
	ACIID        string      `json:"payee_id"` //refers to an optional id provided by https://www.aciworldwide.com/
	PayeeName    string      `json:"payee_name"`
}
