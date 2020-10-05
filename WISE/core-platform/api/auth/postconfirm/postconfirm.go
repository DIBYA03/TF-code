package postconfirm

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/pkg/errors"
	"github.com/wiseco/core-platform/api/auth"
	"github.com/wiseco/core-platform/api/auth/sesmailer"
	idsrv "github.com/wiseco/core-platform/identity"
	usersrv "github.com/wiseco/core-platform/services/user"
	"github.com/wiseco/core-platform/shared"
	"golang.org/x/crypto/bcrypt"
)

// GenerateToken generates verification token
func GenerateToken(email string) (*string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(email), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	hasher := md5.New()
	hasher.Write(hash)
	hexVal := hex.EncodeToString(hasher.Sum(nil))
	return &hexVal, nil
}

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
		return errors.New(auth.IdMissing)
	}

	uID, err := shared.ParseUserID(shared.StringValue(id))
	if err != nil {
		return err
	}

	// If user is in db then don't create
	user, _ := usersrv.NewUserService(auth.NewPostConfirmSourceRequest(event)).GetById(uID)
	if user != nil {
		return nil
	}

	phone, ok := auth.GetUserAttributeByKey(event.Request.UserAttributes, auth.CognitoPhoneNumber)
	if !ok {
		return errors.New(auth.PhoneMissing)
	}

	phoneVerified, _ := auth.GetUserAttributeByKey(event.Request.UserAttributes, auth.CognitoPhoneNumberVerified)
	email, ok := auth.GetUserAttributeByKey(event.Request.UserAttributes, auth.CognitoEmail)

	// Check for email or phone - if match don't create user
	_, err = usersrv.NewUserService(auth.NewPostConfirmSourceRequest(event)).GetUserIDWithPhone(*phone)
	if err == nil {
		return errors.New(auth.UserPhoneExists)
	}

	// Create identity
	identityID, err := idsrv.NewIdentityService(auth.NewPostConfirmIdentitySourceRequest(event)).Create(
		idsrv.IdentityCreate{
			ProviderID:     idsrv.ProviderID(*id),
			ProviderName:   idsrv.ProviderNameCognito,
			ProviderSource: idsrv.ProviderSource(event.UserPoolID),
			Phone:          *phone,
		},
	)
	if err != nil {
		return err
	}

	// Create user
	_, err = usersrv.NewUserService(auth.NewPostConfirmSourceRequest(event)).CreateFromAuth(
		usersrv.UserAuthCreate{
			IdentityID:    *identityID,
			Phone:         *phone,
			PhoneVerified: phoneVerified != nil && *phoneVerified == "true",
			Email:         email,
		},
	)

	if err != nil {
		//Let's roll back cognito creation
		iErr := idsrv.NewIdentityService(auth.NewPostConfirmIdentitySourceRequest(event)).Delete(*identityID)

		if iErr != nil {
			err = errors.Wrap(iErr, err.Error())
		}
	}

	return err
}

func email(event events.CognitoEventUserPoolsPostConfirmation) {
	emailVerified, exists := event.Request.UserAttributes[auth.CognitoEmailVerified]
	emailID := event.Request.UserAttributes[auth.CognitoEmail]
	htmlTemplate := ""
	if exists && emailVerified == "true" {
		htmlTemplate = "<h1>Wise account confirmed</h1><p>Your Wise account is confirmed, you can now use wise.us " +
			"<a href='https://aws.amazon.com/ses/'>Amazon SES</a> using the " +
			"<a href='https://aws.amazon.com/sdk-for-go/'>AWS SDK for Go</a>.</p>"
		//var htmlTemplate string = "<html><head><title>Wise account confirmed</title><style>h1{color:#f00;}</style></head><body><h1>Hello ${event.Request.UserAttributes[common.CognitoFirstName]},</h1><div>Thank you for activation your email.</div></body></html>"
	} else {
		verificationToken, err := GenerateToken(event.Request.UserAttributes[auth.CognitoSub])
		if err != nil {
			log.Println(err)
		}

		//var htmlTemplate string = "<html><head><title>Wise account confirmation needed</title><style>h1{color:#f00;}</style></head><body><h1>Hello ${event.Request.UserAttributes[common.CognitoFirstName]},</h1><div>You should activate your email but i don't have a link for you yet.</div></body></html>"
		htmlTemplate = "<h1>Wise account confirmation details</h1><p>Please click link below to confirm your email address" +
			"<a target='_new' href='" + fmt.Sprintf("%s/%s", os.Getenv("ACTIVATION_ENDPOINT"), *verificationToken) + "'>Verify my email</a>.</p>"

		//var htmlTemplate string = "<html><head><title>Wise account confirmation needed</title><style>h1{color:#f00;}</style></head><body><h1>Hello ${event.Request.UserAttributes[common.CognitoFirstName]},</h1><div>You should activate your email but i don't have a link for you yet.</div></body></html>"
		htmlTemplate = "<h1>Wise account confirmation details</h1><p>Please click this non existing link :D " +
			"<a href='https://aws.amazon.com/ses/'>Amazon SES</a> using the " +
			"<a href='https://aws.amazon.com/sdk-for-go/'>AWS SDK for Go</a>.</p>"
	}

	sesmailer.SetConfig("us-west-2")

	emailData := sesmailer.Email{
		To:      emailID,
		From:    os.Getenv("FROM_EMAIL"),
		Text:    "Hi this is the text message body",
		HTML:    htmlTemplate,
		Subject: os.Getenv("WELCOME_SUBJECT"),
		ReplyTo: os.Getenv("REPLY_TO"),
	}

	resp := sesmailer.SendEmail(emailData)

	log.Println("Mail Sent: ", resp)
}
