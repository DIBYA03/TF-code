package bbva

import "github.com/wiseco/core-platform/partner/bank"

type consumerBank struct {
	client *client
}

type businessBank struct {
	client *client
}

type proxyBank struct {
	client *client
}

type proxyService struct {
	request bank.APIRequest
	client  *client
}

func (b *proxyBank) ProxyService(request bank.APIRequest) bank.ProxyService {
	return &proxyService{
		request: request,
		client:  b.client,
	}
}

func (p *proxyService) GetAccessToken() (*string, error) {
	token, err := p.client.oAuth.Token()
	if err != nil {
		return nil, err
	}

	return &token.AccessToken, nil
}

func (p *proxyService) GetBaseAPIURL() string {
	return getBaseAPIURL()
}

func newBBVAConsumer() *consumerBank {
	return &consumerBank{client: newClient()}
}

func newBBVABusiness() *businessBank {
	return &businessBank{client: newClient()}
}

func newBBVAProxy() *proxyBank {
	return &proxyBank{client: newClient()}
}

func init() {
	bank.AddBank(bank.ProviderNameBBVA, newBBVAConsumer(), newBBVABusiness(), newBBVAProxy())
}
