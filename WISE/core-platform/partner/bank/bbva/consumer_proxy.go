package bbva

import (
	bank "github.com/wiseco/core-platform/partner/bank"
	"github.com/wiseco/core-platform/partner/bank/data"
)

func (p *proxyService) GetConsumerBankID(id bank.ConsumerID) (*bank.ConsumerBankID, error) {
	c, err := data.NewConsumerService(p.request, bank.ProviderNameBBVA).GetByConsumerID(id)
	if err != nil {
		return nil, err
	}

	return &c.BankID, err
}

func (p *proxyService) GetConsumerID(id bank.ConsumerBankID) (*bank.ConsumerID, error) {
	c, err := data.NewConsumerService(p.request, bank.ProviderNameBBVA).GetByBankID(id)
	if err != nil {
		return nil, err
	}

	return &c.ConsumerID, err
}
