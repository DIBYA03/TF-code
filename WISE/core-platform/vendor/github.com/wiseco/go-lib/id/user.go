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
 * User prefix + uuid will be used internally when passed between services and
 * removed when passing to/from the database.
 */
type UserID uuid.UUID

func (id UserID) Prefix() IDPrefix {
	return IDPrefixUser
}

func (id UserID) String() string {
	return string(IDPrefixUser) + id.UUIDString()
}

func (id UserID) UUIDString() string {
	return uuid.UUID(id).String()
}

func (id UserID) IsZero() bool {
	return uuid.UUID(id) == uuid.Nil
}

// Returns empty string when zero for omitempty
func (id UserID) JSONString() string {
	if id.IsZero() {
		return ""
	}

	return id.String()
}

func (id UserID) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.String())
}

func (id *UserID) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	uid, err := ParseUserID(s)
	if err != nil {
		return err
	}

	*id = uid
	return nil
}

// SQL value marshaller
func (id UserID) Value() (driver.Value, error) {
	return id.UUIDString(), nil
}

// SQL scanner
func (id *UserID) Scan(value interface{}) error {
	if value == nil {
		*id = UserID{}
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

	*id = UserID(uid)
	return nil
}

func ParseUserID(id string) (UserID, error) {
	// Return nil id on empty string
	if id == "" {
		return UserID{}, nil
	}

	// Check prefix
	if !strings.HasPrefix(id, string(IDPrefixUser)) {
		return UserID{}, errors.New("invalid user id prefix")
	}

	// Validate UUID
	uid, err := uuid.Parse(strings.TrimPrefix(id, string(IDPrefixUser)))
	if err != nil {
		return UserID{}, err
	}

	return UserID(uid), nil
}
