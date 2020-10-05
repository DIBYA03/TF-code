package main

import (
	"encoding/json"
	"os"

	"github.com/google/uuid"
	"github.com/wiseco/core-platform/partner/bank/bbva"
	"github.com/wiseco/core-platform/shared"
)

func main() {
	region := os.Getenv("SQS_FIFO_REGION_BBVA_NOTIFICATION")
	url := os.Getenv("SQS_FIFO_URL_BBVA_NOTIFICATION")
	if url == "" {
		panic("BBVA notification sqs url missing")
	}

	queue, err := shared.NewSQSMessageQueueFromURL(url, region)
	if err != nil {
		panic(err)
	}

	id := uuid.New().String()
	groupID := uuid.New().String()

	a := bbva.ConsumerEntityNotification{
		CustomerID: "CO-dffeb35f-b653-4ee5-80c1-80053bfaa315",
		PersonalInfo: bbva.ConsumerPersonalInfoNotification{
			SSN:                "345282222",
			FirstName:          "Suresh",
			LastName:           "Venkatraman",
			BirthDate:          "1986-01-01 ",
			CitizenStatus:      "US_CITIZEN",
			CitizenshipCountry: "USA",
		},
		IdentityDocuments: []bbva.ConsumerIdentityDocumentNotification{
			bbva.ConsumerIdentityDocumentNotification{
				DocumentType:   "drivers_license",
				State:          "CA",
				IssueDate:      "2014-02-22",
				ExpirationDate: "2020-02-22",
				DocumentNumber: "123456789012345",
			},
		},
	}

	b, err := json.Marshal(a)
	if err != nil {
		panic(err)
	}

	n := bbva.Notification{
		EventID:           uuid.New().String(),
		EventType:         bbva.EventTypePrefix + bbva.EventTypeConsumerCreate,
		EventTypeVersion:  bbva.NotificationVersion,
		Subscriber:        "app.open.wise.pre",
		CustomerID:        json.RawMessage("\"CO-dffeb35f-b653-4ee5-80c1-80053bfaa315\""),
		CreationTimestamp: 1558025772000,
		Payload:           json.RawMessage(b),
	}

	b, err = json.Marshal(n)
	if err != nil {
		panic(err)
	}

	message := shared.Message{
		ID:      id,
		GroupID: &groupID,
		Body:    json.RawMessage(b),
	}

	_, err = queue.SendMessages([]shared.Message{message})
	if err != nil {
		panic(err)
	}
}
