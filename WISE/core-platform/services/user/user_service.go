/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

// Package for all user related services
package user

import (
	"fmt"
	"log"
	"net/mail"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/ttacon/libphonenumber"
	"github.com/wiseco/core-platform/auth"
	partnerbank "github.com/wiseco/core-platform/partner/bank"
	_ "github.com/wiseco/core-platform/partner/bank/bbva"
	"github.com/wiseco/core-platform/partner/service/segment"
	"github.com/wiseco/core-platform/services"
	"github.com/wiseco/core-platform/services/data"
	"github.com/wiseco/core-platform/shared"
)

type userDatastore struct {
	sourceReq services.SourceRequest
	*sqlx.DB
}

type UserService interface {

	// Fetch operations
	GetById(shared.UserID) (*User, error)
	GetByIdInternal(shared.UserID) (*User, error)
	GetUserIDWithPhone(phone string) (*shared.UserID, error)
	GetUserIDWithIdentity(shared.IdentityID) (*shared.UserID, error)
	GetUserWithIdentity(shared.IdentityID) (*User, error)
	GetUserIDWithConsumer(shared.ConsumerID) (shared.UserID, error)

	// Create user
	Create(UserCreate) (*shared.UserID, error)

	// Create user from auth data
	CreateFromAuth(UserAuthCreate) (*shared.UserID, error)

	// Update user
	Update(UserUpdate) (*User, error)

	// Update user subscription
	UpdateSubscription(shared.UserID, services.SubscriptionStatus) (*User, error)

	// Deactivate user by id
	Deactivate(shared.UserID) error

	// Dev env only
	DeleteById(shared.UserID) error

	// User verification (KYC)
	Submit(shared.UserID) (*User, error)
	StartVerification(shared.UserID) (*UserKYCResponse, error)
	GetVerification(shared.UserID) (*UserKYCResponse, error)

	//User Notification Settings
	UpdateNotification(update *UserNotificationUpdate) (UserNotification, error)

	// Update user partner code
	UpdatePartnerCode(UserPartnerUpdate) (*User, error)
	GetPartnerCodeCount(shared.PartnerID) (int, error)
}

func NewUserService(r services.SourceRequest) UserService {
	return &userDatastore{r, data.DBWrite}
}

// NewUserServiceWithout returns an user service without a source request
func NewUserServiceWithout() UserService {
	return &userDatastore{services.SourceRequest{}, data.DBWrite}
}

