package url

import (
	"fmt"
	"os"
)

const (
	baseInternalSrvApiCspDevURL   = "dev-csp-api.dev.us-west-2.internal.wise.us"
	baseInternalSrvApiCspQAURL    = "qa-csp-api.dev.us-west-2.internal.wise.us"
	baseInternalSrvApiCspStageURL = "csp-api.staging.us-west-2.internal.wise.us"
	baseInternalSrvApiCspSbxURL   = "csp-api.sbx.us-west-2.internal.wise.us"
	baseInternalSrvApiCspProdURL  = "csp-api.prod.us-west-2.internal.wise.us"

	baseSrvApiCspDevURL   = "dev-csp-api.wise.us"
	baseSrvApiCspQAURL    = "qa-csp-api.wise.us"
	baseSrvApiCspStageURL = "staging-csp-api.wise.us"
	baseSrvApiCspSbxURL   = "sbx-csp-api.wise.us"
	baseSrvApiCspProdURL  = "csp-api.wise.us"
)

//BuildAbsoluteForSrvApiCsp build absolute url for service verification
func BuildAbsoluteForSrvApiCsp(path string, urlParams Params) (string, error) {
	return BuildAbsolute(path, getBaseSrvApiCspURLForENV(), urlParams)
}

//GetSrvApiCspConnectionString returns the
func GetSrvApiCspConnectionString() string {
	var r string

	port := os.Getenv("GRPC_SERVICE_PORT")
	switch os.Getenv("API_ENV") {
	case envQA:
		r = fmt.Sprintf("%s:%s", baseInternalSrvApiCspQAURL, port)
	case envStg, envStaging:
		r = fmt.Sprintf("%s:%s", baseInternalSrvApiCspStageURL, port)
	case envSbx:
		r = fmt.Sprintf("%s:%s", baseInternalSrvApiCspSbxURL, port)
	case envPrd, envProd:
		r = fmt.Sprintf("%s:%s", baseInternalSrvApiCspProdURL, port)
	default: // envDev
		r = fmt.Sprintf("%s:%s", baseInternalSrvApiCspDevURL, port)
	}

	return r
}

func getBaseSrvApiCspURLForENV() string {
	baseURL := protocol

	switch os.Getenv("API_ENV") {
	case envQA:
		baseURL = baseURL + baseSrvApiCspQAURL
	case envStg, envStaging:
		baseURL = baseURL + baseSrvApiCspStageURL
	case envSbx:
		baseURL = baseURL + baseSrvApiCspSbxURL
	case envPrd, envProd:
		baseURL = baseURL + baseSrvApiCspProdURL
	default:
		baseURL = baseURL + baseSrvApiCspDevURL
	}

	return baseURL
}
