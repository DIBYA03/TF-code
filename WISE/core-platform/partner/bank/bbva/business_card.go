package bbva

import (
	"errors"
	"fmt"
	"strings"

	"github.com/wiseco/core-platform/partner/bank"
	"github.com/wiseco/core-platform/partner/bank/data"
)

type businessCardService struct {
	request  bank.APIRequest
	business data.Business
	consumer data.Consumer
	client   *client
}

func (b *businessBank) CardService(request bank.APIRequest, bid bank.BusinessID, cid bank.ConsumerID) (bank.BusinessCardService, error) {
	bus, err := data.NewBusinessService(request, bank.ProviderNameBBVA).GetByBusinessID(bid)
	if err != nil {
		return nil, err
	}

	c, err := data.NewConsumerService(request, bank.ProviderNameBBVA).GetByConsumerID(cid)
	if err != nil {
		return nil, err
	}

	return &businessCardService{
		request:  request,
		business: *bus,
		consumer: *c,
		client:   b.client,
	}, nil
}

func (s *businessCardService) Create(preq bank.CreateCardRequest) (*bank.GetCardResponse, error) {
	// Get address property
	addrType, ok := partnerConsumerAddressFromMap[preq.Address]
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
		CardholderName: stripCardName(preq.CardholderName),
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
	proxyConfig := s.client.FetchVGSProxyConfig()
	if err := s.client.do(req, &resp, proxyConfig); err != nil {
		return nil, err
	}

	cards, err := s.GetAll()
	if err != nil {
		return nil, err
	}
	c := getCard(resp.CardID, &cards, resp.CardNumber)
	if c == nil {
		return nil, errors.New("Error creating card")
	}
	return c, nil
}

func (s *businessCardService) Get(id bank.CardBankID) (*bank.GetCardResponse, error) {
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

func (s *businessCardService) getAll(nextPageKey *string) (*GetAllCardsResponse, error) {
	url := "cards/v3.0"
	if nextPageKey != nil {
		url = fmt.Sprintf("%s?page_key=%s", url, *nextPageKey)
	}
	req, err := s.client.get(url, s.request)
	if err != nil {
		return nil, err
	}

	req.Header.Set("OP-User-Id", string(s.consumer.BankID))
	var resp GetAllCardsResponse
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (s *businessCardService) GetAll() ([]bank.GetCardResponse, error) {
	var pageKey *string = nil
	var pr []bank.GetCardResponse
	for true {
		resp, err := s.getAll(pageKey)
		if err != nil {
			return nil, err
		}
		for _, cr := range resp.Cards {
			cardResp, err := cr.partnerCardResponseTo(nil)
			if err != nil {
				return pr, err
			}

			pr = append(pr, *cardResp)
		}
		if resp.PageData.HasMore == "true" {
			pageKey = &resp.PageData.NextPageKey
		} else {
			break
		}
	}
	return pr, nil
}

func (s *businessCardService) GetLimit(id bank.CardBankID) (*bank.GetCardLimitResponse, error) {
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

func (s *businessCardService) matchBusinessCardPANSuffix(cid bank.CardBankID, panSuffix string) (bool, error) {
	// Check PAN match
	req, err := s.client.get("cards/v3.0", s.request)
	if err != nil {
		return false, err
	}

	req.Header.Set("OP-User-Id", string(s.consumer.BankID))
	var resp GetAllCardsResponse
	if err := s.client.do(req, &resp); err != nil {
		return false, err
	}

	for _, cr := range resp.Cards {
		if string(cid) == cr.CardID {
			cardNumber := stripSpaces(cr.CardNumber)
			if strings.HasSuffix(cardNumber, panSuffix) {
				return true, nil
			}
		}

	}

	return false, nil
}

func (s *businessCardService) Reissue(preq bank.ReissueCardRequest) (*bank.GetCardResponse, error) {
	// Get address property
	addrType, ok := partnerConsumerAddressFromMap[preq.Address]
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

	num, err := preq.Phone.String()
	if err != nil {
		return nil, err
	}

	request := &ReissueCardRequest{
		Reason:      reason,
		Delivery:    CardDeliveryStandard,
		Packaging:   CardPackagingRegular,
		AddressID:   string(prop.BankID),
		PhoneNumber: num,
	}

	switch reason {
	case CardReissueReasonNameChange:
		name, ok := partnerCardBusinessNameFrom[preq.BusinessName]
		if !ok {
			return nil, errors.New("invalid business name type")
		}

		request.CardholderName = stripCardName(preq.CardholderName)
		request.BusinessName = name
	default:
		break
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

	cards, err := s.GetAll()
	if err != nil {
		return nil, err
	}
	// card number alias set to empty.
	// doesn't matter as it is not saved in db in this flow
	c := getCard(preq.CardID.String(), &cards, "")
	if c == nil {
		return nil, errors.New("Error reissuing card")
	}
	return c, nil
}

func (s *businessCardService) Activate(preq bank.ActivateCardRequest) (*bank.GetCardResponse, error) {
	// Check PAN match
	panMatch, err := s.matchBusinessCardPANSuffix(preq.CardID, preq.PANLast6)
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
	// card number alias set to empty.
	// doesn't matter as it is not saved in db in this flow
	c := getCard(preq.CardID.String(), &cards, "")
	if c == nil {
		return nil, errors.New("Error activating card")
	}
	return c, nil
}

func (s *businessCardService) SetPIN(preq bank.SetCardPINRequest) error {
	// Check PAN match
	panMatch, err := s.matchBusinessCardPANSuffix(preq.CardID, preq.PANLast6)
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

func (s *businessCardService) Cancel(preq bank.CancelCardRequest) error {
	// Check PAN match
	panMatch, err := s.matchBusinessCardPANSuffix(preq.CardID, preq.PANLast6)
	if err != nil {
		return err
	} else if !panMatch {
		return errors.New("invalid PAN last digits")
	}

	return s.CancelInternal(preq.CardID)
}

func (s *businessCardService) CancelInternal(cardID bank.CardBankID) error {
	// Cancel card
	url := fmt.Sprintf("cards/v3.0/%s", cardID)
	req, err := s.client.delete(url, s.request)
	if err != nil {
		return err
	}

	req.Header.Set("OP-User-Id", string(s.consumer.BankID))
	return s.client.do(req, nil)
}

func (s *businessCardService) Block(preq bank.CardBlockRequest) ([]bank.CardBlockResponse, error) {
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

func (s *businessCardService) Unblock(preq bank.CardUnblockRequest) error {
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

func (s *businessCardService) GetAllBlocks(id bank.CardBankID) ([]bank.CardBlockResponse, error) {
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
