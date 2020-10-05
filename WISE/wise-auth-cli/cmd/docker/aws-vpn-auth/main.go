package main

import (
	"net/http"

	"github.com/wiseco/go-lib/log"
	"github.com/wiseco/go-lib/router"
	"github.com/wiseco/wise-auth-cli/internal/controller"
)

func main() {
	l := log.NewLogger()

	l.Info("### Starting up aws auth cli service ###")

	awsAuthCLICon := controller.NewAWSAuthCLIController()

	r := router.NewRouter()

	r.Handle("/css", http.StripPrefix("/css/", http.FileServer(http.Dir("static/css"))))
	r.Handle("/images", http.StripPrefix("/images/", http.FileServer(http.Dir("static/images"))))
	r.Handle("/fonts", http.StripPrefix("/fonts/", http.FileServer(http.Dir("static/fonts"))))
	r.Handle("/js", http.StripPrefix("/js/", http.FileServer(http.Dir("static/js"))))

	r.HandleHealthCheck()

	//Google SAML requests
	r.Get("/", awsAuthCLICon.GETAWSAuthCLI)
	r.Post("/saml", awsAuthCLICon.POSTAWSAuthCLI)

	err := r.ListenAndServeTLS()
	if err != nil {
		l.ErrorD("Something went wrong with the server:", log.Fields{"err": err.Error()})
	}
}
