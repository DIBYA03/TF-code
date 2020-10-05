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
 * Invoice prefix + uuid will be used internally when passed between services and
 * removed when passing to/from the database.
 */
type InvoiceID uuid.UUID

func (id InvoiceID) Prefix() IDPrefix {
	return IDPrefixInvoice
}

func (id InvoiceID) String() string {
	return string(IDPrefixInvoice) + id.UUIDString()
}

func (id InvoiceID) UUIDString() string {
	return uuid.UUID(id).String()
}

func (id InvoiceID) IsZero() bool {
	return uuid.UUID(id) == uuid.Nil
}

// Returns empty string when zero for omitempty
func (id InvoiceID) JSONString() string {
	if id.IsZero() {
		return ""
	}

	return id.String()
}

func (id InvoiceID) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.String())
}

func (id *InvoiceID) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	uid, err := ParseInvoiceID(s)
	if err != nil {
		return err
	}

	*id = uid
	return nil
}

func NewInvoiceID() (InvoiceID, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return InvoiceID{}, err
	}

	return InvoiceID(id), nil
}


// SQL value marshaller
func (id InvoiceID) Value() (driver.Value, error) {
	return id.UUIDString(), nil
}

// SQL scanner
func (id *InvoiceID) Scan(value interface{}) error {
	if value == nil {
		*id = InvoiceID{}
		return nil
	}

	val, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("unable to scan value %v", value)
	}

	uid, err := uuid.Parse(string(val))
	if err != nil {
		return err
	}

	*id = InvoiceID(uid)
	return nil
}

func ParseInvoiceID(id string) (InvoiceID, error) {
	// Return nil id on empty string
	if id == "" {
		return InvoiceID{}, nil
	}

	// Check prefix
	if !strings.HasPrefix(id, string(IDPrefixInvoice)) {
		return InvoiceID{}, errors.New("invalid invoice id prefix")
	}

	// Validate UUID
	uid, err := uuid.Parse(strings.TrimPrefix(id, string(IDPrefixInvoice)))
	if err != nil {
		return InvoiceID{}, err
	}

	return InvoiceID(uid), nil
}
