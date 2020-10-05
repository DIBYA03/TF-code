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
 * Consumer prefix + uuid will be used internally when passed between services and
 * removed when passing to/from the database.
 */
type ConsumerID uuid.UUID

func (id ConsumerID) Prefix() IDPrefix {
	return IDPrefixConsumer
}

func (id ConsumerID) String() string {
	return string(IDPrefixConsumer) + id.UUIDString()
}

func (id ConsumerID) UUIDString() string {
	return uuid.UUID(id).String()
}

func (id ConsumerID) IsZero() bool {
	return uuid.UUID(id) == uuid.Nil
}

// Returns empty string when zero for omitempty
func (id ConsumerID) JSONString() string {
	if id.IsZero() {
		return ""
	}

	return id.String()
}

func (id ConsumerID) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.String())
}

func (id *ConsumerID) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	uid, err := ParseConsumerID(s)
	if err != nil {
		return err
	}

	*id = uid
	return nil
}

// SQL value marshaller
func (id ConsumerID) Value() (driver.Value, error) {
	return id.UUIDString(), nil
}

// SQL scanner
func (id *ConsumerID) Scan(value interface{}) error {
	if value == nil {
		*id = ConsumerID{}
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

	*id = ConsumerID(uid)
	return nil
}

func ParseConsumerID(id string) (ConsumerID, error) {
	// Return nil id on empty string
	if id == "" {
		return ConsumerID{}, nil
	}

	// Check prefix
	if !strings.HasPrefix(id, string(IDPrefixConsumer)) {
		return ConsumerID{}, errors.New("invalid consumer id prefix")
	}

	// Validate UUID
	uid, err := uuid.Parse(strings.TrimPrefix(id, string(IDPrefixConsumer)))
	if err != nil {
		return ConsumerID{}, err
	}

	return ConsumerID(uid), nil
}
