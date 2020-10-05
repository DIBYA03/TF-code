package bbva

import (
	"errors"
	"fmt"
	"strings"

	"github.com/wiseco/core-platform/partner/bank"
	"github.com/wiseco/core-platform/partner/bank/data"
)

type consumerCardService struct {
	request  bank.APIRequest
	consumer data.Consumer
	client   *client
}

func (b *consumerBank) CardService(r bank.APIRequest, id bank.ConsumerID) (bank.ConsumerCardService, error) {
	c, err := data.NewConsumerService(r, bank.ProviderNameBBVA).GetByConsumerID(id)
	if err != nil {
		return nil, err
	}

	return &consumerCardService{
		request:  r,
		consumer: *c,
		client:   b.client,
	}, nil
}

func (s *consumerCardService) Create(preq bank.CreateCardRequest) (*bank.GetCardResponse, error) {
	// Get address property
	addrType, ok := partnerBusinessAddressFromMap[preq.Address]
	if !ok {
		return nil, errors.New("invalid address type")
	}

	p, ok := partnerAddressPropertyToConsumer[addrType]
	if !ok {
		return nil, errors.New("invalid address type")
	}

	prop, err := data.NewConsumerPropertyService(s.request, bank.ProviderNameBBVA).GetByConsumerID(s.consumer.ID, p)
	if prop == nil {
		return nil, errors.New("Address does not exist")
	}

	name, ok := partnerCardBusinessNameFrom[preq.BusinessName]
	if !ok {
		return nil, errors.New("invalid business name type")
	}

	num, err := preq.Phone.String()
	if err != nil {
		return nil, err
	}

	request := &CreateCardRequest{
		AccountID:      preq.AccountID.String(),
		CardType:       CardTypeDebit,
		CardholderName: strings.ToUpper(preq.CardholderName),
		BusinessName:   name,
		Delivery:       CardDeliveryStandard,
		Packaging:      CardPackagingRegular,
		AddressID:      string(prop.BankID),
		PhoneNumber:    num,
	}

	req, err := s.client.post("cards/v3.0", s.request, request)
	if err != nil {
		return nil, err
	}

	req.Header.Set("OP-User-Id", string(s.consumer.BankID))
	var resp CreateCardResponse
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}

	cards, err := s.GetAll()
	if err != nil {
		return nil, err
	}

	for _, c := range cards {
		if string(c.CardID) == resp.CardID {
			return &c, nil
		}
	}

	return nil, errors.New("Error creating card")
}

func (s *consumerCardService) Get(id bank.CardBankID) (*bank.GetCardResponse, error) {
	url := fmt.Sprintf("cards/v3.0/%s", id)
	req, err := s.client.get(url, s.request)
	if err != nil {
		return nil, err
	}

	req.Header.Set("OP-User-Id", string(s.consumer.BankID))
	var resp GetCardResponse
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}

	limit, err := s.GetLimit(id)
	if err != nil {
		return nil, err
	}

	return resp.partnerCardResponseTo(limit)
}

func (s *consumerCardService) GetAll() ([]bank.GetCardResponse, error) {
	req, err := s.client.get("cards/v3.0", s.request)
	if err != nil {
		return nil, err
	}

	req.Header.Set("OP-User-Id", string(s.consumer.BankID))
	var resp GetAllCardsResponse
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}

	var pr []bank.GetCardResponse
	for _, cr := range resp.Cards {
		cardResp, err := cr.partnerCardResponseTo(nil)
		if err != nil {
			return pr, err
		}

		pr = append(pr, *cardResp)
	}

	return pr, err
}

func (s *consumerCardService) GetLimit(id bank.CardBankID) (*bank.GetCardLimitResponse, error) {
	url := fmt.Sprintf("cards/v3.0/%s/limits", id)
	req, err := s.client.get(url, s.request)
	if err != nil {
		return nil, err
	}

	req.Header.Set("OP-User-Id", string(s.consumer.BankID))
	var resp GetCardLimitResponse
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}

	return &bank.GetCardLimitResponse{
		DailyATMAmount:        resp.ATMDaily,
		DailyPOSAmount:        resp.POSDaily,
		DailyTransactionCount: resp.DailyTransactions,
	}, nil
}

func (s *consumerCardService) matchConsumerPANSuffix(id bank.CardBankID, panSuffix string) (bool, error) {
	// Check PAN match
	url := fmt.Sprintf("cards/v3.0/%s", id)
	req, err := s.client.get(url, s.request)
	if err != nil {
		return false, err
	}

	req.Header.Set("OP-User-Id", string(s.consumer.BankID))
	var resp GetAllCardsResponse
	if err := s.client.do(req, &resp); err != nil {
		return false, err
	}

	for _, cr := range resp.Cards {
		if string(id) == cr.CardID {
			println("checking values ->>>>>>", cr.CardNumber, panSuffix)
			println("testing suffix match ", cr.CardNumber[len(cr.CardNumber)-len(panSuffix):], panSuffix, len(cr.CardNumber), len(panSuffix))
			if strings.HasSuffix(cr.CardNumber, panSuffix) {
				return true, nil
			}
		}

	}

	return false, nil
}

