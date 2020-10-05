package url

import (
	"os"
)

const (
	baseCSPDevURL   = "dev-csp.internal.wise.us"
	baseCSPStageURL = "staging-csp.internal.wise.us"
	baseCSPQAURL    = "qa-csp.internal.wise.us"
	baseCSPProdURL  = "csp.internal.wise.us"
)

//BuildAbsoluteForCSP builds an abosolute path from protocol to params
func BuildAbsoluteForCSP(path string, urlParams Params) (string, error) {
	return BuildAbsolute(path, getBaseCSPURLForENV(), urlParams)
}

func getBaseCSPURLForENV() string {
	baseURL := protocol

	env := os.Getenv("API_ENV")

	switch env {
	case envDev:
		baseURL = baseURL + baseCSPDevURL
	case envStaging:
		baseURL = baseURL + baseCSPStageURL
	case envQA:
		baseURL = baseURL + baseCSPQAURL
	default:
		baseURL = baseURL + baseCSPProdURL
	}

	return baseURL
}
