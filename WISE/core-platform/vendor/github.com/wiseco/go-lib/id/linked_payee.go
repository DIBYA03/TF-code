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
 * LinkedPayee prefix + uuid will be used internally when passed between services and
 * removed when passing to/from the database.
 */
type LinkedPayeeID uuid.UUID

func (id LinkedPayeeID) Prefix() IDPrefix {
	return IDPrefixLinkedPayee
}

func (id LinkedPayeeID) String() string {
	return string(IDPrefixLinkedPayee) + id.UUIDString()
}

func (id LinkedPayeeID) UUIDString() string {
	return uuid.UUID(id).String()
}

func (id LinkedPayeeID) IsZero() bool {
	return uuid.UUID(id) == uuid.Nil
}

// Returns empty string when zero for omitempty
func (id LinkedPayeeID) JSONString() string {
	if id.IsZero() {
		return ""
	}

	return id.String()
}

func (id LinkedPayeeID) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.String())
}

func (id *LinkedPayeeID) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	uid, err := ParseLinkedPayeeID(s)
	if err != nil {
		return err
	}

	*id = uid
	return nil
}

// SQL value marshaller
func (id LinkedPayeeID) Value() (driver.Value, error) {
	return id.UUIDString(), nil
}

// SQL scanner
func (id *LinkedPayeeID) Scan(value interface{}) error {
	if value == nil {
		*id = LinkedPayeeID{}
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

	*id = LinkedPayeeID(uid)
	return nil
}

func ParseLinkedPayeeID(id string) (LinkedPayeeID, error) {
	// Return nil id on empty string
	if id == "" {
		return LinkedPayeeID{}, nil
	}

	// Check prefix
	if !strings.HasPrefix(id, string(IDPrefixLinkedPayee)) {
		return LinkedPayeeID{}, errors.New("invalid linked payee id prefix")
	}

	// Validate UUID
	uid, err := uuid.Parse(strings.TrimPrefix(id, string(IDPrefixLinkedPayee)))
	if err != nil {
		return LinkedPayeeID{}, err
	}

	return LinkedPayeeID(uid), nil
}
