package notification

import (
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/jmoiron/sqlx/types"
)

//Message SQS Message
type Message struct {
	EntityID string
	Category string
	Action   string
	Data     types.JSONText
}

const (
	CSPActionUpdate     = "update"
	CSPCategoryConsumer = "consumer"
	CSPCategoryBusiness = "business"
)

// SendReviewMessage ..
func sendReviewMessage(m Message) error {
	region := os.Getenv("CSP_REVIEW_SQS_REGION")
	queueURL := os.Getenv("CSP_REVIEW_SQS_URL")
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)
	if err != nil {
		log.Printf("error creating session %v", err)
		return err
	}
	svc := sqs.New(sess)
	body, _ := json.Marshal(m)

	_, err = svc.SendMessage(&sqs.SendMessageInput{
		DelaySeconds: aws.Int64(10),
		MessageBody:  aws.String(string(body)),
		QueueUrl:     &queueURL,
	})

	if err != nil {
		log.Printf("error sending review updates to CSP queue %v", err)
	}

	return err
}

// sqs message to update user kyc status
// pass consumerId and status
func sendConsumerKYCStatusChange(id, status string) error {
	data := struct {
		Status string `json:"status"`
	}{Status: status}
	b, _ := json.Marshal(&data)
	m := Message{
		EntityID: id, //the consumerId
		Action:   CSPActionUpdate,
		Category: CSPCategoryConsumer,
		Data:     b,
	}
	return sendReviewMessage(m)
}

// sqs message to update business kyc status
func sendBusinessKYCStatusChange(id, status string) error {
	data := struct {
		Status string `json:"status"`
	}{Status: status}
	b, _ := json.Marshal(&data)
	m := Message{
		EntityID: id, //the consumerId
		Action:   CSPActionUpdate,
		Category: CSPCategoryBusiness,
		Data:     b,
	}
	return sendReviewMessage(m)
}
