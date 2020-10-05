package url

import (
	"fmt"
	"os"
)

const (
	baseInternalSrvVerDevURL   = "dev-verify.dev.us-west-2.internal.wise.us"
	baseInternalSrvVerQAURL    = "qa-verify.dev.us-west-2.internal.wise.us"
	baseInternalSrvVerStageURL = "verify.staging.us-west-2.internal.wise.us"
	baseInternalSrvVerSbxURL   = "verify.sbx.us-west-2.internal.wise.us"
	baseInternalSrvVerProdURL  = "verify.prod.us-west-2.internal.wise.us"

	baseSrvVerDevURL   = "dev-verify.wise.us"
	baseSrvVerQAURL    = "qa-verify.wise.us"
	baseSrvVerStageURL = "staging-verify.wise.us"
	baseSrvVerSbxURL   = "verify.sbx.wise.us"
	baseSrvVerProdURL  = "verify.wise.us"
)

// BuildAbsoluteForSrvVer build absolute url for service verification
func BuildAbsoluteForSrvVer(path string, urlParams Params) (string, error) {
	return BuildAbsolute(path, getBaseSrvVerURLForENV(), urlParams)
}

// GetSrvVerConnectionString returns the
func GetSrvVerConnectionString() string {
	var r string

	port := os.Getenv("GRPC_SERVICE_PORT")
	switch os.Getenv("API_ENV") {
	case envQA:
		r = fmt.Sprintf("%s:%s", baseInternalSrvVerQAURL, port)
	case envStg, envStaging:
		r = fmt.Sprintf("%s:%s", baseInternalSrvVerStageURL, port)
	case envSbx:
		r = fmt.Sprintf("%s:%s", baseInternalSrvVerSbxURL, port)
	case envPrd, envProd:
		r = fmt.Sprintf("%s:%s", baseInternalSrvVerProdURL, port)
	default: // envDev
		r = fmt.Sprintf("%s:%s", baseInternalSrvVerDevURL, port)
	}

	return r
}

func getBaseSrvVerURLForENV() string {
	baseURL := protocol

	switch os.Getenv("API_ENV") {
	case envQA:
		baseURL = baseURL + baseSrvVerQAURL
	case envStg, envStaging:
		baseURL = baseURL + baseSrvVerStageURL
	case envSbx:
		baseURL = baseURL + baseSrvVerSbxURL
	case envPrd, envProd:
		baseURL = baseURL + baseSrvVerProdURL
	default:
		baseURL = baseURL + baseSrvVerDevURL
	}

	return baseURL
}
