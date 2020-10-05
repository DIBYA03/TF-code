package bbva

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"

	"github.com/jmoiron/sqlx/types"
	"github.com/wiseco/core-platform/partner/bank"
	"github.com/wiseco/core-platform/partner/bank/data"
)

type consumerService struct {
	request bank.APIRequest
	client  *client
}

func (b *consumerBank) ConsumerEntityService(request bank.APIRequest) bank.ConsumerService {
	return &consumerService{
		request: request,
		client:  b.client,
	}
}

type UserID string

type CreateConsumerRequest struct {
	UserID             UserID                           `json:"-"`
	FirstName          string                           `json:"first_name"`
	MiddleName         string                           `json:"middle_name,omitempty"`
	LastName           string                           `json:"last_name"`
	SSN                string                           `json:"ssn"`
	DateOfBirth        string                           `json:"dob"`
	Contact            []ContactRequest                 `json:"contact"`
	USResidencyStatus  USResidencyStatus                `json:"citizenship_status"`
	CitizenshipCountry Country3Alpha                    `json:"citizenship_country"`
	Occupation         ConsumerOccupation               `json:"occupation"`
	Income             []Income                         `json:"income"`
	ExpectedActivity   []ExpectedActivity               `json:"expected_activity"`
	Address            []AddressRequest                 `json:"address"`
	Identification     []*ConsumerIdentificationRequest `json:"identification,omitempty"`
	Pep                ConsumerPepRequest               `json:"pep"`
}

type UpdateConsumerRequest struct {
	UserID             UserID                           `json:"-"`
	FirstName          string                           `json:"first_name,omitempty"`
	MiddleName         string                           `json:"middle_name,omitempty"`
	LastName           string                           `json:"last_name,omitempty"`
	SSN                string                           `json:"ssn,omitempty"`
	DateOfBirth        string                           `json:"dob,omitempty"`
	Contact            []ContactRequest                 `json:"contact,omitempty"`
	USResidencyStatus  USResidencyStatus                `json:"citizenship_status,omitempty"`
	CitizenshipCountry Country3Alpha                    `json:"citizenship_country,omitempty"`
	Identification     []*ConsumerIdentificationRequest `json:"identification,omitempty"`
}

func newCreateConsumerRequest(preq bank.CreateConsumerRequest) (*CreateConsumerRequest, error) {
	country, ok := partnerCountryFrom[preq.CitizenshipCountry]
	if !ok {
		return nil, bank.NewErrorFromCode(bank.ErrorCodeInvalidCountryCode)
	}

	at := AddressTypeLegal
	var addresses = []AddressRequest{addressFromPartner(preq.LegalAddress, &at)}

	if preq.MailingAddress != nil {
		at = AddressTypePostal
		addresses = append(addresses, addressFromPartner(*preq.MailingAddress, &at))
	}

	if preq.WorkAddress != nil {
		at = AddressTypeWork
		addresses = append(addresses, addressFromPartner(*preq.WorkAddress, &at))
	}

	residency, err := partnerResidencyFrom(preq.Residency)
	if err != nil {
		return nil, err
	}

	occupation, ok := occupationFromPartnerMap[preq.Occupation]
	if !ok {
		return nil, bank.NewErrorFromCode(bank.ErrorCodeInvalidOccupation)
	}

	income, err := incomeFromPartner(preq.Income)
	if err != nil {
		return nil, err
	}

	activities, err := activityFromPartner(preq.ExpectedActivity)
	if err != nil {
		return nil, err
	}

	var idr []*ConsumerIdentificationRequest
	for _, idReq := range preq.Identification {
		docType, ok := partnerConsumerIdentityDocumentFrom[idReq.DocumentType]
		if !ok {
			return nil, bank.NewErrorFromCode(bank.ErrorCodeInvalidDocumentType)
		}

		country, ok := partnerCountryFrom[idReq.IssueCountry]
		if !ok {
			return nil, bank.NewErrorFromCode(bank.ErrorCodeInvalidCountry)
		}

		if idReq.IssueDate.IsZero() {
			return nil, errors.New("issue date is not set")
		}

		if idReq.ExpirationDate.IsZero() {
			return nil, errors.New("expiration date is not set")
		}

		if !idReq.IssueDate.Before(idReq.ExpirationDate) {
			return nil, errors.New("issue date must be before expiration")
		}

		idr = append(idr, &ConsumerIdentificationRequest{
			Document:       docType,
			Number:         idReq.Number,
			IssuingState:   idReq.IssueState,
			IssuingCountry: country,
			IssuedDate:     idReq.IssueDate.Format("2006-01-02"),
			ExpirationDate: idReq.ExpirationDate.Format("2006-01-02"),
		})
	}

	return &CreateConsumerRequest{
		UserID:      UserID(preq.ConsumerID),
		FirstName:   stripConsumerName(preq.FirstName),
		MiddleName:  stripConsumerName(preq.MiddleName),
		LastName:    stripConsumerName(preq.LastName),
		SSN:         stripTaxID(preq.TaxID),
		DateOfBirth: preq.DateOfBirth.Format("2006-01-02"),
		Contact: []ContactRequest{
			ContactRequest{
				Type:  ContactTypePhone,
				Value: preq.Phone,
			},
			ContactRequest{
				Type:  ContactTypeEmail,
				Value: strings.ToLower(preq.Email),
			},
		}, // Map from req.Contact
		USResidencyStatus:  residency,
		CitizenshipCountry: country,
		Occupation:         occupation,
		Income:             income,
		ExpectedActivity:   activities,
		Address:            addresses,
		Identification:     idr,

		// We don't currently accept PEP
		Pep: ConsumerPepRequest{Association: PepAssociationNo},
	}, nil
}

