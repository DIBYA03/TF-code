/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

package review

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/jmoiron/sqlx/types"
	partnerbank "github.com/wiseco/core-platform/partner/bank"
	"github.com/wiseco/core-platform/services"
	"github.com/wiseco/core-platform/services/csp"
	cspdoc "github.com/wiseco/core-platform/services/csp/document"
	coreDB "github.com/wiseco/core-platform/services/data"
	docs "github.com/wiseco/core-platform/services/document"
	"github.com/wiseco/core-platform/shared"
)

// DocumentUploader BBVA document upload service
type DocumentUploader interface {
	BusinessSingleUpload(shared.BusinessID, shared.BusinessDocumentID) error
	BusinessSingleReUpload(shared.BusinessID) error
	ConsumerSingleUpload(shared.ConsumerID, shared.ConsumerDocumentID) error
	ConsumerSingleReUpload(shared.ConsumerID) error

	BusinessSingle(shared.BusinessID, shared.BusinessDocumentID) error
	ConsumerSingle(shared.ConsumerID, shared.ConsumerDocumentID) error

	BusinessMultiple(shared.BusinessID) error
	ConsumerMultiple(shared.ConsumerID) error

	ConsumerReUpload(shared.ConsumerID) error
	BusinessReUpload(shared.BusinessID) error
}

type uploader struct {
	sourceReq services.SourceRequest
}

type bizdoc struct {
	ID          shared.BusinessDocumentID `db:"id"`
	Key         string                    `db:"storage_key"`
	ContentType string                    `db:"content_type"`
	DocType     string                    `db:"doc_type"`
}

type condoc struct {
	ID          shared.ConsumerDocumentID `db:"id"`
	Key         string                    `db:"storage_key"`
	ContentType string                    `db:"content_type"`
	DocType     string                    `db:"doc_type"`
}

// NewUploader ..
func NewUploader(source services.SourceRequest) DocumentUploader {
	return uploader{source}
}

// Handles api calls and post a sqs message to upload a single document for a business
func (u uploader) BusinessSingleUpload(businessID shared.BusinessID, documentID shared.BusinessDocumentID) error {
	data := csp.BusinessSingleDocumentNotification{
		BusinessID: businessID,
		DocumentID: documentID,
	}
	body, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error marshalling data into %v error :%v ", data, err)
		return err
	}
	if err != nil {
		return err
	}

	sts := csp.DocumentStatusPending
	err = cspdoc.NewCSPBusinessService().CSPBusinessDocumentCreate(cspdoc.CSPBusinessDocument{
		DocumentID: &documentID,
		Status:     &sts,
	})
	if err != nil {
		return err
	}
	m := csp.Message{
		EntityID: string(businessID),
		Action:   csp.ActionUploadSingle,
		Category: csp.CategoryBusinessDocument,
		Data:     body,
	}

	return csp.SendDocumentMessage(m)

}

func (u uploader) BusinessSingleReUpload(businessID shared.BusinessID) error {
	m := csp.Message{
		EntityID: string(businessID),
		Action:   csp.ActionReUpload,
		Category: csp.CategoryBusinessDocument,
	}

	return csp.SendDocumentMessage(m)
}

// Handles api calls and post a sqs message to upload a single document for a consumer
func (u uploader) ConsumerSingleUpload(consumerID shared.ConsumerID, documentID shared.ConsumerDocumentID) error {

	data := csp.ConsumerSingleDocumentNotification{
		ConsumerID: consumerID,
		DocumentID: documentID,
	}
	body, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error marshalling data into %v error :%v ", data, err)
		return err
	}

	if err == nil {
		sts := csp.DocumentStatusPending
		cspdoc.NewCSPConsumerService().CSPConsumerDocumentCreate(cspdoc.CSPConsumerDocument{
			DocumentID: &documentID,
			Status:     &sts,
		})
	}

	m := csp.Message{
		EntityID: string(consumerID),
		Action:   csp.ActionUploadSingle,
		Category: csp.CategoryConsumerDocument,
		Data:     body,
	}

	return csp.SendDocumentMessage(m)

}

