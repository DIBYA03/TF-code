/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

package bank

import (
	"time"

	"github.com/jmoiron/sqlx/types"
)

const (
	NotificationTypeConsumer      = NotificationType("consumer")
	NotificationTypeBusiness      = NotificationType("business")
	NotificationTypeAccount       = NotificationType("account")
	NotificationTypeCard          = NotificationType("card")
	NotificationTypeMoneyTransfer = NotificationType("moneyTransfer")
	NotificationTypeTransaction   = NotificationType("transaction")
)

// Consumer and Business
const (
	NotificationAttributeEmail       = NotificationAttribute("email")
	NotificationAttributePhone       = NotificationAttribute("phone")
	NotificationAttributeAddress     = NotificationAttribute("address")
	NotificationAttributeOccupation  = NotificationAttribute("occupation")
	NotificationAttributeActivity    = NotificationAttribute("activity")
	NotificationAttributeParticipant = NotificationAttribute("participant")
	NotificationAttributeStatus      = NotificationAttribute("status")
	NotificationAttributeBlock       = NotificationAttribute("block")
	NotificationAttributeChargeOff   = NotificationAttribute("chargeOff")
	NotificationAttributeKYC         = NotificationAttribute("kyc")
)

// Business
const (
	NotificationAttributeOwner = NotificationAttribute("owner")
)

const (
	NotificationActionCreate    = NotificationAction("create")
	NotificationActionOpen      = NotificationAction("open")
	NotificationActionActivate  = NotificationAction("activate")
	NotificationActionReissue   = NotificationAction("reissue")
	NotificationActionUpdate    = NotificationAction("update")
	NotificationActionAdd       = NotificationAction("add")
	NotificationActionRemove    = NotificationAction("remove")
	NotificationActionDelete    = NotificationAction("delete")
	NotificationActionCancel    = NotificationAction("cancel")
	NotificationActionCorrected = NotificationAction("corrected")
	NotificationActionAuthorize = NotificationAction("authorize")
	NotificationActionHold      = NotificationAction("hold")
	NotificationActionPosted    = NotificationAction("posted")
)

// Consumer specific notification
type ConsumerNotification struct {
	BankID             ConsumerBankID     `json:"bankId"`
	TaxID              string             `json:"taxId"`
	TaxIDType          ConsumerTaxIDType  `json:"taxIdType"`
	FirstName          string             `json:"firstName"`
	MiddleName         *string            `json:"middleName"`
	LastName           string             `json:"lastName"`
	DateOfBirth        *time.Time         `json:"dateOfBirth"`
	Residency          *ConsumerResidency `json:"residency"`
	CitizenshipCountry *Country           `json:"citizenshipCountry"`
	Documents          []ConsumerDocument `json:"document"`
}

type ConsumerAttributeNotification struct {
	BankID      ConsumerBankID      `json:"bankId"`
	AttributeID string              `json:"attributeId"`
	Phone       *string             `json:"phone,omitempty"`
	Email       *string             `json:"email,omitempty"`
	Address     *AddressResponse    `json:"address,omitempty"`
	Occupation  *ConsumerOccupation `json:"occupation,omitempty"`
}

type ConsumerKYCNotesNotification struct {
	Source string `json:"source"`
	Desc   string `json:"desc"`
}

type ConsumerKYCNotification struct {
	BankID    ConsumerBankID                 `json:"bankId"`
	Risk      KYCRisk                        `json:"risk"`
	Result    KYCResult                      `json:"result"`
	KYCStatus KYCStatus                      `json:"kycStatus"`
	KYCNotes  []ConsumerKYCNotesNotification `json:"kycNotes"`
}

type BusinessFormationNotification struct {
	OriginCountry Country   `json:"originCountry"`
	OriginState   string    `json:"originState"`
	OriginDate    time.Time `json:"originDate"`
}

