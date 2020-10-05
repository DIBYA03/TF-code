/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

package business

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	partnerbank "github.com/wiseco/core-platform/partner/bank"
	"github.com/wiseco/core-platform/services"
	mbsrv "github.com/wiseco/core-platform/services/business"
	"github.com/wiseco/core-platform/services/csp/intercom"
	"github.com/wiseco/core-platform/services/data"
	coreDB "github.com/wiseco/core-platform/services/data"
	"github.com/wiseco/core-platform/shared"
	goLibClear "github.com/wiseco/go-lib/clear"
	"github.com/wiseco/go-lib/grpc"
	golang "github.com/wiseco/protobuf/golang"
	"github.com/wiseco/protobuf/golang/verification/alloy"
	"github.com/wiseco/protobuf/golang/verification/clear"
	"github.com/wiseco/protobuf/golang/verification/email"
	"github.com/wiseco/protobuf/golang/verification/phone"
)

//MemberService ..
type MemberService interface {
	List(id shared.BusinessID, offset, limit int) ([]Member, error)
	GetByID(memberID shared.BusinessMemberID, businessID shared.BusinessID) (*Member, error)
	Create(businessID shared.BusinessID, create mbsrv.BusinessMemberCreate) (*mbsrv.BusinessMember, error)
	Update(memberID shared.BusinessMemberID, businessID shared.BusinessID, update mbsrv.BusinessMemberUpdate) (*mbsrv.BusinessMember, error)
	Submit(memberID shared.BusinessMemberID, businessID shared.BusinessID) (*mbsrv.BusinessMember, error)
	StartVerification(memberID shared.BusinessMemberID, businessID shared.BusinessID) (*mbsrv.MemberKYCResponse, error)
	GetVerification(memberID shared.BusinessMemberID, businessID shared.BusinessID) (*mbsrv.MemberKYCResponse, error)
	PhoneVerification(memberID shared.BusinessMemberID, businessID shared.BusinessID) (string, error)
	GetPhoneVerification(memberID shared.BusinessMemberID, businessID shared.BusinessID) (string, error)
	EmailVerification(memberID shared.BusinessMemberID, businessID shared.BusinessID) (EmailVerification, error)
	GetEmailVerification(memberID shared.BusinessMemberID, businessID shared.BusinessID) (EmailVerification, error)
	RunAlloyKYC(memberID shared.BusinessMemberID, businessID shared.BusinessID) (string, error)
	RunClearKYC(memberID shared.BusinessMemberID, businessID shared.BusinessID, kycType goLibClear.ClearKycType) (string, error)
	GetAlloyKYC(memberID shared.BusinessMemberID, businessID shared.BusinessID) (string, error)
	GetClearKYC(memberID shared.BusinessMemberID, businessID shared.BusinessID, kycType goLibClear.ClearKycType) (string, error)

	RunKYC(string) (*KYCResult, error)
	GetConsumerIDsForBusinessIDsInternal(businessIDs []string) ([]string, error)
}

type memberService struct {
	*sqlx.DB
	sourceReq services.SourceRequest
}

//NewMemberService a new business member service
func NewMemberService(sourceReq services.SourceRequest) MemberService {
	return memberService{data.DBWrite, sourceReq}
}

func (s memberService) List(id shared.BusinessID, offset, limit int) ([]Member, error) {

	var list = make([]Member, 0)
	err := s.Select(&list, `
		SELECT
            business_member.id, business_member.consumer_id, business_member.business_id,
            business_member.title_type, business_member.title_other, business_member.ownership,
            business_member.is_controlling_manager, business_member.deactivated, business_member.created,
            business_member.modified, consumer.first_name, consumer.middle_name, consumer.last_name,
            consumer.date_of_birth, consumer.tax_id AS tax_id_unmasked, consumer.tax_id_type, consumer.kyc_status,
            consumer.legal_address, consumer.mailing_address, consumer.work_address, consumer.residency,
			consumer.citizenship_countries, consumer.occupation, consumer.income_type,
			consumer.phone, consumer.email,
            consumer.activity_type, consumer.is_restricted, wise_user.id as user_id, wise_user.phone as user_phone
        FROM
            business_member
        JOIN
			consumer ON business_member.consumer_id = consumer.id
		LEFT JOIN
            wise_user ON business_member.consumer_id = wise_user.consumer_id
        WHERE
            business_member.deactivated IS NULL AND business_member.business_id = $1`,
		id)

	if err != nil && err == sql.ErrNoRows {
		return list, nil
	}
	return list, err
}

