package activity

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx/types"
	"github.com/wiseco/core-platform/shared"
)

//Type is the Activity Type type e.g `account activity`
// Used to identify major entities for push and stream
// for example: If the activity Type is `user activity`
type Type string

const (
	//TypeUser is the user activity Type
	TypeUser = Type("user")

	//TypeBusiness is the business activity Type
	TypeBusiness = Type("business")

	//TypeConsumer is the consumer Type
	TypeConsumer = Type("consumer")

	//TypeContact is the contacts Type
	TypeContact = Type("contact")

	//TypeDispute is the dispute Type
	TypeDispute = Type("dispute")

	//TypeBankAccount is the bank account activity Type
	TypeBankAccount = Type("bankAccount")

	//TypeCard is the card activity Type
	TypeCard = Type("card")

	//TypePayment is the payment activity Type
	TypePayment = Type("payment")

	//TypeSendMoney is the send money activity Type
	TypeSendMoney = Type("sendMoney")

	//TypeRequestMoney is the request money activity Type
	TypeRequestMoney = Type("requestMoney")

	//TypeFundAccountACH is the fund account ach activity Type
	TypeFundAccountACH = Type("fundAccountACH")

	//TypeTransferTransaction is the transfer transaction Type
	TypeTransferTransaction = Type("transferTransaction")

	//TypeCardTransaction is the card transaction type e.g `debitPosted`
	TypeCardTransaction = Type("cardTransaction")

	// v1.1 Activity Types
	TypeAccountOrigination = Type("accountOriginationCredit")

	TypeCardReaderCredit = Type("cardReaderCredit")

	TypeCardOnlineCredit = Type("cardOnlineCredit")

	TypeBankOnlineCredit = Type("bankOnlineCredit")

	TypeInterestTransferCredit = Type("interestTransferCredit")

	TypeWiseTransferCredit = Type("wiseTransferCredit")

	TypeACHTransferCredit = Type("achTransferCredit")

	TypeACHTransferShopifyCredit = Type("achTransferShopifyCredit")

	TypeWireTransferCredit = Type("wireTransferCredit")

	TypeCheckCredit = Type("checkCredit")

	TypeDepositCredit = Type("depositCredit")

	TypeMerchantRefundCredit = Type("merchantRefundCredit")

	TypeCardPullCredit = Type("cardPullCredit")

	TypeCardPushCredit = Type("cardPushCredit")

	TypeCardReaderPurchaseDebit = Type("cardPurchaseDebit")

	TypeCardReaderPurchaseDebitOnline = Type("cardPurchaseDebitOnline")

	TypeWiseTransferDebit = Type("wiseTransferDebit")

	TypeACHTransferDebit = Type("achTransferDebit")

	TypeACHTransferShopifyDebit = Type("achTransferShopifyDebit")

	TypeCardATMDebit = Type("cardATMDebit")

	TypeFeeDebit = Type("feeDebit")

	TypeCardPushDebit = Type("cardPushDebit")

	TypeCardPullDebit = Type("cardPullDebit")

	TypeCheckDebit = Type("checkDebit")

	TypeHoldApproved = Type("holdApproved")

	TypeHoldReleased = Type("holdReleased")

	TypeOtherCredit = Type("otherCredit")

	TypeCardVisaCredit = Type("cardVisaCredit")
)

//Action is use for any activity that
//handles status change
//e.g  card `status` changed to `activated`
type Action string

