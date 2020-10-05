package bbva

import (
	"encoding/json"
	"time"

	"github.com/wiseco/core-platform/partner/bank"
	"github.com/wiseco/core-platform/partner/bank/data"
)

type AccountNotification struct {
	UserID             string                    `json:"user_id"`
	AccountID          string                    `json:"account_id"`
	StatusChange       AccountNotificationStatus `json:"status_change"`
	StatusChangeReason string                    `json:"status_change_reason"`
}

type AccountNotificationStatus string

const (
	// Account Opened
	AccountNotificationStatusOpened = AccountNotificationStatus("account_opened")

	// Account Data Change
	AccountNotificationStatusDataChange = AccountNotificationStatus("account_data_change")

	// Participant Change
	AccountNotificationStatusParticipantAdded   = AccountNotificationStatus("participant_added")
	AccountNotificationStatusParticipantDeleted = AccountNotificationStatus("participant_deleted")

	// Block Change
	AccountNotificationStatusCreditBlockAdded   = AccountNotificationStatus("credits_block_added")
	AccountNotificationStatusCreditBlockRemoved = AccountNotificationStatus("credits_block_removed")
	AccountNotificationStatusDebitBlockAdded    = AccountNotificationStatus("debits_block_added")
	AccountNotificationStatusDebitBlockRemoved  = AccountNotificationStatus("debits_block_removed")
	AccountNotificationStatusCheckBlockAdded    = AccountNotificationStatus("checks_block_added")
	AccountNotificationStatusCheckBlockRemoved  = AccountNotificationStatus("checks_block_removed")

	// Status Change
	AccountNotificationStatusActive                   = AccountNotificationStatus("active")
	AccountNotificationStatusActiveDormant            = AccountNotificationStatus("active_but_dormant")
	AccountNotificationStatusActiveAbandoned          = AccountNotificationStatus("active_but_abandoned")
	AccountNotificationStatusActivePendingEscheatment = AccountNotificationStatus("active_but_pending_escheatment")
	AccountNotificationStatusInactive                 = AccountNotificationStatus("inactive")
	AccountNotificationStatusClosed                   = AccountNotificationStatus("closed")
	AccountNotificationStatusClosedByCorpSecurity     = AccountNotificationStatus("closed_by_corporate_security")
	AccountNotificationStatusClosedFraud              = AccountNotificationStatus("closed_for_fraud")
	AccountNotificationStatusClosedChargeOff          = AccountNotificationStatus("closed_for_charge_off")
	AccountNotificationStatusClosedEscheatment        = AccountNotificationStatus("closed_for_escheatment")

	// Charge Off Changes
	AccountNotificationStatusChargeOffReached       = AccountNotificationStatus("charge_off_date_reached")
	AccountNotificationStatusChargeOffDateSuspended = AccountNotificationStatus("charge_off_date_suspended")
)

var partnerAccountActionNotificationTo = map[AccountNotificationStatus]bank.NotificationAction{
	AccountNotificationStatusOpened:                   bank.NotificationActionOpen,
	AccountNotificationStatusDataChange:               bank.NotificationActionUpdate,
	AccountNotificationStatusParticipantAdded:         bank.NotificationActionAdd,
	AccountNotificationStatusParticipantDeleted:       bank.NotificationActionDelete,
	AccountNotificationStatusCreditBlockAdded:         bank.NotificationActionAdd,
	AccountNotificationStatusCreditBlockRemoved:       bank.NotificationActionRemove,
	AccountNotificationStatusDebitBlockAdded:          bank.NotificationActionAdd,
	AccountNotificationStatusDebitBlockRemoved:        bank.NotificationActionRemove,
	AccountNotificationStatusCheckBlockAdded:          bank.NotificationActionAdd,
	AccountNotificationStatusCheckBlockRemoved:        bank.NotificationActionRemove,
	AccountNotificationStatusActive:                   bank.NotificationActionUpdate,
	AccountNotificationStatusActiveDormant:            bank.NotificationActionUpdate,
	AccountNotificationStatusActiveAbandoned:          bank.NotificationActionUpdate,
	AccountNotificationStatusActivePendingEscheatment: bank.NotificationActionUpdate,
	AccountNotificationStatusInactive:                 bank.NotificationActionUpdate,
	AccountNotificationStatusClosed:                   bank.NotificationActionUpdate,
	AccountNotificationStatusClosedByCorpSecurity:     bank.NotificationActionUpdate,
	AccountNotificationStatusClosedFraud:              bank.NotificationActionUpdate,
	AccountNotificationStatusClosedChargeOff:          bank.NotificationActionUpdate,
	AccountNotificationStatusClosedEscheatment:        bank.NotificationActionUpdate,
	AccountNotificationStatusChargeOffReached:         bank.NotificationActionUpdate,
	AccountNotificationStatusChargeOffDateSuspended:   bank.NotificationActionUpdate,
}

