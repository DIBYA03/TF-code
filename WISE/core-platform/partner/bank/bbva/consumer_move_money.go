package bbva

import (
	"errors"
	"fmt"

	"github.com/wiseco/core-platform/partner/bank"
	"github.com/wiseco/core-platform/partner/bank/data"
)

type consumerMoneyTransferService struct {
	request  bank.APIRequest
	consumer data.Consumer
	client   *client
}

func (b *consumerBank) MoneyTransferService(request bank.APIRequest, id bank.ConsumerID) (bank.ConsumerMoneyTransferService, error) {
	c, err := data.NewConsumerService(request, bank.ProviderNameBBVA).GetByConsumerID(id)
	if err != nil {
		return nil, err
	}

	return &consumerMoneyTransferService{
		request:  request,
		consumer: *c,
		client:   b.client,
	}, nil
}

func (s *consumerMoneyTransferService) Submit(preq *bank.MoneyTransferRequest) (*bank.MoneyTransferResponse, error) {
	currency, ok := partnerCurrencyFrom[preq.Currency]
	if !ok {
		return nil, errors.New("invalid currency")
	}

	path := "movemoney/v3.0"

	request := &MoveMoneyRequest{
		OriginAccount:      preq.SourceAccountID.String(),
		DestinationAccount: preq.DestAccountID.String(),
		Amount:             preq.Amount,
		Metadata: MoveMoneyMetadata{
			Currency: currency,
		},
	}
	req, err := s.client.post(path, s.request, request)
	if err != nil {
		return nil, err
	}
	req.Header.Set("OP-User-Id", string(string(s.consumer.BankID)))
	var resp = MoveMoneyResponse{}
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}

	return s.Get(bank.MoneyTransferBankID(resp.MoveMoneyID))
}

func (s *consumerMoneyTransferService) Get(id bank.MoneyTransferBankID) (*bank.MoneyTransferResponse, error) {
	path := fmt.Sprintf("movemoney/v3.0/%s", id)
	req, err := s.client.get(path, s.request)
	if err != nil {
		return nil, err
	}

	req.Header.Set("OP-User-Id", string(s.consumer.BankID))
	var moveMoneyResp GetMoveMoneyResponse
	if err := s.client.do(req, &moveMoneyResp); err != nil {
		return nil, err
	}

	return moveMoneyResp.partnerMoneyTransferResponseTo()
}

func (s *consumerMoneyTransferService) GetAll() ([]bank.MoneyTransferResponse, error) {
	/* path := "movemoney/v3.0/"
	   req, err := s.client.get(path, s.request)
	   if err != nil {
	       return nil, err
	   }

	   req.Header.Set("OP-User-Id", string(s.consumer.BankID))
	   var moveMoneyResp GetMoveMoneyResponse
	   if err := s.client.do(req, &moveMoneyResp); err != nil {
	       return nil, err
	   }

	   return moveMoneyResp.partnerMoneyTransferResponseTo() */

	return nil, errors.New("get all not implemented")
}

func (s *consumerMoneyTransferService) Cancel(id bank.MoneyTransferBankID) (*bank.MoneyTransferResponse, error) {
	path := fmt.Sprintf("movemoney/v3.0/%s", id)
	req, err := s.client.delete(path, s.request)
	if err != nil {
		return nil, err
	}

	req.Header.Set("OP-User-Id", string(s.consumer.BankID))
	err = s.client.do(req, nil)
	if err != nil {
		return nil, err
	}

	return s.Get(id)
}
