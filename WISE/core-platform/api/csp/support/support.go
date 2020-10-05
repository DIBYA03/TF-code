package support

import (
	"encoding/json"
	"fmt"

	"github.com/wiseco/core-platform/api"
	"github.com/wiseco/core-platform/services"
	"github.com/wiseco/core-platform/services/csp/consumer"
	csp "github.com/wiseco/core-platform/services/csp/support"
)

func getUserByID(id string, r api.APIRequest) (api.APIResponse, error) {
	usr, err := consumer.NewWithSource(r.SourceRequest()).ByUserID(id)
	if err != nil {
		_, isENF := err.(*services.ErrorNotFound)
		if isENF {
			return api.NotFound(r)
		}

		return api.InternalServerError(r, err)
	}

	resp, err := json.Marshal(usr)
	if err != nil {
		return api.InternalServerError(r, err)
	}
	return api.Success(r, string(resp), false)
}

func getUserByFilter(params map[string]interface{}, r api.APIRequest) (api.APIResponse, error) {
	usr, err := consumer.NewWithSource(r.SourceRequest()).ByFilter(params)
	if err != nil && err == err.(*services.ErrorNotFound) {
		return api.NotFound(r)
	}
	resp, err := json.Marshal(usr)
	if err != nil {
		return api.InternalServerError(r, err)
	}
	return api.Success(r, string(resp), false)
}

func getUserByPhone(phone string, r api.APIRequest) (api.APIResponse, error) {
	usr, err := consumer.NewWithSource(r.SourceRequest()).ByPhone(phone)
	if err != nil {
		_, isENF := err.(*services.ErrorNotFound)
		if isENF {
			return api.NotFound(r)
		}

		return api.InternalServerError(r, err)
	}

	resp, err := json.Marshal(usr)
	if err != nil {
		return api.InternalServerError(r, err)
	}
	return api.Success(r, string(resp), false)
}

func userList(r api.APIRequest) (api.APIResponse, error) {
	list, err := consumer.NewWithSource(r.SourceRequest()).UserList()
	if err != nil {
		_, isENF := err.(*services.ErrorNotFound)
		if isENF {
			return api.NotFound(r)
		}

		return api.InternalServerError(r, err)
	}

	resp, err := json.Marshal(list)
	if err != nil {
		return api.InternalServerError(r, err)
	}
	return api.Success(r, string(resp), false)
}

func getAccount(userID string, r api.APIRequest) (api.APIResponse, error) {
	list, err := csp.ListAccount(userID)
	if err != nil {
		return api.InternalServerError(r, err)
	}
	resp, err := json.Marshal(list)
	if err != nil {
		return api.InternalServerError(r, err)
	}
	return api.Success(r, string(resp), false)
}

func blockAccount(id string, r api.APIRequest) (api.APIResponse, error) {

	var body csp.Blocking
	if err := json.Unmarshal([]byte(r.Body), &body); err != nil {
		return api.BadRequest(r, err)
	}
	if err := csp.NewSupport(r.SourceRequest()).Block(id, body); err != nil {
		return api.InternalServerError(r, err)
	}
	if body.OriginatedFrom == "" || body.Reason == "" || !body.Type.Valid() {
		return api.BadRequest(r, fmt.Errorf("Missing or invalid params. body:%v", body))
	}
	status := struct {
		Message string `json:"message"`
	}{Message: "Success"}
	resp, _ := json.Marshal(status)
	return api.Success(r, string(resp), false)

}

func unblockAccount(accountID, id string, r api.APIRequest) (api.APIResponse, error) {
	if err := csp.NewSupport(r.SourceRequest()).Unblock(accountID, id); err != nil {
		return api.InternalServerError(r, err)
	}
	status := struct {
		Message string `json:"message"`
	}{Message: "Success"}
	resp, _ := json.Marshal(status)
	return api.Success(r, string(resp), false)
}

func listOfAccountBlock(id string, r api.APIRequest) (api.APIResponse, error) {
	list, err := csp.NewSupport(r.SourceRequest()).ListOfAccountBlocks(id)
	if err != nil {
		return api.InternalServerError(r, err)
	}
	resp, _ := json.Marshal(list)
	return api.Success(r, string(resp), false)
}
