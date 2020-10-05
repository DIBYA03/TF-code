package bbva

import (
	"errors"
	"fmt"

	"github.com/wiseco/core-platform/partner/bank"
	partnerbank "github.com/wiseco/core-platform/partner/bank"
	"github.com/wiseco/core-platform/partner/bank/data"
)

type businessLinkedCardService struct {
	request  partnerbank.APIRequest
	business data.Business
	client   *client
}

func (b *businessBank) LinkedCardService(request partnerbank.APIRequest, id bank.BusinessID) (partnerbank.BusinessLinkedCardService, error) {
	bus, err := data.NewBusinessService(request, bank.ProviderNameBBVA).GetByBusinessID(id)
	if err != nil {
		return nil, err
	}

	return &businessLinkedCardService{
		request:  request,
		business: *bus,
		client:   b.client,
	}, nil
}

func (s *businessLinkedCardService) Link(preq *partnerbank.LinkedCardRequest) (*partnerbank.LinkedCardResponse, error) {
	permission, ok := partnerLinkedCardPermissionFrom[preq.Permission]
	if !ok {
		return nil, errors.New("invalid linked account permission")
	}

	at := AddressTypeBilling
	regReq := RegisterCardRequest{
		PrimaryAccountNumber: preq.CardNumber,
		ExpirationDate:       preq.Expiration.Format("2006-01"),
		CVVCode:              preq.CVC,
		NameOnAccount:        stripCardName(preq.AccountHolder),
		Nickname:             preq.Alias,
		Usage:                permission,
		BillingAddress:       addressFromPartner(preq.BillingAddress, &at),
	}

	path := "registered-card-accounts/v3.0"
	req, err := s.client.post(path, s.request, regReq)
	if err != nil {
		return nil, err
	}

	req.Header.Set("OP-User-Id", string(s.business.BankID))
	var cardResp RegisterCardResponse
	if err := s.client.do(req, &cardResp); err != nil {
		return nil, err
	}

	return s.Get(partnerbank.LinkedCardBankID(cardResp.AccountReferenceID))
}

func (s *businessLinkedCardService) Get(id partnerbank.LinkedCardBankID) (*partnerbank.LinkedCardResponse, error) {
	path := fmt.Sprintf("registered-card-accounts/v3.0/%s", id)
	req, err := s.client.get(path, s.request)
	if err != nil {
		return nil, err
	}

	req.Header.Set("OP-User-Id", string(s.business.BankID))
	var cardResp GetRegisterCardResponse
	if err := s.client.do(req, &cardResp); err != nil {
		return nil, err
	}

	return cardResp.partnerLinkedCardResponseTo()
}

func (s *businessLinkedCardService) GetAll() ([]partnerbank.LinkedCardResponse, error) {
	path := "registered-card-accounts/v3.0"
	req, err := s.client.get(path, s.request)
	if err != nil {
		return nil, err
	}

	req.Header.Set("OP-User-Id", string(s.business.BankID))
	var cardsResp GetRegisterCardsResponse
	if err := s.client.do(req, &cardsResp); err != nil {
		return nil, err
	}

	pc := []partnerbank.LinkedCardResponse{}
	for _, card := range cardsResp.RegisteredCards {
		c, err := card.partnerLinkedCardResponseTo()
		if err != nil {
			return nil, err
		}

		pc = append(pc, *c)
	}

	return pc, nil
}

func (s *businessLinkedCardService) Unlink(id partnerbank.LinkedCardBankID) error {
	path := fmt.Sprintf("registered-card-accounts/v3.0/%s", id)
	req, err := s.client.delete(path, s.request)
	if err != nil {
		return err
	}

	req.Header.Set("OP-User-Id", string(s.business.BankID))
	return s.client.do(req, nil)
}
