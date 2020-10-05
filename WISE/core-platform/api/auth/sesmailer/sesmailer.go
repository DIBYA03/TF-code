package sesmailer

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

// Email options.
type Email struct {
	From    string // From source email
	To      string // To destination email(s)
	Subject string // Subject text to send
	Text    string // Text is the text body representation
	HTML    string // Html is the Html body representation
	ReplyTo string // Reply-To email(s)
}

// SetConfig aws configuration
func SetConfig(awsRegion string) {
	os.Setenv("AWS_REGION", awsRegion)
}


//	create a new aws session and returns the session var
//	@returns sess *session.Session
//
func startNewSession() *session.Session {
	sess, err := session.NewSession()
	if err != nil {
		fmt.Println("failed to create session,", err)
	}
	return sess
}

// SendEmail send the email
func SendEmail(emailData Email) *ses.SendEmailOutput {

	// start a new aws session
	sess := startNewSession()
	// start a new ses session
	svc := ses.New(sess)

	// Assemble the email.
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses: []*string{
				aws.String(emailData.To),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(emailData.HTML),
				},
				Text: &ses.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(emailData.Text),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String("UTF-8"),
				Data:    aws.String(emailData.Subject),
			},
		},
		Source: aws.String(emailData.From),
		// Uncomment to use a configuration set
		//ConfigurationSetName: aws.String(ConfigurationSet),
	}

	// Attempt to send the email.
	result, err := svc.SendEmail(input)

	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println("Email Sent to address: " + emailData.To)
	fmt.Println(result)

	return result

}