func (u uploader) ConsumerSingleReUpload(consumerID shared.ConsumerID) error {
	m := csp.Message{
		EntityID: string(consumerID),
		Action:   csp.ActionReUpload,
		Category: csp.CategoryConsumerDocument,
	}

	return csp.SendDocumentMessage(m)
}

// Business single document upload
func (u uploader) BusinessSingle(businessID shared.BusinessID, documentID shared.BusinessDocumentID) error {
	var d bizdoc
	err := coreDB.DBRead.Get(&d, "SELECT storage_key,content_type,doc_type,id  FROM business_document WHERE business_id = $1 AND id = $2 AND deleted IS NULL", businessID, documentID)
	if err != nil {
		log.Printf("Error getting document to upload %v", err)
		return err
	}
	log.Printf("Proccessing single document upload for document id: %s", documentID)
	_, err = u.sendBusinessContent(businessID, d)
	if err != nil {
		log.Printf("error sending document to bank %v", err)
	}
	return err
}

// Consumer single document upload
func (u uploader) ConsumerSingle(consumerID shared.ConsumerID, documentID shared.ConsumerDocumentID) error {
	var d condoc
	err := coreDB.DBRead.Get(&d, "SELECT storage_key,content_type,doc_type,id  FROM consumer_document WHERE consumer_id = $1 AND id = $2 AND deleted IS NULL", consumerID, documentID)
	if err != nil {
		log.Printf("Error getting the list of document to upload %v", err)
		return err
	}
	_, err = u.sendConsumerContent(consumerID, d)
	return err
}

// Send multipe document at once for a giving business
func (u uploader) BusinessMultiple(businessID shared.BusinessID) error {
	log.Printf("Proccessing document upload list for business id :%s", businessID)
	var err error
	var list []bizdoc
	err = coreDB.DBRead.Select(&list, "SELECT storage_key,content_type,doc_type,id  FROM business_document WHERE business_id = $1 AND deleted IS NULL", businessID)
	if err != nil {
		log.Printf("Error getting the list of document to upload %v", err)
		return err
	}

	var wg sync.WaitGroup
	log.Println("business document count to upload", len(list))
	for _, d := range list {
		wg.Add(1)
		go func(d bizdoc) {
			log.Printf("uploading doc %v", d)
			_, err = u.sendBusinessContent(businessID, d)
			wg.Done()
		}(d)
	}

	wg.Wait()
	return err
}

func (u uploader) ConsumerMultiple(consumerID shared.ConsumerID) error {
	var err error
	var list []condoc
	err = coreDB.DBRead.Select(&list, "SELECT storage_key,content_type,doc_type,id  FROM consumer_document WHERE user_id = $1 AND deleted IS NULL", consumerID)
	if err != nil {
		log.Printf("Error getting the list of document to upload %v", err)
		return err
	}

	var wg sync.WaitGroup
	log.Println("count", len(list))
	for _, d := range list {
		wg.Add(1)
		go func(d condoc) {
			log.Printf("uploading doc %v", d)
			_, err = u.sendConsumerContent(consumerID, d)
			wg.Done()
		}(d)
	}

	wg.Wait()
	return err
}

