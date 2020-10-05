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
	"strconv"
	"strings"

	"github.com/microcosm-cc/bluemonday"
	"github.com/wiseco/core-platform/auth"
	"github.com/wiseco/core-platform/shared"
	"mvdan.cc/xurls/v2"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	partnerbank "github.com/wiseco/core-platform/partner/bank"
	"github.com/wiseco/core-platform/services"
	"github.com/wiseco/core-platform/services/banking"
	"github.com/wiseco/core-platform/services/banking/business"
	"github.com/wiseco/core-platform/services/contact"
	"github.com/wiseco/core-platform/services/data"
)

type moneyTransferDatastore struct {
	sourceReq services.SourceRequest
	*sqlx.DB
}

type MoneyTransferService interface {
	// Read
	GetById(string, string, shared.BusinessID) (*business.MoneyTransfer, error)
	GetByContactId(int, int, string, shared.BusinessID) ([]business.MoneyTransfer, error)

	GetContactByTransferId(string, shared.BusinessID) (*ContactTransferDetails, error)
	GetContactByID(string, shared.BusinessID) (*ContactTransferDetails, error)

	//Transfer
	Transfer(*business.TransferInitiate) (*business.MoneyTransfer, error)

	//Cancel
	Cancel(string, string, shared.BusinessID) (*business.TransferCancel, error)
}

func NewMoneyTransferService(r services.SourceRequest) MoneyTransferService {
	return &moneyTransferDatastore{r, data.DBWrite}
}

func (db *moneyTransferDatastore) GetById(id, contactId string, businessID shared.BusinessID) (*business.MoneyTransfer, error) {
	// Check access
	s := db.sourceReq

	err := auth.NewAuthService(s).CheckBusinessAccess(businessID)
	if err != nil {
		return nil, err
	}

	if os.Getenv("USE_BANKING_SERVICE") == "true" {
		bts, err := business.NewBankingTransferService()
		if err != nil {
			return nil, err
		}

		return bts.GetByIDInternal(businessID, id)
	}

	a := business.MoneyTransfer{}

	err = db.Get(
		&a, `
		SELECT * FROM business_money_transfer
		WHERE
			business_money_transfer.id = $1 AND
			business_money_transfer.contact_id = $2 AND
			business_money_transfer.business_id = $3`,
		id,
		contactId,
		businessID,
	)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return &a, err
}

func (db *moneyTransferDatastore) GetByContactId(offset int, limit int, id string, businessID shared.BusinessID) ([]business.MoneyTransfer, error) {
	// Check access
	err := auth.NewAuthService(db.sourceReq).CheckBusinessAccess(businessID)
	if err != nil {
		return nil, err
	}

	rows := []business.MoneyTransfer{}

	if os.Getenv("USE_BANKING_SERVICE") == "true" {
		bts, err := business.NewBankingTransferService()
		if err != nil {
			return rows, err
		}

		return bts.GetByBusinessAndContactID(businessID, id, offset, limit)
	}

	err = db.Select(
		&rows, `
		SELECT * FROM business_money_transfer
		WHERE
			business_money_transfer.contact_id = $1 AND
			business_money_transfer.business_id = $2`,
		id,
		businessID,
	)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return rows, err
}

type ContactTransferDetails struct {
	Contact       contact.Contact `db:"business_contact"`
	Notes         *string         `db:"notes"`
	AccountNumber *string         `db:"account_number"`
	SendEmail     bool            `db:"send_email"`
}

