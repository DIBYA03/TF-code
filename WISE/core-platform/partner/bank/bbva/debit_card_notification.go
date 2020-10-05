package bbva

import (
	"encoding/json"

	"github.com/wiseco/core-platform/partner/bank"
	"github.com/wiseco/core-platform/partner/bank/data"
)

type CardNotificationStatus string

const (
	CardNotificationStatusNewCard      = CardNotificationStatus("new_card")
	CardNotificationStatusActivated    = CardNotificationStatus("activated")
	CardNotificationStatusBlocked      = CardNotificationStatus("blocked")
	CardNotificationStatusUnblocked    = CardNotificationStatus("unblocked")
	CardNotificationStatusReissued     = CardNotificationStatus("reissued")
	CardNotificationStatusCanceled     = CardNotificationStatus("canceled")
	CardNotificationStatusLimitChanged = CardNotificationStatus("limit_changed")
)

var partnerCardNotificationStatusTo = map[CardNotificationStatus]bank.NotificationCardStatus{
	CardNotificationStatusNewCard:      bank.NotificationCardStatusShipped,
	CardNotificationStatusActivated:    bank.NotificationCardStatusActivated,
	CardNotificationStatusBlocked:      bank.NotificationCardStatusBlocked,
	CardNotificationStatusUnblocked:    bank.NotificationCardStatusUnblocked,
	CardNotificationStatusReissued:     bank.NotificationCardStatusReissued,
	CardNotificationStatusCanceled:     bank.NotificationCardStatusCanceled,
	CardNotificationStatusLimitChanged: bank.NotificationCardStatusLimitChanged,
}

type DebitCardNotification struct {
	UserID             string                 `json:"user_id"`
	AccountID          string                 `json:"account_id"`
	CardID             string                 `json:"card_id"`
	StatusChange       CardNotificationStatus `json:"status_change"`
	StatusChangeReason string                 `json:"status_change_reason"`
}

var partnerCardActionNotificationTo = map[CardNotificationStatus]bank.NotificationAction{
	CardNotificationStatusNewCard:      bank.NotificationActionUpdate,
	CardNotificationStatusActivated:    bank.NotificationActionActivate,
	CardNotificationStatusBlocked:      bank.NotificationActionUpdate,
	CardNotificationStatusUnblocked:    bank.NotificationActionUpdate,
	CardNotificationStatusReissued:     bank.NotificationActionReissue,
	CardNotificationStatusCanceled:     bank.NotificationActionCancel,
	CardNotificationStatusLimitChanged: bank.NotificationActionUpdate,
}

type ReissueDebitCardNotificationReason string

const (
	ReissueDebitCardNotificationReasonOriginalCompromised = ReissueDebitCardNotificationReason("original_compromised")
	ReissueDebitCardNotificationReasonCardReplacement     = ReissueDebitCardNotificationReason("hot_card_replacement")
	ReissueDebitCardNotificationReasonBulkReissue         = ReissueDebitCardNotificationReason("bulk_reissue")
	ReissueDebitCardNotificationReasonProductUpgrade      = ReissueDebitCardNotificationReason("product_upgrade")
	ReissueDebitCardNotificationReasonRenewal             = ReissueDebitCardNotificationReason("renewal")
	ReissueDebitCardNotificationReasonNameChange          = ReissueDebitCardNotificationReason("name_change")
)

var partnerReissueReasonNotificationTo = map[ReissueDebitCardNotificationReason]bank.NotificationCardReissueReason{
	ReissueDebitCardNotificationReasonOriginalCompromised: bank.NotificationCardReissueReasonCompromised,
	ReissueDebitCardNotificationReasonCardReplacement:     bank.NotificationCardReissueReasonReplaced,
	ReissueDebitCardNotificationReasonBulkReissue:         bank.NotificationCardReissueReasonReissued,
	ReissueDebitCardNotificationReasonProductUpgrade:      bank.NotificationCardReissueReasonUpgraded,
	ReissueDebitCardNotificationReasonRenewal:             bank.NotificationCardReissueReasonRenewal,
	ReissueDebitCardNotificationReasonNameChange:          bank.NotificationCardReissueReasonNameChange,
}

