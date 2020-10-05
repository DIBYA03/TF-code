package plaid

import (
	"net/http"
	"os"

	"github.com/plaid/plaid-go/plaid"
	"github.com/wiseco/go-lib/log"
)

type plaidService struct {
	log log.Logger
}

func NewAlloyService(l log.Logger) PlaidService {
	return &plaidService{log: l}
}

type PlaidService interface {
	GetAccounts(string) (*PlaidResponse, error)
	GetIdentity(string) (*PlaidResponse, error)
}

func NewPlaidService(l log.Logger) PlaidService {
	return &plaidService{log: l}
}

func getPlaidEnvironment() plaid.Environment {
	var plaidEnv plaid.Environment

	switch os.Getenv("PLAID_ENV") {
	case "production":
		plaidEnv = plaid.Production
		break
	case "development":
		plaidEnv = plaid.Development
		break
	default:
		plaidEnv = plaid.Sandbox
	}

	return plaidEnv
}

func (p *plaidService) GetAccounts(publicToken string) (*PlaidResponse, error) {
	plaidEnv := getPlaidEnvironment()

	clientOptions := plaid.ClientOptions{
		ClientID:    os.Getenv("PLAID_CLIENT_ID"),
		Secret:      os.Getenv("PLAID_SECRET"),
		PublicKey:   os.Getenv("PLAID_PUBLIC_KEY"),
		Environment: plaidEnv,
		HTTPClient:  &http.Client{},
	}

	client, err := plaid.NewClient(clientOptions)
	if err != nil {
		return nil, err
	}

	// Exchange public token to get access token
	response, err := client.ExchangePublicToken(publicToken)
	if err != nil {
		return nil, err
	}

	// Gets list of bank accounts from Plaid
	authResponse, err := client.GetAuth(response.AccessToken)
	if err != nil {
		return nil, err
	}

	resp := transformAccountResponse(authResponse)

	resp.AccessToken = response.AccessToken
	resp.ItemID = response.ItemID
	return resp, nil
}

func (p *plaidService) GetIdentity(publicToken string) (*PlaidResponse, error) {
	plaidEnv := getPlaidEnvironment()

	clientOptions := plaid.ClientOptions{
		ClientID:    os.Getenv("PLAID_CLIENT_ID"),
		Secret:      os.Getenv("PLAID_SECRET"),
		PublicKey:   os.Getenv("PLAID_PUBLIC_KEY"),
		Environment: plaidEnv,
		HTTPClient:  &http.Client{},
	}

	client, err := plaid.NewClient(clientOptions)
	if err != nil {
		return nil, err
	}

	// Exchange public token to get access token
	response, err := client.ExchangePublicToken(publicToken)
	if err != nil {
		return nil, err
	}

	// Gets bank account owner details
	identityResponse, err := client.GetIdentity(response.AccessToken)
	if err != nil {
		return nil, err
	}

	resp := transformIdentityResponse(identityResponse)

	resp.AccessToken = response.AccessToken
	resp.ItemID = response.ItemID
	return resp, nil
}

func transformAccountResponse(authResponse plaid.GetAuthResponse) *PlaidResponse {

	accounts := make([]PlaidAccountResponse, 0)
	for _, a := range authResponse.Accounts {
		plaidAccount := PlaidAccountResponse{
			AccountId:    a.AccountID,
			Mask:         a.Mask,
			Name:         a.Name,
			OfficialName: a.OfficialName,
			Type:         a.Type,
			SubType:      a.Subtype,
			Balance: PlaidAccountBalance{
				Available:              a.Balances.Available,
				Current:                a.Balances.Current,
				ISOCurrencyCode:        a.Balances.ISOCurrencyCode,
				UnOfficialCurrencyCode: a.Balances.UnofficialCurrencyCode,
				Limit:                  a.Balances.Limit,
			},
		}

		accounts = append(accounts, plaidAccount)

	}

	ACHs := make([]PlaidACHResponse, 0)
	for _, a := range authResponse.Numbers.ACH {
		plaidACH := PlaidACHResponse{
			Account:     a.Account,
			AccountId:   a.AccountID,
			Routing:     a.Routing,
			WireRouting: a.WireRouting,
		}

		ACHs = append(ACHs, plaidACH)

	}

	numbers := PlaidNumberResponse{
		ACH: ACHs,
	}

	plaidResponse := PlaidResponse{
		Account:   accounts,
		Numbers:   numbers,
		RequestID: authResponse.RequestID,
	}

	return &plaidResponse
}

func transformIdentityResponse(identityResponse plaid.GetIdentityResponse) *PlaidResponse {
	accounts := make([]PlaidAccountResponse, 0)
	for _, a := range identityResponse.Accounts {
		plaidAccount := PlaidAccountResponse{
			AccountId:    a.AccountID,
			Mask:         a.Mask,
			Name:         a.Name,
			OfficialName: a.OfficialName,
			Type:         a.Type,
			SubType:      a.Subtype,
			Balance: PlaidAccountBalance{
				Available:              a.Balances.Available,
				Current:                a.Balances.Current,
				ISOCurrencyCode:        a.Balances.ISOCurrencyCode,
				UnOfficialCurrencyCode: a.Balances.UnofficialCurrencyCode,
				Limit:                  a.Balances.Limit,
			},
		}

		owners := make([]AccountOwner, 0)
		for _, o := range a.Owners {
			addresses := make([]Address, 0)
			for _, address := range o.Addresses {
				a := Address{
					Primary:       address.Primary,
					StreetAddress: address.Data.Street,
					City:          address.Data.City,
					State:         address.Data.Region,
					Country:       address.Data.Country,
					PostalCode:    address.Data.PostalCode,
				}
				addresses = append(addresses, a)
			}

			emails := make([]Email, 0)
			for _, email := range o.Emails {
				e := Email{
					Primary: email.Primary,
					Email:   email.Data,
					Type:    email.Type,
				}
				emails = append(emails, e)
			}

			phones := make([]Phone, 0)
			for _, phone := range o.PhoneNumbers {
				p := Phone{
					Primary: phone.Primary,
					Phone:   phone.Data,
					Type:    phone.Type,
				}
				phones = append(phones, p)
			}

			owner := AccountOwner{
				Address: addresses,
				Email:   emails,
				Phone:   phones,
				Name:    o.Names,
			}
			owners = append(owners, owner)
		}

		plaidAccount.Owner = owners
		accounts = append(accounts, plaidAccount)
	}

	plaidResponse := PlaidResponse{
		Account:   accounts,
		RequestID: identityResponse.RequestID,
	}

	return &plaidResponse
}
