/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

package payment

import (
	"time"

	"github.com/wiseco/core-platform/shared"
)

type CardReaderCreate struct {
	BusinessID shared.BusinessID `json:"businessId" db:"business_id"`

	Alias *string `json:"alias" db:"alias"`

	DeviceType *string `json:"deviceType" db:"device_type"`

	SerialNumber string `json:"serialNumber" db:"serial_number"`
}

type CardReaderUpdate struct {
	ID shared.CardReaderID `json:"id" db:"id"`

	Alias *string `json:"alias" db:"alias"`

	DeviceType *string `json:"deviceType" db:"device_type"`

	CreatedUserID shared.UserID `json:"createdUserId"`

	BusinessID shared.BusinessID `json:"businessId" db:"business_id"`

	LastConnected *time.Time `json:"lastConnected" db:"last_connected"`
}

type CardReader struct {
	ID shared.CardReaderID `json:"id" db:"id"`

	CreatedUserID shared.UserID `json:"createdUserId"`

	BusinessID shared.BusinessID `json:"businessId" db:"business_id"`

	Alias *string `json:"alias" db:"alias"`

	DeviceType *string `json:"deviceType" db:"device_type"`

	SerialNumber string `json:"serialNumber" db:"serial_number"`

	LastConnected *time.Time `json:"lastConnected" db:"last_connected"`

	Deactivated *time.Time `json:"deactivated" db:"deactivated"`

	Created time.Time `json:"created" db:"created"`

	Modified time.Time `json:"modified" db:"modified"`
}
