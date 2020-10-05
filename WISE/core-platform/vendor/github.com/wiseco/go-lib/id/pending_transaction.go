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
 * PendingTransaction prefix + uuid will be used internally when passed between services and
 * removed when passing to/from the database.
 */
type PendingTransactionID uuid.UUID

func (id PendingTransactionID) Prefix() IDPrefix {
	return IDPrefixPendingTransaction
}

func (id PendingTransactionID) String() string {
	return string(IDPrefixPendingTransaction) + id.UUIDString()
}

func (id PendingTransactionID) UUIDString() string {
	return uuid.UUID(id).String()
}

func (id PendingTransactionID) IsZero() bool {
	return uuid.UUID(id) == uuid.Nil
}

// Returns empty string when zero for omitempty
func (id PendingTransactionID) JSONString() string {
	if id.IsZero() {
		return ""
	}

	return id.String()
}

func (id *PendingTransactionID) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.String())
}

func (id *PendingTransactionID) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	uid, err := ParsePendingTransactionID(s)
	if err != nil {
		return err
	}

	*id = PendingTransactionID(uid)
	return nil
}

// SQL value marshaller
func (id PendingTransactionID) Value() (driver.Value, error) {
	return id.UUIDString(), nil
}

// SQL scanner
func (id *PendingTransactionID) Scan(value interface{}) error {
	if value == nil {
		*id = PendingTransactionID{}
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

	*id = PendingTransactionID(uid)
	return nil
}

func ParsePendingTransactionID(id string) (PendingTransactionID, error) {
	// Return nil id on empty string
	if id == "" {
		return PendingTransactionID{}, nil
	}

	// Check prefix
	if !strings.HasPrefix(id, string(IDPrefixPendingTransaction)) {
		return PendingTransactionID{}, errors.New("invalid pending transaction id prefix")
	}

	// Validate UUID
	uid, err := uuid.Parse(strings.TrimPrefix(id, string(IDPrefixPendingTransaction)))
	if err != nil {
		return PendingTransactionID{}, err
	}

	return PendingTransactionID(uid), nil
}
