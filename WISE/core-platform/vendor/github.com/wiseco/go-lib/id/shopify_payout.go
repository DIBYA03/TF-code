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
 * Shopify payout prefix + uuid will be used internally when passed between services and
 * removed when passing to/from the database.
 */
type ShopifyPayoutID uuid.UUID

func (id ShopifyPayoutID) Prefix() IDPrefix {
	return IDPrefixShopifyPayout
}

func (id ShopifyPayoutID) String() string {
	return string(IDPrefixShopifyPayout) + id.UUIDString()
}

func (id ShopifyPayoutID) UUIDString() string {
	return uuid.UUID(id).String()
}

func (id ShopifyPayoutID) IsZero() bool {
	return uuid.UUID(id) == uuid.Nil
}

func (id ShopifyPayoutID) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.String())
}

func (id *ShopifyPayoutID) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	uid, err := ParseShopifyPayoutID(s)
	if err != nil {
		return err
	}

	*id = uid
	return nil
}

// SQL value marshaller
func (id ShopifyPayoutID) Value() (driver.Value, error) {
	return id.UUIDString(), nil
}

// SQL scanner
func (id *ShopifyPayoutID) Scan(value interface{}) error {
	if value == nil {
		*id = ShopifyPayoutID{}
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

	*id = ShopifyPayoutID(uid)
	return nil
}

func NewShopifyPayoutID() (ShopifyPayoutID, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return ShopifyPayoutID{}, err
	}

	return ShopifyPayoutID(id), nil
}

func ParseShopifyPayoutID(id string) (ShopifyPayoutID, error) {
	// Return nil id on empty string
	if id == "" {
		return ShopifyPayoutID{}, nil
	}

	// Check prefix
	if !strings.HasPrefix(id, string(IDPrefixShopifyPayout)) {
		return ShopifyPayoutID{}, errors.New("invalid shopify payout id prefix")
	}

	// Validate UUID
	uid, err := uuid.Parse(strings.TrimPrefix(id, string(IDPrefixShopifyPayout)))
	if err != nil {
		return ShopifyPayoutID{}, err
	}

	return ShopifyPayoutID(uid), nil
}
