package main

import (
	"context"
	"log"
	"sync"
	"time"

	grpcTypes "github.com/golang/protobuf/ptypes"
	"github.com/wiseco/core-platform/services"
	"github.com/wiseco/core-platform/services/data"
	"github.com/wiseco/core-platform/services/user"
	"github.com/wiseco/core-platform/shared"
	"github.com/wiseco/go-lib/id"
	grpcRoot "github.com/wiseco/protobuf/golang"
	grpcMonitor "github.com/wiseco/protobuf/golang/transaction/monitor"
	"github.com/xtgo/uuid"
)

var KYCStatusToConsumerProto = map[services.KYCStatus]grpcRoot.ConsumerKYCStatus{
	services.KYCStatusNotStarted: grpcRoot.ConsumerKYCStatus_CKS_NOT_STARTED,
	services.KYCStatusSubmitted:  grpcRoot.ConsumerKYCStatus_CKS_SUBMITTED,
	services.KYCStatusReview:     grpcRoot.ConsumerKYCStatus_CKS_REVIEW,
	services.KYCStatusApproved:   grpcRoot.ConsumerKYCStatus_CKS_APPROVED,
	services.KYCStatusDeclined:   grpcRoot.ConsumerKYCStatus_CKS_DECLINED,
}

var OccupationToProto = map[string]grpcRoot.ConsumerOccupation{
	services.OccTypeAgriculture:                grpcRoot.ConsumerOccupation_CO_AGRICULTURE,
	services.OccTypeClergyMinistryStaff:        grpcRoot.ConsumerOccupation_CO_CLERGY_MINISTRY_STAFF,
	services.OccTypeConstructionIndustrial:     grpcRoot.ConsumerOccupation_CO_CONSTRUCTION_INDUSTRIAL,
	services.OccTypeEducation:                  grpcRoot.ConsumerOccupation_CO_EDUCATION,
	services.OccTypeFinanceAccountingTax:       grpcRoot.ConsumerOccupation_CO_FINANCE_ACCOUNTING_TAX,
	services.OccTypeFireFirstResponders:        grpcRoot.ConsumerOccupation_CO_FIRE_FIRST_RESPONDERS,
	services.OccTypeHealthcare:                 grpcRoot.ConsumerOccupation_CO_HEALTHCARE,
	services.OccTypeHomemaker:                  grpcRoot.ConsumerOccupation_CO_HOMEMAKER,
	services.OccTypeLaborGeneral:               grpcRoot.ConsumerOccupation_CO_LABOR_GENERAL,
	services.OccTypeLaborSkilled:               grpcRoot.ConsumerOccupation_CO_LABOR_SKILLED,
	services.OccTypeLawEnforcementSecurity:     grpcRoot.ConsumerOccupation_CO_LAW_ENFORCEMENT_SECURITY,
	services.OccTypeLegalServices:              grpcRoot.ConsumerOccupation_CO_LEGAL_SERVICES,
	services.OccTypeMilitary:                   grpcRoot.ConsumerOccupation_CO_MILITARY,
	services.OccTypeNotaryRegistrar:            grpcRoot.ConsumerOccupation_CO_NOTARY_REGISTRAR,
	services.OccTypePrivateInvestor:            grpcRoot.ConsumerOccupation_CO_PRIVATE_INVESTOR,
	services.OccTypeProfessionalAdministrative: grpcRoot.ConsumerOccupation_CO_PROFESSIONAL_ADMINISTRATIVE,
	services.OccTypeProfessionalManagement:     grpcRoot.ConsumerOccupation_CO_PROFESSIONAL_MANAGEMENT,
	services.OccTypeProfessionalOther:          grpcRoot.ConsumerOccupation_CO_PROFESSIONAL_OTHER,
	services.OccTypeProfessionalTechnical:      grpcRoot.ConsumerOccupation_CO_PROFESSIONAL_TECHNICAL,
	services.OccTypeRetired:                    grpcRoot.ConsumerOccupation_CO_RETIRED,
	services.OccTypeSales:                      grpcRoot.ConsumerOccupation_CO_SALES,
	services.OccTypeSelfEmployed:               grpcRoot.ConsumerOccupation_CO_SELF_EMPLOYED,
	services.OccTypeStudent:                    grpcRoot.ConsumerOccupation_CO_STUDENT,
	services.OccTypeTransportation:             grpcRoot.ConsumerOccupation_CO_TRANSPORTATION,
	services.OccTypeUnemployed:                 grpcRoot.ConsumerOccupation_CO_UNEMPLOYED,
}

