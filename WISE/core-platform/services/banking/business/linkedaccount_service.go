/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

// Package for all internal banking related items such as accounts and cards
package business

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/types"
	"github.com/pkg/errors"
	"github.com/wiseco/core-platform/auth"
	partnerbank "github.com/wiseco/core-platform/partner/bank"
	"github.com/wiseco/core-platform/services"
	"github.com/wiseco/core-platform/services/banking"
	bsrv "github.com/wiseco/core-platform/services/business"
	consrv "github.com/wiseco/core-platform/services/contact"
	"github.com/wiseco/core-platform/services/data"
	usersrv "github.com/wiseco/core-platform/services/user"
	"github.com/wiseco/core-platform/shared"
	"github.com/wiseco/go-lib/id"
	l "github.com/wiseco/go-lib/log"
	"github.com/wiseco/go-lib/plaid"
	grpcBanking "github.com/wiseco/protobuf/golang/banking"
)

type linkedAccountDataStore struct {
	sourceReq services.SourceRequest
	*sqlx.DB
}

type LinkedAccountSource string

const (
	LinkedAccountSourcePlaid = LinkedAccountSource("plaid")
)

type LinkedAccountService interface {
	// Read
	GetById(string, shared.BusinessID) (*LinkedBankAccount, error)
	GetByIDInternal(string) (*LinkedBankAccount, error)
	List(offset int, limit int, businessID shared.BusinessID) ([]LinkedBankAccount, error)
	GetByAccountIDInternal(businessID shared.BusinessID, accountID string) (*LinkedBankAccount, error)
	GetByAccountNumber(businessID shared.BusinessID, accountNumber AccountNumber, routingNumber string) (*LinkedBankAccount, error)
	GetByAccountNumberInternal(businessID shared.BusinessID, accountNumber AccountNumber, routingNumber string) (*LinkedBankAccount, error)

	// Connect and get business bank accounts
	ConnectBankAccount(linkedAccountRequest LinkedBankAccountRequest, businessID shared.BusinessID) ([]LinkedBankAccountBase, error)

	// Register bank account
	LinkExternalBankAccount(*LinkedExternalAccountCreate) (*LinkedBankAccount, error)
	LinkOwnBankAccount(*BankAccount) (*LinkedBankAccount, error)
	LinkToClearingAccount(*ClearingLinkedAccountCreate) (*LinkedBankAccount, error)
	LinkMerchantBankAccount(*MerchantLinkedAccountCreate) (*LinkedBankAccount, error)

	// Unlink bank account
	UnlinkBankAccount(string, shared.BusinessID) (*LinkedBankAccount, error)
}

func NewLinkedAccountService(r services.SourceRequest) LinkedAccountService {
	return &linkedAccountDataStore{r, data.DBWrite}
}

func NewLinkedAccountServiceWithout() LinkedAccountService {
	return &linkedAccountDataStore{services.SourceRequest{}, data.DBWrite}
}

func (db *linkedAccountDataStore) ConnectBankAccount(linkedBankAccountRequest LinkedBankAccountRequest, businessID shared.BusinessID) ([]LinkedBankAccountBase, error) {
	// Check access
	err := auth.NewAuthService(db.sourceReq).CheckBusinessAccess(businessID)
	if err != nil {
		return nil, err
	}

	if linkedBankAccountRequest.PublicToken == "" {
		return nil, errors.New("Public token is required")
	}

	response, err := plaid.NewPlaidService(l.NewLogger()).GetAccounts(linkedBankAccountRequest.PublicToken)
	if err != nil {
		return nil, err
	}

	// Transform plaid response to LinkedAccount that can be sent back to client
	linkedBankAccounts, err := transformPlaidResponseToLinkedAccountBase(response)
	if err != nil {
		return nil, err
	}

	return linkedBankAccounts, nil
}

func transformPlaidResponseToLinkedAccountBase(plaidResponse *plaid.PlaidResponse) ([]LinkedBankAccountBase, error) {

	linkedBankAccounts := make([]LinkedBankAccountBase, 0)

	for _, achNumber := range plaidResponse.Numbers.ACH {

		account, err := getAccountByID(plaidResponse.Account, achNumber.AccountId)
		if err != nil {
			return nil, err
		}

		//-- Only depository checking and depository savings accounts are sent back to client
		if account.Type == PlaidAccountTypeDepository && (account.SubType == banking.BankAccountTypeChecking || account.SubType == banking.BankAccountTypeSavings) {

			linkedBankAccount := LinkedBankAccountBase{
				SourceAccountId: &account.AccountId,
				AccountName:     account.Name,
				AccountNumber:   AccountNumber(achNumber.Account),
			}

			linkedBankAccounts = append(linkedBankAccounts, linkedBankAccount)
		}
	}

	return linkedBankAccounts, nil

}

