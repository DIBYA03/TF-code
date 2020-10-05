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
 * Business prefix + uuid will be used internally when passed between services and
 * removed when passing to/from the database.
 */
type BusinessID uuid.UUID

func (id BusinessID) Prefix() IDPrefix {
	return IDPrefixBusiness
}

func (id BusinessID) String() string {
	return string(IDPrefixBusiness) + id.UUIDString()
}

func (id BusinessID) UUIDString() string {
	return uuid.UUID(id).String()
}

func (id BusinessID) IsZero() bool {
	return uuid.UUID(id) == uuid.Nil
}

// Returns empty string when zero for omitempty
func (id BusinessID) JSONString() string {
	if id.IsZero() {
		return ""
	}

	return id.String()
}

func (id BusinessID) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.String())
}

func (id *BusinessID) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	uid, err := ParseBusinessID(s)
	if err != nil {
		return err
	}

	*id = uid
	return nil
}

// SQL value marshaller
func (id BusinessID) Value() (driver.Value, error) {
	return id.UUIDString(), nil
}

// SQL scanner
func (id *BusinessID) Scan(value interface{}) error {
	if value == nil {
		*id = BusinessID{}
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

	*id = BusinessID(uid)
	return nil
}

func ParseBusinessID(id string) (BusinessID, error) {
	// Return nil id on empty string
	if id == "" {
		return BusinessID{}, nil
	}

	// Check prefix
	if !strings.HasPrefix(id, string(IDPrefixBusiness)) {
		return BusinessID{}, errors.New("invalid business id prefix")
	}

	// Validate UUID
	uid, err := uuid.Parse(strings.TrimPrefix(id, string(IDPrefixBusiness)))
	if err != nil {
		return BusinessID{}, err
	}

	return BusinessID(uid), nil
}
