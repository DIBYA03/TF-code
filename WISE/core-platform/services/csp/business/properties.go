/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

package business

import (
	"strings"
	"time"

	"github.com/jmoiron/sqlx/types"
	partnerbank "github.com/wiseco/core-platform/partner/bank"
	"github.com/wiseco/core-platform/services"
	"github.com/wiseco/core-platform/services/banking/business"
	bsrv "github.com/wiseco/core-platform/services/business"
	mbsrv "github.com/wiseco/core-platform/services/business"
	"github.com/wiseco/core-platform/services/csp"
	"github.com/wiseco/core-platform/shared"
)

//CSPProcessStatus ..
type CSPProcessStatus string

const (
	//CSPStatusNotStarted ..
	CSPStatusNotStarted = CSPProcessStatus("notStarted")

	//CSPStatusApproved ..
	CSPStatusApproved = CSPProcessStatus("approved")

	//CSPStatusInReview ..
	CSPStatusInReview = CSPProcessStatus("inReview")

	//CSPStatusPendingReview  ..
	CSPStatusPendingReview = CSPProcessStatus("pendingReview")

	//CSPStatusDeclined ..
	CSPStatusDeclined = CSPProcessStatus("declined")

	//CSPStatusContinue ..
	CSPStatusContinue = CSPProcessStatus("continue")

	//CSPStatusBankReview ..
	CSPStatusBankReview = CSPProcessStatus("bankReview")

	// CSPStatusBankDeclined ..
	CSPStatusBankDeclined = CSPProcessStatus("bankDeclined")
)

//CSPApproval is the decesion made by csp
type CSPApproval struct {
	Approval CSPProcessStatus `json:"status"`
	Reason   string           `json:"reason"`
	Note     string           `json:"note"`
}

var cspprocess = map[CSPProcessStatus]CSPProcessStatus{
	CSPStatusNotStarted:    CSPStatusNotStarted,
	CSPStatusInReview:      CSPStatusInReview,
	CSPStatusApproved:      CSPStatusApproved,
	CSPStatusPendingReview: CSPStatusPendingReview,
	CSPStatusDeclined:      CSPStatusDeclined,
	CSPStatusContinue:      CSPStatusContinue,
	CSPStatusBankReview:    CSPStatusBankReview,
	CSPStatusBankDeclined:  CSPStatusBankDeclined,
}

func (v CSPProcessStatus) String() string {
	return string(v)
}

//Valid checks if the csp process status is valid
func (v CSPProcessStatus) Valid() bool {
	_, ok := cspprocess[v]
	return ok
}

//CardType ..
type CardType string

const (
	//CardTypeDebit ..
	CardTypeDebit = CardType("debit")

	//CardTypeCredit ..
	CardTypeCredit = CardType("credit")

	//CardTypePrepaid ..
	CardTypePrepaid = CardType("prepaid")

	//CardTypeSingleUse ..
	CardTypeSingleUse = CardType("singleUse")
)

//BankCardCreate  ..
type BankCardCreate struct {
	// Business ID
	BusinessID shared.BusinessID `json:"businessId"`

	// Owner id of this card
	CardholderID shared.UserID `json:"cardholderId"`

	// Bank account id
	BankAccountID string `json:"bankAccountId"`

	// Card type (debit or credit)
	CardType CardType `json:"cardType"`
}

func (c CardType) String() string {
	return string(c)
}

