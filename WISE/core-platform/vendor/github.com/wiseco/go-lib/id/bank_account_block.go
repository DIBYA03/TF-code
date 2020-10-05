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
 * IDPrefixBankAccountBlock prefix + uuid will be used internally when passed between services and
 * removed when passing to/from the database.
 */
type BankAccountBlockID uuid.UUID

func (id BankAccountBlockID) Prefix() IDPrefix {
	return IDPrefixBankAccountBlock
}

func (id BankAccountBlockID) String() string {
	return string(IDPrefixBankAccountBlock) + id.UUIDString()
}

func (id BankAccountBlockID) UUIDString() string {
	return uuid.UUID(id).String()
}

func (id BankAccountBlockID) IsZero() bool {
	return uuid.UUID(id) == uuid.Nil
}

// Returns empty string when zero for omitempty
func (id BankAccountBlockID) JSONString() string {
	if id.IsZero() {
		return ""
	}

	return id.String()
}

func (id BankAccountBlockID) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.String())
}

func (id *BankAccountBlockID) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	uid, err := ParseBankAccountBlockID(s)
	if err != nil {
		return err
	}

	*id = uid
	return nil
}

// SQL value marshaller
func (id BankAccountBlockID) Value() (driver.Value, error) {
	return id.UUIDString(), nil
}

// SQL scanner
func (id *BankAccountBlockID) Scan(value interface{}) error {
	if value == nil {
		*id = BankAccountBlockID{}
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

	*id = BankAccountBlockID(uid)
	return nil
}

func NewBankAccountBlockID() (BankAccountBlockID, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return BankAccountBlockID{}, err
	}

	return BankAccountBlockID(id), nil
}

func ParseBankAccountBlockID(id string) (BankAccountBlockID, error) {
	// Return nil id on empty string
	if id == "" {
		return BankAccountBlockID{}, nil
	}

	// Check prefix
	if !strings.HasPrefix(id, string(IDPrefixBankAccountBlock)) {
		return BankAccountBlockID{}, errors.New("invalid event id prefix")
	}

	// Validate UUID
	uid, err := uuid.Parse(strings.TrimPrefix(id, string(IDPrefixBankAccountBlock)))
	if err != nil {
		return BankAccountBlockID{}, err
	}

	return BankAccountBlockID(uid), nil
}
