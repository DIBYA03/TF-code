package delegate

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/wiseco/core-platform/api"
	"github.com/wiseco/core-platform/delegate"
)

func execute(r api.APIRequest) (api.APIResponse, error) {
	var resource delegate.Resource
	if err := json.Unmarshal([]byte(r.Body), &resource); err != nil {
		return api.BadRequest(r, err)
	}

	resource.SourceRequest = r.SourceRequest()
	resp := delegate.NewProxyService().Execute(resource)

	if resp.Error != nil {
		log.Println("proxy err: ", resp.Error)
		return api.ProxyErrorResponse(r, resp.Body, resp.StatusCode)
	}

	return api.Success(r, string(resp.Body), false)
}

//HandleProxyRequest handles BBVA  requests
func HandleProxyRequest(r api.APIRequest) (api.APIResponse, error) {
	// Disallow for production
	switch os.Getenv("APP_ENV") {
	case "prod", "qa-prod", "beta-prod":
		return api.NotAllowedError(r, nil)
	}

	switch strings.ToUpper(r.HTTPMethod) {
	case http.MethodPost:
		return execute(r)
	default:
		return api.NotSupported(r)
	}
}