//GLobal status change
const (
	//ActionActive is the status of ActionActive or activated
	ActionActive = Action("active")

	//ActionAuthorized is the authorized status for a card
	ActionAuthorize = Action("authorize")

	//ActionAuthReverse is the authorization reversal for a card
	ActionAuthReverse = Action("authReverse")

	//StatusApproved is the status of approved
	ActionApprove = Action("approve")

	//ActionBlock is the status of block
	ActionBlock = Action("block")

	//StatusCreated is the status of create
	ActionCreate = Action("create")

	//ActionClose is the status of close
	ActionClose = Action("close")

	//StatusDeclined is the status of declined
	ActionDecline = Action("decline")

	//StatusDeleted is the status of deleted
	ActionDelete = Action("delete")

	//StatusHold is the status of a hold, for either card or account
	ActionHold = Action("hold")

	//ActionHoldReleased is the expiry of hold set on a card
	ActionHoldReleased = Action("holdReleased")

	//StatusLocked is the status of locked
	ActionLock = Action("lock")

	//StatusPending is the status of pending
	ActionPending = Action("pending")

	//StatusPostedCredit is the status of posted credit
	ActionPostedCredit = Action("postedCredit")

	//StatusInProcessCredit is the status of in process credit
	ActionInProcessCredit = Action("inProcessCredit")

	//StatusPostedDebit is the status of posted debit
	ActionPostedDebit = Action("postedDebit")

	//StatusInProcessDebit is the status of in process debit
	ActionInProcessDebit = Action("inProcessDebit")

	//StatusReviewing is the status of in review
	ActionReviewing = Action("reviewing")

	//StatusRemoved is the status of removed
	ActionRemove = Action("remove")

	//StatusUpdated is the status of updated
	ActionUpdate = Action("update")

	//AccountOriginated is the account origination status
	ActionAccountOriginated = Action("accountOriginated")
)

func (k Type) String() string {
	activities := map[Type]Type{
		TypeUser:                TypeUser,
		TypeBusiness:            TypeBusiness,
		TypeContact:             TypeContact,
		TypeBankAccount:         TypeBankAccount,
		TypeCard:                TypeCard,
		TypePayment:             TypePayment,
		TypeSendMoney:           TypeSendMoney,
		TypeRequestMoney:        TypeRequestMoney,
		TypeFundAccountACH:      TypeFundAccountACH,
		TypeTransferTransaction: TypeTransferTransaction,
		TypeCardTransaction:     TypeCardTransaction,
	}
	return string(activities[k])
}

func (s Action) String() string {
	status := map[Action]Action{
		ActionActive:       ActionActive,
		ActionAuthorize:    ActionAuthorize,
		ActionApprove:      ActionApprove,
		ActionBlock:        ActionBlock,
		ActionCreate:       ActionCreate,
		ActionClose:        ActionClose,
		ActionDecline:      ActionDecline,
		ActionDelete:       ActionDelete,
		ActionHold:         ActionHold,
		ActionLock:         ActionLock,
		ActionPending:      ActionPending,
		ActionPostedCredit: ActionPostedCredit,
		ActionPostedDebit:  ActionPostedDebit,
		ActionReviewing:    ActionReviewing,
		ActionRemove:       ActionRemove,
		ActionUpdate:       ActionUpdate,
	}
	return string(status[s])
}

//IsValid check if the Type is valid within the enum
func (k Type) IsValid() bool {
	return k.String() != ""
}

//IsValid check if the Action is valid within the enum
func (s Action) IsValid() bool {
	return s.String() != ""
}

//Activity is the activity object constructed
// with multiple types and presented to the user
// in either a push notification or the stream activity
type Activity struct {

	//ID the id of the activity
	ID string `json:"id" db:"id"`

	//EntityID the entitity id of the activity e.g `BusinessID`
	EntityID string `json:"entityId" db:"entity_id"`

	//ActivityType is the activity type Type that identifies major activities
	//e.g `PaymentActivity`
	ActivityType Type `json:"activityType" db:"activity_type"`

	//ActivitySubType relates to the changes that can happen
	// on an account or card, naming it `activititySubType` since it could easly change to
	// something totally different later.
	//It's constraint to the `Action` type right now
	Action *Action `json:"Action" db:"activity_action"`

	//Text the text of the activity.
	//Constructed by other types and presented to the user in the stream or push object.
	//This field is not saved on db since it gets constructed on fly
	Text string `json:"text"`

	//Metadata is used to hold any extra details needed to construct a text
	//It cant be forced into a concrete type since it can have different key value pairs
	//for example when a contact is added we need to tell the client the `category`
	//and the `name` of the contact they updated and so this json blob will containt
	//those details that were feeded into the text when we created the activity.
	//When fetching the activity we need to reconstruct the text and so this `category` and `name`
	// fields are needed for the text function to properly construct the original text
	Metadata types.JSONText `json:"-" db:"metadata"`

	//ResourceID e.g `trasaction_id` `contactId` etc
	ResourceID *string `json:"resourceId" db:"resource_id"`

	//Activity Date
	ActivityDate time.Time `json:"activityDate" db:"activity_date"`

	//Created the time when this activity was created
	Created time.Time `json:"created" db:"created"`
}

