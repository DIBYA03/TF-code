package sendgrid

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/wiseco/core-platform/partner/service"
)

type EmailRequest struct {
	SenderEmail   string
	SenderName    string
	ReceiverEmail string
	ReceiverName  string
	BCCEmail      string
	BCCName       string
	Subject       string
	Body          string
}

type EmailAttachmentRequest struct {
	SenderEmail   string
	SenderName    string
	ReceiverEmail string
	ReceiverName  string
	BCCEmail      string
	BCCName       string
	Subject       string
	Body          string
	Attachment    []EmailAttachment
}

type EmailAttachment struct {
	ContentType string
	FileName    string
	ContentID   string
	Attachment  string
}

type EmailResponse struct {
	MessageId string
}

type sendGridService struct {
	request service.APIRequest
}

type SendGridService interface {
	//-- Send email
	SendEmail(EmailRequest) (*EmailResponse, error)

	SendAttachmentEmail(EmailAttachmentRequest) (*EmailResponse, error)
}

func NewSendGridService(request service.APIRequest) SendGridService {
	return &sendGridService{
		request: request,
	}
}

func NewSendGridServiceWithout() SendGridService {
	return &sendGridService{}
}

func (s *sendGridService) SendEmail(request EmailRequest) (*EmailResponse, error) {

	if len(os.Getenv("SENDGRID_API_KEY")) == 0 {
		log.Println("SENDGRID_API_KEY var is missing")

		return nil, errors.New("SENDGRID_API_KEY var is missing")
	}

	email := mail.NewV3Mail()

	from := mail.NewEmail(request.SenderName, request.SenderEmail)
	email.SetFrom(from)

	email.Subject = request.Subject

	content := mail.NewContent("text/html", request.Body)
	email.AddContent(content)

	p := mail.NewPersonalization()
	tos := []*mail.Email{
		mail.NewEmail(request.ReceiverName, request.ReceiverEmail),
	}
	p.AddTos(tos...)

	// Add BCC
	if request.BCCEmail != "" {
		bccs := []*mail.Email{
			mail.NewEmail(request.BCCName, request.BCCEmail),
		}
		p.AddBCCs(bccs...)
	}
	email.AddPersonalizations(p)

	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))

	response, err := client.Send(email)
	if err != nil {
		log.Println(err)
		return nil, err
	} else {
		fmt.Println(response.StatusCode)

		var messageID string
		IDs := response.Headers["X-Message-Id"]
		if len(IDs) > 0 {
			messageID = IDs[0]
		}

		return &EmailResponse{
			MessageId: messageID,
		}, nil
	}
}

func (s *sendGridService) SendAttachmentEmail(request EmailAttachmentRequest) (*EmailResponse, error) {

	if len(os.Getenv("SENDGRID_API_KEY")) == 0 {
		log.Println("SENDGRID_API_KEY var is missing")

		return nil, errors.New("SENDGRID_API_KEY var is missing")
	}

	email := mail.NewV3Mail()

	from := mail.NewEmail(request.SenderName, request.SenderEmail)
	email.SetFrom(from)

	email.Subject = request.Subject

	content := mail.NewContent("text/html", request.Body)
	email.AddContent(content)

	p := mail.NewPersonalization()
	tos := []*mail.Email{
		mail.NewEmail(request.ReceiverName, request.ReceiverEmail),
	}
	p.AddTos(tos...)

	// Add BCC
	if request.BCCEmail != "" {
		bccs := []*mail.Email{
			mail.NewEmail(request.BCCName, request.BCCEmail),
		}
		p.AddBCCs(bccs...)
	}
	email.AddPersonalizations(p)

	for _, attachment := range request.Attachment {
		a := mail.NewAttachment()
		a.SetContent(attachment.Attachment)
		a.SetType(attachment.ContentType)
		a.SetFilename(attachment.FileName)
		a.SetDisposition("attachment")
		a.SetContentID(attachment.ContentID)

		email.AddAttachment(a)
	}

	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))

	response, err := client.Send(email)
	if err != nil {
		log.Println(err)
		return nil, err
	} else {
		fmt.Println(response.StatusCode)

		var messageID string
		IDs := response.Headers["X-Message-Id"]
		if len(IDs) > 0 {
			messageID = IDs[0]
		}

		return &EmailResponse{
			MessageId: messageID,
		}, nil
	}
}
