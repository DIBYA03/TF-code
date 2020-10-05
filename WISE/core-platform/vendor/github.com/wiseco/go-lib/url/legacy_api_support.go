package url

import (
	"os"
)

const (
	baseLegApiSupportDevURL   = "gb89qc7bfk.execute-api.us-west-2.amazonaws.com"
	baseLegApiSupportQAURL    = "afolegltvd.execute-api.us-west-2.amazonaws.com"
	baseLegApiSupportStageURL = "9cxhi8blwi.execute-api.us-west-2.amazonaws.com"
	baseLegApiSupportSbxURL   = "1czb0j4qah.execute-api.us-west-2.amazonaws.com"
	baseLegApiSupportProdURL  = "jtnq8aoiob.execute-api.us-west-2.amazonaws.com"
)

//BuildAbsoluteForLegApiSupport build absolute url for service verification
func BuildAbsoluteForLegApiSupport(path string, urlParams Params) (string, error) {
	return BuildAbsolute(path, getBaseLegApiSupportURLForENV(), urlParams)
}

func getBaseLegApiSupportURLForENV() string {
	baseURL := protocol

	switch os.Getenv("API_ENV") {
	case envQA:
		baseURL = baseURL + baseLegApiSupportQAURL
	case envStg, envStaging:
		baseURL = baseURL + baseLegApiSupportStageURL
	case envSbx:
		baseURL = baseURL + baseLegApiSupportSbxURL
	case envPrd, envProd:
		baseURL = baseURL + baseLegApiSupportProdURL
	default:
		baseURL = baseURL + baseLegApiSupportDevURL
	}

	return baseURL
}
