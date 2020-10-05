package shared

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/google/uuid"
)

const IdentityIDPrefix = "idt-"

type IdentityID string

func (id *IdentityID) MarshalJSON() ([]byte, error) {
	if id == nil {
		return nil, errors.New("invalid id type")
	}

	return json.Marshal(IdentityIDPrefix + string(*id))
}

func (id *IdentityID) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	uid, err := uuid.Parse(strings.TrimPrefix(s, IdentityIDPrefix))
	if err != nil {
		return err
	}

	*id = IdentityID(uid.String())
	return nil
}

func ParseIdentityID(id string) (IdentityID, error) {
	uid, err := uuid.Parse(strings.TrimPrefix(id, IdentityIDPrefix))
	if err != nil {
		return IdentityID(""), err
	}

	return IdentityID(uid.String()), nil
}
