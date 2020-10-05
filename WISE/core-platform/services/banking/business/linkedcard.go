/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

// Package for all internal banking related items such as accounts and cards
package business

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"time"

	"github.com/wiseco/core-platform/services/banking"
	"github.com/wiseco/core-platform/shared"

	"github.com/wiseco/core-platform/services"
)

type CardNumber string

func (n *CardNumber) String() string {
	return string(*n)
}

// Marshal and transform fields as needed
func (n *CardNumber) MarshalJSON() ([]byte, error) {
	return json.Marshal(services.MaskLeft(n.String(), 4))
}

// Linked cards can be used for paying contractors, vendors, etc. via ACH push
// ACH pull requires full verification
type LinkedCard struct {
	// Bank account id
	Id string `json:"id" db:"id"`

	// Business Id
	BusinessID shared.BusinessID `json:"businessId" db:"business_id"`

	// Contact id
	ContactId *string `json:"contactId" db:"contact_id"`

	// Account Reference Id - This Id is from registered bank(bbva for example)
	RegisteredCardId string `json:"registeredCardId" db:"registered_card_id"`

	// Registered bank name - Ex: bbva
	RegisteredBankName string `json:"registeredBankName"  db:"registered_bank_name"`

	// Card number
	CardNumberMasked CardNumber `json:"cardNumberMasked" db:"card_number_masked"`

	// Card brand
	CardBrand string `json:"cardBrand" db:"card_brand"`

	// Card type
	CardType string `json:"cardType" db:"card_type"`

	// Card issuer
	CardIssuer string `json:"cardIssuer" db:"issuer_name"`

	// Fast funds enabled
	FastFundsEnabled bool `json:"fastFundsEnabled" db:"fast_funds_enabled"`

	// Card holder
	CardHolderName string `json:"cardHolderName" db:"card_holder_name"`

	// Account alias
	Alias *string `json:"alias" db:"alias"`

	// Account Usage - Ex: send_and_receive, send_only or receive_only
	Permission banking.LinkedAccountPermission `json:"permission" db:"account_permission"`

	// Billing address
	BillingAddress *services.Address `json:"billingAddress,omitempty" db:"billing_address"`

	// External accounts need to be verified
	Verified bool `json:"verified" db:"verified"`

	// Account usage type
	UsageType *UsageType `json:"usageType" db:"usage_type"`

	CardNumberHashed *string `json:"-" db:"card_number_hashed"`

	// Deactivated
	Deactivated *time.Time `json:"deactivated" db:"deactivated"`

	// Created timestamp
	Created time.Time `json:"created" db:"created"`

	// Modified timestamp
	Modified time.Time `json:"modified" db:"modified"`
}

type LinkedCardCreate struct {
	// User Id
	UserID shared.UserID `json:"userId" db:"user_id"`

	// Business Id
	BusinessID shared.BusinessID `json:"businessId" db:"business_id"`

	// Name on the account
	ContactId *string `json:"contactId" db:"contact_id"`

	// Name on the card
	CardHolderName string `json:"cardHolderName" db:"card_holder_name"`

	// Account alias
	Alias string `json:"alias" db:"alias"`

	// Card number
	CardNumber CardNumber `json:"cardNumber" db:"card_number_masked"`

	// Expiration date
	ExpirationDate shared.ExpDate `json:"expirationDate"`

	// CVC code
	CVVCode string `json:"cvvCode"`

	// Account Usage - Ex: send_and_receive, send_only or receive_only
	Permission banking.LinkedAccountPermission `json:"permission" db:"account_permission"`

	// Billing address
	BillingAddress *services.Address `json:"billingAddress,omitempty" db:"billing_address"`

	// Account usage type
	UsageType *UsageType `json:"usageType" db:"usage_type"`

	// Validate card only for debit pull
	ValidateCard bool 
}

type LinkedCardUpdate struct {
	// Bank account id
	ID string `json:"id" db:"id"`

	// Account usage type
	UsageType *UsageType `json:"usageType" db:"usage_type"`
}

func (c *LinkedCardCreate) HashLinkedCard() *string {
	h := sha256.New()
	card := c.CardNumber.String() + "-" + c.ExpirationDate.String() + "-" + c.CVVCode + "-" + string(c.BusinessID)
	h.Write([]byte(card))
	hash := base64.URLEncoding.EncodeToString(h.Sum(nil))
	return &hash
}
