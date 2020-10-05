package shared

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/google/uuid"
)

type ProcessedSQSMessage struct {
	id      *string
	receipt *string
	err     error
}

func newProcessedSQSMessageChannel(size int) chan ProcessedSQSMessage {
	return make(chan ProcessedSQSMessage, size)
}

type sqsMessageQueue struct {
	srv    *sqs.SQS
	url    *string
	isFIFO bool
}

const maxSQSMessageCount = 10

func NewSQSMessageQueueFromURL(url string, region string) (MessageQueue, error) {
	var sess *session.Session
	var err error
	if region != "" {
		sess, err = session.NewSession(&aws.Config{Region: aws.String(region)})
		// awsConfig := &aws.Config{Region: aws.String(region)}
		// sess, err = session.NewSession(awsConfig.WithLogLevel(aws.LogDebugWithHTTPBody))
	} else {
		sess, err = session.NewSession()
		// sess, err = session.NewSession(aws.NewConfig().WithLogLevel(aws.LogDebugWithHTTPBody))
	}

	if err != nil {
		return nil, err
	}

	return NewSQSMessageQueue(sqs.New(sess), url)
}

func NewSQSMessageQueue(srv *sqs.SQS, url string) (MessageQueue, error) {
	if srv == nil {
		return nil, fmt.Errorf("sqs service must be defined")
	}
	return &sqsMessageQueue{
		srv:    srv,
		url:    aws.String(url),
		isFIFO: strings.HasSuffix(url, ".fifo"),
	}, nil
}

const (
	MessageSentTimestamp   = "SentTimestamp"
	MessageDeduplicationID = "MessageDeduplicationId"
	MessageGroupID         = "MessageGroupId"
	MessageSequenceNumber  = "SequenceNumber"
)

func (q *sqsMessageQueue) ReceiveMessages(ctx context.Context, handler MessageHandler) error {
	// Set up SQS listener
	sentTimestamp := MessageSentTimestamp
	deduplicationID := MessageDeduplicationID
	groupID := MessageGroupID
	sequenceNumber := MessageSequenceNumber
	maxNumberOfMessages := int64(10)
	receiveRequestAttemptId := uuid.New().String()
	visibilityTimeout := int64(30)
	waitTimeSeconds := int64(20)

	recvIn := sqs.ReceiveMessageInput{
		AttributeNames: []*string{
			&sentTimestamp,
			&deduplicationID,
			&groupID,
			&sequenceNumber,
		},
		MaxNumberOfMessages: &maxNumberOfMessages,
		QueueUrl:            q.url,
		VisibilityTimeout:   &visibilityTimeout,
		WaitTimeSeconds:     &waitTimeSeconds,
	}
	if q.isFIFO {
		recvIn.SetReceiveRequestAttemptId(receiveRequestAttemptId)
	}

	// SQS receiver loop
	requestTime := time.Now()
	for ctx.Err() == nil {
		recvOut, err := q.srv.ReceiveMessage(&recvIn)
		if err != nil {
			log.Println("Receive error: ", err)

			ae, ok := err.(awserr.Error)
			if ok && (ae.Code() == sqs.ErrCodeOverLimit || ae.Code() == "MissingRegion") {
				return err
			} else if time.Now().Sub(requestTime).Seconds() > 300 {
				// Reset request id if over limit
				receiveRequestAttemptId := uuid.New().String()
				requestTime = time.Now()
				if q.isFIFO {
					recvIn.SetReceiveRequestAttemptId(receiveRequestAttemptId)
				}
			}

			// Retry on common error and last request < 5 minutes ago
			continue
		}

		// Set new request id on success
		receiveRequestAttemptId := uuid.New().String()
		if q.isFIFO {
			recvIn.SetReceiveRequestAttemptId(receiveRequestAttemptId)
		}

		// If no messages return
		if recvOut.Messages == nil || len(recvOut.Messages) == 0 {
			requestTime = time.Now()
			log.Println("SQS no messages received")
			continue
		}

		log.Printf("SQS records received: %d", len(recvOut.Messages))
		if q.isFIFO {
			// Group Messages by group id
			messageMap := q.groupSQSMessages(recvOut.Messages)
			if len(messageMap) == 0 {
				requestTime = time.Now()
				continue
			}

			// Process as FIFO queue with groups
			q.processGroupedSQSMessages(ctx, handler, messageMap, len(recvOut.Messages))
		} else {
			q.processSQSMessages(ctx, handler, recvOut.Messages)
		}

		requestTime = time.Now()
	}

	return ctx.Err()
}

func (q *sqsMessageQueue) groupSQSMessages(messages []*sqs.Message) map[string][]*sqs.Message {
	var ungroupedMsgs []*sqs.Message
	var messageMap = map[string][]*sqs.Message{}
	for _, message := range messages {
		groupID, ok := message.Attributes[MessageGroupID]
		if ok {
			group, ok := messageMap[*groupID]
			if ok {
				messageMap[*groupID] = append(group, message)
			} else {
				messageMap[*groupID] = []*sqs.Message{message}
			}
		} else {
			ungroupedMsgs = append(ungroupedMsgs, message)
		}
	}

	// Add ungrouped messages as a single linear group
	if len(ungroupedMsgs) > 0 {
		messageMap[uuid.New().String()] = ungroupedMsgs
	}

	return messageMap
}

