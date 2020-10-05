/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

package review

import (
	"fmt"
	"log"

	partnerbank "github.com/wiseco/core-platform/partner/bank"
	"github.com/wiseco/core-platform/services"
	biz "github.com/wiseco/core-platform/services/business"
	docsrv "github.com/wiseco/core-platform/services/csp/document"
	"github.com/wiseco/core-platform/shared"
)

//Continue verification proccess
func (service verification) Continue(id shared.BusinessID) (*Response, error) {
	var b biz.Business
	err := service.Get(&b, "SELECT * FROM business WHERE id = $1", id)
	if err != nil {
		return nil, ErrorResponse{
			Raw:       err,
			ErrorType: KYCErrorTypeOther,
		}
	}

	if kycErrors := checkCommon(b); len(kycErrors) > 0 {
		return nil, ErrorResponse{
			ErrorType: KYCErrorTypeParam,
			Values:    errListString(kycErrors),
			Business:  &b,
		}
	}

	var ef *partnerbank.EntityFormationRequest

	doc, err := docsrv.NewDocumentService().GetByID(b.ID, *b.FormationDocumentID)
	if err != nil {
		log.Printf("error getting document %v", err)
		return nil, ErrorResponse{
			Raw:       fmt.Errorf("Error getting document for formation: businesID: %s docuementID: %v error: %v", b.ID, b.FormationDocumentID, err),
			ErrorType: KYCParamErrorTypeFormationDoc,
			Business:  &b,
		}
	}

	if doc.DocType == nil {
		return nil, ErrorResponse{
			Raw: fmt.Errorf(
				"Document type is required for formation: businesID: %s docuementID: %v",
				b.ID, b.FormationDocumentID,
			),
			ErrorType: KYCParamErrorTypeFormationDoc,
			Business:  &b,
		}
	}

	if doc.Number == nil {
		return nil, ErrorResponse{
			Raw: fmt.Errorf(
				"Document number is required for formation: businesID: %s docuementID: %v",
				b.ID, b.FormationDocumentID,
			),
			ErrorType: KYCParamErrorTypeFormationDoc,
			Business:  &b,
		}
	}

	ef = &partnerbank.EntityFormationRequest{
		DocumentType:   partnerbank.BusinessIdentityDocument(*doc.DocType),
		Number:         *doc.Number,
		IssueDate:      doc.IssuedDate.Time(),
		ExpirationDate: doc.ExpirationDate.Time(),
	}

	ureq := partnerbank.UpdateBusinessRequest{
		BusinessID:      partnerbank.BusinessID(b.ID),
		TaxIDType:       partnerbank.BusinessTaxIDType(*b.TaxIDType),
		TaxID:           b.TaxID.String(),
		EntityType:      partnerbank.BusinessEntity(*b.EntityType),
		IndustryType:    partnerbank.BusinessIndustry(*b.IndustryType),
		Purpose:         *b.Purpose,
		OriginCountry:   partnerbank.Country(*b.OriginCountry),
		OriginState:     *b.OriginState,
		OriginDate:      b.OriginDate.Time(),
		OperationType:   partnerbank.BusinessOperationType(*b.OperationType),
		EntityFormation: ef,
	}

	bank, err := partnerbank.GetBusinessBank(partnerbank.ProviderNameBBVA)
	if err != nil {
		return nil, err
	}

	srv := bank.BusinessEntityService(service.sourceReq.PartnerBankRequest())
	if err != nil {
		return nil, err
	}

	resp, err := srv.Update(ureq)
	if err != nil {
		return nil, ErrorResponse{
			Raw:       err,
			ErrorType: KYCErrorTypeOther,
			Business:  &b,
			Values:    []string{err.Error()},
		}
	}

	var items services.StringArray
	for _, item := range resp.KYC.IDVerifyRequired {
		items = append(items, string(item))
	}

	return &Response{
		Status:      resp.KYC.Status.String(),
		ReviewItems: &items,
	}, nil
}