func (db *moneyTransferDatastore) GetContactByTransferId(transferId string, businessID shared.BusinessID) (*ContactTransferDetails, error) {
	var contact ContactTransferDetails

	if os.Getenv("USE_BANKING_SERVICE") == "true" {
		bts, err := business.NewBankingTransferService()
		if err != nil {
			return &contact, err
		}

		bas, err := business.NewBankingAccountService()
		if err != nil {
			return &contact, err
		}

		t, err := bts.GetByBankID(businessID, transferId)
		if err != nil {
			return &contact, err
		}

		bs, err := bas.GetByBusinessID(businessID, 1, 0)
		if err != nil {
			return &contact, err
		}

		if len(bs) != 1 {
			return &contact, fmt.Errorf("Could not find bank account for transferId:%s and businessID:%s", transferId, businessID)
		}

		ba := bs[0]

		err = db.Get(
			&contact, `
		SELECT
			business_contact.user_id "business_contact.user_id",
			business_contact.first_name "business_contact.first_name",
			business_contact.last_name "business_contact.last_name",
			business_contact.email "business_contact.email",
			business_contact.phone_number "business_contact.phone_number",
			business_contact.id "business_contact.id",
			business_contact.contact_category "business_contact.contact_category",
			business_contact.contact_type "business_contact.contact_type",
			business_contact.engagement "business_contact.engagement",
			business_contact.job_title "business_contact.job_title",
			business_contact.business_name "business_contact.business_name",
			business_contact.created "business_contact.created",
			business_contact.modified "business_contact.modified"
		FROM business_contact
		WHERE business_contact.id = $1`,
			t.ContactId,
		)
		if err != nil {
			return &contact, err
		}

		contact.AccountNumber = &ba.AccountNumber
		contact.Notes = t.Notes
		contact.SendEmail = t.SendEmail

		return &contact, err
	}

	err := db.Get(
		&contact, `
		SELECT
			business_contact.user_id "business_contact.user_id",
			business_contact.first_name "business_contact.first_name",
			business_contact.last_name "business_contact.last_name",
			business_contact.email "business_contact.email",
			business_contact.phone_number "business_contact.phone_number",
			business_contact.id "business_contact.id",
			business_contact.contact_category "business_contact.contact_category",
			business_contact.contact_type "business_contact.contact_type",
			business_contact.engagement "business_contact.engagement",
			business_contact.job_title "business_contact.job_title",
			business_contact.business_name "business_contact.business_name",
			business_contact.created "business_contact.created",
			business_contact.modified "business_contact.modified",
			business_money_transfer.notes "notes",
			business_money_transfer.send_email "send_email",
			business_bank_account.account_number "account_number"
		FROM business_contact
		JOIN business_money_transfer ON business_contact.id = business_money_transfer.contact_id
		JOIN business_bank_account ON business_contact.business_id = business_bank_account.business_id
		WHERE bank_transfer_id = $1 AND business_contact.business_id = $2`,
		transferId,
		businessID,
	)

	return &contact, err

}

func (db *moneyTransferDatastore) GetContactByID(ID string, businessID shared.BusinessID) (*ContactTransferDetails, error) {
	var contact ContactTransferDetails

	if os.Getenv("USE_BANKING_SERVICE") == "true" {
		bts, err := business.NewBankingTransferService()
		if err != nil {
			return &contact, err
		}

		t, err := bts.GetByIDInternal(businessID, ID)
		if err != nil {
			return &contact, err
		}

		err = db.Get(
			&contact, `
		SELECT
			business_contact.user_id "business_contact.user_id",
			business_contact.first_name "business_contact.first_name",
			business_contact.last_name "business_contact.last_name",
			business_contact.email "business_contact.email",
			business_contact.phone_number "business_contact.phone_number",
			business_contact.id "business_contact.id",
			business_contact.contact_category "business_contact.contact_category",
			business_contact.contact_type "business_contact.contact_type",
			business_contact.engagement "business_contact.engagement",
			business_contact.job_title "business_contact.job_title",
			business_contact.business_name "business_contact.business_name",
			business_contact.created "business_contact.created",
			business_contact.modified "business_contact.modified",
			business_money_transfer.notes "notes",
			business_money_transfer.send_email "send_email",
			business_bank_account.account_number "account_number"
		FROM business_contact
		JOIN business_money_transfer ON business_contact.id = business_money_transfer.contact_id
		JOIN business_bank_account ON business_contact.business_id = business_bank_account.business_id
		WHERE business_contact.id = $1`,
			t.ContactId,
		)

		return &contact, err
	}

	err := db.Get(
		&contact, `
		SELECT
			business_contact.user_id "business_contact.user_id",
			business_contact.first_name "business_contact.first_name",
			business_contact.last_name "business_contact.last_name",
			business_contact.email "business_contact.email",
			business_contact.phone_number "business_contact.phone_number",
			business_contact.id "business_contact.id",
			business_contact.contact_category "business_contact.contact_category",
			business_contact.contact_type "business_contact.contact_type",
			business_contact.engagement "business_contact.engagement",
			business_contact.job_title "business_contact.job_title",
			business_contact.business_name "business_contact.business_name",
			business_contact.created "business_contact.created",
			business_contact.modified "business_contact.modified",
			business_money_transfer.notes "notes",
			business_money_transfer.send_email "send_email",
			business_bank_account.account_number "account_number"
		FROM business_contact
		JOIN business_money_transfer ON business_contact.id = business_money_transfer.contact_id
		JOIN business_bank_account ON business_contact.business_id = business_bank_account.business_id
		WHERE business_money_transfer.id = $1 AND business_contact.business_id = $2`,
		ID,
		businessID,
	)

	return &contact, err
}

