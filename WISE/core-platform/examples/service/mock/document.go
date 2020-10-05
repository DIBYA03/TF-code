package mock

import (
	"time"

	"github.com/google/uuid"
	"github.com/wiseco/core-platform/services"
	docsrv "github.com/wiseco/core-platform/services/document"
)

// NewDocument Creates a mock document object
func NewUserDocument(documentId, userId string) docsrv.UserDocument {

	now := time.Now()
	dateNow := services.Date(now)
	IssuingState := "CA"
	contentType := "application/pdf"

	return docsrv.UserDocument{
		OwnerUserId: userId,
		Document: docsrv.Document{
			Id:             documentId,
			Number:         "XYZ123456",
			DocType:        docsrv.UserDocTypeDriversLicense,
			IssuingAuth:    "California Department of Motor Vehicles",
			IssuingState:   &IssuingState,
			IssuingCountry: "US",
			IssuedDate:     dateNow,
			ExpirationDate: &dateNow,
			ContentType:    &contentType,
			Created:        now,
			Updated:        now,
		},
	}
}

// NewBusinessDocument creates a mock document object
func NewBusinessDocument(documentId string, businessId string) docsrv.BusinessDocument {

	now := time.Now()
	dateNow := services.Date(now)
	issuingState := "CA"
	contentType := "application/pdf"

	return docsrv.BusinessDocument{
		BusinessID:    businessId,
		CreatedUserID: uuid.New().String(),
		Document: docsrv.Document{
			Id:             documentId,
			Number:         "XYZ123456",
			DocType:        docsrv.BusinessDocTypeCertOfFormation.String(),
			IssuingAuth:    "California SOS",
			IssuingState:   &issuingState,
			IssuingCountry: "US",
			IssuedDate:     dateNow,
			ExpirationDate: &dateNow,
			ContentType:    &contentType,
			Created:        now,
			Updated:        now,
		},
	}
}
