package data

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/firehose"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/types"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/wiseco/core-platform/partner/bank"
	"github.com/wiseco/core-platform/shared"
)

type RawNotification struct {
	ID       string          `json:"id"`
	BankName string          `json:"bankName"`
	SourceID string          `json:"sourceId"`
	Message  json.RawMessage `json:"message"`
}

type Notification = bank.Notification

type NotificationCreate struct {
	EntityID   bank.NotificationEntityID   `json:"entityId" db:"entity_id"`
	EntityType bank.NotificationEntityType `json:"entityType" db:"entity_type"`
	BankName   bank.ProviderName           `json:"bankName" db:"bank_name"`
	SourceID   bank.SourceID               `json:"sourceId" db:"source_id"`
	Type       bank.NotificationType       `json:"type" db:"notification_type"`
	Action     bank.NotificationAction     `json:"action" db:"notification_action"`
	Attribute  *bank.NotificationAttribute `json:"attribute" db:"notification_attribute"`
	Version    string                      `json:"version" db:"notification_version"`
	Created    time.Time                   `json:"created" db:"created"`
	Data       types.JSONText              `json:"data" db:"notification_data"`
}

type NotificationService interface {
	LogRaw(shared.StreamProvider, []RawNotification) error
	GetByID(bank.NotificationID) (*bank.Notification, error)
	GetBySourceID(bank.SourceID, bank.ProviderName) (*bank.Notification, error)
	Create(*NotificationCreate) (*bank.Notification, error)
	IncrementSend([]bank.NotificationID) error
}

type notificationService struct {
	request bank.APIRequest
	wdb     *sqlx.DB
	rdb     *sqlx.DB
}

func NewNotificationService(r bank.APIRequest) NotificationService {
	return &notificationService{
		request: r,
		wdb:     DBWrite,
		rdb:     DBRead,
	}
}

func (s *notificationService) LogRaw(provider shared.StreamProvider, notifications []RawNotification) error {

	if provider.Type() != shared.StreamProviderTypeKinesis {
		return &bank.Error{
			RawError: errors.New(fmt.Sprintf("unsupported stream type (%s)", provider.Type())),
			Code:     bank.ErrorCodeLogRawNotification,
		}
	}

	sess, err := session.NewSession(&aws.Config{Region: aws.String(provider.Region().String())})
	if err != nil {
		return &bank.Error{
			RawError: err,
			Code:     bank.ErrorCodeLogRawNotification,
		}
	}

	f := firehose.New(sess)

	var records []*firehose.Record
	for _, n := range notifications {
		b, err := json.Marshal(n)
		if err != nil {
			return &bank.Error{
				RawError: err,
				Code:     bank.ErrorCodeLogRawNotification,
			}
		}

		// Append new line to simplify parsing of S3 data
		records = append(records, &firehose.Record{Data: append(b, '\n')})
	}

	name := provider.StreamName()
	in := firehose.PutRecordBatchInput{
		DeliveryStreamName: &name,
		Records:            records,
	}

	_, err = f.PutRecordBatch(&in)
	if err != nil {
		return &bank.Error{
			RawError: err,
			Code:     bank.ErrorCodeLogRawNotification,
		}
	}

	/* TODO: Handle retries?
	log.Printf("Kinesis put record count: %d", len(notifications))
	log.Printf("Kinesis put record failure count: %d", *out.FailedPutCount) */
	return nil
}

func (s *notificationService) GetByID(id bank.NotificationID) (*bank.Notification, error) {
	var n bank.Notification
	err := s.wdb.Get(&n, "SELECT * FROM notification WHERE id = $1", id)
	return &n, err
}

func (s *notificationService) GetBySourceID(sourceID bank.SourceID, name bank.ProviderName) (*bank.Notification, error) {
	var n bank.Notification
	err := s.wdb.Get(&n, "SELECT * FROM notification WHERE source_id = $1 AND bank_name = $2", sourceID, name)
	return &n, err
}

func (s *notificationService) Create(n *NotificationCreate) (*bank.Notification, error) {
	sql := `
        INSERT INTO notification(
			entity_id, entity_type, bank_name, source_id, notification_type, notification_action,
			notification_attribute, notification_version, created, notification_data
		)
        VALUES(
			:entity_id, :entity_type, :bank_name, :source_id, :notification_type, :notification_action,
			:notification_attribute, :notification_version, :created, :notification_data
		)
        RETURNING id`

	stmt, err := s.wdb.PrepareNamed(sql)
	if err != nil {
		return nil, err
	}

	var id bank.NotificationID
	err = stmt.Get(&id, &n)
	if err != nil {
		if pqerr, ok := err.(*pq.Error); ok {
			if pqerr.Code == pq.ErrorCode("23505") {
				return nil, &bank.Error{
					RawError: err,
					Code:     bank.ErrorCodeDuplicateNotification,
				}
			}

			return nil, &bank.Error{
				RawError: err,
				Code:     bank.ErrorCodeInternalDatabaseError,
			}
		}

		return nil, err
	}

	return s.GetByID(id)
}

func (s *notificationService) IncrementSend(ids []bank.NotificationID) error {
	q, args, err := sqlx.In("UPDATE notification SET send_counter=send_counter+1 WHERE id IN (?)", ids)
	if err != nil {
		return err
	}

	q = s.wdb.Rebind(q)
	_, err = s.wdb.Exec(q, args...)
	return err
}
