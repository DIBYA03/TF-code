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
 * BankTransaction prefix + uuid will be used internally when passed between services and
 * removed when passing to/from the database.
 */
type BankTransactionID uuid.UUID

func (id BankTransactionID) Prefix() IDPrefix {
	return IDPrefixBankTransaction
}

func (id BankTransactionID) String() string {
	return string(IDPrefixBankTransaction) + id.UUIDString()
}

func (id BankTransactionID) UUIDString() string {
	return uuid.UUID(id).String()
}

func (id BankTransactionID) IsZero() bool {
	return uuid.UUID(id) == uuid.Nil
}

// Returns empty string when zero for omitempty
func (id BankTransactionID) JSONString() string {
	if id.IsZero() {
		return ""
	}

	return id.String()
}

func (id *BankTransactionID) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.String())
}

func (id *BankTransactionID) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	uid, err := ParseBankTransactionID(s)
	if err != nil {
		return err
	}

	*id = BankTransactionID(uid)
	return nil
}

// SQL value marshaller
func (id BankTransactionID) Value() (driver.Value, error) {
	return id.UUIDString(), nil
}

// SQL scanner
func (id *BankTransactionID) Scan(value interface{}) error {
	if value == nil {
		*id = BankTransactionID{}
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

	*id = BankTransactionID(uid)
	return nil
}

func NewBankTransactionID() (BankTransactionID, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return BankTransactionID{}, err
	}

	return BankTransactionID(id), nil
}

func ParseBankTransactionID(id string) (BankTransactionID, error) {
	// Return nil id on empty string
	if id == "" {
		return BankTransactionID{}, nil
	}

	// Check prefix
	if !strings.HasPrefix(id, string(IDPrefixBankTransaction)) {
		return BankTransactionID{}, errors.New("invalid bank transaction id prefix")
	}

	// Validate UUID
	uid, err := uuid.Parse(strings.TrimPrefix(id, string(IDPrefixBankTransaction)))
	if err != nil {
		return BankTransactionID{}, err
	}

	return BankTransactionID(uid), nil
}