func (db *linkedAccountDataStore) LinkExternalBankAccount(accountCreate *LinkedExternalAccountCreate) (*LinkedBankAccount, error) {

	if accountCreate.PublicToken == "" {
		return nil, errors.New("Public token is required")
	}

	if accountCreate.SourceAccountId == "" {
		return nil, errors.New("Account id is required")
	}

	accountResp, err := plaid.NewPlaidService(l.NewLogger()).GetAccounts(accountCreate.PublicToken)
	if err != nil {
		return nil, err
	}

	identityResp, err := plaid.NewPlaidService(l.NewLogger()).GetIdentity(accountCreate.PublicToken)
	if err != nil {
		return nil, err
	}

	identityResp.Numbers.ACH = accountResp.Numbers.ACH

	err = db.addExternalAccountsAndOwners(accountCreate.BusinessID, identityResp)
	if err != nil {
		return nil, err
	}

	verificationRequest := ExternalAccountVerificationRequest{
		BusinessID:       accountCreate.BusinessID,
		AccessToken:      identityResp.AccessToken,
		PartnerItemID:    identityResp.ItemID,
		PartnerAccountID: accountCreate.SourceAccountId,
	}
	err = NewExternalAccountService(db.sourceReq).Verify(verificationRequest)
	if err != nil {
		log.Println("Verification error ", err)
		return nil, err
	}

	// Transform plaid response to LinkedAccount that can be sent back to client
	linkedAccount, err := transformPlaidResponseToLinkedAccount(
		accountResp,
		accountCreate.UserID,
		accountCreate.BusinessID,
		accountCreate.SourceAccountId,
		accountCreate.ContactID,
	)
	if err != nil {
		return nil, err
	}

	// Check if user has already registered that account with us
	userBankAccount, err := db.GetByAccountNumber(linkedAccount.BusinessID, linkedAccount.AccountNumber, linkedAccount.RoutingNumber)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if userBankAccount.RegisteredAccountId != "" && userBankAccount.Deactivated == nil {
		return userBankAccount, errors.New("Account already registered")
	}

	var accountHolderName string
	if accountCreate.ContactID != nil {
		// Get contact name
		c, err := consrv.NewContactService(db.sourceReq).GetById(*accountCreate.ContactID, accountCreate.BusinessID)
		if err != nil {
			return nil, err
		}

		accountHolderName = c.Name()
	} else {
		//- Get user first name and last name
		u, err := usersrv.NewUserService(db.sourceReq).GetById(accountCreate.UserID)
		if err != nil {
			return nil, err
		}

		accountHolderName = u.Name()
	}

	if os.Getenv("USE_BANKING_SERVICE") == "true" {
		las, err := NewBankingLinkedAccountService()
		if err != nil {
			return nil, err
		}

		q := `SELECT consumer_id FROM wise_user WHERE id = $1`

		var cID shared.ConsumerID

		err = db.Get(&cID, q, linkedAccount.UserID)
		if err != nil {
			return nil, err
		}

		linkedAccount.AccountHolderName = accountHolderName

		la, err := las.LinkBankAccount(linkedAccount, cID)
		if err != nil {
			return nil, err
		}

		laID, err := id.ParseLinkedBankAccountID(la.Id)
		if err != nil {
			return nil, err
		}

		laIDStr := laID.UUIDString()

		// Update external bank account
		cu := ExternalBankAccountUpdate{
			BusinessID:      la.BusinessID,
			AccountNumber:   string(la.AccountNumber),
			RoutingNumber:   la.RoutingNumber,
			LinkedAccountID: &laIDStr,
		}

		_, err = NewExternalAccountService(db.sourceReq).Upsert(cu)
		if err != nil {
			return nil, err
		}

		return la, nil
	}

	request := partnerbank.LinkedBankAccountRequest{
		AccountHolderName: accountHolderName,
		AccountNumber:     string(linkedAccount.AccountNumber),
		AccountType:       partnerbank.AccountType(linkedAccount.AccountType),
		Currency:          partnerbank.Currency(linkedAccount.Currency),
		Permission:        partnerbank.LinkedAccountPermission(banking.LinkedAccountPermissionSendAndRecieve),
		RoutingNumber:     linkedAccount.RoutingNumber,
	}

	// Register bank account with partner bank
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

	linkedAccount = transformWithResponse(resp, linkedAccount)

	// Save registered bank account in database once its registered with BBVA
	la, err := db.CreateLinkedBankAccount(linkedAccount)
	if err != nil {
		return nil, err
	}

	// Update external bank account
	cu := ExternalBankAccountUpdate{
		BusinessID:      la.BusinessID,
		AccountNumber:   string(la.AccountNumber),
		RoutingNumber:   la.RoutingNumber,
		LinkedAccountID: &la.Id,
	}

	_, err = NewExternalAccountService(db.sourceReq).Upsert(cu)
	if err != nil {
		return nil, err
	}

	return la, nil
}

