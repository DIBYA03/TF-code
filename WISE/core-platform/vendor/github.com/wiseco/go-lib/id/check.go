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
 * Check prefix + uuid will be used internally when passed between services and
 * removed when passing to/from the database.
 */
type CheckID uuid.UUID

func (id CheckID) Prefix() IDPrefix {
	return IDPrefixCheck
}

func (id CheckID) String() string {
	return string(IDPrefixCheck) + id.UUIDString()
}

func (id CheckID) UUIDString() string {
	return uuid.UUID(id).String()
}

func (id CheckID) IsZero() bool {
	return uuid.UUID(id) == uuid.Nil
}

// Returns empty string when zero for omitempty
func (id CheckID) JSONString() string {
	if id.IsZero() {
		return ""
	}

	return id.String()
}

func (id CheckID) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.String())
}

func (id *CheckID) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	uid, err := ParseCheckID(s)
	if err != nil {
		return err
	}

	*id = uid
	return nil
}

// SQL value marshaller
func (id CheckID) Value() (driver.Value, error) {
	return id.UUIDString(), nil
}

// SQL scanner
func (id *CheckID) Scan(value interface{}) error {
	if value == nil {
		*id = CheckID{}
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

	*id = CheckID(uid)
	return nil
}

func ParseCheckID(id string) (CheckID, error) {
	// Return nil id on empty string
	if id == "" {
		return CheckID{}, nil
	}

	// Check prefix
	if !strings.HasPrefix(id, string(IDPrefixCheck)) {
		return CheckID{}, errors.New("invalid check id prefix")
	}

	// Validate UUID
	uid, err := uuid.Parse(strings.TrimPrefix(id, string(IDPrefixCheck)))
	if err != nil {
		return CheckID{}, err
	}

	return CheckID(uid), nil
}
