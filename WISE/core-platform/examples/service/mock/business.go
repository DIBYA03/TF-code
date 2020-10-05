package mock

import (
	"time"

	"github.com/google/uuid"
	"github.com/wiseco/core-platform/services"
	bussrv "github.com/wiseco/core-platform/services/business"
)

func NewBusinessMember() bussrv.BusinessMember {
	uId := uuid.New().String()
	return bussrv.BusinessMember{
		Id:                   uuid.New().String(),
		UserId:               &uId,
		KYCStatus:            services.KYCStatusApproved,
		TitleType:            bussrv.TitleTypePresident,
		TitleOther:           nil,
		Ownership:            25.0,
		IsControllingManager: true,
	}
}

// Creates a mock business object
func NewBusiness(businessId string) bussrv.Business {

	entityType := bussrv.EntityTypeUnlistedCorporation
	industryType := bussrv.IndustryTypeHotelMotel
	email := "ownerlbri@example.com"
	phone := "+15625551212"
	address := NewAddress()
	taxIdMasked := services.TaxID("XXXXX9999")
	taxIdType := services.TinTypeEIN
	originCountry := "US"
	originState := "CA"
	now := time.Now()
	dateNow := services.Date(now)
	purpose := "To provide accomodations and related services."

	return bussrv.Business{
		Id:             businessId,
		EmployerNumber: "RI24156",
		LegalName:      "Rodeway Inn Long Beach",
		EntityType:     &entityType,
		IndustryType:   &industryType,
		TaxId:          &taxIdMasked,
		TaxIdType:      &taxIdType,
		OriginCountry:  &originCountry,
		OriginState:    &originState,
		OriginDate:     &dateNow,
		Purpose:        &purpose,
		KYCStatus:      services.KYCStatusApproved,
		Members:        []bussrv.BusinessMember{NewBusinessMember()},
		Email:          &email,
		EmailVerified:  true,
		Phone:          &phone,
		PhoneVerified:  true,
		LegalAddress:   &address,
		MailingAddress: nil,
		Created:        time.Now(),
		Updated:        time.Now(),
	}
}

func NewBusinessAccess(businessId string) bussrv.BusinessAccess {
	return bussrv.BusinessAccess{
		Id:         uuid.New().String(),
		BusinessId: businessId,
		UserId:     uuid.New().String(),
		AccessType: bussrv.BusinessAccessTypeAdmin,
		AccessRole: bussrv.BusinessAccessRoleOfficer,
	}
}
