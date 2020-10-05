package url

const (
	basicAutocompleteURL = "us-autocomplete.api.smartystreets.com"
	proAutocompleteURL   = "us-autocomplete-pro.api.smartystreets.com"
)

//BuildAbsoluteForApp builds an abosolute path from protocol to params
func BuildBasicAutocompleteURLForSmartyStreets(path string, urlParams Params) (string, error) {
	return BuildAbsolute(path, protocol+basicAutocompleteURL, urlParams)
}

func BuildProAutocompleteURLForSmartyStreets(path string, urlParams Params) (string, error) {
	return BuildAbsolute(path, protocol+proAutocompleteURL, urlParams)
}
