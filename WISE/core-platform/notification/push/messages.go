package push

import (
	"context"
	"strconv"

	"firebase.google.com/go/messaging"
	"github.com/wiseco/core-platform/shared"
)

//APNS is the struct for Apple Push Notifications
type APNS struct {
	AlertString      string
	Alert            *APNSAlert
	Badge            *int
	ContentAvailable bool
	CustomData       CustomData
	Token            string
	Body             string
	UserID           shared.UserID
}

//APNSAlert is an alert
//that can be attached to an APN
type APNSAlert struct {
	Title    string
	SubTitle string
	Body     string
	Badge    *int
}

//PushAPNS can be use send a more customized push notification
//defeault priory is 10
func (p *push) PushAPNS(apns APNS) error {
	_, err := client.Send(context.Background(), messageWithAPNS(apns))
	return err
}

//We can customize this a lot better once we know what will be passing
func messageWithAPNS(apns APNS) *messaging.Message {
	return &messaging.Message{
		APNS: &messaging.APNSConfig{
			Headers: map[string]string{
				"apns-priority": "10",
			},
			Payload: &messaging.APNSPayload{
				Aps: &messaging.Aps{
					Alert: &messaging.ApsAlert{
						Title: apns.AlertString,
						Body:  apns.Body,
					},
					Badge: apns.Badge,
				},
				CustomData: apns.CustomData,
			},
		},
		Token: string(apns.Token),
	}
}

//Create a message from a custom payload
func generateMessage(txnID string, businessID shared.BusinessID, isPosted bool, title string, body string, token string) *messaging.Message {
	return &messaging.Message{
		APNS: &messaging.APNSConfig{
			Headers: map[string]string{
				"apns-priority": "10",
			},
			Payload: &messaging.APNSPayload{
				Aps: &messaging.Aps{
					Alert: &messaging.ApsAlert{
						Title: title,
						Body:  body,
					},
				},
				CustomData: map[string]interface{}{
					"transactionId": txnID,
					"isPosted":      isPosted,
					"businessId":    businessID.ToPrefixString(),
				},
			},
		},
		Android: &messaging.AndroidConfig{
			Notification: &messaging.AndroidNotification{
				Body:  body,
				Title: title,
			},
			Data: map[string]string{
				"transactionId": txnID,
				"isPosted":      strconv.FormatBool(isPosted),
				"businessId":    businessID.ToPrefixString(),
			},
		},
		Token: token,
	}
}

func apnsMessage(title, body, token string) *messaging.Message {
	return &messaging.Message{
		APNS: &messaging.APNSConfig{
			Headers: map[string]string{
				"apns-priority": "10",
			},
			Payload: &messaging.APNSPayload{
				Aps: &messaging.Aps{
					Alert: &messaging.ApsAlert{
						Title: title,
						Body:  body,
					},
				},
			},
		},
		Token: token,
	}
}
