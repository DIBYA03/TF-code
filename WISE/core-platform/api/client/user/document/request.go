package document

import (
	"github.com/wiseco/core-platform/api"
)

//DocumentRequest user document create request
func DocumentRequest(r api.APIRequest) (api.APIResponse, error) {
	return api.NotSupported(r)
}

//DocumentURLRequest user document URL request
func DocumentURLRequest(r api.APIRequest) (api.APIResponse, error) {
	return api.NotSupported(r)
}
