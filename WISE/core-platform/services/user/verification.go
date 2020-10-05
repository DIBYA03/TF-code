/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

// Package for all c related services
package user

import (
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/jmoiron/sqlx/types"
	"github.com/wiseco/core-platform/services"
	"github.com/wiseco/core-platform/shared"
)

type consumerReviewNotification struct {
	ConsumerID   shared.ConsumerID     `json:"consumerId"`
	ConsumerName string                `json:"consumerName"`
	Status       string                `json:"status"`
	IDVs         *services.StringArray `json:"idvs"`
	Notes        *types.JSONText       `json:"notes"`
}

const (
	reviewStatus = "submitted"
	category     = "consumer"
	actionCreate = "create"
)

// StartReview will start the verification review internally, client should get a quick response with status `review` back
//Update all the necessary states on review process
func StartReview(c *Consumer) error {

	create := consumerReviewNotification{
		ConsumerName: c.FirstName + " " + c.LastName,
		Status:       reviewStatus,
		ConsumerID:   c.ID,
	}

	useAirstream := false
	if !useAirstream {
		return nil
	}

	return sendUpdateReviewProcess(c.ID, create)
}

func sendUpdateReviewProcess(consumerID shared.ConsumerID, item consumerReviewNotification) error {
	data, err := json.Marshal(item)
	if err != nil {
		log.Printf("Error marshalling item:%v  details:%v", item, err)
		return err
	}
	m := message{
		EntityID: string(consumerID),
		Action:   actionCreate,
		Category: category,
		Data:     data,
	}
	return SendReviewMessage(m)
}

// SendReviewMessage ..
func SendReviewMessage(m message) error {
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
		log.Printf("error sending sqs message to create business to CSP queue. queue url:%s queue region:%s  error:%v", queueURL, region, err)
	}

	return err
}

type message struct {
	EntityID string
	Category string
	Action   string
	Data     types.JSONText
}