func (q *sqsMessageQueue) processGroupedSQSMessages(ctx context.Context, handler MessageHandler, messageMap map[string][]*sqs.Message, totalCount int) {
	// Process as FIFO queue with groups
	var wait sync.WaitGroup
	wait.Add(len(messageMap))
	processed := newProcessedSQSMessageChannel(maxSQSMessageCount)
	for id, m := range messageMap {
		go func(id string, messages []*sqs.Message) {
			defer wait.Done()
			for _, message := range messages {
				if message.Body != nil {
					err := handler.HandleMessage(ctx, Message{*message.MessageId, &id, []byte(*message.Body), nil})
					processed <- ProcessedSQSMessage{message.MessageId, message.ReceiptHandle, err}
				} else {
					processed <- ProcessedSQSMessage{message.MessageId, message.ReceiptHandle, errors.New("message body missing")}
				}
			}
		}(id, m)
	}

	go func() {
		wait.Wait()
		close(processed)
	}()

	q.deleteProcessedSQSMessages(processed)
}

func (q *sqsMessageQueue) processSQSMessages(ctx context.Context, handler MessageHandler, messages []*sqs.Message) {
	var wait sync.WaitGroup
	wait.Add(len(messages))
	processed := newProcessedSQSMessageChannel(maxSQSMessageCount)
	for _, m := range messages {
		go func(message *sqs.Message) {
			defer wait.Done()
			if message.Body != nil {
				err := handler.HandleMessage(ctx, Message{*message.MessageId, nil, json.RawMessage(string(*message.Body)), nil})
				processed <- ProcessedSQSMessage{message.MessageId, message.ReceiptHandle, err}
			} else {
				processed <- ProcessedSQSMessage{message.MessageId, message.ReceiptHandle, errors.New("message body missing")}
			}
		}(m)
	}

	go func() {
		wait.Wait()
		close(processed)
	}()

	q.deleteProcessedSQSMessages(processed)
}

func (q *sqsMessageQueue) deleteProcessedSQSMessages(processed chan ProcessedSQSMessage) {
	// Delete processed messages from queue
	deleteEntries := []*sqs.DeleteMessageBatchRequestEntry{}
	clearEntries := []*sqs.ChangeMessageVisibilityBatchRequestEntry{}
	for pm := range processed {
		if pm.err == nil {
			entry := sqs.DeleteMessageBatchRequestEntry{Id: pm.id, ReceiptHandle: pm.receipt}
			deleteEntries = append(deleteEntries, &entry)
		} else {
			var timeout int64 = 0
			entry := sqs.ChangeMessageVisibilityBatchRequestEntry{Id: pm.id, ReceiptHandle: pm.receipt, VisibilityTimeout: &timeout}
			clearEntries = append(clearEntries, &entry)
			log.Println(pm.err)
		}
	}

	// Delete messages that have been successfully processed
	if len(deleteEntries) > 0 {
		outDelete, err := q.srv.DeleteMessageBatch(&sqs.DeleteMessageBatchInput{Entries: deleteEntries, QueueUrl: q.url})
		if err != nil {
			log.Println(err.Error())
		} else if len(outDelete.Failed) > 0 {
			// TODO: Retry messages when delete fails?
			for _, entry := range outDelete.Failed {
				log.Println(entry.Message)
			}
		}
	}

	// Clear visibility for any messages that have errored out
	if len(clearEntries) > 0 {
		_, err := q.srv.ChangeMessageVisibilityBatch(&sqs.ChangeMessageVisibilityBatchInput{Entries: clearEntries, QueueUrl: q.url})
		if err != nil {
			log.Println(err.Error())
		}
	}
}

func (q *sqsMessageQueue) SendMessages(messages []Message) (SendMessageResult, error) {
	var entries []*sqs.SendMessageBatchRequestEntry
	for _, message := range messages {
		entries = append(entries, q.sqsMessageRequestEntry(message))
	}

	input := sqs.SendMessageBatchInput{
		Entries:  entries,
		QueueUrl: q.url,
	}

	out, err := q.srv.SendMessageBatch(&input)
	if err != nil {
		log.Println(err, " | ", out)
	}

	var sid []string
	for _, entry := range out.Successful {
		if entry != nil {
			sid = append(sid, *entry.MessageId)
		}
	}

	var fid []string
	for _, entry := range out.Failed {
		if entry != nil {
			fid = append(fid, *entry.Id)
		}
	}

	return SendMessageResult{SuccessIDs: sid, FailedIDs: fid}, err
}

func (q *sqsMessageQueue) sqsMessageRequestEntry(message Message) *sqs.SendMessageBatchRequestEntry {
	body := string(message.Body)
	if q.isFIFO {
		dedupID := uuid.New().String()
		return &sqs.SendMessageBatchRequestEntry{
			DelaySeconds:           message.Delay,
			Id:                     &message.ID,
			MessageBody:            &body,
			MessageDeduplicationId: &dedupID,
			MessageGroupId:         message.GroupID,
		}
	} else {
		return &sqs.SendMessageBatchRequestEntry{
			DelaySeconds: message.Delay,
			Id:           &message.ID,
			MessageBody:  &body,
		}
	}
}

func (q *sqsMessageQueue) URL() *string {
	return q.url
}
