package data

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/types"
	"github.com/pkg/errors"
	"github.com/wiseco/core-platform/partner/bank"
)

type ConsumerPropertyID string

type ConsumerProperty struct {
	// Identifier
	ID ConsumerPropertyID `json:"id" db:"id"`

	// Wise partner bank entity id
	ConsumerID ConsumerID `json:"entityId" db:"consumer_id"`

	// Property type
	Type bank.ConsumerPropertyType `json:"propertyType" db:"property_type"`

	// Partner bank property id
	BankID bank.PropertyBankID `json:"bankId" db:"bank_id"`

	// Partner bank name
	BankName bank.ProviderName `json:"bankName" db:"bank_name"`

	// Property value
	Value types.JSONText `json:"propertyValue" db:"property_value"`

	// Created timestamp
	Created time.Time `json:"created" db:"created"`

	// Modified timestamp
	Modified time.Time `json:"modified" db:"modified"`
}

type ConsumerPropertyCreate struct {
	// Wise partner bank entity id
	ConsumerID ConsumerID `json:"entityId" db:"consumer_id"`

	// Property type
	Type bank.ConsumerPropertyType `json:"propertyType" db:"property_type"`

	// Partner bank property id
	BankID bank.PropertyBankID `json:"bankId" db:"bank_id"`

	// Property value
	Value types.JSONText `json:"propertyValue" db:"property_value"`
}

type ConsumerPropertyUpdate struct {
	// Identifier
	ID ConsumerPropertyID `json:"id" db:"id"`

	// Property value
	Value types.JSONText `json:"propertyValue" db:"property_value"`
}

type ConsumerPropertyService interface {
	Create(ConsumerPropertyCreate) (*ConsumerProperty, error)
	Update(ConsumerPropertyUpdate) (*ConsumerProperty, error)
	GetByID(ConsumerPropertyID) (*ConsumerProperty, error)
	GetByConsumerID(ConsumerID, bank.ConsumerPropertyType) (*ConsumerProperty, error)
}

type consumerPropertyService struct {
	sourceReq bank.APIRequest
	bankName  bank.ProviderName
	wdb       *sqlx.DB
	rdb       *sqlx.DB
}

func NewConsumerPropertyService(r bank.APIRequest, name bank.ProviderName) ConsumerPropertyService {
	return &consumerPropertyService{
		sourceReq: r,
		bankName:  name,
		wdb:       DBWrite,
		rdb:       DBRead,
	}
}

func (s *consumerPropertyService) Create(c ConsumerPropertyCreate) (*ConsumerProperty, error) {
	cc := struct {
		BankName bank.ProviderName `db:"bank_name"`
		ConsumerPropertyCreate
	}{
		BankName:               s.bankName,
		ConsumerPropertyCreate: c,
	}

	sql := `
        INSERT INTO consumer_property(consumer_id, property_type, bank_id, bank_name, property_value)
		VALUES(:consumer_id, :property_type, :bank_id, :bank_name, :property_value)
        RETURNING id`

	stmt, err := s.wdb.PrepareNamed(sql)
	if err != nil {
		return nil, err
	}

	var id ConsumerPropertyID
	err = stmt.Get(&id, &cc)
	if err != nil {
		return nil, err
	}

	return s.GetByID(id)
}

func (s *consumerPropertyService) Update(u ConsumerPropertyUpdate) (*ConsumerProperty, error) {
	var columns []string

	if u.Value != nil {
		columns = append(columns, "property_value = :property_value")
	}

	// No changes requested - return user
	if len(columns) == 0 {
		return s.GetByID(u.ID)
	}

	_, err := s.wdb.NamedExec(fmt.Sprintf("UPDATE consumer_property SET %s WHERE id = '%s'", strings.Join(columns, ", "), u.ID), u)
	if err != nil {
		log.Println(err)
		return nil, errors.Cause(err)
	}

	return s.GetByID(u.ID)
}

func (s *consumerPropertyService) GetByID(id ConsumerPropertyID) (*ConsumerProperty, error) {
	e := ConsumerProperty{}

	err := s.wdb.Get(&e, "SELECT * FROM consumer_property WHERE id = $1", id)
	if err != nil {
		return nil, err
	}

	return &e, nil
}

func (s *consumerPropertyService) GetByConsumerID(id ConsumerID, propertType bank.ConsumerPropertyType) (*ConsumerProperty, error) {
	e := ConsumerProperty{}

	err := s.wdb.Get(&e, "SELECT * FROM consumer_property WHERE consumer_id = $1 AND property_type = $2", id, propertType)
	if err != nil {
		return nil, err
	}

	return &e, nil
}
