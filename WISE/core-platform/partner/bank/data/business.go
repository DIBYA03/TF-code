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

type BusinessID string

type Business struct {
	// Identifier
	ID BusinessID `json:"id" db:"id"`

	// Wise entity id
	BusinessID bank.BusinessID `json:"businessId" db:"business_id"`

	// Partner bank name
	BankName bank.ProviderName `json:"bankName" db:"bank_name"`

	// Partner bank entity id
	BankID bank.BusinessBankID `json:"bankId" db:"bank_id"`

	// Partner bank extra data
	BankExtra types.JSONText `json:"-" db:"bank_extra"`

	// KYC Status
	KYCStatus bank.KYCStatus `json:"kycStatus" db:"kyc_status"`

	// Created timestamp
	Created time.Time `json:"created" db:"created"`

	// Modified timestamp
	Modified time.Time `json:"modified" db:"modified"`
}

type BusinessCreate struct {
	// business id
	BusinessID bank.BusinessID `json:"businessId" db:"business_id"`

	// Partner bank entity id
	BankID bank.BusinessBankID `json:"bankId" db:"bank_id"`

	// Partner bank extra data
	BankExtra types.JSONText `json:"-" db:"bank_extra"`

	// KYC Status
	KYCStatus *bank.KYCStatus `json:"kycStatus" db:"kyc_status"`
}

type BusinessUpdate struct {
	// Identifier
	ID BusinessID `json:"id" db:"id"`

	// Partner bank extra data
	BankExtra types.NullJSONText `json:"-" db:"bank_extra"`

	// KYC Status
	KYCStatus *bank.KYCStatus `json:"kycStatus" db:"kyc_status"`
}

type BusinessService interface {
	Create(BusinessCreate) (*Business, error)
	Update(BusinessUpdate) (*Business, error)

	GetByID(BusinessID) (*Business, error)
	GetByBusinessID(bank.BusinessID) (*Business, error)
	GetByBankID(bank.BusinessBankID) (*Business, error)
}

type businessService struct {
	bankName  bank.ProviderName
	sourceReq bank.APIRequest
	rdb       *sqlx.DB
	wdb       *sqlx.DB
}

func NewBusinessService(r bank.APIRequest, n bank.ProviderName) BusinessService {
	return &businessService{
		bankName:  n,
		sourceReq: r,
		rdb:       DBRead,
		wdb:       DBWrite,
	}
}

func (s *businessService) Create(c BusinessCreate) (*Business, error) {
	cc := struct {
		BankName bank.ProviderName `db:"bank_name"`
		BusinessCreate
	}{
		BankName:       s.bankName,
		BusinessCreate: c,
	}

	sql := `
        INSERT INTO business (business_id, bank_name, bank_id, bank_extra, kyc_status)
        VALUES (:business_id, :bank_name, :bank_id, :bank_extra, :kyc_status)
        RETURNING id`

	stmt, err := s.wdb.PrepareNamed(sql)
	if err != nil {
		return nil, err
	}

	var id BusinessID
	err = stmt.Get(&id, &cc)
	if err != nil {
		return nil, err
	}

	return s.GetByID(id)
}

func (s *businessService) Update(u BusinessUpdate) (*Business, error) {
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

	_, err := s.wdb.NamedExec(fmt.Sprintf("UPDATE business SET %s WHERE id = '%s'", strings.Join(columns, ", "), u.ID), u)
	if err != nil {
		log.Println(err)
		return nil, errors.Cause(err)
	}

	return s.GetByID(u.ID)
}

func (s *businessService) GetByID(id BusinessID) (*Business, error) {
	var e Business
	if err := s.wdb.Get(&e, "SELECT * FROM business WHERE id = $1", id); err != nil {
		return nil, err
	}

	return &e, nil
}

func (s *businessService) GetByBusinessID(id bank.BusinessID) (*Business, error) {
	var b Business
	if err := s.wdb.Get(&b, "SELECT * FROM business WHERE business_id = $1 AND bank_name = $2", id, s.bankName); err != nil {
		return nil, err
	}

	return &b, nil
}

func (s *businessService) GetByBankID(id bank.BusinessBankID) (*Business, error) {
	var b Business
	if err := s.wdb.Get(&b, "SELECT * FROM business WHERE bank_id = $1 AND bank_name = $2", id, s.bankName); err != nil {
		return nil, err
	}

	return &b, nil
}
