package document

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/wiseco/core-platform/services"
	"github.com/wiseco/core-platform/services/csp"
	"github.com/wiseco/core-platform/services/csp/data"
	"github.com/wiseco/core-platform/shared"
)

// CSPBusinessService ..
type CSPBusinessService interface {
	CSPBusinessDocumentCreate(CSPBusinessDocument) error
	CSPBusinessDocumentUpdate(shared.BusinessDocumentID, CSPBusinessDocument) (CSPBusinessDocument, error)
	CSPBusinessDocumentExist(shared.BusinessDocumentID) bool
	CSPDocumentByDocumentID(shared.BusinessDocumentID) (CSPBusinessDocument, error)
}

type cspBusinessService struct {
	rdb *sqlx.DB
	wdb *sqlx.DB
}

// NewCSPBusinessService ..
func NewCSPBusinessService() CSPBusinessService {
	return cspBusinessService{wdb: data.DBWrite, rdb: data.DBRead}
}

func (s cspBusinessService) CSPBusinessDocumentCreate(create CSPBusinessDocument) error {
	//Check if we have a record
	var exist bool
	err := s.rdb.Get(&exist, "SELECT EXISTS (SELECT true FROM business_document WHERE document_id = $1)", create.DocumentID)

	if exist {
		log.Printf("record already exist, skipping creating csp business doc for id %v", create.DocumentID)
		return nil
	}

	_, err = s.wdb.Exec(`INSERT INTO business_document(document_id,document_status) VALUES($1,$2)`, create.DocumentID, create.Status)
	if err != nil {
		log.Printf("Error creating business document on CSP db %v", err)
	}

	return err
}

func (s cspBusinessService) CSPBusinessDocumentUpdate(id shared.BusinessDocumentID, update CSPBusinessDocument) (CSPBusinessDocument, error) {
	keys := services.SQLGenForUpdate(update)
	var d CSPBusinessDocument
	q := fmt.Sprintf("UPDATE business_document SET %s WHERE document_id = '%s' RETURNING *", keys, id)
	stmt, err := s.wdb.PrepareNamed(q)
	if err != nil {
		log.Printf("error preparing stmt %v query %v", err, q)
		return d, err
	}

	err = stmt.Get(&d, update)
	if err != nil {
		return d, fmt.Errorf("error keys: %v err: %v", keys, err)
	}

	return d, nil
}

func (s cspBusinessService) CSPBusinessDocumentExist(id shared.BusinessDocumentID) bool {
	var doc CSPBusinessDocument
	err := s.rdb.Get(&doc, "SELECT * FROM business_document WHERE document_id = $1", id)
	if err != nil {
		return false
	}

	if *doc.Status == csp.DocumentStatusUploaded {
		return true
	}

	return false
}

func (s cspBusinessService) CSPDocumentByDocumentID(id shared.BusinessDocumentID) (CSPBusinessDocument, error) {
	var doc CSPBusinessDocument
	err := s.rdb.Get(&doc, "SELECT * FROM business_document WHERE document_id = $1", id)
	if err != nil && err == sql.ErrNoRows {
		return doc, services.ErrorNotFound{}.New("")
	}

	return doc, err
}
