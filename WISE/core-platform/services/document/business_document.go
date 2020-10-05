/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

package document

import (
	"github.com/wiseco/core-platform/shared"
)

//DocumentType the of a document e.g "certificateOfFormation"
type DocumentType string

// Specific to business documents
// We could add all the types including the users maybe?
var documentTypes = map[DocumentType]string{
	BusinessDocTypeArticlesOfIncorporation: "articlesOfIncorporation",
	BusinessDocTypeArticlesOfOrganization:  "articlesOfOrganization",
	BusinessDocTypeAssumedNameCertificate:  "assumedNameCertificate",
	BusinessDocTypeBusinessLicense:         "businessLicense",
	BusinessDocTypePartnershipCertificate:  "certificateOfPartnership",
	BusinessDocTypePartnerAgreement:        "partnershipAgreement",
	BusinessDocTypeCertOfFormation:         "certificateOfFormation",
	BusinessDocTypeDriversLicense:          "driversLicense",
	BusinessDocTypeDriversOther:            "other",
}

const (
	BusinessDocTypeArticlesOfIncorporation = DocumentType("articlesOfIncorporation")  // Corporate charter
	BusinessDocTypeArticlesOfOrganization  = DocumentType("articlesOfOrganization")   // Initial statements to form an LLC
	BusinessDocTypeAssumedNameCertificate  = DocumentType("assumedNameCertificate")   // Certificate for DBA
	BusinessDocTypeBusinessLicense         = DocumentType("businessLicense")          // Business license
	BusinessDocTypePartnershipCertificate  = DocumentType("certificateOfPartnership") // Certificate of partnership
	BusinessDocTypePartnerAgreement        = DocumentType("partnershipAgreement")     // Partnership agreement
	BusinessDocTypeCertOfFormation         = DocumentType("certificateOfFormation")   // Certificate of LLC formation
	BusinessDocTypeDriversLicense          = DocumentType("driversLicense")           // Drivers Lincense of owner (Sole Prop)
	BusinessDocTypeDriversOther            = DocumentType("other")                    // Other business documents
)

//BusinessDocumentUpdate ..
type BusinessDocumentUpdate struct {
	// Document id
	ID shared.BusinessDocumentID `json:"id" db:"id"`

	BusinessID shared.BusinessID `json:"businessId" db:"business_id"`

	// Document number
	Number *DocumentNumber `json:"number" db:"number"`

	// Document type
	DocType *string `json:"docType" db:"doc_type"`

	// Issuing state
	IssuingState *string `json:"issuingState" db:"issuing_state"`

	// Issuing country
	IssuingCountry *string `json:"issuingCountry" db:"issuing_country"`

	// Issuing date
	IssuedDate *shared.Date `json:"issuedDate" db:"issued_date"`

	// Expiration Date
	ExpirationDate *shared.Date `json:"expirationDate" db:"expiration_date"`

	// Mime type
	ContentType *string `json:"contentType" db:"content_type"`

	//document storage key
	StorageKey *string `json:"storage_key,omitempty" db:"storage_key"`
}

//DocumentCreateResponse  ..
type DocumentCreateResponse struct {
	SignedURL *string          `json:"uploadURL"`
	Document  BusinessDocument `json:"document"`
}

//BusinessDocumentCreate ..
type BusinessDocumentCreate struct {
	BusinessID shared.BusinessID `json:"businessId" db:"business_id"`

	// User Created ID
	CreatedUserID shared.UserID `db:"created_user_id"`

	// Document number
	Number *DocumentNumber `json:"number" db:"number"`

	// Document type
	DocType *string `json:"docType" db:"doc_type"`

	// Issuing state
	IssuingState *string `json:"issuingState" db:"issuing_state"`

	// Issuing country
	IssuingCountry *string `json:"issuingCountry" db:"issuing_country"`

	// Issuing date
	IssuedDate *shared.Date `json:"issuedDate" db:"issued_date"`

	// Expiration date
	ExpirationDate *shared.Date `json:"expirationDate" db:"expiration_date"`

	// Mime type
	ContentType string `json:"contentType" db:"content_type"`

	//document storage key
	StorageKey *string `json:"storage_key,omitempty" db:"storage_key"`
}

//BusinessDocument ..
type BusinessDocument struct {
	// Document id
	ID            shared.BusinessDocumentID `json:"id" db:"id"`
	BusinessID    shared.BusinessID         `json:"businessId" db:"business_id"`
	CreatedUserID shared.UserID             `json:"createdUserId" db:"created_user_id"`
	Document
}

func (doc DocumentType) isValid() bool {
	_, ok := documentTypes[doc]
	return ok
}

//NewDocType return a new document type and a boolean value indication if created type is valid
func NewDocType(v string) (DocumentType, bool) {
	d := DocumentType(v)
	return d, d.isValid()
}

//String DocumentType to string
func (doc DocumentType) String() string {
	return string(doc)
}
