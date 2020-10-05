package consumer

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/jmoiron/sqlx/types"
	partnerbank "github.com/wiseco/core-platform/partner/bank"
	"github.com/wiseco/core-platform/services"
	"github.com/wiseco/core-platform/shared"
)

//UserNotification settings object controls how and when push notifications are sent
type UserNotification struct {
	// Send notification on transactions
	Transactions *bool `json:"transactions"`

	// Send notification on transfers
	Transfers *bool `json:"transfers"`

	// Send notification on contact changes
	Contacts *bool `json:"contacts"`
}

// SQL value marshaller
func (n UserNotification) Value() (driver.Value, error) {
	return json.Marshal(n)
}

// SQL scan unmarshaller
func (n *UserNotification) Scan(src interface{}) error {
	source, ok := src.([]byte)
	if !ok {
		return errors.New("Type convertible to []byte")
	}

	var out UserNotification
	err := json.Unmarshal(source, &out)
	if err != nil {
		return err
	}

	*n = out
	return nil
}

type UserNotificationUpdate UserNotification

func (u UserNotificationUpdate) Value() (driver.Value, error) {
	return json.Marshal(u)
}

// Update  Consumer update struct
type Update struct {
	// First name
	FirstName *string `json:"firstName" db:"first_name"`

	// Middle name
	MiddleName *string `json:"middleName" db:"middle_name"`

	// Last name
	LastName *string `json:"lastName" db:"last_name"`

	// Email
	Email *string `json:"email" db:"email"`

	// Phone
	Phone *string `json:"phone" db:"phone"`

	// Date of birth
	DateOfBirth *shared.Date `json:"dateOfBirth" db:"date_of_birth"`

	// Tax id number
	TaxID *string `json:"taxId" db:"tax_id"`

	// Tax id type e.g. SSN or ITIN
	TaxIDType *string `json:"taxIdType" db:"tax_id_type"`

	// Legal address
	LegalAddress *services.Address `json:"legalAddress" db:"legal_address"`

	// Mailing address
	MailingAddress *services.Address `json:"mailingAddress" db:"mailing_address"`

	// Work address
	WorkAddress *services.Address `json:"workAddress" db:"work_address"`

	// Residency
	Residency *services.Residency `json:"residency" db:"residency"`

	// List of citizenships
	CitizenshipCountries *services.StringArray `json:"citizenshipCountries" db:"citizenship_countries"`

	// Consumer's occupation
	Occupation *string `json:"occupation" db:"occupation"`

	// Consumer income source e.g. salary or inheritance
	IncomeType *services.StringArray `json:"incomeType" db:"income_type"`

	// Consumer activity type
	ActivityType *services.StringArray `json:"activityType" db:"activity_type"`
}

