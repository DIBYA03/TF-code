package csp

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
	Category Category
	Action   Action
	Data     types.JSONText
}

// SendDocumentMessage ..
func SendDocumentMessage(m Message) error {
	qURL := os.Getenv("CSP_SQS_URL")
	region := os.Getenv("CSP_SQS_REGION")
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)
	if err != nil {
		log.Printf("error creating session. Queue url:%s queue region:%s  %v", qURL, region, err)
		return err
	}
	svc := sqs.New(sess)
	body, _ := json.Marshal(m)

	_, err = svc.SendMessage(&sqs.SendMessageInput{
		DelaySeconds: aws.Int64(10),
		MessageBody:  aws.String(string(body)),
		QueueUrl:     &qURL,
	})
	if err != nil {
		log.Printf("error sending document upload. CSP queue Document Upload. queue url:%s queue region:%s  error:%v", qURL, region, err)
	}

	return err
}

// SendReviewMessage ..
func SendReviewMessage(m Message) error {
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
		DelaySeconds: aws.Int64(0),
		MessageBody:  aws.String(string(body)),
		QueueUrl:     &queueURL,
	})

	if err != nil {
		log.Printf("error sending review updates to CSP queue. queue url:%s queue region:%s  error:%v", queueURL, region, err)
	}

	return err
}

// Business document upload SQS message
func sendUploadDocuments(businessID string) error {
	m := Message{
		EntityID: businessID,
		Action:   ActionUpload,
		Category: CategoryBusinessDocument,
	}
	return SendDocumentMessage(m)
}

//Consumer document upload SQS message
func sendUploadConsumerDocument(consumerID string) error {
	m := Message{
		EntityID: consumerID,
		Action:   ActionUpload,
		Category: CategoryConsumerDocument,
	}
	return SendDocumentMessage(m)
}

func sendUpdateReviewProcess(businessID string, item BusinessNotification) error {
	data, err := json.Marshal(item)
	if err != nil {
		log.Printf("Error marshalling item:%v  details:%v", item, err)
		return err
	}
	m := Message{
		EntityID: businessID,
		Action:   ActionUpdate,
		Category: CategoryBusiness,
		Data:     data,
	}
	return SendReviewMessage(m)
}
