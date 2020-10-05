package custommessage

import (
	"fmt"
	"os"

	"github.com/wiseco/core-platform/api/auth/customevents"
)

// HandleCases generates verification token
func HandleCases(event customevents.CognitoEventUserPoolsCustomMessage) customevents.CognitoEventUserPoolsCustomMessage {

	eventType := event.TriggerSource

	switch eventType {
	case "CustomMessage_SignUp":
		event.Response.SmsMessage = fmt.Sprintf("%s %s", event.Request.CodeParameter, os.Getenv("SIGNUP_SMS_MESSAGE"))
		event.Response.EmailMessage = fmt.Sprintf(os.Getenv("SIGNUP_EMAIL_MESSAGE")+" %s", event.Request.CodeParameter)
		event.Response.EmailSubject = os.Getenv("SIGNUP_EMAIL_SUBJECT")
	case "CustomMessage_AdminCreateUser":
		fmt.Println("CustomMessage_AdminCreateUser")
	case "CustomMessage_ResendCode":
		event.Response.SmsMessage = fmt.Sprintf("%s %s", event.Request.CodeParameter, os.Getenv("RESEND_SMS_MESSAGE"))
		event.Response.EmailMessage = fmt.Sprintf(os.Getenv("RESEND_EMAIL_MESSAGE")+" %s", event.Request.CodeParameter)
		event.Response.EmailSubject = os.Getenv("RESEND_EMAIL_SUBJECT")
		//fmt.Println("CustomMessage_ResendCode")

	case "CustomMessage_ForgotPassword":
		event.Response.SmsMessage = fmt.Sprintf("%s %s", event.Request.CodeParameter, os.Getenv("FORGOT_SMS_MESSAGE"))
		event.Response.EmailMessage = fmt.Sprintf(os.Getenv("FORGOT_EMAIL_MESSAGE")+" %s", event.Request.CodeParameter)
		event.Response.EmailSubject = os.Getenv("FORGOT_EMAIL_SUBJECT")
		//fmt.Println("CustomMessage_ForgotPassword")

	case "CustomMessage_UpdateUserAttribute":
		fmt.Println("CustomMessage_UpdateUserAttribute")

	case "CustomMessage_VerifyUserAttribute":
		fmt.Println("CustomMessage_VerifyUserAttribute")

	case "CustomMessage_Authentication":
		event.Response.SmsMessage = fmt.Sprintf("%s %s", event.Request.CodeParameter, os.Getenv("AUTH_SMS_MESSAGE"))
		event.Response.EmailMessage = fmt.Sprintf(os.Getenv("AUTH_EMAIL_MESSAGE")+" %s", event.Request.CodeParameter)
		event.Response.EmailSubject = os.Getenv("AUTH_EMAIL_SUBJECT")
		//fmt.Println("CustomMessage_Authentication")
	default:
		fmt.Println("Unknown trigger")
	}
	return event
}

// HandleCognitoCustomMessageRequest is used on custom messages
func HandleCognitoCustomMessageRequest(event customevents.CognitoEventUserPoolsCustomMessage) (customevents.CognitoEventUserPoolsCustomMessage, error) {
	event = HandleCases(event)
	return event, nil
}