type BusinessDocumentNotification struct {
	DocumentType   BusinessIdentityDocument `json:"docType"`
	Number         string                   `json:"number"`
	Issuer         *string                  `json:"issuer"`
	IssueDate      *time.Time               `json:"issueDate"`
	IssueState     *string                  `json:"state"`
	IssueCountry   *Country                 `json:"country"`
	ExpirationDate *time.Time               `json:"expirationDate"`
}

// Business specific notification
type BusinessNotification struct {
	BankID        BusinessBankID                 `json:"bankId"`
	LegalName     string                         `json:"legalName"`
	TaxID         *string                        `json:"taxId"`
	TaxIDType     *BusinessTaxIDType             `json:"taxIdType"`
	EntityType    *BusinessEntity                `json:"entityType"`
	IndustryType  *BusinessIndustry              `json:"industryType"`
	OperationType *BusinessOperationType         `json:"operationType"`
	OriginCountry *Country                       `json:"originCountry"`
	OriginState   *string                        `json:"originState"`
	OriginDate    *time.Time                     `json:"originDate"`
	Documents     []BusinessDocumentNotification `json:"documents"`
}

type BusinessMemberNotification struct {
	MemberBankID         MemberBankID        `json:"memberBankId"`
	ConsumerBankID       ConsumerBankID      `json:"bankEntityId"`
	ConsumerID           ConsumerID          `json:"userId"`
	IsControllingManager bool                `json:"isControllingManager"`
	Ownership            int                 `json:"ownership"`
	Title                BusinessMemberTitle `json:"title"`
	TitleDesc            *string             `json:"titleDesc"`
}

type BusinessMembersNotification struct {
	BankID  BusinessBankID               `json:"bankId"`
	Members []BusinessMemberNotification `json:"members"`
}

type BusinessAttributeNotification struct {
	BankID      BusinessBankID     `json:"bankId"`
	AttributeID string             `json:"attributeId"`
	Phone       *string            `json:"phone,omitempty"`
	Email       *string            `json:"email,omitempty"`
	Address     *AddressResponse   `json:"address,omitempty"`
	Activities  []ExpectedActivity `json:"activities"`
}

type BusinessKYCNotification struct {
	BankID    BusinessBankID `json:"bankId"`
	Risk      KYCRisk        `json:"risk"`
	Result    KYCResult      `json:"result"`
	KYCStatus KYCStatus      `json:"kycStatus"`
}

// Account specific notification
type AccountOpenedNotification struct {
	AccountID AccountBankID `json:"bankAccountId"`
	Opened    *time.Time    `json:"accountType"`
}

type NotificationAccountStatus string

const (
	NotificationAccountStatusActive             = NotificationAccountStatus("active")
	NotificationAccountStatusInactive           = NotificationAccountStatus("inactive")
	NotificationAccountStatusBlocked            = NotificationAccountStatus("blocked")
	NotificationAccountStatusLocked             = NotificationAccountStatus("locked")
	NotificationAccountStatusClosed             = NotificationAccountStatus("closed")
	NotificationAccountStatusClosedBank         = NotificationAccountStatus("closedBank")
	NotificationAccountStatusClosedFraud        = NotificationAccountStatus("closedFraud")
	NotificationAccountStatusDormant            = NotificationAccountStatus("dormant")
	NotificationAccountStatusAbandoned          = NotificationAccountStatus("abandoned")
	NotificationAccountStatusPendingEscheatment = NotificationAccountStatus("pendingEscheatment")
	NotificationAccountStatusEscheated          = NotificationAccountStatus("escheated")
	NotificationAccountStatusChargeOff          = NotificationAccountStatus("chargeOff")
)

type AccountDataChangeNotification struct {
	AccountID AccountBankID `json:"bankAccountId"`
}

type AccountParticipantNotification struct {
	AccountID  AccountBankID          `json:"bankAccountId"`
	EntityID   NotificationEntityID   `json:"entityId"`
	EntityType NotificationEntityType `json:"entityType"`
}

