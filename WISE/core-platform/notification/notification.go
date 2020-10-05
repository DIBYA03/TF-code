package notification

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/wiseco/core-platform/notification/activity"
	"github.com/wiseco/core-platform/partner/service/segment"
	"github.com/wiseco/core-platform/services"
	"github.com/wiseco/core-platform/shared"

	"github.com/wiseco/core-platform/notification/push"

	"github.com/wiseco/core-platform/services/banking"

	"github.com/wiseco/core-platform/services/data"

	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/types"
	"github.com/wiseco/core-platform/services/banking/business"
)

type user struct {
	ID           string         `json:"id" db:"id"`
	Notification types.JSONText `json:"notification" db:"notificaton"`
}

type notificationDataStore struct {
	*sqlx.DB
}

// New returns a new notification creator
func NewNotificationService() NotificationService {
	return notificationDataStore{data.DBWrite}
}

// Service the creation of all notification
type NotificationService interface {

	// Gets user details from entity
	getUserIDByEntity(entityID string, entityType EntityType) (*shared.UserID, error)

	// Get business ID from user ID
	getBusinessIDByEntity(entityID string, entityType EntityType, accountID *string) (*shared.BusinessID, error)
}

func (s notificationDataStore) getUserIDByEntity(entityID string, entityType EntityType) (*shared.UserID, error) {
	var id shared.UserID
	switch entityType {
	case EntityTypeBusiness:
		bID, err := shared.ParseBusinessID(entityID)
		if err != nil {
			return nil, err
		}

		err = s.Get(&id, "SELECT id FROM wise_user WHERE id = (SELECT owner_id FROM business WHERE id = $1)", bID)
		if err != nil {
			return nil, fmt.Errorf("error finding owner at business id: %s", bID)
		}
	case EntityTypeConsumer:
		cID, err := shared.ParseConsumerID(entityID)
		if err != nil {
			return nil, err
		}

		err = s.Get(&id, "SELECT id FROM wise_user WHERE consumer_id = $1", cID)
		if err != nil {
			return nil, fmt.Errorf("error finding user at consumer id: %s", cID)
		}
	case EntityTypeMember:
		return nil, fmt.Errorf("error finding user at member id: %s", entityID)
	}

	return &id, nil
}

func (s notificationDataStore) getBusinessIDByEntity(entityID string, entityType EntityType, accountID *string) (*shared.BusinessID, error) {
	var id shared.BusinessID
	switch entityType {
	case EntityTypeBusiness:
		bID, err := shared.ParseBusinessID(entityID)
		return &bID, err
	case EntityTypeConsumer:
		cID, err := shared.ParseConsumerID(entityID)
		if err != nil {
			return nil, err
		}

		if accountID != nil {
			if os.Getenv("USE_BANKING_SERVICE") == "true" {
				bas, err := business.NewBankingAccountService()
				if err != nil {
					return nil, err
				}

				ba, err := bas.GetByBankAccountId(*accountID)
				if err != nil {
					return nil, err
				}

				return &ba.BusinessID, nil
			}

			err = s.Get(&id, "SELECT business_id FROM business_bank_account WHERE bank_account_id = $1", *accountID)
		} else {
			err = s.Get(&id, "SELECT id FROM business WHERE owner_id = (SELECT id FROM wise_user WHERE consumer_id = $1)", cID)
		}
		if err != nil {
			return nil, fmt.Errorf("error finding user at consumer id: %s", cID)
		}
	case EntityTypeMember:
		return nil, fmt.Errorf("error finding user at member id: %s", entityID)
	}

	return &id, nil
}

func (s notificationDataStore) getBusinessIDByUserID(userID string) (*string, error) {
	var uid string
	err := s.Get(&uid, "SELECT id FROM business WHERE owner_id = $1", userID)
	if err != nil {
		return nil, fmt.Errorf("error finding busness for user id: %s", userID)
	}

	return &uid, nil
}