func (db *linkedAccountDataStore) LinkOwnBankAccount(bankAccount *BankAccount) (*LinkedBankAccount, error) {
	// Removing Check access since this endpoint will only be used on CSP
	// err := auth.NewAuthService(db.sourceReq).CheckBusinessBankAccountAccess(bankAccount.Id)
	// if err != nil {
	// 	log.Println("error getting checking auth access")
	// 	return nil, err
	// }

	// Get business
	b, err := bsrv.NewBusinessService(db.sourceReq).GetByIdInternal(bankAccount.BusinessID)
	if err != nil {
		log.Println("Error getting business for linked own bank account: ", err)
		return nil, err
	}

	linkedAccount := &LinkedBankAccount{
		UserID:                bankAccount.AccountHolderID,
		BusinessID:            bankAccount.BusinessID,
		BusinessBankAccountId: &bankAccount.Id,
		AccountHolderName:     b.Name(),
		AccountNumber:         AccountNumber(bankAccount.AccountNumber),
		Currency:              banking.Currency(bankAccount.Currency),
		AccountType:           banking.AccountType(bankAccount.AccountType),
		UsageType:             &bankAccount.UsageType,
		RoutingNumber:         bankAccount.RoutingNumber,
		WireRouting:           bankAccount.WireRouting,
		Permission:            banking.LinkedAccountPermissionSendAndRecieve,
	}

	if os.Getenv("USE_BANKING_SERVICE") == "true" {
		las, err := NewBankingLinkedAccountService()
		if err != nil {
			return nil, err
		}

		q := `SELECT consumer_id FROM wise_user WHERE id = $1`

		var cID shared.ConsumerID

		err = db.Get(&cID, q, linkedAccount.UserID)
		if err != nil {
			return nil, err
		}

		return las.LinkBankAccount(linkedAccount, cID)
	}

	request := partnerbank.LinkedBankAccountRequest{
		AccountHolderName: b.Name(),
		AccountNumber:     string(bankAccount.AccountNumber),
		AccountType:       partnerbank.AccountType(bankAccount.AccountType),
		Currency:          partnerbank.Currency(bankAccount.Currency),
		Permission:        partnerbank.LinkedAccountPermission(banking.LinkedAccountPermissionSendAndRecieve),
		RoutingNumber:     bankAccount.RoutingNumber,
	}

	// Register bank account with partner bank
	bank, err := partnerbank.GetBusinessBank(partnerbank.ProviderNameBBVA)
	if err != nil {
		return nil, err
	}

	srv, err := bank.LinkedAccountService(db.sourceReq.PartnerBankRequest(), partnerbank.BusinessID(bankAccount.BusinessID))
	if err != nil {
		return nil, err
	}

	resp, err := srv.Link(&request)
	if err != nil {
		log.Println(err)
		log.Println("error linking account services with partner")
		return nil, err
	}

	linkedAccount = transformWithResponse(resp, linkedAccount)

	// Save registered bank account in database once its registered with BBVA
	return db.CreateLinkedBankAccount(linkedAccount)
}

