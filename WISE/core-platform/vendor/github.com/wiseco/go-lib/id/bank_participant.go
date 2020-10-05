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
 * Participant prefix + uuid will be used internally when passed between services and
 * removed when passing to/from the database.
 */
type BankParticipantID uuid.UUID

func (id BankParticipantID) Prefix() IDPrefix {
	return IDPrefixParticipant
}

func (id BankParticipantID) String() string {
	return string(IDPrefixParticipant) + id.UUIDString()
}

func (id BankParticipantID) UUIDString() string {
	return uuid.UUID(id).String()
}

func (id BankParticipantID) IsZero() bool {
	return uuid.UUID(id) == uuid.Nil
}

// Returns empty string when zero for omitempty
func (id BankParticipantID) JSONString() string {
	if id.IsZero() {
		return ""
	}

	return id.String()
}

func (id BankParticipantID) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.String())
}

func (id *BankParticipantID) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	uid, err := ParseBankParticipantID(s)
	if err != nil {
		return err
	}

	*id = uid
	return nil
}

// SQL value marshaller
func (id BankParticipantID) Value() (driver.Value, error) {
	return id.UUIDString(), nil
}

// SQL scanner
func (id *BankParticipantID) Scan(value interface{}) error {
	if value == nil {
		*id = BankParticipantID{}
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

	*id = BankParticipantID(uid)
	return nil
}

func NewBankParticipantID() (BankParticipantID, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return BankParticipantID{}, err
	}

	return BankParticipantID(id), nil
}

func ParseBankParticipantID(id string) (BankParticipantID, error) {
	// Return nil id on empty string
	if id == "" {
		return BankParticipantID{}, nil
	}

	// Check prefix
	if !strings.HasPrefix(id, string(IDPrefixParticipant)) {
		return BankParticipantID{}, errors.New("invalid event id prefix")
	}

	// Validate UUID
	uid, err := uuid.Parse(strings.TrimPrefix(id, string(IDPrefixParticipant)))
	if err != nil {
		return BankParticipantID{}, err
	}

	return BankParticipantID(uid), nil
}
