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
 * PaymentRequest prefix + uuid will be used internally when passed between services and
 * removed when passing to/from the database.
 */
type PaymentRequestID uuid.UUID

func (id PaymentRequestID) Prefix() IDPrefix {
	return IDPrefixPaymentRequest
}

func (id PaymentRequestID) String() string {
	return string(IDPrefixPaymentRequest) + id.UUIDString()
}

func (id PaymentRequestID) UUIDString() string {
	return uuid.UUID(id).String()
}

func (id PaymentRequestID) IsZero() bool {
	return uuid.UUID(id) == uuid.Nil
}

// Returns empty string when zero for omitempty
func (id PaymentRequestID) JSONString() string {
	if id.IsZero() {
		return ""
	}

	return id.String()
}

func (id PaymentRequestID) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.String())
}

func (id *PaymentRequestID) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	uid, err := ParsePaymentRequestID(s)
	if err != nil {
		return err
	}

	*id = uid
	return nil
}

// SQL value marshaller
func (id PaymentRequestID) Value() (driver.Value, error) {
	return id.UUIDString(), nil
}

// SQL scanner
func (id *PaymentRequestID) Scan(value interface{}) error {
	if value == nil {
		*id = PaymentRequestID{}
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

	*id = PaymentRequestID(uid)
	return nil
}

func ParsePaymentRequestID(id string) (PaymentRequestID, error) {
	// Return nil id on empty string
	if id == "" {
		return PaymentRequestID{}, nil
	}

	// Check prefix
	if !strings.HasPrefix(id, string(IDPrefixPaymentRequest)) {
		return PaymentRequestID{}, errors.New("invalid payment request id prefix")
	}

	// Validate UUID
	uid, err := uuid.Parse(strings.TrimPrefix(id, string(IDPrefixPaymentRequest)))
	if err != nil {
		return PaymentRequestID{}, err
	}

	return PaymentRequestID(uid), nil
}