// User User struct
type User struct {
	// User id (uuid)
	ID string `json:"id" db:"id"`

	// ConsumerID is shared by users and members
	ConsumerID shared.ConsumerID `json:"consumerId" db:"consumer_id"`

	// AuthID from auth server
	IdentityID string `json:"-" db:"identity_id"`

	// Partner id
	PartnerID *string `json:"partnerId" db:"partner_id"`

	// Email
	Email *string `json:"email" db:"email"`

	// Email verified
	EmailVerified bool `json:"emailVerified" db:"email_verified"`

	// Phone
	Phone string `json:"phone" db:"phone"`

	// Phone verified
	PhoneVerified bool `json:"phoneVerified" db:"phone_verified"`

	// Notification settings object
	Notification UserNotification `json:"notification" db:"notification"`

	// Deactivated
	Deactivated *time.Time `json:"deactivated" db:"deactivated"`

	// Created time
	Created time.Time `json:"created" db:"created"`

	// Last modified time
	Modified time.Time `json:"modified" db:"modified"`

	// First name
	FirstName string `json:"firstName" db:"first_name"`

	// Middle name
	MiddleName string `json:"middleName" db:"middle_name"`

	// Last name
	LastName string `json:"lastName" db:"last_name"`

	// Date of birth
	DateOfBirth *shared.Date `json:"dateOfBirth" db:"date_of_birth"`

	// Tax id number
	TaxID *services.TaxID `json:"taxIdMasked" db:"tax_id"`

	// Tax id type e.g. SSN or ITIN
	TaxIDType *string `json:"taxIdType" db:"tax_id_type"`

	// KYC Status
	KYCStatus string `json:"kycStatus" db:"kyc_status"`

	// Legal address
	LegalAddress *services.Address `json:"legalAddress" db:"legal_address"`

	// Mailing address
	MailingAddress *services.Address `json:"mailingAddress" db:"mailing_address"`

	// Work address
	WorkAddress *services.Address `json:"workAddress" db:"work_address"`

	// Residency
	Residency *services.Residency `json:"residency" db:"residency"`

	// List of citizenships
	CitizenshipCountries services.StringArray `json:"citizenshipCountries" db:"citizenship_countries"`

	// User's occupation
	Occupation *string `json:"occupation" db:"occupation"`

	// User income source e.g. salary or inheritance
	IncomeType services.StringArray `json:"incomeType" db:"income_type"`

	// User activity type
	ActivityType services.StringArray `json:"activityType" db:"activity_type"`

	// Is this a restricted person
	IsRestricted bool `json:"isRestricted" db:"is_restricted"`

	// BankID bank id (CO)
	BankID *partnerbank.ConsumerBankID `json:"bankId"`

	// Subscription
	SubscriptionStatus *services.SubscriptionStatus `json:"subscriptionStatus" db:"subscription_status"`
}

// CSPConsumer  consumer review item
type CSPConsumer struct {
	ID           string                `json:"id" db:"id"`
	ConsumerName *string               `json:"consumerName" db:"consumer_name"`
	ConsumerID   shared.ConsumerID     `json:"consumerId" db:"consumer_id"`
	Status       string                `json:"status" db:"review_status"`
	IDVs         *services.StringArray `json:"idvs" db:"idvs"`
	Notes        types.NullJSONText    `json:"notes" db:"notes"`
	Submitted    *time.Time            `json:"submitted" db:"submitted"`
	Resolved     *time.Time            `json:"resolved" db:"resolved"`
	Modified     time.Time             `json:"modified" db:"modified"`
	Created      time.Time             `json:"created" db:"created"`
}

// CSPConsumerUpdate ..
type CSPConsumerUpdate struct {
	ConsumerName  *string               `json:"consumerName" db:"consumer_name"`
	ConsumerID    *shared.ConsumerID    `json:"consumerId" db:"consumer_id"`
	Status        *string               `json:"status" db:"review_status"`
	ProcessStatus *string               `json:"processStatus" db:"process_status"`
	IDVs          *services.StringArray `json:"idvs" db:"idvs"`
	Notes         *types.NullJSONText   `json:"notes" db:"notes"`
	Submitted     *time.Time            `json:"submitted" db:"submitted"`
	Resolved      *time.Time            `json:"resolved" db:"resolved"`
}

// CSPConsumerCreate ..
type CSPConsumerCreate struct {
	ConsumerName  *string               `db:"consumer_name"`
	ConsumerID    shared.ConsumerID     `db:"consumer_id"`
	Status        string                `db:"review_status"`
	ProcessStatus *string               `db:"process_status"`
	IDVs          *services.StringArray `db:"idvs"`
	Notes         *types.JSONText       `db:"notes"`
	Submitted     *time.Time            `db:"submitted"`
}

// ConsumerState ..
type ConsumerState struct {
	ID            string    `json:"id" db:"id"`
	ConsumerID    string    `json:"consumerId" db:"consumer_id"`
	ReviewStatus  string    `json:"status" db:"review_status"`
	ProcessStatus string    `json:"processStatus" db:"process_status"`
	Created       time.Time `json:"created" db:"created"`
}

// ConsumerStateCreate ..
type ConsumerStateCreate struct {
	ConsumerID    string `db:"consumer_id"`
	ReviewStatus  string `db:"review_status"`
	ProcessStatus string `db:"process_status"`
}
