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
 * PhoneVerification prefix + uuid will be used internally when passed between services and
 * removed when passing to/from the database.
 */
type PhoneOtpID uuid.UUID

func (id PhoneOtpID) Prefix() IDPrefix {
	return IDPrefixPhoneOTP
}

func (id PhoneOtpID) String() string {
	return string(IDPrefixPhoneOTP) + id.UUIDString()
}

func (id PhoneOtpID) UUIDString() string {
	return uuid.UUID(id).String()
}

func (id PhoneOtpID) IsZero() bool {
	return uuid.UUID(id) == uuid.Nil
}

// Returns empty string when zero for omitempty
func (id PhoneOtpID) JSONString() string {
	if id.IsZero() {
		return ""
	}

	return id.String()
}

func (id PhoneOtpID) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.String())
}

func (id *PhoneOtpID) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	uid, err := ParsePhoneOtpID(s)
	if err != nil {
		return err
	}

	*id = uid
	return nil
}

// SQL value marshaller
func (id PhoneOtpID) Value() (driver.Value, error) {
	return id.UUIDString(), nil
}

// SQL scanner
func (id *PhoneOtpID) Scan(value interface{}) error {
	if value == nil {
		*id = PhoneOtpID{}
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

	*id = PhoneOtpID(uid)
	return nil
}

func NewPhoneOtpID() (PhoneOtpID, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return PhoneOtpID{}, err
	}

	return PhoneOtpID(id), nil
}

func ParsePhoneOtpID(id string) (PhoneOtpID, error) {
	// Return nil id on empty string
	if id == "" {
		return PhoneOtpID{}, nil
	}

	// Check prefix
	if !strings.HasPrefix(id, string(IDPrefixPhoneOTP)) {
		return PhoneOtpID{}, errors.New("invalid phone verification id prefix")
	}

	// Validate UUID
	uid, err := uuid.Parse(strings.TrimPrefix(id, string(IDPrefixPhoneOTP)))
	if err != nil {
		return PhoneOtpID{}, err
	}

	return PhoneOtpID(uid), nil
}
