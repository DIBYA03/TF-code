/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

// Package for all internal banking related items such as accounts and cards
package banking

const (
	// Checking or deposit account
	BankAccountTypeChecking = "checking"

	// Savings or interest bearing account
	BankAccountTypeSavings = "savings"
)

const (
	// Bank account status is active
	BankAccountStatusActive = "active"

	// Bank account status is blocked by bank
	BankAccountStatusBlocked = "blocked"

	// Bank account status is locked by the customer
	BankAccountStatusLocked = "locked"

	// Bank account status is pending close
	BankAccountStatusClosePending = "closePending"

	// Bank account status is closed
	BankAccountStatusClosed = "closed"

	// Bank account status is dormant
	BankAccountStatusDormant = "dormant"

	// Bank account status is abandoned
	BankAccountStatusAbandoned = "abandoned"

	// Bank account status is reverted to the state
	BankAccountStatusEscheated = "escheated"

	// Bank account status is closed and ovedrawn
	BankAccountStatusChargeOff = "chargeOff"
)

const (
	// Allows view only access
	BankAccountAccessView = "view"

	// Allows viewing and exporting reports
	BankAccountAccessAccounting = "accounting"

	// Deposit access
	BankAccountAccessDeposit = "deposit"

	// Withdrawal access
	BankAccountAccessWithdraw = "withdraw"

	// Deposit and withdrawal access
	BankAccountAccessFull = "full"

	// Allows user to update or close account
	BankAccountAccessAdmin = "admin"
)

const (
	// Access role is account owner
	BankAccountAccessRoleOwner = "owner"

	// Access role is spouse
	BankAccountAccessRoleSpouse = "spouse"

	// Access role is minor
	BankAccountAccessRoleMinor = "minor"

	// Access role is attorney
	BankAccountAccessRoleAttorney = "attorney"

	// Access role is beneficiary
	BankAccountAccessRoleBeneficiary = "beneficiary"

	// Access role is custodian
	BankAccountAccessRoleCustodian = "custodian"

	// Access role is conservator
	BankAccountAccessRoleConservator = "conservator"

	// Access role is officer
	BankAccountAccessRoleOfficer = "officer"

	// Access role is
	BankAccountAccessRoleEmployee = "employee"

	// Access role is other
	BankAccountAccessRoleOther = "other"
)

type ConsumerBankAccountAccess struct {
	// Account access id
	Id string `json:"id"`

	// Consumer bank account id
	BankAccountId string `json:"bankAccountId"`

	// Related user id
	UserId string `json:"userId"`

	// Access type e.g. deposit, withdrawal, or admin
	AccessType string `json:"accessType"`

	// Access role e.g. owner, spouse, or employee
	AccessRole string `json:"accessRole"`
}

type BusinessBankAccountAccess struct {
	// Account access id
	Id string `json:"id"`

	// Business id
	BusinessId string `json:"businessId"`

	// Consumer bank account id
	BankAccountId string `json:"bankAccountId"`

	// Related user id
	UserId string `json:"userId"`

	// Access type e.g. deposit, withdrawal, or admin
	AccessType string `json:"accessType"`

	// Access role e.g. owner, spouse, or employee
	AccessRole string `json:"accessRole"`
}
