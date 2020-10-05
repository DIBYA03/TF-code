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
type ItemID uuid.UUID

func (id ItemID) Prefix() IDPrefix {
	return IDPrefixItem
}

func (id ItemID) String() string {
	return string(IDPrefixItem) + id.UUIDString()
}

func (id ItemID) UUIDString() string {
	return uuid.UUID(id).String()
}

func (id ItemID) IsZero() bool {
	return uuid.UUID(id) == uuid.Nil
}

// Returns empty string when zero for omitempty
func (id ItemID) JSONString() string {
	if id.IsZero() {
		return ""
	}

	return id.String()
}

func (id ItemID) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.String())
}

func (id *ItemID) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	uid, err := ParseItemID(s)
	if err != nil {
		return err
	}

	*id = uid
	return nil
}

func NewItemID() (ItemID, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return ItemID{}, err
	}

	return ItemID(id), nil
}


// SQL value marshaller
func (id ItemID) Value() (driver.Value, error) {
	return id.UUIDString(), nil
}

// SQL scanner
func (id *ItemID) Scan(value interface{}) error {
	if value == nil {
		*id = ItemID{}
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

	*id = ItemID(uid)
	return nil
}

func ParseItemID(id string) (ItemID, error) {
	// Return nil id on empty string
	if id == "" {
		return ItemID{}, nil
	}

	// Check prefix
	if !strings.HasPrefix(id, string(IDPrefixItem)) {
		return ItemID{}, errors.New("invalid item id prefix")
	}

	// Validate UUID
	uid, err := uuid.Parse(strings.TrimPrefix(id, string(IDPrefixItem)))
	if err != nil {
		return ItemID{}, err
	}

	return ItemID(uid), nil
}