// HandleNotification is the Notifications entry point. Here is where the raw data coming from SQS will be handle and distributed
func HandleNotification(body *string) error {
	var n Notification

	err := json.Unmarshal([]byte(*body), &n)
	if err != nil {
		return fmt.Errorf("notification error: %v", err)
	}

	// Logging notification details
	attribute := ""
	if n.Attribute != nil {
		attribute = string(*n.Attribute)
	}

	log.Println(fmt.Sprintf("SQS: Notification Type - %s : Action - %s : Attribute - %s", n.Type, n.Action, attribute))

	switch n.Type {
	case TypeTransaction:
		return processTransactionNotification(n)
	case TypeCard:
		return processCardNotification(n)
	case TypeAccount:
		return processAccountNotification(n)
	case TypeConsumer:
		return processConsumerNotification(n)
	case TypeBusiness:
		return processBusinessNotification(n)
	case TypeMoneyTransfer:
		return processMoneyTransferNotification(n)
	case TypePendingTransfer:
		return processPendingMoneyTransferNotification(n)
	default:
		log.Println(fmt.Errorf("invalid notification type: %s", n.Type))
		return nil
	}
}

func processConsumerNotification(n Notification) error {
	switch n.Action {
	case ActionUpdate:
		if n.Attribute != nil && (*n.Attribute == AttributeKYC || *n.Attribute == AttributeEmail || *n.Attribute == AttributePhone || *n.Attribute == AttributeAddress) {
			return processConsumerUpdate(n)
		}

		// Ignore other consumer notifications
		return nil
	default:
		log.Println(fmt.Errorf("invalid consumer notification action: %s", n.Action))
		return nil
	}
}

func processConsumerUpdate(n Notification) error {

	switch *n.Attribute {
	case AttributeKYC:
		var c ConsumerKYCUpdateNotification
		err := json.Unmarshal(n.Data, &c)
		if err != nil {
			return fmt.Errorf("consumer kyc update notification error: %v", err)
		}
		log.Println("Sending consumer kyc update to CSP ", n.EntityID, c.KYCStatus)

		// Send consumer KYC status to CSP
		return sendConsumerKYCStatusChange(n.EntityID, c.KYCStatus)

	default:
		// Map consumer ID to wise user ID
		uID, err := NewNotificationService().getUserIDByEntity(n.EntityID, n.EntityType)
		if err != nil {
			log.Println(err)
			return err
		}

		var c ConsumerContactUpdateNotification
		err = json.Unmarshal(n.Data, &c)
		if err != nil {
			return fmt.Errorf("consumer update notification error: %v", err)
		}

		onConsumerContactUpdate(*uID, c)

	}

	return nil
}

func onConsumerContactUpdate(uID shared.UserID, c ConsumerContactUpdateNotification) {
	var address activity.AddressResponse
	if c.Address != nil {
		address = activity.AddressResponse(*c.Address)
	}

	t := activity.Consumer{
		EntityID: string(uID),
		Email:    c.Email,
		Phone:    c.Phone,
		Address:  &address,
	}
	if t.Email != nil || t.Phone != nil || t.Address != nil {
		err := activity.NewConsumerCreator().Update(t)

		if err != nil {
			log.Println(err)
		}
	}
}

func processBusinessNotification(n Notification) error {
	switch n.Action {
	case ActionUpdate:
		if n.Attribute != nil && (*n.Attribute == "email" || *n.Attribute == "phone" || *n.Attribute == "address") {
			return processBusinessUpdate(n)
		} else if n.Attribute != nil && *n.Attribute == "kyc" {
			return processBusinessKYCUpdate(n)
		}

		return nil
	default:
		log.Println(fmt.Errorf("invalid notification action: %s", n.Action))
		return nil
	}
}

func processBusinessKYCUpdate(n Notification) error {
	var b BusinessKYCNotification
	err := json.Unmarshal(n.Data, &b)
	if err != nil {
		return fmt.Errorf("business kyc notification error: %v", err)
	}

	// Push KYC status to CSP
	sendBusinessKYCStatusChange(n.EntityID, b.KYCStatus)
	return nil
}

