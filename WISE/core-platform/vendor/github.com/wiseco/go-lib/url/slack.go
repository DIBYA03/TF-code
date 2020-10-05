package url

const (
	baseSlackURL = "hooks.slack.com"
)

//BuildAbsoluteForSlack builds an abosolute path from protocol to params
func BuildAbsoluteForSlack(path string, urlParams Params) (string, error) {
	return BuildAbsolute(path, getBaseSlackURLForENV(), urlParams)
}

func getBaseSlackURLForENV() string {
	baseURL := protocol
	baseURL = baseURL + baseSlackURL

	return baseURL
}
