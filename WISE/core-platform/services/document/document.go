/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

package document

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/wiseco/core-platform/services"
	"github.com/wiseco/core-platform/shared"
)

// General document number (DL #, Passport, etc.)
type DocumentNumber string

func (n *DocumentNumber) String() string {
	return string(*n)
}

func (n *DocumentNumber) MarshalJSON() ([]byte, error) {
	return json.Marshal(services.MaskLeft(n.String(), 4))
}

type Document struct {
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
	ContentType string `json:"contentType" db:"content_type"`

	// Created timestamp
	Created time.Time `json:"created" db:"created"`

	//Deleted
	Deleted *time.Time `json:"deleted" db:"deleted"`

	//ContentUploaded indicates content on s3 has been uploaded
	ContentUploaded *time.Time `json:"contentUploaded" db:"content_uploaded"`

	// Update timestamp
	Modified time.Time `json:"modified" db:"modified"`

	// Storage key reference
	StorageKey *string `json:"storage_key,omitempty" db:"storage_key"`
}

type DocumentCreate struct {
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

	// Document storage key
	StorageKey *string `json:"storage_key,omitempty" db:"storage_key"`
}

// If contents exist a new document is created and the old is removed from the database
// Old doc replaced in database - S3 is WORM mode (Write-Once/Read Many)
type DocumentUpdate struct {
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

func (document *Document) Storer() (Storer, error) {
	if document.StorageKey == nil {
		return nil, errors.New("Storage key does not exist")
	}

	storageService, err := NewStorerFromKey(*document.StorageKey)
	if err == nil {
		return nil, err
	}

	return storageService, nil
}