func newUpdateConsumerRequest(c *data.Consumer, preq bank.UpdateConsumerRequest) (*UpdateConsumerRequest, error) {
	country, ok := partnerCountryFrom[preq.CitizenshipCountry]
	if !ok {
		return nil, bank.NewErrorFromCode(bank.ErrorCodeInvalidCountryCode)
	}

	if c.KYCStatus == bank.KYCStatusApproved {
		return nil, errors.New("User has already been approved")
	} else if c.KYCStatus == bank.KYCStatusDeclined {
		return nil, errors.New("User has already been declined")
	} else if c.KYCStatus == bank.KYCStatusReview {
		return &UpdateConsumerRequest{
			UserID:      UserID(c.BankID),
			FirstName:   preq.FirstName,
			MiddleName:  preq.MiddleName,
			LastName:    preq.LastName,
			SSN:         preq.TaxID,
			DateOfBirth: preq.DateOfBirth.Format("2006-01-02"),
		}, nil
	}

	residency, err := partnerResidencyFrom(preq.Residency)
	if err != nil {
		return nil, err
	}

	var idr []*ConsumerIdentificationRequest
	for _, idReq := range preq.Identification {
		docType, ok := partnerConsumerIdentityDocumentFrom[idReq.DocumentType]
		if !ok {
			return nil, bank.NewErrorFromCode(bank.ErrorCodeInvalidDocumentType)
		}

		country, ok := partnerCountryFrom[idReq.IssueCountry]
		if !ok {
			return nil, bank.NewErrorFromCode(bank.ErrorCodeInvalidCountry)
		}

		if idReq.IssueDate.IsZero() {
			return nil, errors.New("issue date is not set")
		}

		if idReq.ExpirationDate.IsZero() {
			return nil, errors.New("expiration date is not set")
		}

		if !idReq.IssueDate.Before(idReq.ExpirationDate) {
			return nil, errors.New("issue date must be before expiration")
		}

		idr = append(idr, &ConsumerIdentificationRequest{
			Document:       docType,
			Number:         idReq.Number,
			IssuingState:   idReq.IssueState,
			IssuingCountry: country,
			IssuedDate:     idReq.IssueDate.Format("2006-01-02"),
			ExpirationDate: idReq.ExpirationDate.Format("2006-01-02"),
		})
	}

	return &UpdateConsumerRequest{
		UserID:             UserID(c.BankID),
		FirstName:          preq.FirstName,
		MiddleName:         preq.MiddleName,
		LastName:           preq.LastName,
		SSN:                preq.TaxID,
		DateOfBirth:        preq.DateOfBirth.Format("2006-01-02"),
		USResidencyStatus:  residency,
		CitizenshipCountry: country,
		Identification:     idr,
	}, nil
}

// ConsumerOccupation
// One of a standard set of values that indicate customer occupation.
type ConsumerOccupation string

const (
	ConsumerOccupationAgriculture                = ConsumerOccupation("agriculture")
	ConsumerOccupationClergyMinistryStaff        = ConsumerOccupation("clergy_ministry_staff")
	ConsumerOccupationConstructionIndustrial     = ConsumerOccupation("construction_industrial")
	ConsumerOccupationEducation                  = ConsumerOccupation("education")
	ConsumerOccupationFinanceAccountingTax       = ConsumerOccupation("finance_accounting_tax")
	ConsumerOccupationFireFirstResponders        = ConsumerOccupation("fire_first_responders")
	ConsumerOccupationHealthcare                 = ConsumerOccupation("healthcare")
	ConsumerOccupationHomemaker                  = ConsumerOccupation("homemaker")
	ConsumerOccupationLaborGeneral               = ConsumerOccupation("labor_general")
	ConsumerOccupationLaborSkilled               = ConsumerOccupation("labor_skilled")
	ConsumerOccupationLawEnforcementSecurity     = ConsumerOccupation("law_enforcement_security")
	ConsumerOccupationLegalServices              = ConsumerOccupation("legal_services")
	ConsumerOccupationMilitary                   = ConsumerOccupation("military")
	ConsumerOccupationNotaryRegistrar            = ConsumerOccupation("notary_registrar")
	ConsumerOccupationPrivateInvestor            = ConsumerOccupation("private_investor")
	ConsumerOccupationProfessionalAdministrative = ConsumerOccupation("professional_administrative")
	ConsumerOccupationProfessionalManagement     = ConsumerOccupation("professional_management")
	ConsumerOccupationProfessionalOther          = ConsumerOccupation("professional_other")
	ConsumerOccupationProfessionalTechnical      = ConsumerOccupation("professional_technical")
	ConsumerOccupationRetired                    = ConsumerOccupation("retired")
	ConsumerOccupationSales                      = ConsumerOccupation("sales")
	ConsumerOccupationSelfEmployed               = ConsumerOccupation("self_employed")
	ConsumerOccupationStudent                    = ConsumerOccupation("student")
	ConsumerOccupationTransportation             = ConsumerOccupation("transportation")
	ConsumerOccupationUnemployed                 = ConsumerOccupation("unemployed")
)

