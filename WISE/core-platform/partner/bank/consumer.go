package bank

import "time"

type ConsumerTaxIDType string

const (
	ConsumerTaxIDTypeSSN  = ConsumerTaxIDType("ssn")
	ConsumerTaxIDTypeITIN = ConsumerTaxIDType("itin")
)

type ConsumerID string
type ConsumerBankID string

type CreateConsumerRequest struct {
	ConsumerID         ConsumerID                       `json:"consumerId"`
	FirstName          string                           `json:"firstName"`
	MiddleName         string                           `json:"middleName"`
	LastName           string                           `json:"lastName"`
	TaxID              string                           `json:"taxId"`
	TaxIDType          ConsumerTaxIDType                `json:"taxIdType"`
	DateOfBirth        time.Time                        `json:"dob"`
	Residency          ConsumerResidency                `json:"residency"`
	CitizenshipCountry Country                          `json:"citizenshipCountry"`
	Occupation         ConsumerOccupation               `json:"occupation"`
	Income             []Income                         `json:"income"`
	ExpectedActivity   []ExpectedActivity               `json:"expectedActivity"`
	Phone              string                           `json:"phone"`
	Email              string                           `json:"email"`
	LegalAddress       AddressRequest                   `json:"legalAddress"`
	MailingAddress     *AddressRequest                  `json:"mailingAddress"`
	WorkAddress        *AddressRequest                  `json:"workAddress"`
	Identification     []*ConsumerIdentificationRequest `json:"identityRequest"`
}

type UpdateConsumerRequest struct {
	ConsumerID         ConsumerID                       `json:"consumerId"`
	FirstName          string                           `json:"firstName"`
	MiddleName         string                           `json:"middleName"`
	LastName           string                           `json:"lastName"`
	TaxID              string                           `json:"taxId"`
	TaxIDType          ConsumerTaxIDType                `json:"taxIdType"`
	DateOfBirth        time.Time                        `json:"dob"`
	Residency          ConsumerResidency                `json:"residency"`
	CitizenshipCountry Country                          `json:"citizenshipCountry"`
	Identification     []*ConsumerIdentificationRequest `json:"identityRequest"`
}

type ConsumerIdentificationRequest struct {
	DocumentType   ConsumerIdentityDocument `json:"docType"`
	Number         string                   `json:"number"`
	Issuer         string                   `json:"issuer"`
	IssueDate      time.Time                `json:"issueDate"`
	IssueState     string                   `json:"state"`
	IssueCountry   Country                  `json:"country"`
	ExpirationDate time.Time                `json:"expirationDate"`
}

type ConsumerResidencyStatus string

const (
	ResidencyStatusCitizen     = ConsumerResidencyStatus("citizen")     // Citizen
	ResidencyStatusResident    = ConsumerResidencyStatus("resident")    // Legal non-citizen resident
	ResidencyStatusNonResident = ConsumerResidencyStatus("nonResident") // Legal temporary non-citizen resident
)

type ConsumerResidency struct {
	Country Country                 `json:"country"`
	Status  ConsumerResidencyStatus `json:"status"`
}

// ConsumerOccupation
// One of a standard set of values that indicate customer occupation.
type ConsumerOccupation string