type BusinessUpdate struct {
	// Company legal name
	LegalName *string `json:"legalName,omitempty" db:"legal_name"`

	// Know your customer status e.g. approved, declined, etc.
	KYCStatus *services.KYCStatus `json:"kycStatus,omitempty" db:"kyc_status"`

	// Entity type e.g. llc, etc
	EntityType *string `json:"entityType,omitempty" db:"entity_type"`

	// Industry type e.g. hotels, etc
	IndustryType *string `json:"industryType,omitempty" db:"industry_type"`

	// DBA name
	DBA *services.StringArray `json:"dba,omitempty" db:"dba"`

	// Tax id number e.g. ein
	TaxID *string `json:"taxId,omitempty" db:"tax_id"`

	// Tax id type e.g. ssn, ein, etc
	TaxIDType *string `json:"taxIdType,omitempty" db:"tax_id_type"`

	// Origin or country of incorporation
	OriginCountry *string `json:"originCountry,omitempty" db:"origin_country"`

	// Origin or state of incorporation
	OriginState *string `json:"originState,omitempty" db:"origin_state"`

	// Origin or date of incorporation
	OriginDate *shared.Date `json:"originDate,omitempty" db:"origin_date"`

	// Freeform business purpose
	Purpose *string `json:"purpose,omitempty" db:"purpose"`

	// Business operations e.g. foreign, domestic, etc.
	OperationType *string `json:"operationType" db:"operation_type"`

	// Business email
	Email *string `json:"email,omitempty" db:"email"`

	// Business website
	Website *services.Website `json:"website,omitempty" db:"website"`

	// Business online info
	OnlineInfo *string `json:"onlineInfo" db:"online_info"`

	// Business phone number
	Phone *string `json:"phone,omitempty" db:"phone"`

	// Business activity type e.g. cash, check, domesticAch
	ActivityType *services.StringArray `json:"activityType,omitempty" db:"activity_type"`

	// Is this a restricted business
	IsRestricted *bool `json:"isRestricted" db:"is_restricted"`

	// Legal address
	LegalAddress *services.Address `json:"legalAddress,omitempty" db:"legal_address"`

	// Headquarter address
	HeadquarterAddress *services.Address `json:"headquarterAddress" db:"headquarter_address"`

	// Mailing address
	MailingAddress *services.Address `json:"mailingAddress,omitempty" db:"mailing_address"`
}

// Business entity object
type Business struct {
	// Business id
	ID shared.BusinessID `json:"id" db:"id"`

	//Owner business id
	OwnerID shared.UserID `json:"ownerId" db:"owner_id"`

	// Generated employer number (6-Digit Code)
	EmployerNumber string `json:"employerNumber" db:"employer_number"`

	// Company legal name
	LegalName *string `json:"legalName" db:"legal_name"`

	// DBA name
	DBA services.StringArray `json:"dba" db:"dba"`

	// Is phone verified?
	PhoneVerified bool `json:"phoneVerified" db:"phone_verified"`

	// Business activity type e.g. cash, check, domesticAch
	ActivityType services.StringArray `json:"activityType" db:"activity_type"`

	//Business Handles cash
	HandlesCash bool `json:"handlesCash" db:"handles_cash"`

	// Is email verified?
	EmailVerified bool `json:"emailVerified" db:"email_verified"`

	// Business members >= 25% owners
	Members []bsrv.BusinessMember `json:"members" db:"members"`

	// Know your customer status e.g. approved, declined, etc.
	KYCStatus services.KYCStatus `json:"kycStatus" db:"kyc_status"`

	// Entity type e.g. llc, etc
	EntityType *string `json:"entityType" db:"entity_type"`

	// Is business restricted
	IsRestricted bool `json:"isRestricted" db:"is_restricted"`

	// Industry type e.g. hotels, etc
	IndustryType *string `json:"industryType" db:"industry_type"`

	// Tax id number e.g. ein
	TaxID *string `json:"taxId" db:"tax_id"`

	// Tax id type e.g. ssn, ein, etc
	TaxIDType *string `json:"taxIdType" db:"tax_id_type"`

	// Origin or country of incorporation
	OriginCountry *string `json:"originCountry" db:"origin_country"`

	// Origin or state of incorporation
	OriginState *string `json:"originState" db:"origin_state"`

	// Origin or date of incorporation
	OriginDate *shared.Date `json:"originDate" db:"origin_date"`

	// Freeform business purpose
	Purpose *string `json:"purpose" db:"purpose"`

	// Business operations e.g. foreign, domestic, etc.
	OperationType *string `json:"operationType" db:"operation_type"`

	// Business email
	Email *string `json:"email" db:"email"`

	// Business website
	Website *services.Website `json:"website,omitempty" db:"website"`

	// Business online info
	OnlineInfo *string `json:"onlineInfo" db:"online_info"`

	// Business phone number
	Phone *string `json:"phone" db:"phone"`

	// Legal address
	LegalAddress *services.Address `json:"legalAddress" db:"legal_address"`

	// Headquarter address
	HeadquarterAddress *services.Address `json:"headquarterAddress" db:"headquarter_address"`

	// Mailing address
	MailingAddress *services.Address `json:"mailingAddress" db:"mailing_address"`

	// Formation Document Id
	FormationDocumentID *shared.BusinessDocumentID `json:"formationDocumentId" db:"formation_document_id"`

	BankID *partnerbank.BusinessBankID `json:"bankId"`

	AvailableBalance *float64 `json:"availableBalance" db:"available_balance"`

	PostedBalance *float64 `json:"postedBalance" db:"posted_balance"`

	//Deactivated
	Deactivated *time.Time `json:"deactivated" db:"deactivated"`

	//Subscription decision date
	SubscriptionDecisionDate *time.Time `json:"subscriptionDecisionDate" db:"subscription_decision_date"`

	// Subscription status
	SubscriptionStatus *services.SubscriptionStatus `json:"subscriptionStatus" db:"subscription_status"`

	// Subscription start date
	SubscriptionStartDate *shared.Date `json:"subscriptionStartDate" db:"subscription_start_date"`

	// Created timestamp
	Created time.Time `json:"created" db:"created"`

	// Modified timestamp
	Modified time.Time `json:"modified" db:"modified"`
}

