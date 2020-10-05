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
 * Email prefix + uuid will be used internally when passed between services and
 * removed when passing to/from the database.
 */
type EmailID uuid.UUID

func (id EmailID) Prefix() IDPrefix {
	return IDPrefixEmail
}

func (id EmailID) String() string {
	return string(IDPrefixEmail) + id.UUIDString()
}

func (id EmailID) UUIDString() string {
	return uuid.UUID(id).String()
}

func (id EmailID) IsZero() bool {
	return uuid.UUID(id) == uuid.Nil
}

// Returns empty string when zero for omitempty
func (id EmailID) JSONString() string {
	if id.IsZero() {
		return ""
	}

	return id.String()
}

func (id EmailID) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.String())
}

func (id *EmailID) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	uid, err := ParseEmailID(s)
	if err != nil {
		return err
	}

	*id = uid
	return nil
}

// SQL value marshaller
func (id EmailID) Value() (driver.Value, error) {
	return id.UUIDString(), nil
}

// SQL scanner
func (id *EmailID) Scan(value interface{}) error {
	if value == nil {
		*id = EmailID{}
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

	*id = EmailID(uid)
	return nil
}

func ParseEmailID(id string) (EmailID, error) {
	// Return nil id on empty string
	if id == "" {
		return EmailID{}, nil
	}

	// Check prefix
	if !strings.HasPrefix(id, string(IDPrefixEmail)) {
		return EmailID{}, errors.New("invalid email id prefix")
	}

	// Validate UUID
	uid, err := uuid.Parse(strings.TrimPrefix(id, string(IDPrefixEmail)))
	if err != nil {
		return EmailID{}, err
	}

	return EmailID(uid), nil
}
