package bbva

import (
	bank "github.com/wiseco/core-platform/partner/bank"
	"github.com/wiseco/core-platform/partner/bank/data"
)

func (p *proxyService) GetBusinessBankID(id bank.BusinessID) (*bank.BusinessBankID, error) {
	b, err := data.NewBusinessService(p.request, bank.ProviderNameBBVA).GetByBusinessID(id)
	if err != nil {
		return nil, err
	}

	return &b.BankID, err
}

func (p *proxyService) GetBusinessID(id bank.BusinessBankID) (*bank.BusinessID, error) {
	b, err := data.NewBusinessService(p.request, bank.ProviderNameBBVA).GetByBankID(id)
	if err != nil {
		return nil, err
	}

	return &b.BusinessID, err
}
