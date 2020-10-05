package review

import (
	"github.com/wiseco/core-platform/services"
	biz "github.com/wiseco/core-platform/services/business"
	svc "github.com/wiseco/core-platform/services/user"
)

type VerificationError string

//NewVerificationError ..
func NewVerificationError(v string) VerificationError {
	return VerificationError(v)
}

const (
	KYCParamErrorIsRestricted     = VerificationError("isRestricted")
	KYCParamErrorFirstName        = VerificationError("firstName")
	KYCParamErrorLastName         = VerificationError("lastName")
	KYCParamErrorEmail            = VerificationError("email")
	KYCParamErrorPhone            = VerificationError("phone")
	KYCParamErrorDateOfBirth      = VerificationError("dateOfBirth")
	KYCParamErrorTaxID            = VerificationError("taxId")
	KYCParamErrorTaxIDType        = VerificationError("taxIdType")
	KYCParamErrorLegalAddress     = VerificationError("legalAddress")
	KYCParamErrorResidency        = VerificationError("residency")
	KYCParamErrorCitizenship      = VerificationError("citizenship")
	KYCParamErrorOccupation       = VerificationError("occupation")
	KYCParamErrorLegalName        = VerificationError("legalName")
	KYCParamErrorIncomeType       = VerificationError("incomeType")
	KYCParamErrorActivityType     = VerificationError("activityType")
	KYCParamErrorEntityType       = VerificationError("entityType")
	KYCParamErrorIndustryType     = VerificationError("industryType")
	KYCParamErrorDeactivated      = VerificationError("deactivated")
	KYCParamErrorOperationType    = VerificationError("operationType")
	KYCParamErrorPurpose          = VerificationError("purpose")
	KYCParamErrorOriginState      = VerificationError("originState")
	KYCParamErrorCountry          = VerificationError("originCountry")
	KYCParamErrorOriginDate       = VerificationError("originDate")
	KYCParamErrorTitleType        = VerificationError("titleType")
	KYCParamErrorTypeFormationDoc = VerificationError("Formation document required")
	KYCReviewErrorLegalAddress    = VerificationError("legalAddress")
	KYCReviewErrorTaxID           = VerificationError("taxId")
	KYCReviewErrorName            = VerificationError("name")
	KYCReviewErrorDOB             = VerificationError("dateOfBirth")
	KYCReviewErrorOFAC            = VerificationError("ofac")
	KYCReviewErrorMismatch        = VerificationError("mismatch")
	KYCReviewErrorDocumentReq     = VerificationError("documentRequired")
	KYCReviewErrorIDVerification  = VerificationError("idVerification")
	KYCErrorTypeMember            = VerificationError("members")
	KYCErrorTypeOther             = VerificationError("other")
	KYCErrorTypeApproved          = VerificationError("Already approved")
	KYCErrorTypeDeclined          = VerificationError("already declined")
	KYCErrorTypeInProgress        = VerificationError("inProgress")
	KYCErrorTypeParam             = VerificationError("param")
	KYCErrorTypeReview            = VerificationError("review")
	KYCErrorTypeDeactivated       = VerificationError("deactivated")
	KYCErrorTypeRestricted        = VerificationError("restricted")
)

func (v VerificationError) String() string {
	return string(v)
}
func (v VerificationError) Error() string {
	return v.String()
}

func checkCommon(b biz.Business) []VerificationError {
	var errors []VerificationError

	if b.IsRestricted {
		errors = append(errors, KYCParamErrorIsRestricted)
	}

	if b.IndustryType == nil || len(*b.IndustryType) == 0 {
		errors = append(errors, KYCParamErrorIndustryType)
	}

	if b.TaxIDType == nil || len(*b.TaxIDType) == 0 {
		errors = append(errors, KYCParamErrorTaxIDType)
	}

	if b.TaxID == nil || len(*b.TaxID) == 0 {
		errors = append(errors, KYCParamErrorTaxID)
	}

	if b.OriginCountry == nil || len(*b.OriginCountry) == 0 {
		errors = append(errors, KYCParamErrorCountry)
	}

	if b.OriginState == nil || len(*b.OriginState) == 0 {
		errors = append(errors, KYCParamErrorOriginState)
	}

	if b.OriginDate == nil {
		errors = append(errors, KYCParamErrorOriginDate)
	}

	if b.Purpose == nil || len(*b.Purpose) == 0 {
		errors = append(errors, KYCParamErrorPurpose)
	}

	if b.OperationType == nil || len(*b.OperationType) == 0 {
		errors = append(errors, KYCParamErrorOperationType)
	}

	if b.Deactivated != nil {
		errors = append(errors, KYCParamErrorDeactivated)
	}

	if b.Email == nil || len(*b.Email) < 6 {
		errors = append(errors, KYCParamErrorEmail)
	}

	if b.Phone == nil || len(*b.Phone) < 10 {
		errors = append(errors, KYCParamErrorPhone)
	}

	if len(b.ActivityType) == 0 {
		errors = append(errors, KYCParamErrorActivityType)
	}

	if b.LegalAddress == nil {
		errors = append(errors, KYCParamErrorLegalAddress)
	}

	if b.EntityType == nil || len(*b.EntityType) == 0 {
		errors = append(errors, KYCParamErrorEntityType)
	}

	if b.FormationDocumentID == nil {
		errors = append(errors, KYCParamErrorTypeFormationDoc)
	}

	if b.LegalName == nil {
		errors = append(errors, KYCParamErrorLegalName)
	}

	return errors
}

func checkConsumerCommons(c *svc.Consumer) []string {
	// Check params
	kycErrors := []string{}
	if len(c.FirstName) < 2 {
		kycErrors = append(kycErrors, services.KYCParamErrorFirstName)
	}

	if len(c.LastName) < 2 {
		kycErrors = append(kycErrors, services.KYCParamErrorLastName)
	}

	if c.Email == nil || len(*c.Email) < 6 {
		kycErrors = append(kycErrors, services.KYCParamErrorEmail)
	}

	if c.Phone == nil || len(*c.Phone) < 10 {
		kycErrors = append(kycErrors, services.KYCParamErrorPhone)
	}

	if c.DateOfBirth == nil {
		kycErrors = append(kycErrors, services.KYCParamErrorDateOfBirth)
	}

	if c.TaxID == nil {
		kycErrors = append(kycErrors, services.KYCParamErrorTaxID)
	}

	if c.TaxIDType == nil {
		kycErrors = append(kycErrors, services.KYCParamErrorTaxIDType)
	}

	if c.Residency == nil {
		kycErrors = append(kycErrors, services.KYCParamErrorResidency)
	}

	if len(c.CitizenshipCountries) == 0 {
		kycErrors = append(kycErrors, services.KYCParamErrorCitizenship)
	}
	return kycErrors
}

func errListString(item []VerificationError) []string {
	var list []string
	for _, i := range item {
		list = append(list, i.String())
	}
	return list
}
