package bbva

const EventTypePrefix = "bbva:events:us:"

type EventType string

const (
	// Consumer
	EventTypeConsumerCreate           = EventType("customers:create")
	EventTypeConsumerUpdate           = EventType("customers:update")
	EventTypeConsumerContactCreate    = EventType("customers:contact-details:create")
	EventTypeConsumerContactUpdate    = EventType("customers:contact-details:update")
	EventTypeConsumerContactDelete    = EventType("customers:contact-details:delete")
	EventTypeConsumerAddressCreate    = EventType("customers:addresses:create")
	EventTypeConsumerAddressUpdate    = EventType("customers:addresses:update")
	EventTypeConsumerAddressDelete    = EventType("customers:addresses:delete")
	EventTypeConsumerOccupationUpdate = EventType("economic:patch")
	EventTypeConsumerKYC              = EventType("customers:kyc")

	// Business
	EventTypeBusinessCreate          = EventType("business:create")
	EventTypeBusinessUpdate          = EventType("business:update")
	EventTypeBusinessOwnerCreate     = EventType("business:owners:create")
	EventTypeBusinessOwnerDelete     = EventType("business:owners:delete")
	EventTypeBusinessIndicatorUpdate = EventType("business:indicators:update")
	EventTypeBusinessContactCreate   = EventType("business:contact-details:create")
	EventTypeBusinessContactUpdate   = EventType("business:contact-details:update")
	EventTypeBusinessContactDelete   = EventType("business:contact-details:delete")
	EventTypeBusinessAddressCreate   = EventType("business:addresses:create")
	EventTypeBusinessAddressUpdate   = EventType("business:addresses:update")
	EventTypeBusinessAddressDelete   = EventType("business:addresses:delete")
	EventTypeBusinessKYC             = EventType("business:kyc")

	// Bank Account
	EventTypeAccountCreate        = EventType("accounts:create")
	EventTypeAccountUpdate        = EventType("accounts:update")
	EventTypeAccountStatus        = EventType("accounts:status")
	EventTypeAccountBlock         = EventType("accounts:blocks")
	EventTypeAccountChargeOff     = EventType("accounts:chargeoff")
	EventTypeAccountChargeOffCard = EventType("accounts:chargeoff-card")

	// Cards
	EventTypeCardCreate = EventType("cards:create")
	EventTypeCardUpdate = EventType("cards:update")

	// Money Transfer
	EventTypePaymentStatusChange = EventType("payments:status:changed")
	EventTypePaymentCorrected    = EventType("payments:corrected:data")

	// 409 Event type does not exists
	// EventTypePaymentTransferUpdate             = EventType("payments:card:transfers:update")

	// Transactions
	EventTypeAuthorization = EventType("transaction:authorization")
	EventTypeFundhold      = EventType("transaction:fundhold")
	EventTypeCardPosted    = EventType("transaction:cardposted")
	EventTypeNonCardPosted = EventType("transaction:noncardposted")

	// Loans
	EventTypeLoanOfferExpired          = EventType("loans:credittailor:expired")
	EventTypeLoanStatusUpdated         = EventType("loans:status:update")
	EventTypeLoanProposalUpdated       = EventType("loans:proposal:status:update")
	EventTypeLoanExpiration            = EventType("loans:expiration")
	EventTypeLoanPaymentStatusChanged  = EventType("loans:payments:status:changed")
	EventTypeLoanActivation            = EventType("loans:activation")
	EventTypeLoanProposalStatusUpdated = EventType("loans:proposal:status:updated")
	EventTypeLoanProposalExpired       = EventType("loans:proposal:expired")
	EventTypeLoanActivated             = EventType("loans:activated")
)

var EventTypes = []EventType{
	// Consumer
	EventTypeConsumerCreate,
	EventTypeConsumerUpdate,
	EventTypeConsumerContactCreate,
	EventTypeConsumerContactUpdate,
	EventTypeConsumerContactDelete,
	EventTypeConsumerAddressCreate,
	EventTypeConsumerAddressUpdate,
	EventTypeConsumerAddressDelete,
	EventTypeConsumerOccupationUpdate,
	EventTypeConsumerKYC,

	// Business
	EventTypeBusinessCreate,
	EventTypeBusinessUpdate,
	EventTypeBusinessOwnerCreate,
	EventTypeBusinessOwnerDelete,
	EventTypeBusinessIndicatorUpdate,
	EventTypeBusinessContactCreate,
	EventTypeBusinessContactUpdate,
	EventTypeBusinessContactDelete,
	EventTypeBusinessAddressCreate,
	EventTypeBusinessAddressUpdate,
	EventTypeBusinessAddressDelete,
	EventTypeBusinessKYC,

	// Account
	EventTypeAccountCreate,
	EventTypeAccountUpdate,
	EventTypeAccountStatus,
	EventTypeAccountBlock,
	EventTypeAccountChargeOff,
	EventTypeAccountChargeOffCard,

	// Cards
	EventTypeCardCreate,
	EventTypeCardUpdate,

	// Money Transfer
	EventTypePaymentStatusChange,
	EventTypePaymentCorrected,
	// EventTypePayment,

	// Transactions
	EventTypeAuthorization,
	EventTypeFundhold,
	EventTypeCardPosted,
	EventTypeNonCardPosted,
}