func processConsumer(monitorClient grpcMonitor.BankTransactionMonitorServiceClient, cID id.ConsumerID) error {
	sharedConID := shared.ConsumerID(uuid.UUID(cID).String())
	c, err := user.NewConsumerServiceWithout().GetByID(sharedConID)
	if err != nil {
		log.Println(err, cID.String())
		return err
	}

	kycStatus, ok := KYCStatusToConsumerProto[c.KYCStatus]
	if !ok {
		log.Printf("Invalid KYC Status: %s %s", c.KYCStatus, cID.String())
		return nil
	}

	if c.KYCStatus != services.KYCStatusApproved {
		log.Printf("Only approved customers supported")
		return nil
	}

	var la *grpcRoot.Address
	if c.LegalAddress != nil {
		la = &grpcRoot.Address{
			Line_1:    c.LegalAddress.StreetAddress,
			Locality:  c.LegalAddress.City,
			AdminArea: c.LegalAddress.State,
			// TODO: Fix errors in address data
			Country:    "US",
			PostalCode: c.LegalAddress.PostalCode,
		}
	}
	created, err := grpcTypes.TimestampProto(c.Created)
	if err != nil {
		log.Println(err, cID.String())
		return err
	}

	modified, err := grpcTypes.TimestampProto(c.Modified)
	if err != nil {
		log.Println(err, cID.String())
		return err
	}

	status := grpcRoot.ConsumerStatus_CS_ACTIVE
	if c.Deactivated != nil {
		status = grpcRoot.ConsumerStatus_CS_INACTIVE
	}

	citizen := "us"
	if c.CitizenshipCountries != nil && len(c.CitizenshipCountries) > 0 && len(c.CitizenshipCountries[0]) == 2 {
		citizen = c.CitizenshipCountries[0]
	}

	occ := grpcRoot.ConsumerOccupation_CO_UNSPECIFIED
	if c.Occupation != nil {
		occ, ok = OccupationToProto[*c.Occupation]
		if !ok {
			log.Printf("Invalid Occupation: %s %s", *c.Occupation, cID.String())
			return nil
		}
	}

	creq := &grpcMonitor.ConsumerRequest{
		Id:                 cID.String(),
		FirstName:          c.FirstName,
		MiddleName:         c.MiddleName,
		LastName:           c.LastName,
		FullName:           c.FullName(),
		Occupation:         occ,
		KycStatus:          kycStatus,
		LegalAddress:       la,
		CitizenshipCountry: citizen,
		Email:              shared.StringValue(c.Email),
		Mobile:             shared.StringValue(c.Phone),
		Status:             status,
		DateOfBirth:        c.DateOfBirth.String(),
		Created:            created,
		Modified:           modified,
	}

	resp, err := monitorClient.AddUpdateConsumer(context.Background() /* client.GetContext() */, creq)
	if err != nil {
		log.Println(err, cID.String())
		return err
	}

	log.Println("Success: ", resp.Id)
	return nil
}

func sendConsumerUpdates(monitorClient grpcMonitor.BankTransactionMonitorServiceClient, dayStart, dayEnd time.Time) {
	// Process in groups of 5
	offset := 0
	limit := 5
	for {
		var consumerIDs []id.ConsumerID
		err := data.DBWrite.Select(
			&consumerIDs,
			`
			SELECT id from consumer
			WHERE
				(created >= $1 AND created < $2) OR
				(modified >= $1 AND modified < $2)
			ORDER BY created ASC OFFSET $3 LIMIT $4`,
			dayStart,
			dayEnd,
			offset,
			limit,
		)
		if err != nil {
			panic(err)
		} else if len(consumerIDs) == 0 {
			log.Println("No more consumers", dayStart, dayEnd)
			break
		}

		wg := sync.WaitGroup{}
		wg.Add(len(consumerIDs))
		for _, cID := range consumerIDs {
			go func(id id.ConsumerID) {
				defer wg.Done()
				_ = processConsumer(monitorClient, id)
			}(cID)
		}

		wg.Wait()
		offset += 5
	}
}
