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
 * Image prefix + uuid will be used internally when passed between services and
 * removed when passing to/from the database.
 */
type ImageID uuid.UUID

func (id ImageID) Prefix() IDPrefix {
	return IDPrefixImage
}

func (id ImageID) String() string {
	return string(IDPrefixImage) + id.UUIDString()
}

func (id ImageID) UUIDString() string {
	return uuid.UUID(id).String()
}

func (id ImageID) IsZero() bool {
	return uuid.UUID(id) == uuid.Nil
}

// Returns empty string when zero for omitempty
func (id ImageID) JSONString() string {
	if id.IsZero() {
		return ""
	}

	return id.String()
}

func (id ImageID) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.String())
}

func (id *ImageID) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	uid, err := ParseImageID(s)
	if err != nil {
		return err
	}

	*id = uid
	return nil
}

// SQL value marshaller
func (id ImageID) Value() (driver.Value, error) {
	return id.UUIDString(), nil
}

// SQL scanner
func (id *ImageID) Scan(value interface{}) error {
	if value == nil {
		*id = ImageID{}
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

	*id = ImageID(uid)
	return nil
}

func ParseImageID(id string) (ImageID, error) {
	// Return nil id on empty string
	if id == "" {
		return ImageID{}, nil
	}

	// Check prefix
	if !strings.HasPrefix(id, string(IDPrefixImage)) {
		return ImageID{}, errors.New("invalid image id prefix")
	}

	// Validate UUID
	uid, err := uuid.Parse(strings.TrimPrefix(id, string(IDPrefixImage)))
	if err != nil {
		return ImageID{}, err
	}

	return ImageID(uid), nil
}
