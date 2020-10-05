package url

import (
	"fmt"
	"os"
)

const (
	baseInternalSrvApiClientDevURL   = "dev-client-api.dev.us-west-2.internal.wise.us"
	baseInternalSrvApiClientQAURL    = "qa-client-api.dev.us-west-2.internal.wise.us"
	baseInternalSrvApiClientStageURL = "client-api.staging.us-west-2.internal.wise.us"
	baseInternalSrvApiClientSbxURL   = "client-api.sbx.us-west-2.internal.wise.us"
	baseInternalSrvApiClientProdURL  = "client-api.prod.us-west-2.internal.wise.us"

	baseSrvApiClientDevURL   = "dev-client-api.wise.us"
	baseSrvApiClientQAURL    = "qa-client-api.wise.us"
	baseSrvApiClientStageURL = "staging-client-api.wise.us"
	baseSrvApiClientSbxURL   = "sbx-client-api.wise.us"
	baseSrvApiClientProdURL  = "client-api.wise.us"
)

//BuildAbsoluteForSrvApiClient build absolute url for service verification
func BuildAbsoluteForSrvApiClient(path string, urlParams Params) (string, error) {
	return BuildAbsolute(path, getBaseSrvApiClientURLForENV(), urlParams)
}

//GetSrvApiClientConnectionString returns the
func GetSrvApiClientConnectionString() string {
	var r string

	port := os.Getenv("GRPC_SERVICE_PORT")
	switch os.Getenv("API_ENV") {
	case envQA:
		r = fmt.Sprintf("%s:%s", baseInternalSrvApiClientQAURL, port)
	case envStg, envStaging:
		r = fmt.Sprintf("%s:%s", baseInternalSrvApiClientStageURL, port)
	case envSbx:
		r = fmt.Sprintf("%s:%s", baseInternalSrvApiClientSbxURL, port)
	case envPrd, envProd:
		r = fmt.Sprintf("%s:%s", baseInternalSrvApiClientProdURL, port)
	default: // envDev
		r = fmt.Sprintf("%s:%s", baseInternalSrvApiClientDevURL, port)
	}

	return r
}

func getBaseSrvApiClientURLForENV() string {
	baseURL := protocol

	switch os.Getenv("API_ENV") {
	case envQA:
		baseURL = baseURL + baseSrvApiClientQAURL
	case envStg, envStaging:
		baseURL = baseURL + baseSrvApiClientStageURL
	case envSbx:
		baseURL = baseURL + baseSrvApiClientSbxURL
	case envPrd, envProd:
		baseURL = baseURL + baseSrvApiClientProdURL
	default:
		baseURL = baseURL + baseSrvApiClientDevURL
	}

	return baseURL
}
