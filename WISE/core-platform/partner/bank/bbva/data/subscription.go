package data

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/wiseco/core-platform/partner/bank"
	"github.com/wiseco/core-platform/partner/bank/bbva"
)

type SubscriptionService interface {
	// Subscribe to notifications
	SubscribeAll(SubscriptionChannelType, ChannelURL) ([]Subscription, error)
	Subscribe(SubscriptionChannelType, ChannelURL, map[bbva.EventType]EventTypeConfig) ([]Subscription, error)
	GetAll(ChannelURL) ([]Subscription, error)
	GetAllByEvents(ChannelURL, []bbva.EventType) ([]Subscription, error)
	UnsubscribeAllByChannel(ChannelURL) error
	Unsubscribe(ChannelURL, []bbva.EventType) error
	UnsubscribeAll() error
}

type subscriptionService struct {
	request bank.APIRequest
	client  *client
	wdb     *sqlx.DB
	rdb     *sqlx.DB
}

func GetSubscriptionService(r bank.APIRequest) SubscriptionService {
	return &subscriptionService{
		request: r,
		client:  newClient(),
		wdb:     DBWrite,
		rdb:     DBRead,
	}
}

type SubscriptionChannelType string

const (
	SubscriptionChannelTypeHTTP = SubscriptionChannelType("http")
	SubscriptionChannelTypeSQS  = SubscriptionChannelType("sqs")
)

const SubscriptionFormatPRAPI = "PRAPI"

type SubscriptionRequest struct {
	EventType   bbva.EventType          `json:"event_type"`
	Version     bbva.Version            `json:"event_version"`
	ChannelType SubscriptionChannelType `json:"channel"`
	CallbackURL ChannelURL              `json:"callback"`
	Format      string                  `json:"format,omitempty"`
}

type SubscriptionResponse struct {
	SubscriptionID   string                  `json:"id"`
	EventID          string                  `json:"event_id"`
	EventDescription string                  `json:"event_description"`
	CreationDate     time.Time               `json:"creation_date"`
	EventType        bbva.EventType          `json:"event_type"`
	ChannelType      SubscriptionChannelType `json:"channel"`
	CallbackURL      ChannelURL              `json:"callback"`
}

type GetSubscriptionResponse struct {
	SubscriptionID   string         `json:"id"`
	EventID          string         `json:"event_id"`
	EventType        bbva.EventType `json:"event_type"`
	EventDescription string         `json:"event_description"`
	CreationDate     time.Time      `json:"creation_date"`
}

func (s *subscriptionService) SubscribeAll(channelType SubscriptionChannelType, channelURL ChannelURL) ([]Subscription, error) {
	path := "notifications/v3.0/subscriptions"
	var subs []SubscriptionResponse
	for _, ec := range AllEventConfigTypes {
		for e, config := range ec {

			request := &SubscriptionRequest{
				EventType:   e,
				Version:     bbva.NotificationVersion,
				ChannelType: channelType,
				CallbackURL: channelURL,
				Format:      config.Format,
			}
			req, err := s.client.post(path, s.request, request)
			if err != nil {
				return nil, err
			}

			var resp = SubscriptionResponse{}
			if err := s.client.do(req, &resp); err != nil {
				log.Println(fmt.Sprintf("Error subscribing to event type (%s): ", e), err)
			} else {
				resp.EventType = e
				resp.ChannelType = channelType
				resp.CallbackURL = channelURL
				subs = append(subs, resp)
			}
		}
	}

	return s.createMany(subs)
}

func (s *subscriptionService) Subscribe(channelType SubscriptionChannelType, channelURL ChannelURL, events map[bbva.EventType]EventTypeConfig) ([]Subscription, error) {
	path := "notifications/v3.0/subscriptions"
	var subs []SubscriptionResponse

	for e, config := range events {
		request := &SubscriptionRequest{
			EventType:   e,
			Version:     bbva.NotificationVersion,
			ChannelType: channelType,
			CallbackURL: channelURL,
			Format:      config.Format,
		}
		req, err := s.client.post(path, s.request, request)
		if err != nil {
			return nil, err
		}

		var resp = SubscriptionResponse{}
		if err := s.client.do(req, &resp); err != nil {
			log.Println(fmt.Sprintf("Error subscribing to event type (%s): ", e), err)
		} else {
			resp.EventType = e
			resp.ChannelType = channelType
			resp.CallbackURL = channelURL
			subs = append(subs, resp)
		}
	}

	return s.createMany(subs)
}

func (s *subscriptionService) GetAll(channelURL ChannelURL) ([]Subscription, error) {
	return s.getAllByChannel(channelURL)
}

func (s *subscriptionService) GetAllByEvents(channelURL ChannelURL, events []bbva.EventType) ([]Subscription, error) {
	return s.getByEvents(channelURL, events)
}

func (s *subscriptionService) UnsubscribeAllByChannel(channelURL ChannelURL) error {
	path := "notifications/v3.0/subscriptions"

	subs, err := s.getAllByChannel(channelURL)
	if err != nil {
		return err
	}

	for _, sub := range subs {
		url := fmt.Sprintf(path, sub.SubscriptionID)
		req, err := s.client.delete(url, s.request)
		if err != nil {
			return err
		}

		err = s.client.do(req, nil)
		if err != nil {
			log.Println(fmt.Sprintf("Error unsubscribing to event type (%s): ", sub.EventType), err)
		}
	}

	return s.deleteAllByChannel(channelURL)
}

