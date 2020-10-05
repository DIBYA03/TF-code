package shared

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/google/uuid"
)

const CardReaderPrefix = "crr-"

type CardReaderID string

func (id *CardReaderID) MarshalJSON() ([]byte, error) {
	if id == nil {
		return nil, errors.New("invalid id type")
	}

	return json.Marshal(CardReaderPrefix + string(*id))
}

func (id *CardReaderID) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	uid, err := uuid.Parse(strings.TrimPrefix(s, CardReaderPrefix))
	if err != nil {
		return err
	}
	*id = CardReaderID(uid.String())

	return nil
}

func ParseCardReaderID(id string) (CardReaderID, error) {
	uid, err := uuid.Parse(strings.TrimPrefix(id, CardReaderPrefix))
	if err != nil {
		return CardReaderID(""), err
	}

	return CardReaderID(uid.String()), nil
}

const PaymentRequestPrefix = "pmr-"

type PaymentRequestID string

func (id *PaymentRequestID) MarshalJSON() ([]byte, error) {
	if id == nil {
		return nil, errors.New("invalid id type")
	}

	return json.Marshal(PaymentRequestPrefix + string(*id))
}
func (id *PaymentRequestID) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	uid, err := uuid.Parse(strings.TrimPrefix(s, PaymentRequestPrefix))
	if err != nil {
		return err
	}
	*id = PaymentRequestID(uid.String())

	return nil
}

func (id *PaymentRequestID) ToPrefixString() string {
	if id == nil {
		return ""
	}

	return PaymentRequestPrefix + string(*id)
}

func (id *PaymentRequestID) ToUUIDString() string {
	if id == nil {
		return ""
	}

	return string(*id)
}

func ParsePaymentRequestID(id string) (PaymentRequestID, error) {
	uid, err := uuid.Parse(strings.TrimPrefix(id, PaymentRequestPrefix))
	if err != nil {
		return PaymentRequestID(""), err
	}

	return PaymentRequestID(uid.String()), nil
}
