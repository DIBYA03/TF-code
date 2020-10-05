/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

// Package for all user related services
package user

import (
	"time"

	"github.com/wiseco/core-platform/shared"
)

type DeviceType string

const (
	// Android device
	DeviceTypeAndroid = DeviceType("android")

	// iOS device
	DeviceTypeIOS = DeviceType("ios")

	// Web browser
	DeviceTypeWeb = DeviceType("web")
)

const (
	// FCM token
	TokenTypeFCM = "fcm"

	// APNS token
	TokenTypeAPNS = "apns"

	// Web push token
	TokenTypeWeb = "web"
)

type DeviceKey string

// Push registration object for iOS, Android, and Web
type PushRegistration struct {
	// Push registration id
	ID shared.UserDeviceID `json:"id" db:"id"`

	UserID shared.UserID `json:"userId" db:"user_id"`

	// Device token type e.g. ios or android
	DeviceType DeviceType `json:"deviceType" db:"device_type"`

	//DeviceKey
	DeviceKey DeviceKey `json:"deviceKey" db:"device_key"`

	//user language
	Language string `json:"language" db:"language"`

	// Token type e.g. apns or fcm
	TokenType string `json:"tokenType" db:"token_type"`

	// Device token
	Token string `json:"token" db:"token"`

	// Created time
	Created time.Time `json:"created" db:"created"`

	// Last modified time
	Modified time.Time `json:"modified" db:"modified"`
}

type PushRegistrationCreate struct {
	// User id
	UserID shared.UserID `json:"userId" db:"user_id"`

	// Device token type e.g. ios or android
	DeviceType DeviceType `json:"deviceType" db:"device_type"`

	// Token type e.g. apns or fcm
	TokenType string `json:"tokenType" db:"token_type"`

	//user language
	Language string `json:"language" db:"language"`

	// DeviceKey of local device (if any)
	DeviceKey DeviceKey `json:"deviceKey" db:"device_key"`

	// Device token
	Token string `json:"token" db:"token"`
}

type DeviceLogout struct {
	// DeviceKey of local device (if any)
	DeviceKey DeviceKey `json:"deviceKey"`
}
