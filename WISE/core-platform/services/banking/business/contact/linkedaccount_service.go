/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

// Package for business contacts
package contact

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
	"github.com/wiseco/core-platform/services/banking"
	"github.com/wiseco/core-platform/services/banking/business"
	"github.com/wiseco/core-platform/services/contact"
	"github.com/wiseco/core-platform/services/data"
	"github.com/wiseco/core-platform/shared"
	grpcBanking "github.com/wiseco/protobuf/golang/banking"
)

type linkedAccountDatastore struct {
	sourceReq services.SourceRequest
	*sqlx.DB
}

type LinkedAccountService interface {
	// Read
	GetById(string, string, shared.BusinessID) (*business.LinkedBankAccount, error)
	List(offset int, limit int, contactID string, businessID shared.BusinessID) ([]*business.LinkedBankAccount, error)

	// Create/Modify
	Create(*business.ContactLinkedAccountCreate) (*business.LinkedBankAccount, error)
	UpdateLinkedAccountUsageType(u *business.LinkedAccountUpdate) error

	// Deactivate
	Deactivate(string, string, shared.BusinessID) error
	DeactivateAll(string, shared.BusinessID) error
}

func NewLinkedAccountService(r services.SourceRequest) LinkedAccountService {
	return &linkedAccountDatastore{r, data.DBWrite}
}

func (db *linkedAccountDatastore) List(offset int, limit int, contactId string, businessID shared.BusinessID) ([]*business.LinkedBankAccount, error) {
	// Check access
	err := auth.NewAuthService(db.sourceReq).CheckBusinessAccess(businessID)
	if err != nil {
		return nil, err
	}

	rows := []*business.LinkedBankAccount{}

	if os.Getenv("USE_BANKING_SERVICE") == "true" {
		las, err := business.NewBankingLinkedAccountService()
		if err != nil {
			return rows, err
		}

		cID := shared.ContactID(contactId)

		stfs := []grpcBanking.LinkedSubtype{
			grpcBanking.LinkedSubtype_LST_CONTACT,
		}

		as, err := las.ListWithContact(businessID, cID, stfs, limit, offset)
		if err != nil {
			return rows, err
		}

		for _, a := range as {
			rows = append(rows, &a)
		}

		return rows, nil
	}

	err = db.Select(&rows, `SELECT * FROM business_linked_bank_account WHERE contact_id = $1 AND business_id = $2 AND 
	usage_type is distinct from $3 AND deactivated IS NULL`, contactId, businessID, business.UsageTypeContactInvisible)

	switch {
	case err == sql.ErrNoRows:
		return rows, nil
	case err != nil:
		return nil, errors.Cause(err)
	}

	return rows, err
}

func (db *linkedAccountDatastore) GetById(id string, contactId string, businessID shared.BusinessID) (*business.LinkedBankAccount, error) {
	// Check access
	err := auth.NewAuthService(db.sourceReq).CheckBusinessAccess(businessID)
	if err != nil {
		return nil, err
	}

	if os.Getenv("USE_BANKING_SERVICE") == "true" {
		las, err := business.NewBankingLinkedAccountService()
		if err != nil {
			return nil, err
		}

		return las.GetById(id)
	}

	a := business.LinkedBankAccount{}

	err = db.Get(&a, "SELECT * FROM business_linked_bank_account WHERE id = $1 AND contact_id = $2 AND business_id = $3", id, contactId, businessID)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return &a, err
}

func (db *linkedAccountDatastore) GetByAccountId(businessID shared.BusinessID, accountNumber business.AccountNumber, routingNumber string) (*business.LinkedBankAccount, error) {
	// Check access
	err := auth.NewAuthService(db.sourceReq).CheckBusinessAccess(businessID)
	if err != nil {
		return nil, err
	}

	if os.Getenv("USE_BANKING_SERVICE") == "true" {
		las, err := business.NewBankingLinkedAccountService()
		if err != nil {
			return nil, err
		}

		return las.GetByAccountNumber(businessID, accountNumber.String(), routingNumber)
	}

	a := business.LinkedBankAccount{}

	err = db.Get(&a, "SELECT * FROM business_linked_bank_account WHERE business_id = $1 AND account_number = $2 AND routing_number = $3 AND deactivated IS NULL", businessID, accountNumber, routingNumber)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return &a, err
}

