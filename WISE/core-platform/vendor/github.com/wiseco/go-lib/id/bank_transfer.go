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
 * BankTransfer prefix + uuid will be used internally when passed between services and
 * removed when passing to/from the database.
 */
type BankTransferID uuid.UUID

func (id BankTransferID) Prefix() IDPrefix {
	return IDPrefixBankTransfer
}

func (id BankTransferID) String() string {
	return string(IDPrefixBankTransfer) + id.UUIDString()
}

func (id BankTransferID) UUIDString() string {
	return uuid.UUID(id).String()
}

func (id BankTransferID) IsZero() bool {
	return uuid.UUID(id) == uuid.Nil
}

// Returns empty string when zero for omitempty
func (id BankTransferID) JSONString() string {
	if id.IsZero() {
		return ""
	}

	return id.String()
}

func (id BankTransferID) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.String())
}

func (id *BankTransferID) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	uid, err := ParseBankTransferID(s)
	if err != nil {
		return err
	}

	*id = uid
	return nil
}

// SQL value marshaller
func (id BankTransferID) Value() (driver.Value, error) {
	return id.UUIDString(), nil
}

// SQL scanner
func (id *BankTransferID) Scan(value interface{}) error {
	if value == nil {
		*id = BankTransferID{}
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

	*id = BankTransferID(uid)
	return nil
}

func ParseBankTransferID(id string) (BankTransferID, error) {
	// Return nil id on empty string
	if id == "" {
		return BankTransferID{}, nil
	}

	// Check prefix
	if !strings.HasPrefix(id, string(IDPrefixBankTransfer)) {
		return BankTransferID{}, errors.New("invalid bank transfer id prefix")
	}

	// Validate UUID
	uid, err := uuid.Parse(strings.TrimPrefix(id, string(IDPrefixBankTransfer)))
	if err != nil {
		return BankTransferID{}, err
	}

	return BankTransferID(uid), nil
}
