package subscription

import (
	"github.com/wiseco/core-platform/services"
	"github.com/wiseco/core-platform/shared"
)

type SubscriptionUpdate struct {
	UserID             shared.UserID               `json:"userId"`
	BusinessID         shared.BusinessID           `json:"businessId"`
	SubscriptionStatus services.SubscriptionStatus `json:"subscriptionStatus"`
}
