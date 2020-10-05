package messages

import (
	"github.com/jmoiron/sqlx/types"
	"github.com/wiseco/core-platform/services"
	"github.com/wiseco/core-platform/services/csp"
	"github.com/wiseco/core-platform/shared"
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
	// CategoryBusiness ..
	CategoryBusiness = Category("business")

	// CategoryAccount ..
	CategoryAccount = Category("account")

	// CategoryCard ..
	CategoryCard = Category("card")

	//CategoryBusinessDocument  ..
	CategoryBusinessDocument = Category("businessDocument")

	//CategoryConsumerDocument ..
	CategoryConsumerDocument = Category("consumerDocument")

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

	//ActionUpload  ..
	ActionUpload = Action("upload")

	// ActionCopy will copy a document from business to consumer and vice versa
	ActionCopy = Action("copy")
	// ActionReUpload  ..
	ActionReUpload = Action("reUpload")

	//ActionUploadSingle ..
	ActionUploadSingle = Action("uploadSingle")
)

//ConsumerNotification  ..
type ConsumerNotification struct {
	ConsumerName *string               `json:"consumerName" db:"consumer_name"`
	ConsumerID   string                `json:"consumerId" db:"consumer_id"`
	Status       string                `json:"status" db:"review_status"`
	IDVs         *services.StringArray `json:"idvs" db:"idvs"`
	Notes        *types.JSONText       `json:"notes" db:"notes"`
}

//BusinessNotification ..
type BusinessNotification struct {
	BusinessID    string                `json:"businessId"`
	BusinessName  string                `json:"businessName"`
	EntityType    *string               `json:"entityType" db:"entity_type"`
	ProcessStatus csp.ProcessStatus     `json:"processStatus"`
	Status        csp.Status            `json:"reviewStatus"`
	IDVs          *services.StringArray `json:"idvs"`
	Notes         *types.JSONText       `json:"notes"`
}

// BusinessSingleDocumentNotification ..
type BusinessSingleDocumentNotification struct {
	BusinessID string `json:"businessId"`
	DocumentID string `json:"documentId"`
}

// ConsumerSingleDocumentNotification  ..
type ConsumerSingleDocumentNotification struct {
	ConsumerID shared.ConsumerID `json:"businessId"`
	DocumentID string            `json:"documentId"`
}
