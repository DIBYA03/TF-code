package url

import (
	"os"
)

const (
	baseMiddeskSandboxURL = "api-sandbox.middesk.com"
	baseMiddeskProdURL    = "api.middesk.com"
)

//BuildAbsoluteForApp builds an abosolute path from protocol to params
func BuildAbsoluteForMiddesk(path string, urlParams Params) (string, error) {
	return BuildAbsolute(path, getBaseMiddeskURLForENV(), urlParams)
}

func getBaseMiddeskURLForENV() string {
	baseURL := protocol

	// Get URL by environment
	env := os.Getenv("API_ENV")
	switch env {
	case envDev:
		baseURL = baseURL + baseMiddeskSandboxURL
	case envStaging:
		baseURL = baseURL + baseMiddeskSandboxURL
	case envQA:
		baseURL = baseURL + baseMiddeskSandboxURL
	case envProd:
		baseURL = baseURL + baseMiddeskProdURL
	default:
		baseURL = baseURL + baseMiddeskSandboxURL
	}

	return baseURL
}
