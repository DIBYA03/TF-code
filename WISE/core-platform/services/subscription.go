package services

type SubscriptionStatus string

const (
	SubscriptionStatusPendingAcceptance = SubscriptionStatus("pending_acceptance")
	SubscriptionStatusActive            = SubscriptionStatus("active")
	SubscriptionStatusUnpaid            = SubscriptionStatus("unpaid")
	SubscriptionStatusCanceled          = SubscriptionStatus("canceled")
)

var SubscriptionStatusValidator = map[SubscriptionStatus]SubscriptionStatus{
	SubscriptionStatusActive: SubscriptionStatusActive,
	SubscriptionStatusUnpaid: SubscriptionStatusUnpaid,
}
