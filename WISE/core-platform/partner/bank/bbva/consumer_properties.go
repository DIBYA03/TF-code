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

type ConsumerContactResponse struct {
	Contacts []ContactEntityResponse `json:"contact"`
}

var partnerContactPropertyToConsumer = map[ContactType]bank.ConsumerPropertyType{
	ContactTypeEmail: bank.ConsumerPropertyTypeContactEmail,
	ContactTypePhone: bank.ConsumerPropertyTypeContactPhone,
}

var partnerContactPropertyFromConsumer = map[bank.ConsumerPropertyType]ContactType{
	bank.ConsumerPropertyTypeContactEmail: ContactTypeEmail,
	bank.ConsumerPropertyTypeContactPhone: ContactTypePhone,
}

var partnerConsumerAddressFromMap = map[bank.AddressRequestType]AddressType{
	bank.AddressRequestTypeLegal:   AddressTypeLegal,
	bank.AddressRequestTypeMailing: AddressTypePostal,
	bank.AddressRequestTypeWork:    AddressTypeWork,
}

var partnerConsumerAddressTo = map[AddressType]bank.AddressRequestType{
	AddressTypeLegal:  bank.AddressRequestTypeLegal,
	AddressTypePostal: bank.AddressRequestTypeMailing,
	AddressTypeWork:   bank.AddressRequestTypeWork,
}

type ConsumerAddressEntitiesResponse struct {
	Addresses []AddressEntityResponse `json:"addresses"`
}

var partnerAddressPropertyToConsumer = map[AddressType]bank.ConsumerPropertyType{
	AddressTypeLegal:  bank.ConsumerPropertyTypeAddressLegal,
	AddressTypePostal: bank.ConsumerPropertyTypeAddressMailing,
	AddressTypeWork:   bank.ConsumerPropertyTypeAddressWork,
}

func createConsumerContacts(s *client, r bank.APIRequest, c *data.Consumer, ce []ContactEntityResponse) error {
	// Map contacts
	urlBase := "consumer/v3.0/contact"
	req, err := s.get(urlBase, r)
	if err != nil {
		return err
	}

	var contacts []ContactEntityResponse
	req.Header.Set("OP-User-Id", string(c.BankID))
	var resp ConsumerContactResponse
	if err = s.do(req, &resp); err == nil {
		contacts = resp.Contacts
	}

	// Fetch by id if call fails
	if err != nil {
		for _, cr := range ce {
			urlBase = fmt.Sprintf("%s/%s", urlBase, cr.ID)
			req, err := s.get(urlBase, r)
			if err != nil {
				return err
			}

			var resp ContactIDResponse
			req.Header.Set("OP-User-Id", string(c.BankID))
			if err = s.do(req, &resp); err != nil {
				return err
			}

			contacts = append(contacts, resp.Contact)
		}
	}

	for _, cr := range contacts {
		// Only handle known consumer property types
		ct, ok := partnerContactPropertyToConsumer[cr.Type]
		if !ok {
			log.Printf("Unknown consumer contact type: %s", cr.Type)
			continue
		}

		var b []byte
		if cr.Value != "" {
			b, _ = json.Marshal(cr.Value)
		} else {
			b, _ = json.Marshal(cr.Contact)
		}
		_, err := data.NewConsumerPropertyService(r, bank.ProviderNameBBVA).Create(
			data.ConsumerPropertyCreate{
				ConsumerID: c.ID,
				Type:       ct,
				BankID:     bank.PropertyBankID(cr.ID),
				Value:      b,
			},
		)
		if err != nil {
			log.Printf("Error saving %s contact (%s): %s", cr.Type, cr.ID, err.Error())
		}
	}

	return nil
}

func updateConsumerContact(s *client, r bank.APIRequest, c *data.Consumer, p bank.ConsumerPropertyType, val string) error {
	_, ok := partnerContactPropertyFromConsumer[p]
	if !ok {
		return errors.New("invalid contact type")
	}

	prop, err := data.NewConsumerPropertyService(r, bank.ProviderNameBBVA).GetByConsumerID(c.ID, p)
	if err != nil {
		return errors.New("contact does not exist")
	}

	// Update on BBVA
	urlBase := fmt.Sprintf("consumer/v3.0/contact/%s", prop.BankID)
	u := ContactUpdateRequest{
		Contact: ContactUpdateRequestValue{
			Value: val,
		},
	}

	req, err := s.patch(urlBase, r, u)
	if err != nil {
		return err
	}

	req.Header.Set("OP-User-Id", string(c.BankID))
	err = s.do(req, nil)
	if err != nil {
		return err
	}

	// BBVA doesn't return value so get
	req, err = s.get(urlBase, r)
	if err != nil {
		return err
	}

	req.Header.Set("OP-User-Id", string(c.BankID))
	var resp ContactIDResponse
	if err = s.do(req, &resp); err != nil {
		return err
	}

	// Update in DB
	var b []byte
	if resp.Contact.Value != "" {
		b, _ = json.Marshal(resp.Contact.Value)
	} else {
		b, _ = json.Marshal(resp.Contact.Contact)
	}

	up := data.ConsumerPropertyUpdate{
		ID:    prop.ID,
		Value: b,
	}
	_, err = data.NewConsumerPropertyService(r, bank.ProviderNameBBVA).Update(up)
	return err
}

