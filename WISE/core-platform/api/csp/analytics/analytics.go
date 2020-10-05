package analytics

import (
	"encoding/json"

	"github.com/wiseco/core-platform/api"
	csp "github.com/wiseco/core-platform/services/csp/analytics"
)

func handleAnalytics(r api.APIRequest) (api.APIResponse, error) {
	metrics, err := csp.NewAnlytics().Metrics()
	if err != nil {
		return api.InternalServerError(r, err)
	}
	resp, err := json.Marshal(metrics)
	if err != nil {
		return api.InternalServerError(r, err)
	}
	return api.Success(r, string(resp), false)
}
