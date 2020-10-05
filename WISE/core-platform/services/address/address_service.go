/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

// Package for address
package address

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/wiseco/core-platform/services"
	"github.com/wiseco/core-platform/services/data"
	"github.com/wiseco/core-platform/shared"
)

type addressDatastore struct {
	sourceReq services.SourceRequest
	*sqlx.DB
}

type AddressService interface {
	// Read
	GetByID(shared.AddressID) (*Address, error)
	GetByConsumerID(shared.ConsumerID, int, int) ([]Address, error)
	GetByContactID(shared.ContactID, int, int) ([]Address, error)
	GetByBusinessID(shared.BusinessID, int, int) ([]Address, error)

	// Create/modify
	Create(*AddressCreate) (*Address, error)
	Update(*AddressUpdate) (*Address, error)

	// Deactivate address
	Deactivate(shared.AddressID) error
	DeactivateAllForContact(shared.ContactID) error

	//DEBUG
	DEBUGDeleteAllForBusiness(shared.BusinessID) error
	DEBUGDeleteAllForContact(shared.ContactID) error
	DEBUGDeleteAllForConsumer(shared.ConsumerID) error
}

func NewAddressService(r services.SourceRequest) AddressService {
	return &addressDatastore{r, data.DBWrite}
}

//TODO Auth should be moved really
//We should restrict access based on the api not at the data level
//We can check access per method and id type but gets weird with consumer and contact

func (db *addressDatastore) GetByID(id shared.AddressID) (*Address, error) {
	a := new(Address)

	err := db.Get(a, "SELECT * FROM address WHERE id = $1 and address_state = $2", id, AddressStateActive)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return a, err
}

func (db *addressDatastore) GetByConsumerID(consumerID shared.ConsumerID, limit int, offset int) ([]Address, error) {
	rows := []Address{}

	err := db.Select(&rows, "SELECT * FROM address WHERE consumer_id = $1 and address_state = $2 LIMIT $3 OFFSET $4", consumerID, AddressStateActive, limit, offset)

	switch {
	case err == sql.ErrNoRows:
		return rows, nil
	case err != nil:
		return nil, err
	}

	return rows, err
}

func (db *addressDatastore) GetByContactID(contactID shared.ContactID, limit int, offset int) ([]Address, error) {
	rows := []Address{}

	err := db.Select(&rows, "SELECT * FROM address WHERE contact_id = $1 and address_state = $2 LIMIT $3 OFFSET $4", contactID, AddressStateActive, limit, offset)

	switch {
	case err == sql.ErrNoRows:
		return rows, nil
	case err != nil:
		return nil, err
	}

	return rows, err
}

func (db *addressDatastore) GetByBusinessID(businessID shared.BusinessID, limit int, offset int) ([]Address, error) {
	// Check access
	rows := []Address{}

	err := db.Select(&rows, "SELECT * FROM address WHERE business_id = $1 and address_state = $2 LIMIT $2 OFFSET $3", businessID, AddressStateActive, limit, offset)

	switch {
	case err == sql.ErrNoRows:
		return rows, nil
	case err != nil:
		return nil, err
	}

	return rows, err
}

func (db *addressDatastore) Create(addressCreate *AddressCreate) (*Address, error) {

	// Default/mandatory fields
	columns := []string{
		"street", "locality", "admin_area", "country", "postal_code", "latitude", "longitude", "address_type", "address_state",
	}
	// Default/mandatory values
	values := []string{
		":street", ":locality", ":admin_area", ":country", ":postal_code", ":latitude", ":longitude", ":address_type", ":address_state",
	}

	if addressCreate.AddressLine2 != "" {
		columns = append(columns, "line2")
		values = append(values, ":line2")
	}

	if addressCreate.ConsumerID != "" {
		columns = append(columns, "consumer_id")
		values = append(values, ":consumer_id")
	} else if addressCreate.ContactID != "" {
		columns = append(columns, "contact_id")
		values = append(values, ":contact_id")
	} else if addressCreate.BusinessID != "" {
		columns = append(columns, "business_id")
		values = append(values, ":business_id")
	} else {
		return nil, errors.New("Must pass either a ConsumerID, ContactID or a BusinessID")
	}

	addressCreate.AddressState = AddressStateActive

	sql := fmt.Sprintf("INSERT INTO address(%s) VALUES(%s) RETURNING *", strings.Join(columns, ", "), strings.Join(values, ", "))

	stmt, err := db.PrepareNamed(sql)
	if err != nil {
		return nil, err
	}

	address := new(Address)

	err = stmt.Get(address, &addressCreate)
	if err != nil {
		return nil, err
	}

	return address, nil
}

func (db *addressDatastore) Update(addressUpdate *AddressUpdate) (*Address, error) {
	a, err := db.GetByID(addressUpdate.ID)
	if err != nil {
		return nil, err
	}

	var columns []string

	if addressUpdate.StreetAddress != "" {
		columns = append(columns, "street = :street")
	}

	if addressUpdate.AddressLine2 != "" {
		columns = append(columns, "line2 = :line2")
	}

	if addressUpdate.Locality != "" {
		columns = append(columns, "locality = :locality")
	}

	if addressUpdate.AdminArea != "" {
		columns = append(columns, "admin_area = :admin_area")
	}

	if addressUpdate.Country != "" {
		columns = append(columns, "country = :country")
	}

	if addressUpdate.PostalCode != "" {
		columns = append(columns, "postal_code = :postal_code")
	}

	if addressUpdate.Latitude != float64(0) {
		columns = append(columns, "latitude = :latitude")
	}

	if addressUpdate.Longitude != float64(0) {
		columns = append(columns, "longitude = :longitude")
	}

	// No changes requested - return address
	if len(columns) == 0 {
		return a, nil
	}

	// Update Address
	_, err = db.NamedExec(
		fmt.Sprintf(
			"UPDATE address SET %s WHERE id = '%s'",
			strings.Join(columns, ", "),
			a.ID,
		), addressUpdate,
	)

	if err != nil {
		return nil, errors.Cause(err)
	}

	return db.GetByID(a.ID)
}

func (db *addressDatastore) Deactivate(ID shared.AddressID) error {
	_, err := db.Exec("UPDATE address SET address_state = $1 WHERE id = $2", AddressStateInactive, ID)

	return err
}

func (db *addressDatastore) DeactivateAllForContact(contactID shared.ContactID) error {
	as, err := db.GetByContactID(contactID, 50, 0)

	if err != nil {
		return err
	}

	return db.deactivateAll(as)
}

func (db *addressDatastore) deactivateAll(as []Address) error {
	//Nothing to do
	if len(as) == 0 {
		return nil
	}

	for _, a := range as {
		err := db.Deactivate(a.ID)

		if err != nil {
			return err
		}
	}

	return nil
}

//ONLY for debug
func (db *addressDatastore) DEBUGDeleteAllForBusiness(businessID shared.BusinessID) error {
	_, err := db.Exec("DELETE FROM address where business_id = $1", businessID)

	return err
}

//ONLY for debug
func (db *addressDatastore) DEBUGDeleteAllForContact(contactID shared.ContactID) error {
	_, err := db.Exec("DELETE FROM address where contact_id = $1", contactID)

	return err
}

//ONLY for debug
func (db *addressDatastore) DEBUGDeleteAllForConsumer(consumerID shared.ConsumerID) error {
	_, err := db.Exec("DELETE FROM address where consumer_id = $1", consumerID)

	return err
}
