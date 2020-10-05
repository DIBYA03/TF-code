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
 * DebitCard prefix + uuid will be used internally when passed between services and
 * removed when passing to/from the database.
 */
type DebitCardID uuid.UUID

func (id DebitCardID) Prefix() IDPrefix {
	return IDPrefixDebitCard
}

func (id DebitCardID) String() string {
	return string(IDPrefixDebitCard) + id.UUIDString()
}

func (id DebitCardID) UUIDString() string {
	return uuid.UUID(id).String()
}

func (id DebitCardID) IsZero() bool {
	return uuid.UUID(id) == uuid.Nil
}

// Returns empty string when zero for omitempty
func (id DebitCardID) JSONString() string {
	if id.IsZero() {
		return ""
	}

	return id.String()
}

func (id DebitCardID) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.String())
}

func (id *DebitCardID) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	uid, err := ParseDebitCardID(s)
	if err != nil {
		return err
	}

	*id = uid
	return nil
}

// SQL value marshaller
func (id DebitCardID) Value() (driver.Value, error) {
	return id.UUIDString(), nil
}

// SQL scanner
func (id *DebitCardID) Scan(value interface{}) error {
	if value == nil {
		*id = DebitCardID{}
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

	*id = DebitCardID(uid)
	return nil
}

func ParseDebitCardID(id string) (DebitCardID, error) {
	// Return nil id on empty string
	if id == "" {
		return DebitCardID{}, nil
	}

	// Check prefix
	if !strings.HasPrefix(id, string(IDPrefixDebitCard)) {
		return DebitCardID{}, errors.New("invalid debit card id prefix")
	}

	// Validate UUID
	uid, err := uuid.Parse(strings.TrimPrefix(id, string(IDPrefixDebitCard)))
	if err != nil {
		return DebitCardID{}, err
	}

	return DebitCardID(uid), nil
}
