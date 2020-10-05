/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

// Package for all business related services
package business

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/wiseco/core-platform/auth"
	"github.com/wiseco/core-platform/services/banking"
	"github.com/wiseco/core-platform/shared"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	partnerbank "github.com/wiseco/core-platform/partner/bank"
	"github.com/wiseco/core-platform/services"
	"github.com/wiseco/core-platform/services/data"
)

type bankCardDatastore struct {
	sourceReq services.SourceRequest
	*sqlx.DB
}

type BankCardService interface {
	// Read
	GetById(id string, businessID shared.BusinessID, userID shared.UserID) (*BankCard, error)
	GetByIDInternal(id string) (*BankCard, error)
	GetByBankCardId(bankCardId string, userID shared.UserID) (*BankCard, error)
	GetByBusinessID(offset int, limit int, businessID shared.BusinessID, userID shared.UserID) ([]BankCard, error)
	GetByAccountId(offset int, limit int, accountId string, businessID shared.BusinessID, userID shared.UserID) ([]BankCard, error)
	GetByAccountInternal(offset int, limit int, accountId string, businessID shared.BusinessID) ([]*BankCard, error)

	// Create business debit card
	CreateBankCard(*BankCardCreate) (*BankCard, error)

	// Activate business debit card
	ActivateBankCard(*BankCardActivate) (*BankCard, error)

	// Block business debit card
	GetCardBlocks(string) ([]*BankCardBlock, error)
	GetByBlockID(shared.BankCardBlockID) (*BankCardBlock, error)
	GetByPartnerBlockID(banking.CardBlockID, string) (*BankCardBlock, error)
	GetByBlockIDAndDate(blockID string, blockDate time.Time) (*BankCardBlock, error)
	BlockBankCard(*BankCardBlockCreate) (*BankCard, error)
	UnBlockBankCard(*BankCardBlockDelete) (*BankCard, error)
	CheckBlockStatus(string) (*BankCard, error)

	// Update card status
	UpdateCardStatus(cardID string, userID shared.UserID, status string) error

	// Cancel Card
	CancelCardInternal(cardID string) (*BankCard, error)
}

func NewCardService(r services.SourceRequest) BankCardService {
	return &bankCardDatastore{r, data.DBWrite}
}

func NewCardServiceWithout() BankCardService {
	return &bankCardDatastore{services.NewSourceRequest(), data.DBWrite}
}

func (db *bankCardDatastore) GetById(id string, businessID shared.BusinessID, userID shared.UserID) (*BankCard, error) {
	// Check access
	err := auth.NewAuthService(db.sourceReq).CheckBusinessBankCardAccess(id)
	if err != nil {
		return nil, err
	}

	c := BankCard{}

	err = db.Get(&c, "SELECT * FROM business_bank_card WHERE id = $1 AND business_id = $2 AND cardholder_id = $3", id, businessID, userID)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if c.CardStatus == banking.CardStatusBlocked {
		blocks, err := db.GetCardBlocks(c.Id)
		if err != nil {
			return nil, err
		}

		c.CardBlock = blocks
	}

	return &c, nil
}

func (db *bankCardDatastore) GetByIDInternal(id string) (*BankCard, error) {
	c := &BankCard{}
	err := db.Get(c, "SELECT * FROM business_bank_card WHERE id = $1", id)
	if err != nil {
		return nil, err
	}

	if c.CardStatus == banking.CardStatusBlocked {
		blocks, err := db.GetCardBlocks(c.Id)
		if err != nil {
			return nil, err
		}

		c.CardBlock = blocks
	}

	return c, nil
}

func (db *bankCardDatastore) GetByBankCardId(bankCardId string, userID shared.UserID) (*BankCard, error) {

	c := BankCard{}

	err := db.Get(&c, "SELECT * FROM business_bank_card WHERE bank_card_id = $1 AND cardholder_id = $2", bankCardId, userID)
	if err != nil {
		log.Println("Error getting bank card", bankCardId, userID)
		return nil, err
	}

	if c.CardStatus == banking.CardStatusBlocked {
		blocks, err := db.GetCardBlocks(c.Id)
		if err != nil {
			return nil, err
		}

		c.CardBlock = blocks
	}

	return &c, nil
}

