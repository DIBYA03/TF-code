package bot

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type Status string

const (
	StatusNew      = Status("new")
	StatusApproved = Status("approved")
	StatusDeclined = Status("declined")
	CSPUserCreated = Status("csp_user_created")
)

// PayLoad notification payload
type PayLoad struct {
	Channel string `json:"channel"`
	Text    string `json:"text"`
}

// SendNotification post a slack message to csp-notification channel
func SendNotification(status Status, name string) error {
	channel := os.Getenv("CSP_NOTIFICATION_SLACK_CHANNEL")
	url := os.Getenv("CSP_NOTIFICATION_SLACK_URL")
	if url == "" || channel == "" {
		return errors.New("Env var missing for url or channel")
	}
	p := PayLoad{Channel: channel,
		Text: makeText(status, name)}
	body, err := json.Marshal(p)
	if err != nil {
		return err
	}
	fmt.Println(body)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Printf("error sending slack message %v", err)
		return err
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response body %v", err)
	}
	log.Printf("Response from slack %s", string(respBody))
	return nil

}

func makeText(status Status, name string) (text string) {

	switch status {
	case StatusNew:
		text = fmt.Sprintf("A new application has been submitted. Business or owner name: %s", name)
	case StatusApproved:
		text = fmt.Sprintf("Application for %s has been approved", name)
	case StatusDeclined:
		text = fmt.Sprintf("Application for %s has been declined", name)
	case CSPUserCreated:
		text = fmt.Sprintf(":wise: A new CSP user signed up for access: %s", name)
	default:
		return ""
	}
	return text
}
