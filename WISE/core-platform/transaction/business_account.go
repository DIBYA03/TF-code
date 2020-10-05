package transaction

import (
	"github.com/jmoiron/sqlx"
)

//AccountService the business transaction service
type AccountService interface {
	GetByID(id string, accountID string) (BusinessPostedTransaction, error)
	List(businessID string, accountID string) ([]BusinessPostedTransaction, error)
}

type accountService struct {
	*sqlx.DB
}

//NewAccountService returns a new business transaction service
func NewAccountService() AccountService {
	return &accountService{DBWrite}
}

func (a accountService) GetByID(id string, accountID string) (BusinessPostedTransaction, error) {
	var t BusinessPostedTransaction
	query := `SELECT * FROM business_transaction
		JOIN business_card_transaction 
		ON business_transaction.id = business_card_transaction.transaction_id
		JOIN business_hold_transaction 
		ON business_transaction.id = business_hold_transaction.transaction_id
	    WHERE  id = $1 AND account_id = $2`

	err := a.Get(&t, query, id, accountID)
	return t, err
}

func (a accountService) List(businessID string, accountID string) ([]BusinessPostedTransaction, error) {
	var list []BusinessPostedTransaction
	query := `SELECT * FROM business_transaction
		JOIN business_card_transaction 
		ON business_transaction.id = business_card_transaction.transaction_id
		JOIN business_hold_transaction 
		ON business_transaction.id = business_hold_transaction.transaction_id
	    WHERE  business_id = $1 AND account_id = $2`

	err := a.Select(&list, query, businessID, accountID)
	return list, err
}
