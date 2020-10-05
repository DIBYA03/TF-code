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
 * Address prefix + uuid will be used internally when passed between services and
 * removed when passing to/from the database.
 */
type AddressID uuid.UUID

func (id AddressID) Prefix() IDPrefix {
	return IDPrefixAddress
}

func (id AddressID) String() string {
	return string(IDPrefixAddress) + id.UUIDString()
}

func (id AddressID) UUIDString() string {
	return uuid.UUID(id).String()
}

func (id AddressID) IsZero() bool {
	return uuid.UUID(id) == uuid.Nil
}

// Returns empty string when zero for omitempty
func (id AddressID) JSONString() string {
	if id.IsZero() {
		return ""
	}

	return id.String()
}

func (id AddressID) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.String())
}

func (id *AddressID) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	uid, err := ParseAddressID(s)
	if err != nil {
		return err
	}

	*id = uid
	return nil
}

// SQL value marshaller
func (id AddressID) Value() (driver.Value, error) {
	return id.UUIDString(), nil
}

// SQL scanner
func (id *AddressID) Scan(value interface{}) error {
	if value == nil {
		*id = AddressID{}
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

	*id = AddressID(uid)
	return nil
}

func ParseAddressID(id string) (AddressID, error) {
	// Return nil id on empty string
	if id == "" {
		return AddressID{}, nil
	}

	// Check prefix
	if !strings.HasPrefix(id, string(IDPrefixAddress)) {
		return AddressID{}, errors.New("invalid address id prefix")
	}

	// Validate UUID
	uid, err := uuid.Parse(strings.TrimPrefix(id, string(IDPrefixAddress)))
	if err != nil {
		return AddressID{}, err
	}

	return AddressID(uid), nil
}
