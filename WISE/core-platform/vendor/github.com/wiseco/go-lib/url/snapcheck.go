package url

import "os"

const (
	baseSnapCheckSandboxURL = "api-sandbox.gosnapcheck.com/v3"
	baseSnapCheckProdURL    = "api-prod.gosnapcheck.com/v3"
)

func BuildURLForSnapCheck(path string, urlParams Params) (string, error) {
	return BuildAbsolute(path, getBaseSnapCheckURLForENV(), urlParams)
}

func getBaseSnapCheckURLForENV() string {
	baseURL := protocol

	// Get URL by environment
	env := os.Getenv("API_ENV")
	switch env {
	case envDev:
		baseURL = baseURL + baseSnapCheckSandboxURL
	case envStaging:
		baseURL = baseURL + baseSnapCheckSandboxURL
	case envQA:
		baseURL = baseURL + baseSnapCheckSandboxURL
	case envPrd, envProd:
		baseURL = baseURL + baseSnapCheckProdURL
	default:
		panic("Unhandled environment " + env)
	}

	return baseURL
}
