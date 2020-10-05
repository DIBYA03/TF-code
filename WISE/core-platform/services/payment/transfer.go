package payment

import (
	"time"

	"github.com/wiseco/core-platform/services"
	"github.com/wiseco/core-platform/services/contact"
	"github.com/wiseco/core-platform/shared"
)

type RequestMode string

const (
	ReqeuestModeSMS     = RequestMode("sms")
	RequestModeEmail    = RequestMode("email")
	RequestModeSMSEmail = RequestMode("smsEmail")
)

var RequestModeToRequestMode = map[RequestMode]RequestMode{
	ReqeuestModeSMS:     ReqeuestModeSMS,
	RequestModeEmail:    RequestModeEmail,
	RequestModeSMSEmail: RequestModeSMSEmail,
}

type TransferRequestCreate struct {
	CreatedUserID  shared.UserID     `json:"createdUserId" db:"created_user_id"`
	BusinessID     shared.BusinessID `json:"businessId" db:"business_id"`
	ContactID      string            `json:"contactId" db:"contact_id"`
	RequestMode    RequestMode       `json:"requestMode" db:"request_mode"`
	ExpirationDate *time.Time        `json:"-" db:"expiration_date"`
	Notes          string            `json:"notes" db:"notes"`
	Amount         shared.Decimal    `json:"amount" db:"amount"`
}

type TransferRequestUpdate struct {
	ID              string  `json:"id" db:"id"`
	MoneyTransferID *string `json:"moneyTransferId" db:"money_transfer_id"`
	PaymentToken    *string `json:"-" db:"payment_token"`
}

type TransferRequest struct {
	ID              string            `json:"id" db:"id"`
	CreatedUserID   shared.UserID     `json:"createdUserId" db:"created_user_id"`
	BusinessID      shared.BusinessID `json:"businessId" db:"business_id"`
	ContactID       string            `json:"contactId" db:"contact_id"`
	RequestMode     RequestMode       `json:"requestMode" db:"request_mode"`
	MoneyTransferID *string           `json:"moneyTransferId" db:"money_transfer_id"`
	PaymentToken    *string           `json:"-" db:"payment_token"`
	ExpirationDate  *time.Time        `json:"-" db:"expiration_date"`
	Notes           string            `json:"notes" db:"notes"`
	Amount          shared.Decimal    `json:"amount" db:"amount"`
	Created         time.Time         `json:"created" db:"created"`
	Modified        time.Time         `json:"modified" db:"modified"`
}

type BusinessContact struct {
	LegalName           string               `json:"legalName" db:"legal_name"`
	DBA                 services.StringArray `json:"dba" db:"dba"`
	ContactFirstName    *string              `json:"firstName" db:"first_name"`
	ContactLastName     *string              `json:"lastName" db:"last_name"`
	ContactBusinessName *string              `json:"businessName" db:"business_name"`
	ContactPhone        string               `json:"phone" db:"phone_number"`
	ContactEmail        string               `json:"email" db:"email"`
	ContactType         contact.ContactType  `json:"contactType" db:"contact_type"`
}

// Response from payments table
type TransferRequestResponse struct {
	BusinessID            shared.BusinessID    `db:"business.id"`
	UserID                shared.UserID        `db:"owner_id"`
	ContactID             *string              `db:"contact_id"`
	ContactFirstName      *string              `db:"business_contact.first_name"`
	ContactLastName       *string              `db:"business_contact.last_name"`
	ContactBusinessName   *string              `db:"business_contact.business_name"`
	ContactType           contact.ContactType  `db:"business_contact.contact_type"`
	TransferRequestID     string               `db:"money_transfer_request.id"`
	BusinessBankAccountID *string              `db:"business_bank_account_id"`
	RegisteredAccountID   string               `db:"business_linked_bank_account.id"`
	LegalName             *string              `db:"legal_name"`
	DBA                   services.StringArray `db:"dba"`
	BusinessName          string
	Amount                float64 `db:"amount"`
	Notes                 string  `db:"notes"`
}
