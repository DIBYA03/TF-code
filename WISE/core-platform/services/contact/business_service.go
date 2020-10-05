/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

// Package for user contacts
package contact

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
	"github.com/wiseco/core-platform/notification/activity"
	"github.com/wiseco/core-platform/services"
	"github.com/wiseco/core-platform/services/data"
	"github.com/wiseco/core-platform/services/email"
	"github.com/wiseco/core-platform/shared"
)

type contactDatastore struct {
	sourceReq services.SourceRequest
	*sqlx.DB
}

type ContactService interface {
	// Read
	GetById(string, shared.BusinessID) (*Contact, error)
	GetByIDInternal(string) (*Contact, error)
	GetByPhoneEmailInternal(*string, *string, shared.BusinessID) (*Contact, error)
	List(offset int, limit int, businessID shared.BusinessID) ([]Contact, error)

	// Create/Modify
	Create(*ContactCreate) (*Contact, error)
	Update(*ContactUpdate) (*Contact, error)

	// Deactivate contact
	Deactivate(string, shared.BusinessID) error
}

func NewContactService(r services.SourceRequest) ContactService {
	return &contactDatastore{r, data.DBWrite}
}

func NewContactServiceWithout() ContactService {
	return &contactDatastore{services.SourceRequest{}, data.DBWrite}
}

func (db *contactDatastore) List(offset int, limit int, businessID shared.BusinessID) ([]Contact, error) {
	// Check access
	err := auth.NewAuthService(db.sourceReq).CheckBusinessAccess(businessID)
	if err != nil {
		return nil, err
	}

	rows := []Contact{}
	err = db.Select(&rows, "SELECT * FROM business_contact WHERE deactivated IS NULL AND business_id = $1", businessID)
	switch {
	case err == sql.ErrNoRows:
		return rows, nil
	case err != nil:
		return nil, errors.Cause(err)
	}

	return rows, err
}

func (db *contactDatastore) GetById(id string, businessID shared.BusinessID) (*Contact, error) {
	// Check access
	err := auth.NewAuthService(db.sourceReq).CheckBusinessAccess(businessID)
	if err != nil {
		return nil, err
	}

	contact := &Contact{}
	err = db.Get(contact, "SELECT * FROM business_contact WHERE id = $1 AND business_id = $2", id, businessID)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return contact, err
}

func (db *contactDatastore) GetByIDInternal(id string) (*Contact, error) {
	contact := &Contact{}
	err := db.Get(contact, "SELECT * FROM business_contact WHERE id = $1", id)
	if err != nil {
		return nil, err
	}

	return contact, err
}

func (db *contactDatastore) GetByPhoneEmailInternal(phone *string, email *string, bID shared.BusinessID) (*Contact, error) {
	contact := &Contact{}

	businessIDClause := "business_id = '" + string(bID) + "'"
	var whereClause string
	if phone != nil && len(*phone) > 0 {
		whereClause = "phone_number = '" + *phone + "'"
	}

	if email != nil && len(*email) > 0 {
		if len(whereClause) > 0 {
			whereClause = whereClause + "OR email = '" + *email + "'"
		} else {
			whereClause = "email = '" + *email + "'"
		}
	}

	if len(whereClause) == 0 {
		return nil, errors.New("Both email and phone can't be empty")
	}

	err := db.Get(contact, "SELECT * FROM business_contact WHERE "+businessIDClause+" AND ("+whereClause+")")
	if err != nil {
		log.Println("Error in fetching contact", err, "SELECT * FROM business_contact WHERE "+whereClause)
		return nil, err
	}

	return contact, err
}

func (db *contactDatastore) Create(contactCreate *ContactCreate) (*Contact, error) {
	// Check access
	err := auth.NewAuthService(db.sourceReq).CheckBusinessAccess(contactCreate.BusinessID)
	if err != nil {
		return nil, err
	}

	// Default/mandatory fields
	columns := []string{
		"user_id", "business_id", "contact_type", "engagement", "job_title", "business_name", "phone_number", "email", "email_id",
	}
	// Default/mandatory values
	values := []string{
		":user_id", ":business_id", ":contact_type", ":engagement", ":job_title", ":business_name", ":phone_number", ":email", ":email_id",
	}

	if contactCreate.Category != nil {
		columns = append(columns, "contact_category")
		values = append(values, ":contact_category")
	}

	if contactCreate.Type == "" {
		return nil, errors.New("Contact type is required")
	}

	switch contactCreate.Type {
	case ContactTypePerson:
		if contactCreate.FirstName == nil {
			return nil, errors.New("First name is required")
		}

		if contactCreate.LastName == nil {
			return nil, errors.New("Last name is required")
		}

		columns = append(columns, "first_name")
		values = append(values, ":first_name")

		columns = append(columns, "last_name")
		values = append(values, ":last_name")
	case ContactTypeBusiness:
		if contactCreate.BusinessName == nil {
			return nil, errors.New("Business name is required")
		}
		break
	default:
		return nil, errors.New("Contact type is invalid")
	}

	// Validate phone no default
	if contactCreate.PhoneNumber != "" {
		ph, err := libphonenumber.Parse(contactCreate.PhoneNumber, "")
		if err != nil {
			return nil, errors.New("Invalid phone number")
		}

		contactCreate.PhoneNumber = libphonenumber.Format(ph, libphonenumber.E164)
	}

	// Validate email
	if contactCreate.Email != "" {
		pe, err := mail.ParseAddress(contactCreate.Email)
		if err != nil {
			return nil, errors.New("Invalid email address")
		}
		contactCreate.Email = pe.Address
	}

	ec := email.EmailCreate{
		EmailAddress: email.EmailAddress(contactCreate.Email),
		Status:       email.StatusActive,
		Type:         email.TypeContact,
	}

	e, err := email.NewEmailService(db.sourceReq).Create(&ec)

	if err != nil {
		return nil, err
	}

	contactCreate.EmailID = e.ID

	sql := fmt.Sprintf("INSERT INTO business_contact(%s) VALUES(%s) RETURNING *", strings.Join(columns, ", "), strings.Join(values, ", "))

	stmt, err := db.PrepareNamed(sql)
	if err != nil {
		return nil, err
	}

	contact := &Contact{}

	err = stmt.Get(contact, &contactCreate)
	if err != nil {
		return nil, err
	}

	db.onContactCreate(*contact)

	return contact, nil

}

