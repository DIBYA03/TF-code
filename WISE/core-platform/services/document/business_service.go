/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

package document

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/wiseco/core-platform/auth"
	"github.com/wiseco/core-platform/services"
	"github.com/wiseco/core-platform/services/data"
	"github.com/wiseco/core-platform/shared"
)

//BusinessDocumentService ...
type BusinessDocumentService interface {

	//Create handles the creation of a business documents
	// params: `businessID` and `BusinessDocumentCreate`
	Create(BusinessDocumentCreate) (*DocumentCreateResponse, error)
	CreateInternal(BusinessDocumentCreate) (*BusinessDocument, error)

	//GetByID gets a business document by id
	//Params: `businessID`
	GetByID(shared.BusinessDocumentID, shared.BusinessID) (*BusinessDocument, error)

	//Update updates a document by id
	Update(BusinessDocumentUpdate) (*DocumentCreateResponse, error)

	//List will list return a list of documents for the business
	//Params: `businessID`, `orderBy`, `limit` `offset`
	List(shared.BusinessID, string, int, int) ([]BusinessDocument, error)

	//SignULR gets the signed url for a business document
	SignedURL(shared.BusinessID, shared.BusinessDocumentID) (*string, error)
}

type businessDocumentStore struct {
	sourceReq services.SourceRequest
	*sqlx.DB
}

// NewBusinessDocumentService ..
func NewBusinessDocumentService(r services.SourceRequest) BusinessDocumentService {
	return &businessDocumentStore{r, data.DBWrite}
}

//NewDocumentService creates a business document service without a source request
func NewDocumentService() BusinessDocumentService {
	return &businessDocumentStore{services.SourceRequest{}, data.DBWrite}
}

func (store businessDocumentStore) uploadDocument(businessID, contentType, docType string, content *string) (*string, error) {

	storer, err := NewStorerFromContent(StorageProviderAWSS3, BusinessPrefix, businessID, contentType, content)
	if err != nil {
		log.Printf("Error uploading document to s3 err: %v", err)
		return nil, err
	}

	return storer.Key()
}

func (store businessDocumentStore) Create(doc BusinessDocumentCreate) (*DocumentCreateResponse, error) {
	// Check access
	err := auth.NewAuthService(store.sourceReq).CheckBusinessAccess(doc.BusinessID)
	if err != nil {
		return nil, err
	}

	if doc.DocType != nil {
		if _, ok := NewDocType(*doc.DocType); !ok {
			return nil, errors.New("invalid business document type")
		}
	}

	storer, err := NewAWSS3DocStorage(string(doc.BusinessID), BusinessPrefix)

	key, err := storer.Key()
	if err != nil {
		return nil, err
	}
	doc.StorageKey = key

	keys := services.SQLGenInsertKeys(doc)
	values := services.SQLGenInsertValues(doc)
	var insertedDoc BusinessDocument

	q := fmt.Sprintf("INSERT INTO business_document (%s) VALUES(%s) RETURNING *", keys, values)
	stmt, err := store.PrepareNamed(q)
	if err != nil {
		return nil, err
	}
	err = stmt.Get(&insertedDoc, doc)
	if err != nil {
		return nil, err
	}

	// Update business
	signed, err := storer.PutSignedURL()
	if err != nil || signed == nil {
		return nil, err
	}
	return &DocumentCreateResponse{SignedURL: signed, Document: insertedDoc}, err
}

func (store businessDocumentStore) CreateInternal(doc BusinessDocumentCreate) (*BusinessDocument, error) {
	keys := services.SQLGenInsertKeys(doc)
	values := services.SQLGenInsertValues(doc)
	var insertedDoc BusinessDocument

	q := fmt.Sprintf("INSERT INTO business_document (%s) VALUES(%s) RETURNING *", keys, values)
	stmt, err := store.PrepareNamed(q)
	if err != nil {
		return nil, err
	}

	err = stmt.Get(&insertedDoc, doc)
	if err != nil {
		return nil, err
	}

	return &insertedDoc, nil
}

func (store businessDocumentStore) GetByID(id shared.BusinessDocumentID, businessID shared.BusinessID) (*BusinessDocument, error) {
	// Check access
	err := auth.NewAuthService(store.sourceReq).CheckBusinessAccess(businessID)
	if err != nil {
		return nil, err
	}

	q := `SELECT * FROM business_document WHERE id = $1 AND business_id = $2`
	var document BusinessDocument
	err = store.Get(&document, q, id, businessID)
	if err != nil && err != sql.ErrNoRows {
		return nil, services.ErrorNotFound{}.New("")
	}
	return &document, err
}

func (store businessDocumentStore) Update(doc BusinessDocumentUpdate) (*DocumentCreateResponse, error) {
	if err := auth.NewAuthService(store.sourceReq).CheckBusinessAccess(doc.BusinessID); err != nil {
		log.Println("error getting authorization access")
		return nil, err
	}

	var document BusinessDocument
	keys := services.SQLGenForUpdate(doc)

	q := fmt.Sprintf("UPDATE business_document SET %s WHERE id = '%s' RETURNING *", keys, doc.ID)
	stmt, err := store.PrepareNamed(q)
	if err != nil {
		log.Printf("error preparing stmt %v query %v", err, q)
		return nil, err
	}

	err = stmt.Get(&document, doc)
	if err != nil {
		return nil, fmt.Errorf("error keys: %v err: %v", keys, err)
	}

	// Update business
	storer, err := NewStorerFromKey(*document.StorageKey)
	url, err := storer.PutSignedURL()
	if err != nil {
		log.Printf("error creating signed url but updating was ok %v", err)
	}
	return &DocumentCreateResponse{
		SignedURL: url,
		Document:  document,
	}, nil
}

func (store businessDocumentStore) List(businessID shared.BusinessID, orderBy string, limit, offset int) ([]BusinessDocument, error) {
	// Check access
	err := auth.NewAuthService(store.sourceReq).CheckBusinessAccess(businessID)
	if err != nil {
		return nil, err
	}

	docs := make([]BusinessDocument, 0)
	q := `SELECT * FROM business_document WHERE business_id = $1 AND deleted IS NULL`

	err = store.Select(&docs, q, businessID)
	if err != nil && err == sql.ErrNoRows {
		return docs, nil
	}

	return docs, err
}

func (store businessDocumentStore) SignedURL(businessID shared.BusinessID, docID shared.BusinessDocumentID) (*string, error) {
	// Check access
	err := auth.NewAuthService(store.sourceReq).CheckBusinessAccess(businessID)
	if err != nil {
		return nil, err
	}

	var key string
	notFound := services.ErrorNotFound{}.New("")
	err = store.Get(&key, "SELECT storage_key from business_document WHERE id = $1 AND business_id = $2", docID, businessID)
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
