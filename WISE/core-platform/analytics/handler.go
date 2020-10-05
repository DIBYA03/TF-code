package analytics

import (
	"encoding/json"
	"log"

	"github.com/jmoiron/sqlx/types"
	b "github.com/wiseco/core-platform/services/banking"
	banking "github.com/wiseco/core-platform/services/banking/business"
	business "github.com/wiseco/core-platform/services/business"
	usr "github.com/wiseco/core-platform/services/user"
	"github.com/wiseco/core-platform/shared"
)

type Action string
type Category string

const (
	// CategoryConsumer
	CategoryConsumer = Category("consumer")

	// CategoryBusiness
	CategoryBusiness = Category("business")

	// CategoryAccount
	CategoryAccount = Category("account")

	// CategoryCard
	CategoryCard = Category("card")
)

const (
	// ActionCreate
	ActionCreate = Action("create")

	// ActionUpdate
	ActionUpdate = Action("update")

	// ActionKYC
	ActionKYC = Action("kyc")

	// ActionSubscription
	ActionSubscription = Action("subscription")

	ActionCSP = "csp"
)

type Message struct {
	UserID   shared.UserID
	Category Category
	Action   Action
	Data     types.JSONText
}

// TODO needs to be accessed by interface
func HandleMessages(body *string) error {
	var m Message

	err := json.Unmarshal([]byte(*body), &m)
	if err != nil {
		log.Printf("Error unmarshal sqs message body into %v error:%v", m, err)
		return err
	}

	switch m.Category {
	case CategoryConsumer:
		return unmarshalConsumer(m)
	case CategoryBusiness:
		return unmarshalBusiness(m)
	case CategoryAccount:
		return unmarshalAccount(m)
	case CategoryCard:
		return unmarshalCard(m)
	}

	return nil

}

func unmarshalConsumer(m Message) error {

	switch m.Action {
	case ActionCreate:
		return unmarshallConsumerCreate(m)
	case ActionUpdate:
		return unmarshallConsumerUpdate(m)
	case ActionKYC:
		return unmarshallConsumerKYC(m)
	case ActionSubscription:
		return unmarshallConsumerSubscription(m)
	}

	return nil
}

func unmarshallConsumerKYC(m Message) error {
	var user usr.UserVerificationUpdate

	err := json.Unmarshal(m.Data, &user)
	if err != nil {
		log.Printf("Error unmarshal segment message body into %v error:%v", m, err)
		return err
	}

	traits := make(map[string]interface{})

	traits[ConsumerKYCStatus] = user.KYCStatus

	err = Identify(m.UserID, traits)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil

}