type AccountStatusNotification struct {
	AccountID  AccountBankID             `json:"bankAccountId"`
	Status     NotificationAccountStatus `json:"status"`
	Reason     string                    `json:"reason"`
	StatusDate *time.Time                `json:"statusDate"`
}

type NotificationAccountBlock string

const (
	NotificationAccountBlockCreditAdded   = NotificationAccountBlock("creditBlockAdded")
	NotificationAccountBlockCreditRemoved = NotificationAccountBlock("creditBlockAddRemoved")
	NotificationAccountBlockDebitAdded    = NotificationAccountBlock("debitBlockAdded")
	NotificationAccountBlockDebitRemoved  = NotificationAccountBlock("debitBlockRemoved")
	NotificationAccountBlockCheckAdded    = NotificationAccountBlock("checkBlockAdded")
	NotificationAccountBlockCheckRemoved  = NotificationAccountBlock("checkBlockRemvoed")
)

type AccountBlockNotification struct {
	AccountID AccountBankID            `json:"bankAccountId"`
	Block     NotificationAccountBlock `json:"block"`
	Reason    string                   `json:"reason"`
}

type NotificationAccountChargeOff string

const (
	NotificationAccountChargeOffSuspended = NotificationAccountChargeOff("chargeOffSuspended")
	NotificationAccountChargeOffReached   = NotificationAccountChargeOff("chargeOffReached")
)

type AccountChargeOffNotification struct {
	AccountID AccountBankID                `json:"bankAccountId"`
	ChargeOff NotificationAccountChargeOff `json:"chargeOff"`
	Reason    string                       `json:"reason"`
}

type NotificationCardReissueReason string

// Card specific notification
type NewCardNotification struct {
	CardID    CardBankID                     `json:"cardId"`
	AccountID AccountBankID                  `json:"bankAccountId"`
	Reason    *NotificationCardReissueReason `json:"reason"`
}

const (
	NotificationCardReissueReasonCompromised = NotificationCardReissueReason("compromised")
	NotificationCardReissueReasonReplaced    = NotificationCardReissueReason("replaced")
	NotificationCardReissueReasonReissued    = NotificationCardReissueReason("reissued")
	NotificationCardReissueReasonUpgraded    = NotificationCardReissueReason("upgraded")
	NotificationCardReissueReasonRenewal     = NotificationCardReissueReason("renewed")
	NotificationCardReissueReasonNameChange  = NotificationCardReissueReason("nameChanged")
)

type CardReissueNotification struct {
	CardID    CardBankID                     `json:"cardId"`
	AccountID AccountBankID                  `json:"bankAccountId"`
	Reason    *NotificationCardReissueReason `json:"reason"`
}

type NotificationCardStatus string

const (
	NotificationCardStatusPrinting     = NotificationCardStatus("printing")
	NotificationCardStatusShipped      = NotificationCardStatus("shipped")
	NotificationCardStatusActivated    = NotificationCardStatus("activated")
	NotificationCardStatusCanceled     = NotificationCardStatus("canceled")
	NotificationCardStatusBlocked      = NotificationCardStatus("blocked")
	NotificationCardStatusUnblocked    = NotificationCardStatus("unblocked")
	NotificationCardStatusReissued     = NotificationCardStatus("reissued")
	NotificationCardStatusLimitChanged = NotificationCardStatus("limitChanged")
)

type NotificationCardBlockReason string

