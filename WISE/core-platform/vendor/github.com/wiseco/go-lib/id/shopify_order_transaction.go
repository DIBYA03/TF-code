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
 * ShopifyOrderTransaction prefix + uuid will be used internally when passed between services and
 * removed when passing to/from the database.
 */
type ShopifyOrderTransactionID uuid.UUID

func (id ShopifyOrderTransactionID) Prefix() IDPrefix {
	return IDPrefixShopifyOrderTransaction
}

func (id ShopifyOrderTransactionID) String() string {
	return string(IDPrefixShopifyOrderTransaction) + id.UUIDString()
}

func (id ShopifyOrderTransactionID) UUIDString() string {
	return uuid.UUID(id).String()
}

func (id ShopifyOrderTransactionID) IsZero() bool {
	return uuid.UUID(id) == uuid.Nil
}

func (id ShopifyOrderTransactionID) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.String())
}

func (id *ShopifyOrderTransactionID) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	uid, err := ParseShopifyOrderTransactionID(s)
	if err != nil {
		return err
	}

	*id = uid
	return nil
}

// SQL value marshaller
func (id ShopifyOrderTransactionID) Value() (driver.Value, error) {
	return id.UUIDString(), nil
}

// SQL scanner
func (id *ShopifyOrderTransactionID) Scan(value interface{}) error {
	if value == nil {
		*id = ShopifyOrderTransactionID{}
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

	*id = ShopifyOrderTransactionID(uid)
	return nil
}

func NewShopifyOrderTransactionID() (ShopifyOrderTransactionID, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return ShopifyOrderTransactionID{}, err
	}

	return ShopifyOrderTransactionID(id), nil
}

func ParseShopifyOrderTransactionID(id string) (ShopifyOrderTransactionID, error) {
	// Return nil id on empty string
	if id == "" {
		return ShopifyOrderTransactionID{}, nil
	}

	// Check prefix
	if !strings.HasPrefix(id, string(IDPrefixShopifyOrderTransaction)) {
		return ShopifyOrderTransactionID{}, errors.New("invalid shopify order id prefix")
	}

	// Validate UUID
	uid, err := uuid.Parse(strings.TrimPrefix(id, string(IDPrefixShopifyOrderTransaction)))
	if err != nil {
		return ShopifyOrderTransactionID{}, err
	}

	return ShopifyOrderTransactionID(uid), nil
}
