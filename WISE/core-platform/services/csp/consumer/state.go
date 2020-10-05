package consumer

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/wiseco/core-platform/services"
	"github.com/wiseco/core-platform/services/csp/data"
)

// StateService ...
type StateService interface {
	Create(ConsumerStateCreate) (ConsumerState, error)
	List(string) ([]ConsumerState, error)
	ByID(string) (ConsumerState, error)
}

type stateService struct {
	rdb *sqlx.DB
	wdb *sqlx.DB
}

// NewStateService ..
func NewStateService() StateService {
	return stateService{wdb: data.DBWrite, rdb: data.DBRead}
}

func (s stateService) Create(create ConsumerStateCreate) (ConsumerState, error) {
	var item ConsumerState
	keys := services.SQLGenInsertKeys(create)
	values := services.SQLGenInsertValues(create)

	q := fmt.Sprintf("INSERT INTO consumer_state (%s) VALUES(%s) RETURNING *", keys, values)
	stmt, err := s.wdb.PrepareNamed(q)
	if err != nil {
		return item, err
	}
	err = stmt.Get(&item, create)

	return item, err
}

func (s stateService) List(consumerID string) ([]ConsumerState, error) {
	list := make([]ConsumerState, 0)
	err := s.rdb.Select(&list, "SELECT * FROM consumer_state WHERE consumer_id = $1", consumerID)
	if err != nil && err == sql.ErrNoRows {
		return list, nil
	}

	return list, err
}

func (s stateService) ByID(id string) (ConsumerState, error) {
	var item ConsumerState
	err := s.rdb.Get(&item, "SELECT * FROM consumer_state WHERE id = $1", id)
	if err == sql.ErrNoRows {
		return item, services.ErrorNotFound{}.New("")
	}
	return item, err
}
