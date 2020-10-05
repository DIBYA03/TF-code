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
 * Event prefix + uuid will be used internally when passed between services and
 * removed when passing to/from the database.
 */
type EventID uuid.UUID

func (id EventID) Prefix() IDPrefix {
	return IDPrefixEvent
}

func (id EventID) String() string {
	return string(IDPrefixEvent) + id.UUIDString()
}

func (id EventID) UUIDString() string {
	return uuid.UUID(id).String()
}

func (id EventID) IsZero() bool {
	return uuid.UUID(id) == uuid.Nil
}

// Returns empty string when zero for omitempty
func (id EventID) JSONString() string {
	if id.IsZero() {
		return ""
	}

	return id.String()
}

func (id EventID) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.String())
}

func (id *EventID) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	uid, err := ParseEventID(s)
	if err != nil {
		return err
	}

	*id = uid
	return nil
}

// SQL value marshaller
func (id EventID) Value() (driver.Value, error) {
	return id.UUIDString(), nil
}

// SQL scanner
func (id *EventID) Scan(value interface{}) error {
	if value == nil {
		*id = EventID{}
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

	*id = EventID(uid)
	return nil
}

func NewEventID() (EventID, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return EventID{}, err
	}

	return EventID(id), nil
}

func ParseEventID(id string) (EventID, error) {
	// Return nil id on empty string
	if id == "" {
		return EventID{}, nil
	}

	// Check prefix
	if !strings.HasPrefix(id, string(IDPrefixEvent)) {
		return EventID{}, errors.New("invalid event id prefix")
	}

	// Validate UUID
	uid, err := uuid.Parse(strings.TrimPrefix(id, string(IDPrefixEvent)))
	if err != nil {
		return EventID{}, err
	}

	return EventID(uid), nil
}