func (db *linkedAccountDatastore) Create(accountCreate *business.ContactLinkedAccountCreate) (*business.LinkedBankAccount, error) {
	// Check access
	err := auth.NewAuthService(db.sourceReq).CheckBusinessAccess(accountCreate.BusinessID)
	if err != nil {
		return nil, err
	}

	if accountCreate.ContactId == "" {
		return nil, errors.New("contact id required")
	}

	if accountCreate.AccountType == "" {
		return nil, errors.New("account type is required")
	}

	if accountCreate.AccountNumber == "" {
		return nil, errors.New("account number is required")
	}

	if accountCreate.RoutingNumber == "" {
		return nil, errors.New("routing number is required")
	}

	if accountCreate.Currency == "" {
		return nil, errors.New("currency is required")
	}

	if accountCreate.Permission == "" {
		return nil, errors.New("account permission is required")
	}

	account, err := db.GetByAccountId(accountCreate.BusinessID, accountCreate.AccountNumber, accountCreate.RoutingNumber)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if account.RegisteredAccountId != "" {
		return account, errors.New("account already registered")
	}

	c, err := contact.NewContactService(db.sourceReq).GetById(accountCreate.ContactId, accountCreate.BusinessID)
	if err != nil {
		return nil, errors.New("contact not found")
	}

	if os.Getenv("USE_BANKING_SERVICE") == "true" {
		las, err := business.NewBankingLinkedAccountService()
		if err != nil {
			return nil, err
		}

		q := `SELECT consumer_id FROM wise_user WHERE id = $1`

		var cID shared.ConsumerID

		err = db.Get(&cID, q, accountCreate.UserID)
		if err != nil {
			return nil, err
		}

		ut := business.UsageTypeContact
		if accountCreate.UsageType != nil {
			ut = *accountCreate.UsageType
		}

		linkedAccount := &business.LinkedBankAccount{
			BusinessID: accountCreate.BusinessID,
			//BusinessBankAccountId: &lbap.AccountId,
			ContactId:           &accountCreate.ContactId,
			RegisteredAccountId: accountCreate.RegisteredAccountId,
			RegisteredBankName:  accountCreate.RegisteredBankName,
			AccountHolderName:   c.Name(),
			AccountNumber:       accountCreate.AccountNumber,
			Currency:            banking.CurrencyUSD,
			AccountType:         accountCreate.AccountType,
			UsageType:           &ut,
			RoutingNumber:       accountCreate.RoutingNumber,
			Permission:          accountCreate.Permission,
		}

		return las.LinkBankAccount(linkedAccount, cID)
	}

	request := partnerbank.LinkedBankAccountRequest{
		AccountHolderName: c.Name(),
		AccountNumber:     string(accountCreate.AccountNumber),
		AccountType:       partnerbank.AccountType(accountCreate.AccountType),
		Currency:          partnerbank.Currency(accountCreate.Currency),
		Permission:        partnerbank.LinkedAccountPermission(accountCreate.Permission),
		RoutingNumber:     accountCreate.RoutingNumber,
	}

	// Register account with partner bank
	bank, err := partnerbank.GetBusinessBank(partnerbank.ProviderNameBBVA)
	if err != nil {
		return nil, err
	}

	srv, err := bank.LinkedAccountService(db.sourceReq.PartnerBankRequest(), partnerbank.BusinessID(accountCreate.BusinessID))
	if err != nil {
		return nil, err
	}

	resp, err := srv.Link(&request)
	if err != nil {
		return nil, err
	}

	linkedAccount := transformResponse(resp)
	linkedAccount.UserID = accountCreate.UserID
	linkedAccount.BusinessID = accountCreate.BusinessID
	linkedAccount.ContactId = &accountCreate.ContactId
	linkedAccount.UsageType = accountCreate.UsageType

	// Set usage type to contact - limited capabilities
	if linkedAccount.UsageType == nil {
		ut := business.UsageTypeContact
		linkedAccount.UsageType = &ut
	}

	// Default/mandatory fields
	columns := []string{
		"user_id", "business_id", "contact_id", "registered_account_id", "registered_bank_name",
		"account_holder_name", "currency", "usage_type", "account_type", "account_number",
		"routing_number", "account_permission",
	}
	// Default/mandatory values
	values := []string{
		":user_id", ":business_id", ":contact_id", ":registered_account_id", ":registered_bank_name",
		":account_holder_name", ":currency", ":usage_type", ":account_type", ":account_number",
		":routing_number", ":account_permission",
	}

	sql := fmt.Sprintf("INSERT INTO business_linked_bank_account(%s) VALUES(%s) RETURNING *", strings.Join(columns, ", "), strings.Join(values, ", "))

	stmt, err := db.PrepareNamed(sql)
	if err != nil {
		return nil, err
	}

	l := &business.LinkedBankAccount{}

	err = stmt.Get(l, &linkedAccount)
	if err != nil {
		return nil, err
	}

	return l, nil

}

