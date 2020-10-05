package document

import (
	"time"

	"github.com/jmoiron/sqlx/types"
	"github.com/wiseco/core-platform/services/csp"
	docs "github.com/wiseco/core-platform/services/document"
	"github.com/wiseco/core-platform/shared"
)

// BaseDocument ..
type BaseDocument struct {
	CreatedUserID shared.UserID `json:"createdUserId" db:"created_user_id"`

	// Document number
	Number *string `json:"number" db:"number"`

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
	ContentType string `json:"contentType" db:"content_type"`

	// Created timestamp
	Created time.Time `json:"created" db:"created"`

	//ContentUploaded indicates content on s3 has been uploaded
	ContentUploaded *time.Time `json:"contentUploaded" db:"content_uploaded"`

	//Deleted
	Deleted *time.Time `json:"deleted" db:"deleted"`

	// Update timestamp
	Modified time.Time `json:"modified" db:"modified"`

	// Storage key reference
	StorageKey *string `json:"storage_key,omitempty" db:"storage_key"`
}

// BusinessDocumentResponse ..
type BusinessDocumentResponse struct {
	SignedURL *string          `json:"uploadURL"`
	Document  BusinessDocument `json:"document"`
}

// BusinessDocumentUpdate ..
type BusinessDocumentUpdate struct {

	// Document number
	Number *string `json:"number" db:"number"`

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

	//Mark the document for conent update on s3
	UpdatingContent *bool `json:"updatingContent" `

	UseFormation *bool `json:"useFormation"`

	//document storage key
	StorageKey *string `json:"storage_key,omitempty" db:"storage_key"`
}

type BusinessDocument struct {
	BaseDocument
	// Document id
	ID shared.BusinessDocumentID `json:"id" db:"id"`

	BusinessID string `json:"businessId" db:"business_id"`
}

type BusinessDocumentCreate struct {
	BusinessID *shared.BusinessID `json:"businessId" db:"business_id"`

	CreatedUserID shared.UserID `json:"createdUserId" db:"created_user_id"`
	// Document number
	Number *string `json:"number" db:"number"`

	// Document type
	DocType *string `json:"docType" db:"doc_type"`

	// Issuing state
	IssuingState *string `json:"issuingState" db:"issuing_state"`

	// Issuing country
	IssuingCountry string `json:"issuingCountry" db:"issuing_country"`

	// Issuing date
	IssuedDate *shared.Date `json:"issuedDate" db:"issued_date"`

	// Expiration Date
	ExpirationDate *shared.Date `json:"expirationDate" db:"expiration_date"`

	// Mime type
	ContentType string `json:"contentType" db:"content_type"`

	// Storage key reference
	StorageKey *string `json:"storage_key,omitempty" db:"storage_key"`

	// File contents
	Content *string `json:"content"`

	//Use this document for formation
	UseFormation *bool `json:"useFormation"`
}

//UserDocumentResponse  ..
type UserDocumentResponse struct {
	SignedURL *string      `json:"uploadURL"`
	Document  UserDocument `json:"document"`
}

type UserDocument struct {
	BaseDocument
	// Document id
	ID shared.UserDocumentID `json:"id" db:"id"`

	UserID string `json:"userId" db:"user_id"`
}

type UserDocumentCreate struct {

	//UserId ..
	UserID *string `json:"userId" db:"user_id"`

	// Document number
	Number *string `json:"number" db:"number"`

	// Document type
	DocType *string `json:"docType" db:"doc_type"`

	// Issuing state
	IssuingState *string `json:"issuingState" db:"issuing_state"`

	// Issuing country
	IssuingCountry string `json:"issuingCountry" db:"issuing_country"`

	// Issuing date
	IssuedDate shared.Date `json:"issuedDate" db:"issued_date"`

	// Expiration date
	ExpirationDate *shared.Date `json:"expirationDate" db:"expiration_date"`

	// Mime type
	ContentType string `json:"contentType" db:"content_type"`

	//document storage key
	StorageKey *string `json:"storage_key,omitempty" db:"storage_key"`
}

