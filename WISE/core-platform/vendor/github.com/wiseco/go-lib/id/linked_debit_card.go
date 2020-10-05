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
 * LinkedDebitCard prefix + uuid will be used internally when passed between services and
 * removed when passing to/from the database.
 */
type LinkedDebitCardID uuid.UUID

func (id LinkedDebitCardID) Prefix() IDPrefix {
	return IDPrefixLinkedDebitCard
}

func (id LinkedDebitCardID) String() string {
	return string(IDPrefixLinkedDebitCard) + id.UUIDString()
}

func (id LinkedDebitCardID) UUIDString() string {
	return uuid.UUID(id).String()
}

func (id LinkedDebitCardID) IsZero() bool {
	return uuid.UUID(id) == uuid.Nil
}

// Returns empty string when zero for omitempty
func (id LinkedDebitCardID) JSONString() string {
	if id.IsZero() {
		return ""
	}

	return id.String()
}

func (id LinkedDebitCardID) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.String())
}

func (id *LinkedDebitCardID) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	uid, err := ParseLinkedDebitCardID(s)
	if err != nil {
		return err
	}

	*id = uid
	return nil
}

// SQL value marshaller
func (id LinkedDebitCardID) Value() (driver.Value, error) {
	return id.UUIDString(), nil
}

// SQL scanner
func (id *LinkedDebitCardID) Scan(value interface{}) error {
	if value == nil {
		*id = LinkedDebitCardID{}
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

	*id = LinkedDebitCardID(uid)
	return nil
}

func ParseLinkedDebitCardID(id string) (LinkedDebitCardID, error) {
	// Return nil id on empty string
	if id == "" {
		return LinkedDebitCardID{}, nil
	}

	// Check prefix
	if !strings.HasPrefix(id, string(IDPrefixLinkedDebitCard)) {
		return LinkedDebitCardID{}, errors.New("invalid linked debit card id prefix")
	}

	// Validate UUID
	uid, err := uuid.Parse(strings.TrimPrefix(id, string(IDPrefixLinkedDebitCard)))
	if err != nil {
		return LinkedDebitCardID{}, err
	}

	return LinkedDebitCardID(uid), nil
}