// Transform partner bank layer response
func transformResponse(response *partnerbank.LinkedBankAccountResponse) business.LinkedBankAccount {

	permission := banking.LinkedAccountPermission(response.Permission)
	currency := banking.Currency(response.Currency)

	linkedAccount := business.LinkedBankAccount{
		AccountNumber:       business.AccountNumber(response.AccountNumber),
		RoutingNumber:       response.RoutingNumber,
		WireRouting:         response.WireRouting,
		BankName:            &response.AccountBankName,
		AccountType:         banking.AccountType(response.AccountType),
		AccountHolderName:   response.AccountHolderName,
		Alias:               &response.Alias,
		Permission:          permission,
		Currency:            currency,
		RegisteredAccountId: string(response.AccountID),
		RegisteredBankName:  "bbva",
	}

	return linkedAccount
}

func (db *linkedAccountDatastore) DeactivateAll(contactID string, businessID shared.BusinessID) error {

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
		err = db.Deactivate(acc.Id, contactID, businessID)
		if err != nil {
			log.Println(err)
			return err
		}
	}

	return nil

}

func (db *linkedAccountDatastore) Deactivate(ID string, contactID string, businessID shared.BusinessID) error {
	// Check access
	err := auth.NewAuthService(db.sourceReq).CheckBusinessAccess(businessID)
	if err != nil {
		return err
	}

	// Get linked account
	la, err := db.GetById(ID, contactID, businessID)
	if err != nil {
		return err
	}

	if os.Getenv("USE_BANKING_SERVICE") == "true" {
		las, err := business.NewBankingLinkedAccountService()
		if err != nil {
			return err
		}

		_, err = las.UnlinkBankAccount(ID)

		return err
	}

	// Unlink card
	bank, err := partnerbank.GetBusinessBank(partnerbank.ProviderNameBBVA)
	if err != nil {
		return err
	}

	srv, err := bank.LinkedAccountService(db.sourceReq.PartnerBankRequest(), partnerbank.BusinessID(businessID))
	if err != nil {
		return nil
	}

	err = srv.Unlink(partnerbank.LinkedAccountBankID(la.RegisteredAccountId))
	if err != nil {
		return err
	}

	// Deactivate card
	_, err = db.Exec("UPDATE business_linked_bank_account SET deactivated = CURRENT_TIMESTAMP WHERE id = $1 AND contact_id = $2 AND business_id = $3", ID, contactID, businessID)
	return err
}

func (db *linkedAccountDatastore) UpdateLinkedAccountUsageType(u *business.LinkedAccountUpdate) error {
	var columns []string

	if os.Getenv("USE_BANKING_SERVICE") == "true" {
		las, err := business.NewBankingLinkedAccountService()
		if err != nil {
			return err
		}

		_, err = las.Update(u)

		return err
	}

	if u.UsageType != nil {
		columns = append(columns, "usage_type = :usage_type")
	}

	_, err := db.NamedExec(fmt.Sprintf("UPDATE business_linked_bank_account SET %s WHERE id = '%s'", strings.Join(columns, ", "), u.ID), u)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
