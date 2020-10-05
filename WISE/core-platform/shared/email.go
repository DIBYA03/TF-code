package shared

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

const EmailIDPrefix = "eml-"

type EmailID string

func (id EmailID) MarshalJSON() ([]byte, error) {
	prefix := EmailIDPrefix

	if len(id) == 0 {
		prefix = ""
	}

	return json.Marshal(prefix + string(id))
}

func (id *EmailID) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	uid, err := uuid.Parse(strings.TrimPrefix(s, EmailIDPrefix))
	if err != nil {
		return err
	}

	*id = EmailID(uid.String())
	return nil
}

func ParseEmailID(id string) (EmailID, error) {
	uid, err := uuid.Parse(strings.TrimPrefix(id, EmailIDPrefix))
	if err != nil {
		return EmailID(""), err
	}

	return EmailID(uid.String()), nil
}

func (id *EmailID) ToPrefixString() string {
	if id == nil {
		return ""
	}

	return EmailIDPrefix + string(*id)
}

func (id *EmailID) Scan(value interface{}) error {
	if value == nil {
		*id = EmailID("")

		return nil
	}

	val, ok := value.([]byte)

	if !ok {
		return fmt.Errorf("Unable to scan value %v", value)
	}

	*id = EmailID(val)

	return nil
}
