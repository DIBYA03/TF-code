package url

const (
	baseIntercomURL = "api.intercom.io"
)

//BuildAbsoluteForApp builds an abosolute path from protocol to params
func BuildAbsoluteForIntercom(path string, urlParams Params) (string, error) {
	return BuildAbsolute(path, getBaseIntercomURLForENV(), urlParams)
}

func getBaseIntercomURLForENV() string {
	baseURL := protocol
	baseURL = baseURL + baseIntercomURL

	return baseURL
}
