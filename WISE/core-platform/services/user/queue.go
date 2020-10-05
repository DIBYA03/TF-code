package user

import (
	"encoding/json"
	"log"
	"os"

	"github.com/wiseco/core-platform/services"

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

//Action SQS Notification action
type Action string

//Category  SQS Notification category
type Category string

const (

	//CategoryConsumer ..
	CategoryConsumer = Category("consumer")
)

const (
	//ActionCreate ..
	ActionCreate = Action("create")

	// ActionUpdate ..
	ActionUpdate = Action("update")

	// ActionStatus ..
	ActionStatus = Action("status")
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
		DelaySeconds: aws.Int64(0),
		MessageBody:  aws.String(string(body)),
		QueueUrl:     &queueURL,
	})

	if err != nil {
		log.Printf("error sending sqs message to create consumer to CSP queue %v", err)
	}

	return err
}

//sendConsumerVerification ..
func sendConsumerVerification(id, name string, status services.KYCStatus, idvs services.StringArray, ac Action) error {
	data := struct {
		Status       services.KYCStatus   `json:"status"`
		IDVS         services.StringArray `json:"idvs"`
		ConsumerName *string              `json:"consumerName"`
	}{
		Status:       status,
		IDVS:         idvs,
		ConsumerName: &name,
	}
	b, _ := json.Marshal(&data)
	m := Message{
		EntityID: id, //the consumerId
		Action:   ac,
		Category: CategoryConsumer,
		Data:     b,
	}
	return sendReviewMessage(m)
}
