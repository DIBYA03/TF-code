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
 * Refresh token prefix + uuid will be used internally when passed between services and
 * removed when passing to/from the database.
 */
type RefreshTokenID uuid.UUID

func (id RefreshTokenID) Prefix() IDPrefix {
	return IDPrefixRefreshToken
}

func (id RefreshTokenID) String() string {
	return string(IDPrefixRefreshToken) + id.UUIDString()
}

func (id RefreshTokenID) UUIDString() string {
	return uuid.UUID(id).String()
}

func (id RefreshTokenID) IsZero() bool {
	return uuid.UUID(id) == uuid.Nil
}

// Returns empty string when zero for omitempty
func (id RefreshTokenID) JSONString() string {
	if id.IsZero() {
		return ""
	}

	return id.String()
}

func (id RefreshTokenID) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.String())
}

func (id *RefreshTokenID) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	uid, err := ParseRefreshTokenID(s)
	if err != nil {
		return err
	}

	*id = uid
	return nil
}

// SQL value marshaller
func (id RefreshTokenID) Value() (driver.Value, error) {
	return id.UUIDString(), nil
}

// SQL scanner
func (id *RefreshTokenID) Scan(value interface{}) error {
	if value == nil {
		*id = RefreshTokenID{}
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

	*id = RefreshTokenID(uid)
	return nil
}

func ParseRefreshTokenID(id string) (RefreshTokenID, error) {
	// Return nil id on empty string
	if id == "" {
		return RefreshTokenID{}, nil
	}

	// Check prefix
	if !strings.HasPrefix(id, string(IDPrefixRefreshToken)) {
		return RefreshTokenID{}, errors.New("invalid refresh token id prefix")
	}

	// Validate UUID
	uid, err := uuid.Parse(strings.TrimPrefix(id, string(IDPrefixRefreshToken)))
	if err != nil {
		return RefreshTokenID{}, err
	}

	return RefreshTokenID(uid), nil
}
