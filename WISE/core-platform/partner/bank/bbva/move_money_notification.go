package bbva

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/wiseco/core-platform/partner/bank"
	"github.com/wiseco/core-platform/partner/bank/data"
)

// NACHA return code
type StatusChangeReasonCode string

type MoveMoneyNotification struct {
	MoveMoneyID             string                  `json:"move_money_id"`
	UserID                  string                  `json:"user_id"`
	OriginAccount           string                  `json:"origin_account"`
	DestinationAccount      string                  `json:"destination_account"`
	Amount                  string                  `json:"amount"`
	TransactionStatus       MoveMoneyStatus         `json:"transaction_status"`
	StatusDate              time.Time               `json:"date_last_status_change"`
	StatusChangeReasonCode  *StatusChangeReasonCode `json:"status_change_reason_code"`
	StatusChangeDescription *string                 `json:"status_change_reason_description"`
}

func (s *notificationService) processMoveMoneyStatusNotification(n Notification, d MoveMoneyNotification, e *NotificationEntity) error {
	n.EventType = EventTypePaymentStatusChange

	nt, ok := partnerNotificationTypeTo[n.Type]
	if !ok {
		return bank.NewErrorFromCode(bank.ErrorCodeInvalidNotificationType)
	}

	status, ok := partnerMoneyTransferStatusTo[d.TransactionStatus]
	if !ok {
		return bank.NewErrorFromCode(bank.ErrorCodeInvalidMoneyTransferStatus)
	}

	var reasonCode *bank.NotificationMoneyTransferReasonCode
	if d.StatusChangeReasonCode != nil {
		rc := bank.NotificationMoneyTransferReasonCode(*d.StatusChangeReasonCode)
		reasonCode = &rc
	}

	amount, err := strconv.ParseFloat(d.Amount, 64)
	if err != nil {
		return &bank.Error{
			RawError: err,
			Code:     bank.ErrorCodeProcessTransactionNotification,
		}
	}

	mn := bank.MoneyTransferStatusNotification{
		MoneyTransferID:        bank.MoneyTransferBankID(d.MoveMoneyID),
		OriginAccountID:        bank.MoneyTransferAccountBankID(d.OriginAccount),
		OriginAccountType:      partnerMoneyTransferAccountTypeFromID(d.OriginAccount),
		DestinationAccountID:   bank.MoneyTransferAccountBankID(d.DestinationAccount),
		DestinationAccountType: partnerMoneyTransferAccountTypeFromID(d.DestinationAccount),
		Amount:                 amount,
		Currency:               bank.CurrencyUSD,
		Status:                 status,
		BankStatus:             string(d.TransactionStatus),
		StatusUpdated:          d.StatusDate,
		ReasonCode:             reasonCode,
		ReasonDescription:      d.StatusChangeDescription,
	}

	b, err := json.Marshal(mn)
	if err != nil {
		return err
	}

	pn := data.NotificationCreate{
		EntityID:   bank.NotificationEntityID(e.EntityID),
		EntityType: bank.NotificationEntityType(e.EntityType),
		BankName:   bank.ProviderNameBBVA,
		SourceID:   bank.SourceID(n.notificationID()),
		Type:       nt,
		Action:     bank.NotificationActionUpdate,
		Version:    n.Version,
		Created:    n.SentDate,
		Data:       b,
	}

	return s.forwardNotification(pn)
}

func (s *notificationService) processMoveMoneyCorrectedNotification(n Notification, d MoveMoneyNotification, e *NotificationEntity) error {
	n.EventType = EventTypePaymentCorrected

	nt, ok := partnerNotificationTypeTo[n.Type]
	if !ok {
		return bank.NewErrorFromCode(bank.ErrorCodeInvalidNotificationType)
	}

	mn := bank.MoneyTransferCorrectedNotification{
		MoneyTransferID: bank.MoneyTransferBankID(d.MoveMoneyID),
	}

	b, err := json.Marshal(mn)
	if err != nil {
		return err
	}

	pn := data.NotificationCreate{
		EntityID:   bank.NotificationEntityID(e.EntityID),
		EntityType: bank.NotificationEntityType(e.EntityType),
		BankName:   bank.ProviderNameBBVA,
		SourceID:   bank.SourceID(n.notificationID()),
		Type:       nt,
		Action:     bank.NotificationActionCorrected,
		Version:    n.Version,
		Created:    n.SentDate,
		Data:       b,
	}

	return s.forwardNotification(pn)
}

func (s *notificationService) processMoveMoneyNotificationMessage(n Notification) error {
	var d MoveMoneyNotification
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

	switch n.Reason {
	case NotificationReasonStatusChange:
		return s.processMoveMoneyStatusNotification(n, d, e)
	case NotificationReasonCorrectedData:
		return s.processMoveMoneyCorrectedNotification(n, d, e)
	}

	return bank.NewErrorFromCode(bank.ErrorCodeInvalidNotificationFormat)
}