var occupationFromPartnerMap = map[bank.ConsumerOccupation]ConsumerOccupation{
	bank.ConsumerOccupationAgriculture:                ConsumerOccupationAgriculture,
	bank.ConsumerOccupationClergyMinistryStaff:        ConsumerOccupationClergyMinistryStaff,
	bank.ConsumerOccupationConstructionIndustrial:     ConsumerOccupationConstructionIndustrial,
	bank.ConsumerOccupationEducation:                  ConsumerOccupationEducation,
	bank.ConsumerOccupationFinanceAccountingTax:       ConsumerOccupationFinanceAccountingTax,
	bank.ConsumerOccupationFireFirstResponders:        ConsumerOccupationFireFirstResponders,
	bank.ConsumerOccupationHealthcare:                 ConsumerOccupationHealthcare,
	bank.ConsumerOccupationHomemaker:                  ConsumerOccupationHomemaker,
	bank.ConsumerOccupationLaborGeneral:               ConsumerOccupationLaborGeneral,
	bank.ConsumerOccupationLaborSkilled:               ConsumerOccupationLaborSkilled,
	bank.ConsumerOccupationLawEnforcementSecurity:     ConsumerOccupationLawEnforcementSecurity,
	bank.ConsumerOccupationLegalServices:              ConsumerOccupationLegalServices,
	bank.ConsumerOccupationMilitary:                   ConsumerOccupationMilitary,
	bank.ConsumerOccupationNotaryRegistrar:            ConsumerOccupationNotaryRegistrar,
	bank.ConsumerOccupationPrivateInvestor:            ConsumerOccupationPrivateInvestor,
	bank.ConsumerOccupationProfessionalAdministrative: ConsumerOccupationProfessionalAdministrative,
	bank.ConsumerOccupationProfessionalManagement:     ConsumerOccupationProfessionalManagement,
	bank.ConsumerOccupationProfessionalOther:          ConsumerOccupationProfessionalOther,
	bank.ConsumerOccupationProfessionalTechnical:      ConsumerOccupationProfessionalTechnical,
	bank.ConsumerOccupationRetired:                    ConsumerOccupationRetired,
	bank.ConsumerOccupationSales:                      ConsumerOccupationSales,
	bank.ConsumerOccupationSelfEmployed:               ConsumerOccupationSelfEmployed,
	bank.ConsumerOccupationStudent:                    ConsumerOccupationStudent,
	bank.ConsumerOccupationTransportation:             ConsumerOccupationTransportation,
	bank.ConsumerOccupationUnemployed:                 ConsumerOccupationUnemployed,
}

var occupationToPartnerMap = map[ConsumerOccupation]bank.ConsumerOccupation{
	ConsumerOccupationAgriculture:                bank.ConsumerOccupationAgriculture,
	ConsumerOccupationClergyMinistryStaff:        bank.ConsumerOccupationClergyMinistryStaff,
	ConsumerOccupationConstructionIndustrial:     bank.ConsumerOccupationConstructionIndustrial,
	ConsumerOccupationEducation:                  bank.ConsumerOccupationEducation,
	ConsumerOccupationFinanceAccountingTax:       bank.ConsumerOccupationFinanceAccountingTax,
	ConsumerOccupationFireFirstResponders:        bank.ConsumerOccupationFireFirstResponders,
	ConsumerOccupationHealthcare:                 bank.ConsumerOccupationHealthcare,
	ConsumerOccupationHomemaker:                  bank.ConsumerOccupationHomemaker,
	ConsumerOccupationLaborGeneral:               bank.ConsumerOccupationLaborGeneral,
	ConsumerOccupationLaborSkilled:               bank.ConsumerOccupationLaborSkilled,
	ConsumerOccupationLawEnforcementSecurity:     bank.ConsumerOccupationLawEnforcementSecurity,
	ConsumerOccupationLegalServices:              bank.ConsumerOccupationLegalServices,
	ConsumerOccupationMilitary:                   bank.ConsumerOccupationMilitary,
	ConsumerOccupationNotaryRegistrar:            bank.ConsumerOccupationNotaryRegistrar,
	ConsumerOccupationPrivateInvestor:            bank.ConsumerOccupationPrivateInvestor,
	ConsumerOccupationProfessionalAdministrative: bank.ConsumerOccupationProfessionalAdministrative,
	ConsumerOccupationProfessionalManagement:     bank.ConsumerOccupationProfessionalManagement,
	ConsumerOccupationProfessionalOther:          bank.ConsumerOccupationProfessionalOther,
	ConsumerOccupationProfessionalTechnical:      bank.ConsumerOccupationProfessionalTechnical,
	ConsumerOccupationRetired:                    bank.ConsumerOccupationRetired,
	ConsumerOccupationSales:                      bank.ConsumerOccupationSales,
	ConsumerOccupationSelfEmployed:               bank.ConsumerOccupationSelfEmployed,
	ConsumerOccupationStudent:                    bank.ConsumerOccupationStudent,
	ConsumerOccupationTransportation:             bank.ConsumerOccupationTransportation,
	ConsumerOccupationUnemployed:                 bank.ConsumerOccupationUnemployed,
}

