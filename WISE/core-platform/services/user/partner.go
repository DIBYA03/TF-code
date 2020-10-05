/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

// Package for all user related services
package user

import (
	"time"

	"github.com/wiseco/core-platform/shared"
)

// Partner entity
type Partner struct {

	// Partner id (uuid)
	ID shared.PartnerID `json:"id" db:"id"`

	// Partner name
	Name string `json:"name" db:"channel_name"`

	// Partner code
	Code string `json:"code" db:"code"`

	// License count given to partner
	GrantedLicenseCount int `json:"grantedLicenseCount" db:"granted_license_count"`

	// Deactivated
	Deactivated *time.Time `json:"deactivated" db:"deactivated"`

	// Created time
	Created time.Time `json:"created" db:"created"`

	// Last modified time
	Modified time.Time `json:"modified" db:"modified"`
}

// Partner verification
type PartnerVerification struct {

	// Partner code
	Code string `json:"code" db:"code"`

	// User ID
	UserID shared.UserID
}
