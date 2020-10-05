package url

//BuildAbsoluteForApp builds an abosolute path from protocol to params
func BuildAbsoluteForShopify(shopURL string, path string, urlParams Params) (string, error) {
	return BuildAbsolute(path, getBaseShopifyURLForENV(shopURL), urlParams)
}

func getBaseShopifyURLForENV(shopURL string) string {
	baseURL := protocol
	baseURL = baseURL + shopURL

	return baseURL
}