func (u uploader) sendBusinessContent(bizID shared.BusinessID, d bizdoc) (*partnerbank.IdentityDocumentResponse, error) {
	// check if this document has been already uploaded to BBVA
	if cspdoc.NewCSPBusinessService().CSPBusinessDocumentExist(d.ID) {
		log.Printf("Document with id %s has been submitted already", d.ID)
		return nil, nil
	}
	// Create csp document
	sts := csp.DocumentStatusPending
	cspdoc.NewCSPBusinessService().CSPBusinessDocumentCreate(cspdoc.CSPBusinessDocument{
		DocumentID: &d.ID,
		Status:     &sts,
	})

	s, err := docs.NewStorerFromKey(d.Key)
	c, err := s.Content()
	if err != nil {
		log.Printf("Error Sendind Document Content err: %v", err)
		return nil, err
	}

	bank, err := partnerbank.GetBusinessBank(partnerbank.ProviderNameBBVA)
	if err != nil {
		log.Printf("Error getting partnet bank  %v", err)
		return nil, err
	}

	srv := bank.BusinessEntityService(u.sourceReq.PartnerBankRequest())

	currentTime := time.Now()
	//Set default IDVs to formation if no items
	IDVs := businessIDVS(bizID)
	if len(IDVs) == 0 {
		req := businessDocumentContentRequest(bizID, d.ContentType, c, []partnerbank.IDVerify{partnerbank.IDVerifyFormationDoc}, d.DocType)
		resp, err := srv.UploadIdentityDocument(req)

		if err == nil {
			st := csp.DocumentStatusUploaded
			cspdoc.NewCSPBusinessService().CSPBusinessDocumentUpdate(d.ID, cspdoc.CSPBusinessDocument{
				BankDocumentID: &resp.Id,
				Status:         &st,
				Submitted:      &currentTime,
			})
		}
		if err != nil {
			st := csp.DocumentStatusFailed
			response := responseFromError(err)
			cspdoc.NewCSPBusinessService().CSPBusinessDocumentUpdate(d.ID, cspdoc.CSPBusinessDocument{
				DocumentID: &d.ID,
				Status:     &st,
				Response:   &response,
				Submitted:  &currentTime,
			})
		}

		return resp, err
	}

	req := businessDocumentContentRequest(bizID, d.ContentType, c, IDVs, d.DocType)
	resp, err := srv.UploadIdentityDocument(req)
	if err == nil {
		st := csp.DocumentStatusUploaded
		cspdoc.NewCSPBusinessService().CSPBusinessDocumentUpdate(d.ID, cspdoc.CSPBusinessDocument{
			DocumentID:     &d.ID,
			BankDocumentID: &resp.Id,
			Status:         &st,
			Submitted:      &currentTime,
		})
	}
	if err != nil {
		st := csp.DocumentStatusFailed
		response := responseFromError(err)
		cspdoc.NewCSPBusinessService().CSPBusinessDocumentUpdate(d.ID, cspdoc.CSPBusinessDocument{
			DocumentID: &d.ID,
			Status:     &st,
			Response:   &response,
			Submitted:  &currentTime,
		})
	}

	return resp, err
}

func (u uploader) sendConsumerContent(consumerID shared.ConsumerID, d condoc) (*partnerbank.IdentityDocumentResponse, error) {
	if cspdoc.NewCSPConsumerService().CSPConsumerDocumentExist(d.ID) {
		log.Printf("Document with id %s has been submitted already", d.ID)
		return nil, nil
	}
	store, err := docs.NewStorerFromKey(d.Key)
	c, err := store.Content()
	if err != nil {
		log.Printf("Error Sendind Document Content err: %v", err)
		return nil, err
	}

	bank, err := partnerbank.GetConsumerBank(partnerbank.ProviderNameBBVA)

	if err != nil {
		log.Printf("error getting consumer bank %v", err)
		return nil, err
	}

	srv := bank.ConsumerEntityService(u.sourceReq.PartnerBankRequest())

	currentTime := time.Now()
	//Set default IDVs to formation if no items
	IDVs := consumerIDVS(consumerID)
	if len(IDVs) == 0 {
		req := consumerDocumentContentRequest(consumerID, d.ContentType, c, []partnerbank.IDVerify{partnerbank.IDVerifyOther}, d.DocType)
		resp, err := srv.UploadIdentityDocument(req)
		if err == nil {
			st := csp.DocumentStatusUploaded
			cspdoc.NewCSPConsumerService().CSPConsumerDocumentUpdate(d.ID, cspdoc.CSPConsumerDocument{
				BankDocumentID: &resp.Id,
				Submitted:      &currentTime,
				Status:         &st,
			})
		}

		if err != nil {
			st := csp.DocumentStatusFailed
			response := responseFromError(err)
			cspdoc.NewCSPConsumerService().CSPConsumerDocumentUpdate(d.ID, cspdoc.CSPConsumerDocument{
				Status:    &st,
				Response:  &response,
				Submitted: &currentTime,
			})
		}
		return resp, err
	}

	req := consumerDocumentContentRequest(consumerID, d.ContentType, c, IDVs, d.DocType)
	resp, err := srv.UploadIdentityDocument(req)
	if err == nil {
		st := csp.DocumentStatusUploaded
		cspdoc.NewCSPConsumerService().CSPConsumerDocumentUpdate(d.ID, cspdoc.CSPConsumerDocument{
			BankDocumentID: &resp.Id,
			Submitted:      &currentTime,
			Status:         &st,
		})
	}

	if err != nil {
		st := csp.DocumentStatusFailed
		response := responseFromError(err)
		cspdoc.NewCSPConsumerService().CSPConsumerDocumentUpdate(d.ID, cspdoc.CSPConsumerDocument{
			Status:    &st,
			Response:  &response,
			Submitted: &currentTime,
		})
	}

	return resp, err
}

