package service

type EntityID string

func (id EntityID) String() string {
	return string(id)
}

type ProviderName string

func (n ProviderName) String() string {
	return string(n)
}

const (
	// Send Grid
	ProviderNameSendGrid = ProviderName("sendgrid")
)
