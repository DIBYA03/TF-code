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
 * BankValidatorFailure prefix + uuid will be used internally when passed between services and
 * removed when passing to/from the database.
 */
type BankValidatorFailureID uuid.UUID

func (id BankValidatorFailureID) Prefix() IDPrefix {
	return IDPrefixBankValidatorFailure
}

func (id BankValidatorFailureID) String() string {
	return string(IDPrefixBankValidatorFailure) + id.UUIDString()
}

func (id BankValidatorFailureID) UUIDString() string {
	return uuid.UUID(id).String()
}

func (id BankValidatorFailureID) IsZero() bool {
	return uuid.UUID(id) == uuid.Nil
}

// Returns empty string when zero for omitempty
func (id BankValidatorFailureID) JSONString() string {
	if id.IsZero() {
		return ""
	}

	return id.String()
}

func (id BankValidatorFailureID) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.String())
}

func (id *BankValidatorFailureID) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	uid, err := ParseBankValidatorFailureID(s)
	if err != nil {
		return err
	}

	*id = uid
	return nil
}

// SQL value marshaller
func (id BankValidatorFailureID) Value() (driver.Value, error) {
	return id.UUIDString(), nil
}

// SQL scanner
func (id *BankValidatorFailureID) Scan(value interface{}) error {
	if value == nil {
		*id = BankValidatorFailureID{}
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

	*id = BankValidatorFailureID(uid)
	return nil
}

func ParseBankValidatorFailureID(id string) (BankValidatorFailureID, error) {
	// Return nil id on empty string
	if id == "" {
		return BankValidatorFailureID{}, nil
	}

	// Check prefix
	if !strings.HasPrefix(id, string(IDPrefixBankValidatorFailure)) {
		return BankValidatorFailureID{}, errors.New("invalid bank account id prefix")
	}

	// Validate UUID
	uid, err := uuid.Parse(strings.TrimPrefix(id, string(IDPrefixBankValidatorFailure)))
	if err != nil {
		return BankValidatorFailureID{}, err
	}

	return BankValidatorFailureID(uid), nil
}