func (db *bankCardDatastore) GetByBusinessID(offset int, limit int, businessID shared.BusinessID, userID shared.UserID) ([]BankCard, error) {
	// Check access
	if userID != db.sourceReq.UserID {
		return nil, errors.New("unauthorized")
	}

	rows := []BankCard{}

	err := db.Select(&rows, "SELECT * FROM business_bank_card WHERE business_id = $1 AND cardholder_id = $2", businessID, userID)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return rows, err
}

func (db *bankCardDatastore) GetByAccountId(offset int, limit int, accountId string, businessID shared.BusinessID, userID shared.UserID) ([]BankCard, error) {
	// Check access
	if userID != db.sourceReq.UserID {
		return nil, errors.New("unauthorized")
	}

	rows := []BankCard{}

	err := db.Select(&rows, "SELECT * FROM business_bank_card WHERE bank_account_id = $1 AND business_id = $2 AND cardholder_id = $3", accountId, businessID, userID)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return rows, err
}

func (db *bankCardDatastore) GetByAccountInternal(offset int, limit int, accountId string, businessID shared.BusinessID) ([]*BankCard, error) {
	rows := []*BankCard{}
	err := db.Select(&rows, "SELECT * FROM business_bank_card WHERE bank_account_id = $1 AND business_id = $2", accountId, businessID)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return rows, err
}

func (db *bankCardDatastore) GetCardBlocks(cardID string) ([]*BankCardBlock, error) {
	rows := []*BankCardBlock{}

	err := db.Select(&rows, "SELECT * FROM business_bank_card_block WHERE card_id = $1 AND block_status = $2", cardID, banking.BlockStatusActive)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return rows, err
}

func (db *bankCardDatastore) CreateBankCard(cardCreate *BankCardCreate) (*BankCard, error) {
	if cardCreate.BankAccountId == "" {
		return nil, errors.New("Account id is required")
	}

	if cardCreate.CardType == "" {
		return nil, errors.New("Card type is required")
	}

	// One debit card per user
	cards, err := db.GetByAccountId(0, 10, cardCreate.BankAccountId, cardCreate.BusinessID, cardCreate.CardholderID)
	if err != nil {
		return nil, err
	}
	if len(cards) > 0 {
		return nil, errors.New("only one debit card allowed per user")
	}

	var b struct {
		FirstName      string            `db:"first_name"`
		LastName       string            `db:"last_name"`
		Phone          string            `db:"phone"`
		MailingAddress *services.Address `db:"mailing_address"`
		ConsumerID     shared.ConsumerID `db:"consumer_id"`
		BankAccountId  string            `db:"bank_account_id"`
		UsageType      UsageType         `db:"usage_type"`
	}

	err = db.Get(
		&b, `
		SELECT
			consumer.first_name, consumer.last_name, wise_user.phone, consumer.mailing_address,
			wise_user.consumer_id, business_bank_account.bank_account_id, business_bank_account.usage_type FROM wise_user
		JOIN business_bank_account ON wise_user.id = business_bank_account.account_holder_id
		JOIN consumer ON wise_user.consumer_id = consumer.id
		WHERE wise_user.id = $1 AND business_bank_account.id = $2`,
		cardCreate.CardholderID,
		cardCreate.BankAccountId,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("business bank account with id: %s not found", cardCreate.BankAccountId)
		}

		return nil, err
	}

	// Only primary accounts can have a card
	if b.UsageType != UsageTypePrimary {
		return nil, fmt.Errorf("cannot create card for business bank account usage type of: %s", b.UsageType)
	}

	address := partnerbank.AddressRequestTypeLegal
	if b.MailingAddress != nil {
		address = partnerbank.AddressRequestTypeMailing
	}

	request := partnerbank.CreateCardRequest{
		AccountID:      partnerbank.AccountBankID(b.BankAccountId),
		CardholderName: b.FirstName + " " + b.LastName,
		Packaging:      partnerbank.CardPackagingRegular,
		Delivery:       partnerbank.CardDeliveryStandard,
		BusinessName:   partnerbank.CardBusinessNameLegal,
		Phone:          partnerbank.PhoneE164(b.Phone),
		Address:        address,
	}

	// Create card with partner bank
	bank, err := partnerbank.GetBusinessBank(partnerbank.ProviderNameBBVA)
	if err != nil {
		return nil, err
	}

	srv, err := bank.CardService(
		db.sourceReq.PartnerBankRequest(),
		partnerbank.BusinessID(cardCreate.BusinessID),
		partnerbank.ConsumerID(b.ConsumerID),
	)
	if err != nil {
		return nil, err
	}

	resp, err := srv.Create(request)
	if err != nil {
		return nil, err
	}

	c := transformCardResponse(resp)

	// Get card limits
	limits, err := srv.GetLimit(partnerbank.CardBankID(c.BankCardId))
	if err != nil {
		return nil, err
	}
	c.DailyTransactionLimit = &limits.DailyTransactionCount
	c.DailyATMLimit = &limits.DailyATMAmount
	c.DailyPOSLimit = &limits.DailyPOSAmount
	c.CardholderID = cardCreate.CardholderID
	c.BusinessID = cardCreate.BusinessID
	c.BankAccountId = cardCreate.BankAccountId

	columns := []string{
		"cardholder_id", "business_id", "bank_account_id", "card_type", "cardholder_name", "is_virtual", "bank_name", "bank_card_id",
		"card_number_masked", "card_brand", "currency", "card_status", "alias", "daily_withdrawal_limit", "daily_pos_limit", "daily_transaction_limit",
	}

	// Default/mandatory values
	values := []string{
		":cardholder_id", ":business_id", ":bank_account_id", ":card_type", ":cardholder_name", ":is_virtual", ":bank_name", ":bank_card_id",
		":card_number_masked", ":card_brand", ":currency", ":card_status", ":alias", ":daily_withdrawal_limit", ":daily_pos_limit", ":daily_transaction_limit",
	}

	sql := fmt.Sprintf("INSERT INTO business_bank_card(%s) VALUES(%s) RETURNING *", strings.Join(columns, ", "), strings.Join(values, ", "))

	stmt, err := db.PrepareNamed(sql)
	if err != nil {
		return nil, err
	}

	card := &BankCard{}

	err = stmt.Get(card, &c)
	if err != nil {
		return nil, err
	} else {
		return card, nil
	}
}