const (
	ConsumerOccupationAgriculture                = ConsumerOccupation("agriculture")
	ConsumerOccupationClergyMinistryStaff        = ConsumerOccupation("clergyMinistryStaff")
	ConsumerOccupationConstructionIndustrial     = ConsumerOccupation("constructionIndustrial")
	ConsumerOccupationEducation                  = ConsumerOccupation("education")
	ConsumerOccupationFinanceAccountingTax       = ConsumerOccupation("financeAccountingTax")
	ConsumerOccupationFireFirstResponders        = ConsumerOccupation("fireFirstResponders")
	ConsumerOccupationHealthcare                 = ConsumerOccupation("healthcare")
	ConsumerOccupationHomemaker                  = ConsumerOccupation("homemaker")
	ConsumerOccupationLaborGeneral               = ConsumerOccupation("laborGeneral")
	ConsumerOccupationLaborSkilled               = ConsumerOccupation("laborSkilled")
	ConsumerOccupationLawEnforcementSecurity     = ConsumerOccupation("lawEnforcementSecurity")
	ConsumerOccupationLegalServices              = ConsumerOccupation("legalServices")
	ConsumerOccupationMilitary                   = ConsumerOccupation("military")
	ConsumerOccupationNotaryRegistrar            = ConsumerOccupation("notaryRegistrar")
	ConsumerOccupationPrivateInvestor            = ConsumerOccupation("privateInvestor")
	ConsumerOccupationProfessionalAdministrative = ConsumerOccupation("professionalAdministrative")
	ConsumerOccupationProfessionalManagement     = ConsumerOccupation("professionalManagement")
	ConsumerOccupationProfessionalOther          = ConsumerOccupation("professionalOther")
	ConsumerOccupationProfessionalTechnical      = ConsumerOccupation("professionalTechnical")
	ConsumerOccupationRetired                    = ConsumerOccupation("retired")
	ConsumerOccupationSales                      = ConsumerOccupation("sales")
	ConsumerOccupationSelfEmployed               = ConsumerOccupation("selfEmployed")
	ConsumerOccupationStudent                    = ConsumerOccupation("student")
	ConsumerOccupationTransportation             = ConsumerOccupation("transportation")
	ConsumerOccupationUnemployed                 = ConsumerOccupation("unemployed")
)

// Income
// One of a standard set of values that indicate the customers source of income.
type Income string

const (
	IncomeInheritance       = Income("inheritance")       // Expected income is inheritance
	IncomeSalary            = Income("salary")            // Expected income is salary.
	IncomeSaleOfCompany     = Income("saleOfCompany")     // Expected income is from sale of company.
	IncomeSaleOfProperty    = Income("saleOfProperty")    // Expected income is from sale of property.
	IncomeInvestments       = Income("investments")       // Expected income is from investments.
	IncomeLifeInsurance     = Income("lifeInsurance")     // Expected income is from life insurance.
	IncomeDivorceSettlement = Income("divorceSettlement") // Expected income is from a divorce settlement.
	IncomeOther             = Income("other")             // Other type of income.
)

type ExpectedActivity string

const (
	// Expected activity on the account will be cash.
	ExpectedActivityCash = ExpectedActivity("cash")

	// Expected activity on the account will be checks.
	ExpectedActivityCheck = ExpectedActivity("check")

	// Expected activity on the account will be from a domestic wire transfer.
	ExpectedActivityDomesticWireTransfer = ExpectedActivity("domesticWireTransfer")

	// Expected activity on the account will be from an international wire transfer.
	ExpectedActivityInternalWireTransfer = ExpectedActivity("internationalWireTransfer")

	// Expected activity on the account will be from a domestic ach payment.
	ExpectedActivityDomesticACH = ExpectedActivity("domesticACH")

	// Expected activity on the account will be from an international ach payment.
	ExpectedActivityInterntaionalACH = ExpectedActivity("internationalACH")
)

// Identity Document Type.
type ConsumerIdentityDocument string

