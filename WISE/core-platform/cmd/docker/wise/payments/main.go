package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/wiseco/core-platform/shared"
)

type Request struct {
	UserID              string
	ClientSecret        string
	BusinessName        string
	Amount              string
	Notes               string
	StripeKey           string
	ErrorMessage        string
	PlaidEnv            string
	PublicPlaidKey      string
	PaymentDate         string
	PaymentReceivedDate string
	AllowCardPayment    bool
	AllowBankPayment    bool
	PaymentToken        string
	SegmentKey          string
	ShowBankDetails     bool
	DueDate             *timestamp.Timestamp
	InvoiceTitle        string
	InvoiceNumber       int64
	AccountNumber       string
	RoutingNumber       string
	BusinessLogo        string
	Currency            string
	Country             string
	AmountStripe        *int64
}

type TransferRequest struct {
	UserID       string
	BusinessName string
	Amount       string
	Notes        string
	ErrorMessage string
	PaymentToken string
	SegmentKey   string
}

type Receipt struct {
	PaymentDate  string
	BusinessName string
	PaymentMode  string
	CardBrand    string
	CardLast4    string
	Amount       string
	Location     string
	ErrorMessage string
	StripeKey    string
	ClientSecret string
	Address      string
	Latitude     string
	Longitude    string
	BusinessID   shared.BusinessID
	ReceiptID    string
}

type InvoiceSuccessResponse struct {
	Amount                string
	CardType              string
	CardNumber            string
	InvoiceTitle          string
	CustomerName          string
	BusinessName          string
	WalletType            string
	ReceiptNumber         string
	IsCardInfoAvailable   bool
	IsWalletInfoAvailable bool
}

type HTMLTemplate string

const (
	HTMLTemplateIndex              = HTMLTemplate("index.html")
	HTMLTemplateCardPayment        = HTMLTemplate("card-payment.html")
	HTMLTemplateBankAccountPayment = HTMLTemplate("bank-account-payment.html")
	HTMLTemplateCardSuccess        = HTMLTemplate("card-success.html")
	HTMLTemplateBankAccountSuccess = HTMLTemplate("bank-success.html")
	HTMLTemplateInvoicePaid        = HTMLTemplate("invoice-paid.html")
	HTMLTemplateError              = HTMLTemplate("error.html")
	HTMLTemplatePlaid              = HTMLTemplate("plaid-details")

	HTMLTemplateGetStarted         = HTMLTemplate("pay-instant/get-started.html")
	HTMLTemplateCardPay            = HTMLTemplate("pay-instant/card-pay.html")
	HTMLTemplateBankPay            = HTMLTemplate("pay-instant/bank-pay.html")
	HTMLTemplateCardPaymentSuccess = HTMLTemplate("pay-instant/card-payment-success.html")
	HTMLTemplateBankPaymentSuccess = HTMLTemplate("pay-instant/bank-payment-success.html")
	HTMLTemplateLinkedAccount      = HTMLTemplate("linked-account")
	HTMLTemplateLinkedCard         = HTMLTemplate("linked-card")
	HTMLTemplateCardNotSupported   = HTMLTemplate("pay-instant/card-not-supported.html")

	HTMLTemplateInvoiceRequestIndex = HTMLTemplate("invoice-new.html")
	HTMLTemplateInvoiceView         = HTMLTemplate("pay-invoice.html")
	HTMLTemplateInvoiceSuccess      = HTMLTemplate("invoice-success.html")
)

func main() {
	f := http.FileServer(http.Dir("fonts"))
	http.Handle("/fonts/", http.StripPrefix("/fonts/", f))

	i := http.FileServer(http.Dir("images"))
	http.Handle("/images/", http.StripPrefix("/images/", i))

	css := http.FileServer(http.Dir("css"))
	http.Handle("/css/", http.StripPrefix("/css/", css))

	components := http.FileServer(http.Dir("components"))
	http.Handle("/components/", http.StripPrefix("/components/", components))

	// Needed for 404 and maintenance page
	http.HandleFunc("/favicon.ico", favIconIcoHandler)
	http.HandleFunc("/favicon.png", favIconPngHandler)

	// if in maintenance, only show maintenance page
	if os.Getenv("MAINTENANCE_ENABLED") == "true" {
		maintenanceMode()
	} else {
		liveMode()
	}
}

func maintenanceMode() {
	http.HandleFunc("/", maintenanceHandler)
	listenAndServe()
}

func liveMode() {
	http.HandleFunc("/request", requestHandler)
	http.HandleFunc("/receipt", receiptHandler)
	http.HandleFunc("/invoice", invoiceHandler)

	http.HandleFunc("/bank-pay", bankPayHandler)
	http.HandleFunc("/bank-payment-success", bankPaymentSuccessHandler)
	http.HandleFunc("/bank-success", bankSuccessHandler)

	http.HandleFunc("/card-not-supported", cardNotSupportedHandler)
	http.HandleFunc("/card-pay", cardPayHandler)
	http.HandleFunc("/card-payment", cardPaymentHandler)
	http.HandleFunc("/card-payment-success", cardPaymentSuccessHandler)
	http.HandleFunc("/card-success", cardSuccessHandler)

	http.HandleFunc("/invoice-paid", invoicePaidHandler)

	http.HandleFunc("/linked-account", linkedAccountHandler)
	http.HandleFunc("/linked-card", linkedCardHandler)

	http.HandleFunc("/plaid-details", plaidDetailshandler)
	http.HandleFunc("/transfer", transferHandler)

	// handling for invoice 2.0 payments
	// card-payment-invoice
	http.HandleFunc("/invoice-request", invoiceRequestHandler)
	http.HandleFunc("/invoice-view", invoiceViewHandler)
	http.HandleFunc("/invoice-receipt", invoiceReceiptHandler)

	http.HandleFunc("/healthcheck.html", healthCheckHandler)
	http.HandleFunc("/", http404Handler) // all pages that are 404
	http.HandleFunc("/vue.js", vueJsHandler)
	http.HandleFunc("/main.js", mainJsHandler)
	http.HandleFunc("/.well-known/apple-developer-merchantid-domain-association", appleFileHandler)

	listenAndServe()
}

func listenAndServe() {
	containerPort := os.Getenv("CONTAINER_LISTEN_PORT")
	err := http.ListenAndServeTLS(fmt.Sprintf(":%s", containerPort), "./ssl/cert.pem", "./ssl/key.pem", nil)
	if err != nil {
		log.Fatal("ListenAndServeTLS: ", err)
	}
}
