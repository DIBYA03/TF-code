package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
)

func cleanUpOldSnapshots(svc *rds.RDS, instanceID string) error {
	snapshotCountLimit, err := strconv.Atoi(os.Getenv("DB_SNAPSHOT_COUNT_LIMIT"))
	if err != nil {
		return err
	} else if snapshotCountLimit <= 0 {
		return errors.New("`DB_SNAPSHOT_COUNT_LIMIT` must be of value > 1")
	}

	// Get all manual snapshots for back up intance only
	listInput := &rds.DescribeDBSnapshotsInput{
		DBInstanceIdentifier: aws.String(instanceID),
		SnapshotType:         aws.String("manual"),
	}

	listResults, err := svc.DescribeDBSnapshots(listInput)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			log.Print(aerr.Error())
			return aerr
		}

		log.Print(err.Error())
		return err
	}

	snapshotList := listResults.DBSnapshots

	// keep for x days
	// AWS also does restore from a specific time as well, which is txns in S3
	if len(snapshotList) > snapshotCountLimit {
		snapshotRemoveCount := len(snapshotList) - snapshotCountLimit
		snapshotsToRemove := snapshotList[:snapshotRemoveCount]

		for _, snapshot := range snapshotsToRemove {
			snapshotID := snapshot.DBSnapshotIdentifier

			deleteInput := &rds.DeleteDBSnapshotInput{
				DBSnapshotIdentifier: snapshotID,
			}

			deleteResults, err := svc.DeleteDBSnapshot(deleteInput)
			if err != nil {
				if aerr, ok := err.(awserr.Error); ok {
					log.Print(aerr.Error())
				}

				log.Print(err.Error())
			}

			log.Printf("Snapshot %s status update: %s", *snapshotID, *deleteResults.DBSnapshot.Status)
		}
	}

	return nil
}

func createSnapshot(svc *rds.RDS, instanceID string) error {
	t := time.Now()
	timestamp := t.Format("20060102150405")

	input := &rds.CreateDBSnapshotInput{
		DBInstanceIdentifier: aws.String(instanceID),
		DBSnapshotIdentifier: aws.String(fmt.Sprintf("%s-%s", instanceID, timestamp)),
	}

	result, err := svc.CreateDBSnapshot(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			log.Print(aerr.Error())
			return aerr
		}

		log.Print(err.Error())
		return err
	}

	log.Print("create snapshot started:", result)
	return nil
}

// HandleRequest handles the request to the lambda function
func HandleRequest(ctx context.Context) (string, error) {
	instanceID := os.Getenv("DB_INSTANCE_IDENTIFIER")
	if instanceID == "" {
		return "", errors.New("Missing `DB_INSTANCE_IDENTIFIER` env var")
	}

	svc := rds.New(session.New())

	err := cleanUpOldSnapshots(svc, instanceID)
	if err != nil {
		log.Println(err)
		return "", err
	}

	err = createSnapshot(svc, instanceID)
	if err != nil {
		log.Println(err)
		return "", err
	}

	return "ok", nil
}

func main() {
	lambda.Start(HandleRequest)
}
