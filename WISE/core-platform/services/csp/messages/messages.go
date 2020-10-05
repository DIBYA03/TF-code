/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

package messages

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/wiseco/core-platform/services"
	"github.com/wiseco/core-platform/services/csp"
	"github.com/wiseco/core-platform/services/csp/airstream"
	"github.com/wiseco/core-platform/services/csp/bot"
	"github.com/wiseco/core-platform/services/csp/business"
	"github.com/wiseco/core-platform/services/csp/consumer"
	"github.com/wiseco/core-platform/services/csp/document"
	"github.com/wiseco/core-platform/services/csp/mail"
	"github.com/wiseco/core-platform/services/csp/review"
	coreDB "github.com/wiseco/core-platform/services/data"
	"github.com/wiseco/core-platform/shared"
)

// HandleSQS ..
func HandleSQS(body string) error {
	var notification Message
	err := json.Unmarshal([]byte(body), &notification)
	if err != nil {
		log.Printf("Error unmarshal sqs message body into %v error:%v", notification, err)
		return err
	}

	switch notification.Category {
	case CategoryBusinessDocument:
		err = businessDocumentNotification(notification)
	case CategoryConsumerDocument:
		err = consumerDocumentNotification(notification)
	case CategoryBusiness:
		err = handleBusinessNotification(notification)
	case CategoryConsumer:
		err = handleConsumerReviewNotification(notification)
	case CategoryAccount:
		err = handleAccountCreation(notification)
	default:
		return fmt.Errorf("CSP cant handle %s category yet, please make sure to pass a valid category", notification.Category)
	}
	return err
}

//Business document notification
func businessDocumentNotification(n Message) error {
	switch n.Action {
	case ActionUploadSingle:
		return businessSingleDocumentUpload(n)
	case ActionCopy:
		return handleCopyBusinessDocument(n)
	case ActionReUpload:
		return handleBusinessDocReUpload(n)
	}
	return multipleBusinessDocument(n)
}

//Business single document upload
func businessSingleDocumentUpload(n Message) error {
	log.Printf("Proccssing single document upload notification: %v", n)
	var body csp.BusinessSingleDocumentNotification
	err := json.Unmarshal([]byte(n.Data), &body)
	if err != nil {
		log.Printf("Error unmarshal sqs message body into %v error:%v", body, err)
		return err
	}

	return review.NewUploader(services.NoToCSPServiceRequest("")).BusinessSingle(body.BusinessID, body.DocumentID)
}

func multipleBusinessDocument(n Message) error {
	bID, err := shared.ParseBusinessID(n.EntityID)
	if err != nil {
		return err
	}

	log.Printf("Proccssing multipe document upload notification: %v", n)
	if err := review.NewUploader(services.NoToCSPServiceRequest("")).BusinessMultiple(bID); err != nil {
		return err
	}

	return nil
}

func handleBusinessNotification(n Message) error {
	log.Printf("Proccssing csp business notification: %v", n)
	switch n.Action {
	case ActionCreate:
		return createCSPBusiness(n.Data)
	}
	return fmt.Errorf("%s action not handled ", n.Action)
}

func createCSPBusiness(body []byte) error {
	var n csp.BusinessNotification
	if err := json.Unmarshal(body, &n); err != nil {
		log.Printf("error umarshalling data field on notification %v", err)
		return err
	}

	b := struct {
		OwnerID    shared.UserID     `db:"owner_id"`
		ConsumerID shared.ConsumerID `db:"consumer_id"`
	}{}
	err := coreDB.DBWrite.Get(
		&b, `
		 SELECT business.owner_id, wise_user.consumer_id
		 FROM business
		 JOIN wise_user ON business.owner_id = wise_user.id
		 WHERE business.id = $1`,
		n.BusinessID,
	)
	if err != nil {
		log.Printf("Error getting business by id to create business on csp %v", err)
		return err
	}

	// Verify all members of the business and change the status to docReview if all are approved
	hasReviews, err := review.VerifyMembers(b.OwnerID, n.BusinessID, b.ConsumerID)

	// if we get no error while verify members than we can proceed to docReviewStatus
	if err == nil && hasReviews == false {
		n.Status = csp.StatusDocReview
	} else {
		n.Status = csp.StatusMemberReview
	}

	create := business.CSPBusinessCreate{
		BusinessID:    n.BusinessID,
		BusinessName:  n.BusinessName,
		EntityType:    n.EntityType,
		ProcessStatus: n.ProcessStatus,
		Status:        n.Status,
		IDVs:          n.IDVs,
		Notes:         n.Notes,
	}

	if _, err := business.NewCSPService().CSPBusinessCreate(create); err != nil {
		log.Printf("Error creating review Item for business id: %s  details: %v", n.BusinessID, err)
		return fmt.Errorf("Error creating csp business with core business id: %s  details: %v", n.BusinessID, err)
	}

	useAirstream := false
	if useAirstream {
		err = airstream.NewService().StartKYB(n.BusinessID)
		if err != nil {
			log.Printf("Airsteram KYB failed for BusinessID: %v, err: %v", n.BusinessID, err)
		}
	}

	// sending email for review status to business owner
	if err := mail.SendEmail(mail.EmailStatusReview, b.OwnerID, n.BusinessName, services.Address{}); err != nil {
		log.Printf("Error sending review email for business name: %s  details: %v", n.BusinessName, err)

	}
	// post slack message about a new business
	bot.SendNotification(bot.StatusNew, n.BusinessName)
	return nil
}

