package csp

import (
	"encoding/json"

	"github.com/jmoiron/sqlx/types"
	"github.com/wiseco/core-platform/services"
	"github.com/wiseco/core-platform/shared"
)

//ProcessStatus  ..
type ProcessStatus string

//ReviewSubstatus - with values wise / consumer / bank ..
type ReviewSubstatus string

//Status ..
type Status string

// KYCStatus ..
type KYCStatus string

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

	// ActionReUpload  ..
	ActionReUpload = Action("reUpload")

	//ActionUploadSingle ..
	ActionUploadSingle = Action("uploadSingle")
	// ActionCopy will copy a document from business to consumer and vice versa
	ActionCopy = Action("copy")
)

const (
	// KYCStatusApproved ..
	KYCStatusApproved = KYCStatus("approved")

	// KYCStatusDeclined ..
	KYCStatusDeclined = KYCStatus("declined")

	// KYCStatusReview  ..
	KYCStatusReview = KYCStatus("review")
)
const (
	// StatusMemberReview member review status
	StatusMemberReview = Status("memberReview")

	// StatusDocReview document review status
	StatusDocReview = Status("docReview")

	// StatusRiskReview risk review status
	StatusRiskReview = Status("riskReview")

	// StatusApproved ..
	StatusApproved = Status("approved")

	//StatusDeclined ..
	StatusDeclined = Status("declined")

	// StatusContinue use to continue review process and not for real status
	StatusContinue = Status("continue")

	// StatusBankApproved ..
	StatusBankApproved = Status("bankApproved")

	// StatusBankReview ..
	StatusBankReview = Status("bankReview")

	// StatusBankDeclined ..
	StatusBankDeclined = Status("bankDeclined")

	// StatusTraining ..
	StatusTraining = Status("training")

	//StatusTrainingComplete ..
	StatusTrainingComplete = Status("trainingComplete")
)

var reviewStatus = map[Status]Status{
	StatusApproved:         StatusApproved,
	StatusDeclined:         StatusDeclined,
	StatusMemberReview:     StatusMemberReview,
	StatusDocReview:        StatusDocReview,
	StatusRiskReview:       StatusRiskReview,
	StatusContinue:         StatusContinue,
	StatusBankApproved:     StatusBankApproved,
	StatusBankReview:       StatusBankReview,
	StatusBankDeclined:     StatusBankDeclined,
	StatusTraining:         StatusTraining,
	StatusTrainingComplete: StatusTrainingComplete,
}

const (
	// ReviewSubstatusPendingOnBank in review pending on bank
	ReviewSubstatusPendingOnBank = ReviewSubstatus("bank")

	// ReviewSubstatusPendingOnCustomer in review pending on Consumer
	ReviewSubstatusPendingOnCustomer = ReviewSubstatus("customer")

	// ReviewSubstatusPendingOnWise in review pending on wise
	ReviewSubstatusPendingOnWise = ReviewSubstatus("wise")
)

var reviewSubstatus = map[ReviewSubstatus]ReviewSubstatus{
	ReviewSubstatusPendingOnBank:     ReviewSubstatusPendingOnBank,
	ReviewSubstatusPendingOnCustomer: ReviewSubstatusPendingOnCustomer,
	ReviewSubstatusPendingOnWise:     ReviewSubstatusPendingOnWise,
}

//ReviewSubstatusObj ...
var ReviewSubstatusObj = ReviewSubstatusPendingOnBank

//ReviewStatus ...
var ReviewStatus = StatusMemberReview

const (
	//Initiated  ..
	Initiated = ProcessStatus("initiated")

	//PendingDocumentUpload ..
	PendingDocumentUpload = ProcessStatus("documentPending")

	//DocUploaded ..
	DocUploaded = ProcessStatus("documentUploaded")

	//AccountCreated  ..
	AccountCreated = ProcessStatus("accountCreated")

	// AccountCreationFailed ..
	AccountCreationFailed = ProcessStatus("accountCreationFailed")

	//CardCreated ..
	CardCreated = ProcessStatus("cardCreated")

	//CardCreationFailed  ..
	CardCreationFailed = ProcessStatus("cardCreationFailed")
)