// Income
// One of a standard set of values that indicate the customers source of income.
type Income string

const (
	IncomeInheritance       Income = "inheritance"        // Expected income is inheritance
	IncomeSalary            Income = "salary"             // Expected income is salary.
	IncomeSaleOfCompany     Income = "company_sale"       // Expected income is from sale of company.
	IncomeSaleOfProperty    Income = "property_sale"      // Expected income is from sale of property.
	IncomeInvestments       Income = "investments"        // Expected income is from investments.
	IncomeLifeInsurance     Income = "life_insurance"     // Expected income is from life insurance.
	IncomeDivorceSettlement Income = "divorce_settlement" // Expected income is from a divorce settlement.
	IncomeOther             Income = "other"              // Other type of income.
)

var incomeToPartnerMap = map[Income]bank.Income{
	IncomeInheritance:       bank.IncomeInheritance,
	IncomeSalary:            bank.IncomeSalary,
	IncomeSaleOfCompany:     bank.IncomeSaleOfCompany,
	IncomeSaleOfProperty:    bank.IncomeSaleOfProperty,
	IncomeInvestments:       bank.IncomeInvestments,
	IncomeLifeInsurance:     bank.IncomeLifeInsurance,
	IncomeDivorceSettlement: bank.IncomeDivorceSettlement,
	IncomeOther:             bank.IncomeOther,
}

var incomeFromPartnerMap = map[bank.Income]Income{
	bank.IncomeInheritance:       IncomeInheritance,
	bank.IncomeSalary:            IncomeSalary,
	bank.IncomeSaleOfCompany:     IncomeSaleOfCompany,
	bank.IncomeSaleOfProperty:    IncomeSaleOfProperty,
	bank.IncomeInvestments:       IncomeInvestments,
	bank.IncomeLifeInsurance:     IncomeLifeInsurance,
	bank.IncomeDivorceSettlement: IncomeDivorceSettlement,
	bank.IncomeOther:             IncomeOther,
}

func incomeFromPartner(pinc []bank.Income) ([]Income, error) {
	var inc []Income
	for _, p := range pinc {
		income, ok := incomeFromPartnerMap[p]
		if !ok {
			return nil, bank.NewErrorFromCode(bank.ErrorCodeInvalidIncome)
		}

		inc = append(inc, income)
	}

	return inc, nil
}

type ExpectedActivity string

const (
	// Expected activity on the account will be cash
	ExpectedActivityCash ExpectedActivity = "cash"

	// Expected activity on the account will be checks
	ExpectedActivityCheck ExpectedActivity = "check"

	// Expected activity on the account will be from a domestic wire transfer
	ExpectedActivityDomesticWireTransfer ExpectedActivity = "domestic_wire_transfer"

	// Expected activity on the account will be from an international wire transfer
	ExpectedActivityInternalWireTransfer ExpectedActivity = "international_wire_transfer"

	// Expected activity on the account will be from a domestic ach payment
	ExpectedActivityDomesticACH ExpectedActivity = "domestic_ach"

	// Expected activity on the account will be from an international ach payment
	ExpectedActivityInterntaionalACH ExpectedActivity = "international_ach"
)

var activityFromPartnerMap = map[bank.ExpectedActivity]ExpectedActivity{
	bank.ExpectedActivityCash:                 ExpectedActivityCash,
	bank.ExpectedActivityCheck:                ExpectedActivityCheck,
	bank.ExpectedActivityDomesticWireTransfer: ExpectedActivityDomesticWireTransfer,
	bank.ExpectedActivityInternalWireTransfer: ExpectedActivityInternalWireTransfer,
	bank.ExpectedActivityDomesticACH:          ExpectedActivityDomesticACH,
	bank.ExpectedActivityInterntaionalACH:     ExpectedActivityInterntaionalACH,
}

var partnerActivityTo = map[ExpectedActivity]bank.ExpectedActivity{
	ExpectedActivityCash:                 bank.ExpectedActivityCash,
	ExpectedActivityCheck:                bank.ExpectedActivityCheck,
	ExpectedActivityDomesticWireTransfer: bank.ExpectedActivityDomesticWireTransfer,
	ExpectedActivityInternalWireTransfer: bank.ExpectedActivityInternalWireTransfer,
	ExpectedActivityDomesticACH:          bank.ExpectedActivityDomesticACH,
	ExpectedActivityInterntaionalACH:     bank.ExpectedActivityInterntaionalACH,
}

