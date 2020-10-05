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
 * DailyTransaction prefix + uuid will be used internally when passed between services and
 * removed when passing to/from the database.
 */
type DailyTransactionID uuid.UUID

func (id DailyTransactionID) Prefix() IDPrefix {
	return IDPrefixDailyTransaction
}

func (id DailyTransactionID) String() string {
	return string(IDPrefixDailyTransaction) + id.UUIDString()
}

func (id DailyTransactionID) UUIDString() string {
	return uuid.UUID(id).String()
}

func (id DailyTransactionID) IsZero() bool {
	return uuid.UUID(id) == uuid.Nil
}

// Returns empty string when zero for omitempty
func (id DailyTransactionID) JSONString() string {
	if id.IsZero() {
		return ""
	}

	return id.String()
}

func (id DailyTransactionID) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.String())
}

func (id *DailyTransactionID) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	uid, err := ParseDailyTransactionID(s)
	if err != nil {
		return err
	}

	*id = uid
	return nil
}

// SQL value marshaller
func (id DailyTransactionID) Value() (driver.Value, error) {
	return id.UUIDString(), nil
}

// SQL scanner
func (id *DailyTransactionID) Scan(value interface{}) error {
	if value == nil {
		*id = DailyTransactionID{}
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

	*id = DailyTransactionID(uid)
	return nil
}

func ParseDailyTransactionID(id string) (DailyTransactionID, error) {
	// Return nil id on empty string
	if id == "" {
		return DailyTransactionID{}, nil
	}

	// Check prefix
	if !strings.HasPrefix(id, string(IDPrefixDailyTransaction)) {
		return DailyTransactionID{}, errors.New("invalid daily transaction id prefix")
	}

	// Validate UUID
	uid, err := uuid.Parse(strings.TrimPrefix(id, string(IDPrefixDailyTransaction)))
	if err != nil {
		return DailyTransactionID{}, err
	}

	return DailyTransactionID(uid), nil
}
