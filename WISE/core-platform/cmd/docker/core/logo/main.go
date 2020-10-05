package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/wiseco/core-platform/partner/service/clearbit"
)

func main() {

	http.HandleFunc("/logo", func(w http.ResponseWriter, r *http.Request) {

		names, ok := r.URL.Query()["name"]
		if !ok || names[0] == "" {
			log.Println("missing name param")
			w.Header()["Content-Type"] = []string{"application/json"}
			w.WriteHeader(400)
			w.Write([]byte(`{"error": "missing name param"}`))
			return
		}

		logo, err := clearbit.GetLogo(names[0])
		if err == nil && logo.URL != "" {
			resp, err := http.Get(logo.URL)
			if err == nil && resp.StatusCode < 400 {
				w.Header()["Content-Type"] = resp.Header["Content-Type"]
				w.Header()["Content-Length"] = resp.Header["Content-Length"]

				b, _ := ioutil.ReadAll(resp.Body)
				w.Write(b)
				return
			}
		}

		/* TODO: Add mcc fallback
		   var mcc string
		   mccs, ok := r.URL.Query()["mcc"]
		   if ok && mccs[0] != "" {
		       mcc = mccs[0]
		   } */

		w.Header()["Content-Type"] = []string{"application/json"}
		w.WriteHeader(404)
		w.Write([]byte(`{"error": "resource not found"}`))
	})

	// healthcheck
	http.HandleFunc("/healthcheck.html", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "ok")
	})

	containerPort := os.Getenv("CONTAINER_LISTEN_PORT")
	err := http.ListenAndServeTLS(fmt.Sprintf(":%s", containerPort), "./ssl/cert.pem", "./ssl/key.pem", nil)
	if err != nil {
		log.Fatal("ListenAndServeTLS: ", err)
	}
}
