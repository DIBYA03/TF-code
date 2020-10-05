package bank

type IDVerify string

func (v IDVerify) String() string {
	return string(v)
}

const (
	// Postal mailing address
	IDVerifyAddress = IDVerify("address")

	// Tax id
	IDVerifyTaxId = IDVerify("taxId")

	// Full name
	IDVerifyFullName = IDVerify("fullName")

	// Date of birth
	IDVerifyDateOfBirth = IDVerify("dob")

	// Any category that isn't covered by the preceding categories.
	IDVerifyOther = IDVerify("other")

	// United States Office of Foreign Asset Control
	IDVerifyOFAC = IDVerify("usofac")

	// Indicates that our bank partner has an existing consumer record that
	// contains the submitted SSN, but that one or more of the submitted values for
	// name, dob, or citizenship_status do not match the information contained in
	// that record. Review the submitted values of name, dob, or citizenship for
	// errors. If there are no errors, contact our client integration team for
	// assistance updating the existing record.
	IDVerifyMismatch = IDVerify("mismatch")

	// Primary consumer doc
	IDVerifyPrimaryDoc = IDVerify("primaryDoc")

	// Secondary consumer doc
	IDVerifySecondaryDoc = IDVerify("secondaryDoc")

	// Business formation document
	IDVerifyFormationDoc = IDVerify("formationDoc")
)

type KYCNote struct {
	Code string `json:"code"`
	Desc string `json:"desc"`
}

// KYCResponse contains a Know Your Customer (KYC) status that your application
// must process.
type KYCResponse struct {
	// Status status of the customer's application.
	Status KYCStatus `json:"status"`

	// An array containing information required to be verified. Valis when kyc
	// status is in review
	IDVerifyRequired []IDVerify `json:"idVerifyRequired"`

	// KYC Notes
	Notes []KYCNote `json:"notes"`
}

type KYCStatus string

func (s KYCStatus) String() string {
	return string(s)
}

const (
	// Default status
	KYCStatusNotStarted KYCStatus = "notStarted"

	// The new customer record can be created using the information supplied. This
	// is the default status for Open Platform sandbox requests.
	KYCStatusApproved KYCStatus = "approved"

	// Customer must submit additional identification before the consumer record can
	// be created. In the row for cip, see IDVerify.
	KYCStatusReview KYCStatus = "review"

	// The submitted identity information cannot be used to create a consumer record.
	KYCStatusDeclined KYCStatus = "declined"
)

type KYCRisk string

const (
	KYCRiskLow    = KYCRisk("low")
	KYCRiskMedium = KYCRisk("medium")
	KYCRiskHigh   = KYCRisk("high")
)

type KYCResult string

const (
	KYCResultPass = KYCResult("pass")
	KYCResultFail = KYCResult("fail")
)
