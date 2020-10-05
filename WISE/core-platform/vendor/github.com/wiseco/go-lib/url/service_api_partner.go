package url

import (
	"fmt"
	"os"
)

const (
	baseInternalSrvApiPartnerDevURL   = "dev-partner-api.dev.us-west-2.internal.wise.us"
	baseInternalSrvApiPartnerQAURL    = "qa-partner-api.dev.us-west-2.internal.wise.us"
	baseInternalSrvApiPartnerStageURL = "partner-api.staging.us-west-2.internal.wise.us"
	baseInternalSrvApiPartnerSbxURL   = "partner-api.sbx.us-west-2.internal.wise.us"
	baseInternalSrvApiPartnerProdURL  = "partner-api.prod.us-west-2.internal.wise.us"

	baseSrvApiPartnerDevURL   = "dev-partner-api.wise.us"
	baseSrvApiPartnerQAURL    = "qa-partner-api.wise.us"
	baseSrvApiPartnerStageURL = "staging-partner-api.wise.us"
	baseSrvApiPartnerSbxURL   = "sbx-partner-api.wise.us"
	baseSrvApiPartnerProdURL  = "partner-api.wise.us"
)

//BuildAbsoluteForSrvApiPartner build absolute url for service verification
func BuildAbsoluteForSrvApiPartner(path string, urlParams Params) (string, error) {
	return BuildAbsolute(path, getBaseSrvApiPartnerURLForENV(), urlParams)
}

//GetSrvApiPartnerConnectionString returns the
func GetSrvApiPartnerConnectionString() string {
	var r string

	port := os.Getenv("GRPC_SERVICE_PORT")
	switch os.Getenv("API_ENV") {
	case envQA:
		r = fmt.Sprintf("%s:%s", baseInternalSrvApiPartnerQAURL, port)
	case envStg, envStaging:
		r = fmt.Sprintf("%s:%s", baseInternalSrvApiPartnerStageURL, port)
	case envSbx:
		r = fmt.Sprintf("%s:%s", baseInternalSrvApiPartnerSbxURL, port)
	case envPrd, envProd:
		r = fmt.Sprintf("%s:%s", baseInternalSrvApiPartnerProdURL, port)
	default: // envDev
		r = fmt.Sprintf("%s:%s", baseInternalSrvApiPartnerDevURL, port)
	}

	return r
}

func getBaseSrvApiPartnerURLForENV() string {
	baseURL := protocol

	switch os.Getenv("API_ENV") {
	case envQA:
		baseURL = baseURL + baseSrvApiPartnerQAURL
	case envStg, envStaging:
		baseURL = baseURL + baseSrvApiPartnerStageURL
	case envSbx:
		baseURL = baseURL + baseSrvApiPartnerSbxURL
	case envPrd, envProd:
		baseURL = baseURL + baseSrvApiPartnerProdURL
	default:
		baseURL = baseURL + baseSrvApiPartnerDevURL
	}

	return baseURL
}
