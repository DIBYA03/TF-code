package main

import (
	"flag"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

func main() {
	userPoolID := flag.String("poolid", "", "cognito user pool id")

	flag.Parse()
	numbers := flag.Args()

	if *userPoolID == "" {
		log.Panic("poolid flag is required")
	}

	sess, err := session.NewSession(&aws.Config{})
	if err != nil {
		log.Panic("error in creating AWS session:", err)
	}
	svc := cognitoidentityprovider.New(sess)

	var users []cognitoUser
	if len(numbers) == 0 {
		users = getCognitoUserList(svc, *userPoolID, "")
	} else {
		for _, n := range numbers {
			users = append(users, getCognitoUserList(svc, *userPoolID, n)...)
		}
	}

	for a, u := range users {
		authEvents := getCognitoAuthEvents(svc, *userPoolID, u.Username)
		users[a].AuthEvents = authEvents
	}

	users = validateEvents(users)

	generateReport(users)
}
