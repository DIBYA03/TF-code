/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

package payment

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/wiseco/core-platform/auth"
	"github.com/wiseco/core-platform/services"
	"github.com/wiseco/core-platform/services/data"
	"github.com/wiseco/core-platform/shared"
)

type cardReaderDatastore struct {
	sourceReq services.SourceRequest
	*sqlx.DB
}

type CardReaderService interface {
	List(userID shared.UserID, businessID shared.BusinessID) ([]CardReader, error)
	GetById(ID shared.CardReaderID, userID shared.UserID, businessID shared.BusinessID) (*CardReader, error)
	Create(CardReaderCreate) (*CardReader, error)
	Update(CardReaderUpdate) (*CardReader, error)
	Deactivate(ID shared.CardReaderID, userID shared.UserID, businessID shared.BusinessID) error
}

func NewCardReaderService(r services.SourceRequest) CardReaderService {
	return &cardReaderDatastore{r, data.DBWrite}
}

func (db *cardReaderDatastore) List(userID shared.UserID, businessID shared.BusinessID) ([]CardReader, error) {
	// Check access
	err := auth.NewAuthService(db.sourceReq).CheckBusinessAccess(businessID)
	if err != nil {
		return nil, err
	}

	rows := []CardReader{}
	err = db.Select(&rows, "SELECT * FROM card_reader WHERE deactivated IS NULL AND business_id = $1", businessID)
	switch {
	case err == sql.ErrNoRows:
		return rows, nil
	case err != nil:
		return nil, errors.Cause(err)
	}

	return rows, err
}

func (db *cardReaderDatastore) GetById(id shared.CardReaderID, userID shared.UserID, businessID shared.BusinessID) (*CardReader, error) {
	// Check access
	err := auth.NewAuthService(db.sourceReq).CheckBusinessAccess(businessID)
	if err != nil {
		return nil, err
	}

	cardReader := CardReader{}

	err = db.Get(&cardReader, "SELECT * FROM card_reader WHERE id = $1 AND business_id = $2", id, businessID)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return &cardReader, err
}

func (db *cardReaderDatastore) Create(c CardReaderCreate) (*CardReader, error) {
	// Check access
	err := auth.NewAuthService(db.sourceReq).CheckBusinessAccess(c.BusinessID)
	if err != nil {
		return nil, err
	}

	if c.SerialNumber == "" {
		return nil, errors.New("Serial number is required")
	}

	// Default/mandatory fields
	columns := []string{
		"business_id", "alias", "device_type", "serial_number",
	}
	// Default/mandatory values
	values := []string{
		":business_id", ":alias", ":device_type", ":serial_number",
	}

	sql := fmt.Sprintf("INSERT INTO card_reader(%s) VALUES(%s) RETURNING *", strings.Join(columns, ", "), strings.Join(values, ", "))

	stmt, err := db.PrepareNamed(sql)
	if err != nil {
		return nil, err
	}

	cardReader := &CardReader{}

	err = stmt.Get(cardReader, &c)
	if err != nil {
		return nil, err
	}

	return cardReader, nil
}

func (db *cardReaderDatastore) Update(u CardReaderUpdate) (*CardReader, error) {
	// Check access
	err := auth.NewAuthService(db.sourceReq).CheckBusinessAccess(u.BusinessID)
	if err != nil {
		return nil, err
	}

	p, err := db.GetById(u.ID, u.CreatedUserID, u.BusinessID)
	if err != nil {
		return nil, err
	}

	var columns []string

	if u.Alias != nil {
		columns = append(columns, "alias = :alias")
	}

	if u.LastConnected != nil {
		columns = append(columns, "last_connected = :last_connected")
	}

	if u.DeviceType != nil {
		columns = append(columns, "device_type = :device_type")
	}

	// No changes requested - return cardReader
	if len(columns) == 0 {
		return p, nil
	}

	_, err = db.NamedExec(
		fmt.Sprintf(
			"UPDATE card_reader SET %s WHERE id = '%s'",
			strings.Join(columns, ", "),
			u.ID,
		), u,
	)

	if err != nil {
		return nil, errors.Cause(err)
	}

	c, err := db.GetById(u.ID, u.CreatedUserID, u.BusinessID)

	return c, err
}

func (db *cardReaderDatastore) Deactivate(ID shared.CardReaderID, userID shared.UserID, businessID shared.BusinessID) error {
	// Check access
	err := auth.NewAuthService(db.sourceReq).CheckBusinessAccess(businessID)
	if err != nil {
		return err
	}

	_, err = db.Exec("UPDATE card_reader SET deactivated = CURRENT_TIMESTAMP WHERE id = $1 AND business_id = $2", ID, businessID)
	return err
}
