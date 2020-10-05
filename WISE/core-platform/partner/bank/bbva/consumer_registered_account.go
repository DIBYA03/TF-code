package bbva

import (
	"errors"
	"fmt"

	"github.com/wiseco/core-platform/partner/bank"
	partnerbank "github.com/wiseco/core-platform/partner/bank"
	"github.com/wiseco/core-platform/partner/bank/data"
)

type consumerLinkedAccountService struct {
	request  partnerbank.APIRequest
	consumer data.Consumer
	client   *client
}

func (b *consumerBank) LinkedAccountService(request partnerbank.APIRequest, id bank.ConsumerID) (partnerbank.ConsumerLinkedAccountService, error) {
	c, err := data.NewConsumerService(request, bank.ProviderNameBBVA).GetByConsumerID(id)
	if err != nil {
		return nil, err
	}

	return &consumerLinkedAccountService{
		request:  request,
		consumer: *c,
		client:   b.client,
	}, nil
}

// Register
// https://bbvaopenplatform.com/docs/guides%7Capicontent%7C07-register-bank-account?code=674527&token=5c7df9a7e8288600018c9108
func (s *consumerLinkedAccountService) Link(preq *partnerbank.LinkedBankAccountRequest) (*partnerbank.LinkedBankAccountResponse, error) {
	currency, ok := partnerCurrencyFrom[preq.Currency]
	if !ok {
		return nil, errors.New("invalid currency")
	}

	permission, ok := partnerLinkedAccountPermissionFrom[preq.Permission]
	if !ok {
		return nil, errors.New("invalid linked account permission")
	}

	accountType, ok := partnerLinkedAccountTypeFromConsumer[preq.AccountType]
	if !ok {
		return nil, errors.New("invalid linked account type")
	}

	regReq := RegisterBankAccountRequest{
		AccountNumber: preq.AccountNumber,
		RoutingNumber: preq.RoutingNumber,
		AccountType:   accountType,
		NameOnAccount: stripAccountName(preq.AccountHolderName),
		Currency:      currency,
		Nickname:      preq.Alias,
		Usage:         permission,
	}

	path := "registered-bank-accounts/v3.0"
	req, err := s.client.post(path, s.request, regReq)
	if err != nil {
		return nil, err
	}

	req.Header.Set("OP-User-Id", string(s.consumer.BankID))
	var response RegisterBankAccountResponse
	if err := s.client.do(req, &response); err != nil {
		return nil, err
	}

	return s.Get(partnerbank.LinkedAccountBankID(response.AccountReferenceID))
}

func (s *consumerLinkedAccountService) Get(id partnerbank.LinkedAccountBankID) (*partnerbank.LinkedBankAccountResponse, error) {
	path := fmt.Sprintf("registered-bank-accounts/v3.0/%s", id)
	req, err := s.client.get(path, s.request)
	if err != nil {
		return nil, err
	}

	req.Header.Set("OP-User-Id", string(s.consumer.BankID))
	var accountResp GetRegisterBankAccountResponse
	if err := s.client.do(req, &accountResp); err != nil {
		return nil, err
	}

	return accountResp.partnerLinkedBankAccountResponseTo()
}

func (s *consumerLinkedAccountService) GetAll() ([]partnerbank.LinkedBankAccountResponse, error) {
	path := "registered-bank-accounts/v3.0"
	req, err := s.client.get(path, s.request)
	if err != nil {
		return nil, err
	}

	req.Header.Set("OP-User-Id", string(s.consumer.BankID))
	var accountsResp GetRegisterBankAccountsResponse
	if err := s.client.do(req, &accountsResp); err != nil {
		return nil, err
	}

	pa := []partnerbank.LinkedBankAccountResponse{}
	for _, account := range accountsResp.Accounts {
		a, err := account.partnerLinkedBankAccountResponseTo()
		if err != nil {
			return nil, err
		}

		pa = append(pa, *a)
	}

	return pa, nil
}

func (s *consumerLinkedAccountService) Unlink(id partnerbank.LinkedAccountBankID) error {
	path := fmt.Sprintf("registered-bank-accounts/v3.0/%s", id)
	req, err := s.client.delete(path, s.request)
	if err != nil {
		return err
	}

	req.Header.Set("OP-User-Id", string(s.consumer.BankID))
	return s.client.do(req, nil)
}
