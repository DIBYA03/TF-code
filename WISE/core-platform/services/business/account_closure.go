package business

import "github.com/wiseco/core-platform/shared"

//AccountClosureItem ...
type AccountClosureItem struct {
	ID          string            `json:"id" db:"id"`
	BusinessID  shared.BusinessID `json:"businessId,omitempty" db:"business_id"`
	Reason      *string           `json:"reason,omitempty" db:"reason"`
	Description *string           `json:"description,omitempty" db:"description"`
	Status      string            `json:"status,omitempty" db:"status"`
}

//AccountClosureCreate ...
type AccountClosureCreate struct {
	BusinessID  shared.BusinessID `json:"businessId" db:"business_id"`
	Reason      string            `json:"reason" db:"reason"`
	Description string            `json:"description" db:"description"`
}
