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
	"github.com/wiseco/core-platform/services/data"
	"github.com/wiseco/core-platform/shared"
	grpcBanking "github.com/wiseco/protobuf/golang/banking"
)

type linkedCardDatastore struct {
	sourceReq services.SourceRequest
	*sqlx.DB
}

type LinkedCardService interface {
	// Read
	GetByIDInternal(string) (*business.LinkedCard, error)
	GetById(string, string, shared.BusinessID) (*business.LinkedCard, error)
	List(offset int, limit int, contactId string, businessID shared.BusinessID) ([]*business.LinkedCard, error)
	GetByLinkedCardHash(shared.BusinessID, string) ([]business.LinkedCard, error)
	GetByLinkedCardHashAndContactID(shared.BusinessID, string, string) (*business.LinkedCard, error)

	// Create
	Create(*business.LinkedCardCreate) (*business.LinkedCard, error)
	RegisterExistingCard(shared.BusinessID, string, string) (*business.LinkedCard, error)
	UpdateLinkedCardUsageType(*business.LinkedCardUpdate) error

	// Unlink and deactivate card
	Deactivate(string, string, shared.BusinessID) error
	DeactivateAll(string, shared.BusinessID) error
}

func NewLinkedCardService(r services.SourceRequest) LinkedCardService {
	return &linkedCardDatastore{r, data.DBWrite}
}

func NewLinkedCardServiceWithout() LinkedCardService {
	return &linkedCardDatastore{services.SourceRequest{}, data.DBWrite}
}

func (db *linkedCardDatastore) List(offset int, limit int, contactId string, businessID shared.BusinessID) ([]*business.LinkedCard, error) {
	// Check access
	err := auth.NewAuthService(db.sourceReq).CheckBusinessAccess(businessID)
	if err != nil {
		return nil, err
	}

	if os.Getenv("USE_BANKING_SERVICE") == "true" {
		lcs, err := business.NewBankingLinkedCardService()
		if err != nil {
			return nil, err
		}

		stfs := []grpcBanking.LinkedSubtype{
			grpcBanking.LinkedSubtype_LST_CONTACT,
		}

		return lcs.ListWithContact(businessID, contactId, stfs, limit, offset)
	}

	rows := []*business.LinkedCard{}

	err = db.Select(&rows, `SELECT * FROM business_linked_card WHERE contact_id = $1 AND business_id = $2 AND 
	usage_type is distinct from $3 AND deactivated IS NULL`, contactId, businessID, business.UsageTypeContactInvisible)
	switch {
	case err == sql.ErrNoRows:
		return rows, nil
	case err != nil:
		return nil, errors.Cause(err)
	}

	return rows, err
}

func (db *linkedCardDatastore) GetByIDInternal(id string) (*business.LinkedCard, error) {
	c := &business.LinkedCard{}

	if os.Getenv("USE_BANKING_SERVICE") == "true" {
		lcs, err := business.NewBankingLinkedCardService()
		if err != nil {
			return nil, err
		}

		return lcs.GetByID(id)
	}

	err := db.Get(c, "SELECT * FROM business_linked_card WHERE id = $1", id)
	return c, err
}

func (db *linkedCardDatastore) GetById(id string, contactId string, businessID shared.BusinessID) (*business.LinkedCard, error) {
	// Check access
	err := auth.NewAuthService(db.sourceReq).CheckBusinessAccess(businessID)
	if err != nil {
		return nil, err
	}

	if os.Getenv("USE_BANKING_SERVICE") == "true" {
		lcs, err := business.NewBankingLinkedCardService()
		if err != nil {
			return nil, err
		}

		return lcs.GetByID(id)
	}

	a := business.LinkedCard{}

	err = db.Get(&a, "SELECT * FROM business_linked_card WHERE id = $1 AND contact_id = $2 AND business_id = $3", id, contactId, businessID)
	if err != nil && err != sql.ErrNoRows {

		return nil, err
	}

	return &a, err
}

