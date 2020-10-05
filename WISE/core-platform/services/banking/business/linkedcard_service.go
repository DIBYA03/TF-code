/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

// Package for all internal banking related items such as accounts and cards
package business

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/wiseco/core-platform/auth"
	partnerbank "github.com/wiseco/core-platform/partner/bank"
	"github.com/wiseco/core-platform/services"
	"github.com/wiseco/core-platform/services/banking"
	"github.com/wiseco/core-platform/services/business"
	"github.com/wiseco/core-platform/services/data"
	"github.com/wiseco/core-platform/services/user"
	"github.com/wiseco/core-platform/shared"
)

type linkedCardDatastore struct {
	sourceReq services.SourceRequest
	*sqlx.DB
}

type LinkedCardService interface {
	Create(*LinkedCardCreate) (*LinkedCard, error)
	Delete(string, shared.BusinessID) error

	GetByID(string, shared.BusinessID) (*LinkedCard, error)
	GetByLinkedCardHash(shared.BusinessID, string) ([]LinkedCard, error)
	GetByLinkedCardHashAndContactID(shared.BusinessID, string, string) (*LinkedCard, error)
	List(offset int, limit int, businessID shared.BusinessID) ([]*LinkedCard, error)
}

func NewLinkedCardService(r services.SourceRequest) LinkedCardService {
	return &linkedCardDatastore{r, data.DBWrite}
}

