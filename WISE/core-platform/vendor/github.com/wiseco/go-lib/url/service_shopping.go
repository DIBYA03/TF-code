package url

import (
	"fmt"
	"os"
)

const (
	baseInternalSrvShpDevURL   = "dev-shopping.dev.us-west-2.internal.wise.us"
	baseInternalSrvShpQAURL    = "qa-shopping.dev.us-west-2.internal.wise.us"
	baseInternalSrvShpStageURL = "shopping.staging.us-west-2.internal.wise.us"
	baseInternalSrvShpSbxURL   = "shopping.sbx.us-west-2.internal.wise.us"
	baseInternalSrvShpProdURL  = "shopping.prod.us-west-2.internal.wise.us"

	baseSrvShpDevURL   = "dev-shopping.wise.us"
	baseSrvShpQAURL    = "qa-shopping.wise.us"
	baseSrvShpStageURL = "staging-shopping.wise.us"
	baseSrvShpSbxURL   = "shopping.sbx.wise.us"
	baseSrvShpProdURL  = "shopping.wise.us"
)

//BuildAbsoluteForSrvShp build absolute url for service shopping
func BuildAbsoluteForSrvShp(path string, urlParams Params) (string, error) {
	return BuildAbsolute(path, getBaseSrvShpURLForENV(), urlParams)
}

//GetSVConnectionString returns the
func GetSrvShpConnectionString() string {
	var r string

	port := os.Getenv("GRPC_SERVICE_PORT")
	switch os.Getenv("API_ENV") {
	case envQA:
		r = fmt.Sprintf("%s:%s", baseInternalSrvShpQAURL, port)
	case envStg, envStaging:
		r = fmt.Sprintf("%s:%s", baseInternalSrvShpStageURL, port)
	case envSbx:
		r = fmt.Sprintf("%s:%s", baseInternalSrvShpSbxURL, port)
	case envPrd, envProd:
		r = fmt.Sprintf("%s:%s", baseInternalSrvShpProdURL, port)
	default: // envDev
		r = fmt.Sprintf("%s:%s", baseInternalSrvShpDevURL, port)
	}

	return r
}

func getBaseSrvShpURLForENV() string {
	baseURL := protocol

	switch os.Getenv("API_ENV") {
	case envQA:
		baseURL = baseURL + baseSrvShpQAURL
	case envStg, envStaging:
		baseURL = baseURL + baseSrvShpStageURL
	case envSbx:
		baseURL = baseURL + baseSrvShpSbxURL
	case envPrd, envProd:
		baseURL = baseURL + baseSrvShpProdURL
	default:
		baseURL = baseURL + baseSrvShpDevURL
	}

	return baseURL
}