func activityFromPartner(pact []bank.ExpectedActivity) ([]ExpectedActivity, error) {
	var act []ExpectedActivity
	for _, p := range pact {
		activity, ok := activityFromPartnerMap[p]
		if !ok {
			return nil, bank.NewErrorFromCode(bank.ErrorCodeInvalidActivity)
		}

		act = append(act, activity)
	}

	return act, nil
}

type USResidencyStatus string

const (
	USResidencyStatusCitizen     = USResidencyStatus("us_citizen")
	USResidencyStatusResident    = USResidencyStatus("resident")
	USResidencyStatusNonResident = USResidencyStatus("non_resident")
)

var partnerUSResidencyStatusFrom = map[bank.ConsumerResidencyStatus]USResidencyStatus{
	bank.ResidencyStatusCitizen:     USResidencyStatusCitizen,
	bank.ResidencyStatusResident:    USResidencyStatusResident,
	bank.ResidencyStatusNonResident: USResidencyStatusNonResident,
}

func partnerResidencyFrom(res bank.ConsumerResidency) (USResidencyStatus, error) {
	if res.Country != bank.CountryUS {
		return USResidencyStatus(""), bank.NewErrorFromCode(bank.ErrorCodeInvalidCitizenStatus)
	}

	status, ok := partnerUSResidencyStatusFrom[res.Status]
	if !ok {
		return USResidencyStatus(""), bank.NewErrorFromCode(bank.ErrorCodeInvalidCitizenStatus)
	}

	return status, nil
}

var partnerResidencyStatusTo = map[USResidencyStatus]bank.ConsumerResidencyStatus{
	USResidencyStatusCitizen:     bank.ResidencyStatusCitizen,
	USResidencyStatusResident:    bank.ResidencyStatusResident,
	USResidencyStatusNonResident: bank.ResidencyStatusNonResident,
}

func partnerResidencyTo(res USResidencyStatus) (bank.ConsumerResidency, error) {
	status, ok := partnerResidencyStatusTo[res]
	if !ok {
		return bank.ConsumerResidency{}, bank.NewErrorFromCode(bank.ErrorCodeInvalidCitizenStatus)
	}

	return bank.ConsumerResidency{
		Country: bank.CountryUS,
		Status:  status,
	}, nil
}

// An array containing the customer's identity document details.
type ConsumerIdentificationRequest struct {
	// Document type
	Document ConsumerIdentityDocument `json:"document"`

	// Identity Document Number - Minimum of 6 characters and no special characters and spaces are allowed.
	Number string `json:"number"`

	// State issuing the document. This only applies to Driver License and State ID
	// Ex: CA for California Driver License
	IssuingState string `json:"issuing_state,omitempty"`

	// Issuer of the Document (Ex. USA for the US Passport and Driver's License)
	IssuingCountry Country3Alpha `json:"issuing_country"`

	// String based on ISO-8601 for specifying the date of issuance for the document. YYYY-MM-DD
	IssuedDate string `json:"issued_date"`

	// String based on ISO-8601 for specifying the date of expiration for the document. YYYY-MM-DD
	ExpirationDate string `json:"expiration_date"`
}

// Identity Document Type.
type ConsumerIdentityDocument string

const (
	ConsumerIdentityDocumentDriversLicense        = ConsumerIdentityDocument("drivers_license")
	ConsumerIdentityDocumentPassport              = ConsumerIdentityDocument("passport")
	ConsumerIdentityDocumentPassportCard          = ConsumerIdentityDocument("passport_card")
	ConsumerIdentityDocumentResidencyPermit       = ConsumerIdentityDocument("residency_permit")
	ConsumerIdentityDocumentWorkPermit            = ConsumerIdentityDocument("work_permit")
	ConsumerIdentityDocumentSocialSecurityCard    = ConsumerIdentityDocument("social_security_card")
	ConsumerIdentityDocumentStateID               = ConsumerIdentityDocument("state_id")
	ConsumerIdentityDocumentAlienRegistrationCard = ConsumerIdentityDocument("alien_registration_card")
	ConsumerIdentityDocumentUSAVisaH              = ConsumerIdentityDocument("H_visa")
	ConsumerIdentityDocumentUSAVisaL              = ConsumerIdentityDocument("L_visa")
	ConsumerIdentityDocumentUSAVisaO              = ConsumerIdentityDocument("O_visa")
	ConsumerIdentityDocumentUSAVisaE1             = ConsumerIdentityDocument("E1_visa")
	ConsumerIdentityDocumentUSAVisaE3             = ConsumerIdentityDocument("E3_visa")
	ConsumerIdentityDocumentUSAVisaI              = ConsumerIdentityDocument("I_visa")
	ConsumerIdentityDocumentUSAVisaP              = ConsumerIdentityDocument("P_visa")
	ConsumerIdentityDocumentUSAVisaTN             = ConsumerIdentityDocument("TN_visa")
	ConsumerIdentityDocumentUSAVisaTD             = ConsumerIdentityDocument("TD_visa")
	ConsumerIdentityDocumentUSAVisaR1             = ConsumerIdentityDocument("R1_visa")
)

