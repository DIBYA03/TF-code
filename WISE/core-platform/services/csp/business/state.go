package business

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/wiseco/core-platform/services"
	"github.com/wiseco/core-platform/services/csp/data"
)

// StateService ...
type StateService interface {
	Create(create BusinessStateCreate) (BusinessState, error)
	List(string) ([]BusinessState, error)
	ByID(string) (BusinessState, error)
}

type stateService struct {
	rdb *sqlx.DB
	wdb *sqlx.DB
}

// NewStateService ..
func NewStateService() StateService {
	return stateService{wdb: data.DBWrite, rdb: data.DBRead}
}
func (s stateService) Create(create BusinessStateCreate) (BusinessState, error) {
	var item BusinessState
	keys := services.SQLGenInsertKeys(create)
	values := services.SQLGenInsertValues(create)

	q := fmt.Sprintf("INSERT INTO business_state (%s) VALUES(%s) RETURNING *", keys, values)
	stmt, err := s.wdb.PrepareNamed(q)
	if err != nil {
		return item, err
	}
	err = stmt.Get(&item, create)

	return item, err
}

func (s stateService) List(businessID string) ([]BusinessState, error) {
	list := make([]BusinessState, 0)
	err := s.rdb.Select(&list, "SELECT * FROM business_state WHERE business_id = $1", businessID)
	if err != nil && err == sql.ErrNoRows {
		return list, nil
	}

	return list, err
}

func (s stateService) ByID(id string) (BusinessState, error) {
	var item BusinessState
	err := s.rdb.Get(&item, "SELECT * FROM business_state WHERE id = $1", id)
	if err == sql.ErrNoRows {
		return item, services.ErrorNotFound{}.New("")
	}
	return item, err
}
