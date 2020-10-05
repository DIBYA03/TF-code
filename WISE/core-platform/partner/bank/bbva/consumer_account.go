package bbva

import (
	"errors"
	"fmt"

	"github.com/wiseco/core-platform/partner/bank"
	"github.com/wiseco/core-platform/partner/bank/data"
)

type consumerAccountService struct {
	request  bank.APIRequest
	consumer data.Consumer
	client   *client
}

func (b *consumerBank) BankAccountService(request bank.APIRequest, id bank.ConsumerID) (bank.ConsumerBankAccountService, error) {
	c, err := data.NewConsumerService(request, bank.ProviderNameBBVA).GetByConsumerID(id)
	if err != nil {
		return nil, err
	}

	return &consumerAccountService{
		request:  request,
		consumer: *c,
		client:   b.client,
	}, nil
}

func (s *consumerAccountService) Create(preq bank.CreateConsumerBankAccountRequest) (*bank.CreateConsumerBankAccountResponse, error) {
	accountType, ok := partnerAccountTypeFrom[preq.AccountType]
	if !ok {
		return nil, errors.New("only checking accounts are supported by BBVA")
	}

	participants, err := s.partnerAllParticipantsFromConsumer(preq)
	if err != nil {
		return nil, err
	} else if len(participants) == 0 {
		return nil, errors.New("Consumer accounts require at least 1 participant")
	}

	accountReq := CreateAccountRequest{
		AccountType:          accountType,
		MultipleParticipants: len(participants) > 1,
		Participants:         participants,
	}

	// Execute request
	path := "accounts/v3.0"
	req, err := s.client.post(path, s.request, accountReq)
	if err != nil {
		return nil, err
	}

	req.Header.Set("OP-User-Id", string(s.consumer.BankID))
	var response = CreateBankAccountResponse{}
	if err := s.client.do(req, &response); err != nil {
		return nil, err
	}

	// On success fetch account by ID and return full data
	resp, err := s.Get(bank.AccountBankID(response.AccountID))
	if err != nil {
		return nil, err
	}

	// Get all participants
	p, err := s.GetParticipants(bank.AccountBankID(resp.AccountID))
	if err != nil {
		return nil, err
	}

	// Return full response object
	return &bank.CreateConsumerBankAccountResponse{
		AccountID:        resp.AccountID,
		BankName:         bank.ProviderNameBBVA,
		AccountHolderID:  preq.AccountHolderID,
		Participants:     p,
		AccountType:      resp.AccountType,
		AccountNumber:    resp.AccountNumber,
		RoutingNumber:    resp.RoutingNumber,
		Alias:            resp.Alias,
		Status:           resp.Status,
		Opened:           resp.Opened,
		AvailableBalance: resp.AvailableBalance,
		PostedBalance:    resp.PostedBalance,
		Currency:         resp.Currency,
	}, nil
}

func (s *consumerAccountService) Get(id bank.AccountBankID) (*bank.GetBankAccountResponse, error) {
	path := fmt.Sprintf("accounts/v3.0/%s", id)
	req, err := s.client.get(path, s.request)
	if err != nil {
		return nil, err
	}

	req.Header.Set("OP-User-Id", string(s.consumer.BankID))
	var resp = GetBankAccountResponse{}
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}

	return resp.toPartnerGetBankAccountResponse(s.request)
}

func (s *consumerAccountService) Patch(id bank.AccountBankID, preq bank.PatchBankAccountRequest) (*bank.GetBankAccountResponse, error) {
	request := &PatchAccountRequest{
		Alias: preq.Alias,
	}
	path := fmt.Sprintf("accounts/v3.0/%s", id)
	req, err := s.client.patch(path, s.request, request)
	if err != nil {
		return nil, err
	}

	req.Header.Set("OP-User-Id", string(s.consumer.BankID))
	if err := s.client.do(req, nil); err != nil {
		return nil, err
	}

	return s.Get(id)
}

