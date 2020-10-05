package support

import (
	"net/http"
	"strings"

	"github.com/wiseco/core-platform/api"
)

//HandleRequest handle suppor requests
func HandleRequest(r api.APIRequest) (api.APIResponse, error) {

	if strings.ToUpper(r.HTTPMethod) != http.MethodGet {
		return api.NotSupported(r)
	}

	id := r.GetPathParam("userId")
	params := make(map[string]interface{})
	params["firstName"] = r.GetQueryParam("firstName")
	params["phone"] = r.GetQueryParam("phone")
	params["userId"] = r.GetQueryParam("userId")
	params["coId"] = r.GetQueryParam("coId")

	if id != "" {
		return getUserByID(id, r)
	} else if len(params) > 0 {
		return getUserByFilter(params, r)
	}
	return userList(r)
}

//HandlePhoneRequest handle suppor requests
func HandlePhoneRequest(r api.APIRequest) (api.APIResponse, error) {
	phone := r.GetPathParam("phone")

	if strings.ToUpper(r.HTTPMethod) != http.MethodGet {
		return api.NotSupported(r)
	}
	if phone != "" {
		return getUserByPhone(phone, r)
	}
	return userList(r)
}

//HandleAccountRequest handles user bank account request
func HandleAccountRequest(r api.APIRequest) (api.APIResponse, error) {
	userID := r.GetPathParam("userId")
	if strings.ToUpper(r.HTTPMethod) != http.MethodGet {
		return api.NotSupported(r)
	}
	if userID != "" {
		return getAccount(userID, r)
	}
	return api.NotSupported(r)
}

//HandleUnblockAccount ..
func HandleUnblockAccount(r api.APIRequest) (api.APIResponse, error) {
	ID := r.GetPathParam("blockId")
	accountID := r.GetPathParam("accountId")
	if strings.ToUpper(r.HTTPMethod) != http.MethodPost {
		return api.NotSupported(r)
	}

	if ID != "" && accountID != "" {
		return unblockAccount(accountID, ID, r)
	}
	return api.NotSupported(r)
}

//HandleAccountBlocks ..
func HandleAccountBlocks(r api.APIRequest) (api.APIResponse, error) {
	ID := r.GetPathParam("accountId")
	if ID != "" {
		switch strings.ToUpper(r.HTTPMethod) {
		case http.MethodPost:
			return blockAccount(ID, r)
		case http.MethodGet:
			return listOfAccountBlock(ID, r)
		default:
			return api.NotSupported(r)
		}
	}
	return api.NotSupported(r)
}