// MemberUpdate ..
type MemberUpdate struct {
	// Title type
	TitleType *mbsrv.TitleType `json:"titleType,omitempty" db:"title_type"`

	// Title other
	TitleOther *string `json:"titleOther,omitempty" db:"title_other"`

	// Ownership percentage
	Ownership *int `json:"ownership,omitempty" db:"ownership"`

	// Is controlling manager
	IsControllingManager *bool `json:"isControllingManager" db:"is_controlling_manager"`

	/*
	 * Consumer Properties
	 */

	// First name
	FirstName *string `json:"firstName,omitempty" db:"first_name"`

	// Middle name
	MiddleName *string `json:"middleName,omitempty" db:"middle_name"`

	// Last name
	LastName *string `json:"lastName,omitempty" db:"last_name"`

	// Email
	Email *string `json:"email,omitempty" db:"email"`

	// Phone
	Phone string `json:"phone,omitempty" db:"phone"`

	// Date of birth
	DateOfBirth *shared.Date `json:"dateOfBirth,omitempty" db:"date_of_birth"`

	// Tax id number
	TaxID *services.TaxID `json:"taxId,omitempty" db:"tax_id"`

	// Tax id type e.g. SSN or ITIN
	TaxIDType *services.TaxIDType `json:"taxIdType,omitempty"  db:"tax_id_type"`

	// Legal address
	LegalAddress *services.Address `json:"legalAddress,omitempty" db:"legal_address"`

	// Mailing address
	MailingAddress *services.Address `json:"mailingAddress,omitempty" db:"mailing_address"`

	// Residency
	Residency *services.Residency `json:"residency,omitempty" db:"residency"`

	// List of citizenships
	CitizenshipCountries *services.StringArray `json:"citizenshipCountries,omitempty" db:"citizenship_countries"`

	// Consumer's occupation
	Occupation *string `json:"occupation,omitempty" db:"occupation"`

	// Consumer income source e.g. salary or inheritance
	IncomeType *services.StringArray `json:"incomeType,omitempty" db:"income_type"`

	// Consumer activity type
	ActivityType *services.StringArray `json:"activityType,omitempty" db:"activity_type"`
}

// CardReaderCreate ..
type CardReaderCreate struct {
	BusinessID shared.BusinessID `json:"businessId" db:"business_id"`

	Alias *string `json:"alias" db:"alias"`

	DeviceType *string `json:"deviceType" db:"device_type"`

	SerialNumber string `json:"serialNumber" db:"serial_number"`
}

// CardReaderUpdate ..
type CardReaderUpdate struct {
	Alias *string `json:"alias" db:"alias"`

	SerialNumber *string `json:"serialNumber" db:"serial_number"`

	DeviceType *string `json:"deviceType" db:"device_type"`

	LastConnected *time.Time `json:"lastConnected" db:"last_connected"`
}

// CardReader ..
type CardReader struct {
	ID shared.CardReaderID `json:"id" db:"id"`

	BusinessID shared.BusinessID `json:"businessId" db:"business_id"`

	Alias *string `json:"alias" db:"alias"`

	DeviceType *string `json:"deviceType" db:"device_type"`

	SerialNumber string `json:"serialNumber" db:"serial_number"`

	LastConnected *time.Time `json:"lastConnected" db:"last_connected"`

	Deactivated *time.Time `json:"deactivated" db:"deactivated"`

	Created time.Time `json:"created" db:"created"`

	Modified time.Time `json:"modified" db:"modified"`
}