func (s memberService) GetByID(memberID shared.BusinessMemberID, businessID shared.BusinessID) (*Member, error) {
	var m Member
	err := s.Get(&m, `
	SELECT
		business_member.id, business_member.consumer_id, business_member.business_id,
		business_member.title_type, business_member.title_other, business_member.ownership,
		business_member.is_controlling_manager, business_member.deactivated, business_member.created,
		business_member.modified, consumer.first_name, consumer.middle_name, consumer.last_name,
		consumer.email, consumer.phone, consumer.date_of_birth, consumer.tax_id AS tax_id_unmasked, consumer.tax_id_type,
		consumer.kyc_status, consumer.legal_address, consumer.mailing_address, consumer.work_address,
		consumer.residency, consumer.citizenship_countries, consumer.occupation, consumer.income_type,
		consumer.activity_type, consumer.is_restricted, wise_user.id as user_id, wise_user.phone as user_phone
	FROM
		business_member
	LEFT JOIN
		consumer ON business_member.consumer_id = consumer.id
	LEFT JOIN
		wise_user ON business_member.consumer_id = wise_user.consumer_id
	WHERE
		business_member.id = $1 AND business_member.business_id = $2`,
		memberID,
		businessID)

	if err != nil && err == sql.ErrNoRows {
		return nil, services.ErrorNotFound{}.New("")
	}

	if err != nil {
		return nil, err
	}

	bank, err := partnerbank.GetProxyBank(partnerbank.ProviderNameBBVA)
	if err != nil {
		log.Printf("Error getting partner bank %v", err)
	}
	CO, err := bank.ProxyService(s.sourceReq.PartnerBankRequest()).GetConsumerBankID(partnerbank.ConsumerID(m.ConsumerID))
	if err != nil && err != sql.ErrNoRows {
		log.Printf("error getting bankID %v", err)
	} else {
		m.BankID = CO
	}
	return &m, nil
}

func (s memberService) Create(businessID shared.BusinessID, create mbsrv.BusinessMemberCreate) (*mbsrv.BusinessMember, error) {
	sr := services.NewSourceRequest()
	sr.UserID = getOwnerID(businessID)
	resp, err := mbsrv.NewMemberService(sr).Create(&create)
	if err != nil {
		log.Printf("Error verifying member %v", err)
	}

	return resp, err
}

func (s memberService) Submit(memberID shared.BusinessMemberID, businessID shared.BusinessID) (*mbsrv.BusinessMember, error) {
	sr := services.NewSourceRequest()
	sr.UserID = getOwnerID(businessID)
	resp, err := mbsrv.NewMemberService(sr).Submit(memberID, businessID)
	if err != nil {
		log.Printf("Error verifying member %v", err)
	}

	return resp, err
}

func (s memberService) Update(memberID shared.BusinessMemberID, businessID shared.BusinessID, update mbsrv.BusinessMemberUpdate) (*mbsrv.BusinessMember, error) {
	sr := services.NewSourceRequest()
	sr.UserID = getOwnerID(businessID)
	update.ID = memberID
	resp, err := mbsrv.NewMemberService(sr).Update(memberID, businessID, &update)
	if err != nil {
		log.Printf("Error verifying member %v", err)
	}

	return resp, err
}

