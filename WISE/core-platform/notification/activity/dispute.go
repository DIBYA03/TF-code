package activity

import (
	"encoding/json"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/wiseco/core-platform/notification/push"
	"github.com/wiseco/core-platform/services/data"
)

type DisputeCreator interface {
	Create(Dispute) error
	Delete(Dispute) error
	text(Activity, Language) (string, error)
}

type disputeCreator struct {
	composer TextComposer
	pusher   *push.Pusher
	*sqlx.DB
}

//NewDisputeCreator returns a new dispute activity service
func NewDisputeCreator() DisputeCreator {
	return &disputeCreator{NewTextComposer(), nil, data.DBWrite}
}

func (c *disputeCreator) Create(d Dispute) error {

	_, err := c.Exec(`INSERT INTO business_activity(entity_id, activity_type, activity_action, resource_id, metadata) 
			VALUES($1, $2, $3, $4, $5)`, d.EntityID, TypeDispute, ActionCreate, d.TransactionID, d.raw())

	if err != nil {
		log.Printf("Error creating activity, error:%v", err)
		return err
	}

	return err
}

func (c *disputeCreator) Delete(d Dispute) error {
	_, err := c.Exec(`INSERT INTO business_activity(entity_id, activity_type, activity_action, resource_id, metadata) 
			VALUES($1, $2, $3, $4, $5)`, d.EntityID, TypeDispute, ActionDelete, d.TransactionID, d.raw())

	if err != nil {
		log.Printf("Error creating activity, error:%v", err)
		return err
	}

	return err
}

func (c *disputeCreator) text(a Activity, lang Language) (string, error) {
	var dispute Dispute
	err := json.Unmarshal(a.Metadata, &dispute)
	if err != nil {
		log.Printf("Error parsing metadata to construct text for activity:%v error:%v", a, err)
		return "", err
	}

	switch *a.Action {
	case ActionCreate:
		return c.composer.Compose(disputeCreateTempl, dispute).String(lang)
	case ActionDelete:
		return c.composer.Compose(disputeDeleteTempl, dispute).String(lang)
	default:
		return "", nil
	}

}
