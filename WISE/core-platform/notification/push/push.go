package push

import (
	"context"
	"log"
	"os"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"github.com/jmoiron/sqlx"
	"github.com/wiseco/core-platform/services/data"
	"github.com/wiseco/core-platform/shared"
	"google.golang.org/api/option"
)

var client *messaging.Client

type userDevice struct {
	UserID   shared.UserID `db:"user_id"`
	Language string        `db:"language"`
	Token    string        `db:"token"`
}

//CustomData is an key value pair addinational
//data to be sent to along with the notitication
type CustomData map[string]interface{}

//Topic is the push notification topic to subscribe to
type Topic string

//TextProvider provides the text to push
//Passing the user preferred language
type TextProvider interface {
	Title(string) string
	Body(string) string
}

//Pusher is the push notitication service
type Pusher interface {
	PushWithData(shared.UserID, shared.BusinessID, TextProvider, string, bool) error
	PushToUser(shared.UserID, TextProvider)
	PushAPNS(APNS) error
}

type TempTextProvider struct {
	PushTitle string
	PushBody  string
}

func (t TempTextProvider) Title(string) string {
	return t.PushTitle
}

func (t TempTextProvider) Body(string) string {
	return t.PushBody
}

type push struct {
	*sqlx.DB
}

//NewPushNotificationService returns back a push notification services
// that uses FCM to send push notification
func NewPushNotificationService() Pusher {
	config := os.Getenv("FIREBASE_CONFIG")
	opt := option.WithCredentialsJSON([]byte(config))
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Printf("Error initializing app error:%v", err)
	}

	client, err = app.Messaging(context.Background())
	if err != nil {
		log.Printf("Error creating a message client error:%v", err)
	}
	return &push{data.DBWrite}
}

func (p *push) getDevice(id string) (userDevice, error) {
	var device userDevice
	err := p.Get(&device, "SELECT user_id, language, token FROM user_device WHERE user_id = $1", id)
	if err != nil {
		log.Printf("Error fetching token for user id: %s error:%v", id, err)
	}

	return device, err
}

func (p *push) getDevices(id shared.UserID) (*[]userDevice, error) {
	var devices []userDevice
	err := p.Select(&devices, "SELECT user_id, language, token FROM user_device WHERE user_id = $1", id)
	if err != nil {
		log.Printf("Error fetching token for user id: %s error:%v", id, err)
	}

	return &devices, err
}

func (p *push) push(userID shared.UserID, provider TextProvider) {
	devices, err := p.getDevices(userID)
	if err != nil {
		log.Printf("Error pushing to user id:%s", userID)
		return
	}

	for _, device := range *devices {

		msg := apnsMessage(provider.Title(device.Language), provider.Body(device.Language), device.Token)
		_, err = client.Send(context.Background(), msg)
		if err != nil {
			log.Printf("Error sending push notification to token:%s error: %v", device.Token, err)
			p.handleError(device.Token, err)
		} else {
			log.Printf("No error sending push to user")
		}

	}

}

func (p *push) PushWithData(userID shared.UserID, businessID shared.BusinessID, provider TextProvider, txnID string, isPosted bool) error {

	devices, err := p.getDevices(userID)
	if err != nil {
		log.Printf("Error pushing to user id:%s", userID)
		return err
	}

	for _, device := range *devices {

		msg := generateMessage(txnID, businessID, isPosted, provider.Title(device.Language), provider.Body(device.Language), device.Token)
		_, err = client.Send(context.Background(), msg)
		if err != nil {
			log.Printf("Error sending push notification to token:%s error: %v", device.Token, err)
			p.handleError(device.Token, err)
		} else {
			log.Printf("No error sending push to user")
		}

	}

	return err

}

//PushToUser swend push notification to user
func (p *push) PushToUser(userID shared.UserID, provider TextProvider) {
	log.Printf("sending push to user:%s title: %s body:%s", userID, provider.Title("en-US"), provider.Body("en-US"))
	p.push(userID, provider)
}

//This function will handle more than one error and
// act on them
func (p *push) handleError(token string, err error) error {
	if messaging.IsRegistrationTokenNotRegistered(err) {
		log.Printf("Error, token not registered in FCM token:%s error:%v looks like the token is not registered", token, err)
		p.Unregister(token)
		return err
	}
	return err
}

//Unregister deletes a token from db
//This method is ussually called when FCM errors out with
//`IsRegistrationTokenNotRegistered` error
func (p *push) Unregister(token string) (string, error) {
	var t string
	if err := p.Get(&t, "DELETE FROM user_device WHERE token = $1", token); err != nil {
		log.Printf("Error deleting token:%s some mismatch might have happened error:%v", token, err)
		return "", err
	}
	return t, nil
}

//Subscribe subscribes the passed token to the passed topic
func Subscribe(token string, to Topic) error {
	_, err := client.SubscribeToTopic(context.Background(), []string{token}, string(to))
	if err != nil {
		log.Printf("Error subscribing token: %s error:%v", token, err)
	}
	return err
}

//Unsubscribe unsubscribes the passed token to the passed topic
func Unsubscribe(token string, to Topic) error {
	_, err := client.UnsubscribeFromTopic(context.Background(), []string{token}, string(to))
	if err != nil {
		log.Printf("Error Unsubscribing token:%s error:%v", token, err)
	}

	return err
}