func (s memberService) StartVerification(memberID shared.BusinessMemberID, businessID shared.BusinessID) (*mbsrv.MemberKYCResponse, error) {

	s.sourceReq.UserID = getOwnerID(businessID)
	resp, err := mbsrv.NewMemberService(s.sourceReq).StartVerification(memberID, businessID)
	if err != nil {
		log.Printf("Error verifying member %v", err)
	}
	return resp, err
}

func (s memberService) GetVerification(memberID shared.BusinessMemberID, businessID shared.BusinessID) (*mbsrv.MemberKYCResponse, error) {

	sr := services.NewSourceRequest()
	sr.UserID = getOwnerID(businessID)
	resp, err := mbsrv.NewMemberService(sr).GetVerification(memberID, businessID)
	if err != nil {
		log.Printf("Error verifying member %v", err)
	}
	return resp, err
}

func getOwnerID(id shared.BusinessID) shared.UserID {
	var userID shared.UserID
	err := coreDB.DBRead.Get(&userID, "SELECT owner_id FROM business WHERE id = $1", id)
	if err != nil {
		log.Printf("error getting userId by businessId %v", err)
	}

	return userID
}

func (s memberService) PhoneVerification(memberID shared.BusinessMemberID, businessID shared.BusinessID) (string, error) {
	sn, err := grpc.GetConnectionStringForService(grpc.ServiceNameVerification)
	if err != nil {
		return "", err
	}
	client, err := grpc.NewInsecureClient(sn)
	if err != nil {
		return "", err
	}
	defer client.CloseAndCancel()
	c := phone.NewPhoneServiceClient(client.GetConn())

	member, err := s.GetByID(memberID, businessID)
	if err != nil {
		return "", err
	}

	vReq := &phone.VerificationRequest{
		PhoneNumber: member.Phone,
		ConsumerId:  string(member.ConsumerID),
	}
	res, err := c.Verification(context.Background(), vReq)

	if err != nil {
		return "", err
	}

	return res.Raw, nil
}

//TODO finish implementation
func (s memberService) GetPhoneVerification(memberID shared.BusinessMemberID, businessID shared.BusinessID) (string, error) {
	return "", nil
}

func (s memberService) EmailVerification(memberID shared.BusinessMemberID, businessID shared.BusinessID) (EmailVerification, error) {
	sn, err := grpc.GetConnectionStringForService(grpc.ServiceNameVerification)
	if err != nil {
		return EmailVerification{}, err
	}
	client, err := grpc.NewInsecureClient(sn)
	if err != nil {
		return EmailVerification{}, err
	}
	defer client.CloseAndCancel()
	c := email.NewEmailServiceClient(client.GetConn())
	member, err := s.GetByID(memberID, businessID)
	if err != nil {
		return EmailVerification{}, err
	}
	vReq := &email.VerificationRequest{
		//TODO @sandrews use correct email_id
		CoreEmailId:  uuid.New().String(),
		EmailAddress: member.Email,
	}
	res, err := c.Verification(client.GetContext(), vReq)
	if err != nil {
		return EmailVerification{}, err
	}
	return EmailVerification{
		Score:   res.Score,
		Verdict: res.Verdict,
	}, nil
}

//TODO finish implementation
func (s memberService) GetEmailVerification(memberID shared.BusinessMemberID, businessID shared.BusinessID) (EmailVerification, error) {
	return EmailVerification{}, nil
}

