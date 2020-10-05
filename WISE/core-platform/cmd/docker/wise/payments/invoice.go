package main

import (
	"net/http"
	"text/template"
)

// InvoiceHandler handles the http request for invoice page
func invoiceHandler(w http.ResponseWriter, r *http.Request) {
	data := Receipt{}

	tmpl := template.Must(template.ParseFiles("invoice.html"))
	tmpl.Execute(w, data)

}