func (db *bankCardDatastore) ActivateBankCard(cardActivate *BankCardActivate) (*BankCard, error) {

	panSuffix := shared.StripSpaces(cardActivate.PANLast6)
	if len(panSuffix) != 6 {
		return nil, errors.New(fmt.Sprintf("Invalid PAN last digits"))
	}

	c, err := db.GetById(cardActivate.Id, cardActivate.BusinessID, cardActivate.CardholderID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New(fmt.Sprintf("Card with id: %s not found", cardActivate.Id))
		}

		return nil, err
	}

	var cid shared.ConsumerID
	err = db.Get(&cid, "SELECT consumer_id FROM wise_user WHERE id = $1", cardActivate.CardholderID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New(fmt.Sprintf("User with id: %s not found", cardActivate.CardholderID))
		}

		return nil, err
	}

	request := partnerbank.ActivateCardRequest{
		CardID:   partnerbank.CardBankID(c.BankCardId),
		PANLast6: panSuffix,
	}

	// Activate card
	bank, err := partnerbank.GetBusinessBank(partnerbank.ProviderNameBBVA)
	if err != nil {
		return nil, err
	}

	srv, err := bank.CardService(
		db.sourceReq.PartnerBankRequest(),
		partnerbank.BusinessID(cardActivate.BusinessID),
		partnerbank.ConsumerID(cid),
	)
	if err != nil {
		return nil, err
	}

	resp, err := srv.Activate(request)
	if err != nil {
		return nil, err
	}

	card := transformCardResponse(resp)
	_, err = db.NamedExec(
		fmt.Sprintf(
			"UPDATE business_bank_card SET card_status = :card_status WHERE id = '%s'",
			cardActivate.Id,
		), card,
	)

	if err != nil {
		return nil, errors.Cause(err)
	}

	return db.GetById(cardActivate.Id, cardActivate.BusinessID, cardActivate.CardholderID)
}