func (s *notificationService) processAccountOpenNotification(n Notification, an AccountNotification, e *NotificationEntity) error {
	n.EventType = EventTypeAccountCreate

	nt, ok := partnerNotificationTypeTo[n.Type]
	if !ok {
		return bank.NewErrorFromCode(bank.ErrorCodeInvalidNotificationType)
	}

	action, ok := partnerAccountActionNotificationTo[an.StatusChange]
	if !ok {
		return bank.NewErrorFromCode(bank.ErrorCodeInvalidNotificationAction)
	}

	opened, err := time.Parse("2006-01-02", an.StatusChangeReason)
	if err != nil {
		return bank.NewErrorFromCode(bank.ErrorCodeInvalidAccountOpenDate)
	}

	b, err := json.Marshal(
		bank.AccountOpenedNotification{
			AccountID: bank.AccountBankID(an.AccountID),
			Opened:    &opened,
		},
	)
	if err != nil {
		return err
	}

	pn := data.NotificationCreate{
		EntityID:   bank.NotificationEntityID(e.EntityID),
		EntityType: bank.NotificationEntityType(e.EntityType),
		BankName:   bank.ProviderNameBBVA,
		SourceID:   bank.SourceID(n.notificationID()),
		Type:       nt,
		Action:     action,
		Version:    n.Version,
		Created:    n.SentDate,
		Data:       b,
	}

	return s.forwardNotification(pn)
}

func (s *notificationService) processAccountDataChangeNotification(n Notification, an AccountNotification, e *NotificationEntity) error {
	n.EventType = EventTypeAccountUpdate

	nt, ok := partnerNotificationTypeTo[n.Type]
	if !ok {
		return bank.NewErrorFromCode(bank.ErrorCodeInvalidNotificationType)
	}

	action, ok := partnerAccountActionNotificationTo[an.StatusChange]
	if !ok {
		return bank.NewErrorFromCode(bank.ErrorCodeInvalidNotificationAction)
	}

	b, err := json.Marshal(
		bank.AccountDataChangeNotification{
			AccountID: bank.AccountBankID(an.AccountID),
		},
	)
	if err != nil {
		return err
	}

	attr := bank.NotificationAttributeParticipant
	pn := data.NotificationCreate{
		EntityID:   bank.NotificationEntityID(e.EntityID),
		EntityType: bank.NotificationEntityType(e.EntityType),
		BankName:   bank.ProviderNameBBVA,
		SourceID:   bank.SourceID(n.notificationID()),
		Type:       nt,
		Attribute:  &attr,
		Action:     action,
		Version:    n.Version,
		Created:    n.SentDate,
		Data:       b,
	}

	return s.forwardNotification(pn)
}

func (s *notificationService) processAccountParticipantNotification(n Notification, an AccountNotification, e *NotificationEntity) error {
	n.EventType = EventTypeAccountUpdate

	nt, ok := partnerNotificationTypeTo[n.Type]
	if !ok {
		return bank.NewErrorFromCode(bank.ErrorCodeInvalidNotificationType)
	}

	action, ok := partnerAccountActionNotificationTo[an.StatusChange]
	if !ok {
		return bank.NewErrorFromCode(bank.ErrorCodeInvalidNotificationAction)
	}

	// TODO: Need to handle member users going forward
	pe, err := notificationEntityFromCustomerID(an.StatusChangeReason)
	if err != nil || pe.EntityType != bank.NotificationEntityTypeConsumer {
		return bank.NewErrorFromCode(bank.ErrorCodeInvalidEntityID)
	}

	b, err := json.Marshal(
		bank.AccountParticipantNotification{
			AccountID:  bank.AccountBankID(an.AccountID),
			EntityID:   pe.EntityID,
			EntityType: pe.EntityType,
		},
	)
	if err != nil {
		return err
	}

	attr := bank.NotificationAttributeParticipant
	pn := data.NotificationCreate{
		EntityID:   bank.NotificationEntityID(e.EntityID),
		EntityType: bank.NotificationEntityType(e.EntityType),
		BankName:   bank.ProviderNameBBVA,
		SourceID:   bank.SourceID(n.notificationID()),
		Type:       nt,
		Attribute:  &attr,
		Action:     action,
		Version:    n.Version,
		Created:    n.SentDate,
		Data:       b,
	}

	return s.forwardNotification(pn)
}

