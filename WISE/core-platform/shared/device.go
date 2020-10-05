package shared

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/google/uuid"
)

const UserDeviceIDPrefix = "udv-"

type UserDeviceID string

func (id *UserDeviceID) MarshalJSON() ([]byte, error) {
	if id == nil {
		return nil, errors.New("invalid id type")
	}

	return json.Marshal(UserDeviceIDPrefix + string(*id))
}

func (id *UserDeviceID) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	uid, err := uuid.Parse(strings.TrimPrefix(s, UserDeviceIDPrefix))
	if err != nil {
		return err
	}

	*id = UserDeviceID(uid.String())
	return nil
}

func ParseUserDeviceID(id string) (UserDeviceID, error) {
	uid, err := uuid.Parse(strings.TrimPrefix(id, UserDeviceIDPrefix))
	if err != nil {
		return UserDeviceID(""), err
	}

	return UserDeviceID(uid.String()), nil
}