const (
	ConsumerIdentityDocumentDriversLicense        = ConsumerIdentityDocument("driversLicense")
	ConsumerIdentityDocumentPassport              = ConsumerIdentityDocument("passport")
	ConsumerIdentityDocumentPassportCard          = ConsumerIdentityDocument("passportCard")
	ConsumerIdentityDocumentWorkPermit            = ConsumerIdentityDocument("workPermit")
	ConsumerIdentityDocumentSocialSecurityCard    = ConsumerIdentityDocument("socialSecurityCard")
	ConsumerIdentityDocumentStateID               = ConsumerIdentityDocument("stateId")
	ConsumerIdentityDocumentAlienRegistrationCard = ConsumerIdentityDocument("alienRegistrationCard")
	ConsumerIdentityDocumentUSAVisaH1B            = ConsumerIdentityDocument("usaVisaH1B")
	ConsumerIdentityDocumentUSAVisaH1C            = ConsumerIdentityDocument("usaVisaH1C")
	ConsumerIdentityDocumentUSAVisaH2A            = ConsumerIdentityDocument("usaVisaH2A")
	ConsumerIdentityDocumentUSAVisaH2B            = ConsumerIdentityDocument("usaVisaH2B")
	ConsumerIdentityDocumentUSAVisaH3             = ConsumerIdentityDocument("usaVisaH3")
	ConsumerIdentityDocumentUSAVisaL1A            = ConsumerIdentityDocument("usaVisaL1A")
	ConsumerIdentityDocumentUSAVisaL1B            = ConsumerIdentityDocument("usaVisaL1B")
	ConsumerIdentityDocumentUSAVisaO1             = ConsumerIdentityDocument("usaVisaO1")
	ConsumerIdentityDocumentUSAVisaE1             = ConsumerIdentityDocument("usaVisaE1")
	ConsumerIdentityDocumentUSAVisaE3             = ConsumerIdentityDocument("usaVisaE3")
	ConsumerIdentityDocumentUSAVisaI              = ConsumerIdentityDocument("usaVisaI")
	ConsumerIdentityDocumentUSAVisaP              = ConsumerIdentityDocument("usaVisaP")
	ConsumerIdentityDocumentUSAVisaTN             = ConsumerIdentityDocument("usaVisaTN")
	ConsumerIdentityDocumentUSAVisaTD             = ConsumerIdentityDocument("usaVisaTD")
	ConsumerIdentityDocumentUSAVisaR1             = ConsumerIdentityDocument("usaVisaR1")
)

type ConsumerDocument struct {
	DocumentType   ConsumerIdentityDocument `json:"docType"`
	Number         string                   `json:"number"`
	Issuer         string                   `json:"issuer"`
	IssueDate      *time.Time               `json:"issueDate"`
	IssueState     *string                  `json:"state"`
	IssueCountry   *Country                 `json:"country"`
	ExpirationDate *time.Time               `json:"expirationDate"`
}

type IdentityStatusConsumerResponse struct {
	ConsumerID ConsumerID     `json:"consumerID"`
	BankID     ConsumerBankID `json:"bankId"`
	KYC        KYCResponse    `json:"kyc"`
}

// ConsumerIdentityRequest
// To help verify a user identity you will need to post the following information about your customer:
// first name, middle name (optional), last name, social security number (SSN), date of birth (DOB), phone, email and address.
// The address should follow United States Postal Service normalization practices.
type ConsumerIdentityRequest struct {
	FirstName   string           `json:"firstName"`
	MiddleName  string           `json:"middleName,omitempty"`
	LastName    string           `json:"lastName"`
	SSN         string           `json:"ssn"`
	DateOfBirth string           `json:"dob"`
	Phone       string           `json:"phone"`
	Email       string           `json:"email"`
	Address     []AddressRequest `json:"address"`
}

//ConsumerDocumentType ConsumerIdentityResponse
type ConsumerIdentityResponse struct {
	ConsumerID ConsumerID     `json:"consumerID"`
	BankID     ConsumerBankID `json:"bankId"`
	KYC        KYCResponse    `json:"kyc"`
	KYCNotes   []string       `json:"kycNotes"`
}

type ConsumerResponse struct {
	ConsumerID ConsumerID     `json:"consumerID"`
	BankID     ConsumerBankID `json:"bankId"`
	KYCStatus  KYCStatus      `json:"kycStatus"`
	Created    time.Time      `json:"created"`
	Updated    time.Time      `json:"updated"`
}
