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
 * DailyBalance prefix + uuid will be used internally when passed between services and
 * removed when passing to/from the database.
 */
type DailyBalanceID uuid.UUID

func (id DailyBalanceID) Prefix() IDPrefix {
	return IDPrefixDailyBalance
}

func (id DailyBalanceID) String() string {
	return string(IDPrefixDailyBalance) + id.UUIDString()
}

func (id DailyBalanceID) UUIDString() string {
	return uuid.UUID(id).String()
}

func (id DailyBalanceID) IsZero() bool {
	return uuid.UUID(id) == uuid.Nil
}

// Returns empty string when zero for omitempty
func (id DailyBalanceID) JSONString() string {
	if id.IsZero() {
		return ""
	}

	return id.String()
}

func (id DailyBalanceID) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.String())
}

func (id *DailyBalanceID) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	uid, err := ParseDailyBalanceID(s)
	if err != nil {
		return err
	}

	*id = uid
	return nil
}

// SQL value marshaller
func (id DailyBalanceID) Value() (driver.Value, error) {
	return id.UUIDString(), nil
}

// SQL scanner
func (id *DailyBalanceID) Scan(value interface{}) error {
	if value == nil {
		*id = DailyBalanceID{}
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

	*id = DailyBalanceID(uid)
	return nil
}

func ParseDailyBalanceID(id string) (DailyBalanceID, error) {
	// Return nil id on empty string
	if id == "" {
		return DailyBalanceID{}, nil
	}

	// Check prefix
	if !strings.HasPrefix(id, string(IDPrefixDailyBalance)) {
		return DailyBalanceID{}, errors.New("invalid daily balance id prefix")
	}

	// Validate UUID
	uid, err := uuid.Parse(strings.TrimPrefix(id, string(IDPrefixDailyBalance)))
	if err != nil {
		return DailyBalanceID{}, err
	}

	return DailyBalanceID(uid), nil
}
