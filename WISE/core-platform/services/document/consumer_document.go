/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

package document

import (
	partnerbank "github.com/wiseco/core-platform/partner/bank"
	"github.com/wiseco/core-platform/shared"
)

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

var ConsumerDocTypeToBankDocType = map[ConsumerIdentityDocument]partnerbank.ConsumerIdentityDocument{
	ConsumerIdentityDocumentDriversLicense:        partnerbank.ConsumerIdentityDocumentDriversLicense,
	ConsumerIdentityDocumentPassport:              partnerbank.ConsumerIdentityDocumentPassport,
	ConsumerIdentityDocumentPassportCard:          partnerbank.ConsumerIdentityDocumentPassportCard,
	ConsumerIdentityDocumentWorkPermit:            partnerbank.ConsumerIdentityDocumentWorkPermit,
	ConsumerIdentityDocumentSocialSecurityCard:    partnerbank.ConsumerIdentityDocumentSocialSecurityCard,
	ConsumerIdentityDocumentStateID:               partnerbank.ConsumerIdentityDocumentStateID,
	ConsumerIdentityDocumentAlienRegistrationCard: partnerbank.ConsumerIdentityDocumentAlienRegistrationCard,
	ConsumerIdentityDocumentUSAVisaH1B:            partnerbank.ConsumerIdentityDocumentUSAVisaH1B,
	ConsumerIdentityDocumentUSAVisaH1C:            partnerbank.ConsumerIdentityDocumentUSAVisaH1C,
	ConsumerIdentityDocumentUSAVisaH2A:            partnerbank.ConsumerIdentityDocumentUSAVisaH2A,
	ConsumerIdentityDocumentUSAVisaH2B:            partnerbank.ConsumerIdentityDocumentUSAVisaH2B,
	ConsumerIdentityDocumentUSAVisaH3:             partnerbank.ConsumerIdentityDocumentUSAVisaH3,
	ConsumerIdentityDocumentUSAVisaL1A:            partnerbank.ConsumerIdentityDocumentUSAVisaL1A,
	ConsumerIdentityDocumentUSAVisaL1B:            partnerbank.ConsumerIdentityDocumentUSAVisaL1B,
	ConsumerIdentityDocumentUSAVisaO1:             partnerbank.ConsumerIdentityDocumentUSAVisaO1,
	ConsumerIdentityDocumentUSAVisaE1:             partnerbank.ConsumerIdentityDocumentUSAVisaE1,
	ConsumerIdentityDocumentUSAVisaE3:             partnerbank.ConsumerIdentityDocumentUSAVisaE3,
	ConsumerIdentityDocumentUSAVisaI:              partnerbank.ConsumerIdentityDocumentUSAVisaI,
	ConsumerIdentityDocumentUSAVisaP:              partnerbank.ConsumerIdentityDocumentUSAVisaP,
	ConsumerIdentityDocumentUSAVisaTN:             partnerbank.ConsumerIdentityDocumentUSAVisaTN,
	ConsumerIdentityDocumentUSAVisaTD:             partnerbank.ConsumerIdentityDocumentUSAVisaTD,
	ConsumerIdentityDocumentUSAVisaR1:             partnerbank.ConsumerIdentityDocumentUSAVisaR1,
}

type ConsumerDocument struct {
	// Document id
	ID         shared.ConsumerDocumentID `json:"id" db:"id"`
	ConsumerID shared.ConsumerID         `json:"consumerId" db:"consumer_id"`
	Document
}

//UserDocumentResponse  ..
type ConsumerDocumentResponse struct {
	SignedURL *string          `json:"uploadURL"`
	Document  ConsumerDocument `json:"document"`
}

type ConsumerDocumentCreate struct {
	// Document id - if we need to duplicate already uploaded doc
	ID *shared.ConsumerDocumentID `json:"id" db:"id"`

	//UserID
	ConsumerID *shared.ConsumerID `json:"consumerId" db:"consumer_id"`

	// Document number
	Number *DocumentNumber `json:"number" db:"number"`

	// Document type
	DocType *ConsumerIdentityDocument `json:"docType" db:"doc_type"`

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

type ConsumerDocumentUpdate struct {
	// Document id
	ID shared.ConsumerDocumentID `json:"id" db:"id"`

	// Document number
	Number *DocumentNumber `json:"number" db:"number"`

	// Document type
	DocType *ConsumerIdentityDocument `json:"docType" db:"doc_type"`

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