func (s memberService) RunAlloyKYC(memberID shared.BusinessMemberID, businessID shared.BusinessID) (string, error) {
	sn, err := grpc.GetConnectionStringForService(grpc.ServiceNameVerification)
	if err != nil {
		return "", err
	}
	client, err := grpc.NewInsecureClient(sn)
	if err != nil {
		return "", err
	}
	defer client.CloseAndCancel()
	c := alloy.NewAlloyServiceClient(client.GetConn())

	member, err := s.GetByID(memberID, businessID)
	if err != nil {
		return "", err
	}
	ph := strings.Trim(member.Phone, "+")
	if member.TaxID == nil {
		return "", errors.New("Tax id for member cannot be null")
	}
	vReq := &alloy.ConsumerVerificationRequest{
		ConsumerId:    string(member.ConsumerID),
		FirstName:     member.FirstName,
		LastName:      member.LastName,
		Email:         member.Email,
		Dob:           member.DateOfBirth.String(),
		AddressLine_1: member.LegalAddress.StreetAddress,
		City:          member.LegalAddress.City,
		State:         member.LegalAddress.State,
		PostalCode:    member.LegalAddress.PostalCode,
		Country:       member.LegalAddress.Country,
		Ssn:           *member.TaxID,
		Phone:         ph,
	}

	res, err := c.ConsumerVerification(context.Background(), vReq)
	if err != nil {
		log.Println(err)
		return "", err
	}

	return string(res.ResultData), nil
}

func (s memberService) RunClearKYC(memberID shared.BusinessMemberID, businessID shared.BusinessID, kycType goLibClear.ClearKycType) (string, error) {
	log.Println("In RunClearKYC")
	sn, err := grpc.GetConnectionStringForService(grpc.ServiceNameVerification)
	if err != nil {
		return "", err
	}
	client, err := grpc.NewInsecureClient(sn)
	if err != nil {
		return "", err
	}
	defer client.CloseAndCancel()
	c := clear.NewClearServiceClient(client.GetConn())

	member, err := s.GetByID(memberID, businessID)
	if err != nil {
		return "", err
	}
	ph := strings.Trim(member.Phone, "+")
	if member.TaxID == nil {
		return "", errors.New("Tax id for member cannot be null")
	}
	vReq := &clear.ConsumerVerificationRequest{
		ConsumerId:  string(member.ConsumerID),
		FirstName:   member.FirstName,
		LastName:    member.LastName,
		Email:       member.Email,
		DateOfBirth: member.DateOfBirth.String(),
		Address: &golang.AddressRequest{
			Line_1:     member.LegalAddress.StreetAddress,
			Locality:   member.LegalAddress.City,
			AdminArea:  member.LegalAddress.State,
			PostalCode: member.LegalAddress.PostalCode,
		},
		TaxId: *member.TaxID,
		Phone: ph,
	}

	res, err := c.RiskInformConsumerVerification(context.Background(), vReq)
	if err != nil {
		log.Println("Error occured while calling consumerVerification", err)
		return "", err
	}
	log.Println("Successfully completed clear consumerVerification")

	var resultData = make(map[string]interface{})
	err = json.Unmarshal([]byte(res.ResultData), &resultData)
	if err != nil {
		return string(res.ResultData), err
	}

	resultData["created"] = res.GetCreated()
	resultData["modified"] = res.GetModified()
	resBytes, err := json.Marshal(resultData)
	if err != nil {
		return string(res.ResultData), err
	}

	return string(resBytes), nil
}

//TODO finish implementation
func (s memberService) GetAlloyKYC(memberID shared.BusinessMemberID, businessID shared.BusinessID) (string, error) {

	sn, err := grpc.GetConnectionStringForService(grpc.ServiceNameVerification)
	if err != nil {
		return "", err
	}
	client, err := grpc.NewInsecureClient(sn)
	if err != nil {
		return "", err
	}
	defer client.CloseAndCancel()
	c := alloy.NewAlloyServiceClient(client.GetConn())

	member, err := s.GetByID(memberID, businessID)
	if err != nil {
		return "", err
	}
	gReq := &alloy.GetConsumerRequest{
		ConsumerId: string(member.ConsumerID),
	}

	res, err := c.GetConsumer(client.GetContext(), gReq)
	if err != nil {
		return "", err
	}

	return string(res.ResultData), nil
}

