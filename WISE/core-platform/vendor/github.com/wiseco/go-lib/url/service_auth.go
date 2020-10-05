package url

import (
	"fmt"
	"os"
)

const (
	baseInternalSrvAuthDevURL   = "dev-auth.dev.us-west-2.internal.wise.us"
	baseInternalSrvAuthQAURL    = "qa-auth.dev.us-west-2.internal.wise.us"
	baseInternalSrvAuthStageURL = "auth.staging.us-west-2.internal.wise.us"
	baseInternalSrvAuthSbxURL   = "auth.sbx.us-west-2.internal.wise.us"
	baseInternalSrvAuthProdURL  = "auth.prod.us-west-2.internal.wise.us"

	baseSrvAuthDevURL   = "dev-auth.wise.us"
	baseSrvAuthQAURL    = "qa-auth.wise.us"
	baseSrvAuthStageURL = "staging-auth.wise.us"
	baseSrvAuthSbxURL   = "auth.sbx.wise.us"
	baseSrvAuthProdURL  = "auth.wise.us"
)

//BuildAbsoluteForSrvAuth build absolute url for service verification
func BuildAbsoluteForSrvAuth(path string, urlParams Params) (string, error) {
	return BuildAbsolute(path, getBaseSrvAuthURLForENV(), urlParams)
}

//GetSrvAuthConnectionString returns the
func GetSrvAuthConnectionString() string {
	var r string

	port := os.Getenv("GRPC_SERVICE_PORT")
	switch os.Getenv("API_ENV") {
	case envQA:
		r = fmt.Sprintf("%s:%s", baseInternalSrvAuthQAURL, port)
	case envStg, envStaging:
		r = fmt.Sprintf("%s:%s", baseInternalSrvAuthStageURL, port)
	case envSbx:
		r = fmt.Sprintf("%s:%s", baseInternalSrvAuthSbxURL, port)
	case envPrd, envProd:
		r = fmt.Sprintf("%s:%s", baseInternalSrvAuthProdURL, port)
	default: // envDev
		r = fmt.Sprintf("%s:%s", baseInternalSrvAuthDevURL, port)
	}

	return r
}

func getBaseSrvAuthURLForENV() string {
	baseURL := protocol

	switch os.Getenv("API_ENV") {
	case envQA:
		baseURL = baseURL + baseSrvAuthQAURL
	case envStg, envStaging:
		baseURL = baseURL + baseSrvAuthStageURL
	case envSbx:
		baseURL = baseURL + baseSrvAuthSbxURL
	case envPrd, envProd:
		baseURL = baseURL + baseSrvAuthProdURL
	default:
		baseURL = baseURL + baseSrvAuthDevURL
	}

	return baseURL
}
