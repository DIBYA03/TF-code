/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

// Package for handling addresss
package contact

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/wiseco/core-platform/api"
	"github.com/wiseco/core-platform/auth"
	addressService "github.com/wiseco/core-platform/services/address"
	"github.com/wiseco/core-platform/shared"
)

type addressPostBody = addressService.AddressCreate
type addressPatchBody = addressService.AddressUpdate

func getAddress(r api.APIRequest, contactID shared.ContactID) (api.APIResponse, error) {
	as, err := addressService.NewAddressService(r.SourceRequest()).GetByContactID(contactID, 10, 0)
	if err != nil {
		return api.BadRequest(r, err)
	}

	addressJSON, err := json.Marshal(as)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	return api.Success(r, string(addressJSON), false)
}

func addressUpdate(r api.APIRequest, addressID shared.AddressID) (api.APIResponse, error) {
	var requestBody addressPatchBody

	if addressID == "" {
		return api.InternalServerError(r, errors.New("Address id cannot be empty"))
	}

	err := json.Unmarshal([]byte(r.Body), &requestBody)

	if err != nil {
		return api.BadRequest(r, err)
	}

	requestBody.ID = addressID
	address, err := addressService.NewAddressService(r.SourceRequest()).Update(&requestBody)

	if err != nil {
		return api.InternalServerError(r, err)
	}

	addressJSON, _ := json.Marshal(address)
	return api.Success(r, string(addressJSON), false)
}

func deactivateAddress(r api.APIRequest, addressID shared.AddressID) (api.APIResponse, error) {
	if addressID == "" {
		return api.InternalServerError(r, errors.New("Address id cannot be empty"))
	}

	err := addressService.NewAddressService(r.SourceRequest()).Deactivate(addressID)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	return api.Success(r, "", false)
}

func createAddress(r api.APIRequest, contactID shared.ContactID) (api.APIResponse, error) {
	var requestBody addressPostBody

	err := json.Unmarshal([]byte(r.Body), &requestBody)

	if err != nil {
		return api.BadRequest(r, err)
	}

	requestBody.ContactID = contactID

	id, err := addressService.NewAddressService(r.SourceRequest()).Create(&requestBody)

	if err != nil {
		return api.InternalServerError(r, err)
	}

	addressJSON, _ := json.Marshal(id)
	return api.Success(r, string(addressJSON), false)
}

// HandleaddressAPIRequests handles the api request
func HandleContactAddressAPIRequests(request api.APIRequest) (api.APIResponse, error) {
	var method = strings.ToUpper(request.HTTPMethod)
	var addressID shared.AddressID

	if request.BusinessID == nil {
		return api.BadRequestError(request, errors.New("missing header X-Wise-Business-ID"))
	}

	err := auth.NewAuthService(request.SourceRequest()).CheckBusinessAccess(*request.BusinessID)
	if err != nil {
		return api.BadRequest(request, err)
	}

	contactID, err := shared.ParseContactID(request.GetPathParam("contactId"))
	if err != nil {
		return api.BadRequest(request, err)
	}

	if request.GetPathParam("addressId") != "" {
		addressID, err = shared.ParseAddressID(request.GetPathParam("addressId"))
		if err != nil {
			return api.BadRequest(request, err)
		}
	}

	switch method {
	case http.MethodGet:
		return getAddress(request, contactID)
	case http.MethodPatch:
		return addressUpdate(request, addressID)
	case http.MethodDelete:
		return deactivateAddress(request, addressID)
	case http.MethodPost:
		return createAddress(request, contactID)
	default:
		return api.NotSupported(request)
	}

}
