package url

import (
	"os"
)

const (
	/*
	 * AWS gateway host name
	 */

	baseLegApiClientDevURL   = "9yvfh4v7kk.execute-api.us-west-2.amazonaws.com"
	baseLegApiClientQAURL    = "pqkow88jfg.execute-api.us-west-2.amazonaws.com"
	baseLegApiClientStageURL = "7b79pgvqi9.execute-api.us-west-2.amazonaws.com"
	baseLegApiClientSbxURL   = "vs4pve2fec.execute-api.us-west-2.amazonaws.com"
	baseLegApiClientProdURL  = "k5ktqth4pf.execute-api.us-west-2.amazonaws.com"
)

//BuildAbsoluteForLegApiClient build absolute url for service verification
func BuildAbsoluteForLegApiClient(path string, urlParams Params) (string, error) {
	return BuildAbsolute(path, getBaseLegApiClientURLForENV(), urlParams)
}

func getBaseLegApiClientURLForENV() string {
	baseURL := protocol

	switch os.Getenv("API_ENV") {
	case envQA:
		baseURL = baseURL + baseLegApiClientQAURL
	case envStg, envStaging:
		baseURL = baseURL + baseLegApiClientStageURL
	case envSbx:
		baseURL = baseURL + baseLegApiClientSbxURL
	case envPrd, envProd:
		baseURL = baseURL + baseLegApiClientProdURL
	default:
		baseURL = baseURL + baseLegApiClientDevURL
	}

	return baseURL
}
