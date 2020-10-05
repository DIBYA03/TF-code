package signature

import (
	"time"

	"github.com/wiseco/core-platform/services"
	"github.com/wiseco/core-platform/shared"
)

type SignatureRequest struct {
	ID                 shared.SignatureRequestID  `json:"id" db:"id"` // TODO use signature ID
	BusinessID         shared.BusinessID          `json:"businessId" db:"business_id"`
	TemplateType       SignatureRequestTemplate   `json:"templateType" db:"template_type"`
	TemplateProvider   SignatureRequestProvider   `json:"templateProvider" db:"template_provider"`
	SignatureRequestID string                     `json:"signatureRequestId" db:"signature_request_id"`
	SignatureID        string                     `json:"signatureId" db:"signature_id"`
	SignatureStatus    SignatureRequestStatus     `json:"signatureStatus" db:"signature_status"`
	DocumentID         *shared.BusinessDocumentID `json:"documentId" db:"document_id"`
	SignURL            string                     `json:"signURL"`
	Created            time.Time                  `json:"created" db:"created"`
	Modified           time.Time                  `json:"modified" db:"modified"`
}

type SignatureRequestCreate struct {
	BusinessID         shared.BusinessID        `json:"businessId" db:"business_id"`
	TemplateType       SignatureRequestTemplate `json:"templateType" db:"template_type"`
	TemplateProvider   SignatureRequestProvider `json:"templateProvider" db:"template_provider"`
	SignatureRequestID string                   `json:"signatureRequestId" db:"signature_request_id"`
	SignatureStatus    SignatureRequestStatus   `json:"signatureStatus" db:"signature_status"`
	SignatureID        string                   `json:"signatureId" db:"signature_id"`
}

type SignatureRequestJoin struct {
	EntityType   string               `db:"entity_type"`
	LegalName    *string              `db:"legal_name"`
	DBA          services.StringArray `db:"dba"`
	LegalAddress services.Address     `db:"legal_address"`
	TaxID        services.TaxID       `db:"business.tax_id"`
	TaxIDType    services.TaxIDType   `db:"business.tax_id_type"`
	TitleType    string               `db:"title_type"`
	TitleOther   *string              `db:"title_other"`
	FirstName    string               `db:"first_name"`
	MiddleName   *string              `db:"middle_name"`
	LastName     string               `db:"last_name"`
	EmailAddress string               `db:"email"`
	Ownership    int                  `db:"ownership"`
}

type SignatureRequestProvider string

const (
	SignatureRequestProviderHellosign = SignatureRequestProvider("hellosign")
)

type SignatureRequestTemplate string

const (
	SignatureRequestTemplateControlPersonCertification = SignatureRequestTemplate("controlPersonCertification")
)

type SignatureRequestStatus string

const (
	SignatureRequestStatusPending   = SignatureRequestStatus("pending")
	SignatureRequestStatusCanceled  = SignatureRequestStatus("canceled")
	SignatureRequestStatusCompleted = SignatureRequestStatus("completed")
)

type SQSMessage struct {
	SignatureRequestID string    `json:"signatureRequestId"`
	EventType          EventType `json:"eventType"`
}

type EventType string

const (
	EventTypeSignatureRequestViewed       = EventType("eventTypeSignatureRequestViewed")
	EventTypeSignatureRequestSigned       = EventType("eventTypeSignatureRequestSigned")
	EventTypeSignatureRequestDownloadable = EventType("eventTypeSignatureRequestDownloadable")
	EventTypeSignatureRequestSent         = EventType("eventTypeSignatureRequestSent")
	EventTypeSignatureRequestDeclined     = EventType("eventTypeSignatureRequestDeclined")
	EventTypeSignatureRequestReassigned   = EventType("eventTypeSignatureRequestReassigned")
	EventTypeSignatureRequestRemind       = EventType("eventTypeSignatureRequestRemind")
	EventTypeSignatureRequestAllSigned    = EventType("eventTypeSignatureRequestAllSigned")
	EventTypeSignatureRequestEmailBounce  = EventType("eventTypeSignatureRequestEmailBounce")
	EventTypeSignatureRequestInvalid      = EventType("eventTypeSignatureRequestInvalid")
	EventTypeSignatureRequestCanceled     = EventType("eventTypeSignatureRequestCanceled")
	EventTypeSignatureRequestPrepared     = EventType("eventTypeSignatureRequestPrepared")
	EventTypeOther                        = EventType("eventTypeOther")
)

const (
	LegalName                     = "legal_name"
	LimitedLiabilityPartnership   = "llp"
	SingleLimitedLiabilityCompany = "sllc"
	MultiLimitedLiabilityCompany  = "mllc"
	Corporation                   = "corp"
	SoleProp                      = "sp"
	DBA                           = "dba"
	StreetAddress                 = "street_address"
	City                          = "city"
	State                         = "state"
	PostalCode                    = "postal_code"
	EIN                           = "ein"
	OperatingState                = "op_state"
	ControlPersonName             = "name"
	ControlPersonTitle            = "title"
	Day                           = "day"
	Month                         = "mon"
	Year                          = "year"
)

const (
	BusinessEntitySoleProprietor              = "soleProprietor"
	BusinessEntityAssociation                 = "association"
	BusinessEntityProfessionalAssociation     = "professionalAssociation"
	BusinessEntitySingleMemberLLC             = "singleMemberLLC"
	BusinessEntityLimitedLiabilityCompany     = "limitedLiabilityCompany"
	BusinessEntityGeneralPartnership          = "generalPartnership"
	BusinessEntityLimitedPartnership          = "limitedPartnership"
	BusinessEntityLimitedLiabilityPartnership = "limitedLiabilityPartnership"
	BusinessEntityProfessionalCorporation     = "professionalCorporation"
	BusinessEntityUnlistedCorporation         = "unlistedCorporation"
)

var controlPersonTitle = map[string]string{
	"chiefExecutiveOfficer": "Chief Executive Officer",
	"chiefFinancialOfficer": "Chief Financial Officer",
	"chiefOperatingOfficer": "Chief Operating Officer",
	"president":             "President",
	"vicePresident":         "Vice President",
	"seniorVicePresident":   "Senior Vice President",
	"treasurer":             "Treasurer",
	"secretary":             "Secretary",
	"generalPartner":        "General Partner",
	"manager":               "Manager",
	"member":                "Member",
	"owner":                 "Owner",
}

var SignatureRequestTemplateTo = map[SignatureRequestTemplate]string{
	SignatureRequestTemplateControlPersonCertification: "Control Person",
}
