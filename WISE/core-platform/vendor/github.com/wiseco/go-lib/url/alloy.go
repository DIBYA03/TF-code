package url

import (
	"os"
)

const (
	baseAlloySandboxURL = "sandbox.alloy.co"
	baseAlloyProdURL    = "api.alloy.co"
)

//BuildAbsoluteForApp builds an abosolute path from protocol to params
func BuildAbsoluteForAlloy(path string, urlParams Params) (string, error) {
	return BuildAbsolute(path, getBaseAlloyURLForENV(), urlParams)
}

func BaseURLForAlloy() string {
	return getBaseAlloyURLForENV()
}

func getBaseAlloyURLForENV() string {
	baseURL := protocol

	env := os.Getenv("API_ENV")

	switch env {
	case envDev:
		baseURL = baseURL + baseAlloySandboxURL
	case envStaging:
		baseURL = baseURL + baseAlloySandboxURL
	case envQA:
		baseURL = baseURL + baseAlloySandboxURL
	case envProd:
		baseURL = baseURL + baseAlloyProdURL
	default:
		baseURL = baseURL + baseAlloySandboxURL
	}

	return baseURL
}
