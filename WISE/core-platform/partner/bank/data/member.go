package data

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/wiseco/core-platform/partner/bank"
)

type MemberID string

type BusinessMemberService interface {
	Create(BusinessMemberCreate) (*BusinessMember, error)
	Update(BusinessMemberUpdate) (*BusinessMember, error)
	GetByID(MemberID) (*BusinessMember, error)
	GetByConsumerID(bank.BusinessID, bank.ConsumerID) (*BusinessMember, error)
	GetByBankID(bank.MemberBankID) (*BusinessMember, error)
	Delete(MemberID) error
}

type businessMemberService struct {
	sourceReq bank.APIRequest
	bankName  bank.ProviderName
	wdb       *sqlx.DB
	rdb       *sqlx.DB
}

func NewBusinessMemberService(r bank.APIRequest, name bank.ProviderName) BusinessMemberService {
	return &businessMemberService{
		sourceReq: r,
		bankName:  name,
		wdb:       DBWrite,
		rdb:       DBRead,
	}
}

type BusinessMember struct {
	// Identifier
	ID MemberID `json:"id" db:"id"`

	// Related Consumer ID
	ConsumerID ConsumerID `json:"consumerId" db:"consumer_id"`

	// Business id
	BusinessID bank.BusinessID `json:"businessId" db:"business_id"`

	// Property bank id
	BankID bank.MemberBankID `json:"bankId" db:"bank_id"`

	// Consumer bank id
	ConsumerBankID bank.ConsumerBankID `json:"consumerBankId" db:"consumer_bank_id"`

	// Business bank id
	BusinessBankID bank.BusinessBankID `json:"businessBankId" db:"business_bank_id"`

	// Bank name
	BankName bank.ProviderName `json:"bankName" db:"bank_name"`

	// KYC Status
	KYCStatus bank.KYCStatus `json:"kycStatus" db:"kyc_status"`

	// Created timestamp
	Created time.Time `json:"created" db:"created"`

	// Modified timestamp
	Modified time.Time `json:"modified" db:"modified"`
}

type BusinessMemberCreate struct {
	// Business entity id
	BusinessID BusinessID `json:"businessId" db:"business_id"`

	// Related Consumer ID
	ConsumerID ConsumerID `json:"consumerId" db:"consumer_id"`

	// Property bank id
	BankID *bank.MemberBankID `json:"bankId" db:"bank_id"`

	// Bank control person ID
	BankControlID *bank.MemberBankID `json:"bankControlId" db:"bank_control_id"`
}

type BusinessMemberUpdate struct {
	// Identifier
	ID MemberID `json:"id" db:"id"`

	// Property bank id
	BankID *bank.MemberBankID `json:"bankId" db:"bank_id"`

	// Bank control person ID
	BankControlID *bank.MemberBankID `json:"bankControlId" db:"bank_control_id"`
}

func (s *businessMemberService) Create(c BusinessMemberCreate) (*BusinessMember, error) {
	cc := struct {
		BankName bank.ProviderName `db:"bank_name"`
		BusinessMemberCreate
	}{
		BankName:             s.bankName,
		BusinessMemberCreate: c,
	}

	sql := `
        INSERT INTO business_member(consumer_id, business_id, bank_id, bank_name)
		VALUES(:consumer_id, :business_id, :bank_id, :bank_name)
        RETURNING id`

	stmt, err := s.wdb.PrepareNamed(sql)
	if err != nil {
		return nil, err
	}

	var id MemberID
	err = stmt.Get(&id, &cc)
	if err != nil {
		return nil, err
	}

	return s.GetByID(id)
}

func (s *businessMemberService) Update(u BusinessMemberUpdate) (*BusinessMember, error) {
	var columns []string

	if u.BankID != nil {
		columns = append(columns, "bank_id = :bank_id")
	}

	if u.BankControlID != nil {
		columns = append(columns, "bank_control_id = :bank_control_id")
	}

	// No changes requested - return business member
	if len(columns) == 0 {
		return s.GetByID(u.ID)
	}
	_, err := s.wdb.NamedExec(fmt.Sprintf("UPDATE business_member SET %s WHERE id = '%s'", strings.Join(columns, ", "), u.ID), u)
	if err != nil {
		log.Println(err)
		return nil, errors.Cause(err)
	}

	return s.GetByID(u.ID)
}

func (s *businessMemberService) GetByID(id MemberID) (*BusinessMember, error) {
	var m BusinessMember
	err := s.wdb.Get(
		&m, `
		SELECT
			business_member.id,
			business_member.consumer_id, 
			business_member.bank_id,
			consumer.bank_id AS consumer_bank_id,
			business.bank_id AS business_bank_id,
			business_member.bank_name,
			consumer.kyc_status,
			business_member.created,
			business_member.modified
		FROM
			business_member
		JOIN
			consumer ON business_member.consumer_id = consumer.id
		JOIN
			business ON business_member.business_id = business.id
		WHERE
			business_member.id = $1`,
		id,
	)
	if err != nil {
		return nil, err
	}

	return &m, nil
}

func (s *businessMemberService) GetByConsumerID(bid bank.BusinessID, cid bank.ConsumerID) (*BusinessMember, error) {
	var m BusinessMember
	err := s.wdb.Get(
		&m, `
        SELECT
            business_member.id,
            business_member.consumer_id, 
            business_member.bank_id,
            consumer.bank_id AS consumer_bank_id,
            business.bank_id AS business_bank_id,
            business_member.bank_name,
            business_member.created,
            business_member.modified
        FROM
            business_member
        JOIN
            consumer ON business_member.consumer_id = consumer.id
        JOIN
            business ON business_member.business_id = business.id
        WHERE
			business.business_id = $1 AND consumer.consumer_id = $2`,
		bid, cid,
	)
	if err != nil {
		return nil, err
	}

	return &m, nil
}

func (s *businessMemberService) GetByBankID(id bank.MemberBankID) (*BusinessMember, error) {
	var m BusinessMember
	err := s.wdb.Get(
		&m, `
        SELECT
            business_member.id,
            business_member.consumer_id,
            business_member.bank_id,
            consumer.bank_id AS consumer_bank_id,
            business.bank_id AS business_bank_id,
            business_member.bank_name,
            business_member.created,
            business_member.modified
        FROM
            business_member
        JOIN
            consumer ON business_member.consumer_id = consumer.id
        JOIN
            business ON business_member.business_id = business.id
        WHERE
            business_member.bank_id = $1`,
		id,
	)
	if err != nil {
		return nil, err
	}

	return &m, nil
}

func (s *businessMemberService) Delete(id MemberID) error {
	_, err := s.wdb.Exec("DELETE FROM business_member WHERE id = $1 AND bank_name = $2", id, s.bankName)
	return err
}
