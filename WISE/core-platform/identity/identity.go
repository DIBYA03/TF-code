/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/
package identity

import (
	"time"

	"github.com/wiseco/core-platform/shared"
)

type ProviderID string
type ProviderName string
type ProviderSource string

const ProviderNameCognito = ProviderName("cognito")

// Identity object
type Identity struct {

	// Identity server id id (uuid)
	ID shared.IdentityID `json:"id" db:"id"`

	// Identity provider id
	ProviderID ProviderID `json:"providerId" db:"provider_id"`

	// Auth Type e.g. cognito, etc
	ProviderName ProviderName `json:"providerName" db:"provider_name"`

	// Identity Source e.g server or id
	ProviderSource ProviderSource `json:"providerSource" db:"provider_source"`

	// Phone
	Phone string `json:"phone" db:"phone"`

	// Deactivated
	Deactivated *time.Time `json:"deactivated" db:"deactivated"`

	// Created time
	Created time.Time `json:"created" db:"created"`

	// Last modified time
	Modified time.Time `json:"modified" db:"modified"`
}

// IdentityCreate is used to create an user's identity
type IdentityCreate struct {
	// Identity provider id
	ProviderID ProviderID `json:"providerId" db:"provider_id"`

	// Auth Type e.g. cognito, etc
	ProviderName ProviderName `json:"providerName" db:"provider_name"`

	// Identity Source e.g server or id
	ProviderSource ProviderSource `json:"providerSource" db:"provider_source"`

	// Phone
	Phone string `json:"phone" db:"phone"`
}

// IdentityUpdate is used to update an user's identity
type IdentityUpdate struct {
	// Identity server id id (uuid)
	ID shared.IdentityID `json:"id" db:"id"`

	// Phone
	Phone string `json:"phone" db:"phone"`
}
