package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

func cardPaymentHandler(w http.ResponseWriter, r *http.Request) {
	loadHTMLTemplate(HTMLTemplateCardPayment, w, r)
}

func plaidDetailshandler(w http.ResponseWriter, r *http.Request) {
	loadHTMLTemplate(HTMLTemplatePlaid, w, r)
}

func cardSuccessHandler(w http.ResponseWriter, r *http.Request) {
	loadHTMLTemplate(HTMLTemplateCardSuccess, w, r)
}

func bankSuccessHandler(w http.ResponseWriter, r *http.Request) {
	loadHTMLTemplate(HTMLTemplateBankAccountSuccess, w, r)
}

func invoicePaidHandler(w http.ResponseWriter, r *http.Request) {
	loadHTMLTemplate(HTMLTemplateInvoicePaid, w, r)
}

func transferHandler(w http.ResponseWriter, r *http.Request) {
	loadTransferRequestHTMLTemplate(HTMLTemplateGetStarted, w, r)
}

func cardPayHandler(w http.ResponseWriter, r *http.Request) {
	loadTransferRequestHTMLTemplate(HTMLTemplateCardPay, w, r)
}

func bankPayHandler(w http.ResponseWriter, r *http.Request) {
	loadTransferRequestHTMLTemplate(HTMLTemplateBankPay, w, r)
}

func linkedAccountHandler(w http.ResponseWriter, r *http.Request) {
	loadTransferRequestHTMLTemplate(HTMLTemplateLinkedAccount, w, r)
}

func linkedCardHandler(w http.ResponseWriter, r *http.Request) {
	loadTransferRequestHTMLTemplate(HTMLTemplateLinkedCard, w, r)
}

func bankPaymentSuccessHandler(w http.ResponseWriter, r *http.Request) {
	loadTransferRequestHTMLTemplate(HTMLTemplateBankPaymentSuccess, w, r)
}

func cardNotSupportedHandler(w http.ResponseWriter, r *http.Request) {
	loadTransferRequestHTMLTemplate(HTMLTemplateCardNotSupported, w, r)
}

func cardPaymentSuccessHandler(w http.ResponseWriter, r *http.Request) {
	loadTransferRequestHTMLTemplate(HTMLTemplateCardPaymentSuccess, w, r)
}

// Invoice 2.0 handlers

// invoiceViewHandler handles the http request for invoice sharable link
func invoiceViewHandler(w http.ResponseWriter, r *http.Request) {
	loadHTMLTemplate(HTMLTemplateInvoiceView, w, r)
}

// invoiceRequestHandler handles the http request for new invoice payment page
func invoiceRequestHandler(w http.ResponseWriter, r *http.Request) {
	loadHTMLTemplate(HTMLTemplateInvoiceRequestIndex, w, r)
}

// invoiceReceiptHandler  handles http request for invoice receipt
func invoiceReceiptHandler(w http.ResponseWriter, r *http.Request) {
	loadHTMLTemplate(HTMLTemplateInvoiceSuccess, w, r)
}

func favIconIcoHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "favicon.ico")
}

func favIconPngHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "favicon.png")
}

func vueJsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/javascript")
	data, err := ioutil.ReadFile("vue.js")
	if err != nil {
		log.Println("error loading vue.js ", err)
	}

	w.Write(data)
}

func mainJsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/javascript")
	data, err := ioutil.ReadFile("main.js")
	if err != nil {
		log.Println("error loading main.js ", err)
	}

	w.Write(data)
}

func appleFileHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/octet-stream")
	data, err := ioutil.ReadFile("stripe-apple-pay-domain-verification")
	if err != nil {
		log.Println("error loading file ", err)
	}

	w.Write(data)
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "ok")
}

func http404Handler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)

	tmpl := template.Must(template.ParseFiles("404.html"))
	tmpl.Execute(w, "")
}

func maintenanceHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	tmpl := template.Must(template.ParseFiles("maintenance.html"))
	tmpl.Execute(w, "")
}