//Contact ..
type Contact struct {
	//ID the entity id e.g `businessId`
	EntityID string `json:"-"`

	//UserID the user to which we need to send a push notification
	UserID shared.UserID `json:"-"`

	ContactID string `json:"-"`

	//Name the name of the contact added
	Name string `json:"name"`

	//Category the category e.g `vendor`
	Category *string `json:"category"`
}

//Dispute ..
type Dispute struct {
	//ID the entity id e.g `businessId`
	EntityID string `json:"-"`

	//UserID the user to which we need to send a push notification
	UserID shared.UserID `json:"-"`

	TransactionID string `json:"-"`

	//Disputed transaction amount
	Amount string `json:"amount"`

	//Category of the dispute e.g `incorrectly charged, etc..`
	Category string `json:"category"`
}

type AddressResponse AddressRequest

// AddressRequest should follow USPS normalization practices (Example: "St" instead of "street",
// common unit designator "APT" instead of "apartment"). Periods "." are not allowed. Zip+4 is not required.
type AddressRequest struct {
	// One of a standard set of values that indicate the customer's address type.
	// POSSIBLE VALUES:
	// legal: Legal address
	// mailing: Mailing address
	Type string `json:"type"`

	// Customer's Address line 1. The valid character set consists of UPPERCASE and lowercase letters,
	// numbers (i.e., integers), apostrophes, hyphens and the SPACE characters ONLY.
	Line1 string `json:"line1"`

	// Customer's Address line 2. The valid character set consists of UPPERCASE and lowercase letters,
	// numbers (i.e., integers), apostrophes, hyphens and the SPACE characters ONLY.
	Line2 string `json:"line2,omitempty"`

	// Customer's Residential City. The valid character set consists of UPPERCASE and lowercase letters,
	// numbers (i.e., integers), apostrophes, hyphens and the SPACE characters ONLY
	City string `json:"city"`

	// USA two-letter State abbreviation as defined by the USPS format.
	State string `json:"state"`

	// Customer's Address Postal Code. Format is 5 numbers for US zip codes, or 12345-1234 for zip codes with extension.
	ZipCode string `json:"zipCode"`

	// Country
	Country string `json:"country,omitempty"`
}

type Consumer struct {
	//ID the entity id e.g `businessId`
	EntityID string `json:"-"`

	// Consumer phone number
	Phone *string `json:"phone,omitempty"`

	// Consumer email address
	Email *string `json:"email,omitempty"`

	// Consumer address
	Address *AddressResponse `json:"address,omitempty"`
}

type Business struct {
	//ID the entity id e.g `businessId`
	EntityID string `json:"-"`

	//UserID the user to which we need to send a push notification
	UserID shared.UserID `json:"-"`

	// Business phone number
	Phone *string `json:"phone,omitempty"`

	// Business email address
	Email *string `json:"email,omitempty"`

	// Business address
	Address *AddressResponse `json:"address,omitempty"`
}

//CardTransaction ..
type CardTransaction struct {
	EntityID      string        `json:"-"`
	UserID        shared.UserID `json:"-"`
	Number        string        `json:"number"`
	Amount        string        `json:"amount"`
	Merchant      string        `json:"merchant"`
	AccountNumber string        `json:"accountNumber"`
	BusinessName  string        `json:"businessName"`
}

//Card status update. EntityID is user id in case of cards
type CardStatus struct {
	BusinessName *string `json:"businessName"`
	EntityID     string  `json:"-"`
	Number       string  `json:"number"`
	Status       string  `json:"status"`
	CardID       string  `json:"cardID"`
}

//AccountTransaction ..
type AccountTransaction struct {
	//ID the entity id e.g `businessId`
	EntityID string `json:"-"`

	//UserID the user to which we need to send a push notification
	UserID shared.UserID `json:"-"`

	// Contact that money was sent to
	ContactName *string `json:"name"`

	// Amount a generic amount for any transaction activity
	Amount AccountTransactionAmount `json:"amount"`

	//Origin the origin account
	Origin string `json:"origin"`

	//Destination the destination account
	Destination string `json:"destination"`

	InterestEarnedMonth *string `json:"interestEarnedMonth"`

	BusinessName *string `json:"businessName"`

	//Created the time when this activity was created
	TransactionDate time.Time `json:"transactionDate"`
}

