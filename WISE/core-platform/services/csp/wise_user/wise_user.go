package wise_user

import (
	"time"

	"github.com/wiseco/core-platform/shared"
)

type OriginatedFrom string

const (
	OriginatedFromCSP = OriginatedFrom("Customer Success Portal")
)

type PhoneChangeRequestCreate struct {
	UserID            shared.UserID  `json:"userId" db:"user_id"`
	OldPhone          string         `json:"oldPhone" db:"old_phone"`
	NewPhone          string         `json:"newPhone" db:"new_phone"`
	OriginatedFrom    OriginatedFrom `json:"originatedFrom" db:"originated_from"`
	CSPUserID         string         `json:"cspUserId" db:"csp_user_id"`
	VerificationNotes string         `json:"verificationNotes" db:"verification_notes"`
}

type PhoneChangeRequest struct {
	ID                string         `json:"id" db:"id"`
	UserID            shared.UserID  `json:"userId" db:"user_id"`
	OldPhone          string         `json:"oldPhone" db:"old_phone"`
	NewPhone          string         `json:"newPhone" db:"new_phone"`
	OriginatedFrom    OriginatedFrom `json:"originatedFrom" db:"originated_from"`
	CSPUserID         *string        `json:"cspUserId" db:"csp_user_id"`
	VerificationNotes string         `json:"verificationNotes" db:"verification_notes"`
	Created           time.Time      `json:"created" db:"created"`
	Modified          time.Time      `json:"modified" db:"modified"`
	CSPUserName       *string
}
