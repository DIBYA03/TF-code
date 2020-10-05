/*******************************************************************
 * Cop yright 2019 Wise Company
 ********************************************************************/

package services

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"net/url"

	partnerbank "github.com/wiseco/core-platform/partner/bank"
	"github.com/wiseco/core-platform/shared"
)

type TaxIDType string

const TaxIDTypeNone = TaxIDType("")

func TaxIDTypeValue(idType *TaxIDType) TaxIDType {
	if idType == nil {
		return TaxIDTypeNone
	}

	return *idType
}

const (
	TaxIDTypeSSN  = TaxIDType("ssn")  // Social Security Number (SSN)
	TaxIDTypeITIN = TaxIDType("itin") // Individual Taxpayer Identification Number (ITIN)
	TaxIDTypeEIN  = TaxIDType("ein")  // Employer Identification Number (EIN)
)

// Tax ID
type TaxID string

const TaxIDNone = TaxID("")

func (n *TaxID) String() string {
	return string(*n)
}

func (n *TaxID) MarshalJSON() ([]byte, error) {
	return json.Marshal(MaskLeft(n.String(), 4))
}

func TaxIDValue(id *TaxID) TaxID {
	if id == nil {
		return TaxIDNone
	}

	return *id
}

const (
	ActivityTypeCash         = "cash"
	ActivityTypeCheck        = "check"
	ActivityTypeDomesticWire = "domesticWireTransfer"
	ActivityTypeIntlWire     = "internationalWireTransfer"
	ActivityTypeDomesticAch  = "domesticACH"
	ActivityTypeIntlAch      = "internationalACH"
)

type KYCStatus string

const (
	KYCStatusNotStarted = KYCStatus("notStarted") // Not yet started (start state)
	KYCStatusSubmitted  = KYCStatus("submitted")  // Submitted for kyc
	KYCStatusReview     = KYCStatus("review")     // KYC in review
	KYCStatusApproved   = KYCStatus("approved")   // KYC approved
	KYCStatusDeclined   = KYCStatus("declined")   // KYC declined
)

var PartnerKYCStatusFromMap = map[partnerbank.KYCStatus]KYCStatus{
	partnerbank.KYCStatusApproved: KYCStatusApproved,
	partnerbank.KYCStatusReview:   KYCStatusReview,
	partnerbank.KYCStatusDeclined: KYCStatusDeclined,
}

var PartnerKYCStatusToMap = map[KYCStatus]partnerbank.KYCStatus{
	KYCStatusApproved: partnerbank.KYCStatusApproved,
	KYCStatusReview:   partnerbank.KYCStatusReview,
	KYCStatusDeclined: partnerbank.KYCStatusDeclined,
}

const (
	KYCParamErrorIsRestricted   = "isRestricted"
	KYCParamErrorFirstName      = "firstName"
	KYCParamErrorLastName       = "lastName"
	KYCParamErrorEmail          = "email"
	KYCParamErrorPhone          = "phone"
	KYCParamErrorDateOfBirth    = "dateOfBirth"
	KYCParamErrorTaxID          = "taxId"
	KYCParamErrorTaxIDType      = "taxIdType"
	KYCParamErrorLegalAddress   = "legalAddress"
	KYCParamErrorResidency      = "residency"
	KYCParamErrorCitizenship    = "citizenship"
	KYCParamErrorOccupation     = "occupation"
	KYCParamErrorLegalName      = "legalName"
	KYCParamErrorIncomeType     = "incomeType"
	KYCParamErrorActivityType   = "activityType"
	KYCParamErrorEntityType     = "entityType"
	KYCParamErrorIndustryType   = "industryType"
	KYCParamErrorDeactivated    = "deactivated"
	KYCParamErrorOperationType  = "operationType"
	KYCParamErrorPurpose        = "purpose"
	KYCParamErrorOriginState    = "originState"
	KYCParamErrorCountry        = "originCountry"
	KYCParamErrorOriginDate     = "originDate"
	KYCParamErrorTitleType      = "titleType"
	KYCParamErrorDocType        = "docType"
	KYCParamErrorIssuingCountry = "issuingCountry"
	KYCParamErrorIssuedDate     = "issuingCountry"
	KYCParamErrorExpirationDate = "issuingCountry"
	KYCParamErrorDocNumber      = "number"
)

const (
	KYCReviewErrorLegalAddress   = "legalAddress"
	KYCReviewErrorTaxID          = "taxId"
	KYCReviewErrorName           = "name"
	KYCReviewErrorDOB            = "dateOfBirth"
	KYCReviewErrorOFAC           = "ofac"
	KYCReviewErrorMismatch       = "mismatch"
	KYCReviewErrorDocumentReq    = "documentRequired"
	KYCReviewErrorIdVerification = "idVerification"
)

const (
	KYCErrorTypeOther       = "other"
	KYCErrorTypeInProgress  = "inProgress"
	KYCErrorTypeParam       = "param"
	KYCErrorTypeReview      = "review"
	KYCErrorTypeDeactivated = "deactivated"
	KYCErrorTypeRestricted  = "restricted"
)

