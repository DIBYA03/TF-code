package payment

import (
	"database/sql"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/wiseco/core-platform/auth"
	"github.com/wiseco/core-platform/services"
	"github.com/wiseco/core-platform/services/data"
	"github.com/wiseco/core-platform/services/document"
	"github.com/wiseco/core-platform/shared"
)

type invoiceDatastore struct {
	sourceReq services.SourceRequest
	*sqlx.DB
}

type InvoiceService interface {
	GetSignedURL(invoiceID string, businessID shared.BusinessID) (*string, error)
	GetSignedURLByRequestID(invoiceID string, businessID shared.BusinessID) (*string, error)
}

func NewInvoiceService(r services.SourceRequest) InvoiceService {
	return &invoiceDatastore{r, data.DBWrite}
}

func (db *invoiceDatastore) GetSignedURL(invoiceID string, businessID shared.BusinessID) (*string, error) {
	// Check access
	err := auth.NewAuthService(db.sourceReq).CheckBusinessAccess(businessID)
	if err != nil {
		return nil, err
	}

	var key string
	notFound := services.ErrorNotFound{}.New("")
	err = db.Get(&key, "SELECT storage_key from business_invoice WHERE id = $1 AND business_id = $2", invoiceID, businessID)
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

	storer, err := document.NewStorerFromKey(key)

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

func (db *invoiceDatastore) GetSignedURLByRequestID(requestID string, businessID shared.BusinessID) (*string, error) {
	// Check access
	err := auth.NewAuthService(db.sourceReq).CheckBusinessAccess(businessID)
	if err != nil {
		return nil, err
	}

	var key string
	notFound := services.ErrorNotFound{}.New("")
	err = db.Get(&key, "SELECT storage_key from business_invoice WHERE request_id = $1 AND business_id = $2", requestID, businessID)
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

	storer, err := document.NewStorerFromKey(key)

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