func processBusinessUpdate(n Notification) error {
	var b BusinessContactUpdateNotification
	err := json.Unmarshal(n.Data, &b)
	if err != nil {
		return fmt.Errorf("business update notification error: %v", err)
	}

	bID, err := shared.ParseBusinessID(n.EntityID)
	if err != nil {
		return err
	}

	onBusinessContactUpdate(bID, b)
	return nil
}

func onBusinessContactUpdate(bID shared.BusinessID, b BusinessContactUpdateNotification) {

	var address activity.AddressResponse

	if b.Address != nil {
		address = activity.AddressResponse(*b.Address)
	}

	business := activity.Business{
		EntityID: string(bID),
		Email:    b.Email,
		Phone:    b.Phone,
		Address:  &address,
	}

	if business.Email != nil || business.Phone != nil || business.Address != nil {
		err := activity.NewBusinessCreator().Update(business)

		if err != nil {
			log.Println(err)
		}
	}
}

func processAccountNotification(n Notification) error {
	log.Println("Account notification", n.Action)

	switch n.Action {
	case ActionOpen:
		return processNewAccountNotification(n)
	case ActionAdd:
		fallthrough
	case ActionRemove:
		if n.Attribute != nil && *n.Attribute == AttributeBlock {
			return processAccountBlockNotification(n)
		}

		fallthrough
	default:
		log.Println(fmt.Errorf("invalid account notification action: %s", n.Action))
		return nil
	}
}

func processAccountBlockNotification(n Notification) error {
	var b AccountBlockNotification
	err := json.Unmarshal(n.Data, &b)
	if err != nil {
		log.Println("Error unmarshalling account block notification", err)
		return fmt.Errorf("card notification error: %v", err)
	}

	// Get all account blocks
	bID, err := NewNotificationService().getBusinessIDByEntity(n.EntityID, n.EntityType, nil)
	if err != nil {
		log.Println(err)
		return err
	}

	account, err := business.NewAccountService().GetByBankAccountId(b.BankAccountID, *bID)
	if err != nil {
		log.Println(err)
		return err
	}

	blocks, err := business.NewBankAccountService(services.NewSourceRequest()).GetAllBankAccountBlocks(account.Id, *bID)
	if err != nil {
		log.Println(err)
		return err
	}

	// Iterate through blocks.
	blocked := false
	for _, block := range blocks {
		accountBlock, err := business.NewAccountService().GetByBlockID(block.BlockID.String(), account.Id)
		if err != nil && err != sql.ErrNoRows {
			log.Println(err)
			return err
		}

		if block.Type != banking.AccountBlockTypeCheck {
			if block.Status == banking.AccountBlockStatusActive {
				blocked = true

				if accountBlock == nil {
					c := banking.AccountBlockCreate{
						AccountID:      account.Id,
						BlockID:        block.BlockID.String(),
						OriginatedFrom: "Notification service",
						BlockType:      block.Type,
					}

					_, err := business.NewAccountService().CreateAccountBlock(c)
					if err != nil {
						log.Println(err)
						return err
					}
				}
			} else {
				if accountBlock != nil && accountBlock.Deactivated == nil {
					err := business.NewAccountService().DeactivateAccountBlock(accountBlock.ID)
					if err != nil {
						log.Println(err)
						return err
					}
				}
			}
		}
	}

	// If there is atleast one active block, then set account to blocked
	if blocked && account.AccountStatus != banking.BankAccountStatusBlocked {
		_, err := business.NewAccountService().UpdateBankAccountStatus(account.Id, banking.BankAccountStatusBlocked)
		if err != nil {
			log.Println(err)
			return err
		}
	}
	return nil
}

func processNewAccountNotification(n Notification) error {
	var t AccountOpenedNotification
	err := json.Unmarshal(n.Data, &t)
	if err != nil {
		return fmt.Errorf("card notification error: %v", err)
	}

	uid, err := NewNotificationService().getUserIDByEntity(n.EntityID, n.EntityType)
	if err != nil {
		return err
	}

	// Push to user
	notification := push.Notification{
		UserID:   uid,
		Provider: push.TempTextProvider{PushTitle: "Wise", PushBody: ApplicationApproved},
	}
	push.Notify(notification, false)

	return nil
}

