package url

import "os"

const (
	baseHellosignURL = "api.hellosign.com"
)

//BuildAbsoluteForApp builds an abosolute path from protocol to params
func BuildAbsoluteForHellosign(path string, urlParams Params) (string, error) {
	return BuildAbsolute(path, getBaseHellosignURLForENV(), urlParams)
}

func getBaseHellosignURLForENV() string {
	baseURL := protocol
	baseURL = baseURL + baseHellosignURL

	return baseURL
}

func IsDevEnv() bool {
	// Get URL by environment
	env := os.Getenv("API_ENV")
	switch env {
	case envDev, envStaging, envQA:
		return true
	case envProd:
		return false
	default:
		return true
	}
}
