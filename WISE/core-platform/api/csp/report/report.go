package business

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/wiseco/core-platform/api"

	report "github.com/wiseco/core-platform/services/csp/report"
)

//HandleCSPReportAPIRequests ...
func HandleCSPReportAPIRequests(request api.APIRequest) (api.APIResponse, error) {
	var method = strings.ToUpper(request.HTTPMethod)

	switch method {
	case http.MethodGet:
		return getReportList(request)
	default:
		return api.NotSupported(request)
	}

}

//HandleCSPReportStatisticsRequest ...
func HandleCSPReportStatisticsRequest(request api.APIRequest) (api.APIResponse, error) {
	var method = strings.ToUpper(request.HTTPMethod)

	switch method {
	case http.MethodGet:
		return getReportStatistics(request)
	default:
		return api.NotSupported(request)
	}

}

func getReportList(r api.APIRequest) (api.APIResponse, error) {
	params := report.CSPReportQueryParams{}

	params.Status = r.GetQueryParam("status")
	params.ReviewStatus = r.GetQueryParam("reviewStatus")
	params.ReviewSubstatus = r.GetQueryParam("reviewSubstatus")
	params.StartDate = r.GetQueryParam("startDate")
	params.EndDate = r.GetQueryParam("endDate")
	params.Offset = r.GetQueryParam("offset")
	params.Limit = r.GetQueryParam("limit")

	params.SortField = r.GetQueryParam("sortField")
	params.SortDirection = r.GetQueryParam("sortDirection")

	list, err := report.NewCSPService().CSPReportBusinessList(params)
	if err != nil {
		return api.InternalServerError(r, err)
	}
	jsonList, _ := json.Marshal(list)
	return api.Success(r, string(jsonList), false)
}

func getReportStatistics(r api.APIRequest) (api.APIResponse, error) {
	params := report.CSPReportQueryParams{}

	params.StartDate = r.GetQueryParam("startDate")
	params.EndDate = r.GetQueryParam("endDate")

	stats, err := report.NewCSPService().CSPReportStatistics(params)
	if err != nil {
		return api.InternalServerError(r, err)
	}
	resp, _ := json.Marshal(stats)
	return api.Success(r, string(resp), false)
}
