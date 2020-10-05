package url

import (
	"os"
)

const (
	baseAppDevURL   = "dev-app.wise.us"
	baseAppStageURL = "staging-app.wise.us"
	baseAppQAURL    = "qa-app.wise.us"
	baseAppProdURL  = "app.wise.us"
)

//BuildAbsoluteForApp builds an abosolute path from protocol to params
func BuildAbsoluteForApp(path string, urlParams Params) (string, error) {
	return BuildAbsolute(path, getBaseAppURLForENV(), urlParams)
}

func getBaseAppURLForENV() string {
	baseURL := protocol

	env := os.Getenv("API_ENV")

	switch env {
	case envDev:
		baseURL = baseURL + baseAppDevURL
	case envStaging:
		baseURL = baseURL + baseAppStageURL
	case envQA:
		baseURL = baseURL + baseAppQAURL
	default:
		baseURL = baseURL + baseAppProdURL
	}

	return baseURL
}
