/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

package review

import (
	"errors"
	"log"

	partnerbank "github.com/wiseco/core-platform/partner/bank"
	"github.com/wiseco/core-platform/services"
	biz "github.com/wiseco/core-platform/services/business"
	"github.com/wiseco/core-platform/services/csp"
	"github.com/wiseco/core-platform/services/csp/business"
	"github.com/wiseco/core-platform/shared"
)

// GetStatus makes a call to bank to get status
func (service verification) GetStatus(id shared.BusinessID) (*Response, error) {
	var b biz.Business
	err := service.Get(&b, "SELECT * FROM business WHERE id = $1", id)
	if err != nil {
		log.Printf("error getting business %v", err)
		return nil, ErrorResponse{
			Raw:       err,
			ErrorType: KYCErrorTypeOther,
		}
	}

	if errList := checkCommon(b); len(errList) > 0 {
		return nil, ErrorResponse{
			Raw:       errors.New("One or more required fields are either missing or incorrect"),
			Values:    errListString(errList),
			ErrorType: KYCErrorTypeParam,
		}
	}

	bank, err := partnerbank.GetBusinessBank(partnerbank.ProviderNameBBVA)
	if err != nil {
		return nil, err
	}

	srv := bank.BusinessEntityService(service.sourceReq.PartnerBankRequest())
	if err != nil {
		log.Printf("error creating service %v", err)
		return nil, err
	}

	resp, err := srv.Status(partnerbank.BusinessID(b.ID))
	if err != nil {
		log.Printf("error getting business status from bank %v", err)
		return nil, ErrorResponse{
			Raw:       err,
			ErrorType: KYCErrorTypeOther,
			Business:  &b,
		}
	}

	var items services.StringArray
	for _, item := range resp.KYC.IDVerifyRequired {
		items = append(items, string(item))
	}
	var itemsToUpdate *services.StringArray
	if len(items) > 0 {
		itemsToUpdate = &items
	}

	kyc, st := kycStatusFromResponse(resp.KYC.Status)
	updates := business.CSPBusinessUpdate{
		Status: &st,
		IDVs:   itemsToUpdate,
	}

	// update core business kyc status
	business.New(services.SourceRequest{}).UpdateKYC(id, kyc)

	// update csp business
	business.NewCSPService().CSPBusinessUpdateByBusinessID(id, updates)

	return &Response{
		Status:      resp.KYC.Status.String(),
		ReviewItems: &items,
	}, nil
}

// Convert bank kyc status to csp status
func kycStatusFromResponse(st partnerbank.KYCStatus) (csp.KYCStatus, csp.Status) {
	switch st {
	case partnerbank.KYCStatusApproved:
		return csp.KYCStatusApproved, csp.StatusBankApproved
	case partnerbank.KYCStatusDeclined:
		return csp.KYCStatusDeclined, csp.StatusBankDeclined
	case partnerbank.KYCStatusReview:
		return csp.KYCStatusReview, csp.StatusBankReview
	}
	return csp.KYCStatusReview, csp.StatusBankReview
}
