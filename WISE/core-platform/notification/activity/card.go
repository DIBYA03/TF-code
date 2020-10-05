package activity

import (
	"encoding/json"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/wiseco/core-platform/notification/push"
	"github.com/wiseco/core-platform/services/data"
)

//CardCreator - handles creation, activation, block/unblock of debit cards
type CardCreator interface {
	StatusUpdate(CardStatus) error
	text(Activity, Language) (string, error)
}

type cardCreator struct {
	composer TextComposer
	pusher   *push.Pusher
	*sqlx.DB
}

//NewCardTransationCreator return a card transaction activity service
func NewCardCreator() CardCreator {
	return cardCreator{NewTextComposer(), nil, data.DBWrite}
}

func (c cardCreator) StatusUpdate(s CardStatus) error {
	_, err := c.Exec(`INSERT INTO user_activity(entity_id, activity_type, activity_action, resource_id, metadata) 
			VALUES($1, $2, $3, $4, $5)`, s.EntityID, TypeCard, ActionUpdate, s.CardID, s.raw())

	if err != nil {
		log.Printf("Error creating activity, error:%v", err)
		return err
	}

	return err
}

func (c cardCreator) text(a Activity, lang Language) (string, error) {
	var t CardStatus
	err := json.Unmarshal(a.Metadata, &t)
	if err != nil {
		return "", err
	}

	switch *a.Action {
	case ActionCreate:
		return c.composer.Compose(cardCreateTempl, t).String(lang)
	case ActionUpdate:
		if t.BusinessName != nil {
			return c.composer.Compose(cardStatusUpdateBusinessNameTempl, t).String(lang)
		} else {
			return c.composer.Compose(cardStatusUpdateTempl, t).String(lang)
		}
	default:
		return "", nil
	}

}
