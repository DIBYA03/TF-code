package id

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

/*
 * Contact prefix + uuid will be used internally when passed between services and
 * removed when passing to/from the database.
 */
type ContactID uuid.UUID

func (id ContactID) Prefix() IDPrefix {
	return IDPrefixContact
}

func (id ContactID) String() string {
	return string(IDPrefixContact) + id.UUIDString()
}

func (id ContactID) UUIDString() string {
	return uuid.UUID(id).String()
}

func (id ContactID) IsZero() bool {
	return uuid.UUID(id) == uuid.Nil
}

// Returns empty string when zero for omitempty
func (id ContactID) JSONString() string {
	if id.IsZero() {
		return ""
	}

	return id.String()
}

func (id ContactID) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.String())
}

func (id *ContactID) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	uid, err := ParseContactID(s)
	if err != nil {
		return err
	}

	*id = uid
	return nil
}

// SQL value marshaller
func (id ContactID) Value() (driver.Value, error) {
	return id.UUIDString(), nil
}

// SQL scanner
func (id *ContactID) Scan(value interface{}) error {
	if value == nil {
		*id = ContactID{}
		return nil
	}

	val, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("Unable to scan value %v", value)
	}

	uid, err := uuid.Parse(string(val))
	if err != nil {
		return err
	}

	*id = ContactID(uid)
	return nil
}

func ParseContactID(id string) (ContactID, error) {
	// Return nil id on empty string
	if id == "" {
		return ContactID{}, nil
	}

	// Check prefix
	if !strings.HasPrefix(id, string(IDPrefixContact)) {
		return ContactID{}, errors.New("invalid contact id prefix")
	}

	// Validate UUID
	uid, err := uuid.Parse(strings.TrimPrefix(id, string(IDPrefixContact)))
	if err != nil {
		return ContactID{}, err
	}

	return ContactID(uid), nil
}
