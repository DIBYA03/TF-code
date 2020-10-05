/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

package review

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	partnerbank "github.com/wiseco/core-platform/partner/bank"
	"github.com/wiseco/core-platform/services"
	"github.com/wiseco/core-platform/services/csp"
	"github.com/wiseco/core-platform/services/csp/business"
	"github.com/wiseco/core-platform/shared"
)

// Review Review services
type Review interface {
	UpdateReview(id shared.BusinessID, status csp.Status) (*Response, error)

	Approve(shared.BusinessID) (*Response, error)

	Decline(shared.BusinessID) (*Response, error)
}

type reviewService struct {
	sourceReq services.SourceRequest
}

// NewReviewService ..
func NewReviewService(req services.SourceRequest) Review {
	return reviewService{sourceReq: req}
}

// NewReviewWithout give you a review service without a source request
func NewReviewWithout() Review {
	return reviewService{}
}

func (service reviewService) Approve(id shared.BusinessID) (*Response, error) {
	business, err := business.NewCSPService().ByBusinessID(id)
	var currentStatus csp.Status
	if err != nil {
		return nil, err
	}

	if !business.Status.Valid() {
		return nil, errors.New("Business status is invalid")
	}

	if business.Status == csp.StatusTrainingComplete {
		return nil, errors.New("Business review process as been completed")
	}

	switch business.Status {
	case csp.StatusMemberReview:
		currentStatus = csp.StatusDocReview
	case csp.StatusDocReview:
		currentStatus = csp.StatusRiskReview
	case csp.StatusRiskReview:
		return service.startVerification(id)
	case csp.StatusTraining:
		currentStatus = csp.StatusTrainingComplete
	default:
		return nil, errors.New("Looks like this business cannot be approve, please check its status")
	}

	if currentStatus == csp.StatusDocReview {
		if err := service.sendCopyBusinessDocument(id, business); err != nil {
			return nil, err
		}
	}

	return service.updateStatus(id, &currentStatus, nil)
}

func (service reviewService) sendCopyBusinessDocument(id shared.BusinessID, cspBusiness business.CSPBusiness) error {
	if cspBusiness.EntityType == nil {
		// entitytType is null, cant determine if we can copy the document
		return nil
	}

	if *cspBusiness.EntityType == EntityTypeSoleProprietor {
		return csp.SendDocumentMessage(csp.Message{
			EntityID: string(id),
			Category: csp.CategoryBusinessDocument,
			Action:   csp.ActionCopy,
		})
	}

	return nil
}

func (service reviewService) Decline(id shared.BusinessID) (*Response, error) {
	business, err := business.NewCSPService().ByBusinessID(id)
	var currentStatus csp.Status
	if err != nil {
		return nil, err
	}

	if !business.Status.Valid() {
		return nil, errors.New("Business status is invalid")
	}

	if business.Status == csp.StatusTrainingComplete {
		return nil, errors.New("Business review process as been completed")
	}

	switch business.Status {
	case csp.StatusRiskReview:
		currentStatus = csp.StatusDocReview
	case csp.StatusDocReview:
		currentStatus = csp.StatusMemberReview
	case csp.StatusBankReview:
		currentStatus = csp.StatusRiskReview
	case csp.StatusMemberReview:
		return service.decline(id, false)
	case csp.StatusBankApproved:
		return service.decline(id, false)
	default:
		return nil, errors.New("Looks like this business cannot be decline, please check its status")
	}

	return service.updateStatus(id, &currentStatus, nil)
}

// UpdateReview  will update the review process and start verification on Bank side
// when the right status is sent
func (service reviewService) UpdateReview(id shared.BusinessID, status csp.Status) (*Response, error) {
	business, err := business.NewCSPService().ByBusinessID(id)
	var currentStatus csp.Status
	if err != nil {
		return nil, err
	}

	if !business.Status.Valid() {
		return nil, errors.New("Business status is invalid")
	}

	switch status {
	case csp.StatusTraining:
		currentStatus = csp.StatusTraining
	case csp.StatusContinue:
		return service.continueReview(id)
	default:
		return nil, errors.New("Looks like this business can not be approve, please check its status")
	}
	return service.updateStatus(id, &currentStatus, nil)

}