// CSPBusiness is the business presentation on business on csp, also known as review item
type CSPBusiness struct {
	ID                   string                `json:"id" db:"id"`
	BusinessName         *string               `json:"businessName" db:"business_name"`
	BusinessID           shared.BusinessID     `json:"businessId" db:"business_id"`
	Status               csp.Status            `json:"status" db:"review_status"`
	EntityType           *string               `json:"entityType" db:"entity_type"`
	ProcessStatus        string                `json:"processStatus" db:"process_status"`
	IDVs                 *services.StringArray `json:"idvs" db:"idvs"`
	Notes                types.NullJSONText    `json:"notes" db:"notes"`
	PromoFunded          *time.Time            `json:"promoFunded" db:"promo_funded"`
	EmployeeCount        int64                 `json:"employeeCount" db:"employee_count"`
	LocationCount        int64                 `json:"locationCount" db:"location_count"`
	CustomerType         string                `json:"customerType" db:"customer_type"`
	PromoMoneyTransferID *string               `json:"promoMoneyTransderId" db:"promo_money_transfer_id"`
	AcceptCardPayment    string                `json:"acceptCardPayment" db:"accept_card_payment"`
	Description          *string               `json:"description" db:"business_description"`
	PayrollType          *string               `json:"payrollType" db:"payroll_type"`
	AccountingProvider   *string               `json:"accountingProvider" db:"accounting_provider"`
	PayrollProvider      *string               `json:"payrollProvider" db:"payroll_provider"`
	Submitted            *time.Time            `json:"submitted" db:"submitted"`
	Modified             time.Time             `json:"modified" db:"modified"`
	Created              time.Time             `json:"created" db:"created"`
	ReviewSubstatus      *csp.ReviewSubstatus  `json:"reviewSubstatus" db:"review_substatus"`
	SubscribedAgentID    *string               `json:"subscibedAgentId" db:"subscribed_agent_id"`
}

// CSPBusinessCreate ...
type CSPBusinessCreate struct {
	BusinessID    shared.BusinessID     `db:"business_id"`
	BusinessName  string                `db:"business_name"`
	EntityType    *string               `db:"entity_type"`
	ProcessStatus csp.ProcessStatus     `db:"process_status"`
	Status        csp.Status            `db:"review_status"`
	IDVs          *services.StringArray `db:"idvs"`
	Notes         *types.JSONText       `db:"notes"`
}

// CSPBusinessUpdate ..
type CSPBusinessUpdate struct {
	ProcessStatus        *csp.ProcessStatus    `db:"process_status"`
	Status               *csp.Status           `db:"review_status"`
	Submitted            *time.Time            `db:"submitted"`
	IDVs                 *services.StringArray `db:"idvs"`
	Notes                *types.JSONText       `db:"notes"`
	PromoFunded          *time.Time            `json:"promoFunded" db:"promo_funded"`
	EmployeeCount        *int64                `json:"employeeCount" db:"employee_count"`
	LocationCount        *int64                `json:"locationCount" db:"location_count"`
	CustomerType         *string               `json:"customerType" db:"customer_type"`
	Description          *string               `json:"description" db:"business_description"`
	PayrollType          *string               `json:"payrollType" db:"payroll_type"`
	AccountingProvider   *string               `json:"accountingProvider" db:"accounting_provider"`
	PayrollProvider      *string               `json:"payrollProvider" db:"payroll_provider"`
	PromoMoneyTransferID *string               `json:"promoMoneyTransferId" db:"promo_money_transfer_id"`
	AcceptCardPayment    *string               `json:"acceptCardPayment" db:"accept_card_payment"`
	ReviewSubstatus      *csp.ReviewSubstatus  `json:"reviewSubstatus" db:"review_substatus"`
}

// Member ..
type Member struct {
	mbsrv.BusinessMember
	TaxID *string `json:"taxId" db:"tax_id_unmasked"`
}

// NotesCreate ..
type NotesCreate struct {
	UserID     string `json:"userId" db:"user_id"`
	BusinessID string `json:"businessId" db:"business_id"`
	Notes      string `json:"notes" db:"notes"`
}

// NotesUpdate ..
type NotesUpdate struct {
	Notes string `json:"notes" db:"notes"`
}

// Notes ..
type Notes struct {
	ID            string    `json:"id" db:"id"`
	UserID        string    `json:"userId" db:"user_id"`
	BusinessID    string    `json:"businessId" db:"business_id"`
	Notes         string    `json:"notes" db:"notes"`
	Modified      time.Time `json:"modified" db:"modified"`
	Created       time.Time `json:"created" db:"created"`
	UserPicture   *string   `json:"userPicture" db:"picture"`
	UserFirstName *string   `json:"userFirstName" db:"first_name"`
	UserLastName  *string   `json:"userLastName" db:"last_name"`
}

