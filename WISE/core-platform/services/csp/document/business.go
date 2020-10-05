package document

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/wiseco/core-platform/services"
	"github.com/wiseco/core-platform/services/data"
	docs "github.com/wiseco/core-platform/services/document"
	"github.com/wiseco/core-platform/shared"
)

// BusinessDocumentService the business document service
type BusinessDocumentService interface {

	//Create document
	Create(businessID shared.BusinessID, doc BusinessDocumentCreate) (*BusinessDocumentResponse, error)

	//Update document
	Update(businessID shared.BusinessID, docID shared.BusinessDocumentID, updates BusinessDocumentUpdate) (*BusinessDocumentResponse, error)

	//Get document by id
	GetByID(businessID shared.BusinessID, docID shared.BusinessDocumentID) (BusinessDocument, error)

	//Delete document by id
	Delete(businessID shared.BusinessID, docID shared.BusinessDocumentID) error

	//Get all documents by busines id
	List(businessID shared.BusinessID, limit, offset int) ([]BusinessDocument, error)

	// Signed URL of the docuement from s3
	URL(businessID shared.BusinessID, docID shared.BusinessDocumentID) (*string, error)

	SetAsFormation(businessID shared.BusinessID, docID shared.BusinessDocumentID, removing bool) error

	CreateFromConsumer(shared.ConsumerID, shared.BusinessID, shared.UserID) error
}

type documentService struct {
	*sqlx.DB
}

//NewDocumentService ..
func NewDocumentService() BusinessDocumentService {
	return documentService{data.DBWrite}
}

func (s documentService) Create(businessID shared.BusinessID, create BusinessDocumentCreate) (*BusinessDocumentResponse, error) {
	var doc BusinessDocument
	var err error
	var key *string
	var formation bool

	storer, err := docs.NewAWSS3DocStorage(string(businessID), docs.BusinessPrefix)
	if err != nil {
		log.Printf("document upload failed error:%v", err)
		return nil, err
	}

	key, _ = storer.Key()
	if create.UseFormation != nil {
		formation = true
		create.UseFormation = nil
	}

	create.StorageKey = key
	keys := services.SQLGenInsertKeys(create)
	values := services.SQLGenInsertValues(create)

	q := fmt.Sprintf("INSERT INTO business_document (%s) VALUES(%s) RETURNING *", keys, values)
	if stmt, err := s.PrepareNamed(q); err != nil {
		log.Printf("error preparing stmt %v", err)
		return nil, err
	} else if err := stmt.Get(&doc, create); err != nil {
		log.Printf("error getting doc from stmt %v", err)
		return nil, err
	}
	// replace business formation document
	if formation {
		q := "UPDATE business SET formation_document_id = $1 WHERE id = $2"
		_, err = s.Exec(q, doc.ID, businessID)
		create.UseFormation = nil
	}

	url, err := storer.PutSignedURL()
	if err != nil {
		log.Printf("Error getting pre-signed url:%v", err)
		return nil, err
	}
	return &BusinessDocumentResponse{
		SignedURL: url,
		Document:  doc,
	}, nil
}

func (s documentService) fromConsumer(businessID shared.BusinessID, create BusinessDocumentCreate) (*BusinessDocumentResponse, error) {
	var doc BusinessDocument
	var err error

	keys := services.SQLGenInsertKeys(create)
	values := services.SQLGenInsertValues(create)

	q := fmt.Sprintf("INSERT INTO business_document (%s) VALUES(%s) RETURNING *", keys, values)
	if stmt, err := s.PrepareNamed(q); err != nil {
		log.Printf("error preparing stmt %v", err)
		return nil, err
	} else if err := stmt.Get(&doc, create); err != nil {
		log.Printf("error getting doc from stmt %v", err)
		return nil, err
	}

	return &BusinessDocumentResponse{
		SignedURL: nil,
		Document:  doc,
	}, err
}

func (s documentService) Update(businessID shared.BusinessID, docID shared.BusinessDocumentID, updates BusinessDocumentUpdate) (*BusinessDocumentResponse, error) {
	var doc BusinessDocument
	var err error
	var formation bool
	var updateContent bool
	var key *string

	storer, err := docs.NewAWSS3DocStorage(string(businessID), docs.BusinessPrefix)
	if err != nil {
		log.Printf("document upload failed error:%v", err)
		return nil, err
	}

	var url *string
	if updates.UpdatingContent != nil {
		updateContent = *updates.UpdatingContent
		updates.UpdatingContent = nil
	}
	if updateContent {
		newURL, err := storer.PutSignedURL()
		if err != nil {
			log.Printf("Error getting pre-signed url:%v", err)
			return nil, err
		}
		url = newURL
		key, _ = storer.Key()
		updates.StorageKey = key
	}
	if updates.UseFormation != nil {
		formation = *updates.UseFormation
		updates.UseFormation = nil
	}

	keys := services.SQLGenForUpdate(updates)

	q := fmt.Sprintf("UPDATE business_document SET %s WHERE id = '%s' RETURNING *", keys, docID)
	stmt, err := s.PrepareNamed(q)
	if err != nil {
		log.Printf("error preparing stmt %v query %v", err, q)
		return nil, err
	}

	err = stmt.Get(&doc, updates)
	if err != nil {
		return nil, fmt.Errorf("error keys: %v err: %v", keys, err)
	}

	if formation {
		//Update business
		q := "UPDATE business SET formation_document_id = $1 WHERE id = $2"
		_, err = s.Exec(q, docID, businessID)
	}

	return &BusinessDocumentResponse{
		SignedURL: url,
		Document:  doc,
	}, nil

}

