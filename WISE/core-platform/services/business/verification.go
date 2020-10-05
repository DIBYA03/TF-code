package business

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/jmoiron/sqlx/types"
	"github.com/wiseco/core-platform/services"
	coreDB "github.com/wiseco/core-platform/services/data"
	"github.com/wiseco/core-platform/shared"
)

//Response the response of business verification
type Response struct {
	Status        string                `json:"status"`
	ReviewItems   *services.StringArray `json:"reviewItems"`
	Notes         *types.JSONText       `json:"notes"`
	BusinessName  string                `json:"-"`
	BusinessOwner string                `json:"-"`
	EntityType    *string               `json:"-"`
}

type businessReviewNotification struct {
	BusinessID    shared.BusinessID     `json:"businessId"`
	BusinessName  string                `json:"businessName"`
	EntityType    *string               `json:"entityType" db:"entity_type"`
	ProcessStatus string                `json:"processStatus"`
	Status        string                `json:"reviewStatus"`
	IDVs          *services.StringArray `json:"idvs"`
	Notes         *types.JSONText       `json:"notes"`
}

const (
	processStatus = "initiated"
	status        = "memberReview"
	category      = "business"
	actionCreate  = "create"
)

// StartReview will start the verification review internally, client should get a quick response with status `review` back
//Update all the necessary states on review process
func StartReview(b *Business) error {
	bname := ""
	if b.LegalName != nil {
		bname = *(b.LegalName)
	} else {
		bname = getOwner(b.OwnerID)
	}

	create := businessReviewNotification{
		ProcessStatus: processStatus,
		Status:        status,
		BusinessID:    b.ID,
		BusinessName:  bname,
		EntityType:    b.EntityType,
	}

	return sendUpdateReviewProcess(b.ID, create)
}

func sendUpdateReviewProcess(businessID shared.BusinessID, item businessReviewNotification) error {
	data, err := json.Marshal(item)
	if err != nil {
		log.Printf("Error marshalling item:%v  details:%v", item, err)
		return err
	}
	m := message{
		EntityID: string(businessID),
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

func getOwner(id shared.UserID) string {
	var user = struct {
		FirstName string `db:"first_name"`
		LastName  string `db:"last_name"`
	}{}
	err := coreDB.DBRead.Get(&user, fmt.Sprintf(`SELECT consumer.first_name,consumer.last_name
	FROM wise_user
	JOIN consumer ON wise_user.consumer_id = consumer.id
	WHERE wise_user.id = '%s'
	`, id))
	if err != nil {
		log.Printf("Error getting user with id:%s  details:%v", id, err)
	}
	return user.FirstName + " " + user.LastName
}

type message struct {
	EntityID string
	Category string
	Action   string
	Data     types.JSONText
}
