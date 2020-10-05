package plaid

type PlaidResponse struct {
	AccessToken string
	ItemID      string
	Account     []PlaidAccountResponse
	Numbers     PlaidNumberResponse
	RequestID   string
}

type PlaidAccountResponse struct {
	AccountId    string
	Mask         string
	Name         string
	OfficialName string
	SubType      string
	Type         string
	Balance      PlaidAccountBalance
	Owner        []AccountOwner
}

type PlaidNumberResponse struct {
	ACH []PlaidACHResponse
}

type PlaidACHResponse struct {
	Account     string
	AccountId   string
	Routing     string
	WireRouting string
}

type PlaidAccountBalance struct {
	Available              float64
	Current                float64
	Limit                  float64
	ISOCurrencyCode        string
	UnOfficialCurrencyCode string
}

type AccountOwner struct {
	Name    []string
	Address []Address
	Email   []Email
	Phone   []Phone
}

type Address struct {
	Primary       bool
	StreetAddress string
	AddressLine2  string
	City          string
	State         string
	Country       string
	PostalCode    string
}

type Email struct {
	Primary bool
	Email   string
	Type    string
}

type Phone struct {
	Primary bool
	Phone   string
	Type    string
}