func (s documentService) GetByID(businessID shared.BusinessID, docID shared.BusinessDocumentID) (BusinessDocument, error) {
	q := `SELECT * FROM business_document WHERE id = $1 AND business_id = $2`
	var document BusinessDocument
	err := s.Get(&document, q, docID, businessID)
	if err != nil && err != sql.ErrNoRows {
		return BusinessDocument{}, services.ErrorNotFound{}.New("")
	}
	return document, err
}

func (s documentService) Delete(businessID shared.BusinessID, docID shared.BusinessDocumentID) error {
	_, err := s.Exec("UPDATE business_document SET deleted = CURRENT_TIMESTAMP WHERE id = $1", docID)
	if err == sql.ErrNoRows {
		return services.ErrorNotFound{}.New("")
	}
	s.checkIsFormation(businessID, docID)
	return err
}

func (s documentService) checkIsFormation(businessID shared.BusinessID, docID shared.BusinessDocumentID) {
	b := struct {
		Formation *shared.BusinessDocumentID `db:"formation_document_id"`
	}{}
	if err := s.Get(&b, "SELECT formation_document_id FROM business WHERE id = $1", businessID); err != nil {
		log.Printf("unable to get business by id %v", err)
	}
	if b.Formation == nil {
		return
	}
	if *b.Formation == docID {
		_, err := s.Exec("UPDATE business SET formation_document_id  = 'NULL' WHERE id = $1", businessID)
		if err != nil {
			log.Printf("error updating business formation_document_id %v", err)
		}
	}

}

func (s documentService) List(businessID shared.BusinessID, limit, offset int) ([]BusinessDocument, error) {
	list := make([]BusinessDocument, 0)
	q := `SELECT * FROM business_document WHERE business_id = $1 AND deleted IS NULL`

	err := s.Select(&list, q, businessID)
	if err != nil && err == sql.ErrNoRows {
		return []BusinessDocument{}, nil
	}

	return list, err
}

func (s documentService) URL(businessID shared.BusinessID, docID shared.BusinessDocumentID) (*string, error) {
	var key string
	notFound := services.ErrorNotFound{}.New("")
	err := s.Get(&key, "SELECT storage_key from business_document WHERE id = $1 AND business_id = $2", docID, businessID)
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

	storer, err := docs.NewStorerFromKey(key)

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

func (s documentService) SetAsFormation(businessID shared.BusinessID, docID shared.BusinessDocumentID, removing bool) error {
	if removing {
		q := "UPDATE business SET formation_document_id = 'null' WHERE id = $1"
		_, err := s.Exec(q, businessID)
		return err
	}
	q := "UPDATE business SET formation_document_id = $1 WHERE id = $2"
	_, err := s.Exec(q, docID, businessID)
	return err
}

func uploadDocument(businessID shared.BusinessID, contentType, docType string, content *string) (*string, error) {
	storer, err := docs.NewStorerFromContent(docs.StorageProviderAWSS3, docs.BusinessPrefix, string(businessID), contentType, content)
	if err != nil {
		log.Printf("Error uploading document to s3 err: %v", err)
		return nil, err
	}

	return storer.Key()
}

func (s documentService) CreateFromConsumer(consumerID shared.ConsumerID, businessID shared.BusinessID, userID shared.UserID) error {

	docs, err := NewConsumerDocumentService().List(consumerID, 0, 30)
	if err != nil {
		return err
	}
	var businessDoc BusinessDocumentCreate
	hasDoc := false
	var consumerDocKey *string
	for _, doc := range docs {
		if doc.DocType == nil {
			continue
		}

		if *doc.DocType == driversLicenseDocType || *doc.DocType == passportDocType {
			businessDoc.Number = doc.Number
			businessDoc.DocType = doc.DocType
			consumerDocKey = doc.StorageKey
			businessDoc.BusinessID = &businessID
			businessDoc.CreatedUserID = userID
			businessDoc.ExpirationDate = doc.ExpirationDate
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
	if consumerDocKey == nil {
		return nil
	}
	// create and upload consumer document
	return s.uploadAndCreate(businessDoc, *consumerDocKey, businessID)
}

func (s documentService) uploadAndCreate(doc BusinessDocumentCreate, businessDocKey string, ownerID shared.BusinessID) error {
	storer, err := docs.NewAWSS3DocStoreFromKey(businessDocKey)
	if err != nil {
		return err
	}
	buf, err := storer.ContentBuffer()
	if err != nil {
		return err
	}

	storageKey, err := storer.UploadBuffer(buf, docs.Options{
		Prefix:  docs.BusinessPrefix,
		OwnerID: string(ownerID),
	})
	if err != nil {
		return err
	}
	doc.StorageKey = &storageKey
	_, err = s.fromConsumer(ownerID, doc)
	return err
}
