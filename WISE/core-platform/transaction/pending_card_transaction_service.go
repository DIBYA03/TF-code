package transaction

import (
	"errors"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/wiseco/core-platform/shared"
)

type PendingCardTransactionService interface {
	//Create a business card transaction
	Create(*BusinessCardTransactionCreate) (BusinessCardTransaction, error)

	//Create card hold transaction
	CreateHold(*BusinessHoldTransactionCreate) (BusinessHoldTransaction, error)
}

type pendingCardService struct {
	*sqlx.DB
}

//NewCardService a new card transaction service
func NewPendingCardService() PendingCardTransactionService {
	return pendingCardService{DBWrite}
}

func (s pendingCardService) Create(t *BusinessCardTransactionCreate) (BusinessCardTransaction, error) {

	var transaction BusinessCardTransaction
	if t == nil {
		return transaction, errors.New("transaction can not be nil")
	}

	values := shared.SQLGenInsertValues(*t)
	keys := shared.SQLGenInsertKeys(*t)

	query := fmt.Sprintf("INSERT INTO business_card_pending_transaction(%s) VALUES(%s) RETURNING *", keys, values)
	stmt, err := s.PrepareNamed(query)
	if err != nil {
		log.Println(err)
		return transaction, err
	}

	err = stmt.Get(&transaction, t)
	if err != nil {
		log.Println(err)
		return transaction, err
	}

	return transaction, err
}

func (s pendingCardService) CreateHold(t *BusinessHoldTransactionCreate) (BusinessHoldTransaction, error) {
	var transaction BusinessHoldTransaction
	if t == nil {
		return transaction, errors.New("transaction can not be nil")
	}
	keys := shared.SQLGenInsertKeys(*t)
	values := shared.SQLGenInsertValues(*t)

	query := fmt.Sprintf(" INSERT INTO business_hold_pending_transaction (%s) VALUES(%s) RETURNING *", keys, values)
	stmt, err := s.PrepareNamed(query)
	err = stmt.Get(&transaction, t)
	return transaction, err
}
