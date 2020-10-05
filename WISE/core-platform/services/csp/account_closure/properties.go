package accountclosure

import (
	"time"

	bus "github.com/wiseco/core-platform/services/business"
	user "github.com/wiseco/core-platform/services/user"
	"github.com/wiseco/core-platform/shared"
	"github.com/wiseco/go-lib/num"
)

//AccountClosureStatus ...
type AccountClosureStatus string

//CSPAccountClosureQueryParams ...
type CSPAccountClosureQueryParams struct {
	Status              string
	StartDate           string
	EndDate             string
	BusinessID          string
	BusinessName        string
	OwnerName           string
	AvailableBalanceMin string
	AvailableBalanceMax string
	PostedBalanceMin    string
	PostedBalanceMax    string
	Offset              string
	Limit               string
	SortField           string
	SortDirection       string
}

// CSPClosureRequestItem ...
type CSPClosureRequestItem struct {
	ID                 string            `json:"id" db:"id"`
	Status             *string           `json:"status" db:"status"`
	Reason             *string           `json:"reason" db:"reason"`
	Description        *string           `json:"description" db:"description"`
	Created            time.Time         `json:"created" db:"created"`
	Modified           time.Time         `json:"modified" db:"modified"`
	Closed             *time.Time        `json:"closed" db:"closed"`
	BusinessID         shared.BusinessID `json:"businessId" db:"business_id"`
	BusinessName       *string           `json:"businessName" db:"business_name"`
	AvailableBalance   *num.Decimal      `json:"availableBalance" db:"available_balance"`
	PostedBalance      *num.Decimal      `json:"postedBalance" db:"posted_balance"`
	OwnerConsumerID    shared.ConsumerID `json:"ownerConsumerId" db:"owner_consumer_id"`
	OwnerUserID        shared.UserID     `json:"ownerUserId" db:"owner_user_id"`
	OwnerFirstName     *string           `json:"ownerFirstName" db:"owner_first_name"`
	OwnerMiddleName    *string           `json:"ownerMiddleName" db:"owner_middle_name"`
	OwnerLastName      *string           `json:"ownerLastName" db:"owner_last_name"`
	Business           *bus.Business     `json:"business,omitempty" db:"business"`
	Owner              *user.Consumer    `json:"owner,omitempty" db:"consumer"`
	DigitalCheckNumber *string           `json:"digitalCheckNumber" db:"digital_check_number"`
	RefundAmount       *num.Decimal      `json:"refundAmount" db:"refund_amount"`
	CSPAgentID         *string           `json:"cspAgentId" db:"csp_agent_id"`
	CSPUserName        *string           `json:"cspUserName" db:"csp_user_name"`
}

//CSPAccountClosureCreate ...
type CSPAccountClosureCreate struct {
	BusinessID  shared.BusinessID `json:"businessId" db:"business_id"`
	Reason      string            `json:"reason" db:"reason"`
	Description string            `json:"description" db:"description"`
}

// CSPClosureRequestPatch ...
type CSPClosureRequestPatch struct {
	Status  string `json:"status"`
	AgendID string `json:"agentId"`
}

//CheckAccountClosureStatus ..
func CheckAccountClosureStatus(s string) bool {
	_, ok := accountClosureStatus[AccountClosureStatus(s)]
	return ok
}

const (
	// AccountClosureRequestPending ...
	AccountClosureRequestPending = AccountClosureStatus("pending")

	// AccountClosureRequestApproved ...
	AccountClosureRequestApproved = AccountClosureStatus("approved")

	// AccountClosureRequestClosed ...
	AccountClosureRequestClosed = AccountClosureStatus("account_closed")

	// AccountClosureRequestCanceled ...
	AccountClosureRequestCanceled = AccountClosureStatus("canceled")

	// AccountClosureRefundPending ...
	AccountClosureRefundPending = AccountClosureStatus("refund_pending")

	// AccountClosureFailed ...
	AccountClosureFailed = AccountClosureStatus("failed")

	// AccountClosureFailedRetry ...
	AccountClosureFailedRetry = AccountClosureStatus("failed_retry")
)

