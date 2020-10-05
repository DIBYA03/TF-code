/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

// Package for all business member related services
package business

import (
	"database/sql"
	"fmt"
	"log"
	"net/mail"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/ttacon/libphonenumber"
	"github.com/wiseco/core-platform/auth"
	"github.com/wiseco/core-platform/services"
	"github.com/wiseco/core-platform/services/data"
	"github.com/wiseco/core-platform/services/user"
	"github.com/wiseco/core-platform/shared"
)

type memberDataStore struct {
	sourceRequest services.SourceRequest
	*sqlx.DB
}

type MemberService interface {
	// Fetch
	List(offset int, limit int, id shared.BusinessID) ([]BusinessMember, error)
	ListInternal(offset int, limit int, id shared.BusinessID) ([]BusinessMember, error)

	GetById(id shared.BusinessMemberID, businessId shared.BusinessID) (*BusinessMember, error)

	// Create Member
	Create(create *BusinessMemberCreate) (*BusinessMember, error)
	Update(id shared.BusinessMemberID, businessId shared.BusinessID, u *BusinessMemberUpdate) (*BusinessMember, error)

	// Submit
	Submit(id shared.BusinessMemberID, businessId shared.BusinessID) (*BusinessMember, error)
	StartVerification(id shared.BusinessMemberID, businessId shared.BusinessID) (*MemberKYCResponse, error)
	GetVerification(id shared.BusinessMemberID, businessId shared.BusinessID) (*MemberKYCResponse, error)

	// Deactivate Member
	Deactivate(id shared.BusinessMemberID, businessId shared.BusinessID) error
}

func NewMemberService(r services.SourceRequest) MemberService {
	return &memberDataStore{r, data.DBWrite}
}

func NewMemberServiceWithout() MemberService {
	return &memberDataStore{services.SourceRequest{}, data.DBWrite}
}

func (db *memberDataStore) count() (int, error) {
	var count int

	err := db.QueryRow("SELECT COUNT(*) FROM business_member WHERE deactivated IS NULL").Scan(&count)

	if err != nil {
		return 0, errors.Cause(err)
	}

	return count, err
}

func (db *memberDataStore) List(offset int, limit int, id shared.BusinessID) ([]BusinessMember, error) {
	// Check access
	err := auth.NewAuthService(db.sourceRequest).CheckBusinessAccess(id)
	if err != nil {
		return nil, err
	}

	return db.ListInternal(offset, limit, id)
}

func (db *memberDataStore) ListInternal(offset int, limit int, id shared.BusinessID) ([]BusinessMember, error) {
	rows := []BusinessMember{}
	err := db.Select(
		&rows, `
		SELECT
            business_member.id, business_member.consumer_id, business_member.business_id,
            business_member.title_type, business_member.title_other, business_member.ownership,
            business_member.is_controlling_manager, business_member.deactivated, business_member.created,
            business_member.modified, consumer.first_name, consumer.middle_name, consumer.last_name,
            consumer.email, consumer.phone, consumer.date_of_birth, consumer.tax_id, consumer.tax_id_type,
			consumer.kyc_status, consumer.legal_address, consumer.mailing_address, consumer.work_address,
			consumer.residency, consumer.citizenship_countries, consumer.occupation, consumer.income_type,
            consumer.activity_type, consumer.is_restricted, consumer.email_id
        FROM
            business_member
        JOIN
            consumer ON business_member.consumer_id = consumer.id
		WHERE
            business_member.deactivated IS NULL AND business_member.business_id = $1`,
		id,
	)

	switch {
	case err == sql.ErrNoRows:
		return rows, nil
	case err != nil:
		log.Print(err)
		return nil, errors.Cause(err)
	}

	return rows, err
}

func (db *memberDataStore) GetById(id shared.BusinessMemberID, businessId shared.BusinessID) (*BusinessMember, error) {
	// Check access
	err := auth.NewAuthService(db.sourceRequest).CheckBusinessAccess(businessId)
	if err != nil {
		return nil, err
	}

	return db.getById(id, businessId)
}

func (db *memberDataStore) getById(id shared.BusinessMemberID, businessId shared.BusinessID) (*BusinessMember, error) {
	m := BusinessMember{}
	err := db.Get(
		&m, `
		SELECT
            business_member.id, business_member.consumer_id, business_member.business_id,
			business_member.title_type, business_member.title_other, business_member.ownership,
			business_member.is_controlling_manager, business_member.deactivated, business_member.created,
            business_member.modified, consumer.first_name, consumer.middle_name, consumer.last_name,
			consumer.email, consumer.phone, consumer.date_of_birth, consumer.tax_id, consumer.tax_id_type,
			consumer.kyc_status, consumer.legal_address, consumer.mailing_address, consumer.work_address,
			consumer.residency, consumer.citizenship_countries, consumer.occupation, consumer.income_type,
            consumer.activity_type, consumer.is_restricted, consumer.email_id
        FROM
            business_member
        JOIN
            consumer ON business_member.consumer_id = consumer.id
        WHERE
            business_member.id = $1 AND business_member.business_id = $2`,
		id,
		businessId,
	)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return &m, err
}