var partnerConsumerIdentityDocumentFrom = map[bank.ConsumerIdentityDocument]ConsumerIdentityDocument{
	bank.ConsumerIdentityDocumentDriversLicense:        ConsumerIdentityDocumentDriversLicense,
	bank.ConsumerIdentityDocumentPassport:              ConsumerIdentityDocumentPassport,
	bank.ConsumerIdentityDocumentPassportCard:          ConsumerIdentityDocumentPassportCard,
	bank.ConsumerIdentityDocumentWorkPermit:            ConsumerIdentityDocumentWorkPermit,
	bank.ConsumerIdentityDocumentSocialSecurityCard:    ConsumerIdentityDocumentSocialSecurityCard,
	bank.ConsumerIdentityDocumentAlienRegistrationCard: ConsumerIdentityDocumentAlienRegistrationCard,
	bank.ConsumerIdentityDocumentUSAVisaH1B:            ConsumerIdentityDocumentUSAVisaH,
	bank.ConsumerIdentityDocumentUSAVisaH1C:            ConsumerIdentityDocumentUSAVisaH,
	bank.ConsumerIdentityDocumentUSAVisaH2A:            ConsumerIdentityDocumentUSAVisaH,
	bank.ConsumerIdentityDocumentUSAVisaH2B:            ConsumerIdentityDocumentUSAVisaH,
	bank.ConsumerIdentityDocumentUSAVisaH3:             ConsumerIdentityDocumentUSAVisaH,
	bank.ConsumerIdentityDocumentUSAVisaL1A:            ConsumerIdentityDocumentUSAVisaL,
	bank.ConsumerIdentityDocumentUSAVisaL1B:            ConsumerIdentityDocumentUSAVisaL,
	bank.ConsumerIdentityDocumentUSAVisaO1:             ConsumerIdentityDocumentUSAVisaO,
	bank.ConsumerIdentityDocumentUSAVisaE1:             ConsumerIdentityDocumentUSAVisaE1,
	bank.ConsumerIdentityDocumentUSAVisaE3:             ConsumerIdentityDocumentUSAVisaE3,
	bank.ConsumerIdentityDocumentUSAVisaI:              ConsumerIdentityDocumentUSAVisaI,
	bank.ConsumerIdentityDocumentUSAVisaP:              ConsumerIdentityDocumentUSAVisaP,
	bank.ConsumerIdentityDocumentUSAVisaTN:             ConsumerIdentityDocumentUSAVisaTN,
	bank.ConsumerIdentityDocumentUSAVisaTD:             ConsumerIdentityDocumentUSAVisaTD,
	bank.ConsumerIdentityDocumentUSAVisaR1:             ConsumerIdentityDocumentUSAVisaR1,
}

var partnerConsumerIdentityUploadDocumentFrom = map[bank.ConsumerIdentityDocument]ConsumerIdentityDocument{
	bank.ConsumerIdentityDocumentDriversLicense:        ConsumerIdentityDocumentDriversLicense,
	bank.ConsumerIdentityDocumentPassport:              ConsumerIdentityDocumentPassport,
	bank.ConsumerIdentityDocumentPassportCard:          ConsumerIdentityDocumentPassportCard,
	bank.ConsumerIdentityDocumentWorkPermit:            ConsumerIdentityDocumentWorkPermit,
	bank.ConsumerIdentityDocumentSocialSecurityCard:    ConsumerIdentityDocumentSocialSecurityCard,
	bank.ConsumerIdentityDocumentAlienRegistrationCard: ConsumerIdentityDocumentResidencyPermit,
	bank.ConsumerIdentityDocumentUSAVisaH1B:            ConsumerIdentityDocumentWorkPermit,
	bank.ConsumerIdentityDocumentUSAVisaH1C:            ConsumerIdentityDocumentWorkPermit,
	bank.ConsumerIdentityDocumentUSAVisaH2A:            ConsumerIdentityDocumentWorkPermit,
	bank.ConsumerIdentityDocumentUSAVisaH2B:            ConsumerIdentityDocumentWorkPermit,
	bank.ConsumerIdentityDocumentUSAVisaH3:             ConsumerIdentityDocumentWorkPermit,
	bank.ConsumerIdentityDocumentUSAVisaL1A:            ConsumerIdentityDocumentWorkPermit,
	bank.ConsumerIdentityDocumentUSAVisaL1B:            ConsumerIdentityDocumentWorkPermit,
	bank.ConsumerIdentityDocumentUSAVisaO1:             ConsumerIdentityDocumentWorkPermit,
	bank.ConsumerIdentityDocumentUSAVisaE1:             ConsumerIdentityDocumentWorkPermit,
	bank.ConsumerIdentityDocumentUSAVisaE3:             ConsumerIdentityDocumentWorkPermit,
	bank.ConsumerIdentityDocumentUSAVisaI:              ConsumerIdentityDocumentWorkPermit,
	bank.ConsumerIdentityDocumentUSAVisaP:              ConsumerIdentityDocumentWorkPermit,
	bank.ConsumerIdentityDocumentUSAVisaTN:             ConsumerIdentityDocumentWorkPermit,
	bank.ConsumerIdentityDocumentUSAVisaTD:             ConsumerIdentityDocumentWorkPermit,
	bank.ConsumerIdentityDocumentUSAVisaR1:             ConsumerIdentityDocumentWorkPermit,
}