// will update cps business status and core business kyc status
func (service reviewService) updateStatus(id shared.BusinessID, status *csp.Status, kcyStatus *csp.KYCStatus) (*Response, error) {
	var sts string

	// Update core business KYC Status
	if kcyStatus != nil {
		if err := business.New(services.SourceRequest{}).UpdateKYC(id, *kcyStatus); err != nil {
			log.Printf("Error updating business KYC %v", err)
			return nil, err
		}
		sts = kcyStatus.String()
	}
	if status != nil {
		updates := business.CSPBusinessUpdate{
			Status: status,
		}
		// Update csp business status
		if _, err := business.NewCSPService().CSPBusinessUpdateByBusinessID(id, updates); err != nil {
			log.Printf("Error updating review status %s", err)
			return nil, err
		}
		sts = status.String()
	}

	return &Response{
		Status: sts,
	}, nil
}

func (service reviewService) startVerification(id shared.BusinessID) (*Response, error) {
	bank, err := partnerbank.GetProxyBank(partnerbank.ProviderNameBBVA)
	if err != nil {
		log.Printf("Error getting partner bank %v", err)
	}

	bankID, err := bank.ProxyService(service.sourceReq.PartnerBankRequest()).GetBusinessBankID(partnerbank.BusinessID(id))
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	// If business was previously submitted, only update the review status
	if bankID != nil {
		currentStatus := csp.StatusBankReview
		return service.updateStatus(id, &currentStatus, nil)
	}

	resp, err := New(service.sourceReq).startBankReview(id)
	if err != nil {
		log.Printf("error starting bank review %v", err)
		return nil, err
	}
	var status csp.Status
	switch csp.KYCStatus(resp.Status) {
	case csp.KYCStatusApproved:
		status = csp.StatusBankApproved
	case csp.KYCStatusDeclined:
		return service.decline(id, true)
	case csp.KYCStatusReview:
		status = csp.StatusBankReview
	default:
		return nil, fmt.Errorf("Invalid business status %s", resp.Status)
	}
	submitted := time.Now()
	updates := business.CSPBusinessUpdate{
		IDVs:      resp.ReviewItems,
		Status:    &status,
		Submitted: &submitted,
		Notes:     resp.Notes,
	}

	// Update csp business
	if _, err := business.NewCSPService().CSPBusinessUpdateByBusinessID(id, updates); err != nil {
		log.Printf("Error updating csp business %s", err)
		return nil, err
	}

	//Update business KYC Status
	sts := csp.KYCStatus(resp.Status)
	if _, err := service.updateStatus(id, nil, &sts); err != nil {
		log.Printf("Error updaing business KYC %v", err)
		return nil, err
	}

	csp.SendDocumentMessage(csp.Message{
		EntityID: string(id),
		Action:   csp.ActionUpload,
		Category: csp.CategoryBusiness,
	})
	return resp, nil
}

func (service reviewService) continueReview(id shared.BusinessID) (*Response, error) {
	resp, err := New(service.sourceReq).Continue(id)
	if err != nil {
		return nil, err
	}
	var status csp.Status

	switch csp.KYCStatus(resp.Status) {
	case csp.KYCStatusApproved:
		status = csp.StatusApproved

	case csp.KYCStatusDeclined:
		return service.decline(id, true)
	case csp.KYCStatusReview:
		status = csp.StatusBankReview
	}

	updates := business.CSPBusinessUpdate{
		IDVs:   resp.ReviewItems,
		Notes:  resp.Notes,
		Status: &status,
	}

	// Update csp business
	if _, err := business.NewCSPService().CSPBusinessUpdateByBusinessID(id, updates); err != nil {
		log.Printf("Error updating csp business %s", err)
		return nil, err
	}

	//Update business KYC Status
	if err := business.New(services.SourceRequest{}).UpdateKYC(id, csp.KYCStatus(resp.Status)); err != nil {
		log.Printf("Error updating core business KYC %v", err)
		return nil, err
	}

	return resp, nil
}

func (service reviewService) decline(id shared.BusinessID, bankDeclined bool) (*Response, error) {
	status := csp.StatusDeclined
	if bankDeclined {
		status = csp.StatusBankDeclined
	}

	sts := csp.KYCStatusDeclined
	if _, err := service.updateStatus(id, &status, &sts); err != nil {
		return nil, err
	}

	return &Response{Status: string(status)}, nil
}

func businessIDVS(businessID shared.BusinessID) []partnerbank.IDVerify {
	b, err := business.NewCSPService().ByBusinessID(businessID)
	if err != nil {
		return []partnerbank.IDVerify{}
	}
	var idvs []partnerbank.IDVerify

	if b.IDVs != nil {
		for _, idv := range *b.IDVs {
			idvs = append(idvs, reviewToIDV(idv))
		}
	}
	return idvs
}
