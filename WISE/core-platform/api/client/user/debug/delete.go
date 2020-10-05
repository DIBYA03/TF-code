/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

package debug

import (
	"net/http"
	"os"
	"strings"

	"github.com/wiseco/core-platform/api"
	idsrv "github.com/wiseco/core-platform/identity"
	"github.com/wiseco/core-platform/services/address"
	busBank "github.com/wiseco/core-platform/services/banking/business"
	bussrv "github.com/wiseco/core-platform/services/business"
	conserv "github.com/wiseco/core-platform/services/contact"
	"github.com/wiseco/core-platform/services/email"
	usersrv "github.com/wiseco/core-platform/services/user"
	"github.com/wiseco/core-platform/shared"
	trxsrv "github.com/wiseco/core-platform/transaction"
)

func deleteUser(r api.APIRequest, id shared.UserID) (api.APIResponse, error) {
	u, err := usersrv.NewUserService(r.SourceRequest()).GetById(id)
	if err != nil {
		return api.NotFound(r)
	}

	payeeService := busBank.NewLinkedPayeeService(r.SourceRequest())
	addressService := address.NewAddressService(r.SourceRequest())
	contactService := conserv.NewContactService(r.SourceRequest())
	emailService := email.NewEmailService(r.SourceRequest())

	list, err := bussrv.NewBusinessService(r.SourceRequest()).List(0, 20, id)
	if err == nil && len(list) > 0 {
		var bids []shared.BusinessID

		for _, b := range list {
			err = payeeService.DEBUGDeleteAllForBusiness(b.ID)
			if err != nil {
				return api.BadRequestError(r, err)
			}

			err = addressService.DEBUGDeleteAllForBusiness(b.ID)
			if err != nil {
				return api.BadRequestError(r, err)
			}

			contacts, err := contactService.List(0, 100, b.ID)
			if err != nil {
				return api.BadRequestError(r, err)
			}

			for _, c := range contacts {
				err = addressService.DEBUGDeleteAllForContact(shared.ContactID(c.ID))
				if err != nil {
					return api.BadRequestError(r, err)
				}

				if c.EmailID != shared.EmailID("") {
					err = emailService.DEBUGDeleteByID(c.EmailID)
					if err != nil {
						return api.BadRequestError(r, err)
					}
				}
			}

			bMembers, err := bussrv.NewMemberService(r.SourceRequest()).List(0, 20, b.ID)
			if err != nil {
				return api.BadRequestError(r, err)

			}

			for _, bm := range bMembers {
				if bm.EmailID != shared.EmailID("") {
					//TODO because of cyclic imports we can't import email package into the user service below.
					//This should all get moved out into a business logic layer
					//For now lets just mark these as inactive
					err = emailService.Deactivate(bm.EmailID)
					if err != nil {
						return api.BadRequestError(r, err)
					}
				}
			}

			bids = append(bids, b.ID)
		}

		trxsrv.DeleteBusinessTransactions(bids)
	}

	err = addressService.DEBUGDeleteAllForConsumer(u.ConsumerID)
	if err != nil {
		return api.BadRequestError(r, err)
	}

	usersrv.NewUserService(r.SourceRequest()).DeleteById(id)
	idsrv.NewIdentityService(r.IdentitySourceRequest()).Delete(shared.IdentityID(u.IdentityID))
	return api.Success(r, "", false)
}

func HandleUserDeleteRequest(r api.APIRequest) (api.APIResponse, error) {
	// Disallow for production
	switch os.Getenv("APP_ENV") {
	case "prod", "qa-prod", "beta-prod":
		return api.NotAllowedError(r, nil)
	}

	userId, err := shared.ParseUserID(r.GetPathParam("userId"))
	if err != nil {
		return api.BadRequestError(r, err)
	}

	switch strings.ToUpper(r.HTTPMethod) {
	case http.MethodDelete:
		return deleteUser(r, userId)
	default:
		return api.NotSupported(r)
	}
}
