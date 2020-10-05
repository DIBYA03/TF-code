package push

import (
	"log"

	t "github.com/wiseco/core-platform/partner/service/twilio"
	usr "github.com/wiseco/core-platform/services/user"
	"github.com/wiseco/core-platform/shared"
)

type Notification struct {
	UserID        *shared.UserID
	Provider      TempTextProvider
	TransactionID *string
	BusinessID    *shared.BusinessID
}

func Notify(n Notification, isPosted bool) {
	// send sms
	go sendSMS(n)

	// send push notification
	go sendPushNotification(n, isPosted)

}

func sendSMS(n Notification) {
	// Get phone number
	user, err := usr.NewUserServiceWithout().GetByIdInternal(*n.UserID)
	if err != nil {
		log.Println(err)
		return
	}

	r := t.SMSRequest{
		Body:  n.Provider.PushBody,
		Phone: user.Phone,
	}

	err = t.NewTwilioService().SendSMS(r)
	if err != nil {
		log.Println(err)
	}
}

func sendPushNotification(n Notification, isPosted bool) {
	if n.TransactionID != nil {
		NewPushNotificationService().PushWithData(*n.UserID, *n.BusinessID, n.Provider, *n.TransactionID, isPosted)
	} else {
		NewPushNotificationService().PushToUser(*n.UserID, n.Provider)
	}
}
