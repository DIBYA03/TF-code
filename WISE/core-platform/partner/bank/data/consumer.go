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

type ConsumerID string

type Consumer struct {
	// Identifier
	ID ConsumerID `json:"id" db:"id"`

	// Internal consumer id
	ConsumerID bank.ConsumerID `json:"consumerId" db:"consumer_id"`

	// Partner bank name
	BankName bank.ProviderName `json:"bankName" db:"bank_name"`

	// Partner bank consumer id
	BankID bank.ConsumerBankID `json:"bankId" db:"bank_id"`

	// Partner bank extra data
	BankExtra types.JSONText `json:"-" db:"bank_extra"`

	// KYC Status
	KYCStatus bank.KYCStatus `json:"kycStatus" db:"kyc_status"`

	// Created timestamp
	Created time.Time `json:"created" db:"created"`

	// Modified timestamp
	Modified time.Time `json:"modified" db:"modified"`
}

type ConsumerCreate struct {
	// Internal consumer id
	ConsumerID bank.ConsumerID `json:"consumerId" db:"consumer_id"`

	// Partner bank entity id
	BankID bank.ConsumerBankID `json:"bankId" db:"bank_id"`

	// Partner bank extra data
	BankExtra types.JSONText `json:"-" db:"bank_extra"`

	// KYC Status
	KYCStatus bank.KYCStatus `json:"kycStatus" db:"kyc_status"`
}

type ConsumerUpdate struct {
	// Identifier
	ID ConsumerID `json:"id" db:"id"`

	// Partner bank extra data
	BankExtra types.NullJSONText `json:"-" db:"bank_extra"`

	// KYC Status
	KYCStatus *bank.KYCStatus `json:"kycStatus" db:"kyc_status"`
}

type ConsumerService interface {
	Create(ConsumerCreate) (*Consumer, error)
	Update(ConsumerUpdate) (*Consumer, error)

	GetByID(ConsumerID) (*Consumer, error)
	GetByConsumerID(bank.ConsumerID) (*Consumer, error)
	GetByBankID(bank.ConsumerBankID) (*Consumer, error)

	DeleteByID(bank.ConsumerID) error
}

type consumerService struct {
	bankName  bank.ProviderName
	sourceReq bank.APIRequest
	rdb       *sqlx.DB
	wdb       *sqlx.DB
}

func NewConsumerService(r bank.APIRequest, n bank.ProviderName) ConsumerService {
	return &consumerService{
		bankName:  n,
		sourceReq: r,
		rdb:       DBRead,
		wdb:       DBWrite,
	}
}

func (s *consumerService) Create(c ConsumerCreate) (*Consumer, error) {

	cc := struct {
		BankName bank.ProviderName `db:"bank_name"`
		ConsumerCreate
	}{
		BankName:       s.bankName,
		ConsumerCreate: c,
	}

	sql := `
        INSERT INTO consumer (consumer_id, bank_name, bank_id, bank_extra, kyc_status)
        VALUES (:consumer_id, :bank_name, :bank_id, :bank_extra, :kyc_status)
        RETURNING id`

	stmt, err := s.wdb.PrepareNamed(sql)
	if err != nil {
		return nil, err
	}

	var id ConsumerID
	err = stmt.Get(&id, &cc)
	if err != nil {
		return nil, err
	}

	return s.GetByID(id)
}

func (s *consumerService) Update(u ConsumerUpdate) (*Consumer, error) {
	var columns []string

	if u.BankExtra.Valid {
		columns = append(columns, "bank_extra = :bank_extra")
	}

	if u.KYCStatus != nil {
		columns = append(columns, "kyc_status = :kyc_status")
	}

	// No changes requested - return user
	if len(columns) == 0 {
		return s.GetByID(u.ID)
	}

	_, err := s.wdb.NamedExec(fmt.Sprintf("UPDATE consumer SET %s WHERE id = '%s'", strings.Join(columns, ", "), u.ID), u)
	if err != nil {
		log.Println(err)
		return nil, errors.Cause(err)
	}

	return s.GetByID(u.ID)
}

func (s *consumerService) GetByID(id ConsumerID) (*Consumer, error) {
	e := Consumer{}
	if err := s.wdb.Get(&e, "SELECT * FROM consumer WHERE id = $1", id); err != nil {
		return nil, err
	}

	return &e, nil
}

func (s *consumerService) GetByConsumerID(id bank.ConsumerID) (*Consumer, error) {
	e := Consumer{}
	if err := s.wdb.Get(&e, "SELECT * FROM consumer WHERE consumer_id = $1 AND bank_name = $2", id, s.bankName); err != nil {
		return nil, err
	}

	return &e, nil
}

func (s *consumerService) GetByBankID(id bank.ConsumerBankID) (*Consumer, error) {
	e := Consumer{}
	if err := s.wdb.Get(&e, "SELECT * FROM consumer WHERE bank_id = $1 AND bank_name = $2", id, s.bankName); err != nil {
		return nil, err
	}

	return &e, nil
}

func (s *consumerService) DeleteByID(id bank.ConsumerID) error {
	tx := s.wdb.MustBegin()
	bm := `DELETE FROM business_member WHERE consumer_id IN (SELECT id FROM consumer WHERE consumer_id = $1)`
	tx.MustExec(bm, id)
	cp := `DELETE FROM consumer_property WHERE consumer_id IN (SELECT id FROM consumer WHERE consumer_id = $1)`
	tx.MustExec(cp, id)
	c := `DELETE FROM consumer WHERE consumer_id = $1`
	tx.MustExec(c, id)
	return tx.Commit()
}