func (db *memberDataStore) Create(m *BusinessMemberCreate) (*BusinessMember, error) {
	// Check access
	err := auth.NewAuthService(db.sourceRequest).CheckBusinessAccess(m.BusinessID)
	if err != nil {
		return nil, err
	}

	// Validate phone no default
	ph, err := libphonenumber.Parse(m.Phone, "")
	if err != nil {
		return nil, err
	}

	m.Phone = libphonenumber.Format(ph, libphonenumber.E164)

	// Validate email
	e, err := mail.ParseAddress(m.Email)
	if err != nil {
		return nil, err
	}

	m.Email = e.Address

	var cid *shared.ConsumerID
	if m.UserID != nil {
		u, err := user.NewUserService(db.sourceRequest).GetById(*m.UserID)
		if err != nil {
			return nil, err
		}

		cid = &u.ConsumerID
	} else {
		cr := user.ConsumerCreate{
			FirstName:            m.FirstName,
			MiddleName:           m.MiddleName,
			LastName:             m.LastName,
			Email:                &m.Email,
			Phone:                &m.Phone,
			DateOfBirth:          m.DateOfBirth,
			TaxID:                m.TaxID,
			TaxIDType:            m.TaxIDType,
			LegalAddress:         m.LegalAddress,
			Residency:            m.Residency,
			CitizenshipCountries: m.CitizenshipCountries,
			Occupation:           &m.Occupation,
			IncomeType:           m.IncomeType,
			ActivityType:         m.ActivityType,
			IsBusinessMember:     true,
		}

		cid, err = user.NewConsumerService(db.sourceRequest).Create(cr)
		if err != nil {
			return nil, err
		}
	}

	mc := struct {
		ConsumerID shared.ConsumerID `json:"consumerId" db:"consumer_id"`
		*BusinessMemberCreate
	}{
		ConsumerID:           *cid,
		BusinessMemberCreate: m,
	}

	q := `
        INSERT INTO business_member(
            consumer_id, business_id, title_type, title_other, ownership, is_controlling_manager
        )
        VALUES(
			:consumer_id, :business_id, :title_type, :title_other, :ownership, :is_controlling_manager
        )
        RETURNING id`

	rows, err := db.NamedQuery(q, &mc)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var id shared.BusinessMemberID
		err = rows.Scan(&id)
		rows.Close()

		if err == nil {
			return db.getById(id, m.BusinessID)
		}

		break
	}

	return nil, err
}

func (db *memberDataStore) Update(id shared.BusinessMemberID, businessId shared.BusinessID, u *BusinessMemberUpdate) (*BusinessMember, error) {
	// Check access
	err := auth.NewAuthService(db.sourceRequest).CheckBusinessAccess(businessId)
	if err != nil {
		return nil, err
	}

	m, err := db.getById(id, businessId)
	if err != nil {
		return nil, err
	}

	// Update consumer object
	cu := user.ConsumerUpdate{
		ID:                   m.ConsumerID,
		FirstName:            u.FirstName,
		MiddleName:           u.MiddleName,
		LastName:             u.LastName,
		Email:                u.Email,
		DateOfBirth:          u.DateOfBirth,
		TaxID:                u.TaxID,
		TaxIDType:            u.TaxIDType,
		LegalAddress:         u.LegalAddress,
		MailingAddress:       u.MailingAddress,
		Residency:            u.Residency,
		CitizenshipCountries: u.CitizenshipCountries,
		Occupation:           u.Occupation,
		IncomeType:           u.IncomeType,
		ActivityType:         u.ActivityType,
	}

	// TODO: Enclose updates in transaction
	_, err = user.NewConsumerService(db.sourceRequest).Update(cu)
	if err != nil {
		return nil, err
	}

	var columns []string

	if u.TitleType != nil {
		columns = append(columns, "title_type = :title_type")
	}

	if u.TitleOther != nil {
		columns = append(columns, "title_other = :title_other")
	}

	if u.Ownership != nil {
		columns = append(columns, "ownership = :ownership")
	}

	// No changes requested - return BusinessMember
	if len(columns) == 0 {
		return nil, user.ConsumerKYCError{
			RawError:   errors.New("member has no new changes"),
			ErrorType:  services.KYCErrorTypeOther,
			ConsumerID: &m.ConsumerID,
		}
	}

	_, err = db.NamedExec(fmt.Sprintf("UPDATE business_member SET %s WHERE id = '%s'", strings.Join(columns, ", "), id), u)
	if err != nil {
		return nil, errors.Cause(err)
	}

	return db.GetById(id, businessId)
}

