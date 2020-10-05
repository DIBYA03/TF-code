package mock

import (
	"time"

	"github.com/wiseco/core-platform/services"
	platformservices "github.com/wiseco/core-platform/services"
	usersrv "github.com/wiseco/core-platform/services/user"
)

func NewUserResidency() platformservices.Residency {
	return platformservices.Residency{
		Country: "US",
		Status:  platformservices.ResidencyStatusCitizen,
	}
}

// Returns a mock user notification object
func NewUserNotification() usersrv.UserNotification {
	return usersrv.UserNotification{
		Transactions: makeBool(true),
		Transfers:    makeBool(true),
		Contacts:     makeBool(true),
	}
}
func makeBool(v bool) *bool {
	return &v
}

// Returns a mock user object
func NewUser(userId string) usersrv.User {

	email := "joe.smith@example.com"
	now := time.Now()
	dateNow := services.Date(now)
	taxId := services.TaxID("XXXXX3333")
	taxIdType := platformservices.TinTypeSSN
	occ := platformservices.OccTypeProfessionalManagement
	residency := NewUserResidency()
	address := NewAddress()

	return usersrv.User{
		Id:             userId,
		FirstName:      "Joe",
		MiddleName:     "L",
		LastName:       "Smith",
		Email:          &email,
		EmailVerified:  true,
		Phone:          "+18665551212",
		PhoneVerified:  true,
		KYCStatus:      platformservices.KYCStatusApproved,
		DateOfBirth:    &dateNow,
		TaxId:          &taxId,
		TaxIdType:      &taxIdType,
		LegalAddress:   &address,
		MailingAddress: nil,
		UserResidency:  &residency,
		Occupation:     &occ,
		Notification:   NewUserNotification(),
		Created:        now,
		Updated:        now,
	}
}
