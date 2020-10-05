package shared

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

// Debit card block ID prefix
const BankCardBlockIDPrefix = "dcb-"

type BankCardBlockID string

func (id BankCardBlockID) MarshalJSON() ([]byte, error) {
	prefix := BankCardBlockIDPrefix

	if len(id) == 0 {
		prefix = ""
	}

	return json.Marshal(prefix + string(id))
}

func (id *BankCardBlockID) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	uid, err := uuid.Parse(strings.TrimPrefix(s, BankCardBlockIDPrefix))
	if err != nil {
		return err
	}

	*id = BankCardBlockID(uid.String())
	return nil
}

func (id *BankCardBlockID) Scan(value interface{}) error {
	if value == nil {
		*id = BankCardBlockID("")

		return nil
	}

	val, ok := value.([]byte)

	if !ok {
		return fmt.Errorf("Unable to scan value %v", value)
	}

	*id = BankCardBlockID(val)

	return nil
}

func ParseCardBlockID(id string) (BankCardBlockID, error) {
	uid, err := uuid.Parse(strings.TrimPrefix(id, BankCardBlockIDPrefix))
	if err != nil {
		return BankCardBlockID(""), err
	}

	return BankCardBlockID(uid.String()), nil
}
