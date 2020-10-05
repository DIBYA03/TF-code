package url

import (
	"fmt"
	"os"
)

const (
	baseInternalInvoiceDevURL   = "dev-invoice.dev.us-west-2.internal.wise.us"
	baseInternalInvoiceQAURL    = "qa-invoice.dev.us-west-2.internal.wise.us"
	baseInternalInvoiceStageURL = "invoice.staging.us-west-2.internal.wise.us"
	baseInternalInvoiceSbxURL   = "invoice.sbx.us-west-2.internal.wise.us"
	baseInternalInvoiceProdURL  = "invoice.prod.us-west-2.internal.wise.us"
)

// GetInvoiceConnectionString returns the
func GetSrvInvoiceConnectionString() string {
	var r string

	port := os.Getenv("GRPC_SERVICE_PORT")
	switch os.Getenv("API_ENV") {
	case envQA:
		r = fmt.Sprintf("%s:%s", baseInternalInvoiceQAURL, port)
	case envStg, envStaging:
		r = fmt.Sprintf("%s:%s", baseInternalInvoiceStageURL, port)
	case envPrd, envProd:
		r = fmt.Sprintf("%s:%s", baseInternalInvoiceProdURL, port)
	default: // envDev
		r = fmt.Sprintf("%s:%s", baseInternalInvoiceDevURL, port)
	}

	return r
}
