package activity

import (
	"encoding/json"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/wiseco/core-platform/notification/push"
	"github.com/wiseco/core-platform/services/data"
)

type ConsumerCreator interface {
	Update(Consumer) error
	text(Activity, Language) (string, error)
}

type consumerCreator struct {
	composer TextComposer
	pusher   *push.Pusher
	*sqlx.DB
}

//NewContactCreator returns a new contact activity service
func NewConsumerCreator() ConsumerCreator {
	return &consumerCreator{NewTextComposer(), nil, data.DBWrite}
}

func (c *consumerCreator) Update(consumer Consumer) error {
	_, err := c.Exec(`INSERT INTO user_activity(entity_id, activity_type, activity_action, resource_id, metadata) 
			VALUES($1, $2, $3, $4, $5)`, consumer.EntityID, TypeConsumer, ActionUpdate, consumer.EntityID, consumer.raw())

	if err != nil {
		log.Printf("Error creating activity, error:%v", err)
		return err
	}

	return err
}

func (c *consumerCreator) text(a Activity, lang Language) (string, error) {
	var consumer Consumer
	err := json.Unmarshal(a.Metadata, &consumer)
	if err != nil {
		log.Printf("Error parsing metadata to construct text for activity:%v error:%v", a, err)
		return "", err
	}

	switch *a.Action {
	case ActionUpdate:
		if consumer.Phone != nil {
			return c.composer.Compose(consumerPhoneUpdateTempl, consumer).String(lang)
		} else if consumer.Email != nil {
			return c.composer.Compose(consumerEmailUpdateTempl, consumer).String(lang)
		} else if consumer.Address != nil {
			return c.composer.Compose(consumerAddressUpdateTempl, consumer.Address).String(lang)
		} else {
			return "", nil
		}
	default:
		return "", nil
	}

}
