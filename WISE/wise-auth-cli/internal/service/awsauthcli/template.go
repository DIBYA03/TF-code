package awsauthcli

import (
	"html/template"
	"net/http"
	"strings"
	"time"
)

type awsAccountLoginRole struct {
	Name string
	URL  string
}

type awsAccountLogin struct {
	Name  string
	Roles []awsAccountLoginRole
}

// GenerateAWSTemplate creates the aws page on the localhost
func GenerateAWSTemplate(urls map[string]string, errors []string, expiration time.Time, w http.ResponseWriter) {
	accounts := map[string]awsAccountLogin{}

	for r, u := range urls {
		rS := strings.Split(r, "-")

		accountName := rS[0]
		roleName := rS[len(rS)-1]

		role := awsAccountLoginRole{
			Name: strings.Title(roleName),
			URL:  u,
		}

		// If there's an account already in the list, get it or create a new one
		var account awsAccountLogin
		if _, ok := accounts[accountName]; ok {
			account = accounts[accountName]
		} else {
			// initialize
			account = awsAccountLogin{
				Name: strings.Title(accountName),
			}
		}

		account.Roles = append(account.Roles, role)
		accounts[accountName] = account

	}

	data := struct {
		Profiles   map[string]awsAccountLogin
		Errors     []string
		Expiration time.Time
	}{
		accounts,
		errors,
		expiration.Local(),
	}

	tmpl := template.Must(template.ParseFiles("./templates/aws-console.html"))
	tmpl.Execute(w, data)
}
