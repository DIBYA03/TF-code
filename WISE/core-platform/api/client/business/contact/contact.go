/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

// Package for handling contacts
package contact

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/wiseco/core-platform/api"
	addressService "github.com/wiseco/core-platform/services/address"
	bankingService "github.com/wiseco/core-platform/services/banking/business"
	bankingContacService "github.com/wiseco/core-platform/services/banking/business/contact"
	contactService "github.com/wiseco/core-platform/services/contact"
	"github.com/wiseco/core-platform/shared"
)

type ContactPostBody = contactService.ContactCreate
type ContactPatchBody = contactService.ContactUpdate

type ContactFull struct {
	// Contact
	contactService.Contact

	// Linked account
	LinkedBankAccount []*bankingService.LinkedBankAccount `json:"linkedBankAccount"`

	//Linked card
	LinkedCard []*bankingService.LinkedCard `json:"linkedCard"`

	//Addresses
	Address []addressService.Address `json:"address"`
}

func deactivateContact(r api.APIRequest, id string, businessID shared.BusinessID) (api.APIResponse, error) {

	// Deactivate all linked accounts
	err := bankingContacService.NewLinkedAccountService(r.SourceRequest()).DeactivateAll(id, businessID)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	// Deactivate all linked cards
	err = bankingContacService.NewLinkedCardService(r.SourceRequest()).DeactivateAll(id, businessID)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	err = addressService.NewAddressService(r.SourceRequest()).DeactivateAllForContact(shared.ContactID(id))
	if err != nil {
		return api.InternalServerError(r, err)
	}

	// Deactivate contact
	err = contactService.NewContactService(r.SourceRequest()).Deactivate(id, businessID)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	return api.Success(r, "", false)
}

func getContact(r api.APIRequest, id string, businessID shared.BusinessID) (api.APIResponse, error) {
	contactFull := ContactFull{}

	c, err := contactService.NewContactService(r.SourceRequest()).GetById(id, businessID)
	if err != nil {
		return api.BadRequest(r, err)
	}

	la, err := bankingContacService.NewLinkedAccountService(r.SourceRequest()).List(0, 10, id, businessID)
	if err != nil {
		return api.BadRequest(r, err)
	}

	lc, err := bankingContacService.NewLinkedCardService(r.SourceRequest()).List(0, 10, id, businessID)
	if err != nil {
		return api.BadRequest(r, err)
	}

	as, err := addressService.NewAddressService(r.SourceRequest()).GetByContactID(shared.ContactID(c.ID), 10, 0)
	if err != nil {
		return api.BadRequest(r, err)
	}

	contactFull.Contact = *c
	contactFull.LinkedBankAccount = la
	contactFull.LinkedCard = lc
	contactFull.Address = as

	contactJSON, err := json.Marshal(contactFull)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	return api.Success(r, string(contactJSON), false)
}

func getContacts(r api.APIRequest, businessID shared.BusinessID) (api.APIResponse, error) {
	offset, _ := r.GetQueryIntParamWithDefault("offset", 0)
	limit, _ := r.GetQueryIntParamWithDefault("limit", 10)

	log.Println("headers are ", r.Headers)

	list, err := contactService.NewContactService(r.SourceRequest()).List(offset, limit, businessID)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	contactListJSON, err := json.Marshal(list)
	if err != nil {
		return api.InternalServerError(r, err)
	}
	return api.Success(r, string(contactListJSON), false)
}

func createContact(r api.APIRequest, businessID shared.BusinessID) (api.APIResponse, error) {
	var requestBody ContactPostBody
	userId := r.UserID
	err := json.Unmarshal([]byte(r.Body), &requestBody)

	if err != nil {
		return api.BadRequest(r, err)
	}

	requestBody.UserID = userId
	requestBody.BusinessID = businessID

	id, err := contactService.NewContactService(r.SourceRequest()).Create(&requestBody)

	if err != nil {
		return api.InternalServerError(r, err)
	}

	contactJSON, _ := json.Marshal(id)
	return api.Success(r, string(contactJSON), false)
}

func contactUpdate(r api.APIRequest, id string, businessID shared.BusinessID) (api.APIResponse, error) {
	var requestBody ContactPatchBody
	userId := r.UserID
	err := json.Unmarshal([]byte(r.Body), &requestBody)

	if err != nil {
		return api.BadRequest(r, err)
	}

	requestBody.UserID = userId
	requestBody.ID = id
	requestBody.BusinessID = businessID
	contact, err := contactService.NewContactService(r.SourceRequest()).Update(&requestBody)

	if err != nil {
		return api.InternalServerError(r, err)
	}

	contactJSON, _ := json.Marshal(contact)
	return api.Success(r, string(contactJSON), false)
}

// HandleContactAPIRequests handles the api request
func HandleContactAPIRequests(request api.APIRequest) (api.APIResponse, error) {
	var method = strings.ToUpper(request.HTTPMethod)
	businessID, err := shared.ParseBusinessID(request.GetPathParam("businessId"))
	if err != nil {
		return api.BadRequest(request, err)
	}

	if contactId := request.GetPathParam("contactId"); contactId != "" {
		switch method {
		case http.MethodGet:
			return getContact(request, contactId, businessID)
		case http.MethodPatch:
			return contactUpdate(request, contactId, businessID)
		case http.MethodDelete:
			return deactivateContact(request, contactId, businessID)
		default:
			return api.NotSupported(request)
		}

	}

	switch method {
	case http.MethodGet:
		return getContacts(request, businessID)
	case http.MethodPost:
		return createContact(request, businessID)
	default:
		return api.NotSupported(request)
	}

}
