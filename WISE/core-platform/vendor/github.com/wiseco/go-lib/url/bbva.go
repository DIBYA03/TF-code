package url

import "os"

const (
	baseBBVASandboxURL = "sandbox-apis.bbvaopenplatform.com"
	baseBBVAPreProdURL = "preprod-apis.bbvaopenplatform.com"
	baseBBVAProdURL    = "apis.bbvaopenplatform.com"

	baseBBVAOAUTHSandboxURL = "sbx-paas.bbvacompass.com"
	baseBBVAOAUTHPreProdURL = "pre-paas.bbvacompass.com"
	baseBBVAOAUTHProdURL    = "paas.bbvacompass.com"
)

func BuildURLForBBVA(path string, urlParams Params) (string, error) {
	return BuildAbsolute(path, getBaseBBVAURLForENV(), urlParams)
}

func getBaseBBVAURLForENV() string {
	baseURL := protocol

	// Get URL by environment
	env := os.Getenv("API_ENV")
	switch env {
	case envDev:
		baseURL = baseURL + baseBBVAPreProdURL
	case envStaging:
		baseURL = baseURL + baseBBVAPreProdURL
	case envQA:
		baseURL = baseURL + baseBBVAPreProdURL
	case envProd:
		baseURL = baseURL + baseBBVAProdURL
	default:
		baseURL = baseURL + baseBBVAProdURL
	}

	return baseURL
}

func BuildURLForBBVAOAUTH(path string, urlParams Params) (string, error) {
	return BuildAbsolute(path, getBaseBBVAOAUTHURLForENV(), urlParams)
}

func getBaseBBVAOAUTHURLForENV() string {
	baseURL := protocol

	// Get URL by environment
	env := os.Getenv("API_ENV")
	switch env {
	case envDev:
		baseURL = baseURL + baseBBVAOAUTHPreProdURL
	case envStaging:
		baseURL = baseURL + baseBBVAOAUTHPreProdURL
	case envQA:
		baseURL = baseURL + baseBBVAOAUTHPreProdURL
	case envProd:
		baseURL = baseURL + baseBBVAOAUTHProdURL
	default:
		baseURL = baseURL + baseBBVAOAUTHProdURL
	}

	return baseURL
}
