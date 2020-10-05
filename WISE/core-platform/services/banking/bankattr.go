/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

// Package for all internal banking related items such as accounts and cards
package banking

import (
	partnerbank "github.com/wiseco/core-platform/partner/bank"
	_ "github.com/wiseco/core-platform/partner/bank/bbva"
)

type BankName string

const (
	BankNameBBVA = BankName(partnerbank.ProviderNameBBVA)
)

type PartnerName string

const (
	PartnerNameBBVA   = PartnerName(partnerbank.ProviderNameBBVA)
	PartnerNamePlaid  = PartnerName(partnerbank.ProviderNamePlaid)
	PartnerNameStripe = PartnerName(partnerbank.ProviderNameStripe)
)

var ToPartnerBankName = map[BankName]partnerbank.ProviderName{
	BankNameBBVA: partnerbank.ProviderNameBBVA,
}

type Currency string

const (
	// ISO currency code
	CurrencyUSD = Currency("usd")
)

type LinkedAccountPermission string

func (u LinkedAccountPermission) String() string {
	return string(u)
}

const (
	LinkedAccountPermissionSendOnly       = LinkedAccountPermission("sendOnly")
	LinkedAccountPermissionRecieveOnly    = LinkedAccountPermission("receiveOnly")
	LinkedAccountPermissionSendAndRecieve = LinkedAccountPermission("sendAndReceive")
)

type AccountType string

func (t AccountType) String() string {
	return string(t)
}

const (
	AccountTypeChecking = AccountType("checking")
	AccountTypeSavings  = AccountType("savings")
)
