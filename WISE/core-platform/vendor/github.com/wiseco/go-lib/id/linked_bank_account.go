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
 * LinkedBankAccount prefix + uuid will be used internally when passed between services and
 * removed when passing to/from the database.
 */
type LinkedBankAccountID uuid.UUID

func (id LinkedBankAccountID) Prefix() IDPrefix {
	return IDPrefixLinkedBankAccount
}

func (id LinkedBankAccountID) String() string {
	return string(IDPrefixLinkedBankAccount) + id.UUIDString()
}

func (id LinkedBankAccountID) UUIDString() string {
	return uuid.UUID(id).String()
}

func (id LinkedBankAccountID) IsZero() bool {
	return uuid.UUID(id) == uuid.Nil
}

// Returns empty string when zero for omitempty
func (id LinkedBankAccountID) JSONString() string {
	if id.IsZero() {
		return ""
	}

	return id.String()
}

func (id LinkedBankAccountID) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.String())
}

func (id *LinkedBankAccountID) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	uid, err := ParseLinkedBankAccountID(s)
	if err != nil {
		return err
	}

	*id = uid
	return nil
}

// SQL value marshaller
func (id LinkedBankAccountID) Value() (driver.Value, error) {
	return id.UUIDString(), nil
}

// SQL scanner
func (id *LinkedBankAccountID) Scan(value interface{}) error {
	if value == nil {
		*id = LinkedBankAccountID{}
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

	*id = LinkedBankAccountID(uid)
	return nil
}

func ParseLinkedBankAccountID(id string) (LinkedBankAccountID, error) {
	// Return nil id on empty string
	if id == "" {
		return LinkedBankAccountID{}, nil
	}

	// Check prefix
	if !strings.HasPrefix(id, string(IDPrefixLinkedBankAccount)) {
		return LinkedBankAccountID{}, errors.New("invalid linked bank account id prefix")
	}

	// Validate UUID
	uid, err := uuid.Parse(strings.TrimPrefix(id, string(IDPrefixLinkedBankAccount)))
	if err != nil {
		return LinkedBankAccountID{}, err
	}

	return LinkedBankAccountID(uid), nil
}
