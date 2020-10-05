/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

package segment

import (
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/wiseco/core-platform/shared"
)

type SegmentService interface {
	// Identify
	PushToAnalyticsQueue(shared.UserID, Category, Action, interface{}) error
}

type Action string
type Category string

const (
	// CategoryConsumer
	CategoryConsumer = Category("consumer")

	// CategoryBusiness
	CategoryBusiness = Category("business")

	// CategoryAccount
	CategoryAccount = Category("account")

	// CategoryCard
	CategoryCard = Category("card")
)

const (
	// ActionCreate
	ActionCreate = Action("create")

	// ActionUpdate
	ActionUpdate = Action("update")

	// ActionKYC
	ActionKYC = Action("kyc")

	// ActionSubscription
	ActionSubscription = Action("subscription")

	ActionCSP = "csp"
)

type AnalyticsData struct {
	UserID   shared.UserID
	Category Category
	Action   Action
	Data     interface{}
}

type segmentService struct{}

func NewSegmentService() SegmentService {
	return &segmentService{}
}

func (s *segmentService) PushToAnalyticsQueue(userID shared.UserID, category Category, action Action, data interface{}) error {

	println("Analytics queue ", userID, category, action)

	body := AnalyticsData{
		UserID:   userID,
		Category: category,
		Action:   action,
		Data:     data,
	}

	out, err := json.Marshal(body)
	if err != nil {
		log.Println(err)
		return err
	}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("SQS_REGION"))},
	)

	svc := sqs.New(sess)

	// URL to our queue
	qURL := os.Getenv("SEGMENT_SQS_URL")

	_, err = svc.SendMessage(&sqs.SendMessageInput{
		DelaySeconds: aws.Int64(10),
		MessageBody:  aws.String(string(out)),
		QueueUrl:     &qURL,
	})
	if err != nil {
		log.Println(err)
		return err
	}

	println("Successfully pushed to analytics queue")

	return nil
}
