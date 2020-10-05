package activity

import (
	"encoding/json"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/wiseco/core-platform/notification/push"
	"github.com/wiseco/core-platform/services/data"
)

type BusinessCreator interface {
	Update(Business) error
	text(Activity, Language) (string, error)
}

type businessCreator struct {
	composer TextComposer
	pusher   *push.Pusher
	*sqlx.DB
}

//NewContactCreator returns a new contact activity service
func NewBusinessCreator() BusinessCreator {
	return &businessCreator{NewTextComposer(), nil, data.DBWrite}
}

func (c *businessCreator) Update(business Business) error {
	_, err := c.Exec(`INSERT INTO business_activity(entity_id, activity_type, activity_action, resource_id, metadata) 
			VALUES($1, $2, $3, $4, $5)`, business.EntityID, TypeBusiness, ActionUpdate, business.EntityID, business.raw())

	if err != nil {
		log.Printf("Error creating activity, error:%v", err)
		return err
	}

	return err
}

func (c *businessCreator) text(a Activity, lang Language) (string, error) {
	var business Business
	err := json.Unmarshal(a.Metadata, &business)
	if err != nil {
		log.Printf("Error parsing metadata to construct text for activity:%v error:%v", a, err)
		return "", err
	}

	switch *a.Action {
	case ActionUpdate:
		if business.Phone != nil {
			return c.composer.Compose(businessPhoneUpdateTempl, business).String(lang)
		} else if business.Email != nil {
			return c.composer.Compose(businessEmailUpdateTempl, business).String(lang)
		} else if business.Address != nil {
			return c.composer.Compose(businessAddressUpdateTempl, business.Address).String(lang)
		} else {
			return "", nil
		}
	default:
		return "", nil
	}

}
