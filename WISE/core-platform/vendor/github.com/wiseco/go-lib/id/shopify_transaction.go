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
 * Shopify transaction prefix + uuid will be used internally when passed between services and
 * removed when passing to/from the database.
 */
type ShopifyTransactionID uuid.UUID

func (id ShopifyTransactionID) Prefix() IDPrefix {
	return IDPrefixShopifyTransaction
}

func (id ShopifyTransactionID) String() string {
	return string(IDPrefixShopifyTransaction) + id.UUIDString()
}

func (id ShopifyTransactionID) UUIDString() string {
	return uuid.UUID(id).String()
}

func (id ShopifyTransactionID) IsZero() bool {
	return uuid.UUID(id) == uuid.Nil
}

func (id ShopifyTransactionID) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.String())
}

func (id *ShopifyTransactionID) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	uid, err := ParseShopifyTransactionID(s)
	if err != nil {
		return err
	}

	*id = uid
	return nil
}

// SQL value marshaller
func (id ShopifyTransactionID) Value() (driver.Value, error) {
	return id.UUIDString(), nil
}

// SQL scanner
func (id *ShopifyTransactionID) Scan(value interface{}) error {
	if value == nil {
		*id = ShopifyTransactionID{}
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

	*id = ShopifyTransactionID(uid)
	return nil
}

func NewShopifyTransactionID() (ShopifyTransactionID, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return ShopifyTransactionID{}, err
	}

	return ShopifyTransactionID(id), nil
}

func ParseShopifyTransactionID(id string) (ShopifyTransactionID, error) {
	// Return nil id on empty string
	if id == "" {
		return ShopifyTransactionID{}, nil
	}

	// Check prefix
	if !strings.HasPrefix(id, string(IDPrefixShopifyTransaction)) {
		return ShopifyTransactionID{}, errors.New("invalid shopify transaction id prefix")
	}

	// Validate UUID
	uid, err := uuid.Parse(strings.TrimPrefix(id, string(IDPrefixShopifyTransaction)))
	if err != nil {
		return ShopifyTransactionID{}, err
	}

	return ShopifyTransactionID(uid), nil
}