func (db *bankCardDatastore) BlockBankCard(cardBlock *BankCardBlockCreate) (*BankCard, error) {
	switch cardBlock.OriginatedFrom {
	case banking.OriginatedFromClientApplication, banking.OriginatedFromCSP:
		break
	default:
		return nil, errors.New("originated from should be a valid value")
	}

	c, err := db.GetById(cardBlock.CardID, cardBlock.BusinessID, cardBlock.CardholderID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New(fmt.Sprintf("Card with id: %s not found", cardBlock.CardID))
		}

		return nil, err
	}

	if c.CardStatus != banking.CardStatusActive {
		return nil, errors.New("Only active card can be blocked")
	}

	var cid shared.ConsumerID
	err = db.Get(&cid, "SELECT consumer_id FROM wise_user WHERE id = $1", cardBlock.CardholderID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New(fmt.Sprintf("User with id: %s not found", cardBlock.CardholderID))
		}

		return nil, err
	}

	request := partnerbank.CardBlockRequest{
		CardID: partnerbank.CardBankID(c.BankCardId),
		Reason: partnerbank.CardBlockReason(cardBlock.BlockID),
	}

	// Block card
	bank, err := partnerbank.GetBusinessBank(partnerbank.ProviderNameBBVA)
	if err != nil {
		return nil, err
	}

	srv, err := bank.CardService(
		db.sourceReq.PartnerBankRequest(),
		partnerbank.BusinessID(cardBlock.BusinessID),
		partnerbank.ConsumerID(cid),
	)
	if err != nil {
		return nil, err
	}

	resp, err := srv.Block(request)
	if err != nil {
		return nil, err
	}

	columns := []string{
		"card_id", "block_id", "reason", "originated_from", "block_status", "block_date",
	}

	// Default/mandatory values
	values := []string{
		":card_id", ":block_id", ":reason", ":originated_from", ":block_status", ":block_date",
	}

	var isActive bool
	for _, block := range resp {
		b := banking.BankCardBlockCreate{
			BlockID:   banking.CardBlockID(block.Reason),
			CardID:    cardBlock.CardID,
			BlockDate: block.BlockDate,
		}

		if block.IsActive {
			b.BlockStatus = banking.BlockStatusActive
		} else {
			b.BlockStatus = banking.BlockStatusInactive
		}

		if b.BlockID == cardBlock.BlockID {
			b.Reason = cardBlock.Reason
			b.OriginatedFrom = cardBlock.OriginatedFrom
		}

		sql := fmt.Sprintf("INSERT INTO business_bank_card_block(%s) VALUES(%s) RETURNING *", strings.Join(columns, ", "), strings.Join(values, ", "))

		stmt, err := db.PrepareNamed(sql)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		cardBlock := &BankCardBlock{}

		err = stmt.Get(cardBlock, &b)
		if err != nil {
			return nil, err
		}

		if !isActive {
			isActive = block.IsActive
		}
	}

	if isActive {
		card := BankCard{}
		card.CardStatus = banking.CardStatusBlocked

		_, err = db.NamedExec(
			fmt.Sprintf(
				"UPDATE business_bank_card SET card_status = :card_status WHERE id = '%s'",
				cardBlock.CardID,
			), card,
		)
		if err != nil {
			return nil, errors.Cause(err)
		}
	}

	return db.GetById(cardBlock.CardID, cardBlock.BusinessID, cardBlock.CardholderID)
}

func (db *bankCardDatastore) UnBlockBankCard(blockDelete *BankCardBlockDelete) (*BankCard, error) {
	c, err := db.GetById(blockDelete.CardID, blockDelete.BusinessID, blockDelete.CardholderID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New(fmt.Sprintf("Card with id: %s not found", blockDelete.CardID))
		}

		return nil, err
	}

	if c.CardStatus != banking.CardStatusBlocked {
		return nil, errors.New("Only blocked cards can be unblocked")
	}

	var cid shared.ConsumerID
	err = db.Get(&cid, "SELECT consumer_id FROM wise_user WHERE id = $1", blockDelete.CardholderID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New(fmt.Sprintf("User with id: %s not found", blockDelete.CardholderID))
		}

		return nil, err
	}

	block, err := db.GetByBlockID(blockDelete.ID)
	if err != nil {
		return nil, err
	}

	request := partnerbank.CardUnblockRequest{
		CardID: partnerbank.CardBankID(c.BankCardId),
		Reason: partnerbank.CardBlockReason(block.BlockID),
	}

	bank, err := partnerbank.GetBusinessBank(partnerbank.ProviderNameBBVA)
	if err != nil {
		return nil, err
	}

	srv, err := bank.CardService(
		db.sourceReq.PartnerBankRequest(),
		partnerbank.BusinessID(blockDelete.BusinessID),
		partnerbank.ConsumerID(cid),
	)
	if err != nil {
		return nil, err
	}

	err = srv.Unblock(request)
	if err != nil {
		return nil, err
	}

	// Update card block
	_, err = db.Exec("UPDATE business_bank_card_block SET block_status = $1 WHERE id = $2", banking.BlockStatusInactive, blockDelete.ID)
	if err != nil {
		return nil, err
	}

	// Unblock card
	card := BankCard{}
	card.CardStatus = banking.CardStatusActive

	_, err = db.NamedExec(
		fmt.Sprintf(
			"UPDATE business_bank_card SET card_status = :card_status WHERE id = '%s'",
			c.Id,
		), card,
	)
	if err != nil {
		return nil, errors.Cause(err)
	}

	return db.GetById(blockDelete.CardID, blockDelete.BusinessID, blockDelete.CardholderID)
}

