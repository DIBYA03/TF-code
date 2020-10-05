package activity

import (
	"encoding/json"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/wiseco/core-platform/notification/push"
	"github.com/wiseco/core-platform/services/data"
)

//AccountTransactionCreator - Handle all money transfer transactions
type TransferTransactionCreator interface {
	AccountOriginated(AccountTransaction) (*string, error)
	DebitPosted(AccountTransaction, Type) (*string, error)
	DebitInProcess(AccountTransaction) (*string, error)
	CreditPosted(AccountTransaction, Type) (*string, error)
	CreditInProcess(AccountTransaction) (*string, error)
	HoldApproved(AccountTransaction) (*string, error)
	HoldReleased(AccountTransaction) (*string, error)
	text(Activity, Language) (string, error)
}

type transferCreator struct {
	composer TextComposer
	pusher   *push.Pusher
	*sqlx.DB
}

//NewTransferCreator return a transfer transaction activity creator
func NewTransferCreator() TransferTransactionCreator {
	return transferCreator{NewTextComposer(), nil, data.DBWrite}
}

func (ac transferCreator) AccountOriginated(t AccountTransaction) (*string, error) {

	sqlStatement := `INSERT INTO business_activity(entity_id, activity_type, activity_action, metadata, activity_date) 
			VALUES($1, $2, $3, $4, $5) RETURNING id`

	var id string
	err := ac.QueryRow(sqlStatement, t.EntityID, TypeTransferTransaction, ActionAccountOriginated, t.raw(), t.TransactionDate).Scan(&id)

	if err != nil {
		log.Printf("Error creating activity, error:%v", err)
		return nil, err
	}

	return &id, err
}

func (ac transferCreator) DebitPosted(t AccountTransaction, activityType Type) (*string, error) {

	sqlStatement := `INSERT INTO business_activity(entity_id, activity_type, activity_action, metadata, activity_date) 
			VALUES($1, $2, $3, $4, $5) RETURNING id`

	var id string
	err := ac.QueryRow(sqlStatement, t.EntityID, activityType, ActionPostedDebit, t.raw(), t.TransactionDate).Scan(&id)

	if err != nil {
		log.Printf("Error creating activity, error:%v", err)
		return nil, err
	}

	return &id, err
}

func (ac transferCreator) DebitInProcess(t AccountTransaction) (*string, error) {
	sqlStatement := `INSERT INTO business_activity(entity_id, activity_type, activity_action, metadata, activity_date) 
			VALUES($1, $2, $3, $4, $5) RETURNING id`

	var id string
	err := ac.QueryRow(sqlStatement, t.EntityID, TypeTransferTransaction, ActionInProcessDebit, t.raw(), t.TransactionDate).Scan(&id)

	if err != nil {
		log.Printf("Error creating activity, error:%v", err)
		return nil, err
	}

	return &id, err
}

func (ac transferCreator) CreditPosted(t AccountTransaction, activityType Type) (*string, error) {

	sqlStatement := `INSERT INTO business_activity(entity_id, activity_type, activity_action, metadata, activity_date) 
			VALUES($1, $2, $3, $4, $5) RETURNING id`

	var id string
	err := ac.QueryRow(sqlStatement, t.EntityID, activityType, ActionPostedCredit, t.raw(), t.TransactionDate).Scan(&id)

	if err != nil {
		log.Printf("Error creating activity, error:%v", err)
		return nil, err
	}

	return &id, err
}

func (ac transferCreator) CreditInProcess(t AccountTransaction) (*string, error) {

	sqlStatement := `INSERT INTO business_activity(entity_id, activity_type, activity_action, metadata, activity_date) 
			VALUES($1, $2, $3, $4, $5) RETURNING id`

	var id string
	err := ac.QueryRow(sqlStatement, t.EntityID, TypeTransferTransaction, ActionInProcessCredit, t.raw(), t.TransactionDate).Scan(&id)

	if err != nil {
		log.Printf("Error creating activity, error:%v", err)
		return nil, err
	}

	return &id, err
}

func (ac transferCreator) HoldApproved(t AccountTransaction) (*string, error) {

	sqlStatement := `INSERT INTO business_activity(entity_id, activity_type, activity_action, metadata, activity_date) 
			VALUES($1, $2, $3, $4, $5) RETURNING id`

	var id string
	err := ac.QueryRow(sqlStatement, t.EntityID, TypeHoldApproved, ActionHold, t.raw(), t.TransactionDate).Scan(&id)

	if err != nil {
		log.Printf("Error creating activity, error:%v", err)
		return nil, err
	}

	return &id, err
}

func (ac transferCreator) HoldReleased(t AccountTransaction) (*string, error) {

	sqlStatement := `INSERT INTO business_activity(entity_id, activity_type, activity_action, metadata, activity_date) 
			VALUES($1, $2, $3, $4, $5) RETURNING id`

	var id string
	err := ac.QueryRow(sqlStatement, t.EntityID, TypeHoldReleased, ActionHoldReleased, t.raw(), t.TransactionDate).Scan(&id)

	if err != nil {
		log.Printf("Error creating activity, error:%v", err)
		return nil, err
	}

	return &id, err
}

