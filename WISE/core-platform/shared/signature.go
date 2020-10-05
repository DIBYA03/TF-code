package shared

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/google/uuid"
)

const SignatureRequestPrefix = "sr-"

type SignatureRequestID string

func (id *SignatureRequestID) MarshalJSON() ([]byte, error) {
	if id == nil {
		return nil, errors.New("invalid id type")
	}

	return json.Marshal(SignatureRequestPrefix + string(*id))
}

func (id *SignatureRequestID) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	uid, err := uuid.Parse(strings.TrimPrefix(s, SignatureRequestPrefix))
	if err != nil {
		return err
	}
	*id = SignatureRequestID(uid.String())

	return nil
}

func ParseSignatureRequestID(id string) (SignatureRequestID, error) {
	uid, err := uuid.Parse(strings.TrimPrefix(id, SignatureRequestPrefix))
	if err != nil {
		return SignatureRequestID(""), err
	}

	return SignatureRequestID(uid.String()), nil
}
