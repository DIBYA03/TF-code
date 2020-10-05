package url

import (
	"fmt"
	"os"
)

const (
	baseInternalBankingDevURL   = "dev-banking.dev.us-west-2.internal.wise.us"
	baseInternalBankingQAURL    = "qa-banking.dev.us-west-2.internal.wise.us"
	baseInternalBankingStageURL = "banking.staging.us-west-2.internal.wise.us"
	baseInternalBankingSbxURL   = "banking.sbx.us-west-2.internal.wise.us"
	baseInternalBankingProdURL  = "banking.prod.us-west-2.internal.wise.us"

	baseBankingDevURL   = "dev-banking.wise.us"
	baseBankingQAURL    = "qa-verfify.wise.us"
	baseBankingStageURL = "staging-banking.wise.us"
	baseBankingSbxURL   = "banking.sbx.wise.us"
	baseBankingProdURL  = "banking.wise.us"
)

// BuildAbsoluteForBanking build absolute url for service verification
func BuildAbsoluteForBanking(path string, urlParams Params) (string, error) {
	return BuildAbsolute(path, getBaseBankingURLForENV(), urlParams)
}

// GetBankingConnectionString returns the
func GetSrvBankingConnectionString() string {
	var r string

	port := os.Getenv("GRPC_SERVICE_PORT")
	switch os.Getenv("API_ENV") {
	case envQA:
		r = fmt.Sprintf("%s:%s", baseInternalBankingQAURL, port)
	case envStg, envStaging:
		r = fmt.Sprintf("%s:%s", baseInternalBankingStageURL, port)
	case envPrd, envProd:
		r = fmt.Sprintf("%s:%s", baseInternalBankingProdURL, port)
	default: // envDev
		r = fmt.Sprintf("%s:%s", baseInternalBankingDevURL, port)
	}

	return r
}

func getBaseBankingURLForENV() string {
	baseURL := protocol

	switch os.Getenv("API_ENV") {
	case envQA:
		baseURL = baseURL + baseBankingQAURL
	case envStg, envStaging:
		baseURL = baseURL + baseBankingStageURL
	case envPrd, envProd:
		baseURL = baseURL + baseBankingProdURL
	default:
		baseURL = baseURL + baseBankingDevURL
	}

	return baseURL
}
