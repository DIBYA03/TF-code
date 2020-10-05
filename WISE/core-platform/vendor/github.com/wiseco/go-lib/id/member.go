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
 * Member prefix + uuid will be used internally when passed between services and
 * removed when passing to/from the database.
 */
type MemberID uuid.UUID

func (id MemberID) Prefix() IDPrefix {
	return IDPrefixMember
}

func (id MemberID) String() string {
	return string(IDPrefixMember) + id.UUIDString()
}

func (id MemberID) UUIDString() string {
	return uuid.UUID(id).String()
}

func (id MemberID) IsZero() bool {
	return uuid.UUID(id) == uuid.Nil
}

// Returns empty string when zero for omitempty
func (id MemberID) JSONString() string {
	if id.IsZero() {
		return ""
	}

	return id.String()
}

func (id MemberID) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.String())
}

func (id *MemberID) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	uid, err := ParseMemberID(s)
	if err != nil {
		return err
	}

	*id = uid
	return nil
}

// SQL value marshaller
func (id MemberID) Value() (driver.Value, error) {
	return id.UUIDString(), nil
}

// SQL scanner
func (id *MemberID) Scan(value interface{}) error {
	if value == nil {
		*id = MemberID{}
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

	*id = MemberID(uid)
	return nil
}

func ParseMemberID(id string) (MemberID, error) {
	// Return nil id on empty string
	if id == "" {
		return MemberID{}, nil
	}

	// Check prefix
	if !strings.HasPrefix(id, string(IDPrefixMember)) {
		return MemberID{}, errors.New("invalid member id prefix")
	}

	// Validate UUID
	uid, err := uuid.Parse(strings.TrimPrefix(id, string(IDPrefixMember)))
	if err != nil {
		return MemberID{}, err
	}

	return MemberID(uid), nil
}
