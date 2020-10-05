/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

package email

import (
	"time"

	"github.com/wiseco/core-platform/shared"
)

type Status string
type Type string

type EmailAddress string

const (
	StatusActive   = Status("active")
	StatusInactive = Status("inactive")

	TypeContact        = Type("contact")
	TypeConsumer       = Type("consumer")
	TypeBusiness       = Type("business")
	TypeBusinessMember = Type("business_member")
)

// Address refers to a physical or mailing address
type Email struct {
	// Email id
	ID shared.EmailID `json:"id" db:"id"`

	// Email address
	EmailAddress EmailAddress `json:"emailAddress" db:"email_address"`

	// Email status
	Status Status `json:"emailStatus" db:"email_status"`

	// Email type
	Type Type `json:"emailType" db:"email_type"`

	// Created timestamp
	Created time.Time `json:"created" db:"created"`

	// Modified timestamp
	Modified time.Time `json:"modified" db:"modified"`
}

type EmailCreate struct {
	// Email id
	ID shared.EmailID `json:"id" db:"id"`

	// Email address
	EmailAddress EmailAddress `json:"emailAddress" db:"email_address"`

	// Email status
	Status Status `json:"emailStatus" db:"email_status"`

	// Email type
	Type Type `json:"emailType" db:"email_type"`
}
