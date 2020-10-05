package bbva

import (
	"errors"
	"strings"

	"github.com/wiseco/core-platform/partner/bank"
	partnerbank "github.com/wiseco/core-platform/partner/bank"
)

type IDVerify string

const (
	IDVerifyAddress      = IDVerify("ADDRESS")
	IDVerifyName         = IDVerify("NAME")
	IDVerifyDOB          = IDVerify("DOB")
	IDVerifyMismatch     = IDVerify("MISMATCH")
	IDVerifyIDV          = IDVerify("IDV")
	IDVerifySSN          = IDVerify("SSN")
	IDVerifyOFAC         = IDVerify("OFAC")
	IDVerifyPrimaryDoc   = IDVerify("PRIMARY_DOC")
	IDVerifySecondaryDoc = IDVerify("SECONDARY_DOC")
	IDVerifyDoc          = IDVerify("DOC")
)

var kycIDVerifyPartnerTo = map[IDVerify]partnerbank.IDVerify{
	IDVerifyAddress:      partnerbank.IDVerifyAddress,
	IDVerifyName:         partnerbank.IDVerifyFullName,
	IDVerifyDOB:          partnerbank.IDVerifyDateOfBirth,
	IDVerifyMismatch:     partnerbank.IDVerifyMismatch,
	IDVerifyIDV:          partnerbank.IDVerifyOther,
	IDVerifySSN:          partnerbank.IDVerifyTaxId,
	IDVerifyOFAC:         partnerbank.IDVerifyOFAC,
	IDVerifyPrimaryDoc:   partnerbank.IDVerifyPrimaryDoc,
	IDVerifySecondaryDoc: partnerbank.IDVerifySecondaryDoc,
	IDVerifyDoc:          partnerbank.IDVerifyFormationDoc,
}

func idVerifyPartnerTo(idv []IDVerify) ([]partnerbank.IDVerify, error) {
	var pidv []partnerbank.IDVerify
	for _, id := range idv {
		idvUpper := IDVerify(strings.ToUpper(string(id)))
		val, ok := kycIDVerifyPartnerTo[idvUpper]
		if !ok {
			return pidv, errors.New("invalid idv value")
		}

		pidv = append(pidv, val)
	}

	return pidv, nil
}

var kycIDVerifyPartnerFrom = map[partnerbank.IDVerify]IDVerify{
	partnerbank.IDVerifyAddress:      IDVerifyAddress,
	partnerbank.IDVerifyFullName:     IDVerifyName,
	partnerbank.IDVerifyDateOfBirth:  IDVerifyDOB,
	partnerbank.IDVerifyMismatch:     IDVerifyMismatch,
	partnerbank.IDVerifyOther:        IDVerifyIDV,
	partnerbank.IDVerifyTaxId:        IDVerifySSN,
	partnerbank.IDVerifyOFAC:         IDVerifyOFAC,
	partnerbank.IDVerifyPrimaryDoc:   IDVerifyPrimaryDoc,
	partnerbank.IDVerifySecondaryDoc: IDVerifySecondaryDoc,
	partnerbank.IDVerifyFormationDoc: IDVerifyDoc,
}

func idVerifyPartnerFrom(pidv []partnerbank.IDVerify) ([]IDVerify, error) {
	var idv []IDVerify
	for _, id := range pidv {
		val, ok := kycIDVerifyPartnerFrom[id]
		if !ok {
			return idv, errors.New("invalid idv value")
		}

		idv = append(idv, val)
	}

	return idv, nil
}

type DocumentIDVerify string

const (
	DocumentIDVerifyAddress      = DocumentIDVerify("address")
	DocumentIDVerifyName         = DocumentIDVerify("name")
	DocumentIDVerifyDOB          = DocumentIDVerify("dob")
	DocumentIDVerifyMismatch     = DocumentIDVerify("mismatch")
	DocumentIDVerifyIDV          = DocumentIDVerify("idv")
	DocumentIDVerifySSN          = DocumentIDVerify("ssn")
	DocumentIDVerifyOFAC         = DocumentIDVerify("ofac")
	DocumentIDVerifyPrimaryDoc   = DocumentIDVerify("primary_doc")
	DocumentIDVerifySecondaryDoc = DocumentIDVerify("secondary_doc")
	DocumentIDVerifyDoc          = DocumentIDVerify("doc")
)

