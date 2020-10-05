package url

import (
	"os"
)

const (
	baseClearTestURL = "s2s.beta.thomsonreuters.com"
	baseClearProdURL = "s2s.thomsonreuters.com"
)

//BuildAbsoluteForClear - builds an abosolute path from protocol to params
func BuildAbsoluteForClear(path string, urlParams Params) (string, error) {
	return BuildAbsolute(path, getBaseClearURLForENV(), urlParams)
}

// BaseURLForClear - returns the base url for clear apis
func BaseURLForClear() string {
	return getBaseClearURLForENV()
}

func getBaseClearURLForENV() string {
	baseURL := protocol

	// Get URL by environment
	env := os.Getenv("API_ENV")
	switch env {
	case envDev:
		baseURL = baseURL + baseClearTestURL
	case envStaging:
		baseURL = baseURL + baseClearTestURL
	case envQA:
		baseURL = baseURL + baseClearTestURL
	case envProd:
		baseURL = baseURL + baseClearProdURL
	default:
		baseURL = baseURL + baseClearTestURL
	}

	return baseURL
}
