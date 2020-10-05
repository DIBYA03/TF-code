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
 * BankAccount prefix + uuid will be used internally when passed between services and
 * removed when passing to/from the database.
 */
type BankAccountID uuid.UUID

func (id BankAccountID) Prefix() IDPrefix {
	return IDPrefixBankAccount
}

func (id BankAccountID) String() string {
	return string(IDPrefixBankAccount) + id.UUIDString()
}

func (id BankAccountID) UUIDString() string {
	return uuid.UUID(id).String()
}

func (id BankAccountID) IsZero() bool {
	return uuid.UUID(id) == uuid.Nil
}

// Returns empty string when zero for omitempty
func (id BankAccountID) JSONString() string {
	if id.IsZero() {
		return ""
	}

	return id.String()
}

func (id BankAccountID) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.String())
}

func (id *BankAccountID) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	uid, err := ParseBankAccountID(s)
	if err != nil {
		return err
	}

	*id = uid
	return nil
}

// SQL value marshaller
func (id BankAccountID) Value() (driver.Value, error) {
	return id.UUIDString(), nil
}

// SQL scanner
func (id *BankAccountID) Scan(value interface{}) error {
	if value == nil {
		*id = BankAccountID{}
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

	*id = BankAccountID(uid)
	return nil
}

func ParseBankAccountID(id string) (BankAccountID, error) {
	// Return nil id on empty string
	if id == "" {
		return BankAccountID{}, nil
	}

	// Check prefix
	if !strings.HasPrefix(id, string(IDPrefixBankAccount)) {
		return BankAccountID{}, errors.New("invalid bank account id prefix")
	}

	// Validate UUID
	uid, err := uuid.Parse(strings.TrimPrefix(id, string(IDPrefixBankAccount)))
	if err != nil {
		return BankAccountID{}, err
	}

	return BankAccountID(uid), nil
}
