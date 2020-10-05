package main

import (
	"time"
)

const (
	crossAccountIPAddress string = "IP addresses used on another account"
	invalidCountry        string = "Login attempt from outside United States"
	highRiskLevel         string = "Login tagged as high risk"
)

type flagged struct {
	IsFlagged bool
	Reasons   []string
}

type authEvent struct {
	EventType    string
	CreationDate time.Time
	RiskLevel    string
	RiskDecision string
	IPAddress    string
	City         string
	Country      string
	DeviceName   string
	Flagged      flagged
}

type cognitoUser struct {
	Phone      string
	Username   string
	AuthEvents []authEvent
}
