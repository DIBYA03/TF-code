package bbva

import (
	"errors"
	"fmt"
	"time"

	"github.com/wiseco/core-platform/partner/bank"
	"github.com/wiseco/core-platform/services"
)

type CardType string

const (
	CardTypeDebit = CardType("debit")
)

var partnerCardTypeFrom = map[bank.CardType]CardType{
	bank.CardTypeDebit: CardTypeDebit,
}

var partnerCardTypeTo = map[CardType]bank.CardType{
	CardTypeDebit: bank.CardTypeDebit,
}

type CardDelivery string

const (
	CardDeliveryStandard         = CardDelivery("standard")
	CardDeliveryPriority         = CardDelivery("priority")
	CardDeliveryPrioritySaturday = CardDelivery("priority_saturday")
)

var partnerCardDeliveryFrom = map[bank.CardDelivery]CardDelivery{
	bank.CardDeliveryStandard:  CardDeliveryStandard,
	bank.CardDeliveryExpedited: CardDeliveryPriority,
}

var partnerCardDeliveryTo = map[CardDelivery]bank.CardDelivery{
	CardDeliveryStandard: bank.CardDeliveryStandard,
	CardDeliveryPriority: bank.CardDeliveryExpedited,
}

type CardPackaging string

const (
	CardPackagingRegular = CardPackaging("regular")
)

var partnerCardPackagingFrom = map[bank.CardPackaging]CardPackaging{
	bank.CardPackagingRegular: CardPackagingRegular,
}

var partnerCardPackagingTo = map[CardPackaging]bank.CardPackaging{
	CardPackagingRegular: bank.CardPackagingRegular,
}

type CardBusinessName string

const (
	CardBusinessNameLegal = CardBusinessName("legal")
	CardBusinessNameDBA   = CardBusinessName("dba")
)

var partnerCardBusinessNameFrom = map[bank.CardBusinessName]CardBusinessName{
	bank.CardBusinessNameLegal: CardBusinessNameLegal,
	bank.CardBusinessNameDBA:   CardBusinessNameDBA,
}

var partnerCardBusinessNameTo = map[CardBusinessName]bank.CardBusinessName{
	CardBusinessNameLegal: bank.CardBusinessNameLegal,
	CardBusinessNameDBA:   bank.CardBusinessNameDBA,
}

type CreateCardRequest struct {
	AccountID      string           `json:"account_id"`
	CardType       CardType         `json:"card_type"`
	CardholderName string           `json:"cardholder_name"`
	BusinessName   CardBusinessName `json:"business_name_on_card"`
	Delivery       CardDelivery     `json:"delivery,omitempty"`
	Packaging      CardPackaging    `json:"packaging"`
	AddressID      string           `json:"address_id"`
	PhoneNumber    string           `json:"phone_number,omitempty"`
}

type CreateCardResponse struct {
	CardID     string `json:"card_id"`     // Example: "DC-106206e8-0ef5-4563-a2c8-e50a723c19bd"
	CardNumber string `json:"card_number"` // Example: "4123506696817757"
}

type CardStatus string

const (
	CardStatusActive           = CardStatus("active")
	CardStatusInactive         = CardStatus("inactive")
	CardStatusBlocked          = CardStatus("blocked")
	CardStatusCanceled         = CardStatus("canceled")
	CardStatusPendingEmbossing = CardStatus("pending_embossing")
	CardStatusPendingDelivery  = CardStatus("pending_delivery")
)

var partnerCardStatusTo = map[CardStatus]bank.CardStatus{
	CardStatusActive:           bank.CardStatusActive,
	CardStatusInactive:         bank.CardStatusInactive,
	CardStatusBlocked:          bank.CardStatusBlocked,
	CardStatusCanceled:         bank.CardStatusCanceled,
	CardStatusPendingEmbossing: bank.CardStatusPrinting,
	CardStatusPendingDelivery:  bank.CardStatusShipped,
}

func (s CardStatus) IsActive() bool {
	return s == CardStatusActive
}

type GetCardResponse struct {
	CardID     string     `json:"card_id"`
	CardNumber string     `json:"card_number"`
	Cardholder string     `json:"card_holder"`
	Type       CardType   `json:"type"`
	Status     CardStatus `json:"card_status"`
	Currency   Currency   `json:"card_currency"`
}

type GetCardLimitResponse struct {
	ATMDaily          float64 `json:"atm_daily"`
	POSDaily          float64 `json:"pos_daily"`
	DailyTransactions int     `json:"daily_transactions"`
}

func (r *GetCardResponse) partnerCardResponseTo(limitResp *bank.GetCardLimitResponse) (*bank.GetCardResponse, error) {
	brand, err := bank.CardBrandFromNumber(r.CardNumber)
	if err != nil {
		return nil, err
	}

	cardType, ok := partnerCardTypeTo[r.Type]
	if !ok {
		return nil, errors.New(fmt.Sprintf("invalid card type (%s)", r.Type))
	}

	currency, ok := partnerCurrencyTo[r.Currency]
	if !ok {
		return nil, errors.New(fmt.Sprintf("invalid currency (%s)", r.Currency))
	}

	status, ok := partnerCardStatusTo[r.Status]
	if !ok {
		return nil, errors.New(fmt.Sprintf("invalid card status (%s)", r.Status))
	}
	// Mask PAN
	panMasked := services.MaskLeft(r.CardNumber, 4)

	return &bank.GetCardResponse{
		CardID:         bank.CardBankID(r.CardID),
		PANMasked:      panMasked,
		PANAlias:       r.CardNumber,
		Type:           cardType,
		Brand:          brand,
		Currency:       currency,
		Status:         status,
		CardholderName: r.Cardholder,
		CardLimit:      limitResp,
	}, nil
}

