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

type BusinessPropertyID string

type BusinessProperty struct {
	// Identifier
	ID BusinessPropertyID `json:"id" db:"id"`

	// Wise partner bank entity id
	BusinessID BusinessID `json:"entityId" db:"business_id"`

	// Property type
	Type bank.BusinessPropertyType `json:"propertyType" db:"property_type"`

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

type BusinessPropertyCreate struct {
	// Wise partner bank entity id
	BusinessID BusinessID `json:"entityId" db:"business_id"`

	// Property type
	Type bank.BusinessPropertyType `json:"propertyType" db:"property_type"`

	// Partner bank property id
	BankID bank.PropertyBankID `json:"bankId" db:"bank_id"`

	// Property value
	Value types.JSONText `json:"propertyValue" db:"property_value"`
}

type BusinessPropertyUpdate struct {
	// Identifier
	ID BusinessPropertyID `json:"id" db:"id"`

	// Property value
	Value types.JSONText `json:"propertyValue" db:"property_value"`
}

type BusinessPropertyService interface {
	Create(BusinessPropertyCreate) (*BusinessProperty, error)
	Update(BusinessPropertyUpdate) (*BusinessProperty, error)
	GetByID(BusinessPropertyID) (*BusinessProperty, error)
	GetByBusinessID(BusinessID, bank.BusinessPropertyType) (*BusinessProperty, error)
}

type businessPropertyService struct {
	sourceReq bank.APIRequest
	bankName  bank.ProviderName
	wdb       *sqlx.DB
	rdb       *sqlx.DB
}

func NewBusinessPropertyService(r bank.APIRequest, name bank.ProviderName) BusinessPropertyService {
	return &businessPropertyService{
		sourceReq: r,
		bankName:  name,
		wdb:       DBWrite,
		rdb:       DBRead,
	}
}

func (s *businessPropertyService) Create(c BusinessPropertyCreate) (*BusinessProperty, error) {
	cc := struct {
		BankName bank.ProviderName `db:"bank_name"`
		BusinessPropertyCreate
	}{
		BankName:               s.bankName,
		BusinessPropertyCreate: c,
	}

	sql := `
        INSERT INTO business_property(business_id, property_type, bank_id, bank_name, property_value)
		VALUES(:business_id, :property_type, :bank_id, :bank_name, :property_value)
        RETURNING id`

	stmt, err := s.wdb.PrepareNamed(sql)
	if err != nil {
		return nil, err
	}

	var id BusinessPropertyID
	err = stmt.Get(&id, &cc)
	if err != nil {
		return nil, err
	}

	return s.GetByID(id)
}

func (s *businessPropertyService) Update(u BusinessPropertyUpdate) (*BusinessProperty, error) {
	var columns []string

	if u.Value != nil {
		columns = append(columns, "property_value = :property_value")
	}

	// No changes requested - return user
	if len(columns) == 0 {
		return s.GetByID(u.ID)
	}

	_, err := s.wdb.NamedExec(fmt.Sprintf("UPDATE business_property SET %s WHERE id = '%s'", strings.Join(columns, ", "), u.ID), u)
	if err != nil {
		log.Println(err)
		return nil, errors.Cause(err)
	}

	return s.GetByID(u.ID)
}

func (s *businessPropertyService) GetByID(id BusinessPropertyID) (*BusinessProperty, error) {
	e := BusinessProperty{}

	err := s.wdb.Get(&e, "SELECT * FROM business_property WHERE id = $1", id)
	if err != nil {
		return nil, err
	}

	return &e, nil
}

func (s *businessPropertyService) GetByBusinessID(id BusinessID, propertyType bank.BusinessPropertyType) (*BusinessProperty, error) {
	var p BusinessProperty

	err := s.wdb.Get(&p, "SELECT * FROM business_property WHERE business_id = $1 AND property_type = $2", id, propertyType)
	if err != nil {
		return nil, err
	}

	return &p, nil
}