var reviewProcess = map[ProcessStatus]ProcessStatus{
	Initiated:             Initiated,
	DocUploaded:           DocUploaded,
	AccountCreated:        AccountCreated,
	AccountCreationFailed: AccountCreationFailed,
	CardCreated:           CardCreated,
	CardCreationFailed:    CardCreationFailed,
}

// DocumentStatus ..
type DocumentStatus string

const (
	DocumentStatusNotStarted = DocumentStatus("notStarted")
	DocumentStatusFailed     = DocumentStatus("failed")
	DocumentStatusPending    = DocumentStatus("pending")
	DocumentStatusUploaded   = DocumentStatus("uploaded")
)

//NewStatus ..
func (Status) NewStatus(s string) (Status, bool) {
	r, ok := reviewStatus[Status(s)]
	return r, ok
}

func (v Status) String() string {
	return string(v)
}

//CheckReviewSubstatus ..
func CheckReviewSubstatus(s string) bool {
	_, ok := reviewSubstatus[ReviewSubstatus(s)]
	return ok
}

func (v ReviewSubstatus) String() string {
	return string(v)
}

func (v KYCStatus) String() string {
	return string(v)
}

func (v ProcessStatus) String() string {
	return string(v)
}

//Valid checks if the csp process status is valid
func (v Status) Valid() bool {
	_, ok := reviewStatus[v]
	return ok
}

// Valid checks if the csp process status is valid
func (v ProcessStatus) Valid() bool {
	_, ok := reviewProcess[v]
	return ok
}

// KYCNote ..
type KYCNote struct {
	Code string `json:"code"`
	Desc string `json:"desc"`
}

// KYCNotes csp kyc notes
type KYCNotes []KYCNote

// Raw convert to []byte
func (v KYCNotes) Raw() types.JSONText {
	b, _ := json.Marshal(v)
	return b
}

// ReviewResponse the response of business verification
type ReviewResponse struct {
	Status        Status                `json:"status"`
	ReviewItems   *services.StringArray `json:"reviewItems"`
	Notes         *types.JSONText       `json:"notes"`
	BusinessName  string                `json:"-"`
	BusinessOwner shared.UserID         `json:"-"`
	EntityType    *string               `json:"-"`
}

// ConsumerNotification  ..
type ConsumerNotification struct {
	ConsumerName *string               `json:"consumerName" db:"consumer_name"`
	ConsumerID   shared.ConsumerID     `json:"consumerId" db:"consumer_id"`
	Status       string                `json:"status" db:"review_status"`
	IDVs         *services.StringArray `json:"idvs" db:"idvs"`
	Notes        *types.JSONText       `json:"notes" db:"notes"`
}

// BusinessNotification ..
type BusinessNotification struct {
	BusinessID    shared.BusinessID     `json:"businessId"`
	BusinessName  string                `json:"businessName"`
	EntityType    *string               `json:"entityType" db:"entity_type"`
	ProcessStatus ProcessStatus         `json:"processStatus"`
	Status        Status                `json:"reviewStatus"`
	IDVs          *services.StringArray `json:"idvs"`
	Notes         *types.JSONText       `json:"notes"`
}

// BusinessSingleDocumentNotification ..
type BusinessSingleDocumentNotification struct {
	BusinessID shared.BusinessID         `json:"businessId"`
	DocumentID shared.BusinessDocumentID `json:"documentId"`
}

// BusinessDocumentCopyNotification ..
type BusinessDocumentCopyNotification struct {
	ConsumerID shared.ConsumerID `json:"consumerId"`
}

// ConsumerDocumentCopyNotification ..
type ConsumerDocumentCopyNotification struct {
	BusinessID shared.BusinessID `json:"businessId"`
}

// ConsumerSingleDocumentNotification  ..
type ConsumerSingleDocumentNotification struct {
	ConsumerID shared.ConsumerID         `json:"consumerId"`
	DocumentID shared.ConsumerDocumentID `json:"documentId"`
}
