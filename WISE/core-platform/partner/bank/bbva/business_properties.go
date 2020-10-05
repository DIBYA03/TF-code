package bbva

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/wiseco/core-platform/partner/bank"
	"github.com/wiseco/core-platform/partner/bank/data"
)

type BusinessContactResponse struct {
	Contacts []ContactEntityResponse `json:"contacts"`
}

var partnerContactPropertyToBusiness = map[ContactType]bank.BusinessPropertyType{
	ContactTypeEmail: bank.BusinessPropertyTypeContactEmail,
	ContactTypePhone: bank.BusinessPropertyTypeContactPhone,
}

var partnerContactPropertyFromBusiness = map[bank.BusinessPropertyType]ContactType{
	bank.BusinessPropertyTypeContactEmail: ContactTypeEmail,
	bank.BusinessPropertyTypeContactPhone: ContactTypePhone,
}

var partnerBusinessAddressFromMap = map[bank.AddressRequestType]AddressType{
	bank.AddressRequestTypeLegal:       AddressTypeLegal,
	bank.AddressRequestTypeMailing:     AddressTypeMailing,
	bank.AddressRequestTypeHeadquarter: AddressTypeHeadquarter,
}

var partnerBusinessAddressTo = map[AddressType]bank.AddressRequestType{
	AddressTypeLegal:       bank.AddressRequestTypeLegal,
	AddressTypeMailing:     bank.AddressRequestTypeMailing,
	AddressTypeHeadquarter: bank.AddressRequestTypeHeadquarter,
}

type BusinessAddressEntitiesResponse struct {
	Addresses []AddressEntityResponse `json:"address"`
}

var partnerAddressPropertyToBusiness = map[AddressType]bank.BusinessPropertyType{
	AddressTypeLegal:       bank.BusinessPropertyTypeAddressLegal,
	AddressTypeMailing:     bank.BusinessPropertyTypeAddressMailing,
	AddressTypeHeadquarter: bank.BusinessPropertyTypeAddressHeadquarter,
}