func (db *linkedCardDatastore) Create(cardCreate *LinkedCardCreate) (*LinkedCard, error) {
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

	cardCreate.ValidateCard = true

	// Verify billing address
	err = db.verifyAddress(cardCreate.BusinessID, cardCreate.BillingAddress)
	if err != nil {
		return nil, err
	}

	if os.Getenv("USE_BANKING_SERVICE") == "true" {
		lcs, err := NewBankingLinkedCardService()
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

	// Compare hash to see if card was already registered
	hash := cardCreate.HashLinkedCard()
	lc, err := db.GetByLinkedCardHashAndContactID(cardCreate.BusinessID, "not used outside of banking service", *hash)
	if err != nil && err != sql.ErrNoRows {
		return nil, errors.New("error fetching linked card hash")
	}

	if lc != nil && len(lc.Id) > 0 {
		return nil, errors.New("card already registered with this business")
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

	linkedCard := transformLinkedCardResponse(resp, cardCreate.CardNumber, cardCreate.CardHolderName, cardCreate.CVVCode)
	linkedCard.BusinessID = cardCreate.BusinessID
	linkedCard.CardNumberHashed = cardCreate.HashLinkedCard()

	usageType := UsageTypeExternal
	linkedCard.UsageType = &usageType

	// Default/mandatory fields
	columns := []string{
		"business_id", "registered_card_id", "registered_bank_name", "card_number_masked", "card_brand", "card_type",
		"issuer_name", "fast_funds_enabled", "card_holder_name", "alias", "account_permission", "billing_address", "usage_type", "card_number_hashed",
	}
	// Default/mandatory values
	values := []string{
		":business_id", ":registered_card_id", ":registered_bank_name", ":card_number_masked", ":card_brand", ":card_type",
		":issuer_name", ":fast_funds_enabled", ":card_holder_name", ":alias", ":account_permission", ":billing_address", ":usage_type", ":card_number_hashed",
	}

	sql := fmt.Sprintf("INSERT INTO business_linked_card(%s) VALUES(%s) RETURNING *", strings.Join(columns, ", "), strings.Join(values, ", "))

	stmt, err := db.PrepareNamed(sql)
	if err != nil {
		return nil, err
	}

	l := &LinkedCard{}

	err = stmt.Get(l, &linkedCard)
	if err != nil {
		return nil, err
	}

	return l, nil

}

func (db *linkedCardDatastore) Delete(lcID string, bID shared.BusinessID) error {
	// Check access
	err := auth.NewAuthService(db.sourceReq).CheckBusinessAccess(bID)
	if err != nil {
		return err
	}

	lcs, err := NewBankingLinkedCardService()
	if err != nil {
		return err
	}

	return lcs.Delete(lcID)
}

// Transform response from partner bank layer
func transformLinkedCardResponse(response *partnerbank.LinkedCardResponse, cardNumber CardNumber, holderName, cvv string) LinkedCard {

	permission := banking.LinkedAccountPermission(response.Permission)
	address := services.Address{
		Type:          services.AddressType(response.BillingAddress.Type),
		StreetAddress: response.BillingAddress.Line1,
		AddressLine2:  response.BillingAddress.Line2,
		City:          response.BillingAddress.City,
		State:         response.BillingAddress.State,
		PostalCode:    response.BillingAddress.ZipCode,
	}

	linkedCard := LinkedCard{
		Permission:         permission,
		RegisteredCardId:   string(response.CardID),
		RegisteredBankName: "bbva",
		CardNumberMasked:   CardNumber(services.MaskLeft(cardNumber.String(), 4)),
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

func (db *linkedCardDatastore) GetByID(ID string, businessID shared.BusinessID) (*LinkedCard, error) {
	// Check access
	err := auth.NewAuthService(db.sourceReq).CheckBusinessAccess(businessID)
	if err != nil {
		return nil, err
	}

	if os.Getenv("USE_BANKING_SERVICE") == "true" {
		lcs, err := NewBankingLinkedCardService()
		if err != nil {
			return nil, err
		}

		return lcs.GetByID(ID)
	}

	a := LinkedCard{}

	err = db.Get(&a, "SELECT * FROM business_linked_card WHERE id = $1 AND business_id = $2", ID, businessID)
	if err != nil && err != sql.ErrNoRows {

		return nil, err
	}

	return &a, err
}

func (db *linkedCardDatastore) List(offset int, limit int, businessID shared.BusinessID) ([]*LinkedCard, error) {
	// Check access
	err := auth.NewAuthService(db.sourceReq).CheckBusinessAccess(businessID)
	if err != nil {
		return nil, err
	}

	if os.Getenv("USE_BANKING_SERVICE") == "true" {
		lcs, err := NewBankingLinkedCardService()
		if err != nil {
			return nil, err
		}

		return lcs.List(businessID, limit, offset)
	}

	rows := []*LinkedCard{}

	err = db.Select(&rows, `SELECT * FROM business_linked_card WHERE business_id = $1 AND
	usage_type = $2 AND deactivated IS NULL`, businessID, UsageTypeExternal)
	switch {
	case err == sql.ErrNoRows:
		return rows, nil
	case err != nil:
		return nil, errors.Cause(err)
	}

	return rows, err
}

func (db *linkedCardDatastore) GetByLinkedCardHash(businessID shared.BusinessID, hash string) ([]LinkedCard, error) {
	// Check access
	err := auth.NewAuthService(db.sourceReq).CheckBusinessAccess(businessID)
	if err != nil {
		return nil, err
	}

	if os.Getenv("USE_BANKING_SERVICE") == "true" {
		lcs, err := NewBankingLinkedCardService()
		if err != nil {
			return nil, err
		}

		return lcs.GetByLinkedCardHash(businessID, hash)
	}

	a := LinkedCard{}

	err = db.Get(&a, "SELECT * FROM business_linked_card WHERE card_number_hashed = $1", hash)
	if err != nil {
		return nil, err
	}

	ret := []LinkedCard{
		a,
	}

	return ret, err
}

func (db *linkedCardDatastore) GetByLinkedCardHashAndContactID(businessID shared.BusinessID, cID string, hash string) (*LinkedCard, error) {
	// Check access
	err := auth.NewAuthService(db.sourceReq).CheckBusinessAccess(businessID)
	if err != nil {
		return nil, err
	}

	if os.Getenv("USE_BANKING_SERVICE") == "true" {
		lcs, err := NewBankingLinkedCardService()
		if err != nil {
			return nil, err
		}

		return lcs.GetByLinkedCardHashAndContactID(cID, hash)
	}

	a := LinkedCard{}

	err = db.Get(&a, "SELECT * FROM business_linked_card WHERE card_number_hashed = $1", hash)
	if err != nil {
		return nil, err
	}

	return &a, err
}

func (db *linkedCardDatastore) verifyAddress(bID shared.BusinessID, ba *services.Address) error {
	// Get business
	b, err := business.NewBusinessService(db.sourceReq).GetById(bID)
	if err != nil {
		return err
	}

	// Get user
	u, err := user.NewUserService(db.sourceReq).GetById(b.OwnerID)
	if err != nil {
		return err
	}

	// Verify user address
	if db.verifyUserAddress(u, ba) {
		return nil
	}

	// Verify business address
	if db.verifyBusinessAddress(b, ba) {
		return nil
	}

	return errors.New("Billing address should match with consumer address or business address")
}

func (db *linkedCardDatastore) verifyUserAddress(u *user.User, ba *services.Address) bool {
	if db.addressCheck(u.LegalAddress, ba) {
		return true
	}

	if db.addressCheck(u.MailingAddress, ba) {
		return true
	}

	if db.addressCheck(u.WorkAddress, ba) {
		return true
	}

	return false
}

func (db *linkedCardDatastore) verifyBusinessAddress(b *business.Business, ba *services.Address) bool {
	if db.addressCheck(b.LegalAddress, ba) {
		return true
	}

	if db.addressCheck(b.MailingAddress, ba) {
		return true
	}

	if db.addressCheck(b.HeadquarterAddress, ba) {
		return true
	}

	return false
}

func (db *linkedCardDatastore) addressCheck(ua *services.Address, ba *services.Address) bool {
	if ua == nil {
		return false
	}

	if !strings.EqualFold(strings.TrimSpace(ua.StreetAddress), strings.TrimSpace(ba.StreetAddress)) {
		return false
	}

	if !strings.EqualFold(strings.TrimSpace(ua.AddressLine2), strings.TrimSpace(ba.AddressLine2)) {
		return false
	}

	if !strings.EqualFold(strings.TrimSpace(ua.City), strings.TrimSpace(ba.City)) {
		return false
	}

	if !strings.EqualFold(strings.TrimSpace(ua.State), strings.TrimSpace(ba.State)) {
		return false
	}

	if !strings.EqualFold(strings.TrimSpace(ua.Country), strings.TrimSpace(ba.Country)) {
		return false
	}

	if !strings.EqualFold(strings.TrimSpace(ua.PostalCode), strings.TrimSpace(ba.PostalCode)) {
		return false
	}

	return true
}
