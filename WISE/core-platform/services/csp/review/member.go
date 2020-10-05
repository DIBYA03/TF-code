/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

package review

import (
	"errors"
	"log"
	"time"

	partnerbank "github.com/wiseco/core-platform/partner/bank"
	"github.com/wiseco/core-platform/services"
	biz "github.com/wiseco/core-platform/services/business"
	bsrv "github.com/wiseco/core-platform/services/business"
	"github.com/wiseco/core-platform/services/csp"
	"github.com/wiseco/core-platform/services/csp/business"
	"github.com/wiseco/core-platform/services/csp/consumer"
	usrv "github.com/wiseco/core-platform/services/user"
	"github.com/wiseco/core-platform/shared"
)

func (service verification) getPartnerBusinessMembers(id shared.BusinessID) ([]partnerbank.BusinessMemberRequest, error) {
	var members []biz.BusinessMember
	var requests []partnerbank.BusinessMemberRequest
	err := service.Select(
		&members, `
		SELECT
            business_member.id, business_member.consumer_id, business_member.business_id,
            business_member.title_type, business_member.title_other, business_member.ownership,
            business_member.is_controlling_manager, business_member.deactivated, business_member.created,
            business_member.modified, consumer.first_name, consumer.middle_name, consumer.last_name,
            consumer.date_of_birth, consumer.tax_id, consumer.tax_id_type, consumer.kyc_status,
            consumer.legal_address, consumer.mailing_address, consumer.work_address, consumer.residency,
			consumer.citizenship_countries, consumer.occupation, consumer.income_type,
			consumer.phone, consumer.email,
            consumer.activity_type, consumer.is_restricted
        FROM
            business_member
        JOIN
            consumer ON business_member.consumer_id = consumer.id
        WHERE
            business_member.deactivated IS NULL AND business_member.business_id = $1`,
		id,
	)
	if err != nil {
		return requests, err
	}

	ownership := 0
	isManager := false
	for _, member := range members {
		title, ok := biz.PartnerTitleTypeFrom[member.TitleType]
		if !ok {
			return requests, errors.New("invalid member title type")
		}
		// All members should have approved kyc status
		if member.KYCStatus != services.KYCStatusApproved {
			return requests, errors.New("Kyc status should be approved for all members")

		}
		// Ownership percentage should not exceed 100%
		ownership = ownership + member.Ownership
		if ownership > 100 {
			return requests, errors.New("Total ownership percentage cannot exceed 100%")
		}

		//- Minimum of one(and only one) control manager is required
		if member.IsControllingManager {
			if isManager {
				return requests, errors.New("Only one control manager is allowed")
			}
			isManager = member.IsControllingManager
		}

		requests = append(
			requests,
			partnerbank.BusinessMemberRequest{
				ConsumerID:           partnerbank.ConsumerID(member.ConsumerID),
				IsControllingManager: member.IsControllingManager,
				Ownership:            member.Ownership,
				Title:                title,
				TitleDesc:            member.TitleOther,
			},
		)
	}
	//- Minimum of one(and only one) control manager is required
	if !isManager {
		return requests, errors.New("Atleast one control manager is required")
	}

	return requests, nil
}

// VerifyMembers will verify all the memebers
func VerifyMembers(ownerID shared.UserID, businessID shared.BusinessID, consumerID shared.ConsumerID) (bool, error) {
	hasReview := false

	// Fetch owner/user
	sr := services.NewSourceRequest()
	sr.UserID = ownerID
	u, err := usrv.NewUserService(sr).GetByIdInternal(ownerID)
	if err != nil {
		log.Printf("Error fetching business owner: %v", err)
		return true, errors.New("Error verifying business owner")
	}

	// Fetch members
	members, err := bsrv.NewMemberService(sr).List(0, 20, businessID)
	if err != nil {
		log.Printf("Error fetching members: %v", err)
		return true, errors.New("Error verifying member list")
	}

	/* Disable automated KYC to allow manual verify
	sr = services.NewSourceRequest()
	sr.UserID = ownerID
	ownerResp, err := usrv.NewUserService(sr).StartVerification(ownerID)
	if err != nil {
		hasReview = true
		log.Println(err)
		createConsumerUser(u, services.KYCStatusSubmitted, []string{})
	} else {
		if ownerResp.Status == services.KYCStatusReview || ownerResp.Status == services.KYCStatusDeclined {
			hasReview = true
		}

		createConsumerUser(u, ownerResp.Status, ownerResp.ReviewItems)
	} */

	hasReview = true
	createConsumerUser(u, services.KYCStatusSubmitted, []string{})

	// Verify all non-owner members
	for _, m := range members {
		// Skip owner
		if m.ConsumerID == consumerID {
			continue
		}

		/* Disable automated KYC to allow manual verify
		sr = services.NewSourceRequest()
		sr.UserID = ownerID
		resp, err := bsrv.NewMemberService(sr).StartVerification(m.ID, businessID)
		if err != nil {
			log.Printf("Error verifying member %v", err)
			createConsumerMember(m, services.KYCStatusSubmitted, []string{})
		} else {
			if resp.Status == services.KYCStatusReview || resp.Status == services.KYCStatusDeclined {
				hasReview = true
			}

			createConsumerMember(m, resp.Status, resp.ReviewItems)
		} */

		hasReview = true
		createConsumerMember(m, services.KYCStatusSubmitted, []string{})
	}

	// set core business kyc status to review
	if err := business.New(services.SourceRequest{}).UpdateKYC(businessID, csp.KYCStatusReview); err != nil {
		log.Printf("Error updating business KYC %v", err)
		return true, err
	}

	return hasReview, err
}

func createConsumerUser(u *usrv.User, status services.KYCStatus, items []string) error {
	name := u.Name()
	submitted := time.Now()
	var idvs services.StringArray
	for _, item := range items {
		idvs = append(idvs, item)
	}
	create := consumer.CSPConsumerCreate{
		ConsumerName: &name,
		ConsumerID:   u.ConsumerID,
		Status:       string(status),
		IDVs:         &idvs,
		Submitted:    &submitted,
	}
	_, err := consumer.NewCSPService().CSPConsumerCreate(create)
	return err
}

func createConsumerMember(m bsrv.BusinessMember, status services.KYCStatus, items []string) error {
	name := m.Name()
	submitted := time.Now()
	var idvs services.StringArray
	for _, item := range items {
		idvs = append(idvs, item)
	}
	create := consumer.CSPConsumerCreate{
		ConsumerName: &name,
		ConsumerID:   m.ConsumerID,
		Status:       string(status),
		IDVs:         &idvs,
		Submitted:    &submitted,
	}
	_, err := consumer.NewCSPService().CSPConsumerCreate(create)
	return err
}
