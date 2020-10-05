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
 * PartnerKey prefix + uuid will be used internally when passed between services and
 * removed when passing to/from the database.
 */
type PartnerKeyID uuid.UUID

func (id PartnerKeyID) Prefix() IDPrefix {
	return IDPrefixPartnerKey
}

func (id PartnerKeyID) String() string {
	return string(IDPrefixPartnerKey) + id.UUIDString()
}

func (id PartnerKeyID) UUIDString() string {
	return uuid.UUID(id).String()
}

func (id PartnerKeyID) IsZero() bool {
	return uuid.UUID(id) == uuid.Nil
}

// Returns empty string when zero for omitempty
func (id PartnerKeyID) JSONString() string {
	if id.IsZero() {
		return ""
	}

	return id.String()
}

func (id PartnerKeyID) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.String())
}

func (id *PartnerKeyID) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	uid, err := ParsePartnerKeyID(s)
	if err != nil {
		return err
	}

	*id = uid
	return nil
}

// SQL value marshaller
func (id PartnerKeyID) Value() (driver.Value, error) {
	return id.UUIDString(), nil
}

// SQL scanner
func (id *PartnerKeyID) Scan(value interface{}) error {
	if value == nil {
		*id = PartnerKeyID{}
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

	*id = PartnerKeyID(uid)
	return nil
}

func ParsePartnerKeyID(id string) (PartnerKeyID, error) {
	// Return nil id on empty string
	if id == "" {
		return PartnerKeyID{}, nil
	}

	// Check prefix
	if !strings.HasPrefix(id, string(IDPrefixPartnerKey)) {
		return PartnerKeyID{}, errors.New("invalid partner key id prefix")
	}

	// Validate UUID
	uid, err := uuid.Parse(strings.TrimPrefix(id, string(IDPrefixPartnerKey)))
	if err != nil {
		return PartnerKeyID{}, err
	}

	return PartnerKeyID(uid), nil
}