func (b *Business) HasDBA() bool {
	if b == nil {
		return false
	}

	if len(b.DBA) == 0 {
		return false
	}

	return len(strings.TrimSpace(b.DBA[0])) > 0
}

// Return DBA else LegalName else empty string
func (b *Business) Name() string {
	if b == nil {
		return ""
	}

	if b.HasDBA() {
		// Use first DBA
		return b.DBA[0]
	}

	if b.LegalName != nil {
		return *b.LegalName
	}

	return ""
}

// BusinessState ..
type BusinessState struct {
	ID            string            `json:"id" db:"id"`
	BusinessID    string            `json:"businessId" db:"business_id"`
	ProcessStatus csp.ProcessStatus `json:"processStatus" db:"process_status"`
	Status        csp.Status        `json:"status" db:"review_status"`
	Created       time.Time         `json:"created" db:"created"`
}

// BusinessStateCreate ..
type BusinessStateCreate struct {
	BusinessID    string     `json:"businessId" db:"business_id"`
	ProcessStatus string     `json:"processStatus" db:"process_status"`
	Status        csp.Status `json:"status" db:"review_status"`
}

// EmailVerification  email verification response
type EmailVerification struct {
	Score   float64 `json:"score"`
	Verdict string  `json:"verdict"`
}

type KYCStatus string

const (
	KYCStatusApproved = KYCStatus("approved")
	KYCStatusDeclined = KYCStatus("declined")
	KYCStatusReview   = KYCStatus("review")
	KYCStatusUnKnown  = KYCStatus("unknown")
)

type VerificationStatus string

const (
	VerificationStatusVerified   = VerificationStatus("verified")
	VerificationStatusUnverified = VerificationStatus("unverified")
	VerificationStatusInReview   = VerificationStatus("review")
	VerificationStatusUnknown    = VerificationStatus("unknown")
)

type KYCSummary struct {
	BankPartner     VerificationStatus `json:"bankPartner"`
	IdentityPartner VerificationStatus `json:"identityPartner"`
	PhonePartner    VerificationStatus `json:"phonePartner"`
	LocationPartner VerificationStatus `json:"locationPartner"`
}

type KYCResult struct {
	Result     KYCStatus  `json:"result"`
	KYCSummary KYCSummary `json:"kycSummary"`
}

type AlloyKYCResult struct {
	AlloySummary AlloySummary `json:"summary"`
}

type AlloyOutCome string

const (
	AlloyOutComeApproved     = AlloyOutCome("Approved")
	AlloyOutComeManualReview = AlloyOutCome("Manual Review")
	AlloyOutComeDeclined     = AlloyOutCome("Declined")
)

type AlloySummary struct {
	Result  string       `json:"result"`
	Score   float64      `json:"score"`
	Outcome AlloyOutCome `json:"outcome"`
}

type PhoneType string

const (
	PhoneTypePerson   = PhoneType("person")
	PhoneTypeBusiness = PhoneType("business")
)

type EveryoneResult struct {
	Type   PhoneType `json:"type"`
	Status bool      `json:"status"`
}

type IntercomResponse struct {
	Users []User `json:"users"`
}

type User struct {
	Email        string   `json:"email"`
	LocationData Location `json:"location_data"`
}

type Location struct {
	Type          string `json:"location_data"`
	CityName      string `json:"city_name"`
	ContinentCode string `json:"NA"`
	CountryName   string `json:"country_name"`
	CountruCode   string `json:"country_code"`
	PostalCode    string `json:"postal_code"`
	RegionName    string `json:"region_name"`
	TimeZone      string `json:"timezone"`
}