// ConsumerPepRequest
// An object containing politically exposed persons information from the user
type ConsumerPepRequest struct {
	Association PepAssociation `json:"association"`

	// Name of the politically exposed person when he/she is a relative or a friend.
	// Required if pep.status is relative or friend.
	Name string `json:"name,omitempty"`

	// Position of the politically exposed person when he/she is a relative or a friend.
	// Required if pep.status is relative or friend.
	Position string `json:"position,omitempty"`
}

// PepAssociation
// One of a standard set of values to indicate if the customer is a politically exposed person.
// The value of this field will be answer to the question: Have you or any persons associated
// with you ever held a political office in a foreign country?
type PepAssociation string

const (
	PepAssociationNo       = "no"       // No politically exposed person is known.
	PepAssociationSelf     = "self"     // I am the politically exposed person.
	PepAssociationRelative = "relative" // A relative is the politically exposed person.
	PepAssociationFriend   = "friend"   // A friend is the politically exposed person.
)

// CreateConsumerResponsereateConsumerResponse is a BBVA compatibile structure for comnsumer create
// request responses
type CreateConsumerResponse struct {
	UserID    UserID                  `json:"user_id"`
	Contacts  []ContactEntityResponse `json:"contacts"`
	Addresses []AddressEntityResponse `json:"addresses"`
	KYCStatusResponse
}

func (resp CreateConsumerResponse) toPartnerIdentityStatusConsumerResponse(id bank.ConsumerID) (*bank.IdentityStatusConsumerResponse, error) {
	kyc, err := resp.KYC.toPartnerBankKYCResponse(resp.KYCNotes)
	if err != nil {
		return nil, err
	}

	return &bank.IdentityStatusConsumerResponse{
		ConsumerID: id,
		BankID:     bank.ConsumerBankID(resp.UserID),
		KYC:        *kyc,
	}, nil
}

// Create
// https://bbvaopenplatform.com/docs/reference%7Capiref%7Ccustomers%7C~consumer~v30%7Cpost
func (s *consumerService) Create(preq bank.CreateConsumerRequest) (*bank.IdentityStatusConsumerResponse, error) {
	c, err := data.NewConsumerService(s.request, bank.ProviderNameBBVA).GetByConsumerID(preq.ConsumerID)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	} else if c != nil {
		return nil, errors.New("user has already created")
	}

	consReq, err := newCreateConsumerRequest(preq)
	if err != nil {
		return nil, err
	}

	req, err := s.client.post("consumer/v3.0", s.request, consReq)
	if err != nil {
		return nil, err
	}

	var resp = CreateConsumerResponse{}
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}

	kycStatus, ok := kycStatusToPartnerMap[KYCStatus(strings.ToUpper(string(resp.KYC.Status)))]
	if !ok {
		return nil, errors.New("invalid KYC status")
	}

	b, err := json.Marshal(
		map[string]interface{}{
			"kycNotes":         resp.KYCNotes,
			"digitalFootprint": resp.DigitalFootprint,
		},
	)
	if err != nil {
		return nil, err
	}

	// Save consumer entity to partner entity table
	cr := data.ConsumerCreate{
		ConsumerID: preq.ConsumerID,
		BankID:     bank.ConsumerBankID(resp.UserID),
		BankExtra:  types.JSONText(string(b)),
		KYCStatus:  kycStatus,
	}

	entity, err := data.NewConsumerService(s.request, bank.ProviderNameBBVA).Create(cr)
	if err != nil {
		return nil, err
	}

	// Map contacts
	err = createConsumerContacts(s.client, s.request, entity, resp.Contacts)
	if err != nil {
		return nil, err
	}

	// Update Contacts
	updateConsumerContact(s.client, s.request, entity, bank.ConsumerPropertyTypeContactEmail, preq.Email)
	updateConsumerContact(s.client, s.request, entity, bank.ConsumerPropertyTypeContactPhone, preq.Phone)

	// Map addresses
	err = createConsumerAddresses(s.client, s.request, entity, resp.Addresses)
	if err != nil {
		return nil, err
	}

	// Updated addresses
	updateConsumerAddress(s.client, s.request, entity, bank.ConsumerPropertyTypeAddressLegal, preq.LegalAddress)

	if preq.MailingAddress != nil {
		updateConsumerAddress(s.client, s.request, entity, bank.ConsumerPropertyTypeAddressMailing, *preq.MailingAddress)
	}

	if preq.WorkAddress != nil {
		updateConsumerAddress(s.client, s.request, entity, bank.ConsumerPropertyTypeAddressWork, *preq.WorkAddress)
	}

	// Return response
	return resp.toPartnerIdentityStatusConsumerResponse(preq.ConsumerID)
}