func (db *userDatastore) GetById(id shared.UserID) (*User, error) {

	// Check access
	err := auth.NewAuthService(db.sourceReq).CheckUserAccess(id)
	if err != nil {
		return nil, err
	}

	u, err := db.getById(id)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (db *userDatastore) getById(id shared.UserID) (*User, error) {

	u := User{}
	err := db.Get(
		&u, `
		SELECT
			wise_user.id, wise_user.consumer_id, wise_user.identity_id, wise_user.partner_id,
			wise_user.email, wise_user.email_verified, wise_user.phone, wise_user.phone_verified,
			wise_user.notification, wise_user.deactivated, wise_user.created, wise_user.modified, 
			wise_user.subscription_status, consumer.first_name, consumer.middle_name, consumer.last_name, 
			consumer.date_of_birth, consumer.tax_id, consumer.tax_id_type, consumer.kyc_status, 
			consumer.legal_address, consumer.mailing_address, consumer.work_address, consumer.residency,
			consumer.citizenship_countries, consumer.occupation, consumer.income_type,
			consumer.activity_type, consumer.is_restricted, consumer.email_id
		FROM
			wise_user
		JOIN
			consumer ON wise_user.consumer_id = consumer.id
		WHERE
			wise_user.id = $1`,
		id,
	)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &u, nil
}

func (db *userDatastore) GetUserIDWithPhone(phone string) (*shared.UserID, error) {
	var id shared.UserID
	err := db.Get(&id, "SELECT id FROM wise_user WHERE phone = $1", phone)
	return &id, err
}

func (db *userDatastore) GetUserIDWithIdentity(id shared.IdentityID) (*shared.UserID, error) {
	var uid shared.UserID
	err := db.Get(&uid, "SELECT id FROM wise_user WHERE identity_id = $1", id)
	return &uid, err
}

func (db *userDatastore) GetUserWithIdentity(id shared.IdentityID) (*User, error) {
	resp := &User{}
	userID, err := db.GetUserIDWithIdentity(id)
	if err != nil {
		return resp, err
	}

	return db.getById(*userID)
}

func (db *userDatastore) GetUserIDWithConsumer(id shared.ConsumerID) (shared.UserID, error) {
	var uid shared.UserID
	err := db.Get(&uid, "SELECT id FROM wise_user WHERE consumer_id = $1", id)
	return uid, err
}

func (db *userDatastore) GetByIdInternal(id shared.UserID) (*User, error) {
	return db.getById(id)
}

func (db *userDatastore) Create(u UserCreate) (*shared.UserID, error) {
	// Validate phone no default
	ph, err := libphonenumber.Parse(u.Phone, "")
	if err != nil {
		return nil, err
	}

	u.Phone = libphonenumber.Format(ph, libphonenumber.E164)

	// Validate email
	if u.Email != nil {
		e, err := mail.ParseAddress(*u.Email)
		if err != nil {
			return nil, err
		}

		u.Email = &e.Address
	}

	cr := ConsumerCreate{
		FirstName:            u.FirstName,
		MiddleName:           u.MiddleName,
		LastName:             u.LastName,
		Email:                u.Email,
		Phone:                &u.Phone,
		DateOfBirth:          u.DateOfBirth,
		TaxID:                u.TaxID,
		TaxIDType:            u.TaxIDType,
		LegalAddress:         u.LegalAddress,
		MailingAddress:       u.MailingAddress,
		WorkAddress:          u.WorkAddress,
		Residency:            u.Residency,
		CitizenshipCountries: u.CitizenshipCountries,
		Occupation:           u.Occupation,
		IncomeType:           u.IncomeType,
		ActivityType:         u.ActivityType,
	}

	cid, err := NewConsumerService(db.sourceReq).Create(cr)
	if err != nil {
		return nil, err
	}

	uc := struct {
		ConsumerID shared.ConsumerID `json:"consumerId" db:"consumer_id"`
		UserCreate
	}{
		ConsumerID: *cid,
		UserCreate: u,
	}

	sql := `
		INSERT INTO wise_user(
			consumer_id, identity_id, email, email_verified, phone, phone_verified
		)
		VALUES(
			:consumer_id, :identity_id, :email, :email_verified, :phone, :phone_verified
		)
		RETURNING id`

	stmt, err := db.PrepareNamed(sql)
	if err != nil {
		return nil, err
	}

	var id shared.UserID
	err = stmt.Get(&id, &uc)
	if err != nil {
		return nil, err
	}

	return &id, nil
}

//CreateFromAuth creates a minimal user from cognito
func (db *userDatastore) CreateFromAuth(u UserAuthCreate) (*shared.UserID, error) {
	// Validate phone no default
	ph, err := libphonenumber.Parse(u.Phone, "")
	if err != nil {
		return nil, err
	}

	u.Phone = libphonenumber.Format(ph, libphonenumber.E164)

	// Validate email
	if u.Email != nil {
		e, err := mail.ParseAddress(*u.Email)
		if err != nil {
			return nil, err
		}

		u.Email = &e.Address
	}

	cr := ConsumerAuthCreate{
		Email: u.Email,
		Phone: &u.Phone,
	}

	cid, err := NewConsumerService(db.sourceReq).CreateFromAuth(cr)
	if err != nil {
		return nil, err
	}

	uc := struct {
		ConsumerID shared.ConsumerID `json:"consumerId" db:"consumer_id"`
		UserAuthCreate
	}{
		ConsumerID:     *cid,
		UserAuthCreate: u,
	}

	println("Analytics debug: before insert")

	// Insert statement
	sql := `
        INSERT INTO wise_user(
            consumer_id, identity_id, email, email_verified, phone, phone_verified
        )
        VALUES(
			:consumer_id, :identity_id, :email, :email_verified, :phone, :phone_verified
        )
        RETURNING id`

	println("Analytics debug: after insert")

	// Execute
	stmt, err := db.PrepareNamed(sql)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	println("Analytics debug: after insert 1")

	// Return id
	var id shared.UserID
	err = stmt.Get(&id, &uc)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	println("Analytics debug: before insert 2")

	// Send to segment
	segment.NewSegmentService().PushToAnalyticsQueue(id, segment.CategoryConsumer, segment.ActionCreate, u)

	println("Analytics debug: before insert 3")

	return &id, nil
}

func (db *userDatastore) Update(u UserUpdate) (*User, error) {
	// Check access
	err := auth.NewAuthService(db.sourceReq).CheckUserAccess(u.ID)
	if err != nil {
		return nil, err
	}

	usr, err := db.getById(u.ID)
	if err != nil {
		return nil, err
	}

	// Update consumer object
	cu := ConsumerUpdate{
		ID:                   usr.ConsumerID,
		FirstName:            u.FirstName,
		MiddleName:           u.MiddleName,
		LastName:             u.LastName,
		Email:                u.Email,
		DateOfBirth:          u.DateOfBirth,
		TaxID:                u.TaxID,
		TaxIDType:            u.TaxIDType,
		LegalAddress:         u.LegalAddress,
		MailingAddress:       u.MailingAddress,
		WorkAddress:          u.WorkAddress,
		Residency:            u.Residency,
		CitizenshipCountries: u.CitizenshipCountries,
		Occupation:           u.Occupation,
		IncomeType:           u.IncomeType,
		ActivityType:         u.ActivityType,
	}

	// TODO: Enclose updates in transaction
	_, err = NewConsumerService(db.sourceReq).Update(cu)
	if err != nil {
		return nil, err
	}

	var columns []string

	if u.Email != nil {
		// Validate email
		e, err := mail.ParseAddress(*u.Email)
		if err != nil {
			return nil, err
		}

		u.Email = &e.Address

		columns = append(columns, "email = :email")
	}

	// No changes requested - return user
	if len(columns) == 0 {
		return db.getById(u.ID)
	}

	_, err = db.NamedExec(fmt.Sprintf("UPDATE wise_user SET %s WHERE id = '%s'", strings.Join(columns, ", "), u.ID), u)

	if err != nil {
		log.Println(err)
		return nil, errors.Cause(err)
	}

	user, err := db.getById(u.ID)

	// Send to segment
	segment.NewSegmentService().PushToAnalyticsQueue(u.ID, segment.CategoryConsumer, segment.ActionUpdate, user)

	return user, err
}

func (db *userDatastore) updateVerification(u UserVerificationUpdate) (*User, error) {
	_, err := db.NamedExec(
		fmt.Sprintf("UPDATE wise_user SET kyc_status = :kyc_status WHERE id = '%s'", u.ID), u,
	)
	if err != nil {
		log.Println(err)
		return nil, errors.Cause(err)
	}

	// Send to segment
	segment.NewSegmentService().PushToAnalyticsQueue(u.ID, segment.CategoryConsumer, segment.ActionKYC, u)

	return db.getById(u.ID)
}

func (db *userDatastore) UpdateSubscription(userID shared.UserID, subscriptionStatus services.SubscriptionStatus) (*User, error) {
	_, err := db.Exec("UPDATE wise_user SET subscription_status = $1 WHERE id = $2", subscriptionStatus, userID)
	if err != nil {
		log.Println(err)
		return nil, errors.Cause(err)
	}

	return db.getById(userID)
}

func (db *userDatastore) Deactivate(id shared.UserID) error {
	// Check access
	err := auth.NewAuthService(db.sourceReq).CheckUserAccess(id)
	if err != nil {
		return err
	}

	_, err = db.Exec(fmt.Sprintf("UPDATE wise_user SET deactivated = CURRENT_TIMESTAMP WHERE id = '%s'", id))
	return err
}

func (db *userDatastore) DeleteById(id shared.UserID) error {
	user, err := db.getById(id)
	if err != nil {
		return err
	}

	tx := db.MustBegin()

	device := `DELETE FROM user_device WHERE user_id = $1`
	tx.MustExec(device, id)
	la := `DELETE FROM business_linked_bank_account WHERE user_id = $1 AND registered_bank_name = $2`
	tx.MustExec(la, id, partnerbank.ProviderNameBBVA)
	lc := `DELETE FROM business_linked_card WHERE contact_id IN (select id FROM business_contact WHERE user_id = $1)`
	tx.MustExec(lc, id)
	mt := `DELETE FROM business_money_transfer WHERE created_user_id = $1`
	tx.MustExec(mt, id)
	mtr := `DELETE FROM money_transfer_request WHERE created_user_id = $1`
	tx.MustExec(mtr, id)
	bmt := `DELETE FROM business_money_transfer WHERE money_request_id in (select id FROM business_money_request WHERE created_user_id = $1)`
	tx.MustExec(bmt, id)
	mrp := `DELETE FROM business_money_request_payment WHERE request_id in (select id FROM business_money_request WHERE created_user_id = $1)`
	tx.MustExec(mrp, id)
	br := `DELETE FROM business_receipt WHERE created_user_id = $1`
	tx.MustExec(br, id)
	bi := `DELETE FROM business_invoice WHERE created_user_id = $1`
	tx.MustExec(bi, id)
	mr := `DELETE FROM business_money_request WHERE created_user_id = $1`
	tx.MustExec(mr, id)
	c := `DELETE FROM business_contact WHERE user_id = $1`
	tx.MustExec(c, id)
	ua := `DELETE FROM user_activity WHERE entity_id = $1`
	tx.MustExec(ua, id)
	ba := `DELETE FROM business_activity WHERE entity_id IN (select business_id FROM business_member WHERE consumer_id = $1)`
	tx.MustExec(ba, user.ConsumerID)
	sr := `DELETE FROM signature_request WHERE business_id IN (SELECT id FROM business WHERE owner_id = $1)`
	tx.MustExec(sr, id)
	d := `DELETE FROM business_document WHERE business_id IN (select business_id FROM business_document WHERE created_user_id = $1)`
	tx.MustExec(d, id)
	bc := `DELETE FROM business_bank_card WHERE cardholder_id = $1`
	tx.MustExec(bc, id)
	bb := `DELETE FROM business_bank_account  WHERE account_holder_id = $1`
	tx.MustExec(bb, id)
	m := `DELETE FROM business_member WHERE business_id IN (select business_id FROM business_member WHERE consumer_id = $1)`
	tx.MustExec(m, user.ConsumerID)
	pos := `DELETE FROM card_reader WHERE business_id IN (SELECT id FROM business WHERE owner_id = $1)`
	tx.MustExec(pos, id)
	b := `DELETE FROM business WHERE owner_id = $1`
	tx.MustExec(b, id)
	udoc := `DELETE FROM user_document WHERE user_id = $1`
	tx.MustExec(udoc, id)
	u := `DELETE FROM wise_user WHERE id = $1`
	tx.MustExec(u, id)
	cd := `DELETE FROM consumer_document WHERE consumer_id = $1`
	tx.MustExec(cd, user.ConsumerID)
	co := `DELETE FROM consumer WHERE id = $1`
	tx.MustExec(co, user.ConsumerID)

	if user.EmailID != shared.EmailID("") {
		e := `DELETE FROM email WHERE id = $1`
		tx.MustExec(e, user.EmailID)
	}

	err = tx.Commit()

	bank, err := partnerbank.GetConsumerBank(partnerbank.ProviderNameBBVA)
	if err == nil {
		return bank.ConsumerEntityService(db.sourceReq.PartnerBankRequest()).Delete(partnerbank.ConsumerID(user.ConsumerID))
	}

	return err
}

func (db *userDatastore) Submit(id shared.UserID) (*User, error) {
	// Check access
	err := auth.NewAuthService(db.sourceReq).CheckUserAccess(id)
	if err != nil {
		return nil, err
	}

	user, err := db.getById(id)
	if err != nil {
		return nil, err
	}

	_, err = NewConsumerService(db.sourceReq).Submit(user.ConsumerID)
	if err != nil {
		return nil, err
	}

	return db.getById(id)
}

func (db *userDatastore) StartVerification(id shared.UserID) (*UserKYCResponse, error) {
	// Check access
	err := auth.NewAuthService(db.sourceReq).CheckUserAccess(id)
	if err != nil {
		return nil, err
	}

	user, err := db.getById(id)
	if err != nil {
		return nil, ConsumerKYCError{
			RawError:  err,
			ErrorType: services.KYCErrorTypeOther,
		}
	}

	resp, err := NewConsumerService(db.sourceReq).StartVerification(user.ConsumerID, false)
	if err != nil {
		cerr, ok := err.(*ConsumerKYCError)
		if ok {
			return nil, UserKYCError{
				RawError:  cerr.RawError,
				ErrorType: cerr.ErrorType,
				Values:    cerr.Values,
				UserID:    &id,
			}
		}

		return nil, err
	}

	user, err = db.getById(id)
	if err != nil {
		return nil, UserKYCError{
			RawError:  err,
			ErrorType: services.KYCErrorTypeOther,
			UserID:    &id,
		}
	}

	return &UserKYCResponse{
		Status:      resp.Status,
		ReviewItems: resp.ReviewItems,
		User:        *user,
	}, nil
}

func (db *userDatastore) GetVerification(id shared.UserID) (*UserKYCResponse, error) {
	// Check access
	err := auth.NewAuthService(db.sourceReq).CheckUserAccess(id)
	if err != nil {
		return nil, err
	}

	user, err := db.getById(id)
	if err != nil {
		return nil, ConsumerKYCError{
			RawError:  err,
			ErrorType: services.KYCErrorTypeOther,
		}
	}

	resp, err := NewConsumerService(db.sourceReq).GetVerification(user.ConsumerID)
	if err != nil {
		cerr, ok := err.(*ConsumerKYCError)
		if ok {
			return nil, UserKYCError{
				RawError:  cerr.RawError,
				ErrorType: cerr.ErrorType,
				Values:    cerr.Values,
				UserID:    &id,
			}
		}

		return nil, err
	}

	user, err = db.getById(id)
	if err != nil {
		return nil, UserKYCError{
			RawError:  err,
			ErrorType: services.KYCErrorTypeOther,
			UserID:    &id,
		}
	}

	var items services.StringArray
	for _, i := range resp.ReviewItems {
		items = append(items, i)
	}

	name := user.FirstName + " " + user.LastName
	sendConsumerVerification(string(user.ConsumerID), name, resp.Status, items, ActionUpdate)
	return &UserKYCResponse{
		Status:      resp.Status,
		ReviewItems: resp.ReviewItems,
		User:        *user,
	}, nil
}

func (db *userDatastore) UpdateNotification(update *UserNotificationUpdate) (UserNotification, error) {
	var notification UserNotification

	// Check access
	err := auth.NewAuthService(db.sourceReq).CheckUserAccess(db.sourceReq.UserID)
	if err != nil {
		return notification, err
	}

	var current UserNotification

	err = db.Get(&current, "SELECT notification FROM wise_user WHERE id = $1 ", db.sourceReq.UserID)
	if err != nil {
		return notification, err
	}

	if update.Transactions == nil {
		update.Transactions = current.Transactions
	}

	if update.Transfers == nil {
		update.Transfers = current.Transfers
	}
	if update.Contacts == nil {
		update.Contacts = current.Contacts
	}

	err = db.Get(&notification, "UPDATE wise_user SET notification = $2 WHERE id = $1 RETURNING notification", db.sourceReq.UserID, update)
	return notification, err
}

func (db *userDatastore) UpdatePartnerCode(u UserPartnerUpdate) (*User, error) {
	user, err := db.GetById(u.ID)
	if err != nil {
		return nil, err
	}

	if user.PartnerID != nil {
		return nil, errors.New("Partner code already set")
	}

	_, err = db.NamedExec(
		fmt.Sprintf("UPDATE wise_user SET partner_id = :partner_id WHERE id = '%s'", u.ID), u,
	)
	if err != nil {
		log.Println(err)
		return nil, errors.Cause(err)
	}

	return db.getById(u.ID)
}

func (db *userDatastore) GetPartnerCodeCount(partnerID shared.PartnerID) (int, error) {
	var count int

	err := db.QueryRow("SELECT COUNT(*) FROM wise_user WHERE deactivated IS NULL AND partner_id = $1", partnerID).Scan(&count)

	if err != nil {
		log.Println(err)
		return 0, errors.Cause(err)
	}

	return count, err
}
