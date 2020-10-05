/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

package transaction

import (
	"github.com/wiseco/core-platform/services"
	"github.com/wiseco/core-platform/services/banking/business"
	"github.com/wiseco/core-platform/services/contact"
)

// Contact mini object for notifications
type Contact struct {
	// Contact id
	Id *string `json:"id,omitempty" db:"id"`

	// Contact type
	Category *contact.ContactCategory `json:"category,omitempty" db:"contact_category"`

	// Contact sub type
	Type *contact.ContactType `json:"type,omitempty" db:"contact_type"`

	// Engagement
	Engagement *string `json:"engagement,omitempty" db:"engagement"`

	// Job title
	JobTitle *string `json:"jobTitle,omitempty" db:"job_title"`

	// Business name - used if sub type is business
	BusinessName *string `json:"businessName,omitempty" db:"business_name"`

	// First name
	FirstName *string `json:"firstName,omitempty" db:"first_name"`

	// Last name
	LastName *string `json:"lastName,omitempty" db:"last_name"`

	// Phone number
	PhoneNumber *string `json:"phoneNumber,omitempty" db:"phone_number"`

	// Email
	Email *string `json:"email,omitempty" db:"email"`

	// Mailing address
	MailingAddress *services.Address `json:"mailingAddress,omitempty" db:"mailing_address"`

	// Masked account number
	AccountNumber *contact.AccountNumber `json:"accountNumberMasked,omitempty" db:"account_number"`

	// Routing number
	RoutingNumber *string `json:"routingNumber,omitempty" db:"routing_number"`

	// Bank name
	BankName *string `json:"bankName,omitempty" db:"bank_name"`

	// Masked card number
	CardNumber *business.CardNumber `json:"cardNumberMasked,omitempty" db:"card_number_masked"`

	// card brand
	CardBrand *string `json:"cardBrand,omitempty" db:"card_brand"`
}