func (db *moneyTransferDatastore) Transfer(transfer *business.TransferInitiate) (*business.MoneyTransfer, error) {
	// Check access
	err := auth.NewAuthService(db.sourceReq).CheckBusinessAccess(transfer.BusinessID)
	if err != nil {
		return nil, err
	}

	if transfer.SourceAccountId == "" {
		return nil, errors.New("Source account id is required")
	}

	if transfer.DestAccountId == "" {
		return nil, errors.New("Destination account id is required")
	}

	if transfer.Amount == 0 {
		return nil, errors.New("Amount is required")
	}

	if transfer.SourceType == "" {
		return nil, errors.New("Source type is required")
	}

	if transfer.SourceType != banking.TransferTypeAccount {
		return nil, errors.New("Only account type is supported for source")
	}

	if transfer.DestType == "" {
		return nil, errors.New("Destination type is required")
	}

	//TODO this is obivous but messy, we create a map in the transfer class and test for inclusion
	//or a method on the transfer struct that tests the inlcusion of DestType
	if transfer.DestType != banking.TransferTypeAccount && transfer.DestType != banking.TransferTypeCard && transfer.DestType != banking.TransferTypeCheck {
		return nil, errors.New("Only account and card types are supported for destination")
	}

	if transfer.Notes != nil && len(*transfer.Notes) > 0 {
		rxRelaxed := xurls.Relaxed()
		if len(rxRelaxed.FindString(*transfer.Notes)) > 0 {
			return nil, errors.New("Notes cannot contain urls")
		}

		// sanitize notes
		policy := bluemonday.StrictPolicy()
		notes := policy.Sanitize(
			*transfer.Notes,
		)
		*transfer.Notes = notes
	}

	// Search source account in linked accounts
	sa, err := business.NewLinkedAccountService(db.sourceReq).GetById(transfer.SourceAccountId, transfer.BusinessID)
	if err != nil {
		log.Println("unable to get source account", transfer.SourceAccountId, transfer.BusinessID)
		return nil, err
	}

	// Check for source account usage type - must be primary or clearing
	if sa.UsageType == nil {
		return nil, errors.New("source account type must be primary or clearing")
	}

	if os.Getenv("USE_BANKING_SERVICE") == "true" {
		bts, err := business.NewBankingTransferService()
		if err != nil {
			return nil, err
		}

		return bts.Transfer(transfer, *sa.UsageType, db.sourceReq)
	}

	// Check if account has sufficient funds to transfer
	if sa.BusinessBankAccountId != nil {
		balance, err := business.NewBankAccountService(db.sourceReq).GetBalanceByID(*sa.BusinessBankAccountId, transfer.BusinessID)
		if err != nil {
			log.Println("error fetching account balance", err)
			return nil, err
		}

		if balance.PostedBalance < transfer.Amount {
			err := errors.New("Insufficient funds to initiate move money")
			log.Println(err)
			return nil, err
		}
	}

	sourceAccountId := sa.RegisteredAccountId

	destinationAccountId, destService, err := db.getDestinationIDAndService(transfer)
	if err != nil {
		log.Println("unable to get destination account", transfer.DestAccountId, transfer.BusinessID)
		return nil, err
	}

	switch transfer.DestType {
	case banking.TransferTypeCheck:
		if transfer.AddressID == "" {
			return nil, errors.New("Address ID is required to pay by check")
		}

		maxCheckAmount, err := strconv.ParseFloat(os.Getenv("MAX_CHECK_AMOUNT_ALLOWED"), 64)
		if err != nil {
			return nil, err
		}

		if transfer.Amount > maxCheckAmount {
			e := fmt.Sprintf("Check amount cannot exceed $%s", strconv.FormatFloat(maxCheckAmount, 'f', 2, 64))
			return nil, errors.New(e)
		}
	case banking.TransferTypeAccount:
		da, ok := destService.(*business.LinkedBankAccount)
		if !ok {
			return nil, fmt.Errorf("Destination service is wrong type. Expecting linkedAccountDatastore got: %+v", destService)
		}

		maxACHAmount, err := strconv.ParseFloat(os.Getenv("ACH_MAX_ALLOWED"), 64)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		if !da.IntraBankAccount() && transfer.Amount > maxACHAmount {
			e := fmt.Sprintf("Transfers to external accounts are limited to $%s per transaction.", strconv.FormatFloat(maxACHAmount, 'f', 2, 64))
			return nil, errors.New(e)
		}
	}

	request := partnerbank.MoneyTransferRequest{
		Amount:          transfer.Amount,
		Currency:        partnerbank.Currency(transfer.Currency),
		SourceAccountID: partnerbank.MoneyTransferAccountBankID(sourceAccountId),
		DestAccountID:   partnerbank.MoneyTransferAccountBankID(destinationAccountId),
	}

	bank, err := partnerbank.GetBusinessBank(partnerbank.ProviderNameBBVA)
	if err != nil {
		return nil, err
	}

	srv, err := bank.MoneyTransferService(db.sourceReq.PartnerBankRequest(), partnerbank.BusinessID(transfer.BusinessID))
	if err != nil {
		return nil, err
	}

	// Transfer money
	resp, err := srv.Submit(&request)
	if err != nil {
		return nil, err
	}

	t := transformTransferResponse(resp)
	t.BusinessID = transfer.BusinessID
	t.ContactId = transfer.ContactId
	t.CreatedUserID = transfer.CreatedUserID
	t.SourceAccountId = transfer.SourceAccountId
	t.SourceType = transfer.SourceType
	t.DestAccountId = transfer.DestAccountId
	t.DestType = transfer.DestType
	t.Notes = transfer.Notes
	t.SendEmail = transfer.SendEmail
	t.MoneyRequestID = transfer.MoneyRequestID
	t.MonthlyInterestID = transfer.MonthlyInterestID

	// Default/mandatory fields
	columns := []string{
		"created_user_id", "business_id", "contact_id", "bank_name", "bank_transfer_id",
		"source_account_id", "source_type", "dest_account_id", "dest_type", "amount", "currency", "notes",
		"status", "send_email", "money_request_id", "account_monthly_interest_id",
	}
	// Default/mandatory values
	values := []string{
		":created_user_id", ":business_id", ":contact_id", ":bank_name", ":bank_transfer_id",
		":source_account_id", ":source_type", ":dest_account_id", ":dest_type", ":amount", ":currency",
		":notes", ":status", ":send_email", ":money_request_id", ":account_monthly_interest_id",
	}

	sql := fmt.Sprintf("INSERT INTO business_money_transfer(%s) VALUES(%s) RETURNING *", strings.Join(columns, ", "), strings.Join(values, ", "))

	stmt, err := db.PrepareNamed(sql)
	if err != nil {
		return nil, err
	}

	l := &business.MoneyTransfer{}

	err = stmt.Get(l, &t)
	if err != nil {
		return nil, err
	}

	// Create pending transaction
	if transfer.DestType == banking.TransferTypeAccount {
		da, ok := destService.(*business.LinkedBankAccount)
		if !ok {
			return nil, fmt.Errorf("Destination service is wrong type. Expecting linkedAccountDatastore got: %+v", destService)
		}

		business.NewMoneyTransferService(db.sourceReq).OnMoneyTransfer(l, sa.BusinessBankAccountId, da.BusinessBankAccountId, banking.PartnerNameBBVA)
	}

	return l, nil
}

