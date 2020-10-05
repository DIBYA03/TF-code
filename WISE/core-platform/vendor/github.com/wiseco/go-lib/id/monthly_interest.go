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
 * MonthlyInterest prefix + uuid will be used internally when passed between services and
 * removed when passing to/from the database.
 */
type MonthlyInterestID uuid.UUID

func (id MonthlyInterestID) Prefix() IDPrefix {
	return IDPrefixMonthlyInterest
}

func (id MonthlyInterestID) String() string {
	return string(IDPrefixMonthlyInterest) + id.UUIDString()
}

func (id MonthlyInterestID) UUIDString() string {
	return uuid.UUID(id).String()
}

func (id MonthlyInterestID) IsZero() bool {
	return uuid.UUID(id) == uuid.Nil
}

// Returns empty string when zero for omitempty
func (id MonthlyInterestID) JSONString() string {
	if id.IsZero() {
		return ""
	}

	return id.String()
}

func (id MonthlyInterestID) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.String())
}

func (id *MonthlyInterestID) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	uid, err := ParseMonthlyInterestID(s)
	if err != nil {
		return err
	}

	*id = uid
	return nil
}

// SQL value marshaller
func (id MonthlyInterestID) Value() (driver.Value, error) {
	return id.UUIDString(), nil
}

// SQL scanner
func (id *MonthlyInterestID) Scan(value interface{}) error {
	if value == nil {
		*id = MonthlyInterestID{}
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

	*id = MonthlyInterestID(uid)
	return nil
}

func ParseMonthlyInterestID(id string) (MonthlyInterestID, error) {
	// Return nil id on empty string
	if id == "" {
		return MonthlyInterestID{}, nil
	}

	// Check prefix
	if !strings.HasPrefix(id, string(IDPrefixMonthlyInterest)) {
		return MonthlyInterestID{}, errors.New("invalid monthly interest id prefix")
	}

	// Validate UUID
	uid, err := uuid.Parse(strings.TrimPrefix(id, string(IDPrefixMonthlyInterest)))
	if err != nil {
		return MonthlyInterestID{}, err
	}

	return MonthlyInterestID(uid), nil
}
