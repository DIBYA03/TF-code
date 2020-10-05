/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

// Package for business
package business

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/wiseco/core-platform/auth"
	partnerbank "github.com/wiseco/core-platform/partner/bank"
	"github.com/wiseco/core-platform/services"
	"github.com/wiseco/core-platform/services/address"
	"github.com/wiseco/core-platform/services/contact"
	"github.com/wiseco/core-platform/services/data"
	"github.com/wiseco/core-platform/shared"
)

type linkedPayeeDataStore struct {
	sourceReq services.SourceRequest
	*sqlx.DB
}

type LinkedPayeeService interface {
	// Read
	GetByID(string, shared.ContactID, shared.BusinessID) (*LinkedPayee, error)
	GetByAddressID(shared.AddressID, shared.BusinessID) (*LinkedPayee, error)
	List(offset int, limit int, contactId string, businessID shared.BusinessID) ([]*LinkedPayee, error)

	// Create
	Create(LinkedPayeeCreate) (*LinkedPayee, error)

	// Deactivate payee
	Deactivate(string, string, shared.BusinessID) error
	DeactivateAll(string, shared.BusinessID) error

	//DEBUG
	DEBUGDeleteAllForBusiness(shared.BusinessID) error
}

func NewLinkedPayeeService(r services.SourceRequest) LinkedPayeeService {
	return &linkedPayeeDataStore{r, data.DBWrite}
}

func (db *linkedPayeeDataStore) GetByID(id string, contactID shared.ContactID, businessID shared.BusinessID) (*LinkedPayee, error) {
	// Check access
	err := auth.NewAuthService(db.sourceReq).CheckBusinessAccess(businessID)
	if err != nil {
		return nil, err
	}

	if os.Getenv("USE_BANKING_SERVICE") == "true" {
		lps, err := NewBankingLinkedPayeeService()
		if err != nil {
			return nil, err
		}

		return lps.GetByID(id)
	}

	a := LinkedPayee{}

	err = db.Get(&a, "SELECT * FROM business_linked_payee WHERE id = $1 AND contact_id = $2 AND business_id = $3 AND status = $4", id, contactID, businessID, PayeeStatusActive)

	return &a, err
}

func (db *linkedPayeeDataStore) GetByAddressID(addressID shared.AddressID, businessID shared.BusinessID) (*LinkedPayee, error) {
	// Check access
	err := auth.NewAuthService(db.sourceReq).CheckBusinessAccess(businessID)
	if err != nil {
		return nil, err
	}

	if os.Getenv("USE_BANKING_SERVICE") == "true" {
		lps, err := NewBankingLinkedPayeeService()
		if err != nil {
			return nil, err
		}

		return lps.GetByAddressID(addressID)
	}

	a := LinkedPayee{}

	err = db.Get(&a, "SELECT * FROM business_linked_payee WHERE address_id = $1 AND status = $2", addressID, PayeeStatusActive)

	return &a, err
}

func (db *linkedPayeeDataStore) List(offset int, limit int, contactID string, businessID shared.BusinessID) ([]*LinkedPayee, error) {
	// Check access
	err := auth.NewAuthService(db.sourceReq).CheckBusinessAccess(businessID)
	if err != nil {
		return nil, err
	}

	if os.Getenv("USE_BANKING_SERVICE") == "true" {
		lps, err := NewBankingLinkedPayeeService()
		if err != nil {
			return nil, err
		}

		return lps.List(contactID, businessID, limit, offset)
	}

	rows := []*LinkedPayee{}

	err = db.Select(&rows, `SELECT * FROM business_linked_payee WHERE contact_id = $1 AND business_id = $2 AND status = $3`, contactID, businessID, PayeeStatusActive)
	switch {
	case err == sql.ErrNoRows:
		return rows, nil
	case err != nil:
		return nil, errors.Cause(err)
	}

	return rows, err
}