func processCardNotification(n Notification) error {

	switch n.Action {
	case ActionActivate, ActionUpdate:
		return processUpdateCardNotification(n)
	default:
		log.Println(fmt.Errorf("invalid card notification action: %s", n.Action))
		return nil
	}
}

func processUpdateCardNotification(n Notification) error {

	var t CardStatusNotification
	err := json.Unmarshal(n.Data, &t)

	if err != nil {
		return fmt.Errorf("card notification error: %v", err)
	}

	// Map consumer ID to wise user ID
	uid, err := NewNotificationService().getUserIDByEntity(n.EntityID, n.EntityType)
	if err != nil {
		return err
	}

	// Fetch card number
	card, err := getBankCardDetails(t.BankCardID)
	if err != nil {
		return fmt.Errorf("Error fetching card number for card id %s", t.BankCardID)
	}

	// Update card status in database
	status, ok := bankCardStatusTo[t.Status]
	if !ok {
		return fmt.Errorf("Invalid debit card status %s", t.Status)
	}

	business.NewCardService(services.NewSourceRequest()).UpdateCardStatus(card.ID, *uid, string(status))

	lastFour, err := card.GetCardNumberLastFour()
	if err != nil {
		return err
	}

	businessName := shared.GetBusinessName(card.BusinessLegalName, card.DBA)

	var body string
	switch status {
	case banking.CardStatusActive:
		body = fmt.Sprintf(DebitCardActivated, businessName, lastFour)
		onCardStatusUpdate(*uid, card.ID, lastFour, businessName, CardStatusActivated)
	case banking.CardStatusBlocked:
		body = fmt.Sprintf(DebitCardBlocked, businessName, lastFour)
		onCardStatusUpdate(*uid, card.ID, lastFour, businessName, CardStatusBlocked)

		_, err := business.NewCardService(services.NewSourceRequest()).CheckBlockStatus(card.ID)
		if err != nil {
			log.Println(err)
		}
	case banking.CardStatusUnblocked:
		body = fmt.Sprintf(DebitCardUnblocked, businessName, lastFour)
		onCardStatusUpdate(*uid, card.ID, lastFour, businessName, CardStatusUnblocked)

		_, err := business.NewCardService(services.NewSourceRequest()).CheckBlockStatus(card.ID)
		if err != nil {
			log.Println(err)
		}
	case banking.CardStatusShipped:
		body = fmt.Sprintf(DebitCardShipped, businessName, lastFour)
		onCardStatusUpdate(*uid, card.ID, lastFour, businessName, CardStatusShipped)
	default:
		return nil

	}

	// Send push notification
	notification := push.Notification{
		UserID:   uid,
		Provider: push.TempTextProvider{PushTitle: "Wise", PushBody: body},
	}
	push.Notify(notification, false)

	// Send to segment
	c, err := business.NewCardService(services.NewSourceRequest()).GetByBankCardId(t.BankCardID, *uid)
	if err != nil {
		return fmt.Errorf("Error fetching card number for card id %s", t.BankCardID)
	}

	segment.NewSegmentService().PushToAnalyticsQueue(c.CardholderID, segment.CategoryCard, segment.ActionUpdate, card)
	return nil
}

func onCardStatusUpdate(uID shared.UserID, cardID string, number string, businessName string, status CardStatus) {

	card := activity.CardStatus{
		BusinessName: &businessName,
		EntityID:     string(uID),
		Status:       string(status),
		Number:       number,
		CardID:       cardID,
	}

	err := activity.NewCardCreator().StatusUpdate(card)

	if err != nil {
		log.Println(err)
	}
}

func getBankCardDetails(bankCardID string) (*BankCardDetails, error) {
	card := BankCardDetails{}

	query := `SELECT legal_name, dba, business_bank_card.id "business_bank_card.id", card_number_masked
	FROM business_bank_card
	JOIN business
	ON business.id = business_bank_card.business_id
	WHERE business_bank_card.bank_card_id = $1`

	err := data.DBRead.Get(&card, query, bankCardID)
	if err != nil {
		return nil, err
	}

	return &card, nil
}
