package shared

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

const ContactIDPrefix = ""

type ContactID string

func (id ContactID) MarshalJSON() ([]byte, error) {
	prefix := ContactIDPrefix

	if len(id) == 0 {
		prefix = ""
	}

	return json.Marshal(prefix + string(id))
}

func (id *ContactID) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	uid, err := uuid.Parse(strings.TrimPrefix(s, ContactIDPrefix))
	if err != nil {
		return err
	}

	*id = ContactID(uid.String())
	return nil
}

func (id *ContactID) Scan(value interface{}) error {
	if value == nil {
		*id = ContactID("")

		return nil
	}

	val, ok := value.([]byte)

	if !ok {
		return fmt.Errorf("Unable to scan value %v", value)
	}

	*id = ContactID(val)

	return nil
}

func ParseContactID(id string) (ContactID, error) {
	uid, err := uuid.Parse(strings.TrimPrefix(id, ContactIDPrefix))
	if err != nil {
		return ContactID(""), err
	}

	return ContactID(uid.String()), nil
}

func (id *ContactID) ToPrefixString() string {
	if id == nil {
		return ""
	}

	return ContactIDPrefix + string(*id)
}
