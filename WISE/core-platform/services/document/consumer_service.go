/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

// Package for all user document related services
package document

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/wiseco/core-platform/auth"
	"github.com/wiseco/core-platform/services"
	"github.com/wiseco/core-platform/services/data"
	"github.com/wiseco/core-platform/shared"
)

type consumerDocumentDatastore struct {
	sourceReq services.SourceRequest
	*sqlx.DB
}

type ConsumerDocumentService interface {
	//Get Count of documents
	Count() (int, error)

	//List documents
	List(consumerID shared.ConsumerID, offset int, limit int) ([]ConsumerDocument, error)

	//List documents
	ListInternal(consumerID shared.ConsumerID, offset int, limit int) ([]ConsumerDocument, error)

	//Get Document by id
	GetByID(shared.ConsumerID, shared.ConsumerDocumentID) (*ConsumerDocument, error)

	//Create a user document and return a signed url
	Create(shared.ConsumerID, ConsumerDocumentCreate) (*ConsumerDocumentResponse, error)

	//Update a user document by id
	Update(shared.ConsumerID, ConsumerDocumentUpdate) (*ConsumerDocumentResponse, error)

	// Delete a user document by id
	Delete(shared.ConsumerID, shared.ConsumerDocumentID) error

	//Signed url for download content
	SignedURL(shared.ConsumerID, shared.ConsumerDocumentID) (*string, error)
}

func NewConsumerDocumentService(r services.SourceRequest) ConsumerDocumentService {
	return &consumerDocumentDatastore{r, data.DBWrite}
}

func (db *consumerDocumentDatastore) Count() (int, error) {
	var count int

	err := db.QueryRow("SELECT COUNT(*) FROM consumer_document").Scan(&count)

	if err != nil {
		log.Println(err)
		return 0, errors.Cause(err)
	}

	return count, err
}

func (db *consumerDocumentDatastore) List(consumerID shared.ConsumerID, offset int, limit int) ([]ConsumerDocument, error) {
	// Check access
	err := auth.NewAuthService(db.sourceReq).CheckConsumerAccess(consumerID)
	if err != nil {
		return nil, err
	}

	list := make([]ConsumerDocument, 0)

	err = db.Select(&list, "SELECT * FROM consumer_document WHERE consumer_id = $1 AND deleted IS NULL ORDER BY created ASC LIMIT $2 OFFSET $3", consumerID, limit, offset)
	if err != nil && err == sql.ErrNoRows {
		log.Printf("no documents %v", err)
		return list, nil
	}
	return list, err
}

func (db *consumerDocumentDatastore) ListInternal(consumerID shared.ConsumerID, offset int, limit int) ([]ConsumerDocument, error) {
	list := make([]ConsumerDocument, 0)

	err := db.Select(&list, "SELECT * FROM consumer_document WHERE consumer_id = $1 AND deleted IS NULL ORDER BY created ASC LIMIT $2 OFFSET $3", consumerID, limit, offset)
	if err != nil && err == sql.ErrNoRows {
		log.Printf("no documents %v", err)
		return list, nil
	}
	return list, err
}

