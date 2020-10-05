package bank

import "time"

type LinkedAccountPermission string

func (a LinkedAccountPermission) String() string {
	return string(a)
}

const (
	LinkedAccountPermissionSendOnly       = LinkedAccountPermission("sendOnly")
	LinkedAccountPermissionRecieveOnly    = LinkedAccountPermission("receiveOnly")
	LinkedAccountPermissionSendAndRecieve = LinkedAccountPermission("sendAndReceive")
)

type LinkedBankAccountRequest struct {
	AccountNumber     string                  `json:"accountNumber"`
	RoutingNumber     string                  `json:"routingNumber"`
	WireRouting       string                  `json:"wireRouting"`
	AccountType       AccountType             `json:"accountType"`
	AccountHolderName string                  `json:"nameOnAccount"`
	Currency          Currency                `json:"currency"`
	Alias             string                  `json:"nickname"`
	Permission        LinkedAccountPermission `json:"permission"`
}

type LinkedAccountBankID string

type LinkedBankAccountResponse struct {
	AccountID         LinkedAccountBankID     `json:"bankAccountId"`
	AccountNumber     string                  `json:"accountNumber"`
	RoutingNumber     string                  `json:"routingNumber"`
	WireRouting       *string                 `json:"wireRouting"`
	AccountBankName   string                  `json:"accountBankName"`
	AccountType       AccountType             `json:"accountType"`
	AccountHolderName string                  `json:"nameOnAccount"`
	Currency          Currency                `json:"currency"`
	Alias             string                  `json:"alias"`
	Permission        LinkedAccountPermission `json:"permission"`
	Status            AccountStatus           `json:"status"`
	StatusDate        time.Time               `json:"statusDate"`
}
