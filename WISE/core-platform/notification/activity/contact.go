package activity

import (
	"encoding/json"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/wiseco/core-platform/notification/push"
	"github.com/wiseco/core-platform/services/data"
)

type ContactCreator interface {
	Create(Contact) error
	Update(Contact) error
	Delete(Contact) error
	text(Activity, Language) (string, error)
}

type contactCreator struct {
	composer TextComposer
	pusher   *push.Pusher
	*sqlx.DB
}

//NewContactCreator returns a new contact activity service
func NewContactCreator() ContactCreator {
	return &contactCreator{NewTextComposer(), nil, data.DBWrite}
}

func (c *contactCreator) Create(contact Contact) error {

	_, err := c.Exec(`INSERT INTO business_activity(entity_id, activity_type, activity_action, resource_id, metadata) 
			VALUES($1, $2, $3, $4, $5)`, contact.EntityID, TypeContact, ActionCreate, contact.ContactID, contact.raw())

	if err != nil {
		log.Printf("Error creating activity, error:%v", err)
		return err
	}

	return err
}

func (c *contactCreator) Update(contact Contact) error {
	_, err := c.Exec(`INSERT INTO business_activity(entity_id, activity_type, activity_action, resource_id, metadata) 
			VALUES($1, $2, $3, $4, $5)`, contact.EntityID, TypeContact, ActionUpdate, contact.ContactID, contact.raw())

	if err != nil {
		log.Printf("Error creating activity, error:%v", err)
		return err
	}

	return err
}

func (c *contactCreator) Delete(contact Contact) error {
	_, err := c.Exec(`INSERT INTO business_activity(entity_id, activity_type, activity_action, resource_id, metadata) 
			VALUES($1, $2, $3, $4, $5)`, contact.EntityID, TypeContact, ActionDelete, contact.ContactID, contact.raw())

	if err != nil {
		log.Printf("Error creating activity, error:%v", err)
		return err
	}

	return err
}

func (c *contactCreator) text(a Activity, lang Language) (string, error) {
	var contact Contact
	err := json.Unmarshal(a.Metadata, &contact)
	if err != nil {
		log.Printf("Error parsing metadata to construct text for activity:%v error:%v", a, err)
		return "", err
	}

	switch *a.Action {
	case ActionCreate:
		return c.composer.Compose(contactCreateTempl, contact).String(lang)
	case ActionUpdate:
		return c.composer.Compose(contactUpdateTempl, contact).String(lang)
	case ActionDelete:
		return c.composer.Compose(contactDeleteTempl, contact).String(lang)
	default:
		return "", nil
	}

}