func (db *bankCardDatastore) GetByBlockID(ID shared.BankCardBlockID) (*BankCardBlock, error) {
	block := BankCardBlock{}

	err := db.Get(&block, "SELECT * FROM business_bank_card_block WHERE id = $1", ID)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return &block, err
}

func (db *bankCardDatastore) GetByPartnerBlockID(blockID banking.CardBlockID, cardID string) (*BankCardBlock, error) {
	block := BankCardBlock{}

	query := `SELECT * FROM business_bank_card_block WHERE block_id = $1 AND card_id = $2 AND block_status = $3`

	err := db.Get(&block, query, blockID, cardID, banking.BlockStatusActive)
	if err != nil {
		return nil, err
	}

	return &block, err
}

func (ds bankCardDatastore) GetByBlockIDAndDate(blockID string, blockDate time.Time) (*BankCardBlock, error) {
	block := BankCardBlock{}

	err := ds.Get(&block, "SELECT * FROM business_bank_card_block WHERE block_id = $1 AND block_date = $2", blockID, blockDate)
	if err != nil {
		return nil, err
	}

	return &block, err
}

// Transform partner bank layer response
func transformCardResponse(response *partnerbank.GetCardResponse) BankCard {

	c := BankCard{}

	c.BankName = "bbva"
	c.BankCardId = string(response.CardID)
	c.CardNumberMasked = response.PANMasked
	c.CardBrand = string(response.Brand)
	c.Currency = banking.Currency(response.Currency)
	c.CardStatus = banking.CardStatus(response.Status)
	c.CardholderName = response.CardholderName
	c.CardType = banking.CardType(response.Type)

	return c
}

func (db *bankCardDatastore) UpdateCardStatus(cardID string, userID shared.UserID, status string) error {
	_, err := db.Exec("UPDATE business_bank_card SET card_status = $1 WHERE id = $2 AND cardholder_id = $3", status, cardID, userID)
	if err != nil {
		log.Println("Error updating card status", userID, status)
		return err
	}

	return nil
}

