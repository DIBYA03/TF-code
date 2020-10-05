package url

import (
	"fmt"
	"os"
)

const (
	baseInternalSrvEvntDevURL   = "dev-event.dev.us-west-2.internal.wise.us"
	baseInternalSrvEvntQAURL    = "qa-event.dev.us-west-2.internal.wise.us"
	baseInternalSrvEvntStageURL = "event.staging.us-west-2.internal.wise.us"
	baseInternalSrvEvntSbxURL   = "event.sbx.us-west-2.internal.wise.us"
	baseInternalSrvEvntProdURL  = "event.prod.us-west-2.internal.wise.us"

	baseSrvEvntDevURL   = "dev-event.wise.us"
	baseSrvEvntQAURL    = "qa-event.wise.us"
	baseSrvEvntStageURL = "staging-event.wise.us"
	baseSrvEvntSbxURL   = "event.sbx.wise.us"
	baseSrvEvntProdURL  = "event.wise.us"
)

//BuildAbsoluteForSrvEvnt build absolute url for service verification
func BuildAbsoluteForSrvEvnt(path string, urlParams Params) (string, error) {
	return BuildAbsolute(path, getBaseSrvEvntURLForENV(), urlParams)
}

//GetSrvEvntConnectionString returns the
func GetSrvEvntConnectionString() string {
	var r string

	port := os.Getenv("GRPC_SERVICE_PORT")
	switch os.Getenv("API_ENV") {
	case envQA:
		r = fmt.Sprintf("%s:%s", baseInternalSrvEvntQAURL, port)
	case envStg, envStaging:
		r = fmt.Sprintf("%s:%s", baseInternalSrvEvntStageURL, port)
	case envSbx:
		r = fmt.Sprintf("%s:%s", baseInternalSrvEvntSbxURL, port)
	case envPrd, envProd:
		r = fmt.Sprintf("%s:%s", baseInternalSrvEvntProdURL, port)
	default: // envDev
		r = fmt.Sprintf("%s:%s", baseInternalSrvEvntDevURL, port)
	}

	return r
}

func getBaseSrvEvntURLForENV() string {
	baseURL := protocol

	switch os.Getenv("API_ENV") {
	case envQA:
		baseURL = baseURL + baseSrvEvntQAURL
	case envStg, envStaging:
		baseURL = baseURL + baseSrvEvntStageURL
	case envSbx:
		baseURL = baseURL + baseSrvEvntSbxURL
	case envPrd, envProd:
		baseURL = baseURL + baseSrvEvntProdURL
	default:
		baseURL = baseURL + baseSrvEvntDevURL
	}

	return baseURL
}
