package report

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/wiseco/core-platform/services"
	"github.com/wiseco/core-platform/services/csp/data"
)

//CSPReportService ...
type CSPReportService interface {
	CSPReportBusinessList(params CSPReportQueryParams) ([]CSPBusinessReportItem, error)
	CSPReportStatistics(params CSPReportQueryParams) (StatisticsResponse, error)
}

// NewCSPService the service for CRUD on csp business
func NewCSPService() CSPReportService {
	return cspReportService{rdb: data.DBRead}
}

type cspReportService struct {
	rdb *sqlx.DB
}

func (s cspReportService) CSPReportBusinessList(params CSPReportQueryParams) ([]CSPBusinessReportItem, error) {
	var businesses []CSPBusinessReportItem
	reportQuery := ""
	dateClause := ""
	if params.StartDate != "" {
		dateClause += " AND b.created >= '" + params.StartDate + "'"
	}

	if params.EndDate != "" {
		dateClause += " AND b.created <= '" + params.EndDate + "'"
	}

	limit := "20"
	offset := "0"

	if params.Limit != "" {
		limit = params.Limit
	}

	if params.Offset != "" {
		offset = params.Offset
	}

	pagination := "LIMIT " + limit + " OFFSET " + offset

	sortBy := "b.created"
	sortDirection := "DESC"

	switch params.SortDirection {
	case "desc":
		sortDirection = "DESC"
	case "asc":
		sortDirection = "ASC"
	}

	switch params.SortField {
	case "businessName":
		sortBy = "b.business_name"
	case "received":
		sortBy = "b.created"
	case "approved":
		sortBy = "bs.created"
	case "days":
		sortBy = "days"
	case "reviewStatus":
		sortBy = "b.review_status"
	case "reviewSubstatus":
		sortBy = "b.review_substatus"
	}
	sort := sortBy + " " + sortDirection

	if params.Status == "approved" {
		columns := "b.id, b.business_id, b.business_name, b.created as received, b.review_status, bs.created as approved, date_part('day',age(bs.created, b.created )) + 1 AS days"
		tableJoinRule := "b.id = bs.business_id"
		clause := "bs.process_status = 'initiated' AND bs.review_status = 'bankApproved'" + dateClause
		reportQuery = fmt.Sprintf("SELECT %v FROM public.business b LEFT JOIN public.business_state bs ON %v WHERE %v ORDER BY %v %v;", columns, tableJoinRule, clause, sort, pagination)
	} else if params.Status == "pending" {
		columns := "b.id, b.business_id, b.business_name, b.created as received, b.review_status, b.review_substatus"
		clause := "b.review_status in ('review', 'inReview', 'pendingReview', 'riskReview', 'memberReview', 'docReview', 'bankReview')" + dateClause
		reportQuery = fmt.Sprintf("SELECT %v FROM public.business b WHERE %v ORDER BY %v %v;", columns, clause, sort, pagination)
	} else {
		log.Printf("Invalid status %v", params.ReviewStatus)
		return businesses, services.ErrorNotFound{}.New("Invalid status")
	}

	err := s.rdb.Select(&businesses, reportQuery)
	if err == sql.ErrNoRows {
		log.Printf("no business %v", err)
		return businesses, services.ErrorNotFound{}.New("")
	}
	return businesses, err

}

func (s cspReportService) CSPReportStatistics(params CSPReportQueryParams) (StatisticsResponse, error) {
	var sr StatisticsResponse

	start := params.StartDate
	end := params.EndDate

	currentTime := time.Now()
	twoDayBack := currentTime.AddDate(0, 0, -2).Format("2006-01-02")
	sevenDayBack := currentTime.AddDate(0, 0, -7).Format("2006-01-02")

	rbq := fmt.Sprintf(`(SELECT COUNT(*) FROM business WHERE created > '%v' AND created < '%v') AS received`, start, end) // Received business query

	sbq := fmt.Sprintf(`(SELECT COUNT( DISTINCT business_id) FROM business_state WHERE process_status = 'initiated' AND review_status = 'bankReview' AND created > '%v' AND created < '%v') AS submitted`, start, end)
	apf := `(SELECT COUNT(*) FROM business WHERE review_status IN ('bankApproved', 'training', 'trainingComplete') AND business_id IN (SELECT business_id FROM business_state WHERE process_status = 'initiated' AND review_status = 'bankReview' GROUP BY business_id HAVING COUNT(*) = 1)) AS approved_on_first_attempt` // Approved business on first query
	apr := `(SELECT COUNT(*) FROM business WHERE review_status IN ('bankApproved', 'training', 'trainingComplete') AND business_id IN (SELECT business_id FROM business_state WHERE process_status = 'initiated' AND review_status = 'bankReview' GROUP BY business_id HAVING COUNT(*) > 1)) AS approved_on_resubmit`      // Approved business on resubmit query
	pww := fmt.Sprintf(`(SELECT COUNT(*) FROM business WHERE review_status IN ('memberReview', 'docReview', 'riskReview', 'bankReview') AND review_substatus = 'wise' AND created > '%v' AND created < '%v') AS pending_with_wise`, start, end)
	pw2 := fmt.Sprintf(`(SELECT COUNT(*) FROM business WHERE review_status IN ('memberReview', 'docReview', 'riskReview', 'bankReview') AND review_substatus = 'wise' AND created < '`+twoDayBack+`' AND created > '`+sevenDayBack+`' AND created > '%v' AND created < '%v') AS pending_with_wise_2to7_days`, start, end)
	pw7 := fmt.Sprintf(`(SELECT COUNT(*) FROM business WHERE review_status IN ('memberReview', 'docReview', 'riskReview', 'bankReview') AND review_substatus = 'wise' AND created < '`+sevenDayBack+`' AND created > '%v' AND created < '%v') AS pending_with_wise_more_than_7_days`, start, end)

	queryString := fmt.Sprintf("SELECT %v, %v, %v, %v, %v, %v, %v", rbq, sbq, apf, apr, pww, pw2, pw7)

	err := s.rdb.Get(&sr, queryString)

	if err != nil {
		fmt.Println("Statistics error: ", err)
		return sr, err
	}

	return sr, err

}