const (
	NotificationCardBlockReasonReissue             = NotificationCardBlockReason("reissue")
	NotificationCardBlockReasonFraud               = NotificationCardBlockReason("fraud")
	NotificationCardBlockReasonSuspended           = NotificationCardBlockReason("suspended")
	NotificationCardBlockReasonBadAddress          = NotificationCardBlockReason("badAddress")
	NotificationCardBlockReasonMailReturned        = NotificationCardBlockReason("mailReturned")
	NotificationCardBlockReasonLostStolen          = NotificationCardBlockReason("lostStolen")
	NotificationCardBlockReasonCompromised         = NotificationCardBlockReason("compromised")
	NotificationCardBlockReasonLostATM             = NotificationCardBlockReason("lostATM")
	NotificationCardBlockReasonSecurityBlocked     = NotificationCardBlockReason("securityBlocked")
	NotificationCardBlockReasonFraudBlocked        = NotificationCardBlockReason("fraudBlocked")
	NotificationCardBlockReasonInvestigatorBlocked = NotificationCardBlockReason("investigatorBlocked")
	NotificationCardBlockReasonBadDebt             = NotificationCardBlockReason("badDebt")
	NotificationCardBlockReasonTemporaryBlock      = NotificationCardBlockReason("temporaryBlock")
)

type CardBlockNotification struct {
	CardID    CardBankID                  `json:"cardId"`
	AccountID AccountBankID               `json:"bankAccountId"`
	Status    NotificationCardStatus      `json:"status"`
	Reason    NotificationCardBlockReason `json:"reason"`
}

type CardStatusNotification struct {
	CardID    CardBankID             `json:"cardId"`
	AccountID AccountBankID          `json:"bankAccountId"`
	Status    NotificationCardStatus `json:"status"`
}

// NACHA return code
type NotificationMoneyTransferReasonCode string

// Money transfer
type MoneyTransferStatusNotification struct {
	MoneyTransferID        MoneyTransferBankID                  `json:"moneyTransferId"`
	OriginAccountID        MoneyTransferAccountBankID           `json:"originAccountId"`
	OriginAccountType      MoneyTransferAccountType             `json:"originAccountType"`
	DestinationAccountID   MoneyTransferAccountBankID           `json:"destinationAccountId"`
	DestinationAccountType MoneyTransferAccountType             `json:"destinationAccountType"`
	Amount                 float64                              `json:"amount"`
	Currency               Currency                             `json:"currency,omitempty"`
	Status                 MoneyTransferStatus                  `json:"status"`
	BankStatus             string                               `json:"bank_status"`
	StatusUpdated          time.Time                            `json:"statusUpdated"`
	ReasonCode             *NotificationMoneyTransferReasonCode `json:"reasonCode"`
	ReasonDescription      *string                              `json:"reasonDescription"`
}

type MoneyTransferCorrectedNotification struct {
	MoneyTransferID MoneyTransferBankID `json:"moneyTransferId"`
}

type NotificationID string

func (id NotificationID) String() string {
	return string(id)
}

type NotificationType string
type NotificationAttribute string
type NotificationAction string

type NotificationEntityID string

func (id NotificationEntityID) String() string {
	return string(id)
}

type NotificationEntityType string

const (
	NotificationEntityTypeConsumer = NotificationEntityType("consumer")
	NotificationEntityTypeBusiness = NotificationEntityType("business")
	NotificationEntityTypeMember   = NotificationEntityType("member")
)

type SourceID string

type Notification struct {
	ID          NotificationID         `json:"id" db:"id"`
	EntityID    NotificationEntityID   `json:"entityId" db:"entity_id"`
	EntityType  NotificationEntityType `json:"entityType" db:"entity_type"`
	BankName    ProviderName           `json:"bankName" db:"bank_name"`
	SourceID    SourceID               `json:"sourceId" db:"source_id"`
	Type        NotificationType       `json:"type" db:"notification_type"`
	Action      NotificationAction     `json:"action" db:"notification_action"`
	Attribute   *NotificationAttribute `json:"attribute" db:"notification_attribute"`
	Version     string                 `json:"version" db:"notification_version"`
	SendCounter int                    `json:"-" db:"send_counter"`
	Created     time.Time              `json:"created" db:"created"`
	Data        types.JSONText         `json:"data" db:"notification_data"`
}