func (db *contactDatastore) onContactCreate(c Contact) error {

	contact := activity.Contact{
		EntityID:  string(c.BusinessID),
		UserID:    c.UserID,
		ContactID: c.ID,
	}

	contact.Name = c.Name()
	err := activity.NewContactCreator().Create(contact)
	if err != nil {
		log.Println(err)
	}

	return err
}

func (db *contactDatastore) Update(contactUpdate *ContactUpdate) (*Contact, error) {
	// Check access
	err := auth.NewAuthService(db.sourceReq).CheckBusinessAccess(contactUpdate.BusinessID)
	if err != nil {
		return nil, err
	}

	b, err := db.GetById(contactUpdate.ID, contactUpdate.BusinessID)
	if err != nil {
		return nil, err
	}

	var columns []string

	if contactUpdate.Category != nil {
		columns = append(columns, "contact_category = :contact_category")
	}

	if contactUpdate.Type != nil {
		columns = append(columns, "contact_type = :contact_type")
	}

	if contactUpdate.Engagement != nil {
		columns = append(columns, "engagement = :engagement")
	}

	if contactUpdate.JobTitle != nil {
		columns = append(columns, "job_title = :job_title")
	}

	if contactUpdate.BusinessName != "" {
		columns = append(columns, "business_name = :business_name")
	}

	if contactUpdate.FirstName != nil {
		columns = append(columns, "first_name = :first_name")
	}

	if contactUpdate.LastName != nil {
		columns = append(columns, "last_name = :last_name")
	}

	if contactUpdate.PhoneNumber != nil {
		// Validate phone no default
		ph, err := libphonenumber.Parse(*contactUpdate.PhoneNumber, "")
		if err != nil {
			return nil, errors.New("Invalid phone number")
		}

		phone := libphonenumber.Format(ph, libphonenumber.E164)
		contactUpdate.PhoneNumber = &phone

		columns = append(columns, "phone_number = :phone_number")
	}

	if contactUpdate.Email != nil {
		// Validate email
		pe, err := mail.ParseAddress(*contactUpdate.Email)
		if err != nil {
			return nil, errors.New("Invalid email address")
		}

		contactUpdate.Email = &pe.Address

		columns = append(columns, "email = :email")

		esrvc := email.NewEmailService(db.sourceReq)
		if b.EmailID != shared.EmailID("") {
			err = esrvc.Deactivate(b.EmailID)

			if err != nil {
				return nil, err
			}
		}

		ec := email.EmailCreate{
			EmailAddress: email.EmailAddress(pe.Address),
			Status:       email.StatusActive,
			Type:         email.TypeConsumer,
		}

		e, err := esrvc.Create(&ec)

		if err != nil {
			return nil, err
		}

		contactUpdate.EmailID = e.ID

		columns = append(columns, "email_id = :email_id")
	}

	// No changes requested - return contact
	if len(columns) == 0 {
		return b, nil
	}

	_, err = db.NamedExec(
		fmt.Sprintf(
			"UPDATE business_contact SET %s WHERE id = '%s'",
			strings.Join(columns, ", "),
			contactUpdate.ID,
		), contactUpdate,
	)

	if err != nil {
		return nil, errors.Cause(err)
	}

	c, err := db.GetById(contactUpdate.ID, contactUpdate.BusinessID)

	db.onContactUpdate(*c)

	return c, err
}

func (db *contactDatastore) onContactUpdate(c Contact) error {

	contact := activity.Contact{
		EntityID:  string(c.BusinessID),
		UserID:    c.UserID,
		ContactID: c.ID,
	}

	contact.Name = c.Name()
	err := activity.NewContactCreator().Update(contact)

	if err != nil {
		log.Println(err)
	}

	return err
}

func (db *contactDatastore) Deactivate(ID string, businessID shared.BusinessID) error {
	// Check access
	err := auth.NewAuthService(db.sourceReq).CheckBusinessAccess(businessID)
	if err != nil {
		return err
	}

	_, err = db.Exec("UPDATE business_contact SET deactivated = CURRENT_TIMESTAMP WHERE id = $1 AND business_id = $2", ID, businessID)

	// Add to activity stream
	c, err := db.GetById(ID, businessID)
	if err == nil {
		db.onContactDelete(*c)
	}

	return err
}

func (db *contactDatastore) onContactDelete(c Contact) error {

	contact := activity.Contact{
		EntityID:  string(c.BusinessID),
		UserID:    c.UserID,
		ContactID: c.ID,
	}

	contact.Name = c.Name()
	err := activity.NewContactCreator().Delete(contact)
	if err != nil {
		log.Println(err)
	}

	return err
}
