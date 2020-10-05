package data

import "github.com/wiseco/core-platform/partner/bank/bbva"

type EventTypeConfig struct {
	Format string
}

var EntityEventTypes = []bbva.EventType{
	// Consumer
	bbva.EventTypeConsumerCreate,
	bbva.EventTypeConsumerUpdate,
	bbva.EventTypeConsumerContactCreate,
	bbva.EventTypeConsumerContactUpdate,
	bbva.EventTypeConsumerContactDelete,
	bbva.EventTypeConsumerAddressCreate,
	bbva.EventTypeConsumerAddressUpdate,
	bbva.EventTypeConsumerAddressDelete,
	bbva.EventTypeConsumerOccupationUpdate,
	bbva.EventTypeConsumerKYC,

	// Business
	bbva.EventTypeBusinessCreate,
	bbva.EventTypeBusinessUpdate,
	bbva.EventTypeBusinessOwnerCreate,
	bbva.EventTypeBusinessOwnerDelete,
	bbva.EventTypeBusinessIndicatorUpdate,
}

var OtherBankingEventTypes = []bbva.EventType{
	// Account
	bbva.EventTypeAccountCreate,
	bbva.EventTypeAccountUpdate,
	bbva.EventTypeAccountStatus,
	bbva.EventTypeAccountBlock,
	bbva.EventTypeAccountChargeOff,
	bbva.EventTypeAccountChargeOffCard,

	// Cards
	bbva.EventTypeCardCreate,
	bbva.EventTypeCardUpdate,
}

var MoveMoneyEventTypes = []bbva.EventType{
	bbva.EventTypePaymentStatusChange,
	bbva.EventTypePaymentCorrected,
	// EventTypePayment,
}

var AccountTransactionEventTypes = []bbva.EventType{
	bbva.EventTypeFundhold,
	bbva.EventTypeNonCardPosted,
}

var CardTransactionEventTypes = []bbva.EventType{
	bbva.EventTypeAuthorization,
	bbva.EventTypeCardPosted,
}

var EntityEventConfigTypes = map[bbva.EventType]EventTypeConfig{
	// Consumer
	bbva.EventTypeConsumerCreate:           EventTypeConfig{},
	bbva.EventTypeConsumerUpdate:           EventTypeConfig{},
	bbva.EventTypeConsumerContactCreate:    EventTypeConfig{},
	bbva.EventTypeConsumerContactUpdate:    EventTypeConfig{},
	bbva.EventTypeConsumerContactDelete:    EventTypeConfig{},
	bbva.EventTypeConsumerAddressCreate:    EventTypeConfig{},
	bbva.EventTypeConsumerAddressUpdate:    EventTypeConfig{},
	bbva.EventTypeConsumerAddressDelete:    EventTypeConfig{},
	bbva.EventTypeConsumerOccupationUpdate: EventTypeConfig{},
	bbva.EventTypeConsumerKYC:              EventTypeConfig{},

	// Business
	bbva.EventTypeBusinessCreate:          EventTypeConfig{},
	bbva.EventTypeBusinessUpdate:          EventTypeConfig{},
	bbva.EventTypeBusinessOwnerCreate:     EventTypeConfig{},
	bbva.EventTypeBusinessOwnerDelete:     EventTypeConfig{},
	bbva.EventTypeBusinessIndicatorUpdate: EventTypeConfig{},
	bbva.EventTypeBusinessContactCreate:   EventTypeConfig{},
	bbva.EventTypeBusinessContactUpdate:   EventTypeConfig{},
	bbva.EventTypeBusinessContactDelete:   EventTypeConfig{},
	bbva.EventTypeBusinessAddressCreate:   EventTypeConfig{},
	bbva.EventTypeBusinessAddressUpdate:   EventTypeConfig{},
	bbva.EventTypeBusinessAddressDelete:   EventTypeConfig{},
	bbva.EventTypeBusinessKYC:             EventTypeConfig{},
}

var OtherBankingEventConfigTypes = map[bbva.EventType]EventTypeConfig{
	// Account
	bbva.EventTypeAccountCreate:        EventTypeConfig{SubscriptionFormatPRAPI},
	bbva.EventTypeAccountUpdate:        EventTypeConfig{SubscriptionFormatPRAPI},
	bbva.EventTypeAccountStatus:        EventTypeConfig{SubscriptionFormatPRAPI},
	bbva.EventTypeAccountBlock:         EventTypeConfig{SubscriptionFormatPRAPI},
	bbva.EventTypeAccountChargeOff:     EventTypeConfig{SubscriptionFormatPRAPI},
	bbva.EventTypeAccountChargeOffCard: EventTypeConfig{SubscriptionFormatPRAPI},

	// Cards
	bbva.EventTypeCardCreate: EventTypeConfig{SubscriptionFormatPRAPI},
	bbva.EventTypeCardUpdate: EventTypeConfig{SubscriptionFormatPRAPI},
}

var MoveMoneyEventConfigTypes = map[bbva.EventType]EventTypeConfig{
	bbva.EventTypePaymentStatusChange: EventTypeConfig{SubscriptionFormatPRAPI},
	bbva.EventTypePaymentCorrected:    EventTypeConfig{SubscriptionFormatPRAPI},
	// EventTypePaymentTransferUpdate: EventTypeConfig{SubscriptionFormatPRAPI},
}

var AccountTransactionEventConfigTypes = map[bbva.EventType]EventTypeConfig{
	bbva.EventTypeFundhold:      EventTypeConfig{SubscriptionFormatPRAPI},
	bbva.EventTypeNonCardPosted: EventTypeConfig{SubscriptionFormatPRAPI},
}

var CardTransactionEventConfigTypes = map[bbva.EventType]EventTypeConfig{
	bbva.EventTypeAuthorization: EventTypeConfig{SubscriptionFormatPRAPI},
	bbva.EventTypeCardPosted:    EventTypeConfig{SubscriptionFormatPRAPI},
}

var AllEventConfigTypes = []map[bbva.EventType]EventTypeConfig{
	EntityEventConfigTypes,
	OtherBankingEventConfigTypes,
	MoveMoneyEventConfigTypes,
	AccountTransactionEventConfigTypes,
	CardTransactionEventConfigTypes,
}