// consumerDocumentNotification entity id here is the user id
func consumerDocumentNotification(notification Message) error {
	switch notification.Action {
	case ActionUploadSingle:
		return consumeSingleDocumentUpload(notification)
	case ActionReUpload:
		return handleReuploadConsumerDocument(notification)
	case ActionUpload:
		cID, err := shared.ParseConsumerID(notification.EntityID)
		if err != nil {
			return err
		}
		return review.NewUploader(services.NoToCSPServiceRequest("")).ConsumerMultiple(cID)
	default:
		return fmt.Errorf("Action %s not handled", notification.Action)
	}
}

func consumeSingleDocumentUpload(n Message) error {
	log.Printf("Processing consumer single document upload notification %v", n)
	var body csp.ConsumerSingleDocumentNotification
	if err := json.Unmarshal([]byte(n.Data), &body); err != nil {
		log.Printf("Error unmarshal sqs message body into %v error:%v", body, err)
		return err
	}

	return review.NewUploader(services.NoToCSPServiceRequest("")).ConsumerSingle(body.ConsumerID, body.DocumentID)
}

func handleConsumerReviewNotification(n Message) error {
	log.Printf("Proccssing csp consumer notification: %v", n)
	data := struct {
		Status       string               `json:"status"`
		IDVS         services.StringArray `json:"idvs"`
		ConsumerName *string              `json:"consumerName"`
	}{}
	err := json.Unmarshal(n.Data, &data)

	if err != nil {
		log.Printf("error unmarhalling notification body %v", err)
		return err
	}

	id, err := shared.ParseConsumerID(n.EntityID)
	if err != nil {
		return err
	}

	if n.Action == ActionUpdate {
		if err = json.Unmarshal(n.Data, &data); err != nil {
			log.Printf("error parsing notification %v", err)
			return err
		}
		_, err = consumer.New().UpdateKYC(id, data.Status)
		err = consumer.NewCSPService().UpdateStatus(id, data.Status)
		return err
	}

	if n.Action == ActionCreate {
		c := consumer.CSPConsumerCreate{
			ConsumerName: data.ConsumerName,
			ConsumerID:   id,
			Status:       data.Status,
			IDVs:         &data.IDVS,
		}
		con, err := consumer.NewCSPService().CSPConsumerCreate(c)
		if err != nil {
			fmt.Println("Unable to create CSP consumer: ", err)
			return err
		}
		conID, err := shared.ParseConsumerID(con.ID)

		useAirstream := false
		if useAirstream {
			err = airstream.NewService().StartKYC(conID)
		}

	}
	return err
}

// Create bank account and card once business has been approve by bank
func handleAccountCreation(n Message) error {
	bID, err := shared.ParseBusinessID(n.EntityID)
	if err != nil {
		return err
	}
	sr := services.NewSourceRequest()
	err = business.NewBanking(sr).CreateBankAccount(bID)
	if err != nil {
		return err
	}
	err = business.NewBanking(sr).CreateCard(bID)
	return err
}

func handleCopyBusinessDocument(n Message) error {
	bID, err := shared.ParseBusinessID(n.EntityID)
	if err != nil {
		log.Printf("Error parsing entity id to business id %s", n.EntityID)
		return err
	}
	business, err := business.New(services.SourceRequest{}).ByID(bID)
	if err != nil {
		log.Printf("Error getting business with id %v", bID)
		return err
	}
	consumerID, err := getConsumerIDFromUserID(business.OwnerID)
	if err != nil {
		log.Print("Error getting consumer ID")
		return err
	}
	err = document.NewDocumentService().CreateFromConsumer(consumerID, bID, business.OwnerID)

	return err

}

func handleBusinessDocReUpload(n Message) error {
	bID, err := shared.ParseBusinessID(n.EntityID)
	if err != nil {
		return err
	}
	return review.NewUploader(services.NoToCSPServiceRequest("")).BusinessReUpload(bID)
}

func handleReuploadConsumerDocument(notification Message) error {
	cID, err := shared.ParseConsumerID(notification.EntityID)
	if err != nil {
		return err
	}
	return review.NewUploader(services.NoToCSPServiceRequest("")).ConsumerReUpload(cID)
}
