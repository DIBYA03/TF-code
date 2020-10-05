/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

// Package for transaction services
package transaction

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/wiseco/core-platform/auth"
	"github.com/wiseco/core-platform/services"
	"github.com/wiseco/core-platform/services/document"
	"github.com/wiseco/core-platform/shared"
)

type receiptDatastore struct {
	sourceReq services.SourceRequest
	*sqlx.DB
}

type ReceiptService interface {
	// Read
	GetById(string, shared.PostedTransactionID, shared.BusinessID) (*Receipt, error)
	List(int, int, shared.PostedTransactionID, shared.BusinessID) ([]Receipt, error)

	// Create receipt
	Create(*ReceiptCreate) (*Receipt, error)

	// Signed url
	GetSignedURL(string, shared.PostedTransactionID, shared.BusinessID) (*string, error)

	Delete(receiptID string, transactionID shared.PostedTransactionID, bID shared.BusinessID) error
}

func NewReceiptService(r services.SourceRequest) ReceiptService {
	return &receiptDatastore{r, DBWrite}
}

func (db *receiptDatastore) GetById(ID string, txnID shared.PostedTransactionID, bID shared.BusinessID) (*Receipt, error) {
	// Check access
	err := auth.NewAuthService(db.sourceReq).CheckBusinessAccess(bID)
	if err != nil {
		return nil, err
	}

	receipt := Receipt{}

	err = db.Get(&receipt, "SELECT * FROM business_transaction_attachment WHERE id = $1 AND transaction_id = $2 AND business_id = $3", ID, txnID, bID)
	if err != nil && err != sql.ErrNoRows {
		log.Println(err)
		return nil, err
	}

	return &receipt, err
}

func (db *receiptDatastore) List(offset int, limit int, txnID shared.PostedTransactionID, businessID shared.BusinessID) ([]Receipt, error) {
	// Check access
	err := auth.NewAuthService(db.sourceReq).CheckBusinessAccess(businessID)
	if err != nil {
		return nil, err
	}

	rows := []Receipt{}
	err = db.Select(&rows, "SELECT * FROM business_transaction_attachment WHERE business_id = $1 AND transaction_id = $2", businessID, txnID)
	switch {
	case err == sql.ErrNoRows:
		return rows, nil
	case err != nil:
		return nil, errors.Cause(err)
	}

	return rows, err
}

func (db *receiptDatastore) Create(r *ReceiptCreate) (*Receipt, error) {
	// Check access
	err := auth.NewAuthService(db.sourceReq).CheckBusinessAccess(r.BusinessID)
	if err != nil {
		return nil, err
	}

	storer, err := document.NewAWSS3DocStorage(string(r.BusinessID), document.BusinessPrefix)
	if err != nil {
		log.Printf("document upload failed error:%v", err)
		return nil, err
	}

	key, _ := storer.Key()
	r.StorageKey = key

	// Store in database
	columns := []string{
		"transaction_id", "created_user_id", "business_id", "content_type", "storage_key",
	}
	// Default/mandatory values
	values := []string{
		":transaction_id", ":created_user_id", ":business_id", ":content_type", ":storage_key",
	}

	sql := fmt.Sprintf("INSERT INTO business_transaction_attachment(%s) VALUES(%s) RETURNING *", strings.Join(columns, ", "), strings.Join(values, ", "))

	stmt, err := db.PrepareNamed(sql)
	if err != nil {
		return nil, err
	}

	receipt := &Receipt{}

	err = stmt.Get(receipt, &r)
	if err != nil {
		return nil, err
	}

	url, err := storer.PutSignedURL()
	if err != nil {
		log.Printf("Error getting pre-signed url:%v", err)
		return nil, err
	}

	receipt.SignedURL = url

	// Give back receipt object
	return receipt, nil

}

func (db *receiptDatastore) GetSignedURL(receiptID string, txnID shared.PostedTransactionID, businessID shared.BusinessID) (*string, error) {
	// Check access
	err := auth.NewAuthService(db.sourceReq).CheckBusinessAccess(businessID)
	if err != nil {
		return nil, err
	}

	var key string
	notFound := services.ErrorNotFound{}.New("")
	err = db.Get(&key, "SELECT storage_key from business_transaction_attachment WHERE id = $1 AND transaction_id = $2 AND business_id = $3", receiptID, txnID, businessID)
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

func (db *receiptDatastore) Delete(receiptID string, transactionID shared.PostedTransactionID, bID shared.BusinessID) error {
	_, err := db.Exec("UPDATE business_transaction_attachment SET deleted = CURRENT_TIMESTAMP WHERE id = $1 AND transaction_id = $2", receiptID, transactionID)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