func (ac transferCreator) text(a Activity, lang Language) (string, error) {
	var t AccountTransaction
	err := json.Unmarshal(a.Metadata, &t)
	if err != nil {
		log.Printf("Error parsing metadata to construct text for activity:%v error:%v", a, err)
		return "", err
	}

	switch *a.Action {
	case ActionAccountOriginated:
		return ac.composer.Compose(accountOriginatedTempl, t).String(lang)
	case ActionPostedDebit:
		return ac.handleDebitTransferActivity(a, t, lang)
	case ActionInProcessDebit:
		return ac.composer.Compose(accountInProcessDebitTempl, t).String(lang)
	case ActionPostedCredit:
		return ac.handleCreditTransferActivity(a, t, lang)
	case ActionInProcessCredit:
		return ac.composer.Compose(accountInProcessCreditTempl, t).String(lang)
	case ActionHold:
		return ac.composer.Compose(accountHoldApprovedTempl, t).String(lang)
	case ActionHoldReleased:
		return ac.composer.Compose(accountHoldReleasedTempl, t).String(lang)
	default:
		return "", nil
	}
}

func (ac transferCreator) handleCreditTransferActivity(a Activity, t AccountTransaction, lang Language) (string, error) {
	switch a.ActivityType {
	case TypeCardReaderCredit:
		return ac.composer.Compose(accountCardReaderCreditTempl, t).String(lang)
	case TypeCardOnlineCredit:
		return ac.composer.Compose(accountCardCreditTempl, t).String(lang)
	case TypeBankOnlineCredit:
		return ac.composer.Compose(accountBankCreditTempl, t).String(lang)
	case TypeWiseTransferCredit:
		return ac.composer.Compose(accountWiseTransferCreditTempl, t).String(lang)
	case TypeACHTransferShopifyCredit:
		return ac.composer.Compose(accountACHTransferShopifyCreditTempl, t).String(lang)
	case TypeACHTransferCredit:
		if t.ContactName != nil && *t.ContactName != "" {
			return ac.composer.Compose(accountACHTransferCreditTempl, t).String(lang)
		} else {
			return ac.composer.Compose(accountACHTransferCreditGenericTempl, t).String(lang)
		}
	case TypeWireTransferCredit:
		if t.ContactName != nil && *t.ContactName != "" {
			return ac.composer.Compose(accountWireTransferCreditTempl, t).String(lang)
		} else {
			return ac.composer.Compose(accountWireTransferCreditGenericTempl, t).String(lang)
		}
	case TypeCardPullCredit:
		if t.ContactName != nil && *t.ContactName != "" {
			return ac.composer.Compose(accountDebitPullCreditTempl, t).String(lang)
		} else {
			return ac.composer.Compose(accountDebitPullCreditGenericTempl, t).String(lang)
		}
	case TypeInterestTransferCredit:
		return ac.composer.Compose(accountInterestCreditTempl, t).String(lang)
	case TypeCheckCredit:
		return ac.composer.Compose(accountCheckCreditTempl, t).String(lang)
	case TypeDepositCredit:
		return ac.composer.Compose(accountDepositCreditTempl, t).String(lang)
	case TypeOtherCredit:
		return ac.composer.Compose(accountOtherCreditTempl, t).String(lang)
	default:
		return ac.composer.Compose(accountPostedCreditTempl, t).String(lang)
	}
}

func (ac transferCreator) handleDebitTransferActivity(a Activity, t AccountTransaction, lang Language) (string, error) {
	switch a.ActivityType {
	case TypeWiseTransferDebit:
		if t.ContactName != nil && *t.ContactName != "" {
			return ac.composer.Compose(accountWiseTransferDebitTempl, t).String(lang)
		} else {
			return ac.composer.Compose(accountACHTransferDebitGenericTempl, t).String(lang)
		}
	case TypeACHTransferDebit:
		if t.ContactName != nil && *t.ContactName != "" {
			return ac.composer.Compose(accountACHTransferDebitTempl, t).String(lang)
		} else {
			return ac.composer.Compose(accountACHTransferDebitGenericTempl, t).String(lang)
		}
	case TypeCardPushDebit:
		return ac.composer.Compose(accountPushDebitDebitTempl, t).String(lang)
	case TypeCheckDebit:
		return ac.composer.Compose(accountCheckDebitTempl, t).String(lang)
	case TypeFeeDebit:
		if t.ContactName != nil && *t.ContactName != "" {
			return ac.composer.Compose(accountFeeDebitTempl, t).String(lang)
		} else {
			return ac.composer.Compose(accountFeeDebitGenericTempl, t).String(lang)
		}
	default:
		if t.ContactName != nil && *t.ContactName != "" {
			return ac.composer.Compose(accountPostedDebitTempl, t).String(lang)
		} else {
			return ac.composer.Compose(accountPostedDebitGenericTempl, t).String(lang)
		}
	}
}
