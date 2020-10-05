/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

package review

import (
	"errors"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	partnerbank "github.com/wiseco/core-platform/partner/bank"
	"github.com/wiseco/core-platform/services"
	biz "github.com/wiseco/core-platform/services/business"
	"github.com/wiseco/core-platform/services/csp"
	docs "github.com/wiseco/core-platform/services/csp/document"
	"github.com/wiseco/core-platform/services/data"
	"github.com/wiseco/core-platform/shared"
)

//BankVerification  this service talks to the bank layer directly
type BankVerification interface {
	startBankReview(shared.BusinessID) (*Response, error)
	Continue(shared.BusinessID) (*Response, error)
	GetStatus(shared.BusinessID) (*Response, error)
}

type verification struct {
	*sqlx.DB
	sourceReq services.SourceRequest
}

// New BankVerification service
func New(source services.SourceRequest) BankVerification {
	return verification{data.DBWrite, source}
}

// NewWithout return a new BankVerification without a source request
func NewWithout() BankVerification {
	return verification{data.DBWrite, services.SourceRequest{}}
}

func getOwner(id string) string {
	var user = struct {
		FirstName string `db:"first_name"`
		LastName  string `db:"last_name"`
	}{}
	err := data.DBRead.Get(&user, fmt.Sprintf(`SELECT consumer.first_name,consumer.last_name
	FROM wise_user
	JOIN consumer ON wise_user.consumer_id = consumer.id
	WHERE wise_user.id = '%s'
	`, id))
	if err != nil {
		log.Printf("Error getting user with id:%s  details:%v", id, err)
	}
	return user.FirstName + " " + user.LastName
}

// StartBankReview start the verification process on bank side
func (service verification) startBankReview(id shared.BusinessID) (*Response, error) {
	var b biz.Business
	if err := service.Get(&b, "SELECT * FROM business WHERE id = $1", id); err != nil {
		log.Printf("error getting business %v", err)
		return nil, err
	}

	if errList := checkCommon(b); len(errList) > 0 {
		log.Printf("error missing params %v", errList)
		return nil, ErrorResponse{
			Raw:       KYCErrorTypeParam,
			ErrorType: KYCErrorTypeParam,
			Values:    errListString(errList),
			Business:  &b,
		}
	}

	if err := service.checkStatus(&b); err != nil {
		log.Printf("error checking status %v", err)
		return nil, err
	}
	resp, err := service.sendRequest(&b)
	if err != nil {
		log.Printf("error sending request KYC Request %v", err)
		return nil, err
	}

	var items services.StringArray
	for _, item := range resp.KYC.IDVerifyRequired {
		items = append(items, string(item))
	}
	var notes csp.KYCNotes
	for _, n := range resp.KYC.Notes {
		notes = append(notes, csp.KYCNote(n))
	}
	kycnotes := notes.Raw()
	return &Response{
		Status:        resp.KYC.Status.String(),
		ReviewItems:   &items,
		Notes:         &kycnotes,
		BusinessName:  b.Name(),
		BusinessOwner: b.OwnerID,
		EntityType:    b.EntityType,
	}, nil
}

//checkStatus checks current business status
func (service verification) checkStatus(b *biz.Business) error {
	switch b.KYCStatus {
	case services.KYCStatusApproved:
		return ErrorResponse{
			Raw:       KYCErrorTypeApproved,
			ErrorType: KYCErrorTypeApproved,
			Business:  b,
		}
	case services.KYCStatusDeclined:
		return ErrorResponse{
			Raw:       KYCErrorTypeDeclined,
			ErrorType: KYCErrorTypeDeclined,
			Business:  b,
		}
	}
	return nil
}

func (service verification) sendRequest(b *biz.Business) (*partnerbank.IdentityStatusBusinessResponse, error) {
	r, err := service.buildRequest(b)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	bank, err := partnerbank.GetBusinessBank(partnerbank.ProviderNameBBVA)
	if err != nil {
		return nil, err
	}

	srv := bank.BusinessEntityService(service.sourceReq.PartnerBankRequest())
	if err != nil {
		return nil, err
	}

	resp, err := srv.Create(r)
	if err != nil {
		log.Println(err)
		return nil, ErrorResponse{
			Raw:       err,
			Values:    []string{err.Error()},
			ErrorType: KYCErrorTypeOther,
			Business:  b,
		}
	}

	return resp, err
}