var partnerNotificationBlockStatusTo = map[AccountNotificationStatus]bank.NotificationAccountBlock{
	AccountNotificationStatusCreditBlockAdded:   bank.NotificationAccountBlockCreditAdded,
	AccountNotificationStatusCreditBlockRemoved: bank.NotificationAccountBlockCreditRemoved,
	AccountNotificationStatusDebitBlockAdded:    bank.NotificationAccountBlockDebitAdded,
	AccountNotificationStatusDebitBlockRemoved:  bank.NotificationAccountBlockDebitRemoved,
	AccountNotificationStatusCheckBlockAdded:    bank.NotificationAccountBlockCheckAdded,
	AccountNotificationStatusCheckBlockRemoved:  bank.NotificationAccountBlockCheckRemoved,
}

func (s *notificationService) processAccountBlockNotification(n Notification, an AccountNotification, e *NotificationEntity) error {
	n.EventType = EventTypeAccountBlock

	nt, ok := partnerNotificationTypeTo[n.Type]
	if !ok {
		return bank.NewErrorFromCode(bank.ErrorCodeInvalidNotificationType)
	}

	action, ok := partnerAccountActionNotificationTo[an.StatusChange]
	if !ok {
		return bank.NewErrorFromCode(bank.ErrorCodeInvalidNotificationAction)
	}

	block, ok := partnerNotificationBlockStatusTo[an.StatusChange]
	if !ok {
		return bank.NewErrorFromCode(bank.ErrorCodeInvalidAccountStatus)
	}

	b, err := json.Marshal(
		bank.AccountBlockNotification{
			AccountID: bank.AccountBankID(an.AccountID),
			Block:     block,
			Reason:    an.StatusChangeReason,
		},
	)
	if err != nil {
		return err
	}

	attr := bank.NotificationAttributeBlock
	pn := data.NotificationCreate{
		EntityID:   bank.NotificationEntityID(e.EntityID),
		EntityType: bank.NotificationEntityType(e.EntityType),
		BankName:   bank.ProviderNameBBVA,
		SourceID:   bank.SourceID(n.notificationID()),
		Type:       nt,
		Attribute:  &attr,
		Action:     action,
		Version:    n.Version,
		Created:    n.SentDate,
		Data:       b,
	}

	return s.forwardNotification(pn)
}

var partnerNotificationAccountStatusTo = map[AccountNotificationStatus]bank.NotificationAccountStatus{
	AccountNotificationStatusActive:                   bank.NotificationAccountStatusActive,
	AccountNotificationStatusActiveDormant:            bank.NotificationAccountStatusDormant,
	AccountNotificationStatusActiveAbandoned:          bank.NotificationAccountStatusAbandoned,
	AccountNotificationStatusActivePendingEscheatment: bank.NotificationAccountStatusPendingEscheatment,
	AccountNotificationStatusInactive:                 bank.NotificationAccountStatusInactive,
	AccountNotificationStatusClosed:                   bank.NotificationAccountStatusClosed,
	AccountNotificationStatusClosedByCorpSecurity:     bank.NotificationAccountStatusClosedBank,
	AccountNotificationStatusClosedFraud:              bank.NotificationAccountStatusClosedFraud,
	AccountNotificationStatusClosedChargeOff:          bank.NotificationAccountStatusChargeOff,
	AccountNotificationStatusClosedEscheatment:        bank.NotificationAccountStatusEscheated,
}

func (s *notificationService) processAccountStatusNotification(n Notification, an AccountNotification, e *NotificationEntity) error {
	n.EventType = EventTypeAccountStatus

	nt, ok := partnerNotificationTypeTo[n.Type]
	if !ok {
		return bank.NewErrorFromCode(bank.ErrorCodeInvalidNotificationType)
	}

	action, ok := partnerAccountActionNotificationTo[an.StatusChange]
	if !ok {
		return bank.NewErrorFromCode(bank.ErrorCodeInvalidNotificationAction)
	}

	status, ok := partnerNotificationAccountStatusTo[an.StatusChange]
	if !ok {
		return bank.NewErrorFromCode(bank.ErrorCodeInvalidAccountStatus)
	}

	b, err := json.Marshal(
		bank.AccountStatusNotification{
			AccountID: bank.AccountBankID(an.AccountID),
			Status:    status,
			Reason:    an.StatusChangeReason,
		},
	)
	if err != nil {
		return err
	}

	attr := bank.NotificationAttributeStatus
	pn := data.NotificationCreate{
		EntityID:   bank.NotificationEntityID(e.EntityID),
		EntityType: bank.NotificationEntityType(e.EntityType),
		BankName:   bank.ProviderNameBBVA,
		SourceID:   bank.SourceID(n.notificationID()),
		Type:       nt,
		Attribute:  &attr,
		Action:     action,
		Version:    n.Version,
		Created:    n.SentDate,
		Data:       b,
	}

	return s.forwardNotification(pn)
}

