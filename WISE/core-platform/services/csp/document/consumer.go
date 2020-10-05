package document

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/wiseco/core-platform/services"
	cspDB "github.com/wiseco/core-platform/services/csp/data"
	coreDB "github.com/wiseco/core-platform/services/data"
	docsrv "github.com/wiseco/core-platform/services/document"
	"github.com/wiseco/core-platform/shared"
)

// ConsumerDocumentService document service
type ConsumerDocumentService interface {
	//List documents
	List(consumerID shared.ConsumerID, offset int, limit int) ([]ConsumerDocument, error)

	//Get Document by id
	GetByID(shared.ConsumerDocumentID) (*ConsumerDocument, error)

	//Create a user document and return a signed url
	Create(consumerID shared.ConsumerID, create ConsumerDocumentCreate) (*ConsumerDocumentResponse, error)

	//Update a user document by id
	Update(shared.ConsumerDocumentID, ConsumerDocumentUpdate) (*ConsumerDocumentResponse, error)

	//Signed url for download content
	SignedURL(shared.ConsumerDocumentID) (*string, error)

	Delete(shared.ConsumerDocumentID) error

	Status(consumerID shared.ConsumerID, docID shared.ConsumerDocumentID) (*Status, error)

	CreateFromBusiness(shared.ConsumerID, shared.BusinessID) error
}

type consumerDocumentService struct {
}

// NewConsumerDocumentService new user/consumer document service
func NewConsumerDocumentService() ConsumerDocumentService {
	return consumerDocumentService{}
}

func (s consumerDocumentService) List(consumerID shared.ConsumerID, offset int, limit int) ([]ConsumerDocument, error) {
	list := make([]ConsumerDocument, 0)

	err := coreDB.DBRead.Select(&list, "SELECT * FROM consumer_document WHERE consumer_id = $1 AND deleted IS NULL ORDER BY created ASC LIMIT $2 OFFSET $3", consumerID, limit, offset)
	if err == sql.ErrNoRows {
		log.Printf("no documents %v", err)
		return list, nil
	}
	return list, err
}

func (s consumerDocumentService) GetByID(docID shared.ConsumerDocumentID) (*ConsumerDocument, error) {
	var doc ConsumerDocument
	err := coreDB.DBRead.Get(&doc, "SELECT * FROM consumer_document WHERE id = $1", docID)
	if err != nil {
		return nil, err
	}

	return &doc, nil
}

func (s consumerDocumentService) SignedURL(docID shared.ConsumerDocumentID) (*string, error) {
	var key string
	notFound := services.ErrorNotFound{}.New("")
	err := coreDB.DBRead.Get(&key, "SELECT storage_key from consumer_document WHERE id = $1 ", docID)
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

	storer, err := docsrv.NewStorerFromKey(key)

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

func (s consumerDocumentService) Create(consumerID shared.ConsumerID, doc ConsumerDocumentCreate) (*ConsumerDocumentResponse, error) {

	//add consumer id to document
	doc.ConsumerID = &consumerID

	storer, err := docsrv.NewAWSS3DocStorage(string(consumerID), docsrv.ConsumerPrefix)

	key, err := storer.Key()
	if key == nil || err != nil {
		return nil, err
	}

	doc.StorageKey = key

	keys := services.SQLGenInsertKeys(doc)
	values := services.SQLGenInsertValues(doc)
	var insertedDoc ConsumerDocument

	q := fmt.Sprintf("INSERT INTO consumer_document (%s) VALUES(%s) RETURNING *", keys, values)
	stmt, err := coreDB.DBWrite.PrepareNamed(q)
	if err != nil {
		return nil, err
	}

	err = stmt.Get(&insertedDoc, doc)
	if err != nil {
		return nil, err
	}

	url, err := storer.PutSignedURL()
	if err != nil {
		log.Printf("error getting pre-signed url")
		return nil, err
	}

	return &ConsumerDocumentResponse{
		SignedURL: url,
		Document:  insertedDoc,
	}, nil
}

func (s consumerDocumentService) Update(docID shared.ConsumerDocumentID, doc ConsumerDocumentUpdate) (*ConsumerDocumentResponse, error) {
	var document ConsumerDocument
	var updateContent bool

	if doc.UpdatingContent != nil {
		updateContent = *doc.UpdatingContent
		doc.UpdatingContent = nil
	}
	keys := services.SQLGenForUpdate(doc)

	q := fmt.Sprintf("UPDATE consumer_document SET %s WHERE id = '%s' RETURNING *", keys, docID)
	stmt, err := coreDB.DBWrite.PrepareNamed(q)
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
		storer, _ := docsrv.NewStorerFromKey(*document.StorageKey)
		newURL, _ := storer.PutSignedURL()
		url = newURL
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

func (s consumerDocumentService) Delete(docID shared.ConsumerDocumentID) error {
	_, err := coreDB.DBWrite.Exec("UPDATE consumer_document SET deleted = CURRENT_TIMESTAMP WHERE id = $1", docID)
	if err == sql.ErrNoRows {
		return services.ErrorNotFound{}.New("")
	}
	return nil
}

func (s consumerDocumentService) Status(consumerID shared.ConsumerID, docID shared.ConsumerDocumentID) (*Status, error) {
	var status Status
	err := cspDB.DBWrite.Get(&status, "SELECT * FROM consumer_document WHERE document_id = $1", docID)
	if err != nil && err == sql.ErrNoRows {
		return nil, services.ErrorNotFound{}.New("")
	}
	if err != nil {
		log.Printf("Error getting doc sts %v", err)
		return nil, err
	}
	return &status, nil
}

// assuming this business id is a sole Prop and we only need to fetch the documents from it
// fetching 30 docs and finding one with the driverLicense type
func (s consumerDocumentService) CreateFromBusiness(consumerID shared.ConsumerID, businessID shared.BusinessID) error {
	docs, err := NewDocumentService().List(businessID, 30, 0)
	if err != nil {
		return err
	}
	var consumerDoc ConsumerDocumentCreate
	hasDoc := false
	var businessDocKey *string
	for _, doc := range docs {
		if doc.DocType == nil {
			continue
		}
		if *doc.DocType == driversLicenseDocType {
			consumerDoc.Number = doc.Number
			docType := docsrv.ConsumerIdentityDocumentDriversLicense
			consumerDoc.DocType = &docType
			businessDocKey = doc.StorageKey
			consumerDoc.ExpirationDate = doc.ExpirationDate
			// set doc to upload and break out, we only want one doc
			hasDoc = true
			break
		}
	}
	// return when no doc to create
	if !hasDoc {
		return nil
	}
	// return if we dont have a storage key to fetch doc content
	if businessDocKey == nil {
		return nil
	}
	// create and upload consumer document
	return s.uploadAndCreate(consumerDoc, *businessDocKey, consumerID)
}

func (s consumerDocumentService) uploadAndCreate(doc ConsumerDocumentCreate, businessDocKey string, ownerID shared.ConsumerID) error {
	storer, err := docsrv.NewAWSS3DocStoreFromKey(businessDocKey)
	if err != nil {
		return err
	}
	buf, err := storer.ContentBuffer()
	if err != nil {
		return err
	}

	storageKey, err := storer.UploadBuffer(buf, docsrv.Options{
		Prefix:  docsrv.ConsumerPrefix,
		OwnerID: string(ownerID),
	})
	if err != nil {
		return err
	}
	doc.StorageKey = &storageKey
	_, err = s.Create(ownerID, doc)
	return err
}