func (db *linkedCardDatastore) Create(cardCreate *business.LinkedCardCreate) (*business.LinkedCard, error) {
	// Check access
	err := auth.NewAuthService(db.sourceReq).CheckBusinessAccess(cardCreate.BusinessID)
	if err != nil {
		return nil, err
	}

	if cardCreate.CardHolderName == "" {
		return nil, errors.New("account holder name is required")
	}

	if cardCreate.CardNumber == "" {
		return nil, errors.New("card number is required")
	}

	if cardCreate.ExpirationDate.IsZero() {
		return nil, errors.New("expiration date is required")
	}

	if cardCreate.CVVCode == "" {
		return nil, errors.New("CVV code is required")
	}

	if cardCreate.BillingAddress == nil {
		return nil, errors.New("billing address is required")
	}

	if cardCreate.Alias == "" {
		return nil, errors.New("alias is required")
	}

	cardCreate.ValidateCard = false

	if os.Getenv("USE_BANKING_SERVICE") == "true" {
		lcs, err := business.NewBankingLinkedCardService()
		if err != nil {
			return nil, err
		}

		q := `SELECT consumer_id FROM wise_user WHERE id = $1`

		var cID shared.ConsumerID

		err = db.Get(&cID, q, cardCreate.UserID)
		if err != nil {
			return nil, err
		}

		return lcs.Create(cardCreate, cID)
	}

	addressRequest := cardCreate.BillingAddress.ToPartnerBankAddress(services.AddressTypeNone)
	request := partnerbank.LinkedCardRequest{
		AccountHolder:  cardCreate.CardHolderName,
		CardNumber:     string(cardCreate.CardNumber),
		Expiration:     cardCreate.ExpirationDate.Time(),
		CVC:            cardCreate.CVVCode,
		Permission:     partnerbank.LinkedCardPermission(cardCreate.Permission),
		BillingAddress: addressRequest,
		Alias:          cardCreate.Alias,
	}

	// Register card with partner bank
	bank, err := partnerbank.GetBusinessBank(partnerbank.ProviderNameBBVA)
	if err != nil {
		return nil, err
	}

	srv, err := bank.LinkedCardService(db.sourceReq.PartnerBankRequest(), partnerbank.BusinessID(cardCreate.BusinessID))
	if err != nil {
		return nil, err
	}

	resp, err := srv.Link(&request)
	if err != nil {
		return nil, err
	}

	linkedCard := transformCardResponse(resp, cardCreate.CardNumber, cardCreate.CardHolderName, cardCreate.CVVCode)
	linkedCard.BusinessID = cardCreate.BusinessID
	linkedCard.ContactId = cardCreate.ContactId
	linkedCard.CardNumberHashed = cardCreate.HashLinkedCard()
	linkedCard.UsageType = cardCreate.UsageType

	// Default/mandatory fields
	columns := []string{
		"business_id", "contact_id", "registered_card_id", "registered_bank_name", "card_number_masked", "card_brand", "card_type",
		"issuer_name", "fast_funds_enabled", "card_holder_name", "alias", "account_permission", "billing_address", "usage_type", "card_number_hashed",
	}
	// Default/mandatory values
	values := []string{
		":business_id", ":contact_id", ":registered_card_id", ":registered_bank_name", ":card_number_masked", ":card_brand", ":card_type",
		":issuer_name", ":fast_funds_enabled", ":card_holder_name", ":alias", ":account_permission", ":billing_address", ":usage_type", ":card_number_hashed",
	}

	sql := fmt.Sprintf("INSERT INTO business_linked_card(%s) VALUES(%s) RETURNING *", strings.Join(columns, ", "), strings.Join(values, ", "))

	stmt, err := db.PrepareNamed(sql)
	if err != nil {
		return nil, err
	}

	l := &business.LinkedCard{}

	err = stmt.Get(l, &linkedCard)
	if err != nil {
		return nil, err
	}

	return l, nil

}

func (db *linkedCardDatastore) RegisterExistingCard(businessID shared.BusinessID, contactID string, hash string) (*business.LinkedCard, error) {
	lcs, err := business.NewBankingLinkedCardService()
	if err != nil {
		return nil, err
	}

	return lcs.RegisterExistingCard(businessID, contactID, hash)
}

