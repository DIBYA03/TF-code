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
 * SecretKey prefix + uuid will be used internally when passed between services and
 * removed when passing to/from the database.
 */
type SecretKey uuid.UUID

func (id SecretKey) Prefix() IDPrefix {
	return IDPrefixSecretKey
}

func (id SecretKey) String() string {
	return string(IDPrefixSecretKey) + id.UUIDString()
}

func (id SecretKey) UUIDString() string {
	return uuid.UUID(id).String()
}

func (id SecretKey) IsZero() bool {
	return uuid.UUID(id) == uuid.Nil
}

// Returns empty string when zero for omitempty
func (id SecretKey) JSONString() string {
	if id.IsZero() {
		return ""
	}

	return id.String()
}

func (id SecretKey) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.String())
}

func (id *SecretKey) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	uid, err := ParseSecretKey(s)
	if err != nil {
		return err
	}

	*id = uid
	return nil
}

// SQL value marshaller
func (id SecretKey) Value() (driver.Value, error) {
	return id.UUIDString(), nil
}

// SQL scanner
func (id *SecretKey) Scan(value interface{}) error {
	if value == nil {
		*id = SecretKey{}
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

	*id = SecretKey(uid)
	return nil
}

func ParseSecretKey(id string) (SecretKey, error) {
	// Return nil id on empty string
	if id == "" {
		return SecretKey{}, nil
	}

	// Check prefix
	if !strings.HasPrefix(id, string(IDPrefixSecretKey)) {
		return SecretKey{}, errors.New("invalid secret key prefix")
	}

	// Validate UUID
	uid, err := uuid.Parse(strings.TrimPrefix(id, string(IDPrefixSecretKey)))
	if err != nil {
		return SecretKey{}, err
	}

	return SecretKey(uid), nil
}
