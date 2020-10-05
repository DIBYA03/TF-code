package bbva

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/jung-kurt/gofpdf"
	"github.com/wiseco/core-platform/partner/bank"
	"github.com/wiseco/core-platform/partner/bank/data"
)

type businessAccountService struct {
	request  bank.APIRequest
	business data.Business
	client   *client
}

func (b *businessBank) BankAccountService(r bank.APIRequest, id bank.BusinessID) (bank.BusinessBankAccountService, error) {
	bus, err := data.NewBusinessService(r, bank.ProviderNameBBVA).GetByBusinessID(id)
	if err != nil {
		return nil, err
	}

	return &businessAccountService{
		request:  r,
		business: *bus,
		client:   b.client,
	}, nil
}

func (s *businessAccountService) CreateConsumerAccount(preq bank.CreateConsumerBankAccountRequest) (*bank.CreateConsumerBankAccountResponse, error) {
	return nil, bank.NewErrorFromCode(bank.ErrorCodeNotImplemented)
}

func (s *businessAccountService) Create(preq bank.CreateBusinessBankAccountRequest) (*bank.CreateBusinessBankAccountResponse, error) {
	accountType, ok := partnerAccountTypeFrom[preq.AccountType]
	if !ok {
		return nil, errors.New("only checking accounts are supported by BBVA")
	}

	if preq.IsForeign {
		return nil, errors.New("foreign companies not accepted")
	}

	businessType, ok := partnerBusinessTypeFrom[preq.BusinessType]
	if !ok {
		return nil, errors.New("invalid business type")
	}

	participants, err := s.partnerAllParticipantsFromBusiness(preq)
	if err != nil {
		return nil, err
	} else if len(participants) < 2 {
		return nil, errors.New("Business accounts require 1 or more non-account holder participants")
	}

	accountReq := CreateAccountRequest{
		AccountType:          accountType,
		MultipleParticipants: len(participants) > 1,
		Participants:         participants,
		BusinessType:         businessType,
	}

	// Execute request
	path := "accounts/v3.0"
	req, err := s.client.post(path, s.request, accountReq)
	if err != nil {
		return nil, err
	}

	req.Header.Set("OP-User-Id", string(s.business.BankID))
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
	return &bank.CreateBusinessBankAccountResponse{
		AccountID:        resp.AccountID,
		BankName:         bank.ProviderNameBBVA,
		BusinessID:       preq.BusinessID,
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

func (s *businessAccountService) Get(id bank.AccountBankID) (*bank.GetBankAccountResponse, error) {
	path := fmt.Sprintf("accounts/v3.0/%s", id)
	req, err := s.client.get(path, s.request)
	if err != nil {
		return nil, err
	}

	req.Header.Set("OP-User-Id", string(s.business.BankID))
	var resp = GetBankAccountResponse{}
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}

	return resp.toPartnerGetBankAccountResponse(s.request)
}

func (s *businessAccountService) Patch(id bank.AccountBankID, preq bank.PatchBankAccountRequest) (*bank.GetBankAccountResponse, error) {
	request := &PatchAccountRequest{
		Alias: preq.Alias,
	}
	path := fmt.Sprintf("accounts/v3.0/%s", id)
	req, err := s.client.patch(path, s.request, request)
	if err != nil {
		return nil, err
	}

	req.Header.Set("OP-User-Id", string(s.business.BankID))
	if err := s.client.do(req, nil); err != nil {
		return nil, err
	}

	return s.Get(id)
}

func (s *businessAccountService) Close(id bank.AccountBankID, reason bank.AccountCloseReason) (*bank.GetBankAccountResponse, error) {
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

	req.Header.Set("OP-User-Id", string(s.business.BankID))
	if err := s.client.do(req, nil); err != nil {
		return nil, err
	}

	return s.Get(id)
}

func (s *businessAccountService) Block(preq bank.AccountBlockRequest) (*bank.AccountBlockResponse, error) {
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

	req.Header.Set("OP-User-Id", string(s.business.BankID))
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

func (s *businessAccountService) Unblock(preq bank.AccountUnblockRequest) error {
	path := fmt.Sprintf("accounts/v3.0/%s/blocks/%s", preq.AccountID, preq.BlockID)
	req, err := s.client.delete(path, s.request)
	if err != nil {
		return err
	}

	req.Header.Set("OP-User-Id", string(s.business.BankID))
	return s.client.do(req, nil)
}

func (s *businessAccountService) GetAllBlocks(id bank.AccountBankID) ([]bank.AccountBlockResponse, error) {
	path := fmt.Sprintf("accounts/v3.0/%s/blocks", id)
	req, err := s.client.get(path, s.request)
	if err != nil {
		return nil, err
	}

	var resp GetAccountBlocks

	req.Header.Set("OP-User-Id", string(s.business.BankID))
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

func (s *businessAccountService) GetStatementByID(aid bank.AccountBankID, sid bank.AccountStatementBankID) (*bank.GetAccountStatementDocument, error) {
	// If this is BBVA preprod, return a sample PDF for testing
	if os.Getenv("BBVA_APP_ENV") == "preprod" || os.Getenv("BBVA_APP_ENV") == "sandbox" {
		var b bytes.Buffer

		pdf := gofpdf.New("P", "mm", "A4", "")
		pdf.AddPage()
		pdf.SetFont("Arial", "B", 12)
		pdf.Cell(40, 10, "Wise Statement Sample")

		err := pdf.Output(&b)
		if err != nil {
			return nil, err
		}

		return &bank.GetAccountStatementDocument{
			ContentType: "application/pdf",
			Content:     b.Bytes(),
		}, nil
	}

	path := fmt.Sprintf("accounts/v3.0/%s/statements/%s", aid, sid)
	req, err := s.client.get(path, s.request)
	if err != nil {
		return nil, err
	}

	resp := clientResponse{}
	req.Header.Set("OP-User-Id", string(s.business.BankID))
	req.Header.Set("Accept", "application/pdf")
	if err := s.client.doClientResp(req, &resp); err != nil {
		return nil, err
	}

	contentType := resp.httpResp.Header[http.CanonicalHeaderKey("Content-Type")]
	if len(contentType) == 0 {
		return nil, errors.New("no content type found")
	}

	return &bank.GetAccountStatementDocument{
		ContentType: contentType[0],
		Content:     resp.bytes,
	}, nil
}

func (s *businessAccountService) GetStatements(id bank.AccountBankID) ([]bank.AccountStatementResponse, error) {
	var statements []bank.AccountStatementResponse

	// Create sample data if this isn't prod
	if os.Getenv("BBVA_APP_ENV") == "preprod" || os.Getenv("BBVA_APP_ENV") == "sandbox" {
		statements = append(statements, bank.AccountStatementResponse{
			StatementID: "ST-8e2e3b5f-38bd-3136-b53f-964f4f8ac6dc",
			Description: "Statement for the period ending on 2019-08-30",
			Created:     time.Date(2019, 9, 30, 0, 0, 0, 0, time.UTC),
			PageCount:   5,
		})

		statements = append(statements, bank.AccountStatementResponse{
			StatementID: "ST-acf0cfc7-2b42-3048-8061-3231a5a41d2f",
			Description: "Statement for the period ending on 2019-10-31",
			Created:     time.Date(2019, 10, 31, 0, 0, 0, 0, time.UTC),
			PageCount:   1,
		})

		// sampleTime, _ = time.Parse("2019-11-15'T'12:00:00'Z'", "2019-10-31")
		statements = append(statements, bank.AccountStatementResponse{
			StatementID: "ST-acf0cfc7-2b42-3048-8061-3231a5a41d2f",
			Description: "Statement for the period ending on 2019-11-31",
			Created:     time.Date(2019, 11, 31, 0, 0, 0, 0, time.UTC),
			PageCount:   1,
		})

		return statements, nil
	}

	path := fmt.Sprintf("accounts/v3.0/%s/statements", id)
	req, err := s.client.get(path, s.request)
	if err != nil {
		return nil, err
	}

	resp := clientResponse{body: &GetAccountStatements{}}
	req.Header.Set("OP-User-Id", string(s.business.BankID))
	if err := s.client.doClientResp(req, &resp); err != nil {
		return nil, err
	}

	// Handle 204
	if resp.httpResp.StatusCode == 204 {
		return []bank.AccountStatementResponse{}, nil
	}

	// Return statements
	for _, statement := range resp.body.(*GetAccountStatements).Statements {
		s, err := statement.toPartnerAccountStatement()
		if err != nil {
			return statements, err
		}

		statements = append(statements, *s)
	}

	return statements, nil
}
