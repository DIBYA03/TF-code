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
 * Partner prefix + uuid will be used internally when passed between services and
 * removed when passing to/from the database.
 */
type PartnerID uuid.UUID

func (id PartnerID) Prefix() IDPrefix {
	return IDPrefixPartner
}

func (id PartnerID) String() string {
	return string(IDPrefixPartner) + id.UUIDString()
}

func (id PartnerID) UUIDString() string {
	return uuid.UUID(id).String()
}

func (id PartnerID) IsZero() bool {
	return uuid.UUID(id) == uuid.Nil
}

// Returns empty string when zero for omitempty
func (id PartnerID) JSONString() string {
	if id.IsZero() {
		return ""
	}

	return id.String()
}

func (id PartnerID) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.String())
}

func (id *PartnerID) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	uid, err := ParsePartnerID(s)
	if err != nil {
		return err
	}

	*id = uid
	return nil
}

// SQL value marshaller
func (id PartnerID) Value() (driver.Value, error) {
	return id.UUIDString(), nil
}

// SQL scanner
func (id *PartnerID) Scan(value interface{}) error {
	if value == nil {
		*id = PartnerID{}
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

	*id = PartnerID(uid)
	return nil
}

func ParsePartnerID(id string) (PartnerID, error) {
	// Return nil id on empty string
	if id == "" {
		return PartnerID{}, nil
	}

	// Check prefix
	if !strings.HasPrefix(id, string(IDPrefixPartner)) {
		return PartnerID{}, errors.New("invalid partner id prefix")
	}

	// Validate UUID
	uid, err := uuid.Parse(strings.TrimPrefix(id, string(IDPrefixPartner)))
	if err != nil {
		return PartnerID{}, err
	}

	return PartnerID(uid), nil
}
