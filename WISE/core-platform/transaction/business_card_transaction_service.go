package transaction

import (
	"errors"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/wiseco/core-platform/shared"
)

//CardService
type CardService interface {

	//Create a business card transaction
	Create(*BusinessCardTransactionCreate) (BusinessCardTransaction, error)

	//Create card hold transaction
	CreateHold(*BusinessHoldTransactionCreate) (BusinessHoldTransaction, error)
}

type cardService struct {
	*sqlx.DB
}

//NewCardService a new card transaction service
func NewCardService() CardService {
	return cardService{DBWrite}
}

func (s cardService) Create(t *BusinessCardTransactionCreate) (BusinessCardTransaction, error) {
	var transaction BusinessCardTransaction
	if t == nil {
		return transaction, errors.New("transaction can not be nil")
	}

	values := shared.SQLGenInsertValues(*t)
	keys := shared.SQLGenInsertKeys(*t)

	query := fmt.Sprintf("INSERT INTO business_card_transaction(%s) VALUES(%s) RETURNING *", keys, values)
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

func (s cardService) CreateHold(t *BusinessHoldTransactionCreate) (BusinessHoldTransaction, error) {
	var transaction BusinessHoldTransaction
	if t == nil {
		return transaction, errors.New("transaction can not be nil")
	}
	keys := shared.SQLGenInsertKeys(*t)
	values := shared.SQLGenInsertValues(*t)

	query := fmt.Sprintf(" INSERT INTO business_hold_transaction (%s) VALUES(%s) RETURNING *", keys, values)
	stmt, err := s.PrepareNamed(query)
	err = stmt.Get(&transaction, t)
	return transaction, err
}