var accountClosureStatus = map[AccountClosureStatus]AccountClosureStatus{
	AccountClosureRequestPending:  AccountClosureRequestPending,
	AccountClosureRequestApproved: AccountClosureRequestApproved,
	AccountClosureRequestClosed:   AccountClosureRequestClosed,
	AccountClosureRequestCanceled: AccountClosureRequestCanceled,
	AccountClosureFailed:          AccountClosureFailed,
	AccountClosureFailedRetry:     AccountClosureFailedRetry,
}

//FIELDS
const (
	AccountClosureFieldBusinessName       = "businessName"
	AccountClosureFieldCreated            = "created"
	AccountClosureFieldClosed             = "closed"
	AccountClosureFieldRefundAmount       = "refundAmount"
	AccountClosureFieldDigitalCheckNumber = "digitalCheckNumber"
	AccountClosureFieldOwnerName          = "ownerName"
	AccountClosureFieldStatus             = "status"
	AccountClosureFieldAvailableBalance   = "availableBalance"
	AccountClosureFieldPostedBalance      = "postedBalance"
)

//SORT DIRECTION
const (
	AccountClosureSortAscending  = "asc"
	AccountClosureSortDescending = "desc"
)

//REQUEST API PARAM STATUS
const (
	AccountClosureListRequestPending = "requestPending"
	AccountClosureListRequestClosed  = "requestClosed"
	AccountClosureListRequestFailed  = "requestFailed"
)

//AccountClosureState ...
type AccountClosureState string

// CSPClosureState ...
type CSPClosureState struct {
	ID               string    `json:"id" db:"id"`
	ClosureRequestID string    `json:"closureRequestID" db:"account_closure_request_id"`
	State            *string   `json:"state" db:"closure_state"`
	ItemID           *string   `json:"itemId" db:"item_id"`
	Description      *string   `json:"description" db:"description"`
	Created          time.Time `json:"created" db:"created"`
	Modified         time.Time `json:"modified" db:"modified"`
}

// CSPClosureStatePostBody ...
type CSPClosureStatePostBody struct {
	State            AccountClosureState `json:"state" db:"closure_state"`
	ClosureRequestID string              `json:"closureRequestID" db:"account_closure_request_id"`
	ItemID           *string             `json:"itemId" db:"item_id"`
	Description      *string             `json:"description" db:"description"`
}

// Closure State
const (
	ACStateRequestCreated = AccountClosureState("request_created")

	ACStateRequestApproved       = AccountClosureState("request_approved")
	ACStateRequestCanceled       = AccountClosureState("request_canceled")
	ACStateRequestRetryRequested = AccountClosureState("retry_requested")

	ACStateCancelCardStarted = AccountClosureState("cancel_card_started")
	ACStateCancelCardFailed  = AccountClosureState("cancel_card_failed")
	ACStateCancelCardSuccess = AccountClosureState("cancel_card_success")

	ACStatePullBalanceStarted = AccountClosureState("pull_balance_started")
	ACStatePullBalanceFailed  = AccountClosureState("pull_balance_failed")
	ACStatePullBalanceSuccess = AccountClosureState("pull_balance_success")

	ACStateDeactivateAccountStarted = AccountClosureState("deactivate_account_started")
	ACStateDeactivateAccountFailed  = AccountClosureState("deactivate_account_failed")
	ACStateDeactivateAccountSuccess = AccountClosureState("deactivate_account_success")

	ACStateDeactivateBusinessStarted = AccountClosureState("deactivate_business_started")
	ACStateDeactivateBusinessFailed  = AccountClosureState("deactivate_business_failed")
	ACStateDeactivateBusinessSuccess = AccountClosureState("deactivate_business_success")
)
