package order

import (
	"github.com/wiseco/core-platform/shared"
	"github.com/wiseco/go-lib/id"
)

type ShopifyOrderMessage struct {
	ID         id.ShopifyOrderID `json:"id"`
	OrderID    int64             `json:"orderId"`
	BusinessID shared.BusinessID `json:"businessId"`
	EventType  EventType         `json:"eventType"`
	FirstName  *string           `json:"firstName"`
	LastName   *string           `json:"lastName"`
	Phone      *string           `json:"phone"`
	Email      *string           `json:"email"`
	Amount     string            `json:"amount"`
}

type EventType string

const (
	EventTypeOrderCreated            = EventType("orderCreated")
	EventTypeOrderCanceled           = EventType("orderCanceled")
	EventTypeOrderFulfilled          = EventType("orderFulfilled")
	EventTypeOrderPaid               = EventType("orderPaid")
	EventTypeOrderPartiallyFulfilled = EventType("orderPartiallyFulfilled")
	EventTypeOrderUpdated            = EventType("orderUpdated")
)
