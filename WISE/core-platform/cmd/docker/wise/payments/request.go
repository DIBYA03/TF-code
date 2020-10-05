package main

import "net/http"

// RequestHandler handles the http request for requests page
func requestHandler(w http.ResponseWriter, r *http.Request) {

	loadHTMLTemplate(HTMLTemplateIndex, w, r)
}