func (s *consumerService) Status(id bank.ConsumerID) (*bank.IdentityStatusConsumerResponse, error) {
	c, err := data.NewConsumerService(s.request, bank.ProviderNameBBVA).GetByConsumerID(id)
	if err != nil {
		return nil, err
	}

	req, err := s.client.get("consumer/v3.0/identity", s.request)
	if err != nil {
		return nil, err
	}

	req.Header.Set("OP-User-Id", string(c.BankID))
	var resp = KYCStatusResponse{}
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}

	// Save consumer entity to partner entity table in case of change
	status, ok := kycStatusToPartnerMap[KYCStatus(strings.ToUpper(string(resp.KYC.Status)))]
	if !ok {
		return nil, errors.New("Invalid KYC Status")
	}

	b, err := json.Marshal(
		map[string]interface{}{
			"kycNotes":         resp.KYCNotes,
			"digitalFootprint": resp.DigitalFootprint,
		},
	)
	if err != nil {
		return nil, err
	}

	u := data.ConsumerUpdate{
		ID:        c.ID,
		BankExtra: types.NullJSONText{types.JSONText(string(b)), true},
		KYCStatus: &status,
	}

	_, err = data.NewConsumerService(s.request, bank.ProviderNameBBVA).Update(u)
	if err != nil {
		return nil, err
	}

	kyc, err := resp.KYC.toPartnerBankKYCResponse(resp.KYCNotes)
	if err != nil {
		return nil, err
	}

	return &bank.IdentityStatusConsumerResponse{
		ConsumerID: c.ConsumerID,
		BankID:     c.BankID,
		KYC:        *kyc,
	}, nil
}

func (s *consumerService) Update(preq bank.UpdateConsumerRequest) (*bank.IdentityStatusConsumerResponse, error) {
	c, err := data.NewConsumerService(s.request, bank.ProviderNameBBVA).GetByConsumerID(preq.ConsumerID)
	if err != nil {
		return nil, err
	}

	consReq, err := newUpdateConsumerRequest(c, preq)
	if err != nil {
		return nil, err
	}

	req, err := s.client.patch("consumer/v3.0", s.request, consReq)
	if err != nil {
		return nil, err
	}

	req.Header.Set("OP-User-Id", string(consReq.UserID))
	if err := s.client.do(req, nil); err != nil {
		return nil, err
	}

	return s.Status(c.ConsumerID)
}

/* VerifyIdentity
// https://bbvaopenplatform.com/docs/guides%7Capicontent%7Cverify-a-user-identity
func (s *consumerService) VerifyIdentity(request bank.ConsumerIdentityRequest) (*bank.ConsumerIdentityResponse, error) {
	c, err := data.NewConsumerService(s.request, bank.ProviderNameBBVA).GetByConsumerID(preq.ConsumerID)
    if err != nil {
        return nil, err
    }

	req, err := s.client.post("consumer/v3.0/identity", s.request, request)
	if err != nil {
		return nil, err
	}

	req.Header.Set("OP-User-Id", c.BankID)
	var response = bank.ConsumerIdentityResponse{}
	if err := s.client.do(req, &response); err != nil {
		return nil, err
	}
	return &response, nil
} */

func (s *consumerService) UploadIdentityDocument(preq bank.ConsumerIdentityDocumentRequest) (*bank.IdentityDocumentResponse, error) {
	c, err := data.NewConsumerService(s.request, bank.ProviderNameBBVA).GetByConsumerID(preq.ConsumerID)
	if err != nil {
		return nil, err
	}

	docType, ok := partnerConsumerIdentityUploadDocumentFrom[preq.IdentityDocument]
	if !ok {
		return nil, errors.New("invalid consumer identity document type")
	}

	if docType == ConsumerIdentityDocumentAlienRegistrationCard {
		docType = ConsumerIdentityDocumentResidencyPermit
	}

	fileType, ok := partnerContentTypeFrom[preq.ContentType]
	if !ok {
		return nil, errors.New("invalid consumer identity document content type")
	}

	idvs, err := documentIDVerifyPartnerFrom(preq.IDVerifyRequired)
	if err != nil {
		return nil, err
	}

	docReq := &ConsumerIdentityDocumentRequest{
		File:             base64.StdEncoding.EncodeToString(preq.Content),
		FileType:         fileType,
		IDVerifyRequired: idvs,
		DocType:          docType,
	}

	r, err := s.client.post("consumer/v3.1/identity/document", s.request, docReq)
	if err != nil {
		return nil, err
	}

	r.Header.Set("OP-User-Id", string(c.BankID))
	var response = bank.IdentityDocumentResponse{}
	if err := s.client.do(r, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

/* Update phone or email */
func (s *consumerService) UpdateContact(id bank.ConsumerID, p bank.ConsumerPropertyType, val string) error {
	c, err := data.NewConsumerService(s.request, bank.ProviderNameBBVA).GetByConsumerID(id)
	if err != nil {
		return err
	}

	return updateConsumerContact(s.client, s.request, c, p, val)
}

/* Update address */
func (s *consumerService) UpdateAddress(id bank.ConsumerID, p bank.ConsumerPropertyType, a bank.AddressRequest) error {
	c, err := data.NewConsumerService(s.request, bank.ProviderNameBBVA).GetByConsumerID(id)
	if err != nil {
		return err
	}

	return updateConsumerAddress(s.client, s.request, c, p, a)
}

func (s *consumerService) Delete(id bank.ConsumerID) error {
	return data.NewConsumerService(s.request, bank.ProviderNameBBVA).DeleteByID(id)
}
