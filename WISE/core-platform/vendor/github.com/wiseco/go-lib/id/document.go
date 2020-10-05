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
 * Document prefix + uuid will be used internally when passed between services and
 * removed when passing to/from the database.
 */
type DocumentID uuid.UUID

func (id DocumentID) Prefix() IDPrefix {
	return IDPrefixDocument
}

func (id DocumentID) String() string {
	return string(IDPrefixDocument) + id.UUIDString()
}

func (id DocumentID) UUIDString() string {
	return uuid.UUID(id).String()
}

func (id DocumentID) IsZero() bool {
	return uuid.UUID(id) == uuid.Nil
}

// Returns empty string when zero for omitempty
func (id DocumentID) JSONString() string {
	if id.IsZero() {
		return ""
	}

	return id.String()
}

func (id DocumentID) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.String())
}

func (id *DocumentID) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	uid, err := ParseDocumentID(s)
	if err != nil {
		return err
	}

	*id = uid
	return nil
}

// SQL value marshaller
func (id DocumentID) Value() (driver.Value, error) {
	return id.UUIDString(), nil
}

// SQL scanner
func (id *DocumentID) Scan(value interface{}) error {
	if value == nil {
		*id = DocumentID{}
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

	*id = DocumentID(uid)
	return nil
}

func ParseDocumentID(id string) (DocumentID, error) {
	// Return nil id on empty string
	if id == "" {
		return DocumentID{}, nil
	}

	// Check prefix
	if !strings.HasPrefix(id, string(IDPrefixDocument)) {
		return DocumentID{}, errors.New("invalid document id prefix")
	}

	// Validate UUID
	uid, err := uuid.Parse(strings.TrimPrefix(id, string(IDPrefixDocument)))
	if err != nil {
		return DocumentID{}, err
	}

	return DocumentID(uid), nil
}
