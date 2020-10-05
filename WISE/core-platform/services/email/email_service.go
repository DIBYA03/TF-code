/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

// Package for email
package email

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"

	"github.com/wiseco/core-platform/services"
	"github.com/wiseco/core-platform/services/data"
	"github.com/wiseco/core-platform/shared"
)

type emailDatastore struct {
	sourceReq services.SourceRequest
	*sqlx.DB
}

type EmailService interface {
	// Read
	GetByID(shared.EmailID) (*Email, error)
	GetByEmailAddress(string) (*Email, error)

	// Create
	Create(*EmailCreate) (*Email, error)

	// Deactivate email
	Deactivate(shared.EmailID) error

	// Is the email passed in available for the type?
	IsAvailable(EmailAddress, Type) (bool, error)

	//DEBUG
	DEBUGDeleteByID(shared.EmailID) error
}

func NewEmailService(r services.SourceRequest) EmailService {
	return &emailDatastore{r, data.DBWrite}
}

func (db *emailDatastore) GetByID(id shared.EmailID) (*Email, error) {
	e := new(Email)

	err := db.Get(e, "SELECT * FROM email WHERE id = $1", id)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return e, err
}

func (db *emailDatastore) GetByEmailAddress(emailAddress string) (*Email, error) {
	e := new(Email)

	err := db.Get(e, "SELECT * FROM email WHERE email_address = $1", emailAddress)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return e, err
}

func (db *emailDatastore) Create(emailCreate *EmailCreate) (*Email, error) {

	// Default/mandatory fields
	columns := []string{
		"email_address", "email_status", "email_type",
	}
	// Default/mandatory values
	values := []string{
		":email_address", ":email_status", ":email_type",
	}

	sql := fmt.Sprintf("INSERT INTO email(%s) VALUES(%s) RETURNING *", strings.Join(columns, ", "), strings.Join(values, ", "))

	stmt, err := db.PrepareNamed(sql)
	if err != nil {
		return nil, err
	}

	email := new(Email)

	err = stmt.Get(email, &emailCreate)
	if err != nil {
		return nil, err
	}

	return email, nil
}

func (db *emailDatastore) Deactivate(ID shared.EmailID) error {
	_, err := db.Exec("UPDATE email SET email_status = $1 WHERE id = $2", StatusInactive, ID)

	return err
}

func (db *emailDatastore) IsAvailable(ea EmailAddress, et Type) (bool, error) {
	var exists int

	err := db.QueryRow("SELECT COUNT(*) FROM email WHERE email_address = $1 AND email_status = $2 AND email_type = $3", ea, StatusActive, et).Scan(&exists)

	return exists < 1, err
}

//ONLY for debug
func (db *emailDatastore) DEBUGDeleteByID(emailID shared.EmailID) error {
	_, err := db.Exec("DELETE FROM email where id = $1", emailID)

	return err
}