func createBusinessContacts(s *client, r bank.APIRequest, b *data.Business, c []ContactEntityResponse) error {
	urlBase := "business/v3.1/contacts"

	// Map contacts
	req, err := s.get(urlBase, r)
	if err != nil {
		return err
	}

	var contacts []ContactEntityResponse
	req.Header.Set("OP-User-Id", string(b.BankID))
	var resp BusinessContactResponse
	if err = s.do(req, &resp); err == nil {
		contacts = resp.Contacts
	}

	// Fetch by id if call fails
	if err != nil {
		for _, cr := range c {
			urlBase = fmt.Sprintf("%s/%s", urlBase, cr.ID)
			req, err := s.get(urlBase, r)
			if err != nil {
				return err
			}

			var resp ContactIDResponse
			req.Header.Set("OP-User-Id", string(b.BankID))
			if err = s.do(req, &resp); err != nil {
				return err
			}

			contacts = append(contacts, resp.Contact)
		}
	}

	for _, cr := range contacts {
		ct, ok := partnerContactPropertyToBusiness[cr.Type]
		if !ok {
			log.Printf("Unknown business contact type: %s", cr.Type)
			continue
		}

		var by []byte
		if cr.Value != "" {
			by, _ = json.Marshal(cr.Value)
		} else {
			by, _ = json.Marshal(cr.Contact)
		}
		_, err := data.NewBusinessPropertyService(r, bank.ProviderNameBBVA).Create(
			data.BusinessPropertyCreate{
				BusinessID: b.ID,
				Type:       ct,
				BankID:     bank.PropertyBankID(cr.ID),
				Value:      by,
			},
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func updateBusinessContact(s *client, r bank.APIRequest, b *data.Business, p bank.BusinessPropertyType, val string) error {
	_, ok := partnerContactPropertyFromBusiness[p]
	if !ok {
		return errors.New("invalid contact type")
	}

	prop, err := data.NewBusinessPropertyService(r, bank.ProviderNameBBVA).GetByBusinessID(b.ID, p)
	if err != nil {
		return errors.New("contact does not exist")
	}

	// Update on BBVA
	urlBase := fmt.Sprintf("business/v3.1/contacts/%s", prop.BankID)
	u := ContactUpdateRequest{
		Contact: ContactUpdateRequestValue{
			Value: val,
		},
	}

	req, err := s.patch(urlBase, r, u)
	if err != nil {
		return err
	}

	req.Header.Set("OP-User-Id", string(b.BankID))
	err = s.do(req, nil)
	if err != nil {
		return err
	}

	// Get address
	req, err = s.get(urlBase, r)
	if err != nil {
		return err
	}

	req.Header.Set("OP-User-Id", string(b.BankID))
	var resp ContactIDResponse
	if err = s.do(req, &resp); err != nil {
		return err
	}

	// Update in DB
	var by []byte
	if resp.Contact.Value != "" {
		by, _ = json.Marshal(resp.Contact.Value)
	} else {
		by, _ = json.Marshal(resp.Contact.Contact)
	}

	// Update business_property table
	updatedProperty := data.BusinessPropertyUpdate{
		ID:    prop.ID,
		Value: by,
	}

	_, err = data.NewBusinessPropertyService(r, bank.ProviderNameBBVA).Update(updatedProperty)
	return err
}

func createBusinessAddresses(s *client, r bank.APIRequest, b *data.Business, ae []AddressEntityResponse) error {
	urlBase := "business/v3.1/addresses"

	// Map addresses
	req, err := s.get(urlBase, r)
	if err != nil {
		return err
	}

	var addresses []AddressEntityResponse
	req.Header.Set("OP-User-Id", string(b.BankID))
	var resp BusinessAddressEntitiesResponse
	if err = s.do(req, &resp); err != nil {
		return err
	}

	addresses = resp.Addresses

	// Fetch by id if call fails
	if err != nil {
		for _, ar := range ae {
			urlBase = fmt.Sprintf("%s/%s", urlBase, ar.ID)
			req, err := s.get(urlBase, r)
			if err != nil {
				return err
			}

			var resp AddressIDResponse
			req.Header.Set("OP-User-Id", string(b.BankID))
			if err = s.do(req, &resp); err != nil {
				return err
			}

			addresses = append(addresses, resp.Address)
		}
	}

	for _, ar := range addresses {
		at, ok := partnerAddressPropertyToBusiness[ar.Type]
		if !ok {
			log.Printf("Unknown business address type: %s", ar.Type)
			continue
		}

		by, _ := json.Marshal(ar)
		_, err := data.NewBusinessPropertyService(r, bank.ProviderNameBBVA).Create(
			data.BusinessPropertyCreate{
				BusinessID: b.ID,
				Type:       at,
				BankID:     bank.PropertyBankID(ar.ID),
				Value:      by,
			},
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func updateBusinessAddress(s *client, r bank.APIRequest, b *data.Business, p bank.BusinessPropertyType, a bank.AddressRequest) error {
	// Update or add address
	prop, err := data.NewBusinessPropertyService(r, bank.ProviderNameBBVA).GetByBusinessID(b.ID, p)
	if err == nil {
		// Update on BBVA
		addrType := AddressTypeEmpty
		urlBase := fmt.Sprintf("business/v3.1/addresses/%s", prop.BankID)
		address := AddressCreateRequest{
			Address: addressFromPartner(a, &addrType),
		}
		req, err := s.patch(urlBase, r, address)
		if err != nil {
			return err
		}

		req.Header.Set("OP-User-Id", string(b.BankID))
		err = s.do(req, nil)
		if err != nil {
			return err
		}

		// Get saved address
		req, err = s.get(urlBase, r)
		if err != nil {
			return err
		}

		var resp AddressIDResponse
		req.Header.Set("OP-User-Id", string(b.BankID))
		if err = s.do(req, &resp); err != nil {
			return err
		}

		// Upate business_property table
		var by []byte
		by, _ = json.Marshal(resp.Address)
		updatedProperty := data.BusinessPropertyUpdate{
			ID:    prop.ID,
			Value: by,
		}

		_, err = data.NewBusinessPropertyService(r, bank.ProviderNameBBVA).Update(updatedProperty)
		return err
	} else if err == sql.ErrNoRows {
		addrType, ok := partnerBusinessAddressFromMap[a.Type]
		if !ok {
			return errors.New("Invalid address type")
		}

		// Add address
		urlBase := "business/v3.1/addresses"
		address := AddressCreateRequest{
			Address: addressFromPartner(a, &addrType),
		}
		req, err := s.post(urlBase, r, address)
		if err != nil {
			return err
		}

		var addrResp AddressCreateResponse
		req.Header.Set("OP-User-Id", string(b.BankID))
		err = s.do(req, &addrResp)
		if err != nil {
			return err
		}

		// Get saved address
		url := fmt.Sprintf("%s/%s", urlBase, addrResp.AddressID)
		req, err = s.get(url, r)
		if err != nil {
			return err
		}

		var resp AddressIDResponse
		req.Header.Set("OP-User-Id", string(b.BankID))
		if err = s.do(req, &resp); err != nil {
			return err
		}

		// Store
		pat, ok := partnerBusinessAddressFromMap[a.Type]
		if !ok {
			return errors.New("invalid address type")
		}

		at, ok := partnerAddressPropertyToBusiness[pat]
		if !ok {
			return errors.New("invalid address type")
		}

		by, _ := json.Marshal(resp.Address)
		_, err = data.NewBusinessPropertyService(r, bank.ProviderNameBBVA).Create(
			data.BusinessPropertyCreate{
				BusinessID: b.ID,
				Type:       at,
				BankID:     bank.PropertyBankID(addrResp.AddressID),
				Value:      by,
			},
		)

		return err
	}

	return err
}
