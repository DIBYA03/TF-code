package activity

import (
	"encoding/json"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/wiseco/core-platform/notification/push"
	"github.com/wiseco/core-platform/services/data"
)

//CardTransactionCreator - handles all card transactions like atm withdrawal, purchase, card decline etc..
type CardTransactionCreator interface {
	PostedDebit(CardTransaction, Type) (*string, error)
	PostedCredit(CardTransaction, Type) (*string, error)
	Declined(CardTransaction) (*string, error)
	Authorized(CardTransaction) (*string, error)
	AuthReversed(CardTransaction) (*string, error)
	HoldApproved(CardTransaction) (*string, error)
	HoldExpired(CardTransaction) (*string, error)
	text(Activity, Language) (string, error)
}

type cardTransactionCreator struct {
	composer TextComposer
	pusher   *push.Pusher
	*sqlx.DB
}

//NewCardTransationCreator return a card transaction activity service
func NewCardTransationCreator() CardTransactionCreator {
	return cardTransactionCreator{NewTextComposer(), nil, data.DBWrite}
}

func (c cardTransactionCreator) Declined(t CardTransaction) (*string, error) {

	sqlStatement := `INSERT INTO business_activity(entity_id, activity_type, activity_action, metadata) 
			VALUES($1, $2, $3, $4) RETURNING id`

	var id string
	err := c.QueryRow(sqlStatement, t.EntityID, TypeCardTransaction, ActionDecline, t.raw()).Scan(&id)

	if err != nil {
		log.Printf("Error creating activity, error:%v", err)
		return nil, err
	}

	return &id, nil
}

func (c cardTransactionCreator) AuthReversed(t CardTransaction) (*string, error) {

	sqlStatement := `INSERT INTO business_activity(entity_id, activity_type, activity_action, metadata) 
			VALUES($1, $2, $3, $4) RETURNING id`

	var id string
	err := c.QueryRow(sqlStatement, t.EntityID, TypeCardTransaction, ActionAuthReverse, t.raw()).Scan(&id)

	if err != nil {
		log.Printf("Error creating activity, error:%v", err)
		return nil, err
	}

	return nil, err
}

func (c cardTransactionCreator) HoldExpired(t CardTransaction) (*string, error) {

	sqlStatement := `INSERT INTO business_activity(entity_id, activity_type, activity_action, metadata) 
			VALUES($1, $2, $3, $4) RETURNING id`

	var id string
	err := c.QueryRow(sqlStatement, t.EntityID, TypeCardTransaction, ActionHoldReleased, t.raw()).Scan(&id)

	if err != nil {
		log.Printf("Error creating activity, error:%v", err)
		return nil, err
	}

	return nil, err
}

func (c cardTransactionCreator) HoldApproved(t CardTransaction) (*string, error) {

	sqlStatement := `INSERT INTO business_activity(entity_id, activity_type, activity_action, metadata) 
			VALUES($1, $2, $3, $4) RETURNING id`

	var id string
	err := c.QueryRow(sqlStatement, t.EntityID, TypeCardTransaction, ActionHold, t.raw()).Scan(&id)

	if err != nil {
		log.Printf("Error creating activity, error:%v", err)
		return nil, err
	}

	return nil, err
}

func (c cardTransactionCreator) Authorized(t CardTransaction) (*string, error) {

	sqlStatement := `INSERT INTO business_activity(entity_id, activity_type, activity_action, metadata) 
			VALUES($1, $2, $3, $4) RETURNING id`

	var id string
	err := c.QueryRow(sqlStatement, t.EntityID, TypeCardTransaction, ActionAuthorize, t.raw()).Scan(&id)

	if err != nil {
		log.Printf("Error creating activity, error:%v", err)
		return nil, err
	}

	return nil, err
}

func (c cardTransactionCreator) PostedDebit(t CardTransaction, activityType Type) (*string, error) {
	sqlStatement := `INSERT INTO business_activity(entity_id, activity_type, activity_action, metadata) 
			VALUES($1, $2, $3, $4) RETURNING id`

	var id string
	err := c.QueryRow(sqlStatement, t.EntityID, activityType, ActionPostedDebit, t.raw()).Scan(&id)

	if err != nil {
		log.Printf("Error creating activity, error:%v", err)
		return nil, err
	}

	return &id, err
}

func (c cardTransactionCreator) PostedCredit(t CardTransaction, activityType Type) (*string, error) {

	sqlStatement := `INSERT INTO business_activity(entity_id, activity_type, activity_action, metadata) 
			VALUES($1, $2, $3, $4) RETURNING id`

	var id string
	err := c.QueryRow(sqlStatement, t.EntityID, activityType, ActionPostedCredit, t.raw()).Scan(&id)

	if err != nil {
		log.Printf("Error creating activity, error:%v", err)
		return nil, err
	}

	return &id, err
}

