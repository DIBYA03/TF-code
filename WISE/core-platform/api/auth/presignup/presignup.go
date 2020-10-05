package presignup

import (
	"encoding/json"
	"errors"
	"net/mail"

	"github.com/aws/aws-lambda-go/events"
	"github.com/ttacon/libphonenumber"
	"github.com/wiseco/core-platform/api/auth"
	usersrv "github.com/wiseco/core-platform/services/user"
)

func HandleCognitoPreSignupRequest(event events.CognitoEventUserPoolsPreSignup) (events.CognitoEventUserPoolsPreSignup, error) {
	errorList := []string{}
	errorList = append(errorList, auth.ValidatePhone(event)...)

	if len(errorList) > 0 {
		errorBytes, _ := json.Marshal(errorList)
		return event, errors.New(string(errorBytes))
	}

	// Check for phone/email
	phone, ok := auth.GetUserAttributeByKey(event.Request.UserAttributes, auth.CognitoPhoneNumber)
	if !ok || phone == nil {
		errorList = append(errorList, auth.PhoneMissing)
	}

	_, err := libphonenumber.Parse(*phone, "")
	if err != nil {
		errorList = append(errorList, auth.PhoneInvalid)
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

	// Check for email or phone match and reject if found either
	_, err = usersrv.NewUserService(auth.NewPreSignUpSourceRequest(event)).GetUserIDWithPhone(*phone)
	if err == nil {
		errorBytes, _ := json.Marshal(append(errorList, auth.UserPhoneExists))
		return event, errors.New(string(errorBytes))
	}

	return event, nil
}
