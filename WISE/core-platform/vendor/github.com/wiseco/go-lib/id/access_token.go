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
 * AccessTokenID prefix + uuid will be used internally when passed between services and
 * removed when passing to/from the database.
 */
type AccessTokenID uuid.UUID

func (id AccessTokenID) Prefix() IDPrefix {
	return IDPrefixAccessToken
}

func (id AccessTokenID) String() string {
	return string(IDPrefixAccessToken) + id.UUIDString()
}

func (id AccessTokenID) UUIDString() string {
	return uuid.UUID(id).String()
}

func (id AccessTokenID) IsZero() bool {
	return uuid.UUID(id) == uuid.Nil
}

// Returns empty string when zero for omitempty
func (id AccessTokenID) JSONString() string {
	if id.IsZero() {
		return ""
	}

	return id.String()
}

func (id AccessTokenID) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.String())
}

func (id *AccessTokenID) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	uid, err := ParseAccessTokenID(s)
	if err != nil {
		return err
	}

	*id = uid
	return nil
}

// SQL value marshaller
func (id AccessTokenID) Value() (driver.Value, error) {
	return id.UUIDString(), nil
}

// SQL scanner
func (id *AccessTokenID) Scan(value interface{}) error {
	if value == nil {
		*id = AccessTokenID{}
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

	*id = AccessTokenID(uid)
	return nil
}

func ParseAccessTokenID(id string) (AccessTokenID, error) {
	// Return nil id on empty string
	if id == "" {
		return AccessTokenID{}, nil
	}

	// Check prefix
	if !strings.HasPrefix(id, string(IDPrefixAccessToken)) {
		return AccessTokenID{}, errors.New("invalid access token id prefix")
	}

	// Validate UUID
	uid, err := uuid.Parse(strings.TrimPrefix(id, string(IDPrefixAccessToken)))
	if err != nil {
		return AccessTokenID{}, err
	}

	return AccessTokenID(uid), nil
}
