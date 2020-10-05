package test

import (
	"time"

	"github.com/google/uuid"
	"github.com/wiseco/core-platform/api"
)

//TestRequest creates a test request
func TestRequest(method string) *api.APIRequest {
	request := &api.APIRequest{}
	request.HTTPMethod = method
	request.PathParameters = map[string]string{}
	request.Headers = map[string]string{"accept": "application/json",
		"Accept-Language": "en-us"}
	request.EnvID = "dev"
	request.RequestID = uuid.New().String()
	request.APIKey = "test-invoke-api-key"
	request.SourceIP = "76.21.112.249"
	request.UserAgent = "Go 1.x Wise Test Client"
	request.StartedAt = time.Now()

	return request
}

//TestResource is a resource use for testing
type TestResource struct {
	Method     string
	StatusCode int
	Resource   string
}