func (db *linkedPayeeDataStore) Create(payeeCreate LinkedPayeeCreate) (*LinkedPayee, error) {
	// Check access
	err := auth.NewAuthService(db.sourceReq).CheckBusinessAccess(payeeCreate.BusinessID)

	if err != nil {
		return nil, err
	}

	if payeeCreate.ContactID == "" {
		return nil, errors.New("Contact ID is required")
	}

	if payeeCreate.BusinessID == shared.BusinessID("") {
		return nil, errors.New("Business ID is required")
	}

	if payeeCreate.AddressID == shared.AddressID("") {
		return nil, errors.New("Address ID is required")
	}

	a, err := address.NewAddressService(db.sourceReq).GetByID(payeeCreate.AddressID)

	if err != nil {
		return nil, err
	}

	if a == nil {
		return nil, fmt.Errorf("Could not find address for id:%s", payeeCreate.AddressID)
	}

	c, err := contact.NewContactService(db.sourceReq).GetById(payeeCreate.ContactID, payeeCreate.BusinessID)

	if err != nil {
		return nil, err
	}

	//If none was passed in let's get the name out of the Contact
	if payeeCreate.AccountHolderName == "" {
		payeeCreate.AccountHolderName = c.Name()
	}

	if payeeCreate.PayeeName == "" {
		payeeCreate.PayeeName = c.Name()
	}

	bank, err := partnerbank.GetBusinessBank(partnerbank.ProviderNameBBVA)
	if err != nil {
		return nil, err
	}

	payeeCreate.BankName = partnerbank.ProviderNameBBVA

	if os.Getenv("USE_BANKING_SERVICE") == "true" {
		lps, err := NewBankingLinkedPayeeService()
		if err != nil {
			return nil, err
		}

		return lps.Create(payeeCreate, a)
	}

	srv, err := bank.LinkedPayeeService(db.sourceReq.PartnerBankRequest(), partnerbank.BusinessID(payeeCreate.BusinessID))

	adr := &partnerbank.PayeeAddressRequest{
		Line1:   a.StreetAddress,
		City:    a.Locality,
		State:   a.AdminArea,
		ZipCode: a.PostalCode,
	}

	if a.AddressLine2 != "" {
		adr.Line2 = a.AddressLine2
	}

	resp, err := srv.Link(&partnerbank.LinkedPayeeRequest{
		PayeeName:         payeeCreate.PayeeName,
		AccountNumber:     payeeCreate.ContactID,
		NameOnAccount:     payeeCreate.AccountHolderName,
		RemittanceAddress: adr,
	})

	if err != nil {
		return nil, err
	}

	payeeCreate.BankPayeeID = string(resp.BankPayeeID)

	payeeCreate.Status = PayeeStatusActive

	// Default/mandatory fields
	columns := []string{
		"business_id", "contact_id", "address_id", "bank_payee_id", "bank_name", "account_holder_name", "payee_name", "status",
	}
	// Default/mandatory values
	values := []string{
		":business_id", ":contact_id", ":address_id", ":bank_payee_id", ":bank_name", ":account_holder_name", ":payee_name", ":status",
	}

	sql := fmt.Sprintf("INSERT INTO business_linked_payee(%s) VALUES(%s) RETURNING *", strings.Join(columns, ", "), strings.Join(values, ", "))

	stmt, err := db.PrepareNamed(sql)
	if err != nil {
		return nil, err
	}

	l := &LinkedPayee{}

	err = stmt.Get(l, &payeeCreate)

	return l, err
}

func (db *linkedPayeeDataStore) DeactivateAll(contactID string, businessID shared.BusinessID) error {

	la, err := db.List(0, 10, contactID, businessID)
	if err != nil && err != sql.ErrNoRows {
		log.Println(err)
		return err
	}

	if err == sql.ErrNoRows {
		log.Println(err)
		return nil
	}

	for _, acc := range la {
		err = db.Deactivate(acc.ID, contactID, businessID)
		if err != nil {
			log.Println(err)
			return err
		}
	}

	return nil

}

func (db *linkedPayeeDataStore) Deactivate(ID string, contactID string, businessID shared.BusinessID) error {
	// Check access
	err := auth.NewAuthService(db.sourceReq).CheckBusinessAccess(businessID)
	if err != nil {
		log.Println(err)
		return err
	}

	// Deactivate Payee
	_, err = db.Exec("UPDATE business_linked_payee SET status = $1 WHERE id = $2 AND contact_id = $3 AND business_id = $4", PayeeStatusInactive, ID, contactID, businessID)
	return err
}

//ONLY for debug
func (db *linkedPayeeDataStore) DEBUGDeleteAllForBusiness(businessID shared.BusinessID) error {
	_, err := db.Exec("DELETE FROM business_linked_payee where business_id = $1", businessID)

	return err
}
