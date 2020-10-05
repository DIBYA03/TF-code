package business

import (
	"encoding/json"

	"github.com/wiseco/core-platform/services"
	"github.com/wiseco/core-platform/shared"

	"github.com/wiseco/core-platform/api"
	csp "github.com/wiseco/core-platform/services/csp/business"
)

var service csp.Service

func handleBusinessID(id string, r api.APIRequest) (api.APIResponse, error) {
	bID, err := shared.ParseBusinessID(id)
	if err != nil {
		return api.BadRequestError(r, err)
	}

	service = csp.New(r.SourceRequest())
	b, err := service.ByID(bID)
	if err != nil {
		_, isENF := err.(*services.ErrorNotFound)
		if isENF {
			return api.NotFound(r)
		}

		return api.InternalServerError(r, err)
	}

	resp, err := json.Marshal(b)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	return api.Success(r, string(resp), false)
}

func handleBusinessList(status kycQueryStatus, limit, offset int, r api.APIRequest) (api.APIResponse, error) {
	service = csp.New(r.SourceRequest())
	list, err := service.ListAll(limit, offset)
	if err != nil {
		return api.InternalServerError(r, err)
	}
	resp, err := json.Marshal(list)
	if err != nil {
		return api.InternalServerError(r, err)
	}
	return api.Success(r, string(resp), false)
}

func handleBusinessUpdate(id string, r api.APIRequest) (api.APIResponse, error) {
	bID, err := shared.ParseBusinessID(id)
	if err != nil {
		return api.BadRequestError(r, err)
	}

	service = csp.New(r.SourceRequest())
	var body csp.BusinessUpdate
	if err := json.Unmarshal([]byte(r.Body), &body); err != nil {
		return api.BadRequest(r, err)
	}
	b, err := service.UpdateID(bID, body)
	if err != nil {
		return api.InternalServerError(r, err)
	}
	resp, err := json.Marshal(b)
	if err != nil {
		return api.InternalServerError(r, err)
	}
	return api.Success(r, string(resp), false)
}

func handleCSPBusinessByBusinessID(id shared.BusinessID, r api.APIRequest) (api.APIResponse, error) {
	item, err := csp.NewCSPService().ByBusinessID(id)
	if err != nil {
		_, isENF := err.(*services.ErrorNotFound)
		if isENF {
			return api.NotFound(r)
		}

		return api.InternalServerError(r, err)
	}

	resp, err := json.Marshal(item)
	if err != nil {
		return api.InternalServerError(r, err)
	}
	return api.Success(r, string(resp), false)
}

func handleCSPBusinessUpdate(id shared.BusinessID, r api.APIRequest) (api.APIResponse, error) {
	var body csp.CSPBusinessUpdate
	if err := json.Unmarshal([]byte(r.Body), &body); err != nil {
		return api.BadRequest(r, err)
	}
	b, err := csp.NewCSPService().CSPBusinessUpdateByBusinessID(id, body)
	if err != nil {
		return api.InternalServerError(r, err)
	}
	resp, err := json.Marshal(b)
	if err != nil {
		return api.InternalServerError(r, err)
	}
	return api.Success(r, string(resp), false)
}

func handleMidDeskVerification(id string, r api.APIRequest) (api.APIResponse, error) {
	bID, err := shared.ParseBusinessID(id)
	if err != nil {
		return api.BadRequestError(r, err)
	}
	b, err := csp.New(r.SourceRequest()).RunMiddesk(bID)
	if err != nil {
		return api.InternalServerError(r, err)
	}
	resp, err := json.Marshal(b)
	return api.Success(r, string(resp), false)
}

func handleClearBusinessVerification(id string, r api.APIRequest) (api.APIResponse, error) {
	bID, err := shared.ParseBusinessID(id)
	if err != nil {
		return api.BadRequestError(r, err)
	}
	b, err := csp.New(r.SourceRequest()).RunClearVerification(bID)
	if err != nil {
		return api.InternalServerError(r, err)
	}
	resp, err := json.Marshal(b)
	return api.Success(r, string(resp), false)
}

func handleState(CSPBusinessID string, r api.APIRequest) (api.APIResponse, error) {
	list, err := csp.NewStateService().List(CSPBusinessID)
	if err != nil {
		return api.InternalServerError(r, err)
	}
	resp, err := json.Marshal(list)
	if err != nil {
		return api.InternalServerError(r, err)
	}
	return api.Success(r, string(resp), false)
}

func handleMidDeskGetVerification(id string, r api.APIRequest) (api.APIResponse, error) {
	bID, err := shared.ParseBusinessID(id)
	if err != nil {
		return api.BadRequestError(r, err)
	}
	b, err := csp.New(r.SourceRequest()).GetMiddesk(bID)
	if err != nil {
		return api.InternalServerError(r, err)
	}
	resp, err := json.Marshal(b)
	if err != nil {
		return api.InternalServerError(r, err)
	}
	return api.Success(r, string(resp), false)

}

func handleClearBusinessGetVerification(id string, r api.APIRequest) (api.APIResponse, error) {
	bID, err := shared.ParseBusinessID(id)
	if err != nil {
		return api.BadRequestError(r, err)
	}
	b, err := csp.New(r.SourceRequest()).GetClearVerification(bID)
	if err != nil {
		return api.InternalServerError(r, err)
	}
	resp, err := json.Marshal(b)
	if err != nil {
		return api.InternalServerError(r, err)
	}
	return api.Success(r, string(resp), false)

}