func (v *CardStatus) raw() []byte {
	rw, _ := json.Marshal(v)
	return rw
}

func (v *Business) raw() []byte {
	rw, _ := json.Marshal(v)
	return rw
}

func (v *Consumer) raw() []byte {
	rw, _ := json.Marshal(v)
	return rw
}

func (v *Contact) raw() []byte {
	rw, _ := json.Marshal(v)
	return rw
}

func (v *Dispute) raw() []byte {
	rw, _ := json.Marshal(v)
	return rw
}

func (v *CardTransaction) raw() []byte {
	rw, _ := json.Marshal(v)
	return rw
}

func (v *AccountTransaction) raw() []byte {
	rw, _ := json.Marshal(v)
	return rw
}

//ComposeActivityList - adds text to both user and business activities
func ComposeActivityList(list []Activity, lang string) (*[]Activity, error) {
	var updatedList = make([]Activity, len(list))
	for i, a := range list {
		activity, _ := ComposeActivity(a, lang)
		a = *activity
		updatedList[i] = a
	}
	return &updatedList, nil
}

//ComposeActivity composes text for all activities
func ComposeActivity(a Activity, lang string) (*Activity, error) {
	switch a.ActivityType {
	case TypeContact:
		txt, err := NewContactCreator().text(a, Language(lang))
		a.Text = txt
		return &a, err
	case TypeDispute:
		txt, err := NewDisputeCreator().text(a, Language(lang))
		a.Text = txt
		return &a, err
	case TypeCardReaderPurchaseDebit:
		fallthrough
	case TypeCardReaderPurchaseDebitOnline:
		fallthrough
	case TypeCardATMDebit:
		fallthrough
	case TypeMerchantRefundCredit:
		fallthrough
	case TypeCardPushCredit:
		fallthrough
	case TypeCardVisaCredit:
		fallthrough
	case TypeCardTransaction:
		txt, err := NewCardTransationCreator().text(a, Language(lang))
		a.Text = txt
		return &a, err
	case TypeTransferTransaction:
		fallthrough
	case TypeWiseTransferDebit:
		fallthrough
	case TypeWiseTransferCredit:
		fallthrough
	case TypeCardOnlineCredit:
		fallthrough
	case TypeBankOnlineCredit:
		fallthrough
	case TypeCardReaderCredit:
		fallthrough
	case TypeACHTransferShopifyCredit:
		fallthrough
	case TypeACHTransferCredit:
		fallthrough
	case TypeWireTransferCredit:
		fallthrough
	case TypeCheckCredit:
		fallthrough
	case TypeDepositCredit:
		fallthrough
	case TypeInterestTransferCredit:
		fallthrough
	case TypeCardPullCredit:
		fallthrough
	case TypeOtherCredit:
		fallthrough
	case TypeCardPushDebit:
		fallthrough
	case TypeCardPullDebit:
		fallthrough
	case TypeCheckDebit:
		fallthrough
	case TypeACHTransferShopifyDebit:
		fallthrough
	case TypeACHTransferDebit:
		fallthrough
	case TypeHoldApproved:
		fallthrough
	case TypeHoldReleased:
		fallthrough
	case TypeFeeDebit:
		txt, err := NewTransferCreator().text(a, Language(lang))
		a.Text = txt
		return &a, err
	case TypeConsumer:
		txt, err := NewConsumerCreator().text(a, Language(lang))
		a.Text = txt
		return &a, err
	case TypeBusiness:
		txt, err := NewBusinessCreator().text(a, Language(lang))
		a.Text = txt
		return &a, err
	case TypeCard:
		txt, err := NewCardCreator().text(a, Language(lang))
		a.Text = txt
		return &a, err
	}

	return &a, nil
}

//UserActivityList takes an activity list and add the text
func UserActivityList(list []Activity) (*[]Activity, error) {
	return &list, nil
}

//UserActivity takes an activity and add the text to it
func UserActivity(a Activity) (*Activity, error) {
	return &a, nil
}

type AccountTransactionAmount float64

func (a AccountTransactionAmount) Format(f fmt.State, c rune) {
	f.Write([]byte(strconv.FormatFloat(math.Abs(float64(a)), 'f', 2, 64)))
}
