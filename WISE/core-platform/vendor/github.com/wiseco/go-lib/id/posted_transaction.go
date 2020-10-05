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
 * PostedTransaction prefix + uuid will be used internally when passed between services and
 * removed when passing to/from the database.
 */
type PostedTransactionID uuid.UUID

func (id PostedTransactionID) Prefix() IDPrefix {
	return IDPrefixPostedTransaction
}

func (id PostedTransactionID) String() string {
	return string(IDPrefixPostedTransaction) + id.UUIDString()
}

func (id PostedTransactionID) UUIDString() string {
	return uuid.UUID(id).String()
}

func (id PostedTransactionID) IsZero() bool {
	return uuid.UUID(id) == uuid.Nil
}

// Returns empty string when zero for omitempty
func (id PostedTransactionID) JSONString() string {
	if id.IsZero() {
		return ""
	}

	return id.String()
}

func (id *PostedTransactionID) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.String())
}

func (id *PostedTransactionID) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	uid, err := ParsePostedTransactionID(s)
	if err != nil {
		return err
	}

	*id = PostedTransactionID(uid)
	return nil
}

// SQL value marshaller
func (id PostedTransactionID) Value() (driver.Value, error) {
	return id.UUIDString(), nil
}

// SQL scanner
func (id *PostedTransactionID) Scan(value interface{}) error {
	if value == nil {
		*id = PostedTransactionID{}
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

	*id = PostedTransactionID(uid)
	return nil
}

func ParsePostedTransactionID(id string) (PostedTransactionID, error) {
	// Return nil id on empty string
	if id == "" {
		return PostedTransactionID{}, nil
	}

	// Check prefix
	if !strings.HasPrefix(id, string(IDPrefixPostedTransaction)) {
		return PostedTransactionID{}, errors.New("invalid posted transaction id prefix")
	}

	// Validate UUID
	uid, err := uuid.Parse(strings.TrimPrefix(id, string(IDPrefixPostedTransaction)))
	if err != nil {
		return PostedTransactionID{}, err
	}

	return PostedTransactionID(uid), nil
}
