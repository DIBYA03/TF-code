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
 * LinkedCard prefix + uuid will be used internally when passed between services and
 * removed when passing to/from the database.
 */
type LinkedCardID uuid.UUID

func (id LinkedCardID) Prefix() IDPrefix {
	return IDPrefixLinkedCard
}

func (id LinkedCardID) String() string {
	return string(IDPrefixLinkedCard) + id.UUIDString()
}

func (id LinkedCardID) UUIDString() string {
	return uuid.UUID(id).String()
}

func (id LinkedCardID) IsZero() bool {
	return uuid.UUID(id) == uuid.Nil
}

// Returns empty string when zero for omitempty
func (id LinkedCardID) JSONString() string {
	if id.IsZero() {
		return ""
	}

	return id.String()
}

func (id LinkedCardID) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.String())
}

func (id *LinkedCardID) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	uid, err := ParseLinkedCardID(s)
	if err != nil {
		return err
	}

	*id = uid
	return nil
}

// SQL value marshaller
func (id LinkedCardID) Value() (driver.Value, error) {
	return id.UUIDString(), nil
}

// SQL scanner
func (id *LinkedCardID) Scan(value interface{}) error {
	if value == nil {
		*id = LinkedCardID{}
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

	*id = LinkedCardID(uid)
	return nil
}

func ParseLinkedCardID(id string) (LinkedCardID, error) {
	// Return nil id on empty string
	if id == "" {
		return LinkedCardID{}, nil
	}

	// Check prefix
	if !strings.HasPrefix(id, string(IDPrefixLinkedCard)) {
		return LinkedCardID{}, errors.New("invalid linked card id prefix")
	}

	// Validate UUID
	uid, err := uuid.Parse(strings.TrimPrefix(id, string(IDPrefixLinkedCard)))
	if err != nil {
		return LinkedCardID{}, err
	}

	return LinkedCardID(uid), nil
}
