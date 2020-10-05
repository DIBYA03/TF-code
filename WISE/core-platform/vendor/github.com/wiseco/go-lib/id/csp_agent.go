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
 * CspAgent prefix + uuid will be used internally when passed between services and
 * removed when passing to/from the database.
 */
type CspAgentID uuid.UUID

func (id CspAgentID) Prefix() IDPrefix {
	return IDPrefixCspAgent
}

func (id CspAgentID) String() string {
	return string(IDPrefixCspAgent) + id.UUIDString()
}

func (id CspAgentID) UUIDString() string {
	return uuid.UUID(id).String()
}

func (id CspAgentID) IsZero() bool {
	return uuid.UUID(id) == uuid.Nil
}

// Returns empty string when zero for omitempty
func (id CspAgentID) JSONString() string {
	if id.IsZero() {
		return ""
	}

	return id.String()
}

func (id CspAgentID) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.String())
}

func (id *CspAgentID) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	uid, err := ParseCspAgentID(s)
	if err != nil {
		return err
	}

	*id = uid
	return nil
}

// SQL value marshaller
func (id CspAgentID) Value() (driver.Value, error) {
	return id.UUIDString(), nil
}

// SQL scanner
func (id *CspAgentID) Scan(value interface{}) error {
	if value == nil {
		*id = CspAgentID{}
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

	*id = CspAgentID(uid)
	return nil
}

func ParseCspAgentID(id string) (CspAgentID, error) {
	// Return nil id on empty string
	if id == "" {
		return CspAgentID{}, nil
	}

	// Check prefix
	if !strings.HasPrefix(id, string(IDPrefixCspAgent)) {
		return CspAgentID{}, errors.New("invalid csp agent id prefix")
	}

	// Validate UUID
	uid, err := uuid.Parse(strings.TrimPrefix(id, string(IDPrefixCspAgent)))
	if err != nil {
		return CspAgentID{}, err
	}

	return CspAgentID(uid), nil
}
