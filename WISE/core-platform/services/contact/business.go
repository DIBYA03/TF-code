/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

// Package for user contacts
package contact

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/wiseco/core-platform/services"
	"github.com/wiseco/core-platform/shared"
)

const (
	ContactCategoryContractor = ContactCategory("contractor")
	ContactCategoryVendor     = ContactCategory("vendor")
	ContactCategoryShopify    = ContactCategory("shopify") // Contact added from shopify order
)

type ContactCategory string

const (
	ContactTypePerson   = ContactType("person")
	ContactTypeBusiness = ContactType("business")
)

type AccountNumber string

func (n *AccountNumber) String() string {
	return string(*n)
}

// Marshal and transform fields as needed
func (n *AccountNumber) MarshalJSON() ([]byte, error) {
	return json.Marshal(services.MaskLeft(n.String(), 4))
}

type ContactType string

type Contact struct {
	// Contact id
	ID string `json:"id" db:"id"`

	// User id
	UserID shared.UserID `json:"userId" db:"user_id"`

	// Business id
	BusinessID shared.BusinessID `json:"businessId" db:"business_id"`

	// Contact type
	Category *ContactCategory `json:"category" db:"contact_category"`

	// Contact sub type
	Type ContactType `json:"type" db:"contact_type"`

	// Engagement
	Engagement *string `json:"engagement" db:"engagement"`

	// Job title
	JobTitle *string `json:"jobTitle" db:"job_title"`

	// Business name - used if sub type is business
	BusinessName *string `json:"businessName" db:"business_name"`

	// First name
	FirstName *string `json:"firstName" db:"first_name"`

	// Last name
	LastName *string `json:"lastName" db:"last_name"`

	// Phone number
	PhoneNumber string `json:"phoneNumber" db:"phone_number"`

	// Email
	Email string `json:"email" db:"email"`

	// Email ID
	EmailID shared.EmailID `json:"emailID" db:"email_id"`

	// Mailing address
	MailingAddress *services.Address `json:"mailingAddress" db:"mailing_address"`

	// Deactivated
	Deactivated *time.Time `json:"deactivated" db:"deactivated"`

	// Created timestamp
	Created time.Time `json:"created" db:"created"`

	// Modified timestamp
	Modified time.Time `json:"modified" db:"modified"`
}

func (c *Contact) Name() string {
	var name string
	switch c.Type {
	case ContactTypePerson:
		name = shared.StringValue(c.FirstName) + " " + shared.StringValue(c.LastName)
	case ContactTypeBusiness:
		name = shared.StringValue(c.BusinessName)
	}

	return strings.TrimSpace(name)
}

type ContactCreate struct {
	// User id
	UserID shared.UserID `json:"userId" db:"user_id"`

	// Business id
	BusinessID shared.BusinessID `json:"businessId" db:"business_id"`

	// Contact type
	Category *ContactCategory `json:"category" db:"contact_category"`

	// Contact sub type
	Type ContactType `json:"type" db:"contact_type"`

	// Engagement
	Engagement *string `json:"engagement" db:"engagement"`

	// Job title
	JobTitle *string `json:"jobTitle" db:"job_title"`

	// Business name - used if sub type is business
	BusinessName *string `json:"businessName" db:"business_name"`

	// First name
	FirstName *string `json:"firstName" db:"first_name"`

	// Last name
	LastName *string `json:"lastName" db:"last_name"`

	// Phone number
	PhoneNumber string `json:"phoneNumber" db:"phone_number"`

	// Email
	Email string `json:"email" db:"email"`

	// Email ID
	EmailID shared.EmailID `json:"emailID" db:"email_id"`
}

type ContactUpdate struct {
	// Contact id
	ID string `json:"id" db:"id"`

	// User id
	UserID shared.UserID `json:"userId" db:"user_id"`

	// Business id
	BusinessID shared.BusinessID `json:"businessId" db:"business_id"`

	// Contact type
	Category *ContactCategory `json:"category" db:"contact_category"`

	// Contact sub type
	Type *ContactType `json:"type" db:"contact_type"`

	// Engagement
	Engagement *string `json:"engagement" db:"engagement"`

	// Job title
	JobTitle *string `json:"jobTitle" db:"job_title"`

	// Business name - used if sub type is business
	BusinessName string `json:"businessName" db:"business_name"`

	// First name
	FirstName *string `json:"firstName" db:"first_name"`

	// Last name
	LastName *string `json:"lastName" db:"last_name"`

	// Phone number
	PhoneNumber *string `json:"phoneNumber" db:"phone_number"`

	// Email
	Email *string `json:"email" db:"email"`

	// Email ID
	EmailID shared.EmailID `json:"emailID" db:"email_id"`
}