//Not porting this method over to banking service, looks like its not used, and its noop
func (db *moneyTransferDatastore) Cancel(id, contactId string, businessID shared.BusinessID) (*business.TransferCancel, error) {
	a := business.TransferCancel{}

	err := db.Get(&a, "SELECT * FROM business_money_transfer WHERE id = $1 AND contact_id = $2 AND business_id = $3", id, contactId, businessID)
	if err != nil && err != sql.ErrNoRows {

		return nil, err
	}

	return &a, err
}

func (db *moneyTransferDatastore) getDestinationIDAndService(transfer *business.TransferInitiate) (string, interface{}, error) {

	if transfer.DestType == banking.TransferTypeAccount {
		// Search destination account in linked accounts
		da, err := NewLinkedAccountService(db.sourceReq).GetById(transfer.DestAccountId, *transfer.ContactId, transfer.BusinessID)
		if err != nil {
			log.Println("unable to get dest account", transfer.DestAccountId, *transfer.ContactId, transfer.BusinessID)
			return "", nil, err
		}

		return da.RegisteredAccountId, da, nil

	} else if transfer.DestType == banking.TransferTypeCard {
		// Search destination account in linked cards
		card, err := NewLinkedCardService(db.sourceReq).GetById(transfer.DestAccountId, *transfer.ContactId, transfer.BusinessID)
		if err != nil {
			return "", nil, err
		}

		return card.RegisteredCardId, card, nil

	} else if transfer.DestType == banking.TransferTypeCheck {
		// Search destination account in payees
		payeeService := business.NewLinkedPayeeService(db.sourceReq)

		payee, err := payeeService.GetByAddressID(transfer.AddressID, transfer.BusinessID)

		if err != nil {
			switch err {
			case sql.ErrNoRows:
				payeeCreate := business.LinkedPayeeCreate{
					BusinessID: transfer.BusinessID,
					ContactID:  *transfer.ContactId,
					AddressID:  transfer.AddressID,
				}

				//TODO I don't really like this as a side effect of this method.
				//Should we add payees when the contact is created?
				//Since most folks wont be sending checks to all there contacts it seems wasteful to just implictly add the contact
				//Fix this when we move banking into it's own service
				payee, err = payeeService.Create(payeeCreate)
				if err != nil {
					return "", nil, err
				}

			default:
				return "", nil, err
			}
		}

		//Set this to a relevant ID now that we have it
		transfer.DestAccountId = payee.ID

		return payee.BankPayeeID, payee, nil
	}

	return "", nil, fmt.Errorf("Unable to get destination for type: %s", transfer.DestAccountId)
}

// Transform partner bank layer response
func transformTransferResponse(response *partnerbank.MoneyTransferResponse) business.MoneyTransfer {

	t := business.MoneyTransfer{}

	t.BankName = "bbva"
	t.BankTransferId = string(response.TransferID)
	t.Amount = response.Amount
	t.Currency = banking.Currency(response.Currency)
	t.Status = string(response.Status)

	return t
}