func (s *consumerCardService) Reissue(preq bank.ReissueCardRequest) (*bank.GetCardResponse, error) {
	// Get address property
	addrType, ok := partnerBusinessAddressFromMap[preq.Address]
	if !ok {
		return nil, errors.New("invalid address type")
	}

	p, ok := partnerAddressPropertyToConsumer[addrType]
	if !ok {
		return nil, errors.New("invalid address type")
	}

	prop, err := data.NewConsumerPropertyService(s.request, bank.ProviderNameBBVA).GetByConsumerID(s.consumer.ID, p)
	if prop == nil {
		return nil, errors.New("Address does not exist")
	}

	reason, ok := partnerCardReissueReasonFrom[preq.Reason]
	if !ok {
		return nil, errors.New("invalid reissue reason")
	}

	name, ok := partnerCardBusinessNameFrom[preq.BusinessName]
	if !ok {
		return nil, errors.New("invalid business name type")
	}

	num, err := preq.Phone.String()
	if err != nil {
		return nil, err
	}

	request := &ReissueCardRequest{
		CardholderName: strings.ToUpper(preq.CardholderName),
		Reason:         reason,
		BusinessName:   name,
		Delivery:       CardDeliveryStandard,
		Packaging:      CardPackagingRegular,
		AddressID:      string(prop.BankID),
		PhoneNumber:    num,
	}

	url := fmt.Sprintf("cards/v3.0/%s/reissue", preq.CardID)
	req, err := s.client.post(url, s.request, request)
	if err != nil {
		return nil, err
	}

	req.Header.Set("OP-User-Id", string(s.consumer.BankID))
	if err := s.client.do(req, nil); err != nil {
		return nil, err
	}

	return s.Get(bank.CardBankID(preq.CardID))
}

func (s *consumerCardService) Activate(preq bank.ActivateCardRequest) (*bank.GetCardResponse, error) {
	// Check PAN match
	panMatch, err := s.matchConsumerPANSuffix(preq.CardID, preq.PANLast6)
	if err != nil {
		return nil, err
	} else if !panMatch {
		return nil, errors.New("invalid PAN last digits")
	}

	// Activate card
	url := fmt.Sprintf("cards/v3.0/%s/activation", preq.CardID)
	req, err := s.client.post(url, s.request, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("OP-User-Id", string(s.consumer.BankID))
	if err := s.client.do(req, nil); err != nil {
		return nil, err
	}

	cards, err := s.GetAll()
	if err != nil {
		return nil, err
	}

	for _, c := range cards {
		if c.CardID == preq.CardID {
			return &c, nil
		}
	}

	return nil, errors.New("Error activating card")
}

func (s *consumerCardService) SetPIN(preq bank.SetCardPINRequest) error {
	// Check PAN match
	panMatch, err := s.matchConsumerPANSuffix(preq.CardID, preq.PANLast6)
	if err != nil {
		return err
	} else if !panMatch {
		return errors.New("invalid PAN last digits")
	}

	// Set card pin
	request := SetCardPINRequest{preq.PIN}

	url := fmt.Sprintf("cards/v3.0/%s/pin", preq.CardID)
	req, err := s.client.patch(url, s.request, &request)
	if err != nil {
		return err
	}

	req.Header.Set("OP-User-Id", string(s.consumer.BankID))
	return s.client.do(req, nil)
}

func (s *consumerCardService) Cancel(preq bank.CancelCardRequest) error {
	// Check PAN match
	panMatch, err := s.matchConsumerPANSuffix(preq.CardID, preq.PANLast6)
	if err != nil {
		return err
	} else if !panMatch {
		return errors.New("invalid PAN last digits")
	}

	return s.CancelInternal(preq.CardID)
}

func (s *consumerCardService) CancelInternal(cardID bank.CardBankID) error {
	// Cancel card
	url := fmt.Sprintf("cards/v3.0/%s", cardID)
	req, err := s.client.delete(url, s.request)
	if err != nil {
		return err
	}

	req.Header.Set("OP-User-Id", string(s.consumer.BankID))
	return s.client.do(req, nil)
}

func (s *consumerCardService) Block(preq bank.CardBlockRequest) ([]bank.CardBlockResponse, error) {
	reason, ok := partnerCardBlockReasonFrom[preq.Reason]
	if !ok {
		return nil, errors.New("invalid block reason")
	}

	request := &CardBlockRequest{
		Reason: reason,
	}

	url := fmt.Sprintf("cards/v3.0/%s/blocks", preq.CardID)
	req, err := s.client.post(url, s.request, request)
	if err != nil {
		return nil, err
	}

	req.Header.Set("OP-User-Id", string(s.consumer.BankID))
	if err := s.client.do(req, nil); err != nil {
		return nil, err
	}

	return s.GetAllBlocks(preq.CardID)
}

func (s *consumerCardService) Unblock(preq bank.CardUnblockRequest) error {
	reason, ok := partnerCardBlockReasonFrom[preq.Reason]
	if !ok {
		return errors.New("invalid block reason")
	}

	url := fmt.Sprintf("cards/v3.0/%s/blocks/%s", preq.CardID, reason)
	req, err := s.client.delete(url, s.request)
	if err != nil {
		return err
	}

	req.Header.Set("OP-User-Id", string(s.consumer.BankID))
	return s.client.do(req, nil)
}

func (s *consumerCardService) GetAllBlocks(id bank.CardBankID) ([]bank.CardBlockResponse, error) {
	url := fmt.Sprintf("cards/v3.0/%s/blocks", id)
	req, err := s.client.get(url, s.request)
	if err != nil {
		return nil, err
	}

	var resp GetAllBlocksResponse

	req.Header.Set("OP-User-Id", string(s.consumer.BankID))
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}

	var blocks []bank.CardBlockResponse
	for _, b := range resp.CardBlocks {
		reason, ok := partnerCardBlockReasonTo[b.Reason]
		if !ok {
			return nil, errors.New("invalid block reason")
		}

		pb := bank.CardBlockResponse{
			CardID:    id,
			Reason:    reason,
			BlockDate: b.BlockDate,
			IsActive:  b.IsActive,
		}

		blocks = append(blocks, pb)
	}

	return blocks, nil
}
