package id

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

/*
 * Party prefix + uuid will be used internally when passed between services and
 * removed when passing to/from the database.
 */
type PartyID uuid.UUID

func (id PartyID) Prefix() IDPrefix {
	return IDPrefixParty
}

func (id PartyID) String() string {
	return string(IDPrefixParty) + id.UUIDString()
}

func (id PartyID) UUIDString() string {
	return uuid.UUID(id).String()
}

func (id PartyID) IsZero() bool {
	return uuid.UUID(id) == uuid.Nil
}

// Returns empty string when zero for omitempty
func (id PartyID) JSONString() string {
	if id.IsZero() {
		return ""
	}

	return id.String()
}

func (id PartyID) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.String())
}

func (id *PartyID) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	uid, err := ParsePartyID(s)
	if err != nil {
		return err
	}

	*id = uid
	return nil
}

// SQL value marshaller
func (id PartyID) Value() (driver.Value, error) {
	return id.UUIDString(), nil
}

// SQL scanner
func (id *PartyID) Scan(value interface{}) error {
	if value == nil {
		*id = PartyID{}
		return nil
	}

	val, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("Unable to scan value %v", value)
	}

	uid, err := uuid.Parse(string(val))
	if err != nil {
		return err
	}

	*id = PartyID(uid)
	return nil
}

func (id PartyID) GetLinkedAccountID() LinkedBankAccountID {
	s := id.UUIDString()

	laID, _ := ParseLinkedBankAccountID(string(IDPrefixLinkedBankAccount) + s)

	return laID
}

func ParsePartyID(id string) (PartyID, error) {
	// Return nil id on empty string
	if id == "" {
		return PartyID{}, nil
	}

	// Check prefix
	if !strings.HasPrefix(id, string(IDPrefixParty)) {
		return PartyID{}, errors.New("invalid party id prefix")
	}

	// Validate UUID
	uid, err := uuid.Parse(strings.TrimPrefix(id, string(IDPrefixParty)))
	if err != nil {
		return PartyID{}, err
	}

	return PartyID(uid), nil
}

func PartyIDFromLinkedAccount(laID LinkedBankAccountID) PartyID {
	s := laID.UUIDString()

	pID, _ := ParsePartyID(string(IDPrefixParty) + s)

	return pID
}