func (s *subscriptionService) Unsubscribe(channelURL ChannelURL, events []bbva.EventType) error {
	path := "notifications/v3.0/subscriptions/%s"

	subs, err := s.getByEvents(channelURL, events)
	if err != nil {
		return err
	}

	for _, sub := range subs {
		url := fmt.Sprintf(path, sub.SubscriptionID)
		req, err := s.client.delete(url, s.request)
		if err != nil {
			return err
		}

		err = s.client.do(req, nil)
		if err != nil {
			log.Println(fmt.Sprintf("Error subscribing to event type (%s): ", sub.EventType), err)
		}
	}

	return s.deleteAllByEvents(channelURL, events)
}

func (s *subscriptionService) UnsubscribeAll() error {
	path := "notifications/v3.0/subscriptions"
	req, err := s.client.delete(path, s.request)
	if err != nil {
		return err
	}

	err = s.client.do(req, nil)
	if err != nil {
		return err
	}

	return s.deleteAll()
}

type SubscriptionID string
type ChannelURL string

type Subscription struct {
	ID             SubscriptionID          `json:"id" db:"id"`
	SubscriptionID string                  `json:"bbvaSubscriptionID" db:"subscription_id"`
	EventID        string                  `json:"eventId" db:"event_id"`
	EventType      bbva.EventType          `json:"eventType" db:"event_type"`
	EventDesc      string                  `json:"eventDesc" db:"event_desc"`
	ChannelType    SubscriptionChannelType `json:"channelType" db:"channel_type"`
	ChannelURL     ChannelURL              `json:"channelUrl" db:"channel_url"`
	Created        time.Time               `json:"created" db:"created"`
}

type SubscriptionCreate struct {
	SubscriptionID string                  `json:"id" db:"subscription_id"`
	EventID        string                  `json:"eventId" db:"event_id"`
	EventType      bbva.EventType          `json:"eventType" db:"event_type"`
	EventDesc      string                  `json:"eventDesc" db:"event_desc"`
	ChannelType    SubscriptionChannelType `json:"channelType" db:"channel_type"`
	ChannelURL     ChannelURL              `json:"channelUrl" db:"channel_url"`
	Created        time.Time               `json:"created" db:"created"`
}

func (s *subscriptionService) getByID(id SubscriptionID) (*Subscription, error) {
	var sub Subscription

	err := s.wdb.Get(&sub, "SELECT * FROM notification_subscription WHERE id = $1", id)
	if err != nil {
		return nil, err
	}

	return &sub, nil
}

func (s *subscriptionService) getByEvents(channelURL ChannelURL, events []bbva.EventType) ([]Subscription, error) {
	var subs []Subscription
	q, args, err := sqlx.In("SELECT * FROM notification_subscription WHERE channel_url = ? AND event_type IN (?)", channelURL, events)
	if err != nil {
		return subs, err
	}

	q = s.wdb.Rebind(q)
	err = s.wdb.Select(&subs, q, args...)
	if err == sql.ErrNoRows {
		return []Subscription{}, nil
	} else if err != nil {
		return nil, err
	}

	return subs, nil
}

func (s *subscriptionService) getAllByChannel(channelURL ChannelURL) ([]Subscription, error) {
	var subs []Subscription

	err := s.wdb.Select(&subs, "SELECT * FROM notification_subscription channel_url = $1", channelURL)
	if err == sql.ErrNoRows {
		return []Subscription{}, nil
	} else if err != nil {
		return nil, err
	}

	return subs, nil
}

func (s *subscriptionService) create(resp SubscriptionResponse) (*Subscription, error) {
	c := SubscriptionCreate{
		SubscriptionID: resp.SubscriptionID,
		EventID:        resp.EventID,
		EventType:      resp.EventType,
		EventDesc:      resp.EventDescription,
		ChannelType:    resp.ChannelType,
		ChannelURL:     resp.CallbackURL,
		Created:        resp.CreationDate,
	}

	q := `
        INSERT INTO notification_subscription (subscription_id, event_id, event_type, event_desc, channel_type, channel_url, created)
		VALUES (:subscription_id, :event_id, :event_type, :event_desc, :channel_type, :channel_url, :created)
        RETURNING id`

	stmt, err := s.wdb.PrepareNamed(q)
	if err != nil {
		return nil, err
	}

	var id SubscriptionID
	err = stmt.Get(&id, &c)
	if err != nil {
		return nil, err
	}

	return s.getByID(id)
}

func (s *subscriptionService) createMany(resp []SubscriptionResponse) ([]Subscription, error) {
	var subs []Subscription
	if len(resp) == 0 {
		log.Println("No subscriptions to create")
		return subs, nil
	}

	for _, respSub := range resp {
		sub, err := s.create(respSub)
		if err != nil {
			log.Println(err)
			continue
		}

		subs = append(subs, *sub)
	}

	return subs, nil
}

func (s *subscriptionService) deleteByID(id SubscriptionID) error {
	_, err := s.wdb.Exec("DELETE FROM notification_subscription WHERE id = $1", id)
	return err
}

func (s *subscriptionService) deleteAllByChannel(channelURL ChannelURL) error {
	_, err := s.wdb.Exec("DELETE FROM notification_subscription WHERE channel_url = $1", channelURL)
	return err
}

func (s *subscriptionService) deleteAllByEvents(channelURL ChannelURL, events []bbva.EventType) error {
	q, args, err := sqlx.In("DELETE FROM notification_subscription WHERE channel_url = ? AND event_type IN (?)", channelURL, events)
	if err != nil {
		return err
	}

	q = s.wdb.Rebind(q)
	_, err = s.wdb.Exec(q, args...)
	return err
}

func (s *subscriptionService) deleteAll() error {
	_, err := s.wdb.Exec("DELETE FROM notification_subscription")
	return err
}
