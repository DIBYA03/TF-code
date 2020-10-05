package shared

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/google/uuid"
)

const PostedTransactionPrefix = "pst-"

type PostedTransactionID string

func (id *PostedTransactionID) MarshalJSON() ([]byte, error) {
	if id == nil {
		return nil, errors.New("invalid id type")
	}

	return json.Marshal(PostedTransactionPrefix + string(*id))
}
func (id *PostedTransactionID) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	uid, err := uuid.Parse(strings.TrimPrefix(s, PostedTransactionPrefix))
	if err != nil {
		return err
	}
	*id = PostedTransactionID(uid.String())

	return nil
}

func ParsePostedTransactionID(id string) (PostedTransactionID, error) {
	uid, err := uuid.Parse(strings.TrimPrefix(id, PostedTransactionPrefix))
	if err != nil {
		return PostedTransactionID(""), err
	}
	return PostedTransactionID(uid.String()), nil
}

func (id *PostedTransactionID) ToPrefixString() string {
	if id == nil {
		return ""
	}

	return PostedTransactionPrefix + string(*id)
}

const PendingTransactionPrefix = "pnt-"

type PendingTransactionID string

func (id *PendingTransactionID) MarshalJSON() ([]byte, error) {
	if id == nil {
		return nil, errors.New("invalid id type")
	}

	return json.Marshal(PendingTransactionPrefix + string(*id))
}

func (id *PendingTransactionID) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	uid, err := uuid.Parse(strings.TrimPrefix(s, PendingTransactionPrefix))
	if err != nil {
		return err
	}
	*id = PendingTransactionID(uid.String())

	return nil
}

func ParsePendingTransactionID(id string) (PendingTransactionID, error) {
	uid, err := uuid.Parse(strings.TrimPrefix(id, PendingTransactionPrefix))
	if err != nil {
		return PendingTransactionID(""), err
	}
	return PendingTransactionID(uid.String()), nil
}

func (id *PendingTransactionID) ToPrefixString() string {
	if id == nil {
		return ""
	}

	return PendingTransactionPrefix + string(*id)
}
