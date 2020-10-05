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
 * CardReader prefix + uuid will be used internally when passed between services and
 * removed when passing to/from the database.
 */
type CardReaderID uuid.UUID

func (id CardReaderID) Prefix() IDPrefix {
	return IDPrefixCardReader
}

func (id CardReaderID) String() string {
	return string(IDPrefixCardReader) + id.UUIDString()
}

func (id CardReaderID) UUIDString() string {
	return uuid.UUID(id).String()
}

func (id CardReaderID) IsZero() bool {
	return uuid.UUID(id) == uuid.Nil
}

// Returns empty string when zero for omitempty
func (id CardReaderID) JSONString() string {
	if id.IsZero() {
		return ""
	}

	return id.String()
}

func (id CardReaderID) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.String())
}

func (id *CardReaderID) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	uid, err := ParseCardReaderID(s)
	if err != nil {
		return err
	}

	*id = uid
	return nil
}

// SQL value marshaller
func (id CardReaderID) Value() (driver.Value, error) {
	return id.UUIDString(), nil
}

// SQL scanner
func (id *CardReaderID) Scan(value interface{}) error {
	if value == nil {
		*id = CardReaderID{}
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

	*id = CardReaderID(uid)
	return nil
}

func ParseCardReaderID(id string) (CardReaderID, error) {
	// Return nil id on empty string
	if id == "" {
		return CardReaderID{}, nil
	}

	// Check prefix
	if !strings.HasPrefix(id, string(IDPrefixCardReader)) {
		return CardReaderID{}, errors.New("invalid card reader id prefix")
	}

	// Validate UUID
	uid, err := uuid.Parse(strings.TrimPrefix(id, string(IDPrefixCardReader)))
	if err != nil {
		return CardReaderID{}, err
	}

	return CardReaderID(uid), nil
}