func (service verification) buildRequest(b *biz.Business) (partnerbank.CreateBusinessRequest, error) {
	// Fetch business members
	memberRequests, err := service.getPartnerBusinessMembers(b.ID)
	if err != nil {
		log.Printf("error getting  members %v", err)
		return partnerbank.CreateBusinessRequest{}, ErrorResponse{
			Raw:       err,
			ErrorType: KYCErrorTypeMember,
			Values:    []string{err.Error()},
			Business:  b,
		}
	}

	if b.FormationDocumentID == nil {
		return partnerbank.CreateBusinessRequest{}, ErrorResponse{
			Raw:       errors.New("Business does not have a formation document set"),
			ErrorType: KYCParamErrorTypeFormationDoc,
			Business:  b,
		}
	}

	doc, err := docs.NewDocumentService().GetByID(b.ID, *b.FormationDocumentID)
	if err != nil {
		log.Printf("error getting docs %v", err)
		return partnerbank.CreateBusinessRequest{}, ErrorResponse{
			Raw:       fmt.Errorf("Error getting document for entity formation - business ID: %s documentID: %v error: %v", b.ID, b.FormationDocumentID, err),
			ErrorType: KYCParamErrorTypeFormationDoc,
			Business:  b,
		}
	}

	if doc.Number == nil {
		return partnerbank.CreateBusinessRequest{}, ErrorResponse{
			Raw: fmt.Errorf(
				"document number required for entity formation - business ID: %s documentID: %v", b.ID, b.FormationDocumentID,
			),
			ErrorType: KYCParamErrorTypeFormationDoc,
			Business:  b,
		}
	}

	if doc.DocType == nil {
		return partnerbank.CreateBusinessRequest{}, ErrorResponse{
			Raw: fmt.Errorf(
				"document type required for entity formation - business ID: %s documentID: %v", b.ID, b.FormationDocumentID,
			),
			ErrorType: KYCParamErrorTypeFormationDoc,
			Business:  b,
		}
	}

	var ef partnerbank.EntityFormationRequest

	ef = partnerbank.EntityFormationRequest{
		DocumentType:   partnerbank.BusinessIdentityDocument(*doc.DocType),
		Number:         *doc.Number,
		IssueDate:      doc.IssuedDate.Time(),
		ExpirationDate: doc.ExpirationDate.Time(),
	}

	r := partnerbank.CreateBusinessRequest{
		BusinessID:         partnerbank.BusinessID(b.ID),
		LegalName:          shared.StringValue(b.LegalName),
		DBA:                b.DBA,
		TaxIDType:          partnerbank.BusinessTaxIDType(*b.TaxIDType),
		TaxID:              b.TaxID.String(),
		EntityType:         partnerbank.BusinessEntity(*b.EntityType),
		IndustryType:       partnerbank.BusinessIndustry(*b.IndustryType),
		Phone:              *b.Phone,
		Email:              *b.Email,
		Purpose:            *b.Purpose,
		ExpectedActivities: b.ActivityType.ToPartnerBankActivity(),
		Members:            memberRequests,
		LegalAddress:       b.LegalAddress.ToPartnerBankAddress(services.AddressTypeLegal),
		OriginCountry:      partnerbank.Country(*b.OriginCountry),
		OriginState:        *b.OriginState,
		OriginDate:         b.OriginDate.Time(),
		OperationType:      partnerbank.BusinessOperationType(*b.OperationType),
		EntityFormation:    &ef,
	}
	// Use legal address for headquarters if not available
	if b.HeadquarterAddress != nil {
		address := b.HeadquarterAddress.ToPartnerBankAddress(services.AddressTypeHeadquarter)
		r.HeadquarterAddress = address
	} else {
		r.HeadquarterAddress = r.LegalAddress
		r.HeadquarterAddress.Type = partnerbank.AddressRequestType(services.AddressTypeHeadquarter)
	}
	if b.MailingAddress != nil {
		address := b.MailingAddress.ToPartnerBankAddress(services.AddressTypeMailing)
		r.MailingAddress = &address
	}

	return r, err
}