// UserDocumentUpdate user document update
type UserDocumentUpdate struct {
	// Document number
	Number *string `json:"number" db:"number"`

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

	// flag to know if content is being updated, if so send back a signed url
	UpdatingContent *bool `json:"updatingContent"`

	//document storage key
	StorageKey *string `json:"storage_key,omitempty" db:"storage_key"`
}

// ConsumerDocumentCreate ..
type ConsumerDocumentCreate struct {

	//UserId ..
	ConsumerID *shared.ConsumerID `json:"consumerId" db:"consumer_id"`

	// Document number
	Number *string `json:"number" db:"number"`

	// Document type
	DocType *docs.ConsumerIdentityDocument `json:"docType" db:"doc_type"`

	// Issuing state
	IssuingState *string `json:"issuingState" db:"issuing_state"`

	// Issuing country
	IssuingCountry string `json:"issuingCountry" db:"issuing_country"`

	// Issuing date
	IssuedDate shared.Date `json:"issuedDate" db:"issued_date"`

	// Expiration date
	ExpirationDate *shared.Date `json:"expirationDate" db:"expiration_date"`

	// Mime type
	ContentType string `json:"contentType" db:"content_type"`

	//document storage key
	StorageKey *string `json:"storage_key,omitempty" db:"storage_key"`
}

// ConsumerDocumentUpdate ..
type ConsumerDocumentUpdate struct {
	// Document number
	Number *string `json:"number" db:"number"`

	// Document type
	DocType *docs.ConsumerIdentityDocument `json:"docType" db:"doc_type"`

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

	// flag to know if content is being updated, if so send back a signed url
	UpdatingContent *bool `json:"updatingContent"`

	//document storage key
	StorageKey *string `json:"storage_key,omitempty" db:"storage_key"`
}

// ConsumerDocumentResponse  ..
type ConsumerDocumentResponse struct {
	SignedURL *string          `json:"uploadURL"`
	Document  ConsumerDocument `json:"document"`
}

// ConsumerDocument ..
type ConsumerDocument struct {
	BaseDocument
	// Document id
	ID shared.ConsumerDocumentID `json:"id" db:"id"`

	ConsumerID shared.ConsumerID `json:"consumerId" db:"consumer_id"`
}

//Status document status
type Status struct {
	ID             string          `json:"id" db:"id"`
	DocumentID     *string         `json:"documentId" db:"document_id"`
	Submitted      *time.Time      `json:"submitted" db:"submitted"`
	Status         string          `json:"status" db:"document_status"`
	Response       *types.JSONText `json:"response" db:"response"`
	Created        time.Time       `json:"created" db:"created"`
	BankDocumentID *string         `json:"bankDocumentId" db:"bank_document_id"`
	Modified       time.Time       `json:"modified" db:"modified"`
}

// CSPBusinessDocument ...
type CSPBusinessDocument struct {
	ID             *string                    `json:"id" db:"id"`
	DocumentID     *shared.BusinessDocumentID `json:"documentId" db:"document_id"`
	Submitted      *time.Time                 `json:"submitted" db:"submitted"`
	Status         *csp.DocumentStatus        `json:"status" db:"document_status"`
	Response       *types.JSONText            `json:"response" db:"response"`
	Created        *time.Time                 `json:"created" db:"created"`
	BankDocumentID *string                    `json:"bankDocumentId" db:"bank_document_id"`
	Modified       *time.Time                 `json:"modified" db:"modified"`
}

// CSPConsumerDocument ..
type CSPConsumerDocument struct {
	ID             *string                    `json:"id" db:"id"`
	DocumentID     *shared.ConsumerDocumentID `json:"documentId" db:"document_id"`
	Submitted      *time.Time                 `json:"submitted" db:"submitted"`
	Status         *csp.DocumentStatus        `json:"status" db:"document_status"`
	Response       *types.JSONText            `json:"response" db:"response"`
	Created        *time.Time                 `json:"created" db:"created"`
	BankDocumentID *string                    `json:"bankDocumentId" db:"bank_document_id"`
	Modified       *time.Time                 `json:"modified" db:"modified"`
}

const (
	driversLicenseDocType = "driversLicense"
	passportDocType       = "passport"
)
