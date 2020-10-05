package bbva

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/wiseco/core-platform/partner/bank"
	"github.com/wiseco/core-platform/partner/bank/data"
	"github.com/wiseco/core-platform/shared"
)

type NotificationHandlerInfo struct {
	InQueue        shared.MessageQueue
	OutQueue       shared.MessageQueue
	StreamProvider shared.StreamProvider
}

type NotificationService interface {
	// Replay notifications
	ReplayByID(bank.NotificationID) error
	ReplayAll(time.Time) error

	// Handle notifications
	HandleNotifications(context.Context) error
}

type notificationService struct {
	info NotificationHandlerInfo
}

func NewNotificationService(info NotificationHandlerInfo) NotificationService {
	return &notificationService{info}
}

func (s *notificationService) ReplayByID(id bank.NotificationID) error {
	return errors.New("not implemented")
}

func (s *notificationService) ReplayAll(t time.Time) error {
	return errors.New("not implemented")
}

func (s *notificationService) HandleNotifications(ctx context.Context) error {
	err := s.info.InQueue.ReceiveMessages(ctx, s)
	if err != nil {
		return &bank.Error{
			RawError: err,
			Code:     bank.ErrorCodeHandleNotifications,
		}
	}

	return nil
}

func (s *notificationService) HandleMessage(_ context.Context, m shared.Message) error {
	if m.Body == nil {
		return bank.NewErrorFromCode(bank.ErrorCodeHandleNotificationMessage)
	}

	var n Notification
	err := json.Unmarshal(m.Body, &n)
	if err != nil {
		return err
	}

	raw := data.RawNotification{
		ID:       m.ID,
		BankName: bank.ProviderNameBBVA.String(),
		SourceID: n.notificationID(),
		Message:  m.Body,
	}
	err = data.NewNotificationService(bank.NewAPIRequest()).LogRaw(s.info.StreamProvider, []data.RawNotification{raw})
	if err != nil {
		return &bank.Error{
			RawError: err,
			Code:     bank.ErrorCodeHandleNotificationMessage,
		}
	}

	return s.processNotification(n)
}

var notificationTypeFromEventType = map[EventType]NotificationType{
	// Consumer
	EventTypeConsumerCreate:           NotificationTypeConsumer,
	EventTypeConsumerUpdate:           NotificationTypeConsumer,
	EventTypeConsumerContactCreate:    NotificationTypeConsumer,
	EventTypeConsumerContactUpdate:    NotificationTypeConsumer,
	EventTypeConsumerContactDelete:    NotificationTypeConsumer,
	EventTypeConsumerAddressCreate:    NotificationTypeConsumer,
	EventTypeConsumerAddressUpdate:    NotificationTypeConsumer,
	EventTypeConsumerAddressDelete:    NotificationTypeConsumer,
	EventTypeConsumerOccupationUpdate: NotificationTypeConsumer,
	EventTypeConsumerKYC:              NotificationTypeConsumer,

	// Business
	EventTypeBusinessCreate:          NotificationTypeBusiness,
	EventTypeBusinessUpdate:          NotificationTypeBusiness,
	EventTypeBusinessOwnerCreate:     NotificationTypeBusiness,
	EventTypeBusinessOwnerDelete:     NotificationTypeBusiness,
	EventTypeBusinessIndicatorUpdate: NotificationTypeBusiness,
	EventTypeBusinessContactCreate:   NotificationTypeBusiness,
	EventTypeBusinessContactUpdate:   NotificationTypeBusiness,
	EventTypeBusinessContactDelete:   NotificationTypeBusiness,
	EventTypeBusinessAddressCreate:   NotificationTypeBusiness,
	EventTypeBusinessAddressUpdate:   NotificationTypeBusiness,
	EventTypeBusinessAddressDelete:   NotificationTypeBusiness,
	EventTypeBusinessKYC:             NotificationTypeBusiness,

	// Account
	EventTypeAccountCreate:        NotificationTypeAccount,
	EventTypeAccountUpdate:        NotificationTypeAccount,
	EventTypeAccountStatus:        NotificationTypeAccount,
	EventTypeAccountBlock:         NotificationTypeAccount,
	EventTypeAccountChargeOff:     NotificationTypeAccount,
	EventTypeAccountChargeOffCard: NotificationTypeAccount,

	// Cards
	EventTypeCardCreate: NotificationTypeCard,
	EventTypeCardUpdate: NotificationTypeCard,

	// Money Transfer
	EventTypePaymentStatusChange: NotificationTypeMoveMoney,
	EventTypePaymentCorrected:    NotificationTypeMoveMoney,

	// Transactions
	EventTypeAuthorization: NotificationTypeTransaction,
	EventTypeFundhold:      NotificationTypeTransaction,
	EventTypeCardPosted:    NotificationTypeTransaction,
	EventTypeNonCardPosted: NotificationTypeTransaction,
}