func unmarshallConsumerSubscription(m Message) error {
	var u usr.User

	err := json.Unmarshal(m.Data, &u)
	if err != nil {
		log.Printf("Error unmarshal segment message body into %v error:%v", m, err)
		return err
	}

	traits := make(map[string]interface{})

	traits[SubscriptionStatus] = u.SubscriptionStatus

	err = Identify(m.UserID, traits)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func unmarshallConsumerUpdate(m Message) error {
	var user usr.User

	err := json.Unmarshal(m.Data, &user)
	if err != nil {
		log.Printf("Error unmarshal segment message body into %v error:%v", m, err)
		return err
	}

	traits := make(map[string]interface{})

	traits[ConsumerFirstName] = user.FirstName

	traits[ConsumerMiddleName] = user.MiddleName

	traits[ConsumerLastName] = user.LastName

	if user.Email != nil {
		traits[ConsumerEmail] = user.Email
	}

	if user.DateOfBirth != nil {
		traits[ConsumerDateOfBirth] = user.DateOfBirth
	}

	err = Identify(m.UserID, traits)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil

}

func unmarshallConsumerCreate(m Message) error {
	var user usr.UserAuthCreate

	err := json.Unmarshal(m.Data, &user)
	if err != nil {
		log.Printf("Error unmarshal segment message body into %v error:%v", m, err)
		return err
	}

	traits := make(map[string]interface{})

	traits[ConsumerPhone] = user.Phone
	traits[ConsumerPhoneVerified] = user.PhoneVerified

	err = Identify(m.UserID, traits)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil

}

func unmarshalBusiness(m Message) error {

	switch m.Action {
	case ActionCreate, ActionUpdate:
		return unmarshalBusinessCreate(m)
	case ActionKYC:
		return unmarshallBusinessKYC(m)
	case ActionCSP:
		return handleBusinessCSPUpdates(m)

	}

	return nil

}

func unmarshallBusinessKYC(m Message) error {
	var status string

	err := json.Unmarshal(m.Data, &status)
	if err != nil {
		log.Printf("Error unmarshal segment message body into %v error:%v", m, err)
		return err
	}

	traits := make(map[string]interface{})

	traits[BusinessKYCStatus] = status

	err = Identify(m.UserID, traits)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil

}

func unmarshalBusinessCreate(m Message) error {
	var b business.Business

	err := json.Unmarshal(m.Data, &b)
	if err != nil {
		log.Printf("Error unmarshal segment message body into %v error:%v", m, err)
		return err
	}

	traits := make(map[string]interface{})

	traits[BusinessId] = b.ID
	traits[BusinessLegalName] = b.LegalName
	traits[BusinessDBA] = b.DBA
	traits[BusinessEntityType] = b.EntityType
	traits[BusinessIndustryType] = b.IndustryType
	traits[BusinessKYCStatus] = b.KYCStatus

	err = Identify(m.UserID, traits)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil

}

func unmarshalCard(m Message) error {

	switch m.Action {
	case ActionCreate:
		return unmarshalCardCreate(m)
	case ActionUpdate:
		return unmarshalCardUpdate(m)

	}

	return nil

}

func unmarshalCardCreate(m Message) error {
	var c b.BankCardMini

	err := json.Unmarshal(m.Data, &c)
	if err != nil {
		log.Printf("Error unmarshal segment message body into %v error:%v", m, err)
		return err
	}

	traits := make(map[string]interface{})

	traits[BusinessCardId] = c.BankCardID
	traits[BusinessCardStatus] = c.CardStatus
	traits[BusinessTransactionCount] = c.DailyTransactionLimit

	err = Identify(m.UserID, traits)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil

}

func unmarshalCardUpdate(m Message) error {
	var c banking.BankCard

	err := json.Unmarshal(m.Data, &c)
	if err != nil {
		log.Printf("Error unmarshal segment message body into %v error:%v", m, err)
		return err
	}

	traits := make(map[string]interface{})

	traits[BusinessCardStatus] = c.CardStatus

	err = Identify(m.UserID, traits)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil

}

func unmarshalAccount(m Message) error {

	switch m.Action {
	case ActionCreate:
		return unmarshalAccountCreate(m)

	}

	return nil

}

func unmarshalAccountCreate(m Message) error {
	var a banking.BankAccount

	err := json.Unmarshal(m.Data, &a)
	if err != nil {
		log.Printf("Error unmarshal segment message body into %v error:%v", m, err)
		return err
	}

	traits := make(map[string]interface{})

	traits[BusinessAccountId] = a.BankAccountId
	traits[BusinessAccountType] = a.AccountType
	traits[BusinessAccountStatus] = a.AccountStatus
	traits[BusinessAccountAlias] = a.Alias
	traits[BusinessAccountBankName] = a.BankName
	traits[BusinessAccountOpened] = a.Opened

	err = Identify(m.UserID, traits)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil

}

func handleBusinessCSPUpdates(m Message) error {

	var b CSPBusinessUpdate

	err := json.Unmarshal(m.Data, &b)
	if err != nil {
		log.Printf("Error unmarshal segment message body into %v error:%v", m, err)
		return err
	}

	traits := make(map[string]interface{})
	if b.KYCBStatus != nil {
		traits[BusinessKYCStatus] = b.KYCBStatus
	}
	if b.PromoFunded != nil && b.Amount != nil {
		traits[BusinessNewPromoFunding] = b.Amount
	}
	err = Identify(m.UserID, traits)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