func createConsumerAddresses(s *client, r bank.APIRequest, c *data.Consumer, ae []AddressEntityResponse) error {
	urlBase := "consumer/v3.0/address"

	// Map addresses
	req, err := s.get(urlBase, r)
	if err != nil {
		return err
	}

	var addresses []AddressEntityResponse
	req.Header.Set("OP-User-Id", string(c.BankID))
	var resp ConsumerAddressEntitiesResponse
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
			req.Header.Set("OP-User-Id", string(c.BankID))
			if err = s.do(req, &resp); err != nil {
				return err
			}

			addresses = append(addresses, resp.Address)
		}
	}

	for _, ar := range addresses {
		at, ok := partnerAddressPropertyToConsumer[ar.Type]
		if !ok {
			log.Printf("Unknown consumer address type: %s", ar.Type)
			continue
		}

		by, _ := json.Marshal(ar)
		_, err := data.NewConsumerPropertyService(r, bank.ProviderNameBBVA).Create(
			data.ConsumerPropertyCreate{
				ConsumerID: c.ID,
				Type:       at,
				BankID:     bank.PropertyBankID(ar.ID),
				Value:      by,
			},
		)
		if err != nil {
			log.Printf("Error saving %s address (%s): %s", at, ar.ID, err.Error())
		}
	}

	return nil
}

func updateConsumerAddress(s *client, r bank.APIRequest, c *data.Consumer, p bank.ConsumerPropertyType, a bank.AddressRequest) error {
	prop, err := data.NewConsumerPropertyService(r, bank.ProviderNameBBVA).GetByConsumerID(c.ID, p)
	if err == nil {
		// Update on BBVA
		addrType := AddressTypeEmpty
		urlBase := fmt.Sprintf("consumer/v3.0/address/%s", prop.BankID)
		address := AddressCreateRequest{
			Address: addressFromPartner(a, &addrType),
		}
		req, err := s.patch(urlBase, r, address)
		if err != nil {
			return err
		}

		req.Header.Set("OP-User-Id", string(c.BankID))
		err = s.do(req, nil)
		if err != nil {
			return err
		}

		// Get Address
		req, err = s.get(urlBase, r)
		if err != nil {
			return err
		}

		var resp AddressIDResponse
		req.Header.Set("OP-User-Id", string(c.BankID))
		if err = s.do(req, &resp); err != nil {
			return err
		}

		// Update in DB
		b, _ := json.Marshal(resp.Address)
		up := data.ConsumerPropertyUpdate{
			ID:    prop.ID,
			Value: b,
		}
		_, err = data.NewConsumerPropertyService(r, bank.ProviderNameBBVA).Update(up)
		return err
	} else if err == sql.ErrNoRows {
		addrType, ok := partnerConsumerAddressFromMap[a.Type]
		if !ok {
			return errors.New("Invalid address type")
		}

		// Add address
		urlBase := "consumer/v3.0/address"
		address := AddressCreateRequest{
			Address: addressFromPartner(a, &addrType),
		}
		req, err := s.post(urlBase, r, address)
		if err != nil {
			return err
		}

		var addrResp AddressCreateResponse
		req.Header.Set("OP-User-Id", string(c.BankID))
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
		req.Header.Set("OP-User-Id", string(c.BankID))
		if err = s.do(req, &resp); err != nil {
			return err
		}

		// Store
		pat, ok := partnerConsumerAddressFromMap[a.Type]
		if !ok {
			return errors.New("invalid address type")
		}

		at, ok := partnerAddressPropertyToConsumer[pat]
		if !ok {
			return errors.New("invalid address type")
		}

		by, _ := json.Marshal(resp.Address)
		_, err = data.NewConsumerPropertyService(r, bank.ProviderNameBBVA).Create(
			data.ConsumerPropertyCreate{
				ConsumerID: c.ID,
				Type:       at,
				BankID:     bank.PropertyBankID(addrResp.AddressID),
				Value:      by,
			},
		)

		return err
	}

	return err
}
