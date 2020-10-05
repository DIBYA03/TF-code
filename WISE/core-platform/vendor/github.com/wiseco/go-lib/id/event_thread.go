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
 * EventThread prefix + uuid will be used internally when passed between services and
 * removed when passing to/from the database.
 */
type EventThreadID uuid.UUID

func (id EventThreadID) Prefix() IDPrefix {
	return IDPrefixEventThread
}

func (id EventThreadID) String() string {
	return string(IDPrefixEventThread) + id.UUIDString()
}

func (id EventThreadID) UUIDString() string {
	return uuid.UUID(id).String()
}

func (id EventThreadID) IsZero() bool {
	return uuid.UUID(id) == uuid.Nil
}

// Returns empty string when zero for omitempty
func (id EventThreadID) JSONString() string {
	if id.IsZero() {
		return ""
	}

	return id.String()
}

func (id EventThreadID) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.String())
}

func (id *EventThreadID) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	uid, err := ParseEventThreadID(s)
	if err != nil {
		return err
	}

	*id = uid
	return nil
}

// SQL value marshaller
func (id EventThreadID) Value() (driver.Value, error) {
	return id.UUIDString(), nil
}

// SQL scanner
func (id *EventThreadID) Scan(value interface{}) error {
	if value == nil {
		*id = EventThreadID{}
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

	*id = EventThreadID(uid)
	return nil
}

func NewEventThreadID() (EventThreadID, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return EventThreadID{}, err
	}

	return EventThreadID(id), nil
}

func ParseEventThreadID(id string) (EventThreadID, error) {
	// Return nil id on empty string
	if id == "" {
		return EventThreadID{}, nil
	}

	// Check prefix
	if !strings.HasPrefix(id, string(IDPrefixEventThread)) {
		return EventThreadID{}, errors.New("invalid event thread id prefix")
	}

	// Validate UUID
	uid, err := uuid.Parse(strings.TrimPrefix(id, string(IDPrefixEventThread)))
	if err != nil {
		return EventThreadID{}, err
	}

	return EventThreadID(uid), nil
}
