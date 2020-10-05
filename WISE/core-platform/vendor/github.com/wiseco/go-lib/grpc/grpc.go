package grpc

import "fmt"

const (
	template = "./ssl/%s.pem"
)

const (
	// Service name passed into a NewClient or NewServer
	ServiceNameVerification = serviceName("verification")
	ServiceNameTransaction  = serviceName("transaction")
	ServiceNameInvoice      = serviceName("invoice")
	ServiceNameEvent        = serviceName("event")
	ServiceNameShopping     = serviceName("shopping")
	ServiceNameBanking      = serviceName("banking")
	ServiceNameAuth         = serviceName("auth")

	// GRPC api services
	ServiceNameApiPartner = serviceName("apiPartner")
	ServiceNameApiClient  = serviceName("apiClient")
	ServiceNameApiCsp     = serviceName("apiCsp")
)

type serviceName string

func getKeyLocation(sn serviceName) (string, error) {
	switch sn {
	case ServiceNameVerification:
		return fmt.Sprintf(template, "service-verification-key"), nil
	case ServiceNameTransaction:
		return fmt.Sprintf(template, "service-transaction-key"), nil
	case ServiceNameInvoice:
		return fmt.Sprintf(template, "service-invoice-key"), nil
	case ServiceNameEvent:
		return fmt.Sprintf(template, "service-event-key"), nil
	case ServiceNameShopping:
		return fmt.Sprintf(template, "service-shopping-key"), nil
	case ServiceNameBanking:
		return fmt.Sprintf(template, "service-banking-key"), nil
	case ServiceNameAuth:
		return fmt.Sprintf(template, "service-auth-key"), nil
	case ServiceNameApiPartner:
		return fmt.Sprintf(template, "service-api-partner-key"), nil
	case ServiceNameApiClient:
		return fmt.Sprintf(template, "service-api-client-key"), nil
	case ServiceNameApiCsp:
		return fmt.Sprintf(template, "service-api-csp-key"), nil
	}

	return "", fmt.Errorf("Service name must be registered. Unknown:%s", sn)
}

func getCertLocation(sn serviceName) (string, error) {
	switch sn {
	case ServiceNameVerification:
		return fmt.Sprintf(template, "service-verification-cert"), nil
	case ServiceNameTransaction:
		return fmt.Sprintf(template, "service-transaction-cert"), nil
	case ServiceNameInvoice:
		return fmt.Sprintf(template, "service-invoice-cert"), nil
	case ServiceNameEvent:
		return fmt.Sprintf(template, "service-event-cert"), nil
	case ServiceNameShopping:
		return fmt.Sprintf(template, "service-shopping-cert"), nil
	case ServiceNameBanking:
		return fmt.Sprintf(template, "service-banking-cert"), nil
	case ServiceNameAuth:
		return fmt.Sprintf(template, "service-auth-cert"), nil
	case ServiceNameApiPartner:
		return fmt.Sprintf(template, "service-api-partner-cert"), nil
	case ServiceNameApiClient:
		return fmt.Sprintf(template, "service-api-client-cert"), nil
	case ServiceNameApiCsp:
		return fmt.Sprintf(template, "service-api-csp-cert"), nil
	}

	return "", fmt.Errorf("Service name must be registered. Unknown:%s", sn)
}
