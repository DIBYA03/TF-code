package document

import (
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/jmoiron/sqlx/types"
)

//SendUploadConsumerDocument Consumer document upload SQS message
func SendUploadConsumerDocument(userID string) error {
	m := Message{
		EntityID: userID,
		Action:   "upload",
		Category: "consumerDocument",
	}
	return SendDocumentMessage(m)
}

// SendDocumentMessage ..
func SendDocumentMessage(m Message) error {
	qURL := os.Getenv("CSP_SQS_URL")
	region := os.Getenv("CSP_SQS_REGION")
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
		DelaySeconds: aws.Int64(30),
		MessageBody:  aws.String(string(body)),
		QueueUrl:     &qURL,
	})
	log.Printf("error sending document upload to CSP queue %v", err)
	return err
}

//Message SQS Message
type Message struct {
	EntityID string
	Category string
	Action   string
	Data     types.JSONText
}