var partnerNotificationAccountChargeOffTo = map[AccountNotificationStatus]bank.NotificationAccountChargeOff{
	AccountNotificationStatusChargeOffReached:       bank.NotificationAccountChargeOffSuspended,
	AccountNotificationStatusChargeOffDateSuspended: bank.NotificationAccountChargeOffReached,
}

func (s *notificationService) processAccountChargeOffNotification(n Notification, an AccountNotification, e *NotificationEntity) error {
	n.EventType = EventTypeAccountChargeOff

	nt, ok := partnerNotificationTypeTo[n.Type]
	if !ok {
		return bank.NewErrorFromCode(bank.ErrorCodeInvalidNotificationType)
	}

	action, ok := partnerAccountActionNotificationTo[an.StatusChange]
	if !ok {
		return bank.NewErrorFromCode(bank.ErrorCodeInvalidNotificationAction)
	}

	chargeOff, ok := partnerNotificationAccountChargeOffTo[an.StatusChange]
	if !ok {
		return bank.NewErrorFromCode(bank.ErrorCodeInvalidAccountStatus)
	}

	b, err := json.Marshal(
		bank.AccountChargeOffNotification{
			AccountID: bank.AccountBankID(an.AccountID),
			ChargeOff: chargeOff,
			Reason:    an.StatusChangeReason,
		},
	)
	if err != nil {
		return err
	}

	attr := bank.NotificationAttributeChargeOff
	pn := data.NotificationCreate{
		EntityID:   bank.NotificationEntityID(e.EntityID),
		EntityType: bank.NotificationEntityType(e.EntityType),
		BankName:   bank.ProviderNameBBVA,
		SourceID:   bank.SourceID(n.notificationID()),
		Type:       nt,
		Attribute:  &attr,
		Action:     action,
		Version:    n.Version,
		Created:    n.SentDate,
		Data:       b,
	}

	return s.forwardNotification(pn)
}

func (s *notificationService) processAccountNotificationMessage(n Notification) error {
	if n.Reason != NotificationReasonStatusChange {
		return bank.NewErrorFromCode(bank.ErrorCodeProcessTransactionNotification)
	}

	var an AccountNotification
	err := n.unmarshalData(&an)
	if err != nil {
		return &bank.Error{
			RawError: err,
			Code:     bank.ErrorCodeProcessTransactionNotification,
		}
	}

	e, err := notificationEntityFromCustomerID(an.UserID)
	if err != nil {
		return bank.NewErrorFromCode(bank.ErrorCodeInvalidEntityID)
	}

	switch an.StatusChange {
	case AccountNotificationStatusOpened:
		return s.processAccountOpenNotification(n, an, e)
	case AccountNotificationStatusDataChange:
		return s.processAccountDataChangeNotification(n, an, e)
	case AccountNotificationStatusParticipantAdded, AccountNotificationStatusParticipantDeleted:
		return s.processAccountParticipantNotification(n, an, e)
	case
		AccountNotificationStatusCreditBlockAdded,
		AccountNotificationStatusCreditBlockRemoved,
		AccountNotificationStatusDebitBlockAdded,
		AccountNotificationStatusDebitBlockRemoved,
		AccountNotificationStatusCheckBlockAdded,
		AccountNotificationStatusCheckBlockRemoved:
		return s.processAccountBlockNotification(n, an, e)
	case
		AccountNotificationStatusActive,
		AccountNotificationStatusActiveDormant,
		AccountNotificationStatusActiveAbandoned,
		AccountNotificationStatusActivePendingEscheatment,
		AccountNotificationStatusInactive,
		AccountNotificationStatusClosed,
		AccountNotificationStatusClosedByCorpSecurity,
		AccountNotificationStatusClosedFraud,
		AccountNotificationStatusClosedChargeOff,
		AccountNotificationStatusClosedEscheatment:
		return s.processAccountStatusNotification(n, an, e)
	case AccountNotificationStatusChargeOffReached, AccountNotificationStatusChargeOffDateSuspended:
		return s.processAccountChargeOffNotification(n, an, e)
	}

	return bank.NewErrorFromCode(bank.ErrorCodeInvalidNotificationFormat)
}