func (s memberService) GetClearKYC(memberID shared.BusinessMemberID, businessID shared.BusinessID, kycType goLibClear.ClearKycType) (string, error) {

	sn, err := grpc.GetConnectionStringForService(grpc.ServiceNameVerification)
	if err != nil {
		return "", err
	}
	client, err := grpc.NewInsecureClient(sn)
	if err != nil {
		return "", err
	}
	defer client.CloseAndCancel()
	c := clear.NewClearServiceClient(client.GetConn())

	member, err := s.GetByID(memberID, businessID)
	if err != nil {
		return "", err
	}
	gReq := &clear.GetConsumerRequest{
		ConsumerId:   string(member.ConsumerID),
		ClearKycType: string(kycType),
	}
	res, err := c.GetConsumer(client.GetContext(), gReq)
	if err != nil {
		return "", err
	}

	var resultData = make(map[string]interface{})
	err = json.Unmarshal([]byte(res.ResultData), &resultData)
	if err != nil {
		return string(res.ResultData), err
	}

	resultData["created"] = res.GetCreated()
	resultData["modified"] = res.GetModified()
	resBytes, err := json.Marshal(resultData)
	if err != nil {
		return string(res.ResultData), err
	}

	return string(resBytes), nil
}

func (s memberService) RunKYC(ssn string) (*KYCResult, error) {

	m, err := s.GetByTaxID(ssn)
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup

	result := KYCResult{}
	summary := KYCSummary{}

	// Alloy KYC
	wg.Add(1)
	go func() {
		defer wg.Done()
		summary.IdentityPartner = s.runAlloyKYC(m.ID, m.BusinessID)
	}()

	// BBVA KYC
	wg.Add(1)
	go func() {
		defer wg.Done()
		summary.BankPartner = s.runBankKYC(m.KYCStatus)
	}()

	// Phone verification
	wg.Add(1)
	go func() {
		defer wg.Done()
		summary.PhonePartner = s.runPhoneVerification(m.ID, m.BusinessID)
	}()

	// Intercom verification
	wg.Add(1)
	go func() {
		defer wg.Done()
		summary.LocationPartner = s.runLocationVerification(m)
	}()

	wg.Wait()

	result.KYCSummary = summary

	if (result.KYCSummary.IdentityPartner == VerificationStatusUnverified) || (result.KYCSummary.BankPartner == VerificationStatusUnverified) {
		result.Result = KYCStatusDeclined
		return &result, nil
	}

	identityVerified := result.KYCSummary.IdentityPartner == VerificationStatusVerified
	bankVerified := result.KYCSummary.BankPartner == VerificationStatusVerified
	locationVerified := result.KYCSummary.LocationPartner == VerificationStatusVerified
	phoneVerified := result.KYCSummary.PhonePartner == VerificationStatusVerified

	if identityVerified && bankVerified && locationVerified && phoneVerified {
		result.Result = KYCStatusApproved
	} else {
		result.Result = KYCStatusReview
	}

	return &result, nil
}

func (s memberService) GetByTaxID(ssn string) (*Member, error) {
	var m Member
	err := s.Get(&m, `
	SELECT
		business_member.id, business_member.consumer_id, business_member.business_id,
		business_member.title_type, business_member.title_other, business_member.ownership,
		business_member.is_controlling_manager, business_member.deactivated, business_member.created,
		business_member.modified, consumer.first_name, consumer.middle_name, consumer.last_name,
		consumer.email, consumer.phone, consumer.date_of_birth, consumer.tax_id AS tax_id_unmasked, consumer.tax_id_type,
		consumer.kyc_status, consumer.legal_address, consumer.mailing_address, consumer.work_address,
		consumer.residency, consumer.citizenship_countries, consumer.occupation, consumer.income_type,
		consumer.activity_type, consumer.is_restricted
	FROM
		business_member
	LEFT JOIN
		consumer ON business_member.consumer_id = consumer.id
	WHERE
		tax_id = $1`,
		ssn)

	if err != nil && err == sql.ErrNoRows {
		return nil, services.ErrorNotFound{}.New("")
	}

	if err != nil {
		return nil, err
	}

	return &m, nil
}

