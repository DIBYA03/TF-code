package bbva

import (
	"errors"
	"fmt"

	"github.com/wiseco/core-platform/partner/bank"
	"github.com/wiseco/core-platform/partner/bank/data"
)

func (s *businessAccountService) partnerParticipantRequestFrom(p bank.AccountParticipantRequest) (*ParticipantRequest, error) {
	if s.business.KYCStatus != bank.KYCStatusApproved {
		return nil, errors.New("all participants must be approved")
	}

	role, ok := partnerParticipantRoleFrom[p.Role]
	if !ok {
		return nil, errors.New("invalid participant role requested")
	}

	c, err := data.NewConsumerService(s.request, bank.ProviderNameBBVA).GetByConsumerID(p.ConsumerID)
	if err != nil {
		return nil, err
	}

	return &ParticipantRequest{
		UserID: string(c.BankID),
		Role:   role,
	}, nil
}

func (s *businessAccountService) partnerAllParticipantsFromBusiness(preq bank.CreateBusinessBankAccountRequest) ([]ParticipantRequest, error) {
	p := []ParticipantRequest{
		ParticipantRequest{
			UserID: string(s.business.BankID),
			Role:   ParticipantRoleHolder,
		},
	}

	for _, pp := range preq.ExtraParticipants {
		participant, err := s.partnerParticipantRequestFrom(pp)
		if err != nil {
			return nil, err
		}

		p = append(p, *participant)
	}

	return p, nil
}

func (s *businessAccountService) GetParticipants(id bank.AccountBankID) ([]bank.AccountParticipantResponse, error) {
	path := fmt.Sprintf("accounts/v3.0/%s/participants", id)
	req, err := s.client.get(path, s.request)
	if err != nil {
		return nil, err
	}

	req.Header.Set("OP-User-Id", string(s.business.BankID))
	var response AccountParticipantsResponse
	if err := s.client.do(req, &response); err != nil {
		return nil, err
	}

	return partnerParticipantsTo(s.request, response)
}

func (s *businessAccountService) AddParticipants(id bank.AccountBankID, p []bank.AccountParticipantRequest) ([]bank.AccountParticipantResponse, error) {
	pp := []bank.AccountParticipantResponse{}
	for _, participant := range p {
		npp, err := s.AddParticipant(id, participant)
		if err != nil {
			return pp, err
		}

		pp = append(pp, *npp)
	}

	return pp, nil
}

func (s *businessAccountService) AddParticipant(id bank.AccountBankID, p bank.AccountParticipantRequest) (*bank.AccountParticipantResponse, error) {
	c, err := data.NewConsumerService(s.request, bank.ProviderNameBBVA).GetByConsumerID(p.ConsumerID)
	if err != nil {
		return nil, err
	}

	np := ParticipantRequest{
		UserID: string(c.BankID),
		Role:   partnerParticipantRoleFrom[p.Role],
	}

	path := fmt.Sprintf("accounts/v3.0/%s/participants", id)
	req, err := s.client.post(path, s.request, &np)
	if err != nil {
		return nil, err
	}

	req.Header.Set("OP-User-Id", string(s.business.BankID))
	var resp = AccountParticipantResponse{}
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}

	return resp.partnerParticipantResponseTo(s.request)
}

func (s *businessAccountService) RemoveParticipant(aid bank.AccountBankID, cid bank.ConsumerID) error {
	c, err := data.NewConsumerService(s.request, bank.ProviderNameBBVA).GetByConsumerID(cid)
	if err != nil {
		return err
	}

	path := fmt.Sprintf("accounts/v3.0/%s/participants/%s", aid, c.BankID)
	req, err := s.client.delete(path, s.request)

	req.Header.Set("OP-User-Id", string(s.business.BankID))
	if err := s.client.do(req, nil); err != nil {
		return err
	}

	return nil
}