func (c cardTransactionCreator) text(a Activity, lang Language) (string, error) {
	var t CardTransaction
	err := json.Unmarshal(a.Metadata, &t)
	if err != nil {
		return "", err
	}

	switch *a.Action {
	case ActionDecline:
		return c.handleCardDeclinedActivity(t, a, lang)
	case ActionAuthorize:
		return c.handleCardAuthorizationActivity(t, a, lang)
	case ActionAuthReverse:
		return c.handleCardAuthReversalActivity(t, a, lang)
	case ActionPostedDebit:
		return c.handleCardDebitActivity(t, a, lang)
	case ActionPostedCredit:
		return c.handleCardCreditActivity(t, a, lang)
	default:
		return "", nil
	}

}

func (c cardTransactionCreator) handleCardDebitActivity(t CardTransaction, a Activity, lang Language) (string, error) {
	switch a.ActivityType {
	case TypeCardReaderPurchaseDebitOnline:
		fallthrough
	case TypeCardReaderPurchaseDebit:
		if t.Merchant != "" {
			return c.composer.Compose(cardPostedDebitCardReaderTempl, t).String(lang)
		} else {
			return c.composer.Compose(cardPostedDebitCardReaderGenericTempl, t).String(lang)
		}
	case TypeCardATMDebit:
		return c.composer.Compose(cardPostedDebitCardATMTempl, t).String(lang)
	default:
		if t.Merchant != "" {
			return c.composer.Compose(cardPostedDebitTempl, t).String(lang)
		} else {
			return c.composer.Compose(cardPostedDebitGenericTempl, t).String(lang)
		}
	}
}

func (c cardTransactionCreator) handleCardCreditActivity(t CardTransaction, a Activity, lang Language) (string, error) {
	switch a.ActivityType {
	case TypeMerchantRefundCredit:
		return c.composer.Compose(cardPostedCreditMerchantRefundTempl, t).String(lang)
	case TypeCardPushCredit:
		if len(t.Merchant) > 0 {
			return c.composer.Compose(cardPushDebitCreditTempl, t).String(lang)
		} else {
			return c.composer.Compose(cardPushDebitCreditGenericTempl, t).String(lang)
		}
	case TypeCardVisaCredit:
		if len(t.Merchant) > 0 {
			return c.composer.Compose(cardVisaCreditTempl, t).String(lang)
		} else {
			return c.composer.Compose(cardVisaCreditGenericTempl, t).String(lang)
		}
	default:
		return c.composer.Compose(cardPostedCreditTempl, t).String(lang)
	}
}

func (c cardTransactionCreator) handleCardAuthorizationActivity(t CardTransaction, a Activity, lang Language) (string, error) {
	if t.BusinessName != "" {
		if t.Merchant != "" {
			return c.composer.Compose(cardAuthorizeBusinessNameTempl, t).String(lang)
		} else {
			return c.composer.Compose(cardAuthorizeBusinessNameGenericTempl, t).String(lang)
		}
	} else {
		if t.Merchant != "" {
			return c.composer.Compose(cardAuthorizeTempl, t).String(lang)
		} else {
			return c.composer.Compose(cardAuthorizeGenericTempl, t).String(lang)
		}
	}
}

func (c cardTransactionCreator) handleCardAuthReversalActivity(t CardTransaction, a Activity, lang Language) (string, error) {
	if t.Merchant != "" {
		return c.composer.Compose(cardAuthReversalTempl, t).String(lang)
	} else {
		return c.composer.Compose(cardAuthReversalGenericTempl, t).String(lang)
	}
}

func (c cardTransactionCreator) handleHoldApproveActivity(t CardTransaction, a Activity, lang Language) (string, error) {
	if t.Merchant != "" {
		return c.composer.Compose(cardHoldApproveTempl, t).String(lang)
	} else {
		return c.composer.Compose(cardHoldApproveGenericTempl, t).String(lang)
	}
}

func (c cardTransactionCreator) handleCardDeclinedActivity(t CardTransaction, a Activity, lang Language) (string, error) {
	if t.BusinessName != "" {
		if t.Merchant != "" {
			return c.composer.Compose(cardDeclineBusinessNameTempl, t).String(lang)
		} else {
			return c.composer.Compose(cardDeclineBusinessNameGenericTempl, t).String(lang)
		}
	} else {
		if t.Merchant != "" {
			return c.composer.Compose(cardDeclineTempl, t).String(lang)
		} else {
			return c.composer.Compose(cardDeclineGenericTempl, t).String(lang)
		}
	}
}
