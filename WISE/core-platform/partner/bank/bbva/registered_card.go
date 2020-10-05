package bbva

import (
	"errors"
	"time"

	"github.com/wiseco/core-platform/partner/bank"
)

type RegisterCardType string

const (
	RegisterCardTypeDebit        = "debit"
	RegisterCardTypePrepaidDebit = "reloadable_prepaid"
	RegisterCardTypeUnknown      = "unknown"
)

var partnerLinkedCardTypeFrom = map[bank.LinkedCardType]RegisterCardType{
	bank.LinkedCardTypeDebit:        RegisterCardTypeDebit,
	bank.LinkedCardTypePrepaidDebit: RegisterCardTypePrepaidDebit,
	bank.LinkedCardTypeUnknown:      RegisterCardTypeUnknown,
}

var partnerLinkedCardTypeTo = map[RegisterCardType]bank.LinkedCardType{
	RegisterCardTypeDebit:        bank.LinkedCardTypeDebit,
	RegisterCardTypePrepaidDebit: bank.LinkedCardTypePrepaidDebit,
	RegisterCardTypeUnknown:      bank.LinkedCardTypeUnknown,
}

//TODO when we move all of this to the partner/client service remove wise specific types. the front end should work with the types we get back from bbva
func GetWiseLinkedCardType(rt RegisterCardType) bank.LinkedCardType {
	lct, ok := partnerLinkedCardTypeTo[rt]
	if !ok {
		lct = bank.LinkedCardTypeUnknown
	}

	return lct
}

type RegisterCardBrand string

const (
	RegisterCardBrandVisa       = "visa"
	RegisterCardBrandMastercard = "mastercard"
)

var partnerLinkedCardBrandFrom = map[bank.LinkedCardBrand]RegisterCardBrand{
	bank.LinkedCardBrandVisa:       RegisterCardBrandVisa,
	bank.LinkedCardBrandMastercard: RegisterCardBrandMastercard,
}

var partnerLinkedCardBrandTo = map[RegisterCardBrand]bank.LinkedCardBrand{
	RegisterCardBrandVisa:       bank.LinkedCardBrandVisa,
	RegisterCardBrandMastercard: bank.LinkedCardBrandMastercard,
}

type RegisteredCardUsage string

func (u RegisteredCardUsage) String() string {
	return string(u)
}

const (
	RegisteredCardUsageSendOnly       = RegisteredCardUsage("send_only")
	RegisteredCardUsageRecieveOnly    = RegisteredCardUsage("receive_only")
	RegisteredCardUsageSendAndRecieve = RegisteredCardUsage("send_or_receive")
)

var partnerLinkedCardPermissionFrom = map[bank.LinkedCardPermission]RegisteredCardUsage{
	bank.LinkedCardPermissionSendOnly:       RegisteredCardUsageSendOnly,
	bank.LinkedCardPermissionRecieveOnly:    RegisteredCardUsageRecieveOnly,
	bank.LinkedCardPermissionSendAndRecieve: RegisteredCardUsageSendAndRecieve,
}

var partnerLinkedCardPermissionTo = map[RegisteredCardUsage]bank.LinkedCardPermission{
	RegisteredCardUsageSendOnly:       bank.LinkedCardPermissionSendOnly,
	RegisteredCardUsageRecieveOnly:    bank.LinkedCardPermissionRecieveOnly,
	RegisteredCardUsageSendAndRecieve: bank.LinkedCardPermissionSendAndRecieve,
}

type RegisterCardRequest struct {
	PrimaryAccountNumber string              `json:"primary_account_number"`
	ExpirationDate       string              `json:"expiration_date"`
	CVVCode              string              `json:"cvv2_cvc_code"`
	NameOnAccount        string              `json:"name_on_account"`
	Nickname             string              `json:"nickname"`
	Usage                RegisteredCardUsage `json:"usage"`
	BillingAddress       AddressRequest      `json:"billing_address"`
}

type RegisterCardResponse struct {
	AccountReferenceID string            `json:"account_reference_id"`
	AccountLast4       string            `json:"account_last4"`
	RegisterCardBrand  RegisterCardBrand `json:"card_brand"`
	IssuerName         string            `json:"issuer_name"`
	RegisterCardType   RegisterCardType  `json:"card_type"`
	FastFundsEnabled   string            `json:"fast_funds_enabled"`
}

type GetRegisterCardsResponse struct {
	RegisteredCards []GetRegisterCardResponse `json:"registered_cards"`
}

type RegisterCardStatus string

const (
	RegisterCardStatusActive  = RegisterCardStatus("active")
	RegisterCardStatusDeleted = RegisterCardStatus("deleted")
)

func (s RegisterCardStatus) IsActive() bool {
	return s == RegisterCardStatusActive
}

type FastFundsEnabledType string

const (
	FastFundsEnabledTypeYes     = "yes"
	FastFundsEnabledTypeNo      = "no"
	FastFundsEnabledTypeUnknown = "unknown"
)

type GetRegisterCardResponse struct {
	AccountReferenceID   string              `json:"account_reference_id"`
	PrimaryAccountNumber string              `json:"primary_account_number"`
	ExpirationDate       string              `json:"expiration_date"`
	RegisterCardBrand    RegisterCardBrand   `json:"card_brand"`
	IssuerName           string              `json:"issuer_name"`
	RegisterCardType     RegisterCardType    `json:"card_type"`
	Currency             Currency            `json:"currency"`
	NameOnAccount        string              `json:"name_on_account"`
	Nickname             string              `json:"nickname"`
	BillingAddress       AddressResponse     `json:"billing_address"`
	FastFundsEnabled     string              `json:"fast_funds_enabled"`
	Usage                RegisteredCardUsage `json:"usage"`
	Status               RegisterCardStatus  `json:"status"`
	StatusDate           time.Time           `json:"status_date"`
}

func (a *GetRegisterCardResponse) partnerLinkedCardResponseTo() (*bank.LinkedCardResponse, error) {
	cardType, ok := partnerLinkedCardTypeTo[a.RegisterCardType]
	if !ok {
		return nil, errors.New("invalid card type")
	}

	cardBrand, ok := partnerLinkedCardBrandTo[a.RegisterCardBrand]
	if !ok {
		return nil, errors.New("invalid card brand")
	}

	currency, ok := partnerCurrencyTo[a.Currency]
	if !ok {
		return nil, errors.New("invalid currency")
	}

	perm, ok := partnerLinkedCardPermissionTo[a.Usage]
	if !ok {
		return nil, errors.New("invalid permission")
	}

	at := AddressTypeBilling
	return &bank.LinkedCardResponse{
		CardID:           bank.LinkedCardBankID(a.AccountReferenceID),
		CardType:         cardType,
		BrandType:        cardBrand,
		IssuingBank:      a.IssuerName,
		AccountHolder:    a.NameOnAccount,
		Alias:            a.Nickname,
		CardNumber:       a.PrimaryAccountNumber,
		Expiration:       a.ExpirationDate,
		Currency:         currency,
		FastFundsEnabled: a.FastFundsEnabled == FastFundsEnabledTypeYes,
		Permission:       perm,
		BillingAddress:   a.BillingAddress.addressToPartner(&at),
		Active:           a.Status.IsActive(),
		LastStatusChange: &a.StatusDate,
	}, nil
}