func (db *consumerDocumentDatastore) GetByID(consumerID shared.ConsumerID, docID shared.ConsumerDocumentID) (*ConsumerDocument, error) {
	// Check access
	err := auth.NewAuthService(db.sourceReq).CheckConsumerAccess(consumerID)
	if err != nil {
		return nil, err
	}

	u := ConsumerDocument{}

	err = db.Get(&u, "SELECT * FROM consumer_document WHERE id = $1", docID)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (db *consumerDocumentDatastore) SignedURL(consumerID shared.ConsumerID, docID shared.ConsumerDocumentID) (*string, error) {
	// Check access
	err := auth.NewAuthService(db.sourceReq).CheckConsumerAccess(consumerID)
	if err != nil {
		return nil, err
	}

	var key string

	notFound := services.ErrorNotFound{}.New("")

	err = db.Get(&key, "SELECT storage_key from consumer_document WHERE id = $1 ", docID)
	if err != nil && err == sql.ErrNoRows {
		return nil, notFound
	}

	if key == "" {
		return nil, notFound
	}

	if err != nil {
		log.Printf("Error getting document storage_key  error:%v", err)
		return nil, err
	}

	storer, err := NewStorerFromKey(key)
	if err != nil {
		log.Printf("Error creating storer  error:%v", err)
		return nil, err
	}

	url, err := storer.SignedUrl()
	if url == nil {
		log.Printf("no url url:%v err:%v", url, err)
		return nil, notFound
	}

	return url, err
}

func (db *consumerDocumentDatastore) Create(consumerID shared.ConsumerID, doc ConsumerDocumentCreate) (*ConsumerDocumentResponse, error) {
	// Check access
	err := auth.NewAuthService(db.sourceReq).CheckConsumerAccess(consumerID)
	if err != nil {
		return nil, err
	}

	if doc.DocType == nil {
		return nil, errors.New("Doc type is required")
	}

	_, ok := ConsumerDocTypeToBankDocType[*doc.DocType]
	if !ok {
		return nil, errors.New("Invalid doc type")
	}

	if doc.ID != nil {
		d, err := db.GetByID(consumerID, *doc.ID)
		if err != nil {
			return nil, err
		}

		doc.Number = d.Number
		doc.IssuingState = d.IssuingState
		doc.IssuingCountry = d.IssuingCountry
		doc.IssuedDate = d.IssuedDate
		doc.ExpirationDate = d.ExpirationDate
		doc.StorageKey = d.StorageKey
		doc.ContentType = d.ContentType
	}

	doc.ConsumerID = &consumerID

	var url *string
	// Create signed url only when document ID is not present
	if doc.ID == nil {
		storer, err := NewAWSS3DocStorage(string(*doc.ConsumerID), ConsumerPrefix)

		key, err := storer.Key()
		if err != nil {
			return nil, err
		}

		doc.StorageKey = key

		url, err = storer.PutSignedURL()
		if err != nil {
			log.Printf("error getting pre-signed url")
			return nil, err
		}
	}

	doc.ID = nil

	keys := services.SQLGenInsertKeys(doc)
	values := services.SQLGenInsertValues(doc)
	var insertedDoc ConsumerDocument

	q := fmt.Sprintf("INSERT INTO consumer_document (%s) VALUES(%s) RETURNING *", keys, values)
	stmt, err := db.PrepareNamed(q)
	if err != nil {
		return nil, err
	}

	err = stmt.Get(&insertedDoc, doc)
	if err != nil {
		return nil, err
	}

	return &ConsumerDocumentResponse{
		SignedURL: url,
		Document:  insertedDoc,
	}, nil
}

func (db *consumerDocumentDatastore) Update(consumerID shared.ConsumerID, doc ConsumerDocumentUpdate) (*ConsumerDocumentResponse, error) {
	// Check access
	err := auth.NewAuthService(db.sourceReq).CheckConsumerAccess(consumerID)
	if err != nil {
		return nil, err
	}

	var document ConsumerDocument
	var updateContent bool

	if doc.DocType != nil {
		_, ok := ConsumerDocTypeToBankDocType[*doc.DocType]
		if !ok {
			return nil, errors.New("Invalid doc type")
		}
	}

	if doc.UpdatingContent != nil {
		updateContent = *doc.UpdatingContent
		doc.UpdatingContent = nil
	}
	keys := services.SQLGenForUpdate(doc)

	q := fmt.Sprintf("UPDATE consumer_document SET %s WHERE id = '%s' RETURNING *", keys, doc.ID)
	stmt, err := db.PrepareNamed(q)
	if err != nil {
		log.Printf("error preparing stmt %v query %v", err, q)
		return nil, err
	}
	err = stmt.Get(&document, doc)
	if err != nil {
		return nil, fmt.Errorf("error keys: %v err: %v", keys, err)
	}

	var url *string
	if updateContent {
		storer, _ := NewStorerFromKey(*document.StorageKey)
		url, _ = storer.PutSignedURL()
	}

	if err != nil {
		log.Printf("error creating signed url but updating was ok %v", err)
		return nil, err
	}
	return &ConsumerDocumentResponse{
		SignedURL: url,
		Document:  document,
	}, nil
}

func (db *consumerDocumentDatastore) Delete(consumerID shared.ConsumerID, docID shared.ConsumerDocumentID) error {
	// Check access
	err := auth.NewAuthService(db.sourceReq).CheckConsumerAccess(consumerID)
	if err != nil {
		return err
	}

	// Delete only if consumer kyc status is in submitted or notStarted state
	var kycStatus services.KYCStatus
	if err != nil {
		err := db.Get(&kycStatus, "SELECT kyc_status FROM consumer WHERE id = $1", consumerID)
		return err
	}

	switch kycStatus {
	case services.KYCStatusNotStarted:
		fallthrough
	case services.KYCStatusSubmitted:
		_, err = db.Exec("UPDATE consumer_document SET deleted = CURRENT_TIMESTAMP WHERE id = $1", docID)
		if err == sql.ErrNoRows {
			return services.ErrorNotFound{}.New("")
		}
	default:
		return errors.New("Consumer kyc status doesn't allow deleting documents")
	}

	return nil
}