func (db *linkedAccountDataStore) LinkToClearingAccount(c *ClearingLinkedAccountCreate) (*LinkedBankAccount, error) {
	if c.AccountType == "" {
		return nil, errors.New("Account type is required")
	}

	if c.AccountNumber == "" {
		return nil, errors.New("Account number is required")
	}

	if c.RoutingNumber == "" {
		return nil, errors.New("Routing number is required")
	}

	if c.Currency == "" {
		return nil, errors.New("Currency is required")
	}

	if c.Permission == "" {
		return nil, errors.New("Account permission is required")
	}

	// Get business
	b, err := bsrv.NewBusinessService(db.sourceReq).GetById(c.BusinessID)
	if err != nil {
		log.Println("error getting business for linked own bank account")
		return nil, err
	}

	account, err := db.GetByAccountNumber(c.BusinessID, c.AccountNumber, c.RoutingNumber)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if account.RegisteredAccountId != "" {
		return nil, errors.New("Account already registered")
	}

	request := partnerbank.LinkedBankAccountRequest{
		AccountHolderName: b.Name(),
		AccountNumber:     string(c.AccountNumber),
		AccountType:       partnerbank.AccountType(c.AccountType),
		Currency:          partnerbank.Currency(c.Currency),
		Permission:        partnerbank.LinkedAccountPermission(c.Permission),
		RoutingNumber:     c.RoutingNumber,
	}

	// Register account with partner bank
	bank, err := partnerbank.GetBusinessBank(partnerbank.ProviderNameBBVA)
	if err != nil {
		return nil, err
	}

	srv, err := bank.LinkedAccountService(db.sourceReq.PartnerBankRequest(), partnerbank.BusinessID(c.BusinessID))
	if err != nil {
		return nil, err
	}

	resp, err := srv.Link(&request)
	if err != nil {
		return nil, err
	}

	linkedAccount := transformResponse(resp)
	linkedAccount.UserID = c.UserID
	linkedAccount.BusinessID = c.BusinessID
	linkedAccount.AccountHolderName = b.Name()

	// Set usage type to clearing
	ut := UsageTypeClearing
	linkedAccount.UsageType = &ut

	// Default/mandatory fields
	columns := []string{
		"user_id", "business_id", "business_bank_account_id", "registered_account_id",
		"registered_bank_name", "account_holder_name", "currency", "account_type",
		"usage_type", "account_number", "routing_number", "account_permission",
	}
	// Default/mandatory values
	values := []string{
		":user_id", ":business_id", ":business_bank_account_id", ":registered_account_id",
		":registered_bank_name", ":account_holder_name", ":currency", ":account_type",
		":usage_type", ":account_number", ":routing_number", ":account_permission",
	}

	sql := fmt.Sprintf("INSERT INTO business_linked_bank_account(%s) VALUES(%s) RETURNING *", strings.Join(columns, ", "), strings.Join(values, ", "))

	stmt, err := db.PrepareNamed(sql)
	if err != nil {
		return nil, err
	}

	l := &LinkedBankAccount{}

	err = stmt.Get(l, &linkedAccount)
	if err != nil {
		return nil, err
	}

	return l, nil
}

