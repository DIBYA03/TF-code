package id

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

/*
 * UUID will be used internally when passed between services and
 * removed when passing to/from the database.
 */
type UUID uuid.UUID

func (id UUID) Prefix() IDPrefix {
	return IDPrefixNone
}

func (id UUID) String() string {
	return id.UUIDString()
}

func (id UUID) UUIDString() string {
	return uuid.UUID(id).String()
}

func (id UUID) IsZero() bool {
	return uuid.UUID(id) == uuid.Nil
}

// Returns empty string when zero for omitempty
func (id UUID) JSONString() string {
	if id.IsZero() {
		return ""
	}

	return id.String()
}

func (id UUID) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.String())
}

func (id *UUID) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	uid, err := ParseUUID(s)
	if err != nil {
		return err
	}

	*id = uid
	return nil
}

// SQL value marshaller
func (id UUID) Value() (driver.Value, error) {
	return uuid.UUID(id).String(), nil
}

// SQL scanner
func (id *UUID) Scan(value interface{}) error {
	if value == nil {
		*id = UUID{}
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

	*id = UUID(uid)
	return nil
}

func NewUUID() (UUID, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return UUID{}, err
	}

	return UUID(id), nil
}

func ParseUUID(id string) (UUID, error) {
	// Return nil id on empty string
	if id == "" {
		return UUID{}, nil
	}

	// Validate UUID
	uid, err := uuid.Parse(id)
	if err != nil {
		return UUID{}, err
	}

	return UUID(uid), nil
}
