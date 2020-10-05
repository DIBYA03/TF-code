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
 * ShopifyOrder prefix + uuid will be used internally when passed between services and
 * removed when passing to/from the database.
 */
type ShopifyOrderID uuid.UUID

func (id ShopifyOrderID) Prefix() IDPrefix {
	return IDPrefixShopifyOrder
}

func (id ShopifyOrderID) String() string {
	return string(IDPrefixShopifyOrder) + id.UUIDString()
}

func (id ShopifyOrderID) UUIDString() string {
	return uuid.UUID(id).String()
}

func (id ShopifyOrderID) IsZero() bool {
	return uuid.UUID(id) == uuid.Nil
}

func (id ShopifyOrderID) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.String())
}

func (id *ShopifyOrderID) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	uid, err := ParseShopifyOrderID(s)
	if err != nil {
		return err
	}

	*id = uid
	return nil
}

// SQL value marshaller
func (id ShopifyOrderID) Value() (driver.Value, error) {
	return id.UUIDString(), nil
}

// SQL scanner
func (id *ShopifyOrderID) Scan(value interface{}) error {
	if value == nil {
		*id = ShopifyOrderID{}
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

	*id = ShopifyOrderID(uid)
	return nil
}

func NewShopifyOrderID() (ShopifyOrderID, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return ShopifyOrderID{}, err
	}

	return ShopifyOrderID(id), nil
}

func ParseShopifyOrderID(id string) (ShopifyOrderID, error) {
	// Return nil id on empty string
	if id == "" {
		return ShopifyOrderID{}, nil
	}

	// Check prefix
	if !strings.HasPrefix(id, string(IDPrefixShopifyOrder)) {
		return ShopifyOrderID{}, errors.New("invalid shopify order id prefix")
	}

	// Validate UUID
	uid, err := uuid.Parse(strings.TrimPrefix(id, string(IDPrefixShopifyOrder)))
	if err != nil {
		return ShopifyOrderID{}, err
	}

	return ShopifyOrderID(uid), nil
}