type NotificationType string

const (
	NotificationTypeConsumer    = NotificationType("consumers")
	NotificationTypeBusiness    = NotificationType("business")
	NotificationTypeAccount     = NotificationType("accounts")
	NotificationTypeCard        = NotificationType("cards")
	NotificationTypeMoveMoney   = NotificationType("move_money")
	NotificationTypeTransaction = NotificationType("transactions")
)

var partnerNotificationTypeTo = map[NotificationType]bank.NotificationType{
	NotificationTypeConsumer:    bank.NotificationTypeConsumer,
	NotificationTypeBusiness:    bank.NotificationTypeBusiness,
	NotificationTypeAccount:     bank.NotificationTypeAccount,
	NotificationTypeCard:        bank.NotificationTypeCard,
	NotificationTypeMoveMoney:   bank.NotificationTypeMoneyTransfer,
	NotificationTypeTransaction: bank.NotificationTypeTransaction,
}

type Version string

const NotificationVersion = Version("1.0.0")

type Timestamp int64

func (t *Timestamp) UnmarshalJSON(data []byte) error {
	var ts string
	err := json.Unmarshal(data, &ts)
	if err == nil {
		tsnum, err := strconv.ParseInt(ts, 10, 64)
		if err != nil {
			return err
		}

		*t = Timestamp(tsnum)
	} else {
		var tsnum int64
		err := json.Unmarshal(data, &tsnum)
		if err != nil {
			return err
		}

		*t = Timestamp(tsnum)
	}

	return nil
}

type Notification struct {
	ID                string             `json:"notification_id"`
	Type              NotificationType   `json:"notification_type"`
	Version           string             `json:"notification_version"`
	Reason            NotificationReason `json:"notification_reason"`
	SentDate          time.Time          `json:"notification_sent_date"`
	EventID           string             `json:"eventId"`
	EventTypeFull     string             `json:"eventTypeFull"`
	EventType         EventType          `json:"eventType"`
	EventTypeVersion  Version            `json:"eventTypeVersion"`
	Subscriber        string             `json:"subscriber"`
	CustomerID        json.RawMessage    `json:"customerId"`
	CreationTimestamp Timestamp          `json:"creationTimestamp"`
	Data              json.RawMessage    `json:"notification_data"`
	Payload           json.RawMessage    `json:"payload"`
}

func (n *Notification) customerID() (string, error) {
	// Customer ID usually comes in a string
	var id string
	err := json.Unmarshal(n.CustomerID, &id)
	if err != nil {
		// or an array with a single string value
		var ids []string
		err = json.Unmarshal(n.CustomerID, &ids)
		if err != nil || len(ids) == 0 {
			return id, bank.NewErrorFromCode(bank.ErrorCodeInvalidEntityID)
		}

		id = ids[0]
		if notificationEntityTypeFromCustomerID(id) == NotificationEntityTypeNone {
			return id, bank.NewErrorFromCode(bank.ErrorCodeInvalidEntityID)
		}

		return id, nil
	}

	if notificationEntityTypeFromCustomerID(id) != NotificationEntityTypeNone {
		return id, nil
	}

	// Sometimes comes in an array in string
	var ids []string
	err = json.Unmarshal([]byte(id), &ids)
	if err != nil || len(ids) == 0 {
		return id, bank.NewErrorFromCode(bank.ErrorCodeInvalidEntityID)
	}

	id = ids[0]
	if notificationEntityTypeFromCustomerID(id) == NotificationEntityTypeNone {
		return id, bank.NewErrorFromCode(bank.ErrorCodeInvalidEntityID)
	}

	return id, nil
}

func (n *Notification) notificationID() string {
	if n.ID != "" {
		return n.ID
	}

	return n.EventID
}

