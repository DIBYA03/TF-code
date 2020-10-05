package id

/*
 * User prefix + uuid will be used internally when passed between services and
 * removed when passing to/from the database.
 */

type IDPrefix string

var uuidRE = "[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}"

const (
	// Base UUID - useful when dealing with ambiguous type
	IDPrefixNone = IDPrefix("")

	// Logging
	IDPrefixRequest = IDPrefix("req-")

	// Auth
	IDPrefixIdentity     = IDPrefix("idn-")
	IDPrefixPartner      = IDPrefix("ptr-")
	IDPrefixPartnerKey   = IDPrefix("ptk-")
	IDPrefixClientKey    = IDPrefix("pck-")
	IDPrefixSecretKey    = IDPrefix("psk-")
	IDPrefixAccessToken  = IDPrefix("act-")
	IDPrefixRefreshToken = IDPrefix("rft-")

	// Agent
	IDPrefixCspAgent = IDPrefix("csa-")

	// Core
	IDPrefixBusiness = IDPrefix("bus-")
	IDPrefixConsumer = IDPrefix("con-")
	IDPrefixUser     = IDPrefix("usr-")
	IDPrefixContact  = IDPrefix("cnt-")
	IDPrefixAddress  = IDPrefix("adr-")
	IDPrefixDocument = IDPrefix("doc-")
	IDPrefixMember   = IDPrefix("mem-")
	IDPrefixEmail    = IDPrefix("eml-")

	// Banking
	IDPrefixBankAccount          = IDPrefix("bac-")
	IDPrefixLinkedBankAccount    = IDPrefix("lba-")
	IDPrefixDebitCard            = IDPrefix("dbc-")
	IDPrefixLinkedDebitCard      = IDPrefix("ldc-")
	IDPrefixLinkedPayee          = IDPrefix("lpy-")
	IDPrefixLinkedCard           = IDPrefix("lca-")
	IDPrefixBankTransfer         = IDPrefix("btr-")
	IDPrefixParticipant          = IDPrefix("ppt-")
	IDPrefixBankAccountBlock     = IDPrefix("bab-")
	IDPrefixBankValidatorFailure = IDPrefix("bvf-")
	IDPrefixBankTransaction      = IDPrefix("btx-")
	IDPrefixParty                = IDPrefix("pty-")

	// Batch
	IDPrefixDailyBalance     = IDPrefix("dbl-")
	IDPrefixDailyTransaction = IDPrefix("dtx-")
	IDPrefixMonthlyInterest  = IDPrefix("mit-")

	// TODO: Deprecate
	IDPrefixPendingTransaction = IDPrefix("pnt-")
	IDPrefixPostedTransaction  = IDPrefix("pst-")

	// Event
	IDPrefixEvent       = IDPrefix("evn-")
	IDPrefixEventThread = IDPrefix("evt-")

	// Payments
	IDPrefixCardReader     = IDPrefix("cdr-")
	IDPrefixPaymentRequest = IDPrefix("pmr-")

	// Shopify
	IDPrefixShopifyBusiness         = IDPrefix("shb-")
	IDPrefixShopifyPayout           = IDPrefix("shp-")
	IDPrefixShopifyTransaction      = IDPrefix("sht-")
	IDPrefixShopifyOrder            = IDPrefix("sho-")
	IDPrefixShopifyOrderTransaction = IDPrefix("sot-")

	// Verification
	IDPrefixPhoneOTP = IDPrefix("ptp-")

	// Invoice
	IDPrefixInvoice = IDPrefix("inv-")

	// Item associated with invoice
	IDPrefixItem = IDPrefix("itm-")

	//Image
	IDPrefixImage = IDPrefix("img-")

	// Check
	IDPrefixCheck = IDPrefix("chk-")
)

func (pr IDPrefix) String() string {
	return string(pr)
}

func (pr IDPrefix) REMatch() string {
	return string(pr) + uuidRE
}

type ID interface {
	Prefix() IDPrefix
	String() string
	UUIDString() string
	IsZero() bool
	JSONString() string
}
