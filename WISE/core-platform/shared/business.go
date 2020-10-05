package shared

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

const BusinessIDPrefix = "bus-"

type BusinessID string
type BusinessIDs []string

func (bIDs BusinessIDs) Join(delimiter string) string {
	return strings.Join(bIDs, delimiter)
}

func (id BusinessID) MarshalJSON() ([]byte, error) {
	prefix := BusinessIDPrefix

	if len(id) == 0 {
		prefix = ""
	}

	return json.Marshal(prefix + string(id))
}

func (id *BusinessID) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	uid, err := uuid.Parse(strings.TrimPrefix(s, BusinessIDPrefix))
	if err != nil {
		return err
	}

	*id = BusinessID(uid.String())
	return nil
}

func (id *BusinessID) Scan(value interface{}) error {
	if value == nil {
		*id = BusinessID("")

		return nil
	}

	val, ok := value.([]byte)

	if !ok {
		return fmt.Errorf("Unable to scan value %v", value)
	}

	*id = BusinessID(val)

	return nil
}

func ParseBusinessID(id string) (BusinessID, error) {
	uid, err := uuid.Parse(strings.TrimPrefix(id, BusinessIDPrefix))
	if err != nil {
		return BusinessID(""), err
	}

	return BusinessID(uid.String()), nil
}

func (id *BusinessID) ToPrefixString() string {
	if id == nil {
		return ""
	}

	return BusinessIDPrefix + string(*id)
}

const BusinessMemberIDPrefix = "bme-"

type BusinessMemberID string

func (id *BusinessMemberID) MarshalJSON() ([]byte, error) {
	if id == nil {
		return nil, errors.New("invalid id type")
	}

	return json.Marshal(BusinessMemberIDPrefix + string(*id))
}

func (id *BusinessMemberID) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	uid, err := uuid.Parse(strings.TrimPrefix(s, BusinessMemberIDPrefix))
	if err != nil {
		return err
	}

	*id = BusinessMemberID(uid.String())
	return nil
}

func ParseBusinessMemberID(id string) (BusinessMemberID, error) {
	uid, err := uuid.Parse(strings.TrimPrefix(id, BusinessMemberIDPrefix))
	if err != nil {
		return BusinessMemberID(""), err
	}

	return BusinessMemberID(uid.String()), nil
}

// Return DBA else LegalName else empty string
func GetBusinessName(legalName *string, dba []string) string {
	name := GetDBAName(dba)
	if len(name) > 0 {
		return name
	}

	if legalName != nil {
		return *legalName
	}

	return ""
}

func GetDBAName(dba []string) string {
	if dba != nil && len(dba) > 0 && len(strings.TrimSpace(dba[0])) > 0 {
		// Use first DBA
		return strings.TrimSpace(dba[0])
	}

	return ""
}
