package main

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/wiseco/core-platform/shared"
)

type sqsEnv struct {
	region string `json:"region"`
	url    string `json:"url"`
}

func main() {
	// BBVA SQS External
	regionSQS := os.Getenv("SQS_BBVA_REGION")
	urlSQS := os.Getenv("SQS_BBVA_URL")
	if urlSQS == "" {
		panic("BBVA sqs notification url missing")
	}

	queue, err := shared.NewSQSMessageQueueFromURL(urlSQS, regionSQS)
	if err != nil {
		panic(err)
	}

	// BBVA SNS publishes to various SQS queues
	regionSNS := os.Getenv("SNS_BBVA_REGION")
	topicARN := os.Getenv("SNS_BBVA_ARN")
	if topicARN == "" {
		panic("BBVA SNS topic missing")
	}

	var sess *session.Session
	if regionSNS != "" {
		sess, err = session.NewSession(&aws.Config{Region: aws.String(regionSNS)})
	} else {
		sess, err = session.NewSession()
	}

	if err != nil {
		panic(err)
	}

	h := handlerInfo{
		srv:      sns.New(sess),
		topicARN: topicARN,
	}

	err = queue.ReceiveMessages(context.Background(), &h)
	panic(err)
}

type handlerInfo struct {
	srv      *sns.SNS
	topicARN string
}

func (h *handlerInfo) HandleMessage(_ context.Context, m shared.Message) error {
	if m.Body == nil {
		return errors.New("empty message body")
	}

	body := string(m.Body)
	in := sns.PublishInput{
		Message:  &body,
		TopicArn: aws.String(h.topicARN),
	}

	out, err := h.srv.Publish(&in)
	if err != nil {
		return err
	}

	log.Println(*out.MessageId)
	return nil
}
