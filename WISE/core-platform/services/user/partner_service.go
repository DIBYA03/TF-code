/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

// Package for all user related services
package user

import (
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/wiseco/core-platform/auth"
	_ "github.com/wiseco/core-platform/partner/bank/bbva"
	"github.com/wiseco/core-platform/services"
	"github.com/wiseco/core-platform/services/data"
	"github.com/wiseco/core-platform/shared"
)

type partnerDatastore struct {
	sourceReq services.SourceRequest
	*sqlx.DB
}

type PartnerService interface {

	// Verify partner code
	Verify(*PartnerVerification) (*User, error)
}

func NewPartnerService(r services.SourceRequest) PartnerService {
	return &partnerDatastore{r, data.DBWrite}
}

func (db *partnerDatastore) Verify(v *PartnerVerification) (*User, error) {
	if v.Code == "" {
		return nil, errors.New("Partner code cannot be empty")
	}

	p, err := db.getByCode(v.UserID, v.Code)
	if err == sql.ErrNoRows {
		return nil, errors.New("Invalid partner code")
	}

	if err != nil {
		return nil, err
	}

	// Check partner's remaining license count
	count, err := NewUserService(db.sourceReq).GetPartnerCodeCount(p.ID)
	if err != nil {
		return nil, err
	}

	if (p.GrantedLicenseCount - count) <= 0 {
		return nil, errors.New("Invalid partner code")
	}

	// Update user's partner code id
	u := UserPartnerUpdate{
		ID:        v.UserID,
		PartnerID: p.ID,
	}

	user, err := NewUserService(db.sourceReq).UpdatePartnerCode(u)
	if err != nil {
		return nil, err
	}

	return user, nil

}

func (db *partnerDatastore) getByCode(userID shared.UserID, code string) (*Partner, error) {
	// Check access
	err := auth.NewAuthService(db.sourceReq).CheckUserAccess(userID)
	if err != nil {
		return nil, err
	}

	partner := Partner{}

	err = db.Get(&partner, "SELECT * FROM channel_partner WHERE code = $1", code)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return &partner, err
}
