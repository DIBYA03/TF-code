package transaction

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/firehose"
)

type BusinessTransactionLog struct {
	*BusinessPostedTransactionCreate
	CardDetails *BusinessCardTransactionCreate `json:"cardTransaction,omitempty"`
	HoldDetails *BusinessHoldTransactionCreate `json:"holdTransaction,omitempty"`
}

func (store transactionStore) Log(t *BusinessPostedTransactionCreate, c *BusinessCardTransactionCreate, h *BusinessHoldTransactionCreate) error {
	sess, err := session.NewSession(&aws.Config{Region: aws.String(LogProvider.Region().String())})
	if err != nil {
		return err
	}

	// Generate id for transaction record
	txn := BusinessTransactionLog{
		BusinessPostedTransactionCreate: t,
		CardDetails:                     c,
		HoldDetails:                     h,
	}

	b, err := json.Marshal(txn)
	if err != nil {
		return err
	}

	// Append new line to simplify parsing of S3 data
	record := &firehose.Record{Data: append(b, '\n')}

	name := LogProvider.StreamName()
	in := firehose.PutRecordInput{
		DeliveryStreamName: &name,
		Record:             record,
	}

	f := firehose.New(sess)
	_, err = f.PutRecord(&in)
	if err != nil {
		return err
	}

	/* TODO: Handle retries?
	   log.Printf("Kinesis put record count: %d", len(notifications))
	   log.Printf("Kinesis put record failure count: %d", *out.FailedPutCount) */
	return nil
}
