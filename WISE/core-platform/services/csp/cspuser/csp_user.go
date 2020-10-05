/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

// Package cspuser for all csp user related services
package cspuser

import (
	"time"

	"github.com/wiseco/go-lib/id"
)

// CSPUser is the object defining an csp user
type CSPUser struct {

	// User id (uuid)
	ID id.CspAgentID `json:"id" db:"id"`

	// Cognito id (uuid)
	CognitoID string `json:"cognito_id" db:"cognito_id"`

	// FirstName is split off of `name` from cognito
	FirstName string `json:"firstName" db:"first_name"`

	// MiddleName is split off of `name` from cognito
	MiddleName string `json:"middleName" db:"middle_name"`

	// LastName is split off of `name` from cognito
	LastName string `json:"lastName" db:"last_name"`

	// Email
	Email *string `json:"email" db:"email"`

	// Email verified
	EmailVerified bool `json:"emailVerified" db:"email_verified"`

	// Phone
	Phone string `json:"phone" db:"phone"`

	// Phone verified
	PhoneVerified bool `json:"phoneVerified" db:"phone_verified"`

	// Picutre from Google of csp user
	Picture string `json:"picture" db:"picture"`

	// Created time
	Created time.Time `json:"created" db:"created"`

	// Last modified time
	Modified time.Time `json:"modified" db:"modified"`

	Active bool `json:"active" db:"active"`
}
