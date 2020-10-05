package analytics

import (
	"net/http"
	"strings"

	"github.com/wiseco/core-platform/api"
)

//HandleRequest analytics request
func HandleRequest(r api.APIRequest) (api.APIResponse, error) {
	if strings.ToUpper(r.HTTPMethod) != http.MethodGet {
		return api.NotSupported(r)
	}
	return handleAnalytics(r)
}
