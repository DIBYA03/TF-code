package mail

import (
	"errors"
	"fmt"
	"log"
	"os"

	mailer "github.com/wiseco/core-platform/partner/service/sendgrid"
	srv "github.com/wiseco/core-platform/services"
	"github.com/wiseco/core-platform/services/csp"
	"github.com/wiseco/core-platform/services/data"
	"github.com/wiseco/core-platform/shared"
)

var wiseEmail = os.Getenv("WISE_SUPPORT_EMAIL")
var wiseSenderName = os.Getenv("WISE_SUPPORT_NAME")

type user struct {
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
	Email     string `db:"email"`
}

// EmailStatus ..
type EmailStatus string

const (
	// EmailStatusReview ..
	EmailStatusReview = EmailStatus("review")

	// EmailStatusApproved  ..
	EmailStatusApproved = EmailStatus("approved")

	// EmailStatusDeclined ..
	EmailStatusDeclined = EmailStatus("declined")
)

func sendOnReview(userID shared.UserID) error {
	if wiseEmail == "" {
		wiseEmail = "chat@wise.us"
	}

	if wiseSenderName == "" {
		wiseSenderName = "Wise"
	}

	receiver := struct {
		FirstName string `db:"first_name"`
		LastName  string `db:"last_name"`
		Email     string `db:"email"`
	}{}

	receiver, err := getUser(userID)
	body := fmt.Sprintf(srv.BusinessReviewEmail, receiver.FirstName)
	if err != nil {
		return err
	}
	email := mailer.EmailRequest{
		Subject:       srv.BusinessReviewEmailSubject,
		Body:          body,
		SenderName:    wiseSenderName,
		SenderEmail:   wiseEmail,
		ReceiverEmail: receiver.Email,
		ReceiverName:  receiver.FirstName,
	}
	resp, err := mailer.NewSendGridServiceWithout().SendEmail(email)
	if err != nil {
		log.Printf("error sending review email to owner %s  %v", receiver.FirstName, err)
	}
	log.Printf("Response from sendgrid %v", resp)
	return err
}

func sendOnApproved(userID shared.UserID, businessName string, mailingAddress srv.Address) error {

	if wiseEmail == "" {
		wiseEmail = "chat@wise.us"
	}

	if wiseSenderName == "" {
		wiseSenderName = "Wise"
	}

	receiver := struct {
		FirstName string `db:"first_name"`
		LastName  string `db:"last_name"`
		Email     string `db:"email"`
	}{}

	// create mailing address block
	addressLine2 := ""
	if mailingAddress.AddressLine2 != "" {
		addressLine2 = fmt.Sprintf("%s<br />", mailingAddress.AddressLine2)
	}

	addressBlock := fmt.Sprintf(
		srv.BusinessAddressBlock,
		mailingAddress.StreetAddress,
		addressLine2,
		mailingAddress.City,
		mailingAddress.State,
		mailingAddress.PostalCode,
	)

	receiver, err := getUser(userID)
	body := fmt.Sprintf(srv.BusinessApprovedEmail, receiver.FirstName, businessName, addressBlock)
	if err != nil {
		return err
	}
	email := mailer.EmailRequest{
		Subject:       srv.BusinessApprovedEmailSubject,
		Body:          body,
		SenderName:    wiseSenderName,
		SenderEmail:   wiseEmail,
		ReceiverEmail: receiver.Email,
		ReceiverName:  receiver.FirstName,
	}
	resp, err := mailer.NewSendGridServiceWithout().SendEmail(email)
	if err != nil {
		log.Printf("error sending approval email to business %s  %v", businessName, err)
	}
	log.Printf("Response from sendgrid %v", resp)
	return err
}

func sendOnDeclined(userID shared.UserID, businessName string) error {
	if wiseEmail == "" {
		wiseEmail = "chat@wise.us"
	}

	if wiseSenderName == "" {
		wiseSenderName = "Wise"
	}
	receiver, err := getUser(userID)
	body := fmt.Sprintf(srv.BusinessDeclinedEmail, receiver.FirstName, businessName)
	if err != nil {
		return err
	}
	email := mailer.EmailRequest{
		Subject:       fmt.Sprintf(srv.BusinessDeclinedEmailSubject),
		Body:          body,
		SenderName:    wiseSenderName,
		SenderEmail:   wiseEmail,
		ReceiverEmail: receiver.Email,
		ReceiverName:  receiver.FirstName,
	}
	resp, err := mailer.NewSendGridServiceWithout().SendEmail(email)
	if err != nil {
		log.Printf("error sending decline email to business %s  %v", businessName, err)
	}
	log.Printf("Response from sendgrid %v", resp)
	return err
}

// SendEmail ..
func SendEmail(status EmailStatus, userID shared.UserID, businessName string, mailingAddress srv.Address) error {
	switch status {
	case EmailStatusReview:
		return sendOnReview(userID)
	case EmailStatusApproved:
		return sendOnApproved(userID, businessName, mailingAddress)
	case EmailStatusDeclined:
		return sendOnDeclined(userID, businessName)
	}
	return errors.New("Invalid email status")
}

// getUser ..
func getUser(id shared.UserID) (user, error) {
	var usr user
	err := data.DBRead.Get(&usr, fmt.Sprintf(`SELECT wise_user.email,consumer.first_name,consumer.last_name
	FROM wise_user
	JOIN consumer ON wise_user.consumer_id = consumer.id
	WHERE wise_user.id = '%s'
	`, id))
	if err != nil {
		log.Printf("Error getting user with id:%s  details:%v", id, err)
	}
	return usr, err
}

// WiseTeam use to email team about a status change
type WiseTeam struct {
	Name  string
	Email string
}

// EmailWiseTeam ..
func EmailWiseTeam(businessName string, status csp.KYCStatus, team []WiseTeam) {
	if wiseEmail == "" {
		wiseEmail = "chat@wise.us"
	}

	if wiseSenderName == "" {
		wiseSenderName = "Wise"
	}
	for _, member := range team {
		body := fmt.Sprintf(srv.BusinessStatusChange, member.Name, businessName, status)

		email := mailer.EmailRequest{
			Subject:       fmt.Sprintf("%s's Wise Banking application %s", businessName, status),
			Body:          body,
			SenderName:    wiseSenderName,
			SenderEmail:   wiseEmail,
			ReceiverEmail: member.Email,
			ReceiverName:  member.Name,
		}

		resp, err := mailer.NewSendGridServiceWithout().SendEmail(email)
		if err != nil {
			log.Printf("error email wise team %v", err)
		}
		log.Printf("response from sendgrid %v", resp)
	}

}
