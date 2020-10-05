package twilio

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/wiseco/core-platform/partner/service"
)

type SMSRequest struct {
	Phone string
	Body  string
}

type twilioService struct {
	request service.APIRequest
}

type TwilioService interface {
	//-- Send sms
	SendSMS(SMSRequest) error
}

func NewTwilioService() TwilioService {
	return &twilioService{}
}

func (t *twilioService) SendSMS(r SMSRequest) error {
	// Set initial variables

	accountSID := os.Getenv("TWILIO_ACCOUNT_SID")
	if accountSID == "" {
		panic(errors.New("twilio account sid missing"))
	}

	apiKey := os.Getenv("TWILIO_API_SID")
	if apiKey == "" {
		panic(errors.New("twilio api key missing"))
	}

	apiSecret := os.Getenv("TWILIO_API_SECRET")
	if apiKey == "" {
		panic(errors.New("twilio api secret missing"))
	}

	senderPhone := os.Getenv("TWILIO_SENDER_PHONE")
	if senderPhone == "" {
		panic(errors.New("twilio sender phone missing"))
	}

	urlStr := "https://api.twilio.com/2010-04-01/Accounts/" + accountSID + "/Messages.json"

	// Build out the data for our message
	v := url.Values{}
	v.Set("To", r.Phone)
	v.Set("From", senderPhone)
	v.Set("Body", r.Body)
	rb := *strings.NewReader(v.Encode())

	// Create client
	client := &http.Client{}

	req, _ := http.NewRequest("POST", urlStr, &rb)
	req.SetBasicAuth(apiKey, apiSecret)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Make request
	resp, _ := client.Do(req)
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		log.Println("SMS sent successfully")
		return nil
	} else {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("Unable to send SMS ", err)
		}

		return errors.New(string(body))
	}

}
