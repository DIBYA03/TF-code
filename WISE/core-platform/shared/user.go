package shared

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/google/uuid"
)

const UserIDPrefix = "usr-"

type UserID string

const UserIDEmpty = UserID("00000000-0000-0000-0000-000000000000")

func (id *UserID) MarshalJSON() ([]byte, error) {
	if id == nil {
		return nil, errors.New("invalid id type")
	}

	return json.Marshal(UserIDPrefix + string(*id))
}

func (id *UserID) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	uid, err := uuid.Parse(strings.TrimPrefix(s, UserIDPrefix))
	if err != nil {
		return err
	}

	*id = UserID(uid.String())
	return nil
}

func ParseUserID(id string) (UserID, error) {
	uid, err := uuid.Parse(strings.TrimPrefix(id, UserIDPrefix))
	if err != nil {
		return UserID(""), err
	}

	return UserID(uid.String()), nil
}

func (id *UserID) ToPrefixString() string {
	if id == nil {
		return ""
	}

	return UserIDPrefix + string(*id)
}

const PartnerIDPrefix = "ptc-"

type PartnerID string

func (id *PartnerID) MarshalJSON() ([]byte, error) {
	if id == nil {
		return nil, errors.New("invalid id type")
	}

	return json.Marshal(PartnerIDPrefix + string(*id))
}

func (id *PartnerID) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	uid, err := uuid.Parse(strings.TrimPrefix(s, PartnerIDPrefix))
	if err != nil {
		return err
	}

	*id = PartnerID(uid.String())
	return nil
}

func ParsePartnerID(id string) (PartnerID, error) {
	uid, err := uuid.Parse(strings.TrimPrefix(id, PartnerIDPrefix))
	if err != nil {
		return PartnerID(""), err
	}

	return PartnerID(uid.String()), nil
}
