package shared

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/google/uuid"
)

const BusinessDocumentIDPrefix = "bdc-"

type BusinessDocumentID string

func (id *BusinessDocumentID) MarshalJSON() ([]byte, error) {
	if id == nil {
		return nil, errors.New("invalid id type")
	}

	return json.Marshal(BusinessDocumentIDPrefix + string(*id))
}

func (id *BusinessDocumentID) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	uid, err := uuid.Parse(strings.TrimPrefix(s, BusinessDocumentIDPrefix))
	if err != nil {
		return err
	}

	*id = BusinessDocumentID(uid.String())
	return nil
}

func ParseBusinessDocumentID(id string) (BusinessDocumentID, error) {
	uid, err := uuid.Parse(strings.TrimPrefix(id, BusinessDocumentIDPrefix))
	if err != nil {
		return BusinessDocumentID(""), err
	}

	return BusinessDocumentID(uid.String()), nil
}

const UserDocumentIDPrefix = "udc-"

type UserDocumentID string

func (id *UserDocumentID) MarshalJSON() ([]byte, error) {
	if id == nil {
		return nil, errors.New("invalid id type")
	}

	return json.Marshal(UserDocumentIDPrefix + string(*id))
}

func (id *UserDocumentID) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	uid, err := uuid.Parse(strings.TrimPrefix(s, UserDocumentIDPrefix))
	if err != nil {
		return err
	}

	*id = UserDocumentID(uid.String())
	return nil
}

func ParseUserDocumentID(id string) (UserDocumentID, error) {
	uid, err := uuid.Parse(strings.TrimPrefix(id, UserDocumentIDPrefix))
	if err != nil {
		return UserDocumentID(""), err
	}

	return UserDocumentID(uid.String()), nil
}

const ConsumerDocumentIDPrefix = "cdc-"

type ConsumerDocumentID string

func (id *ConsumerDocumentID) MarshalJSON() ([]byte, error) {
	if id == nil {
		return nil, errors.New("invalid id type")
	}

	return json.Marshal(ConsumerDocumentIDPrefix + string(*id))
}

func (id *ConsumerDocumentID) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	uid, err := uuid.Parse(strings.TrimPrefix(s, ConsumerDocumentIDPrefix))
	if err != nil {
		return err
	}

	*id = ConsumerDocumentID(uid.String())
	return nil
}

func ParseConsumerDocumentID(id string) (ConsumerDocumentID, error) {
	uid, err := uuid.Parse(strings.TrimPrefix(id, ConsumerDocumentIDPrefix))
	if err != nil {
		return ConsumerDocumentID(""), err
	}

	return ConsumerDocumentID(uid.String()), nil
}
