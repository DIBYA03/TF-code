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
 * Request prefix + uuid will be used internally when passed between services and
 * removed when passing to/from the database.
 */
type RequestID uuid.UUID

func (id RequestID) Prefix() IDPrefix {
	return IDPrefixRequest
}

func (id RequestID) String() string {
	return string(IDPrefixRequest) + id.UUIDString()
}

func (id RequestID) UUIDString() string {
	return uuid.UUID(id).String()
}

func (id RequestID) IsZero() bool {
	return uuid.UUID(id) == uuid.Nil
}

// Returns empty string when zero for omitempty
func (id RequestID) JSONString() string {
	if id.IsZero() {
		return ""
	}

	return id.String()
}

func (id RequestID) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.String())
}

func (id *RequestID) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	uid, err := ParseRequestID(s)
	if err != nil {
		return err
	}

	*id = uid
	return nil
}

// SQL value marshaller
func (id RequestID) Value() (driver.Value, error) {
	return id.UUIDString(), nil
}

// SQL scanner
func (id *RequestID) Scan(value interface{}) error {
	if value == nil {
		*id = RequestID{}
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

	*id = RequestID(uid)
	return nil
}

func NewRequestID() (RequestID, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return RequestID{}, err
	}

	return RequestID(id), nil
}

func ParseRequestID(id string) (RequestID, error) {
	// Return nil id on empty string
	if id == "" {
		return RequestID{}, nil
	}

	// Check prefix
	if !strings.HasPrefix(id, string(IDPrefixRequest)) {
		return RequestID{}, errors.New("invalid request id prefix")
	}

	// Validate UUID
	uid, err := uuid.Parse(strings.TrimPrefix(id, string(IDPrefixRequest)))
	if err != nil {
		return RequestID{}, err
	}

	return RequestID(uid), nil
}
