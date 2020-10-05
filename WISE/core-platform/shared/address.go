package shared

import (
	"encoding/json"
	"strings"

	"github.com/google/uuid"
)

const AddressIDPrefix = "adr-"

type AddressID string

func (id AddressID) MarshalJSON() ([]byte, error) {
	prefix := AddressIDPrefix

	if len(id) == 0 {
		prefix = ""
	}

	return json.Marshal(prefix + string(id))
}

func (id *AddressID) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	uid, err := uuid.Parse(strings.TrimPrefix(s, AddressIDPrefix))
	if err != nil {
		return err
	}

	*id = AddressID(uid.String())
	return nil
}

func ParseAddressID(id string) (AddressID, error) {
	uid, err := uuid.Parse(strings.TrimPrefix(id, AddressIDPrefix))
	if err != nil {
		return AddressID(""), err
	}

	return AddressID(uid.String()), nil
}

func (id *AddressID) ToPrefixString() string {
	if id == nil {
		return ""
	}

	return AddressIDPrefix + string(*id)
}
