package business

import (
	"github.com/wiseco/core-platform/shared"
)

type PartnerName string

func (p PartnerName) String() string {
	return string(p)
}

const (
	PartnerNameShopify = PartnerName("shopify")
)

type Partner struct {
	BusinessID shared.BusinessID `json:"businessId"`

	// Partner name
	Name PartnerName `json:"name"`

	// Activation code
	ActivationCode string `json:"activationCode"`
}