func (ds bankCardDatastore) CheckBlockStatus(cardID string) (*BankCard, error) {
	card, err := ds.getCardByID(cardID)
	if err != nil {
		return nil, err
	}

	var cid shared.ConsumerID
	err = ds.Get(&cid, "SELECT consumer_id FROM wise_user WHERE id = $1", card.CardholderID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New(fmt.Sprintf("User with id: %s not found", card.CardholderID))
		}

		return nil, err
	}

	bank, err := partnerbank.GetBusinessBank(partnerbank.ProviderNameBBVA)
	if err != nil {
		return nil, err
	}

	srv, err := bank.CardService(
		ds.sourceReq.PartnerBankRequest(),
		partnerbank.BusinessID(card.BusinessID),
		partnerbank.ConsumerID(cid),
	)
	if err != nil {
		return nil, err
	}

	blocks, err := srv.GetAllBlocks(partnerbank.CardBankID(card.BankCardId))
	if err != nil {
		log.Println(err)
		return nil, err
	}

	blocked := false

	// Deactivate all blocks
	err = ds.deactivateAllCardBlock(cardID)
	if err != nil {
		return nil, err
	}

	// Iterate each block
	for _, b := range blocks {
		block, err := ds.GetByBlockIDAndDate(string(b.Reason), b.BlockDate)
		if err != nil && err != sql.ErrNoRows {
			log.Println(err)
			return nil, err
		}

		if b.IsActive {
			blocked = true

			if block == nil {
				create := BankCardBlockCreate{
					BusinessID: card.BusinessID,
				}

				create.CardID = card.Id
				create.CardholderID = card.CardholderID
				create.BlockID = banking.CardBlockID(b.Reason)
				create.OriginatedFrom = banking.OriginatedFromBank

				reason := ""
				create.Reason = &reason
				create.BlockDate = b.BlockDate
				create.BlockStatus = banking.BlockStatusActive

				_, err := ds.addCardBlock(create)
				if err != nil {
					log.Println(err)
					return nil, err
				}
			} else if block.BlockStatus == banking.BlockStatusInactive {
				query := `UPDATE business_bank_card_block SET block_status = $1 WHERE id = $2`

				_, err := ds.Exec(query, banking.BlockStatusActive, block.ID)
				if err != nil {
					log.Println(err)
					return nil, err
				}
			}
		}
	}

	bankCard := BankCard{}
	if blocked {
		bankCard.CardStatus = banking.CardStatusBlocked
	} else {
		bankCard.CardStatus = banking.CardStatusActive
	}

	_, err = ds.NamedExec(
		fmt.Sprintf(
			"UPDATE business_bank_card SET card_status = :card_status WHERE id = '%s'",
			card.Id,
		), bankCard,
	)

	if err != nil {
		return nil, err
	}

	sr := services.NewSourceRequest()
	sr.UserID = card.CardholderID
	return NewCardService(sr).GetById(card.Id, card.BusinessID, card.CardholderID)
}

func (ds bankCardDatastore) getCardByID(ID string) (*BankCard, error) {
	c := BankCard{}

	err := ds.Get(&c, "SELECT * FROM business_bank_card WHERE id = $1", ID)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (ds bankCardDatastore) getActiveBlockID(cardID string) (*BankCardBlock, error) {
	c := BankCardBlock{}

	err := ds.Get(&c, "SELECT * FROM business_bank_card_block WHERE card_id = $1 AND block_status = 'active'", cardID)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (ds bankCardDatastore) addCardBlock(c BankCardBlockCreate) (*BankCardBlock, error) {

	columns := []string{
		"card_id", "block_id", "reason", "originated_from", "block_status", "block_date",
	}

	// Default/mandatory values
	values := []string{
		":card_id", ":block_id", ":reason", ":originated_from", ":block_status", ":block_date",
	}

	sql := fmt.Sprintf("INSERT INTO business_bank_card_block(%s) VALUES(%s) RETURNING *", strings.Join(columns, ", "), strings.Join(values, ", "))

	stmt, err := ds.PrepareNamed(sql)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	cardBlock := &BankCardBlock{}

	err = stmt.Get(cardBlock, &c)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return cardBlock, nil
}

func (ds bankCardDatastore) deactivateAllCardBlock(cardID string) error {
	block := `UPDATE business_bank_card_block SET block_status = $1 WHERE card_id = $2`

	_, err := ds.Exec(block, banking.BlockStatusInactive, cardID)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (db *bankCardDatastore) CancelCardInternal(cardID string) (*BankCard, error) {
	card, err := db.GetByIDInternal(cardID)
	if err != nil {
		log.Println("Error fetching card:", cardID)
		return nil, err
	}

	var cID shared.ConsumerID
	err = db.Get(&cID, "SELECT wise_user.consumer_id FROM wise_user WHERE wise_user.id = $1", card.CardholderID)
	if err != nil {
		return nil, err
	}

	bank, err := partnerbank.GetBusinessBank(partnerbank.ProviderNameBBVA)
	if err != nil {
		return nil, err
	}

	srv, err := bank.CardService(
		db.sourceReq.PartnerBankRequest(),
		partnerbank.BusinessID(card.BusinessID),
		partnerbank.ConsumerID(cID),
	)

	err = srv.CancelInternal(partnerbank.CardBankID(card.BankCardId))
	if err != nil {
		return nil, err
	}

	// Update Status
	err = db.UpdateCardStatus(card.Id, card.CardholderID, string(banking.CardStatusCanceled))
	if err != nil {
		return nil, err
	}

	return db.GetByIDInternal(cardID)
}
