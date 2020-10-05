package report

import (
	"time"

	csp "github.com/wiseco/core-platform/services/csp"
	"github.com/wiseco/core-platform/shared"
)

//CSPReportQueryParams ...
type CSPReportQueryParams struct {
	Status          string
	ReviewStatus    string
	ReviewSubstatus string
	StartDate       string
	EndDate         string
	Offset          string
	Limit           string
	SortField       string
	SortDirection   string
}

// CSPBusinessReportItem for Report Feature
type CSPBusinessReportItem struct {
	ID              string               `json:"id" db:"id"`
	BusinessName    *string              `json:"businessName" db:"business_name"`
	BusinessID      shared.BusinessID    `json:"businessId" db:"business_id"`
	ReviewStatus    csp.Status           `json:"reviewStatus" db:"review_status"`
	EntityType      *string              `json:"entityType" db:"entity_type"`
	ProcessStatus   string               `json:"processStatus" db:"process_status"`
	Received        time.Time            `json:"receivedOn" db:"received"`
	Approved        time.Time            `json:"approvedOn" db:"approved"`
	ReviewSubstatus *csp.ReviewSubstatus `json:"reviewSubstatus" db:"review_substatus"`
	DaysForApproval string               `json:"daysForApproval" db:"days"`
}

//StatisticsResponse  ..
type StatisticsResponse struct {
	ReceivedApplications               int64 `json:"receivedApplications" db:"received"`
	SubmittedApplications              int64 `json:"submittedApplications" db:"submitted"`
	ApprovedApplicationsOnFirstAttempt int64 `json:"approvedOnFirstAttempt" db:"approved_on_first_attempt"`
	ApprovedApplicationsOnResbumit     int64 `json:"approvedOnResubmit" db:"approved_on_resubmit"`
	PendingApplicaionsWithWise         int64 `json:"pendingWithWise" db:"pending_with_wise"`
	PendingApplicaions2to7MoreDays     int64 `json:"pendingWithWise2to7Days" db:"pending_with_wise_2to7_days"`
	PendingApplicaionsMoreThan7Days    int64 `json:"pendingWithWiseMoreThan7Days" db:"pending_with_wise_more_than_7_days"`
}