var StateMap = map[string]string{
	"Alabama":        "AL",
	"Alaska":         "AK",
	"Arizona":        "AZ",
	"Arkansas":       "AR",
	"California":     "CA",
	"Colorado":       "CO",
	"Connecticut":    "CT",
	"Delaware":       "DE",
	"Florida":        "FL",
	"Georgia":        "GA",
	"Hawaii":         "HI",
	"Idaho":          "ID",
	"Illinois":       "IL",
	"Indiana":        "IN",
	"Iowa":           "IA",
	"Kansas":         "KS",
	"Kentucky":       "KY",
	"Louisiana":      "LA",
	"Maine":          "ME",
	"Maryland":       "MD",
	"Massachusetts":  "MA",
	"Michigan":       "MI",
	"Minnesota":      "MN",
	"Mississippi":    "MS",
	"Missouri":       "MO",
	"Montana":        "MT",
	"Nebraska":       "NE",
	"Nevada":         "NV",
	"New Hampshire":  "NH",
	"New Jersey":     "NJ",
	"New Mexico":     "NM",
	"New York":       "NY",
	"North Carolina": "NC",
	"North Dakota":   "ND",
	"Ohio":           "OH",
	"Oklahoma":       "OK",
	"Oregon":         "OR",
	"Pennsylvania":   "PA",
	"Rhode Island":   "RI",
	"South Carolina": "SC",
	"South Dakota":   "SD",
	"Tennessee":      "TN",
	"Texas":          "TX",
	"Utah":           "UT",
	"Vermont":        "VT",
	"Virginia":       "VA",
	"Washington":     "WA",
	"West Virginia":  "WV",
	"Wisconsin":      "WI",
	"Wyoming":        "WY",

	// Territories
	"American Samoa":                 "AS",
	"District of Columbia":           "DC",
	"Federated States of Micronesia": "FM",
	"Guam":                           "GU",
	"Marshall Islands":               "MH",
	"Northern Mariana Islands":       "MP",
	"Palau":                          "PW",
	"Puerto Rico":                    "PR",
	"Virgin Islands":                 "VI",

	// Armed Forces (AE includes Europe, Africa, Canada, and the Middle East)
	"Armed Forces Americas": "AA",
	"Armed Forces Europe":   "AE",
	"Armed Forces Pacific":  "AP",
}

type KYBResult struct {
	WiseScore  float64    `json:"wiseScore"`
	Result     KYCStatus  `json:"result"`
	KYBSummary KYBSummary `json:"kybSummary"`
}

type MiddeskResponse struct {
	Formation *Formation `json:"formation"`
	Summary   []Summary  `json:"summary"`
}

type Formation struct {
	EntityType     string `json:"entity_type"`
	FormationDate  string `json:"formation_date"`
	FormationState string `json:"formation_state"`
}

type Summary struct {
	Name    string        `json:"name"`
	Status  SummaryStatus `json:"status"`
	Message string        `json:"message"`
}

type SummaryStatus string

const (
	SummaryStatusSuccess = SummaryStatus("success")
	SummaryStatusFailure = SummaryStatus("failure")
	SummaryStatusWarning = SummaryStatus("warning")
)

type KYBSummary struct {
	Address    VerificationStatus `json:"address"`
	Tin        VerificationStatus `json:"tin"`
	Formation  VerificationStatus `json:"formation"`
	Watchlist  WatchListStatus    `json:"watchlist"`
	InBusiness VerificationStatus `json:"inBusiness"`
}

type WatchListStatus string

const (
	WatchListStatusFound    = WatchListStatus("found")
	WatchListStatusNotFound = WatchListStatus("notFound")
	WatchListStatusUnKnown  = WatchListStatus("unknown")
)

const (
	TinMatchWeightage           = 10
	AddressMatchWeightage       = 10
	WatchlistNotFoundWeightage  = 10
	FormationMatchWeightage     = 10
	MaxYearsInBusinessWeightage = 10
)

type ExternalBankAccount struct {
	business.ExternalBankAccount
	Owners []business.ExternalBankAccountOwner `json:"owners"`
}

type Subscription struct {
	UserID                     shared.UserID                `json:"userId"`
	UserSubscriptionStatus     *services.SubscriptionStatus `json:"userSubscriptionStatus"`
	BusinessID                 shared.BusinessID            `json:"businessId"`
	BusinessSubscriptionStatus *services.SubscriptionStatus `json:"businessSubscriptionStatus"`
	SubscriptionStartDate      *shared.Date                 `json:"subscriptionStartDate"`
	SubscriptionDecisionDate   *time.Time                   `json:"subscriptionDecisionDate"`
	SubscribedAgentName        *string                      `json:"subscribedAgentName"`
}

type SubscriptionUpdate struct {
	BusinessID            shared.BusinessID            `json:"businessId"`
	SubscriptionStatus    *services.SubscriptionStatus `json:"subscriptionStatus"`
	SubscriptionStartDate *shared.Date                 `json:"subscriptionStartDate"`
}
