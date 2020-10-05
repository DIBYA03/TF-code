/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

// Package cspuser for all csp user related services
package cspuser

import (
	"fmt"
	"log"
	"net/mail"

	"github.com/jmoiron/sqlx"
	"github.com/wiseco/core-platform/services/csp/auth"
	"github.com/wiseco/core-platform/services/csp/data"
	"github.com/wiseco/core-platform/services/csp/services"
	"github.com/wiseco/go-lib/id"
)

type cspUserDatastore struct {
	sourceReq services.SourceRequest
	*sqlx.DB
}

type CspUserService interface {

	// Fetch operations
	GetById(string) (*CSPUser, error)
	GetByIdInternal(string) (*CSPUser, error)
	GetIdByCognitoID(string) (id.CspAgentID, error)

	ByCognitoID(string) (string, error)

	// Create user
	Create(CSPUser) (*CSPUser, error)
	GetUserByEmail(string) (*CSPUser, error)
	Update(CSPUser) (*CSPUser, error)

	// Deactivate user by id
	Deactivate(string) error
}

func NewUserService(r services.SourceRequest) CspUserService {
	return &cspUserDatastore{r, data.DBWrite}
}

func (db *cspUserDatastore) GetById(id string) (*CSPUser, error) {
	// Check access
	err := auth.NewService(db.sourceReq).CheckUserAccess(id)
	if err != nil {
		return nil, err
	}

	return db.getById(id)
}

func (db *cspUserDatastore) getById(id string) (*CSPUser, error) {
	u := CSPUser{}
	err := db.Get(&u, `SELECT * FROM csp_user WHERE id = $1`, id)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &u, nil
}

func (db *cspUserDatastore) GetUserByEmail(email string) (*CSPUser, error) {
	var u CSPUser
	err := db.Get(&u, "SELECT * FROM csp_user WHERE email = $1", email)
	return &u, err
}

func (db *cspUserDatastore) GetByIdInternal(id string) (*CSPUser, error) {
	return db.getById(id)
}

func (db *cspUserDatastore) Create(u CSPUser) (*CSPUser, error) {
	// Validate email
	if u.Email != nil {
		e, err := mail.ParseAddress(*u.Email)
		if err != nil {
			return nil, err
		}

		u.Email = &e.Address
	}

	sql := `
		INSERT INTO csp_user(
			cognito_id, first_name, middle_name, last_name, email, email_verified,
				phone, phone_verified, picture
		)
		VALUES(
			:cognito_id, :first_name, :middle_name, :last_name, :email, :email_verified,
				:phone, :phone_verified, :picture
		)
		RETURNING id`

	_, err := db.NamedExec(sql, &u)
	if err != nil {
		return &u, err
	}

	return &u, nil
}

func (db *cspUserDatastore) Update(u CSPUser) (*CSPUser, error) {
	// Validate email
	if u.Email != nil {
		e, err := mail.ParseAddress(*u.Email)
		if err != nil {
			return nil, err
		}

		u.Email = &e.Address
	}

	sql := `
		UPDATE csp_user
			SET
				first_name  = :first_name,
				middle_name = :middle_name,
				last_name   = :last_name,
				phone       = :phone,
				picture     = :picture
			WHERE
				id = :id`

	_, err := db.NamedExec(sql, &u)
	if err != nil {
		return &u, err
	}

	return &u, nil
}

func (db *cspUserDatastore) Deactivate(id string) error {
	// Check access
	err := auth.NewService(db.sourceReq).CheckUserAccess(id)
	if err != nil {
		return err
	}

	_, err = db.Exec(fmt.Sprintf("UPDATE csp_user SET active = false WHERE id = '%s'", id))
	return err
}

func (db *cspUserDatastore) ByCognitoID(cognitoID string) (string, error) {
	var userID string
	err := db.Get(&userID, `SELECT id FROM csp_user WHERE cognito_id = $1`, cognitoID)
	return userID, err
}

func (db *cspUserDatastore) GetIdByCognitoID(cognitoID string) (id.CspAgentID, error) {
	var userID id.CspAgentID
	err := db.Get(&userID, `SELECT id FROM csp_user WHERE cognito_id = $1`, cognitoID)
	return userID, err
}
