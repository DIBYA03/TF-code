package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/wiseco/core-platform/services/payment"
	"github.com/wiseco/core-platform/shared"
)

// ReceiptHandler handles the http request for receipt page
func receiptHandler(w http.ResponseWriter, r *http.Request) {
	params, ok := r.URL.Query()["token"]

	if !ok || len(params[0]) < 1 {

		data := Receipt{
			ErrorMessage: "requested url is invalid",
		}

		tmpl := template.Must(template.ParseFiles("error.html"))
		tmpl.Execute(w, data)

		return
	}

	token := params[0]

	receipt, err := payment.NewPaymentServiceInternal().GetReceiptInfo(token)
	if err != nil {
		// Return error html page
		log.Println("invalid token error", err)

		data := Receipt{
			ErrorMessage: err.Error(),
		}

		tmpl := template.Must(template.ParseFiles("error.html"))
		tmpl.Execute(w, data)

		return
	}

	amount := shared.FormatFloatAmount(receipt.Amount)

	paidDate := receipt.PaymentDate.Format("2006-01-02T15:04:05")

	cardBrand := strings.Title(strings.ToLower(*receipt.CardBrand))

	data := Receipt{
		ClientSecret: receipt.ClientSecret,
		BusinessName: receipt.BusinessName,
		Amount:       amount,
		PaymentMode:  "Card reader",
		StripeKey:    receipt.StripeKey,
		PaymentDate:  paidDate,
		CardBrand:    cardBrand,
		CardLast4:    *receipt.CardLast4,
		Address:      receipt.PurchaseAddress.City + ", " + receipt.PurchaseAddress.State + " " + receipt.PurchaseAddress.PostalCode,
		Latitude:     fmt.Sprintf("%f", receipt.PurchaseAddress.Latitude),
		Longitude:    fmt.Sprintf("%f", receipt.PurchaseAddress.Longitude),
		BusinessID:   receipt.BusinessID,
		ReceiptID:    receipt.ReceiptID,
	}

	tmpl, err := template.New("receipt.html").Funcs(template.FuncMap{
		"GenerateReceiptUrlByID": generateReceiptUrlByID,
	}).ParseFiles("receipt.html")
	if err != nil {
		log.Println("error is ", err)
	}

	//tmpl := template.Must(template.ParseFiles("receipt.html"))
	err = tmpl.Execute(w, data)
	if err != nil {
		log.Println("error executing template ", err)
	}

}

func generateReceiptUrlByID(businessID shared.BusinessID, receiptID string) string {
	url, err := payment.NewReceiptServiceInternal().GetSignedURL(receiptID, businessID)
	if err != nil {
		return ""
	}

	return *url

}
