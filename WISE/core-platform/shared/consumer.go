package shared

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

const ConsumerIDPrefix = "con-"

type ConsumerID string

func (id ConsumerID) ToPrefixString() string {
	prefix := ConsumerIDPrefix

	if len(id) == 0 {
		prefix = ""
	}

	return prefix + string(id)
}

func (id ConsumerID) MarshalJSON() ([]byte, error) {
	prefix := ConsumerIDPrefix

	if len(id) == 0 {
		prefix = ""
	}

	return json.Marshal(prefix + string(id))
}

func (id *ConsumerID) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	uid, err := uuid.Parse(strings.TrimPrefix(s, ConsumerIDPrefix))
	if err != nil {
		return err
	}

	*id = ConsumerID(uid.String())
	return nil
}

func (id *ConsumerID) Scan(value interface{}) error {
	if value == nil {
		*id = ConsumerID("")

		return nil
	}

	val, ok := value.([]byte)

	if !ok {
		return fmt.Errorf("Unable to scan value %v", value)
	}

	*id = ConsumerID(val)

	return nil
}

func ParseConsumerID(id string) (ConsumerID, error) {
	uid, err := uuid.Parse(strings.TrimPrefix(id, ConsumerIDPrefix))
	if err != nil {
		return ConsumerID(""), err
	}

	return ConsumerID(uid.String()), nil
}

// BBVA mandated visa card name max length
const maxVisaCardName = 20

func GetVisaCardHolderName(firstName, middleName, lastName string) (string, error) {
	fn := strings.TrimSpace(firstName)
	mn := strings.TrimSpace(middleName)
	ln := strings.TrimSpace(lastName)

	if len(fn) == 0 || len(ln) == 0 {
		return "", errors.New("invalid card holder name")
	}

	// Try first middle last
	var name string
	if len(mn) > 0 {
		// Try middle initial
		name = fn + " " + mn[0:1] + " " + ln
		if len(name) <= maxVisaCardName {
			return name, nil
		}
	}

	// Try first and last
	name = fn + " " + ln
	if len(name) <= maxVisaCardName {
		return name, nil
	}

	// If space for at least first initial
	if len(ln) <= maxVisaCardName-2 {
		return fn[0:maxVisaCardName-len(ln)-1] + " " + ln, nil
	}

	// Use first initial and truncate last name
	return fn[0:1] + " " + ln[0:maxVisaCardName-2], nil
}
