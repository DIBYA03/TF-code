package business

import (
	"time"

	"github.com/jmoiron/sqlx/types"
	"github.com/wiseco/core-platform/services/banking"
	"github.com/wiseco/core-platform/shared"
)

type UsageType string

const UsageTypeNone = UsageType("")

const (
	UsageTypeDemo             = UsageType("demo")             // Demo account
	UsageTypePrimary          = UsageType("primary")          // Primary account
	UsageTypeClearing         = UsageType("clearing")         // Clearing account
	UsageTypeExternal         = UsageType("external")         // External account (Business own plaid accounts)
	UsageTypeContact          = UsageType("contact")          // Contact account (Contact account)
	UsageTypeContactInvisible = UsageType("contactInvisible") // Contact account but not visible to business(Contact account)
	UsageTypeMerchant         = UsageType("merchant")         // Merchant account used for card payments
)

// BusinessBankAccount bank account for a business
type BankAccount struct {
	// Business id
	BusinessID shared.BusinessID `json:"businessId" db:"business_id"` // Related business account for this bank account

	// Account usage type
	UsageType UsageType `json:"usageType" db:"usage_type"`

	// Basic account struct
	banking.BankAccount
}

func (a *BankAccount) toAccountBalance() BankAccountBalance {
	return BankAccountBalance{
		AccountID:        a.Id,
		AvailableBalance: a.AvailableBalance,
		PostedBalance:    a.PostedBalance,
		Currency:         a.Currency,
		Modified:         a.Modified,
	}
}

// BankAccountCreate use to create a business bank account
type BankAccountCreateFull struct {

	//Partner bank Account id (e.g BBVA account id)
	BankAccountId string `json:"bankAccountId" db:"bank_account_id"`

	// Partner bank account extra data
	BankExtra types.JSONText `json:"-" db:"bank_extra"`

	// Account holder user id
	AccountHolderID shared.UserID `json:"accountHolderId" db:"account_holder_id"`

	// Account status e.g. active or closed
	AccountStatus string `json:"accountStatus" db:"account_status"`

	// Primary account number
	AccountNumber string `json:"accountNumber" db:"account_number"`

	// Primary account routing number
	RoutingNumber string `json:"routingNumber" db:"routing_number"`

	// Wire Routing Number
	WireRouting *string `json:"wireRouting" db:"wire_routing"`

	// Available balance for widthdrawals
	AvailableBalance float64 `json:"availableBalance" db:"available_balance"`

	// Current posted balance
	PostedBalance float64 `json:"postedBalance" db:"posted_balance"`

	// Denominated currency
	Currency banking.Currency `json:"currency" db:"currency"`

	// Timestamp opened (UTC)
	Opened time.Time `json:"opened" db:"opened"`

	// Basic account create struct
	BankAccountCreate
}

type BankAccountCreate struct {
	// Business id
	BusinessID shared.BusinessID `json:"businessId" db:"business_id"`

	// Account usage type
	UsageType UsageType `json:"usageType" db:"usage_type"`

	// Basic account create struct
	banking.BankAccountCreate
}

type BankAccountUpdate struct {
	// Business Id
	BusinessID shared.BusinessID `json:"businessId" db:"business_id"`

	// Basic update struct
	banking.BankAccountUpdate
}

type BankAccountStatement struct {
	// Bank account id
	BusinessID shared.BusinessID `json:"businessId"`

	// Basic account statement struct
	banking.BankAccountStatement
}

type BankAccountBalance struct {
	// Bank account id
	AccountID string `json:"accountId" db:"id"`

	// Available balance for widthdrawals
	AvailableBalance float64 `json:"availableBalance" db:"available_balance"`

	// Current posted balance
	PostedBalance float64 `json:"postedBalance" db:"posted_balance"`

	// Current AvailableBalance - PendingDebitBalance
	ActualBalance float64 `json:"actualBalance" db:"-"`

	// Denominated currency
	Currency banking.Currency `json:"currency" db:"currency"`

	// Timestamp modified (UTC)
	Modified time.Time `json:"modified" db:"modified"`
}
