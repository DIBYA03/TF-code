package business

import (
	"github.com/wiseco/core-platform/services/banking"
	"github.com/wiseco/core-platform/shared"
)

// BankCardPartial contains limited data about the card
// Sensitive info like Card number is removed
type BankCardPartial struct {
	BusinessID shared.BusinessID `json:"businessId" db:"business_id"` // Related business account for this card

	// Card block
	CardBlock []*BankCardBlock `json:"cardBlock"`

	// Card reissue history
	CardReissueHistory []CardReissueHistory `json:"cardReissueHistory"`

	banking.BankCard
}

// Business debit or credit card
type BankCard struct {
	BankCardPartial
	CardNumberAlias string `json:"cardNumberAlias" db:"card_number_alias"`
}

type BankCardCreate struct {
	BusinessID shared.BusinessID `json:"businessId"` // Related business account for this card

	banking.BankCardCreate
}

type BankCardActivate struct {
	BusinessID shared.BusinessID `json:"businessId"` // Related business account for this card

	banking.BankCardActivate
}

type BankCardBlockCreate struct {
	BusinessID shared.BusinessID `json:"businessId"` // Related business account for this card

	banking.BankCardBlockCreate
}

type BankCardBlockDelete struct {
	BusinessID shared.BusinessID `json:"businessId"` // Related business account for this card

	banking.BankCardBlockDelete
}

type BankCardBlock struct {
	BusinessID shared.BusinessID `json:"businessId"` // Related business account for this card

	banking.BankCardBlock
}

type BankCardUpdate = banking.BankCardUpdate

type CardReissueRequest struct {
	BusinessID shared.BusinessID `db:"business_id"`
	banking.CardReissueRequest
}

type CardReissueHistory struct {
	BusinessID shared.BusinessID `json:"businessId" db:"business_id"`
	banking.CardReissueHistory
}