func (s *consumerAccountService) Close(id bank.AccountBankID, reason bank.AccountCloseReason) (*bank.GetBankAccountResponse, error) {
	cr, ok := partnerAccountCloseReasonFrom[reason]
	if !ok {
		return nil, errors.New("invalid close reason")
	}

	request := &PatchAccountRequest{
		Status:       AccountStatusClosed,
		StatusReason: cr,
	}
	path := fmt.Sprintf("accounts/v3.0/%s", id)
	req, err := s.client.patch(path, s.request, request)
	if err != nil {
		return nil, err
	}

	req.Header.Set("OP-User-Id", string(s.consumer.BankID))
	if err := s.client.do(req, nil); err != nil {
		return nil, err
	}

	return s.Get(id)
}

func (s *consumerAccountService) Block(preq bank.AccountBlockRequest) (*bank.AccountBlockResponse, error) {
	blockType, ok := partnerAccountBlockTypeFrom[preq.Type]
	if !ok {
		return nil, errors.New("invalid block type")
	}

	request := &CreateAccountBlockRequest{
		BlockType: blockType,
		Reason:    preq.Reason,
	}

	path := fmt.Sprintf("accounts/v3.0/%s/blocks", preq.AccountID)
	req, err := s.client.post(path, s.request, request)
	if err != nil {
		return nil, err
	}

	var resp CreateAccountBlockResponse

	req.Header.Set("OP-User-Id", string(s.consumer.BankID))
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}

	blocksResp, err := s.GetAllBlocks(preq.AccountID)
	if err != nil {
		return nil, err
	}

	for _, block := range blocksResp {
		if block.BlockID == bank.AccountBlockBankID(resp.BlockID) {
			return &block, nil
		}
	}

	return nil, errors.New("block creation failed")
}

func (s *consumerAccountService) Unblock(preq bank.AccountUnblockRequest) error {
	path := fmt.Sprintf("accounts/v3.0/%s/blocks/%s", preq.AccountID, preq.BlockID)
	req, err := s.client.delete(path, s.request)
	if err != nil {
		return err
	}

	req.Header.Set("OP-User-Id", string(s.consumer.BankID))
	return s.client.do(req, nil)
}

func (s *consumerAccountService) GetAllBlocks(id bank.AccountBankID) ([]bank.AccountBlockResponse, error) {
	path := fmt.Sprintf("accounts/v3.0/%s/blocks", id)
	req, err := s.client.get(path, s.request)
	if err != nil {
		return nil, err
	}

	var resp GetAccountBlocks

	req.Header.Set("OP-User-Id", string(s.consumer.BankID))
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}

	var blocks []bank.AccountBlockResponse
	for _, block := range resp.Blocks {
		b, err := block.toPartnerAccountBlockResponse()
		if err != nil {
			return blocks, err
		}

		blocks = append(blocks, *b)
	}

	return blocks, nil
}

func (s *consumerAccountService) GetStatementByID(aid bank.AccountBankID, sid bank.AccountStatementBankID) (*bank.GetAccountStatementDocument, error) {
	path := fmt.Sprintf("/accounts/v3.0/%s/statements/%s", aid, sid)
	req, err := s.client.get(path, s.request)
	if err != nil {
		return nil, err
	}

	var pdf []byte

	req.Header.Set("OP-User-Id", string(s.consumer.BankID))
	if err := s.client.do(req, &pdf); err != nil {
		return nil, err
	}

	return &bank.GetAccountStatementDocument{
		ContentType: "application/pdf",
		Content:     pdf,
	}, nil
}

func (s *consumerAccountService) GetStatements(id bank.AccountBankID) ([]bank.AccountStatementResponse, error) {
	path := fmt.Sprintf("accounts/v3.0/%s/statements", id)
	req, err := s.client.get(path, s.request)
	if err != nil {
		return nil, err
	}

	var resp GetAccountStatements

	req.Header.Set("OP-User-Id", string(s.consumer.BankID))
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}

	var statements []bank.AccountStatementResponse
	for _, statement := range resp.Statements {
		s, err := statement.toPartnerAccountStatement()
		if err != nil {
			return statements, err
		}

		statements = append(statements, *s)
	}

	return statements, nil
}
