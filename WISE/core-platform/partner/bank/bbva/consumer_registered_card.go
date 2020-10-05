package bbva

import (
	"errors"
	"fmt"

	"github.com/wiseco/core-platform/partner/bank"
	partnerbank "github.com/wiseco/core-platform/partner/bank"
	"github.com/wiseco/core-platform/partner/bank/data"
)

type consumerLinkedCardService struct {
	request  partnerbank.APIRequest
	consumer data.Consumer
	client   *client
}

func (b *consumerBank) LinkedCardService(request partnerbank.APIRequest, id bank.ConsumerID) (partnerbank.ConsumerLinkedCardService, error) {
	c, err := data.NewConsumerService(request, bank.ProviderNameBBVA).GetByConsumerID(id)
	if err != nil {
		return nil, err
	}

	return &consumerLinkedCardService{
		request:  request,
		consumer: *c,
		client:   b.client,
	}, nil
}

func (s *consumerLinkedCardService) Link(preq *partnerbank.LinkedCardRequest) (*partnerbank.LinkedCardResponse, error) {
	permission, ok := partnerLinkedCardPermissionFrom[preq.Permission]
	if !ok {
		return nil, errors.New("invalid linked account permission")
	}

	at := AddressTypeEmpty
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

	req.Header.Set("OP-User-Id", string(s.consumer.BankID))
	var cardResp RegisterCardResponse
	if err := s.client.do(req, &cardResp); err != nil {
		return nil, err
	}

	return s.Get(partnerbank.LinkedCardBankID(cardResp.AccountReferenceID))
}

func (s *consumerLinkedCardService) Get(id partnerbank.LinkedCardBankID) (*partnerbank.LinkedCardResponse, error) {
	path := fmt.Sprintf("registered-card-accounts/v3.0/%s", id)
	req, err := s.client.get(path, s.request)
	if err != nil {
		return nil, err
	}

	req.Header.Set("OP-User-Id", string(s.consumer.BankID))
	var cardResp GetRegisterCardResponse
	if err := s.client.do(req, &cardResp); err != nil {
		return nil, err
	}

	return cardResp.partnerLinkedCardResponseTo()
}

func (s *consumerLinkedCardService) GetAll() ([]partnerbank.LinkedCardResponse, error) {
	path := "registered-card-accounts/v3.0"
	req, err := s.client.get(path, s.request)
	if err != nil {
		return nil, err
	}

	req.Header.Set("OP-User-Id", string(s.consumer.BankID))
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

func (s *consumerLinkedCardService) Unlink(id partnerbank.LinkedCardBankID) error {
	path := fmt.Sprintf("registered-card-accounts/v3.0/%s", id)
	req, err := s.client.delete(path, s.request)
	if err != nil {
		return err
	}

	req.Header.Set("OP-User-Id", string(s.consumer.BankID))
	return s.client.do(req, nil)
}
