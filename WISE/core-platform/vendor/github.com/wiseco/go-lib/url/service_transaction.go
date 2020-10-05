package url

import (
	"fmt"
	"os"
)

const (
	baseInternalSrvTxnDevURL   = "dev-transaction.dev.us-west-2.internal.wise.us"
	baseInternalSrvTxnQAURL    = "qa-transaction.dev.us-west-2.internal.wise.us"
	baseInternalSrvTxnStageURL = "transaction.staging.us-west-2.internal.wise.us"
	baseInternalSrvTxnSbxURL   = "transaction.sbx.us-west-2.internal.wise.us"
	baseInternalSrvTxnProdURL  = "transaction.prod.us-west-2.internal.wise.us"

	baseSrvTxnDevURL   = "dev-transaction.wise.us"
	baseSrvTxnQAURL    = "qa-transaction.wise.us"
	baseSrvTxnStageURL = "staging-transaction.wise.us"
	baseSrvTxnSbxURL   = "transaction.sbx.wise.us"
	baseSrvTxnProdURL  = "transaction.wise.us"
)

//BuildAbsoluteForSrvTxn build absolute url for service verification
func BuildAbsoluteForSrvTxn(path string, urlParams Params) (string, error) {
	return BuildAbsolute(path, getBaseSrvTxnURLForENV(), urlParams)
}

//GetSrvTxnConnectionString returns the
func GetSrvTxnConnectionString() string {
	var r string

	port := os.Getenv("GRPC_SERVICE_PORT")
	switch os.Getenv("API_ENV") {
	case envQA:
		r = fmt.Sprintf("%s:%s", baseInternalSrvTxnQAURL, port)
	case envStg, envStaging:
		r = fmt.Sprintf("%s:%s", baseInternalSrvTxnStageURL, port)
	case envSbx:
		r = fmt.Sprintf("%s:%s", baseInternalSrvTxnSbxURL, port)
	case envPrd, envProd:
		r = fmt.Sprintf("%s:%s", baseInternalSrvTxnProdURL, port)
	default: // envDev
		r = fmt.Sprintf("%s:%s", baseInternalSrvTxnDevURL, port)
	}

	return r
}

func getBaseSrvTxnURLForENV() string {
	baseURL := protocol

	switch os.Getenv("API_ENV") {
	case envQA:
		baseURL = baseURL + baseSrvTxnQAURL
	case envStg, envStaging:
		baseURL = baseURL + baseSrvTxnStageURL
	case envSbx:
		baseURL = baseURL + baseSrvTxnSbxURL
	case envPrd, envProd:
		baseURL = baseURL + baseSrvTxnProdURL
	default:
		baseURL = baseURL + baseSrvTxnDevURL
	}

	return baseURL
}
