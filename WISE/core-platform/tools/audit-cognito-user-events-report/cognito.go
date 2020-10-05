package main

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

func getCognitoUserAttr(name string, attributes []*cognitoidentityprovider.AttributeType) string {
	for _, attr := range attributes {
		if *attr.Name == name {
			return *attr.Value
		}
	}

	return "not found"
}

// either send a phone number or empty string for all users
func getCognitoUserList(svc *cognitoidentityprovider.CognitoIdentityProvider, userPoolID string, filter string) []cognitoUser {
	// Simple way to get all users if needed
	if filter != "" {
		filter = fmt.Sprintf("phone_number = \"%s\"", filter)
	}

	input := &cognitoidentityprovider.ListUsersInput{
		AttributesToGet: []*string{
			aws.String("phone_number"),
		},
		Filter:     aws.String(filter),
		UserPoolId: aws.String(userPoolID),
	}

	var users []cognitoUser
	err := svc.ListUsersPages(input, func(
		page *cognitoidentityprovider.ListUsersOutput, lastPage bool) bool {
		for _, u := range page.Users {
			users = append(users, cognitoUser{
				Username: *u.Username,
				Phone:    getCognitoUserAttr("phone_number", u.Attributes),
			})
		}

		return !lastPage
	})

	if err != nil {
		log.Panic("error getting user list by number:", err.Error())
	}

	if len(users) < 1 {
		log.Print("no users found with filter:", filter)
		return []cognitoUser{}
	}

	return users
}

func getCognitoAuthEvents(svc *cognitoidentityprovider.CognitoIdentityProvider, userPoolID string, username string) []authEvent {
	log.Println("getting auth events for", username)

	input := &cognitoidentityprovider.AdminListUserAuthEventsInput{
		UserPoolId: aws.String(userPoolID),
		Username:   aws.String(username),
	}

	var events []authEvent
	err := svc.AdminListUserAuthEventsPages(input, func(
		page *cognitoidentityprovider.AdminListUserAuthEventsOutput, lastPage bool) bool {
		events = append(events, processEvents(page.AuthEvents)...)
		return !lastPage
	})
	if err != nil {
		log.Println("error getting auth events for", username, ":", err.Error())
		return []authEvent{}
	}

	return events
}
