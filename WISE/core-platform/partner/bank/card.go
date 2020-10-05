package bank

import (
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/ttacon/libphonenumber"
)

type CardBankID string

func (id CardBankID) String() string {
	return string(id)
}

type CardType string

const (
	CardTypeDebit = CardType("debit")
)

type CardBrand string

const (
	CardBrandVisa       = CardBrand("visa")
	CardBrandMastercard = CardBrand("mc")
	CardBrandAmex       = CardBrand("amex")
)

func isVisa(cardNumber string) bool {
	return strings.HasPrefix(cardNumber, "4")
}

func isMasterCard(cardNumber string) bool {
	switch cardNumber[:1] {
	case "50", "51", "52", "53", "54", "55":
		return true
	}

	return false
}

func isAmex(cardNumber string) bool {
	switch cardNumber[:1] {
	case "34", "37":
		return true
	}

	return false
}

const MinCardNumberLen = 12
const MaxCardNumberLen = 19

func CardBrandFromNumber(cardNumber string) (CardBrand, error) {
	if len(cardNumber) < 12 || len(cardNumber) > 19 {
		return CardBrand(""), errors.New("invalid card number")
	}

	if isVisa(cardNumber) {
		return CardBrandVisa, nil
	} else if isMasterCard(cardNumber) {
		return CardBrandMastercard, nil
	} else if isAmex(cardNumber) {
		return CardBrandAmex, nil
	}

	return CardBrand(""), errors.New("card brand unknown")
}

type CardDelivery string

const (
	CardDeliveryStandard  = CardDelivery("standard")
	CardDeliveryExpedited = CardDelivery("expedited")
)

type CardPackaging string

const (
	CardPackagingRegular = CardPackaging("regular")
)

type CardBusinessName string

const (
	CardBusinessNameLegal = CardBusinessName("legal")
	CardBusinessNameDBA   = CardBusinessName("dba")
)

type PhoneE164 string

func (p PhoneE164) String() (string, error) {
	ph, err := libphonenumber.Parse(string(p), "US")
	if err != nil {
		return "", err
	}

	return libphonenumber.Format(ph, libphonenumber.E164), nil
}

func (p PhoneE164) USNational() (string, error) {
	ph, err := libphonenumber.Parse(string(p), "US")
	if err != nil {
		return "", err
	}

	num := libphonenumber.Format(ph, libphonenumber.E164)
	if len(num) != 12 {
		return "", fmt.Errorf("unexpected number format %s", num)
	}

	return num[2:12], nil
}

type CreateCardRequest struct {
	AccountID      AccountBankID      `json:"accountId"`
	Type           CardType           `json:"cardType"`
	CardholderName string             `json:"cardholderName"`
	BusinessName   CardBusinessName   `json:"businessName"`
	Delivery       CardDelivery       `json:"delivery"`
	Packaging      CardPackaging      `json:"packaging"`
	Phone          PhoneE164          `json:"phone"`
	Address        AddressRequestType `json:"address"`
}

type CardStatus string

const (
	CardStatusPrinting = CardStatus("printing")
	CardStatusShipped  = CardStatus("shipped")
	CardStatusActive   = CardStatus("active")
	CardStatusInactive = CardStatus("inactive")
	CardStatusBlocked  = CardStatus("blocked")
	CardStatusCanceled = CardStatus("canceled")
)

type GetCardLimitResponse struct {
	DailyATMAmount        float64 `json:"dailyATMAmount"`
	DailyPOSAmount        float64 `json:"dailyPOSAmount"`
	DailyTransactionCount int     `json:"dailyTransactionCount"`
}

// TODO: Use vault tokenizer for credit card handling
type GetCardResponse struct {
	CardID         CardBankID            `json:"cardId"`
	PANMasked      string                `json:"panMasked"`
	PANAlias       string                `json:"panAlias"`
	Type           CardType              `json:"cardType"`
	Brand          CardBrand             `json:"cardBrand"`
	Currency       Currency              `json:"currency"`
	Status         CardStatus            `json:"status"`
	CardholderName string                `json:"cardholderName"`
	CardLimit      *GetCardLimitResponse `json:"cardLimit"`
}

type CardReissueReason string

const (
	CardReissueReasonNameChange  = CardReissueReason("nameChange")
	CardReissueReasonDamaged     = CardReissueReason("damaged")
	CardReissueReasonLost        = CardReissueReason("lost")
	CardReissueReasonNotReceived = CardReissueReason("notReceived")
	CardReissueReasonNotWorking  = CardReissueReason("notWorking")
	CardReissueReasonStolen      = CardReissueReason("stolen")
	CardReissueReasonUpgrade     = CardReissueReason("upgrade")
	CardReissueReasonResendPin   = CardReissueReason("resendPin")
	CardReissueReasonOther       = CardReissueReason("other")
)

type ReissueCardRequest struct {
	CardID         CardBankID         `json:"cardId"`
	CardholderName string             `json:"cardholderName"`
	BusinessName   CardBusinessName   `json:"businessName"`
	Reason         CardReissueReason  `json:""`
	Delivery       CardDelivery       `json:"delivery"`
	Packaging      CardPackaging      `json:"packaging"`
	Phone          PhoneE164          `json:"phone"`
	Address        AddressRequestType `json:"address"`
}

type ActivateCardRequest struct {
	CardID   CardBankID `json:"cardId"`
	PANLast6 string     `json:"panLast6"`
}

type SetCardPINRequest struct {
	CardID   CardBankID `json:"cardId"`
	PANLast6 string     `json:"panLast6"`
	PIN      string     `json:"PIN"`
}

type CancelCardRequest struct {
	CardID   CardBankID `json:"cardId"`
	PANLast6 string     `json:"panLast6"`
}

type CardBlockReason string

const (
	CardBlockReasonLost      = CardBlockReason("lost")
	CardBlockReasonStolen    = CardBlockReason("stolen")
	CardBlockReasonDispute   = CardBlockReason("dispute")
	CardBlockReasonFraud     = CardBlockReason("fraud")
	CardBlockReasonInternal  = CardBlockReason("internal")
	CardBlockReasonLocked    = CardBlockReason("locked")
	CardBlockReasonChargeOff = CardBlockReason("chargeOff")
)

type CardBlockRequest struct {
	CardID CardBankID      `json:"cardId"`
	Reason CardBlockReason `json:"reason"`
}

type CardUnblockRequest struct {
	CardID CardBankID      `json:"cardId"`
	Reason CardBlockReason `json:"reason"`
}

type CardBlockResponse struct {
	CardID    CardBankID      `json:"cardId"`
	Reason    CardBlockReason `json:"reason"`
	BlockDate time.Time       `json:"blockDate"`
	IsActive  bool            `json:"isActive"`
}
