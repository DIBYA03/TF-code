package bank

import "time"

type LinkedCardPermission string

func (c LinkedCardPermission) String() string {
	return string(c)
}

const (
	LinkedCardPermissionSendOnly       = LinkedCardPermission("sendOnly")
	LinkedCardPermissionRecieveOnly    = LinkedCardPermission("receiveOnly")
	LinkedCardPermissionSendAndRecieve = LinkedCardPermission("sendAndReceive")
)

type LinkedCardRequest struct {
	AccountHolder  string               `json:"accountHolder"`
	Alias          string               `json:"alias"`
	CardNumber     string               `json:"cardNumber"`
	Expiration     time.Time            `json:"expiration"`
	CVC            string               `json:"cvc"`
	Currency       Currency             `json:"currency"`
	Permission     LinkedCardPermission `json:"permission"`
	BillingAddress AddressRequest       `json:"billingAddress"`
}

type LinkedCardType string

const (
	LinkedCardTypeUnknown      = LinkedCardType("unknown")
	LinkedCardTypeDebit        = LinkedCardType("debit")
	LinkedCardTypePrepaidDebit = LinkedCardType("prepaidDebit")
	LinkedCardTypeCredit       = LinkedCardType("credit")
	LinkedCardTypeCharge       = LinkedCardType("charge")
)

type LinkedCardBrand string

const (
	LinkedCardBrandUnknown    = LinkedCardBrand("unknown")
	LinkedCardBrandVisa       = LinkedCardBrand("visa")
	LinkedCardBrandMastercard = LinkedCardBrand("mc")
	LinkedCardBrandDiscover   = LinkedCardBrand("discover")
	LinkedCardBrandAmex       = LinkedCardBrand("amex")
)

type LinkedCardBankID string

type LinkedCardResponse struct {
	CardID           LinkedCardBankID     `json:"cardId"`
	CardType         LinkedCardType       `json:"cardType"`
	BrandType        LinkedCardBrand      `json:"brandType"`
	IssuingBank      string               `json:"issuingBank"`
	AccountHolder    string               `json:"accountHolder"`
	Alias            string               `json:"alias"`
	CardNumber       string               `json:"cardNumber"`
	Expiration       string               `json:"expiration"`
	Currency         Currency             `json:"currency"`
	FastFundsEnabled bool                 `json:"fastFundsEnabled"`
	Permission       LinkedCardPermission `json:"permission"`
	BillingAddress   AddressResponse      `json:"billingAddress"`
	Active           bool                 `json:"active"`
	LastStatusChange *time.Time           `json:"lastStatusChange"`
}
