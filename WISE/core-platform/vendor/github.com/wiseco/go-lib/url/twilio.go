package url

const (
	baseTwilioURL = "api.twilio.com"
)

func BuildAbsoluteForTwilio(path string, urlParams Params) (string, error) {
	return BuildAbsolute(path, getBaseTwilioURLForENV(), urlParams)
}

func BaseURLForTwilio() string {
	return getBaseTwilioURLForENV()
}

func getBaseTwilioURLForENV() string {
	return protocol + baseTwilioURL
}
