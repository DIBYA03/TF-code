package support

import (
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/wiseco/core-platform/services"
	"github.com/wiseco/core-platform/services/banking/business"

	"github.com/wiseco/core-platform/services/data"
)

type CardSupportService interface {
	Block(*business.BankCardBlockCreate) (*business.BankCardPartial, error)
	UnBlock(string) (*business.BankCardPartial, error)
	ListBlocks(string) ([]*business.BankCardBlock, error)

	CheckBlockStatus(string) (*business.BankCardPartial, error)
}

type cardSupportDataStore struct {
	*sqlx.DB
	sourceReq services.SourceRequest
}

func NewCardSupportService(sourceReq services.SourceRequest) CardSupportService {
	return cardSupportDataStore{data.DBWrite, sourceReq}
}

func (ds cardSupportDataStore) Block(c *business.BankCardBlockCreate) (*business.BankCardPartial, error) {
	card, err := ds.getCardByID(c.CardID)
	if err != nil {
		return nil, err
	}

	c.CardholderID = card.CardholderID
	c.BusinessID = card.BusinessID

	sr := services.NewSourceRequest()
	sr.UserID = card.CardholderID

	card, err = business.NewCardService(sr).BlockBankCard(c)
	if err != nil {
		return nil, err
	}

	return &card.BankCardPartial, nil
}

func (ds cardSupportDataStore) ListBlocks(cardID string) ([]*business.BankCardBlock, error) {
	blocks, err := business.NewCardService(ds.sourceReq).GetCardBlocks(cardID)
	if err != nil {
		return nil, err
	}

	return blocks, nil
}

func (ds cardSupportDataStore) UnBlock(cardID string) (*business.BankCardPartial, error) {
	card, err := ds.getCardByID(cardID)
	if err != nil {
		return nil, err
	}

	block, err := ds.getActiveBlockID(cardID)
	if err != nil {
		return nil, err
	}

	del := business.BankCardBlockDelete{
		BusinessID: card.BusinessID,
	}
	del.ID = block.ID
	del.CardholderID = card.CardholderID
	del.CardID = card.Id

	sr := services.NewSourceRequest()
	sr.UserID = card.CardholderID

	card, err = business.NewCardService(sr).UnBlockBankCard(&del)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &card.BankCardPartial, nil
}

func (ds cardSupportDataStore) CheckBlockStatus(cardID string) (*business.BankCardPartial, error) {
	card, err := business.NewCardService(ds.sourceReq).CheckBlockStatus(cardID)
	if err != nil {
		return nil, err
	}
	return &card.BankCardPartial, nil
}

func (ds cardSupportDataStore) getCardByID(ID string) (*business.BankCard, error) {
	c := business.BankCard{}

	err := ds.Get(&c, "SELECT * FROM business_bank_card WHERE id = $1", ID)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (ds cardSupportDataStore) getActiveBlockID(cardID string) (*business.BankCardBlock, error) {
	c := business.BankCardBlock{}

	err := ds.Get(&c, "SELECT * FROM business_bank_card_block WHERE card_id = $1 AND block_status = 'active'", cardID)
	if err != nil {
		return nil, err
	}

	return &c, nil
}