func (db *linkedCardDatastore) DeactivateAll(contactID string, businessID shared.BusinessID) error {

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

func (db *linkedCardDatastore) Deactivate(ID string, contactID string, businessID shared.BusinessID) error {
	// Check access
	err := auth.NewAuthService(db.sourceReq).CheckBusinessAccess(businessID)
	if err != nil {
		log.Println(err)
		return err
	}

	if os.Getenv("USE_BANKING_SERVICE") == "true" {
		lcs, err := business.NewBankingLinkedCardService()
		if err != nil {
			return err
		}

		return lcs.Delete(ID)
	}

	// Get linked card
	lc, err := db.GetById(ID, contactID, businessID)
	if err != nil {
		log.Println(err)
		return err
	}

	// Unlink card
	bank, err := partnerbank.GetBusinessBank(partnerbank.ProviderNameBBVA)
	if err != nil {
		log.Println(err)
		return err
	}

	srv, err := bank.LinkedCardService(db.sourceReq.PartnerBankRequest(), partnerbank.BusinessID(businessID))
	if err != nil {
		log.Println(err)
		return err
	}

	err = srv.Unlink(partnerbank.LinkedCardBankID(lc.RegisteredCardId))
	if err != nil {
		log.Println(err)
		return err
	}

	// Deactivate card
	_, err = db.Exec("UPDATE business_linked_card SET deactivated = CURRENT_TIMESTAMP WHERE id = $1 AND contact_id = $2 AND business_id = $3", ID, contactID, businessID)
	return err
}

func (db *linkedCardDatastore) UpdateLinkedCardUsageType(u *business.LinkedCardUpdate) error {
	var columns []string

	if os.Getenv("USE_BANKING_SERVICE") == "true" {
		lcs, err := business.NewBankingLinkedCardService()
		if err != nil {
			return err
		}

		_, err = lcs.UpdateLinkedCardUsageType(u.ID, u.UsageType)
		return err
	}

	if u.UsageType != nil {
		columns = append(columns, "usage_type = :usage_type")
	}

	_, err := db.NamedExec(fmt.Sprintf("UPDATE business_linked_card SET %s WHERE id = '%s'", strings.Join(columns, ", "), u.ID), u)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (db *linkedCardDatastore) GetByLinkedCardHash(businessID shared.BusinessID, hash string) ([]business.LinkedCard, error) {
	// Check access
	err := auth.NewAuthService(db.sourceReq).CheckBusinessAccess(businessID)
	if err != nil {
		return nil, err
	}

	if os.Getenv("USE_BANKING_SERVICE") == "true" {
		lcs, err := business.NewBankingLinkedCardService()
		if err != nil {
			return nil, err
		}

		return lcs.GetByLinkedCardHash(businessID, hash)
	}

	a := business.LinkedCard{}

	err = db.Get(&a, "SELECT * FROM business_linked_card WHERE card_number_hashed = $1", hash)
	if err != nil {
		return nil, err
	}

	ret := []business.LinkedCard{
		a,
	}

	return ret, err
}

func (db *linkedCardDatastore) GetByLinkedCardHashAndContactID(businessID shared.BusinessID, cID string, hash string) (*business.LinkedCard, error) {
	// Check access
	err := auth.NewAuthService(db.sourceReq).CheckBusinessAccess(businessID)
	if err != nil {
		return nil, err
	}

	if os.Getenv("USE_BANKING_SERVICE") == "true" {
		lcs, err := business.NewBankingLinkedCardService()
		if err != nil {
			return nil, err
		}

		return lcs.GetByLinkedCardHashAndContactID(cID, hash)
	}

	a := business.LinkedCard{}

	err = db.Get(&a, "SELECT * FROM business_linked_card WHERE card_number_hashed = $1", hash)
	if err != nil {
		return nil, err
	}

	return &a, err
}

// Transform response from partner bank layer
func transformCardResponse(response *partnerbank.LinkedCardResponse, cardNumber business.CardNumber, holderName, cvv string) business.LinkedCard {

	permission := banking.LinkedAccountPermission(response.Permission)
	address := services.Address{
		Type:          services.AddressType(response.BillingAddress.Type),
		StreetAddress: response.BillingAddress.Line1,
		AddressLine2:  response.BillingAddress.Line2,
		City:          response.BillingAddress.City,
		State:         response.BillingAddress.State,
		PostalCode:    response.BillingAddress.ZipCode,
	}

	linkedCard := business.LinkedCard{
		Permission:         permission,
		RegisteredCardId:   string(response.CardID),
		RegisteredBankName: "bbva",
		CardNumberMasked:   business.CardNumber(services.MaskLeft(cardNumber.String(), 4)),
		CardBrand:          string(response.BrandType),
		CardType:           string(response.CardType),
		CardIssuer:         response.IssuingBank,
		BillingAddress:     &address,
		CardHolderName:     holderName,
		Alias:              &response.Alias,
		FastFundsEnabled:   response.FastFundsEnabled,
	}

	return linkedCard
}