func businessDocumentContentRequest(bizID shared.BusinessID, ctype string, content []byte, identities []partnerbank.IDVerify, docType string) partnerbank.BusinessIdentityDocumentRequest {
	return partnerbank.BusinessIdentityDocumentRequest{
		BusinessID:       partnerbank.BusinessID(bizID),
		ContentType:      docContentType(ctype),
		IDVerifyRequired: identities,
		IdentityDocument: partnerbank.BusinessIdentityDocument(docType),
		Content:          content,
	}
}

func consumerDocumentContentRequest(consumerID shared.ConsumerID, ctype string, content []byte,
	identities []partnerbank.IDVerify, docType string) partnerbank.ConsumerIdentityDocumentRequest {
	return partnerbank.ConsumerIdentityDocumentRequest{
		ConsumerID:       partnerbank.ConsumerID(consumerID),
		IdentityDocument: partnerbank.ConsumerIdentityDocument(docType), //Convert consumer document type to partner identity document
		ContentType:      docContentType(ctype),
		IDVerifyRequired: identities,
		Content:          content,
	}
}

func docContentType(content string) partnerbank.ContentType {
	switch content {
	case "application/pdf":
		return partnerbank.ContentTypePDF
	case "image/png":
		return partnerbank.ContentTypePNG
	case "image/jpeg":
		return partnerbank.ContentTypeJPEG
	}
	return partnerbank.ContentTypePDF
}

func responseFromError(e error) types.JSONText {
	v := struct {
		ErrorMessage string `json:"errorMessage"`
	}{ErrorMessage: e.Error()}
	b, _ := json.Marshal(&v)
	return types.JSONText(b)
}

func (u uploader) ConsumerReUpload(consumerID shared.ConsumerID) error {
	c, err := docs.DownloadBBVAReQueueContent()
	if err != nil {
		return err
	}
	bank, err := partnerbank.GetConsumerBank(partnerbank.ProviderNameBBVA)

	if err != nil {
		return err
	}

	srv := bank.ConsumerEntityService(u.sourceReq.PartnerBankRequest())

	IDVs := consumerIDVS(consumerID)
	if len(IDVs) == 0 {
		req := consumerDocumentContentRequest(consumerID, "application/pdf", c, []partnerbank.IDVerify{partnerbank.IDVerifyOther}, "driverLicense")
		_, err := srv.UploadIdentityDocument(req)

		return err
	}

	req := consumerDocumentContentRequest(consumerID, "application/pdf", c, IDVs, "driversLicense")
	_, err = srv.UploadIdentityDocument(req)
	if err != nil {
		log.Printf("Error sendind document %v", err)
	}
	return err
}

func (u uploader) BusinessReUpload(bizID shared.BusinessID) error {
	c, err := docs.DownloadBBVAReQueueContent()
	if err != nil {
		return err
	}
	bank, err := partnerbank.GetBusinessBank(partnerbank.ProviderNameBBVA)
	if err != nil {
		return err
	}

	srv := bank.BusinessEntityService(u.sourceReq.PartnerBankRequest())

	//Set default IDVs to formation if no items
	IDVs := businessIDVS(bizID)
	if len(IDVs) == 0 {
		req := businessDocumentContentRequest(bizID, "application/pdf", c, []partnerbank.IDVerify{partnerbank.IDVerifyFormationDoc}, "other")
		_, err := srv.UploadIdentityDocument(req)
		return err
	}

	req := businessDocumentContentRequest(bizID, "application/pdf", c, IDVs, "other")
	_, err = srv.UploadIdentityDocument(req)

	return err
}
