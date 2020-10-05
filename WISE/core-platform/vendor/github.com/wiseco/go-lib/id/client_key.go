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
 * ClientKey prefix + uuid will be used internally when passed between services and
 * removed when passing to/from the database.
 */
type ClientKey uuid.UUID

func (id ClientKey) Prefix() IDPrefix {
	return IDPrefixClientKey
}

func (id ClientKey) String() string {
	return string(IDPrefixClientKey) + id.UUIDString()
}

func (id ClientKey) UUIDString() string {
	return uuid.UUID(id).String()
}

func (id ClientKey) IsZero() bool {
	return uuid.UUID(id) == uuid.Nil
}

// Returns empty string when zero for omitempty
func (id ClientKey) JSONString() string {
	if id.IsZero() {
		return ""
	}

	return id.String()
}

func (id ClientKey) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.String())
}

func (id *ClientKey) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	uid, err := ParseClientKey(s)
	if err != nil {
		return err
	}

	*id = uid
	return nil
}

// SQL value marshaller
func (id ClientKey) Value() (driver.Value, error) {
	return id.UUIDString(), nil
}

// SQL scanner
func (id *ClientKey) Scan(value interface{}) error {
	if value == nil {
		*id = ClientKey{}
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

	*id = ClientKey(uid)
	return nil
}

func ParseClientKey(id string) (ClientKey, error) {
	// Return nil id on empty string
	if id == "" {
		return ClientKey{}, nil
	}

	// Check prefix
	if !strings.HasPrefix(id, string(IDPrefixClientKey)) {
		return ClientKey{}, errors.New("invalid client key id prefix")
	}

	// Validate UUID
	uid, err := uuid.Parse(strings.TrimPrefix(id, string(IDPrefixClientKey)))
	if err != nil {
		return ClientKey{}, err
	}

	return ClientKey(uid), nil
}
