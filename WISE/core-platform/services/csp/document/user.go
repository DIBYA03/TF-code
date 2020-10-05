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

// UserDocumentService document service
type UserDocumentService interface {
	//List documents
	List(consumerID string, offset int, limit int) ([]*UserDocument, error)

	//Get Document by id
	GetByID(shared.UserDocumentID) (*UserDocument, error)

	//Create a user document and return a signed url
	Create(consumerID string, create UserDocumentCreate) (*UserDocumentResponse, error)

	//Update a user document by id
	Update(shared.UserDocumentID, UserDocumentUpdate) (*UserDocumentResponse, error)

	//Signed url for download content
	SignedURL(shared.UserDocumentID) (*string, error)

	// Delete an user document by id

	Delete(shared.UserDocumentID) error

	Status(consumerID string, docID shared.UserDocumentID) (*Status, error)
}

type userDocumentService struct {
}

// NewUserDocumentService new user/consumer document service
func NewUserDocumentService() UserDocumentService {
	return userDocumentService{}
}

func (s userDocumentService) List(consumerID string, offset int, limit int) ([]*UserDocument, error) {
	list := make([]*UserDocument, 0)
	var userID string
	err := coreDB.DBRead.Get(&userID, "SELECT id FROM wise_user WHERE consumer_id = $1", consumerID)
	if err != nil {
		log.Printf("error getting user using consumer id %v", err)
		return list, err
	}
	err = coreDB.DBRead.Select(&list, "SELECT * FROM user_document WHERE user_id = $1 AND deleted IS NULL ORDER BY created ASC LIMIT $2 OFFSET $3", userID, limit, offset)
	if err == sql.ErrNoRows {
		log.Printf("no documents %v", err)
		return list, nil
	}
	return list, err
}

func (s userDocumentService) GetByID(docID shared.UserDocumentID) (*UserDocument, error) {
	u := UserDocument{}

	err := coreDB.DBRead.Get(&u, "SELECT * FROM user_document WHERE id = $1", docID)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (s userDocumentService) SignedURL(docID shared.UserDocumentID) (*string, error) {
	var key string
	notFound := services.ErrorNotFound{}.New("")
	err := coreDB.DBRead.Get(&key, "SELECT storage_key from user_document WHERE id = $1 ", docID)
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

func (s userDocumentService) Create(consumerID string, doc UserDocumentCreate) (*UserDocumentResponse, error) {
	var userID string
	err := coreDB.DBRead.Get(&userID, "SELECT id FROM wise_user WHERE consumer_id = $1", consumerID)
	if err != nil {
		log.Printf("error getting user using consumer id %v", err)
		return nil, err
	}
	//add user id to document
	doc.UserID = &userID

	storer, err := docsrv.NewAWSS3DocStorage(userID, docsrv.UserPrefix)

	key, err := storer.Key()
	if key == nil || err != nil {
		return nil, err
	}

	doc.StorageKey = key

	keys := services.SQLGenInsertKeys(doc)
	values := services.SQLGenInsertValues(doc)
	var insertedDoc UserDocument

	q := fmt.Sprintf("INSERT INTO user_document (%s) VALUES(%s) RETURNING *", keys, values)
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

	return &UserDocumentResponse{
		SignedURL: url,
		Document:  insertedDoc,
	}, nil
}

func (s userDocumentService) Update(docID shared.UserDocumentID, doc UserDocumentUpdate) (*UserDocumentResponse, error) {
	var document UserDocument
	var updateContent bool

	if doc.UpdatingContent != nil {
		updateContent = *doc.UpdatingContent
		doc.UpdatingContent = nil
	}
	keys := services.SQLGenForUpdate(doc)

	q := fmt.Sprintf("UPDATE user_document SET %s WHERE id = '%s' RETURNING *", keys, docID)
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

	return &UserDocumentResponse{
		SignedURL: url,
		Document:  document,
	}, nil
}

func (s userDocumentService) Delete(docID shared.UserDocumentID) error {
	_, err := coreDB.DBWrite.Exec("UPDATE user_document SET deleted = CURRENT_TIMESTAMP WHERE id = $1", docID)
	if err == sql.ErrNoRows {
		return services.ErrorNotFound{}.New("")
	}
	return nil
}

func (s userDocumentService) Status(consumerID string, docID shared.UserDocumentID) (*Status, error) {
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
