package bbva

import (
	"errors"
	"time"

	"github.com/wiseco/core-platform/partner/bank"
)

type RegisteredAccountUsage string

func (u RegisteredAccountUsage) String() string {
	return string(u)
}

const (
	RegisteredAccountUsageSendOnly       = RegisteredAccountUsage("send_only")
	RegisteredAccountUsageRecieveOnly    = RegisteredAccountUsage("receive_only")
	RegisteredAccountUsageSendAndRecieve = RegisteredAccountUsage("send_and_receive")
)

var partnerLinkedAccountPermissionFrom = map[bank.LinkedAccountPermission]RegisteredAccountUsage{
	bank.LinkedAccountPermissionSendOnly:       RegisteredAccountUsageSendOnly,
	bank.LinkedAccountPermissionRecieveOnly:    RegisteredAccountUsageRecieveOnly,
	bank.LinkedAccountPermissionSendAndRecieve: RegisteredAccountUsageSendAndRecieve,
}

var partnerLinkedAccountPermissionTo = map[RegisteredAccountUsage]bank.LinkedAccountPermission{
	RegisteredAccountUsageSendOnly:       bank.LinkedAccountPermissionSendOnly,
	RegisteredAccountUsageRecieveOnly:    bank.LinkedAccountPermissionRecieveOnly,
	RegisteredAccountUsageSendAndRecieve: bank.LinkedAccountPermissionSendAndRecieve,
}

type RegisteredAccountType string

const (
	RegisteredAccountTypeConsumerChecking = RegisteredAccountType("consumer_checking")
	RegisteredAccountTypeConsumerSavings  = RegisteredAccountType("consumer_savings")
	RegisteredAccountTypeBusinessChecking = RegisteredAccountType("business_checking")
	RegisteredAccountTypeBusinessSavings  = RegisteredAccountType("business_savings")
)

var partnerLinkedAccountTypeFromConsumer = map[bank.AccountType]RegisteredAccountType{
	bank.AccountTypeChecking: RegisteredAccountTypeConsumerChecking,
	bank.AccountTypeSavings:  RegisteredAccountTypeConsumerSavings,
}

var partnerLinkedAccountTypeFromBusiness = map[bank.AccountType]RegisteredAccountType{
	bank.AccountTypeChecking: RegisteredAccountTypeBusinessChecking,
	bank.AccountTypeSavings:  RegisteredAccountTypeBusinessSavings,
}

var partnerLinkedAccountTypeToConsumer = map[RegisteredAccountType]bank.AccountType{
	RegisteredAccountTypeConsumerChecking: bank.AccountTypeChecking,
	RegisteredAccountTypeConsumerSavings:  bank.AccountTypeSavings,
}

var partnerLinkedAccountTypeToBusiness = map[RegisteredAccountType]bank.AccountType{
	RegisteredAccountTypeBusinessChecking: bank.AccountTypeChecking,
	RegisteredAccountTypeBusinessSavings:  bank.AccountTypeSavings,
}

type RegisterBankAccountRequest struct {
	AccountNumber string                 `json:"account_number"`
	RoutingNumber string                 `json:"routing_number"`
	AccountType   RegisteredAccountType  `json:"account_type"`
	NameOnAccount string                 `json:"name_on_account"`
	Currency      Currency               `json:"currency"`
	Nickname      string                 `json:"nickname"`
	Usage         RegisteredAccountUsage `json:"usage"`
}

type RegisterBankAccountResponse struct {
	AccountReferenceID string `json:"account_reference_id"`
	AccountLast4       string `json:"account_last4"`
}

type GetRegisterBankAccountsResponse struct {
	Accounts []GetRegisterBankAccountResponse `json:"registered_accounts"`
}

type GetRegisterBankAccountResponse struct {
	AccountReferenceID string                 `json:"account_reference_id"`
	AccountNumber      string                 `json:"account_number"`
	RoutingNumber      string                 `json:"routing_number"`
	BankName           string                 `json:"bank_name"`
	AccountType        RegisteredAccountType  `json:"account_type"`
	NameOnAccount      string                 `json:"name_on_account"`
	Currency           Currency               `json:"currency"`
	Nickname           string                 `json:"nickname"`
	Usage              RegisteredAccountUsage `json:"usage"`
	Status             AccountStatus          `json:"status"`
	StatusDate         time.Time              `json:"status_date"`
}

func (a *GetRegisterBankAccountResponse) partnerLinkedBankAccountResponseTo() (*bank.LinkedBankAccountResponse, error) {
	accountType, ok := partnerLinkedAccountTypeToConsumer[a.AccountType]
	if !ok {
		accountType, ok = partnerLinkedAccountTypeToBusiness[a.AccountType]
		if !ok {
			return nil, errors.New("invalid account type")
		}
	}

	currency, ok := partnerCurrencyTo[a.Currency]
	if !ok {
		return nil, errors.New("invalid currency")
	}

	perm, ok := partnerLinkedAccountPermissionTo[a.Usage]
	if !ok {
		return nil, errors.New("invalid permission")
	}

	status, ok := partnerAccountStatusTo[a.Status]
	if !ok {
		return nil, errors.New("invalid account status")
	}

	return &bank.LinkedBankAccountResponse{
		AccountID:         bank.LinkedAccountBankID(a.AccountReferenceID),
		AccountNumber:     a.AccountNumber,
		RoutingNumber:     a.RoutingNumber,
		AccountBankName:   a.BankName,
		AccountType:       accountType,
		AccountHolderName: a.NameOnAccount,
		Currency:          currency,
		Alias:             a.Nickname,
		Permission:        perm,
		Status:            status,
		StatusDate:        a.StatusDate,
	}, nil
}
