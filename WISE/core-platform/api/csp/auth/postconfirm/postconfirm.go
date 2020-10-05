package postconfirm

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/wiseco/core-platform/api/csp/auth"
	"github.com/wiseco/core-platform/services/csp/bot"
	cspusersrv "github.com/wiseco/core-platform/services/csp/cspuser"
)

// HandleCognitoPostConfirmRequest is used on postsignup
func HandleCognitoPostConfirmRequest(event events.CognitoEventUserPoolsPostConfirmation) (events.CognitoEventUserPoolsPostConfirmation, error) {
	// Make sure we only try to create a user on signup and no other time
	if event.TriggerSource != "PostConfirmation_ConfirmSignUp" {
		return event, nil
	}

	if err := createUser(event); err != nil {
		return event, err
	}

	// email(event)
	return event, nil
}

func createUser(event events.CognitoEventUserPoolsPostConfirmation) error {

	id, ok := auth.GetUserAttributeByKey(event.Request.UserAttributes, auth.CognitoSub)
	if !ok {
		return errors.New(auth.IDMissing)
	}

	// If user is in db then don't create
	user, _ := cspusersrv.NewUserService(auth.NewPostConfirmSourceRequest(event)).GetById(*id)
	if user != nil {
		return nil
	}

	email, ok := auth.GetUserAttributeByKey(event.Request.UserAttributes, auth.CognitoEmail)
	if !ok {
		return errors.New(auth.EmailMissing)
	}

	emailVerifiedStr, ok := auth.GetUserAttributeByKey(event.Request.UserAttributes, auth.CognitoEmailVerified)
	if !ok {
		return errors.New(auth.EmailVerifiedMissing)
	}
	emailVerified, err := strconv.ParseBool(*emailVerifiedStr)
	if err != nil {
		return err
	}

	var userImage string
	userPicture, _ := auth.GetUserAttributeByKey(event.Request.UserAttributes, auth.CognitoPicture)
	if userPicture != nil {
		userImage = *userPicture
	} else {
		userImage = "https://app.wise.us/img/wise-logo.9d3d8da3.png"
	}

	// Check for email - if match don't create user
	_, err = cspusersrv.NewUserService(auth.NewPostConfirmSourceRequest(event)).GetUserByEmail(*email)
	if err == nil {
		return errors.New(auth.UserEmailExists)
	}

	// Get first, middle, and last name from `name`
	fullName, ok := auth.GetUserAttributeByKey(event.Request.UserAttributes, auth.CognitoName)
	if !ok {
		return errors.New(auth.EmailVerifiedMissing)
	}

	firstName, middleName, lastName := convertFullName(*fullName)

	// Create user
	_, err = cspusersrv.NewUserService(auth.NewPostConfirmSourceRequest(event)).Create(
		cspusersrv.CSPUser{
			CognitoID:     *id,
			FirstName:     firstName,
			MiddleName:    middleName,
			LastName:      lastName,
			Email:         email,
			EmailVerified: emailVerified,
			Picture:       userImage,
		},
	)

	if err != nil {
		return err
	}

	err = bot.SendNotification(bot.CSPUserCreated, *email)
	if err != nil {
		log.Println(err)
	}

	return fmt.Errorf("%s is not active to use csp", *email)
}

func convertFullName(name string) (string, string, string) {
	nameSlice := strings.Split(name, " ")

	firstName := nameSlice[0]
	// If there's only a first name
	if len(nameSlice) <= 1 {
		return firstName, "", ""
	}

	lastName := nameSlice[len(nameSlice)-1]
	// If there is only a first and last name
	if len(nameSlice) <= 2 {
		return firstName, "", lastName
	}

	middleName := strings.Join(nameSlice[1:len(nameSlice)-1], " ")

	return firstName, middleName, lastName
}
