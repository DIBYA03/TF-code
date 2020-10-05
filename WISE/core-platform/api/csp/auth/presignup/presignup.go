package presignup

import (
	"encoding/json"
	"errors"
	"net/mail"

	"github.com/aws/aws-lambda-go/events"
	"github.com/wiseco/core-platform/api/csp/auth"
	cspusersrv "github.com/wiseco/core-platform/services/csp/cspuser"
)

func HandleCognitoPreSignupRequest(event events.CognitoEventUserPoolsPreSignup) (events.CognitoEventUserPoolsPreSignup, error) {
	errorList := []string{}
	errorList = append(errorList, auth.ValidateEmail(event)...)

	if len(errorList) > 0 {
		errorBytes, _ := json.Marshal(errorList)
		return event, errors.New(string(errorBytes))
	}

	// Check email
	email, ok := auth.GetUserAttributeByKey(event.Request.UserAttributes, auth.CognitoEmail)
	if ok {
		_, err := mail.ParseAddress(*email)
		if err != nil {
			errorList = append(errorList, auth.EmailInvalid)
		}
	}

	if len(errorList) > 0 {
		errorBytes, _ := json.Marshal(errorList)
		return event, errors.New(string(errorBytes))
	}

	// Check for email match and reject if found
	_, err := cspusersrv.NewUserService(auth.NewPreSignUpSourceRequest(event)).GetUserByEmail(*email)
	if err == nil {
		errorBytes, _ := json.Marshal(append(errorList, auth.UserEmailExists))
		return event, errors.New(string(errorBytes))
	}

	return event, nil
}
