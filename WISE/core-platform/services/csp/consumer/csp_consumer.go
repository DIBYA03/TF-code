package consumer

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/wiseco/core-platform/services"
	"github.com/wiseco/core-platform/services/csp/data"
	"github.com/wiseco/core-platform/shared"
)

// CSPService csp consumer
type CSPService interface {
	UpdateStatus(consumerID shared.ConsumerID, status string) error
	GetAll(params map[string]interface{}) ([]CSPConsumer, error)
	CSPConsumerByID(id string) (CSPConsumer, error)
	ByConsumerID(id string) (*CSPConsumer, error)
	CSPConsumerCreate(CSPConsumerCreate) (CSPConsumer, error)
	CSPConsumerByConsumerID(id shared.ConsumerID) (CSPConsumer, error)
	CSPConsumerUpdateByConsumerID(shared.ConsumerID, CSPConsumerUpdate) (CSPConsumer, error)
	List(status string, limit, offset int) ([]CSPConsumer, error)
}

type cspConsumerService struct {
	rdb *sqlx.DB
	wdb *sqlx.DB
}

//NewCSPService ..
func NewCSPService() CSPService {
	return cspConsumerService{wdb: data.DBWrite, rdb: data.DBRead}
}

func (s cspConsumerService) CSPConsumerCreate(create CSPConsumerCreate) (CSPConsumer, error) {
	var consumer CSPConsumer
	keys := services.SQLGenInsertKeys(create)
	values := services.SQLGenInsertValues(create)

	q := fmt.Sprintf("INSERT INTO consumer (%s) VALUES(%s) RETURNING *", keys, values)
	stmt, err := s.wdb.PrepareNamed(q)
	if err != nil {
		return consumer, err
	}
	err = stmt.Get(&consumer, create)
	if err == nil {
		if _, err := NewStateService().Create(ConsumerStateCreate{
			ConsumerID:   consumer.ID,
			ReviewStatus: consumer.Status,
		}); err != nil {
			log.Printf("Error creating consumer state %v", err)
		}
	}
	return consumer, err
}

func (s cspConsumerService) GetAll(params map[string]interface{}) ([]CSPConsumer, error) {
	clause := ""

	if params["submitStart"] != nil && params["submitStart"] != "" {
		clause = "created >= '" + params["submitStart"].(string) + "'"
	}

	if params["submitEnd"] != nil && params["submitEnd"] != "" {
		endQuery := "created <= '" + params["submitEnd"].(string) + "'"
		if len(clause) > 0 {
			clause = clause + " AND " + endQuery
		} else {
			clause = endQuery
		}

	}

	if len(clause) > 0 {
		clause = " WHERE " + clause
	}
	var items []CSPConsumer

	q := "SELECT * FROM consumer" + clause
	err := s.rdb.Select(&items, q)
	if err == sql.ErrNoRows {
		log.Printf("no business %v", err)
		return items, services.ErrorNotFound{}.New("")
	}
	return items, err
}

func (s cspConsumerService) CSPConsumerByID(id string) (CSPConsumer, error) {
	var item CSPConsumer
	err := s.rdb.Get(&item, "SELECT * FROM consumer WHERE id = $1", id)
	if err == sql.ErrNoRows {
		log.Printf("no business %v", err)
		return item, services.ErrorNotFound{}.New("")
	}
	return item, err
}

func (s cspConsumerService) ByConsumerID(id string) (*CSPConsumer, error) {
	var item CSPConsumer
	err := s.rdb.Get(&item, "SELECT * FROM consumer WHERE consumer_id = $1", id)
	if err == sql.ErrNoRows {
		log.Printf("no business %v", err)
		return nil, services.ErrorNotFound{}.New("")
	}
	return &item, err
}

func (s cspConsumerService) CSPConsumerByConsumerID(id shared.ConsumerID) (CSPConsumer, error) {
	var item CSPConsumer
	err := s.rdb.Get(&item, "SELECT * FROM consumer WHERE consumer_id = $1", id)
	return item, err
}

func (s cspConsumerService) List(status string, limit, offset int) ([]CSPConsumer, error) {
	var list []CSPConsumer
	err := s.rdb.Select(&list, "SELECT * FROM consumer WHERE review_status = $1 ORDER BY created DESC LIMIT $2 OFFSET $3", status, limit, offset)
	return list, err
}

// Update csp consumer status
func (s cspConsumerService) UpdateStatus(consumerID shared.ConsumerID, status string) error {
	updates := CSPConsumerUpdate{
		Status: &status,
	}
	_, err := s.CSPConsumerUpdateByConsumerID(consumerID, updates)
	return err
}

func (s cspConsumerService) CSPConsumerUpdateByConsumerID(id shared.ConsumerID, updates CSPConsumerUpdate) (CSPConsumer, error) {
	var consumer CSPConsumer
	current, err := s.CSPConsumerByConsumerID(id)
	if err != nil {
		return consumer, err
	}
	keys := services.SQLGenForUpdate(updates)
	q := fmt.Sprintf("UPDATE consumer SET %s WHERE consumer_id = '%s' RETURNING *", keys, id)
	stmt, err := s.wdb.PrepareNamed(q)
	if err != nil {
		log.Printf("error preparing stmt %v query %v", err, q)
		return consumer, err
	}

	err = stmt.Get(&consumer, updates)
	if err != nil {
		return consumer, err
	}

	if current.Status != consumer.Status {
		if updates.Status != nil {
			// cant update status with empty string
			if *updates.Status == "" {
				return consumer, nil
			}
			if _, err := NewStateService().Create(ConsumerStateCreate{
				ConsumerID:   consumer.ID,
				ReviewStatus: consumer.Status,
			}); err != nil {
				log.Printf("Error creating consumer state %v", err)
			}
		}
	}

	return consumer, nil
}