func (s *notificationService) processNewDebitCardNotification(n Notification, d DebitCardNotification, e *NotificationEntity) error {
	nt, ok := partnerNotificationTypeTo[n.Type]
	if !ok {
		return bank.NewErrorFromCode(bank.ErrorCodeInvalidNotificationType)
	}

	action, ok := partnerCardActionNotificationTo[CardNotificationStatus(d.StatusChange)]
	if !ok {
		return bank.NewErrorFromCode(bank.ErrorCodeInvalidNotificationAction)
	}

	var statusReason *bank.NotificationCardReissueReason
	reason, ok := partnerReissueReasonNotificationTo[ReissueDebitCardNotificationReason(d.StatusChangeReason)]
	if ok {
		statusReason = &reason
	}

	cn := bank.NewCardNotification{
		CardID:    bank.CardBankID(d.CardID),
		AccountID: bank.AccountBankID(d.AccountID),
		Reason:    statusReason,
	}

	b, err := json.Marshal(cn)
	if err != nil {
		return err
	}

	pn := data.NotificationCreate{
		EntityID:   bank.NotificationEntityID(e.EntityID),
		EntityType: bank.NotificationEntityType(e.EntityType),
		BankName:   bank.ProviderNameBBVA,
		SourceID:   bank.SourceID(n.ID),
		Type:       nt,
		Action:     action,
		Version:    n.Version,
		Created:    n.SentDate,
		Data:       b,
	}

	return s.forwardNotification(pn)
}

func (s *notificationService) processReissueDebitCardNotification(n Notification, d DebitCardNotification, e *NotificationEntity) error {
	nt, ok := partnerNotificationTypeTo[n.Type]
	if !ok {
		return bank.NewErrorFromCode(bank.ErrorCodeInvalidNotificationType)
	}

	action, ok := partnerCardActionNotificationTo[d.StatusChange]
	if !ok {
		return bank.NewErrorFromCode(bank.ErrorCodeInvalidNotificationAction)
	}

	var statusReason *bank.NotificationCardReissueReason
	reason, ok := partnerReissueReasonNotificationTo[ReissueDebitCardNotificationReason(d.StatusChangeReason)]
	if ok {
		statusReason = &reason
	}

	cn := bank.CardReissueNotification{
		CardID:    bank.CardBankID(d.CardID),
		AccountID: bank.AccountBankID(d.AccountID),
		Reason:    statusReason,
	}

	b, err := json.Marshal(cn)
	if err != nil {
		return err
	}

	pn := data.NotificationCreate{
		EntityID:   bank.NotificationEntityID(e.EntityID),
		EntityType: bank.NotificationEntityType(e.EntityType),
		BankName:   bank.ProviderNameBBVA,
		SourceID:   bank.SourceID(n.ID),
		Type:       nt,
		Action:     action,
		Version:    n.Version,
		Created:    n.SentDate,
		Data:       b,
	}

	return s.forwardNotification(pn)
}

type NotificationCardBlockStatusReason string

const (
	NotificationCardBlockStatusReasonPotentialFraud      = NotificationCardBlockStatusReason("potential_fraud")
	NotificationCardBlockStatusReasonSuspended           = NotificationCardBlockStatusReason("suspended")
	NotificationCardBlockStatusReasonBadAddress          = NotificationCardBlockStatusReason("bad_address")
	NotificationCardBlockStatusReasonReturnedMail        = NotificationCardBlockStatusReason("returned_mail")
	NotificationCardBlockStatusReasonLostStolen          = NotificationCardBlockStatusReason("reported_lost_or_stolen")
	NotificationCardBlockStatusReasonReportedCompromised = NotificationCardBlockStatusReason("reported_compromised")
	NotificationCardBlockStatusReasonCapturedAtATM       = NotificationCardBlockStatusReason("captured_at_ATM”")
	NotificationCardBlockStatusReasonBlockedBySecurity   = NotificationCardBlockStatusReason("blocked_by_security")
	NotificationCardBlockStatusReasonBulkReissue         = NotificationCardBlockStatusReason("bulk_reissue")
	NotificationCardBlockStatusReasonBlockedBankFraud    = NotificationCardBlockStatusReason("blocked_by_bank_fraud_system")
	NotificationCardBlockStatusReasonBlockedInvestiagtor = NotificationCardBlockStatusReason("blocked_manually_by_investigator")
	NotificationCardBlockStatusReasonBadDebt             = NotificationCardBlockStatusReason("bad_debt")
	NotificationCardBlockStatusReasonTemporaryBlock      = NotificationCardBlockStatusReason("temporary_block")
)

var parterCardNotificationStatusReasonTo = map[NotificationCardBlockStatusReason]bank.NotificationCardBlockReason{
	NotificationCardBlockStatusReasonPotentialFraud:      bank.NotificationCardBlockReasonFraud,
	NotificationCardBlockStatusReasonSuspended:           bank.NotificationCardBlockReasonSuspended,
	NotificationCardBlockStatusReasonBadAddress:          bank.NotificationCardBlockReasonBadAddress,
	NotificationCardBlockStatusReasonReturnedMail:        bank.NotificationCardBlockReasonMailReturned,
	NotificationCardBlockStatusReasonLostStolen:          bank.NotificationCardBlockReasonLostStolen,
	NotificationCardBlockStatusReasonReportedCompromised: bank.NotificationCardBlockReasonCompromised,
	NotificationCardBlockStatusReasonCapturedAtATM:       bank.NotificationCardBlockReasonLostATM,
	NotificationCardBlockStatusReasonBlockedBySecurity:   bank.NotificationCardBlockReasonSecurityBlocked,
	NotificationCardBlockStatusReasonBulkReissue:         bank.NotificationCardBlockReasonReissue,
	NotificationCardBlockStatusReasonBlockedBankFraud:    bank.NotificationCardBlockReasonFraudBlocked,
	NotificationCardBlockStatusReasonBlockedInvestiagtor: bank.NotificationCardBlockReasonInvestigatorBlocked,
	NotificationCardBlockStatusReasonBadDebt:             bank.NotificationCardBlockReasonBadDebt,
	NotificationCardBlockStatusReasonTemporaryBlock:      bank.NotificationCardBlockReasonTemporaryBlock,
}

