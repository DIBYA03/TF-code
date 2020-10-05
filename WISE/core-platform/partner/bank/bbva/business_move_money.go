package bbva

import (
	"errors"
	"fmt"

	"github.com/wiseco/core-platform/partner/bank"
	partnerbank "github.com/wiseco/core-platform/partner/bank"
	"github.com/wiseco/core-platform/partner/bank/data"
)

type businessMoneyTransferService struct {
	request  partnerbank.APIRequest
	business data.Business
	client   *client
}

func (b *businessBank) MoneyTransferService(request partnerbank.APIRequest, id bank.BusinessID) (partnerbank.BusinessMoneyTransferService, error) {
	bus, err := data.NewBusinessService(request, bank.ProviderNameBBVA).GetByBusinessID(id)
	if err != nil {
		return nil, err
	}

	return &businessMoneyTransferService{
		request:  request,
		business: *bus,
		client:   b.client,
	}, nil
}

func (s *businessMoneyTransferService) Submit(preq *partnerbank.MoneyTransferRequest) (*partnerbank.MoneyTransferResponse, error) {
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
	req.Header.Set("OP-User-Id", string(string(s.business.BankID)))
	var resp = MoveMoneyResponse{}
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}

	return s.Get(partnerbank.MoneyTransferBankID(resp.MoveMoneyID))
}

func (s *businessMoneyTransferService) Get(id partnerbank.MoneyTransferBankID) (*partnerbank.MoneyTransferResponse, error) {
	path := fmt.Sprintf("movemoney/v3.0/%s", id)
	req, err := s.client.get(path, s.request)
	if err != nil {
		return nil, err
	}

	req.Header.Set("OP-User-Id", string(s.business.BankID))
	var moveMoneyResp GetMoveMoneyResponse
	if err := s.client.do(req, &moveMoneyResp); err != nil {
		return nil, err
	}

	return moveMoneyResp.partnerMoneyTransferResponseTo()
}

func (s *businessMoneyTransferService) GetAll() ([]partnerbank.MoneyTransferResponse, error) {
	/* path := "movemoney/v3.0/"
	   req, err := s.client.get(path, s.request)
	   if err != nil {
	       return nil, err
	   }

	   req.Header.Set("OP-User-Id", string(s.business.BankID))
	   var moveMoneyResp GetMoveMoneyResponse
	   if err := s.client.do(req, &moveMoneyResp); err != nil {
	       return nil, err
	   }

	   return moveMoneyResp.partnerMoneyTransferResponseTo() */

	return nil, errors.New("get all not implemented")
}

func (s *businessMoneyTransferService) Cancel(id partnerbank.MoneyTransferBankID) (*partnerbank.MoneyTransferResponse, error) {
	path := fmt.Sprintf("movemoney/v3.0/%s", id)
	req, err := s.client.delete(path, s.request)
	if err != nil {
		return nil, err
	}

	req.Header.Set("OP-User-Id", string(s.business.BankID))
	err = s.client.do(req, nil)
	if err != nil {
		return nil, err
	}

	return s.Get(id)
}