const (
	OccTypeAgriculture                = "agriculture"
	OccTypeClergyMinistryStaff        = "clergyMinistryStaff"
	OccTypeConstructionIndustrial     = "constructionIndustrial"
	OccTypeEducation                  = "education"
	OccTypeFinanceAccountingTax       = "financeAccountingTax"
	OccTypeFireFirstResponders        = "fireFirstResponders"
	OccTypeHealthcare                 = "healthcare"
	OccTypeHomemaker                  = "homemaker"
	OccTypeLaborGeneral               = "laborGeneral"
	OccTypeLaborSkilled               = "laborSkilled"
	OccTypeLawEnforcementSecurity     = "lawEnforcementSecurity"
	OccTypeLegalServices              = "legalServices"
	OccTypeMilitary                   = "military"
	OccTypeNotaryRegistrar            = "notaryRegistrar"
	OccTypePrivateInvestor            = "privateInvestor"
	OccTypeProfessionalAdministrative = "professionalAdministrative"
	OccTypeProfessionalManagement     = "professionalManagement"
	OccTypeProfessionalOther          = "professionalOther"
	OccTypeProfessionalTechnical      = "professionalTechnical"
	OccTypeRetired                    = "retired"
	OccTypeSales                      = "sales"
	OccTypeSelfEmployed               = "selfEmployed"
	OccTypeStudent                    = "student"
	OccTypeTransportation             = "transportation"
	OccTypeUnemployed                 = "unemployed"
)

const (
	IncomeTypeInheritance       = "inheritance"
	IncomeTypeSalary            = "salary"
	IncomeTypeSaleOfCompany     = "saleOfCompany"
	IncomeTypeSaleOfProperty    = "saleOfProperty"
	IncomeTypeInvestments       = "investments"
	IncomeTypeLifeInsurance     = "lifeInsurance"
	IncomeTypeDivorceSettlement = "divorceSettlement"
	IncomeTypeOther             = "other"
)

const (
	ResidencyStatusCitizen     = "citizen"     // Citizen
	ResidencyStatusResident    = "resident"    // Legal non-citizen resident
	ResidencyStatusNonResident = "nonResident" // Legal temporary non-citizen resident
)

type Residency struct {
	Country string `json:"country"`
	Status  string `json:"status"`
}

// SQL value marshaller
func (res Residency) Value() (driver.Value, error) {
	return json.Marshal(res)
}

// SQL scan unmarshaller
func (res *Residency) Scan(src interface{}) error {
	source, ok := src.([]byte)
	if !ok {
		return errors.New("Type convertible to []byte")
	}

	var out Residency
	err := json.Unmarshal(source, &out)
	if err != nil {
		return err
	}

	*res = out
	return nil
}

type StringArray []string

// SQL value marshaller
func (sa StringArray) Value() (driver.Value, error) {
	return json.Marshal(sa)
}

// SQL scan unmarshaller
func (sa *StringArray) Scan(src interface{}) error {
	source, ok := src.([]byte)
	if !ok {
		return errors.New("Type convertible to []byte")
	}

	var out StringArray
	err := json.Unmarshal(source, &out)
	if err != nil {
		return err
	}

	*sa = out
	return nil
}

func (sa StringArray) ToArray() []string {
	return []string(sa)
}

func (sa StringArray) ToPartnerBankIncome() []partnerbank.Income {
	// Transform into PartnerBank - TODO: use string not types
	incs := []partnerbank.Income{}
	for _, inc := range sa {
		incs = append(incs, partnerbank.Income(inc))
	}

	return incs
}

func (sa StringArray) ToPartnerBankActivity() []partnerbank.ExpectedActivity {
	acts := []partnerbank.ExpectedActivity{}
	for _, a := range sa {
		acts = append(acts, partnerbank.ExpectedActivity(a))
	}

	return acts
}

func ValidateTaxID(id *TaxID, idType *TaxIDType) (*TaxID, error) {
	// Ignore if tax id is nil
	if id == nil {
		return nil, nil
	}

	// If type is nil but not id then return error
	if idType == nil {
		return nil, errors.New("invalid tax id or type")
	}

	// Only SSN and EIN's are handled
	switch *idType {
	case TaxIDTypeSSN, TaxIDTypeEIN:
		taxID := TaxID(shared.StripNonDigits(string(*id)))
		if len(taxID) != 9 {
			return nil, errors.New("invalid tax id or type")
		}

		return &taxID, nil
	default:
		return nil, errors.New("invalid tax id or type")
	}
}

type Website string

func (w *Website) String() string {
	return string(*w)
}

func (w *Website) IsValidURL() bool {
	if w == nil {
		return true
	}

	if len(*w) == 0 {
		return true
	}

	_, err := url.ParseRequestURI(w.String())
	if err != nil {
		return false
	}

	u, err := url.Parse(w.String())
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}

	return true
}