func (s memberService) runAlloyKYC(mID shared.BusinessMemberID, bID shared.BusinessID) VerificationStatus {
	var a AlloyKYCResult
	alloyResp, err := s.RunAlloyKYC(mID, bID)
	if err != nil {
		return VerificationStatusUnknown
	}

	err = json.Unmarshal([]byte(alloyResp), &a)
	if err != nil {
		return VerificationStatusUnknown
	}

	// Alloy results
	switch a.AlloySummary.Outcome {
	case AlloyOutComeApproved:
		return VerificationStatusVerified
	case AlloyOutComeManualReview:
		return VerificationStatusInReview
	case AlloyOutComeDeclined:
		return VerificationStatusUnverified
	default:
		return VerificationStatusUnknown
	}
}

func (s memberService) runBankKYC(status services.KYCStatus) VerificationStatus {
	// BBVA results
	switch status {
	case services.KYCStatusApproved:
		return VerificationStatusVerified
	case services.KYCStatusDeclined:
		return VerificationStatusUnverified
	default:
		return VerificationStatusInReview
	}
}

func (s memberService) runPhoneVerification(mID shared.BusinessMemberID, bID shared.BusinessID) VerificationStatus {
	var e EveryoneResult
	phoneResp, err := s.PhoneVerification(mID, bID)
	if err != nil {
		log.Println(err)
		return VerificationStatusUnknown
	}

	err = json.Unmarshal([]byte(phoneResp), &e)
	if err != nil {
		log.Println(err)
		return VerificationStatusUnknown
	}

	// Everyone results
	switch e.Type {
	case PhoneTypeBusiness, PhoneTypePerson:
		return VerificationStatusVerified
	default:
		return VerificationStatusUnverified
	}
}

func (s memberService) runLocationVerification(m *Member) VerificationStatus {
	var i IntercomResponse
	intercomResp, err := intercom.New(s.sourceReq).GetByEmailID(m.Email)
	if err != nil {
		return VerificationStatusUnknown
	}

	err = json.Unmarshal([]byte(*intercomResp), &i)
	if err != nil {
		return VerificationStatusUnknown
	}

	// Intercom results
	if len(i.Users) > 0 {
		region := i.Users[0].LocationData.RegionName
		state, ok := StateMap[region]
		if !ok {
			return VerificationStatusUnverified
		}

		var stateExpression string
		if m.MailingAddress != nil {
			stateExpression = m.MailingAddress.State
		}

		if m.LegalAddress != nil {
			if len(stateExpression) > 0 {
				stateExpression = stateExpression + "|"
			}

			stateExpression = stateExpression + m.LegalAddress.State
		}

		if m.WorkAddress != nil {
			if len(stateExpression) > 0 {
				stateExpression = stateExpression + "|"
			}

			stateExpression = stateExpression + m.WorkAddress.State
		}

		if len(stateExpression) == 0 {
			return VerificationStatusUnverified
		}

		regEx := regexp.MustCompile(stateExpression)

		if regEx.MatchString(state) {
			return VerificationStatusVerified
		} else {
			return VerificationStatusUnverified
		}
	} else {
		return VerificationStatusUnverified
	}
}

// For tools/Projectlock
func (s memberService) GetConsumerIDsForBusinessIDsInternal(businessIDs []string) ([]string, error) {
	var consumerIds = []string{}
	businessIdsStr := "'" + strings.Join(businessIDs[:], "','") + "'"
	q := fmt.Sprintf("SELECT consumer_id FROM business_member WHERE business_id IN (%v);", businessIdsStr)
	err := s.Select(&consumerIds, q)
	if err != nil {
		fmt.Println("Get consumer id err: ", err)
	}
	return consumerIds, err
}