func (db *linkedAccountDataStore) LinkMerchantBankAccount(accountCreate *MerchantLinkedAccountCreate) (*LinkedBankAccount, error) {
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

	account, err := db.GetByAccountNumberInternal(accountCreate.BusinessID, accountCreate.AccountNumber, accountCreate.RoutingNumber)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if account.RegisteredAccountId != "" {
		return account, errors.New("account already registered")
	}

	if os.Getenv("USE_BANKING_SERVICE") == "true" {
		las, err := NewBankingLinkedAccountService()
		if err != nil {
			return nil, err
		}

		q := `SELECT consumer_id FROM wise_user WHERE id = $1`

		var cID shared.ConsumerID

		err = db.Get(&cID, q, accountCreate.UserID)
		if err != nil {
			return nil, err
		}

		linkedAccount := &LinkedBankAccount{
			UserID:            accountCreate.UserID,
			BusinessID:        accountCreate.BusinessID,
			AccountHolderName: accountCreate.AccountHolderName,
			AccountNumber:     accountCreate.AccountNumber,
			Currency:          accountCreate.Currency,
			AccountType:       accountCreate.AccountType,
			UsageType:         &accountCreate.UsageType,
			RoutingNumber:     accountCreate.RoutingNumber,
			Permission:        accountCreate.Permission,
		}

		return las.LinkBankAccount(linkedAccount, cID)
	}

	request := partnerbank.LinkedBankAccountRequest{
		AccountHolderName: accountCreate.AccountHolderName,
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

	// Set usage type to merchant - limited capabilities
	ut := UsageTypeMerchant
	linkedAccount.UsageType = &ut

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

	l := &LinkedBankAccount{}
	err = stmt.Get(l, &linkedAccount)
	if err != nil {
		return nil, err
	}

	return l, nil
}

func (db *linkedAccountDataStore) CreateLinkedBankAccount(linkedBankAccount *LinkedBankAccount) (*LinkedBankAccount, error) {

	if len(linkedBankAccount.RegisteredAccountId) == 0 {
		return nil, errors.New("Reference id missing")
	}

	// Default/mandatory fields
	columns := []string{
		"user_id", "business_id", "business_bank_account_id", "contact_id", "registered_account_id",
		"registered_bank_name", "account_holder_name", "account_name", "currency", "account_type",
		"usage_type", "account_number", "bank_name", "routing_number", "wire_routing", "source_account_id",
		"source_id", "source_name", "account_permission", "alias",
	}

	// Default/mandatory values
	values := []string{
		":user_id", ":business_id", ":business_bank_account_id", ":contact_id", ":registered_account_id",
		":registered_bank_name", ":account_holder_name", ":account_name", ":currency", ":account_type",
		":usage_type", ":account_number", ":bank_name", ":routing_number", ":wire_routing",
		":source_account_id ", ":source_id", ":source_name", ":account_permission", ":alias",
	}

	sql := fmt.Sprintf("INSERT INTO business_linked_bank_account(%s) VALUES(%s) RETURNING *", strings.Join(columns, ", "), strings.Join(values, ", "))

	stmt, err := db.PrepareNamed(sql)
	if err != nil {
		return nil, err
	}

	l := &LinkedBankAccount{}

	err = stmt.Get(l, &linkedBankAccount)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	return l, nil
}

func (db *linkedAccountDataStore) GetByAccountIDInternal(businessID shared.BusinessID, accountID string) (*LinkedBankAccount, error) {
	/* Don't check access for internal calls
	err := auth.NewAuthService(db.sourceReq).CheckBusinessAccess(businessID)
	if err != nil {
		return nil, err
	} */

	if os.Getenv("USE_BANKING_SERVICE") == "true" {
		las, err := NewBankingLinkedAccountService()
		if err != nil {
			return nil, err
		}

		return las.GetByAccountIDInternal(accountID)
	}

	a := LinkedBankAccount{}
	err := db.Get(
		&a, `
		SELECT * FROM business_linked_bank_account
		WHERE business_id = $1 AND business_bank_account_id = $2 AND deactivated IS NULL`,
		businessID,
		accountID,
	)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return &a, err
}

func (db *linkedAccountDataStore) GetByAccountNumber(businessID shared.BusinessID, accountNumber AccountNumber, routingNumber string) (*LinkedBankAccount, error) {
	// Check access
	err := auth.NewAuthService(db.sourceReq).CheckBusinessAccess(businessID)
	if err != nil {
		return nil, err
	}

	return db.GetByAccountNumberInternal(businessID, accountNumber, routingNumber)
}

func (db *linkedAccountDataStore) GetByAccountNumberInternal(businessID shared.BusinessID, accountNumber AccountNumber, routingNumber string) (*LinkedBankAccount, error) {
	if os.Getenv("USE_BANKING_SERVICE") == "true" {
		las, err := NewBankingLinkedAccountService()
		if err != nil {
			return nil, err
		}

		return las.GetByAccountNumber(businessID, accountNumber.String(), routingNumber)
	}

	a := LinkedBankAccount{}
	err := db.Get(
		&a, `
		SELECT * FROM business_linked_bank_account
		WHERE business_id = $1 AND account_number = $2 AND routing_number = $3 AND deactivated IS NULL`,
		businessID,
		accountNumber,
		routingNumber,
	)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return &a, err
}

func (db *linkedAccountDataStore) List(offset int, limit int, businessID shared.BusinessID) ([]LinkedBankAccount, error) {
	// Check access
	err := auth.NewAuthService(db.sourceReq).CheckBusinessAccess(businessID)
	if err != nil {
		return nil, err
	}

	if os.Getenv("USE_BANKING_SERVICE") == "true" {
		las, err := NewBankingLinkedAccountService()
		if err != nil {
			return nil, err
		}

		stfs := []grpcBanking.LinkedSubtype{
			grpcBanking.LinkedSubtype_LST_PRIMARY,
			grpcBanking.LinkedSubtype_LST_EXTERNAL,
		}

		return las.List(businessID, stfs, limit, offset)
	}

	rows := []LinkedBankAccount{}
	err = db.Select(
		&rows,
		`SELECT * FROM business_linked_bank_account
		WHERE business_id = $1 AND deactivated IS NULL AND usage_type IN ($2, $3)`,
		businessID,
		UsageTypePrimary,
		UsageTypeExternal,
	)
	switch {
	case err == sql.ErrNoRows:
		return rows, nil
	case err != nil:
		return nil, errors.Cause(err)
	}

	return rows, err
}

func (db *linkedAccountDataStore) GetById(id string, businessID shared.BusinessID) (*LinkedBankAccount, error) {
	// Check access
	err := auth.NewAuthService(db.sourceReq).CheckBusinessLinkedAccountAccess(id)
	if err != nil {
		return nil, err
	}

	if os.Getenv("USE_BANKING_SERVICE") == "true" {
		las, err := NewBankingLinkedAccountService()
		if err != nil {
			return nil, err
		}

		return las.GetById(id)
	}

	u := &LinkedBankAccount{}
	err = db.Get(u, "SELECT * FROM business_linked_bank_account WHERE id = $1 AND business_id = $2", id, businessID)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (db *linkedAccountDataStore) GetByIDInternal(id string) (*LinkedBankAccount, error) {
	if os.Getenv("USE_BANKING_SERVICE") == "true" {
		las, err := NewBankingLinkedAccountService()
		if err != nil {
			return nil, err
		}

		return las.GetById(id)
	}

	u := &LinkedBankAccount{}
	err := db.Get(u, "SELECT * FROM business_linked_bank_account WHERE id = $1", id)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func transformPlaidResponseToLinkedAccount(
	plaidResponse *plaid.PlaidResponse,
	userID shared.UserID,
	businessID shared.BusinessID,
	accountID string,
	contactID *string,
) (*LinkedBankAccount, error) {

	for _, achNumber := range plaidResponse.Numbers.ACH {
		if accountID == achNumber.AccountId {
			linkedBankAccount := LinkedBankAccount{
				RoutingNumber: achNumber.Routing, WireRouting: &achNumber.WireRouting}
			linkedBankAccount.SourceAccountId = &achNumber.AccountId

			account, err := getAccountByID(plaidResponse.Account, achNumber.AccountId)
			if err != nil {
				return nil, err
			}

			linkedBankAccount.BusinessID = businessID
			linkedBankAccount.UserID = userID

			currency := banking.Currency(strings.ToLower(account.Balance.ISOCurrencyCode))
			linkedBankAccount.Currency = currency

			linkedBankAccount.AccountNumber = AccountNumber(achNumber.Account)
			linkedBankAccount.AccountName = &account.Name
			linkedBankAccount.AccountType = banking.AccountType(account.SubType)
			linkedBankAccount.SourceId = &plaidResponse.RequestID

			sourceName := string(LinkedAccountSourcePlaid)
			linkedBankAccount.SourceName = &sourceName

			// Permission is set to send and receive
			linkedBankAccount.Permission = banking.LinkedAccountPermissionSendAndRecieve

			// Set usage type as external
			var ut UsageType
			if contactID != nil {
				ut = UsageTypeContact
			} else {
				ut = UsageTypeExternal
			}
			linkedBankAccount.UsageType = &ut

			linkedBankAccount.ContactId = contactID

			return &linkedBankAccount, nil
		}

	}

	return nil, errors.New("Account Does not exist")

}

func getAccountByID(accounts []plaid.PlaidAccountResponse, accountId string) (*plaid.PlaidAccountResponse, error) {
	for _, account := range accounts {
		if account.AccountId == accountId {
			return &account, nil
		}
	}

	return nil, errors.New("Account Does not exist")

}

func transformWithResponse(response *partnerbank.LinkedBankAccountResponse, linkedAccount *LinkedBankAccount) *LinkedBankAccount {

	permission := banking.LinkedAccountPermission(response.Permission)
	currency := banking.Currency(response.Currency)

	linkedAccount.AccountNumber = AccountNumber(response.AccountNumber)
	linkedAccount.RoutingNumber = response.RoutingNumber
	linkedAccount.WireRouting = response.WireRouting
	linkedAccount.BankName = &response.AccountBankName
	linkedAccount.AccountType = banking.AccountType(response.AccountType)
	linkedAccount.AccountHolderName = response.AccountHolderName
	linkedAccount.Alias = &response.Alias
	linkedAccount.Permission = permission
	linkedAccount.Currency = currency
	linkedAccount.RegisteredAccountId = string(response.AccountID)
	linkedAccount.RegisteredBankName = "bbva"

	return linkedAccount
}

// Transform partner bank layer response
func transformResponse(response *partnerbank.LinkedBankAccountResponse) *LinkedBankAccount {

	permission := banking.LinkedAccountPermission(response.Permission)
	currency := banking.Currency(response.Currency)

	linkedAccount := &LinkedBankAccount{
		AccountNumber:       AccountNumber(response.AccountNumber),
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

func (db *linkedAccountDataStore) UnlinkBankAccount(id string, businessID shared.BusinessID) (*LinkedBankAccount, error) {

	// Get linked bank account
	la, err := db.GetById(id, businessID)
	if err != nil {
		return nil, err
	}

	// Only accounts linked via external sources(like plaid) can be unlinked
	if la.SourceId == nil {
		return nil, errors.New("This account cannot be unlinked")
	}

	if os.Getenv("USE_BANKING_SERVICE") == "true" {
		las, err := NewBankingLinkedAccountService()
		if err != nil {
			return nil, err
		}

		return las.UnlinkBankAccount(id)
	}

	// Unregister bank account with partner bank
	bank, err := partnerbank.GetBusinessBank(partnerbank.ProviderNameBBVA)
	if err != nil {
		return nil, err
	}

	srv, err := bank.LinkedAccountService(db.sourceReq.PartnerBankRequest(), partnerbank.BusinessID(businessID))
	if err != nil {
		return nil, err
	}

	err = srv.Unlink(partnerbank.LinkedAccountBankID(la.RegisteredAccountId))
	if err != nil {
		return nil, err
	}

	// Deactivate registered account
	_, err = db.Exec("UPDATE business_linked_bank_account SET deactivated = CURRENT_TIMESTAMP WHERE id = $1 AND business_id = $2", id, businessID)
	if err != nil {
		return nil, err
	}

	la, err = db.GetById(id, businessID)
	if err != nil {
		return nil, err
	}

	return la, nil
}

func (db *linkedAccountDataStore) addExternalAccountsAndOwners(bID shared.BusinessID, resp *plaid.PlaidResponse) error {
	for _, achNumber := range resp.Numbers.ACH {

		account, err := getAccountByID(resp.Account, achNumber.AccountId)
		if err != nil {
			return err
		}

		loginTime := time.Now()

		cu := ExternalBankAccountUpdate{
			BusinessID:          bID,
			PartnerAccountID:    &account.AccountId,
			PartnerName:         string(LinkedAccountSourcePlaid),
			AccountName:         account.Name,
			OfficialAccountName: account.OfficialName,
			AccountType:         account.Type,
			AccountSubtype:      account.SubType,
			AccountNumber:       achNumber.Account,
			RoutingNumber:       achNumber.Routing,
			WireRouting:         achNumber.WireRouting,
			AvailableBalance:    &account.Balance.Available,
			PostedBalance:       &account.Balance.Current,
			Currency:            &account.Balance.ISOCurrencyCode,
			LastLogin:           &loginTime,
		}

		owners := make([]ExternalBankAccountOwnerCreate, 0)
		for _, o := range account.Owner {
			addressJSON, err := json.Marshal(o.Address)
			if err != nil {
				return err
			}

			emailJSON, err := json.Marshal(o.Email)
			if err != nil {
				return err
			}

			phoneJSON, err := json.Marshal(o.Phone)
			if err != nil {
				return err
			}

			owner := ExternalBankAccountOwnerCreate{
				OwnerAddress:      types.JSONText(string(addressJSON)),
				Email:             types.JSONText(string(emailJSON)),
				Phone:             types.JSONText(string(phoneJSON)),
				AccountHolderName: o.Name,
			}

			owners = append(owners, owner)
		}

		cu.Owner = owners

		ea, err := NewExternalAccountService(db.sourceReq).Upsert(cu)
		if err != nil {
			return err
		}

		err = NewExternalAccountService(db.sourceReq).CreateOwners(ea.ID, cu.Owner)
		if err != nil {
			return err
		}

	}

	return nil
}