func (db *memberDataStore) Submit(id shared.BusinessMemberID, businessId shared.BusinessID) (*BusinessMember, error) {
	// Check access
	err := auth.NewAuthService(db.sourceRequest).CheckBusinessAccess(businessId)
	if err != nil {
		return nil, err
	}

	m, err := db.getById(id, businessId)
	if err != nil {
		return nil, err
	}

	_, err = user.NewConsumerService(db.sourceRequest).Submit(m.ConsumerID)
	if err != nil {
		return nil, err
	}

	return db.getById(id, businessId)
}

func (db *memberDataStore) StartVerification(id shared.BusinessMemberID, businessId shared.BusinessID) (*MemberKYCResponse, error) {
	// Check access
	err := auth.NewAuthService(db.sourceRequest).CheckBusinessAccess(businessId)
	if err != nil {
		return nil, err
	}

	m, err := db.getById(id, businessId)
	if err != nil {
		return nil, err
	}

	var resp *user.ConsumerKYCResponse
	switch m.KYCStatus {
	case services.KYCStatusNotStarted:
		return nil, MemberKYCError{
			RawError:  errors.New("member not submitted"),
			ErrorType: services.KYCErrorTypeOther,
			MemberID:  &id,
		}
	case services.KYCStatusReview:
		return nil, MemberKYCError{
			RawError:  errors.New("member already in review"),
			ErrorType: services.KYCErrorTypeOther,
			MemberID:  &id,
		}
	case services.KYCStatusApproved, services.KYCStatusDeclined:
		return nil, MemberKYCError{
			RawError:  errors.New("member already approved or declined"),
			ErrorType: services.KYCErrorTypeOther,
			MemberID:  &id,
		}
	}

	resp, err = user.NewConsumerService(db.sourceRequest).StartVerification(m.ConsumerID, false)
	if err != nil {
		cerr, ok := err.(*user.ConsumerKYCError)
		if ok {
			return nil, MemberKYCError{
				RawError:  cerr.RawError,
				ErrorType: cerr.ErrorType,
				Values:    cerr.Values,
				MemberID:  &id,
			}
		}

		return nil, err
	}

	member, err := db.getById(id, m.BusinessID)
	if err != nil {
		return nil, MemberKYCError{
			RawError:  err,
			ErrorType: services.KYCErrorTypeOther,
			MemberID:  &id,
		}
	}

	return &MemberKYCResponse{
		Status:      resp.Status,
		ReviewItems: resp.ReviewItems,
		Member:      *member,
	}, nil
}

func (db *memberDataStore) GetVerification(id shared.BusinessMemberID, businessId shared.BusinessID) (*MemberKYCResponse, error) {
	// Check access
	err := auth.NewAuthService(db.sourceRequest).CheckBusinessAccess(businessId)
	if err != nil {
		return nil, err
	}

	m, err := db.getById(id, businessId)
	if err != nil {
		return nil, err
	}

	var resp *user.ConsumerKYCResponse
	switch m.KYCStatus {
	case services.KYCStatusNotStarted:
		return nil, MemberKYCError{
			RawError:  errors.New("member not submitted"),
			ErrorType: services.KYCErrorTypeOther,
			MemberID:  &id,
		}
	case services.KYCStatusSubmitted:
		return nil, MemberKYCError{
			RawError:  errors.New("member submitted but not in review"),
			ErrorType: services.KYCErrorTypeOther,
			MemberID:  &id,
		}
	case services.KYCStatusApproved, services.KYCStatusDeclined:
		return nil, MemberKYCError{
			RawError:  errors.New("member already approved or declined"),
			ErrorType: services.KYCErrorTypeOther,
			MemberID:  &id,
		}
	}

	resp, err = user.NewConsumerService(db.sourceRequest).GetVerification(m.ConsumerID)
	if err != nil {
		cerr, ok := err.(*user.ConsumerKYCError)
		if ok {
			return nil, MemberKYCError{
				RawError:  cerr.RawError,
				ErrorType: cerr.ErrorType,
				Values:    cerr.Values,
				MemberID:  &id,
			}
		}

		return nil, err
	}

	member, err := db.getById(id, m.BusinessID)
	if err != nil {
		return nil, MemberKYCError{
			RawError:  err,
			ErrorType: services.KYCErrorTypeOther,
			MemberID:  &id,
		}
	}

	return &MemberKYCResponse{
		Status:      resp.Status,
		ReviewItems: resp.ReviewItems,
		Member:      *member,
	}, nil
}

func (db *memberDataStore) Deactivate(id shared.BusinessMemberID, businessId shared.BusinessID) error {
	// Check access
	err := auth.NewAuthService(db.sourceRequest).CheckBusinessAccess(businessId)
	if err != nil {
		return err
	}

	_, err = db.Exec("DELETE FROM business_member WHERE id = $1 AND business_id = $2", id, businessId)
	return err
}
