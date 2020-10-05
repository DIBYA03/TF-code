package preauthentication

import (
	"errors"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/wiseco/core-platform/api/csp/auth"
	cspusersrv "github.com/wiseco/core-platform/services/csp/cspuser"
)

// HandleCognitoPostConfirmRequest is used on postsignup
func HandleCognitoPreAuthenticationRequest(event events.CognitoEventUserPoolsPreAuthentication) (events.CognitoEventUserPoolsPreAuthentication, error) {
	if event.TriggerSource != "PreAuthentication_Authentication" {
		return event, nil
	}

	if err := validateUser(event); err != nil {
		return event, err
	}

	return event, nil
}

func validateUser(event events.CognitoEventUserPoolsPreAuthentication) error {

	email, ok := auth.GetUserAttributeByKey(event.Request.UserAttributes, auth.CognitoEmail)
	if !ok {
		return errors.New(auth.EmailMissing)
	}

	// Get the user from db
	user, err := cspusersrv.NewUserService(auth.NewPreAuthenticationSourceRequest(event)).GetUserByEmail(*email)
	if err != nil {
		return err
	}

	// Check if user is active to use CSP
	if !user.Active {
		return fmt.Errorf("%s is not active to use csp", *email)
	}

	// User is active, let's make sure profile pic is updated from Google
	picture, ok := auth.GetUserAttributeByKey(event.Request.UserAttributes, auth.CognitoPicture)
	if !ok {
		// No picture from Google
		return nil
	}

	// Update picture if needed
	if user.Picture != *picture {
		user.Picture = *picture

		_, err = cspusersrv.NewUserService(auth.NewPreAuthenticationSourceRequest(event)).Update(*user)
		if err != nil {
			log.Print(err)
			return nil
		}
	}

	return nil
}
