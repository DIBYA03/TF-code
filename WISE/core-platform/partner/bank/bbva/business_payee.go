package bbva

import (
	"fmt"

	"github.com/wiseco/core-platform/partner/bank"
	partnerbank "github.com/wiseco/core-platform/partner/bank"
	"github.com/wiseco/core-platform/partner/bank/data"
)

type businessLinkedPayeeService struct {
	request  partnerbank.APIRequest
	business data.Business
	client   *client
}

func (b *businessBank) LinkedPayeeService(request partnerbank.APIRequest, id bank.BusinessID) (partnerbank.BusinessLinkedPayeeService, error) {
	bus, err := data.NewBusinessService(request, bank.ProviderNameBBVA).GetByBusinessID(id)
	if err != nil {
		return nil, err
	}

	return &businessLinkedPayeeService{
		request:  request,
		business: *bus,
		client:   b.client,
	}, nil
}

func (s *businessLinkedPayeeService) Link(pr *partnerbank.LinkedPayeeRequest) (*partnerbank.LinkedPayeeResponse, error) {
	var resp partnerbank.LinkedPayeeResponse

	path := "registered-payees/v3.0"

	req, err := s.client.post(path, s.request, pr)
	if err != nil {
		return &resp, err
	}

	req.Header.Set("OP-User-Id", string(s.business.BankID))

	err = s.client.do(req, &resp)

	return &resp, err
}

func (s *businessLinkedPayeeService) Get(id partnerbank.BankPayeeID) (*partnerbank.LinkedPayeeResponse, error) {
	var resp partnerbank.LinkedPayeeResponse

	path := fmt.Sprintf("registered-payees/v3.0/%s", id)
	req, err := s.client.get(path, s.request)
	if err != nil {
		return &resp, err
	}

	req.Header.Set("OP-User-Id", string(s.business.BankID))
	err = s.client.do(req, &resp)

	return &resp, err
}

func (s *businessLinkedPayeeService) Unlink(id partnerbank.BankPayeeID) error {
	path := fmt.Sprintf("registered-payees/v3.0/%s", id)

	req, err := s.client.delete(path, s.request)
	if err != nil {
		return err
	}

	req.Header.Set("OP-User-Id", string(s.business.BankID))

	return s.client.do(req, nil)
}