var kycDocumentIDVerifyPartnerFrom = map[partnerbank.IDVerify]DocumentIDVerify{
	partnerbank.IDVerifyAddress:      DocumentIDVerifyAddress,
	partnerbank.IDVerifyFullName:     DocumentIDVerifyName,
	partnerbank.IDVerifyDateOfBirth:  DocumentIDVerifyDOB,
	partnerbank.IDVerifyMismatch:     DocumentIDVerifyMismatch,
	partnerbank.IDVerifyOther:        DocumentIDVerifyIDV,
	partnerbank.IDVerifyTaxId:        DocumentIDVerifySSN,
	partnerbank.IDVerifyOFAC:         DocumentIDVerifyOFAC,
	partnerbank.IDVerifyFormationDoc: DocumentIDVerifyDoc,
	partnerbank.IDVerifyPrimaryDoc:   DocumentIDVerifyPrimaryDoc,
	partnerbank.IDVerifySecondaryDoc: DocumentIDVerifySecondaryDoc,
}

func documentIDVerifyPartnerFrom(pidv []partnerbank.IDVerify) ([]DocumentIDVerify, error) {
	var idv []DocumentIDVerify
	for _, id := range pidv {
		val, ok := kycDocumentIDVerifyPartnerFrom[id]
		if !ok {
			return idv, errors.New("invalid idv value")
		}

		idv = append(idv, val)
	}

	return idv, nil
}

// KYCResponse contains a Know Your Customer (KYC) status that your application must process.
type KYCResponse struct {
	Status KYCStatus `json:"status"`

	// Customer Identification Program status.
	// PASS: Customer information is sufficient to enable the creation of the new consumer record.
	// IDV: IF KYC STATUS equals REVIEW, IDV lists what additional identification must be submitted for review. Possible values are address, ssn, name, dob, idv, ofac or mismatch.
	// 		IF IDV equals mismatch, BBVA’s bank partner has an existing consumer record that contains the submitted ssn, but that one or more of the submitted values for name, dob or
	//		citizenship_status do not match the information contained in that record. Review the submitted values of name, dob or citizenship for errors. If there are no errors, contact
	//		our client integration team for assistance updating the existing record.
	// FAIL: The submitted customer information corresponds to information in the records of the US Treasury’s Office of Foreign Assets Control (OFAC).
	CIP string `json:"cip,omitempty"`

	// Overall KYC risk of LOW or MEDIUM, based on the combination of identity information provided in the request body.
	Risk string `json:"risk,omitempty"`

	IDVerifyRequired []IDVerify `json:"idv_required,omitempty"`
}

type KYCStatus string

const (
	KYCStatusApproved KYCStatus = "APPROVED" // The new customer record can be created using the information supplied. This is the default status for Open Platform sandbox requests.
	KYCStatusReview   KYCStatus = "REVIEW"   // Customer must submit additional identification before the consumer record can be created. In the row for cip, see IDV.
	KYCStatusDeclined KYCStatus = "DECLINED" // The submitted identity information cannot be used to create a consumer record.
)

var kycStatusToPartnerMap = map[KYCStatus]partnerbank.KYCStatus{
	KYCStatusApproved: partnerbank.KYCStatusApproved,
	KYCStatusReview:   partnerbank.KYCStatusReview,
	KYCStatusDeclined: partnerbank.KYCStatusDeclined,
}

var kycStatusFromPartnerMap = map[partnerbank.KYCStatus]KYCStatus{
	partnerbank.KYCStatusApproved: KYCStatusApproved,
	partnerbank.KYCStatusReview:   KYCStatusReview,
	partnerbank.KYCStatusDeclined: KYCStatusDeclined,
}

type KYCNoteResponse struct {
	Code   string `json:"code"`
	Detail string `json:"detail"`
}

func (r KYCResponse) toPartnerBankKYCResponse(respNotes []KYCNoteResponse) (*partnerbank.KYCResponse, error) {
	idvs, err := idVerifyPartnerTo(r.IDVerifyRequired)
	if err != nil {
		return nil, err
	}

	status, ok := kycStatusToPartnerMap[KYCStatus(strings.ToUpper(string(r.Status)))]
	if !ok {
		return nil, bank.NewErrorFromCode(bank.ErrorCodeInvalidKYCStatus)
	}

	var notes []partnerbank.KYCNote
	if respNotes != nil {
		for _, n := range respNotes {
			notes = append(notes, partnerbank.KYCNote{Code: n.Code, Desc: n.Detail})
		}
	}

	return &partnerbank.KYCResponse{
		Status:           status,
		IDVerifyRequired: idvs,
		Notes:            notes,
	}, nil
}

type KYCStatusResponse struct {
	KYC              KYCResponse       `json:"kyc"`
	KYCNotes         []KYCNoteResponse `json:"kyc_notes,omitempty"`
	DigitalFootprint interface{}       `json:"digital_footprint,omitempty"`
}
