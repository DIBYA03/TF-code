package document

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/wiseco/core-platform/services"
	"github.com/wiseco/core-platform/services/csp"
	"github.com/wiseco/core-platform/services/csp/data"
	"github.com/wiseco/core-platform/shared"
)

// CSPConsumerService ..
type CSPConsumerService interface {
	CSPConsumerDocumentCreate(CSPConsumerDocument)
	CSPConsumerDocumentUpdate(shared.ConsumerDocumentID, CSPConsumerDocument) (CSPConsumerDocument, error)
	CSPConsumerDocumentExist(id shared.ConsumerDocumentID) bool
}

type cspConsumerService struct {
	rdb *sqlx.DB
	wdb *sqlx.DB
}

// NewCSPConsumerService ..
func NewCSPConsumerService() CSPConsumerService {
	return cspConsumerService{wdb: data.DBWrite, rdb: data.DBRead}
}

func (s cspConsumerService) CSPConsumerDocumentCreate(create CSPConsumerDocument) {
	//Check if we have a record
	var exist bool
	err := s.rdb.Get(&exist, "SELECT EXISTS (SELECT true FROM consumer_document WHERE document_id = $1)", create.DocumentID)

	if exist {
		log.Printf("record already exist, skipping creating csp consumer doc for id %v", create.DocumentID)
		return
	}

	_, err = s.wdb.Exec(`
	INSERT INTO consumer_document(document_id,document_status) VALUES($1,$2)`, create.DocumentID, create.Status)
	if err != nil {
		log.Printf("Error creating consumer document on CSP db %v", err)
	}
}

func (s cspConsumerService) CSPConsumerDocumentUpdate(id shared.ConsumerDocumentID, update CSPConsumerDocument) (CSPConsumerDocument, error) {
	keys := services.SQLGenForUpdate(update)
	var d CSPConsumerDocument
	q := fmt.Sprintf("UPDATE consumer_document SET %s WHERE document_id = '%s' RETURNING *", keys, id)
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

func (s cspConsumerService) CSPConsumerDocumentExist(id shared.ConsumerDocumentID) bool {
	var doc CSPConsumerDocument
	err := s.rdb.Get(&doc, "SELECT * FROM consumer_document WHERE document_id = $1", id)
	if err != nil {
		return false
	}

	if *doc.Status == csp.DocumentStatusUploaded {
		return true
	}

	return false
}