func getCard(cardId string, cards *[]bank.GetCardResponse, cardNumberAlias string) *bank.GetCardResponse {
	for _, c := range *cards {
		if string(c.CardID) == cardId {
			c.PANAlias = cardNumberAlias
			return &c
		}
	}
	return nil
}

type Pages struct {
	TotalItems  int    `json:"total_items"`
	NextPageKey string `json:"next_page_key"`
	HasMore     string `json:"has_more"`
}

type GetAllCardsResponse struct {
	Cards    []GetCardResponse `json:"cards"`
	PageData Pages             `json:"pages"`
}

type CardReissueReason string

const (
	CardReissueReasonNameChange  = CardReissueReason("name_change")
	CardReissueReasonDamaged     = CardReissueReason("damaged_card")
	CardReissueReasonNotReceived = CardReissueReason("not_received")
	CardReissueReasonNotWorking  = CardReissueReason("atm_or_pos_not_working")
	CardReissueReasonUpgrade     = CardReissueReason("upgrade")
	CardReissueReasonRemoveImage = CardReissueReason("remove_image")
	CardReissueReasonResendPin   = CardReissueReason("resend_pin")
)

var partnerCardReissueReasonFrom = map[bank.CardReissueReason]CardReissueReason{
	bank.CardReissueReasonNameChange:  CardReissueReasonNameChange,
	bank.CardReissueReasonDamaged:     CardReissueReasonDamaged,
	bank.CardReissueReasonLost:        CardReissueReasonNotReceived,
	bank.CardReissueReasonNotReceived: CardReissueReasonNotReceived,
	bank.CardReissueReasonResendPin:   CardReissueReasonResendPin,
	bank.CardReissueReasonNotWorking:  CardReissueReasonNotWorking,
	bank.CardReissueReasonStolen:      CardReissueReasonNotReceived,
	bank.CardReissueReasonUpgrade:     CardReissueReasonUpgrade,
	bank.CardReissueReasonOther:       CardReissueReasonUpgrade,
}

var partnerCardReissueReasonTo = map[CardReissueReason]bank.CardReissueReason{
	CardReissueReasonNameChange:  bank.CardReissueReasonNameChange,
	CardReissueReasonDamaged:     bank.CardReissueReasonDamaged,
	CardReissueReasonNotReceived: bank.CardReissueReasonNotReceived,
	CardReissueReasonResendPin:   bank.CardReissueReasonResendPin,
	CardReissueReasonNotWorking:  bank.CardReissueReasonNotWorking,
	CardReissueReasonUpgrade:     bank.CardReissueReasonUpgrade,
	CardReissueReasonRemoveImage: bank.CardReissueReasonOther,
}

type ReissueCardRequest struct {
	CardholderName string            `json:"name,omitempty"`
	Reason         CardReissueReason `json:"reason"`
	BusinessName   CardBusinessName  `json:"business_name_on_card,omitempty"`
	Delivery       CardDelivery      `json:"delivery"`
	Packaging      CardPackaging     `json:"packaging"`
	AddressID      string            `json:"address_id"`
	PhoneNumber    string            `json:"phone_number,omitempty"`
}

type SetCardPINRequest struct {
	PIN string `json:"new_pin"`
}

type CardBlockReason string

const (
	CardBlockReasonLost      = CardBlockReason("lost")
	CardBlockReasonStolen    = CardBlockReason("stolen")
	CardBlockReasonFraud     = CardBlockReason("fraud")
	CardBlockReasonSecurity  = CardBlockReason("security")
	CardBlockReasonInternal  = CardBlockReason("internal")
	CardBlockReasonTempBlock = CardBlockReason("temporary_block")
	CardBlockReasonBadDebt   = CardBlockReason("bad_debt")
)

var partnerCardBlockReasonFrom = map[bank.CardBlockReason]CardBlockReason{
	bank.CardBlockReasonLost:      CardBlockReasonLost,
	bank.CardBlockReasonStolen:    CardBlockReasonStolen,
	bank.CardBlockReasonDispute:   CardBlockReasonFraud,
	bank.CardBlockReasonFraud:     CardBlockReasonFraud,
	bank.CardBlockReasonInternal:  CardBlockReasonInternal,
	bank.CardBlockReasonLocked:    CardBlockReasonTempBlock,
	bank.CardBlockReasonChargeOff: CardBlockReasonBadDebt,
}

var partnerCardBlockReasonTo = map[CardBlockReason]bank.CardBlockReason{
	CardBlockReasonLost:      bank.CardBlockReasonLost,
	CardBlockReasonStolen:    bank.CardBlockReasonStolen,
	CardBlockReasonFraud:     bank.CardBlockReasonFraud,
	CardBlockReasonSecurity:  bank.CardBlockReasonFraud,
	CardBlockReasonInternal:  bank.CardBlockReasonInternal,
	CardBlockReasonTempBlock: bank.CardBlockReasonLocked,
	CardBlockReasonBadDebt:   bank.CardBlockReasonChargeOff,
}

type CardBlockRequest struct {
	Reason CardBlockReason `json:"block_id"`
}

type GetBlockResponse struct {
	CardID    string          `json:"card_id"`
	Reason    CardBlockReason `json:"block_id"`
	BlockDate time.Time       `json:"block_date"`
	IsActive  bool            `json:"is_active"`
}

type GetAllBlocksResponse struct {
	CardBlocks []GetBlockResponse `json:"blocks"`
}