func (s *notificationService) processDebitCardBlockedNotification(n Notification, d DebitCardNotification, e *NotificationEntity) error {
	nt, ok := partnerNotificationTypeTo[n.Type]
	if !ok {
		return bank.NewErrorFromCode(bank.ErrorCodeInvalidNotificationType)
	}

	action, ok := partnerCardActionNotificationTo[d.StatusChange]
	if !ok {
		return bank.NewErrorFromCode(bank.ErrorCodeInvalidNotificationAction)
	}

	status, ok := partnerCardNotificationStatusTo[d.StatusChange]
	if !ok {
		return bank.NewErrorFromCode(bank.ErrorCodeInvalidNotificationStatus)
	}

	reason, ok := parterCardNotificationStatusReasonTo[NotificationCardBlockStatusReason(d.StatusChangeReason)]
	if !ok {
		return bank.NewErrorFromCode(bank.ErrorCodeInvalidCardStatus)
	}

	cn := bank.CardBlockNotification{
		CardID:    bank.CardBankID(d.CardID),
		AccountID: bank.AccountBankID(d.AccountID),
		Status:    status,
		Reason:    reason,
	}

	b, err := json.Marshal(cn)
	if err != nil {
		return err
	}

	pn := data.NotificationCreate{
		EntityID:   bank.NotificationEntityID(e.EntityID),
		EntityType: bank.NotificationEntityType(e.EntityType),
		BankName:   bank.ProviderNameBBVA,
		SourceID:   bank.SourceID(n.ID),
		Type:       nt,
		Action:     action,
		Version:    n.Version,
		Created:    n.SentDate,
		Data:       b,
	}

	return s.forwardNotification(pn)
}

func (s *notificationService) processDebitCardStatusChangeNotification(n Notification, d DebitCardNotification, e *NotificationEntity) error {
	nt, ok := partnerNotificationTypeTo[n.Type]
	if !ok {
		return bank.NewErrorFromCode(bank.ErrorCodeInvalidNotificationType)
	}

	action, ok := partnerCardActionNotificationTo[d.StatusChange]
	if !ok {
		return bank.NewErrorFromCode(bank.ErrorCodeInvalidNotificationAction)
	}

	status, ok := partnerCardNotificationStatusTo[d.StatusChange]
	if !ok {
		return bank.NewErrorFromCode(bank.ErrorCodeInvalidNotificationStatus)
	}

	cn := bank.CardStatusNotification{
		CardID:    bank.CardBankID(d.CardID),
		AccountID: bank.AccountBankID(d.AccountID),
		Status:    status,
	}

	b, err := json.Marshal(cn)
	if err != nil {
		return err
	}

	pn := data.NotificationCreate{
		EntityID:   bank.NotificationEntityID(e.EntityID),
		EntityType: bank.NotificationEntityType(e.EntityType),
		BankName:   bank.ProviderNameBBVA,
		SourceID:   bank.SourceID(n.ID),
		Type:       nt,
		Action:     action,
		Version:    n.Version,
		Created:    n.SentDate,
		Data:       b,
	}

	return s.forwardNotification(pn)
}

func (s *notificationService) processDebitCardNotification(n Notification) error {
	var d DebitCardNotification
	err := n.unmarshalData(&d)
	if err != nil {
		return &bank.Error{
			RawError: err,
			Code:     bank.ErrorCodeProcessTransactionNotification,
		}
	}

	e, err := notificationEntityFromCustomerID(d.UserID)
	if err != nil {
		return bank.NewErrorFromCode(bank.ErrorCodeInvalidEntityID)
	}

	switch d.StatusChange {
	case CardNotificationStatusReissued:
		return s.processReissueDebitCardNotification(n, d, e)
	case CardNotificationStatusBlocked:
		return s.processDebitCardBlockedNotification(n, d, e)
	case
		CardNotificationStatusNewCard, CardNotificationStatusActivated,
		CardNotificationStatusUnblocked, CardNotificationStatusCanceled,
		CardNotificationStatusLimitChanged:
		return s.processDebitCardStatusChangeNotification(n, d, e)
	}

	return bank.NewErrorFromCode(bank.ErrorCodeInvalidNotificationFormat)
}