func (n *Notification) unmarshalData(v interface{}) error {
	if len(n.Data) > 0 {
		return json.Unmarshal(n.Data, v)
	} else if len(n.Payload) > 0 {
		return json.Unmarshal(n.Payload, v)
	}

	return bank.NewErrorFromCode(bank.ErrorCodeNoPayload)
}

func (s *notificationService) processNotification(n Notification) error {
	// Extract event type and notificaion type
	if n.EventType != "" {
		n.EventTypeFull = string(n.EventType)
		n.EventType = EventType(strings.TrimPrefix(n.EventTypeFull, EventTypePrefix))
		nType, ok := notificationTypeFromEventType[n.EventType]
		if !ok {
			return fmt.Errorf("event type unknown (%s)", n.EventTypeFull)
		}

		n.Type = nType
	}

	if n.CreationTimestamp > 0 {
		n.SentDate = time.Unix(int64(n.CreationTimestamp/1000), 0)
	}

	if n.EventTypeVersion != "" {
		n.Version = string(n.EventTypeVersion)
	}

	switch n.Type {
	case NotificationTypeConsumer:
		return s.processConsumerNotificationMessage(n)
	case NotificationTypeBusiness:
		return s.processBusinessNotificationMessage(n)
	case NotificationTypeAccount:
		return s.processAccountNotificationMessage(n)
	case NotificationTypeCard:
		return s.processDebitCardNotification(n)
	case NotificationTypeMoveMoney:
		return s.processMoveMoneyNotificationMessage(n)
	case NotificationTypeTransaction:
		return s.processTransactionNotificationMessage(n)
	}

	return fmt.Errorf("event type unknown (%s)", n.Type)
}

type NotificationReason string

const (
	NotificationReasonStatusChange          = NotificationReason("status_change")
	NotificationReasonCorrectedData         = NotificationReason("corrected_data")
	NotificationReasonAuthorizationApproved = NotificationReason("authorization_approved")
	NotificationReasonAuthorizationDeclined = NotificationReason("authorization_declined")
	NotificationReasonAuthorizationReversal = NotificationReason("authorization_reversal")
	NotificationReasonFundsHoldSet          = NotificationReason("funds_hold_set")
	NotificationReasonFundsHoldReleased     = NotificationReason("funds_hold_released")
	NotificationReasonCardDebitPosted       = NotificationReason("card_debit_posted")
	NotificationReasonCardCreditPosted      = NotificationReason("card_credit_posted")
	NotificationReasonNonCardDebitPosted    = NotificationReason("noncard_debit_posted")
	NotificationReasonNonCardCreditPosted   = NotificationReason("noncard_credit_posted")
)

