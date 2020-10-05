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
 * Shopify business prefix + uuid will be used internally when passed between services and
 * removed when passing to/from the database.
 */
type ShopifyBusinessID uuid.UUID

func (id ShopifyBusinessID) Prefix() IDPrefix {
	return IDPrefixShopifyBusiness
}

func (id ShopifyBusinessID) String() string {
	return string(IDPrefixShopifyBusiness) + id.UUIDString()
}

func (id ShopifyBusinessID) UUIDString() string {
	return uuid.UUID(id).String()
}

func (id ShopifyBusinessID) IsZero() bool {
	return uuid.UUID(id) == uuid.Nil
}

func (id ShopifyBusinessID) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.String())
}

func (id *ShopifyBusinessID) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	uid, err := ParseShopifyBusinessID(s)
	if err != nil {
		return err
	}

	*id = uid
	return nil
}

// SQL value marshaller
func (id ShopifyBusinessID) Value() (driver.Value, error) {
	return id.UUIDString(), nil
}

// SQL scanner
func (id *ShopifyBusinessID) Scan(value interface{}) error {
	if value == nil {
		*id = ShopifyBusinessID{}
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

	*id = ShopifyBusinessID(uid)
	return nil
}

func NewShopifyBusinessID() (ShopifyBusinessID, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return ShopifyBusinessID{}, err
	}

	return ShopifyBusinessID(id), nil
}

func ParseShopifyBusinessID(id string) (ShopifyBusinessID, error) {
	// Return nil id on empty string
	if id == "" {
		return ShopifyBusinessID{}, nil
	}

	// Check prefix
	if !strings.HasPrefix(id, string(IDPrefixShopifyBusiness)) {
		return ShopifyBusinessID{}, errors.New("invalid shopify business id prefix")
	}

	// Validate UUID
	uid, err := uuid.Parse(strings.TrimPrefix(id, string(IDPrefixShopifyBusiness)))
	if err != nil {
		return ShopifyBusinessID{}, err
	}

	return ShopifyBusinessID(uid), nil
}