func (s *notificationService) forwardNotification(pn data.NotificationCreate) error {
	// Save notification to database
	notification, err := data.NewNotificationService(bank.NewAPIRequest()).Create(&pn)
	if err != nil {
		if perr, ok := err.(*bank.Error); ok && perr.Code == bank.ErrorCodeDuplicateNotification {
			notification, err = data.NewNotificationService(bank.NewAPIRequest()).GetBySourceID(pn.SourceID, pn.BankName)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	// Send to internal SQS queue
	b, err := json.Marshal(notification)
	if err != nil {
		return &bank.Error{
			RawError: err,
			Code:     bank.ErrorCodeProcessTransactionNotification,
		}
	}

	groupID := notification.EntityID.String()
	message := shared.Message{
		ID:      notification.ID.String(),
		GroupID: &groupID,
		Body:    json.RawMessage(b),
	}

	_, err = s.info.OutQueue.SendMessages([]shared.Message{message})
	if err != nil {
		return &bank.Error{
			RawError: err,
			Code:     bank.ErrorCodeProcessTransactionNotification,
		}
	}

	return data.NewNotificationService(bank.NewAPIRequest()).IncrementSend([]bank.NotificationID{notification.ID})
}

type NotificationEntity struct {
	EntityID   bank.NotificationEntityID   `json:"entityId"`
	EntityType bank.NotificationEntityType `json:"entityType"`
	BankID     string                      `json:"bankId"`
	KYCStatus  bank.KYCStatus              `json:"kycStatus"`
}

func notificationEntityFromCustomerID(id string) (*NotificationEntity, error) {
	entityType := notificationEntityTypeFromCustomerID(id)
	switch entityType {
	case NotificationEntityTypeConsumer:
		c, err := data.NewConsumerService(bank.NewAPIRequest(), bank.ProviderNameBBVA).GetByBankID(bank.ConsumerBankID(id))
		if err != nil {
			return nil, bank.NewErrorFromCode(bank.ErrorCodeInvalidEntityID)
		}

		return &NotificationEntity{
			EntityID:   bank.NotificationEntityID(c.ConsumerID),
			EntityType: bank.NotificationEntityTypeConsumer,
			BankID:     string(c.BankID),
			KYCStatus:  c.KYCStatus,
		}, nil
	case NotificationEntityTypeBusiness:
		b, err := data.NewBusinessService(bank.NewAPIRequest(), bank.ProviderNameBBVA).GetByBankID(bank.BusinessBankID(id))
		if err != nil {
			return nil, bank.NewErrorFromCode(bank.ErrorCodeInvalidEntityID)
		}

		return &NotificationEntity{
			EntityID:   bank.NotificationEntityID(b.BusinessID),
			EntityType: bank.NotificationEntityTypeBusiness,
			BankID:     string(b.BankID),
			KYCStatus:  b.KYCStatus,
		}, nil
	case NotificationEntityTypeMember:
		m, err := data.NewBusinessMemberService(bank.NewAPIRequest(), bank.ProviderNameBBVA).GetByBankID(bank.MemberBankID(id))
		if err != nil {
			return nil, bank.NewErrorFromCode(bank.ErrorCodeInvalidEntityID)
		}

		return &NotificationEntity{
			EntityID:   bank.NotificationEntityID(m.ConsumerBankID),
			EntityType: bank.NotificationEntityTypeMember,
			BankID:     string(m.ConsumerBankID),
			KYCStatus:  m.KYCStatus,
		}, nil
	}

	return nil, bank.NewErrorFromCode(bank.ErrorCodeInvalidEntityID)
}

type NotificationEntityType string

const NotificationEntityTypeNone = NotificationEntityType("")

const (
	NotificationEntityTypeConsumer = NotificationEntityType("consumer")
	NotificationEntityTypeBusiness = NotificationEntityType("business")
	NotificationEntityTypeMember   = NotificationEntityType("member")
)

const (
	EntityPrefixConsumer       = "CO-"
	EntityPrefixBusiness       = "NC-"
	EntityPrefixBusinessMember = "OW-"
)

func notificationEntityTypeFromCustomerID(customerID string) NotificationEntityType {
	if strings.HasPrefix(customerID, EntityPrefixConsumer) {
		return NotificationEntityTypeConsumer
	} else if strings.HasPrefix(customerID, EntityPrefixBusiness) {
		return NotificationEntityTypeBusiness
	} else if strings.HasPrefix(customerID, EntityPrefixBusinessMember) {
		return NotificationEntityTypeMember
	}

	return NotificationEntityTypeNone
}

const (
	AccountPrefix           = "AC-"
	RegisteredAccountPrefix = "RA-"
	CardPrefix              = "DC-"
	RegisteredCardPrefix    = "RC-"
)

type NotificationMoneyTransferAccountType string

const (
	NotificationMoneyTransferAccountTypeBank = NotificationMoneyTransferAccountType("bank")
	NotificationMoneyTransferAccountTypeCard = NotificationMoneyTransferAccountType("card")
)

func notificationMoneyTransferAccountTypeFromID(accountID string) NotificationMoneyTransferAccountType {
	if strings.HasPrefix(accountID, AccountPrefix) || strings.HasPrefix(accountID, RegisteredAccountPrefix) {
		return NotificationMoneyTransferAccountTypeBank
	} else if strings.HasPrefix(accountID, CardPrefix) || strings.HasPrefix(accountID, RegisteredCardPrefix) {
		return NotificationMoneyTransferAccountTypeBank
	}

	return NotificationMoneyTransferAccountType("")
}

func partnerMoneyTransferAccountTypeFromID(accountID string) bank.MoneyTransferAccountType {
	if strings.HasPrefix(accountID, AccountPrefix) || strings.HasPrefix(accountID, RegisteredAccountPrefix) {
		return bank.MoneyTransferAccountTypeBank
	} else if strings.HasPrefix(accountID, CardPrefix) || strings.HasPrefix(accountID, RegisteredCardPrefix) {
		return bank.MoneyTransferAccountTypeCard
	}

	return bank.MoneyTransferAccountType("")
}
