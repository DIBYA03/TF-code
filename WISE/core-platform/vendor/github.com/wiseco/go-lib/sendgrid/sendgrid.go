package sendgrid

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/sendgrid/sendgrid-go"
	sg "github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
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
	CCEmails      []CCEmail
}

type CCEmail struct {
	Name  string
	Email string
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
	CCEmails      []CCEmail
	Attachment    []EmailAttachment
}

type EmailAttachment struct {
	ContentType string
	FileName    string
	ContentID   string
	Attachment  string
	Disposition string
}

//GetDisposition returns attachment disposition type, defautls to attachment
func (e *EmailAttachment) GetDisposition() string {
	if e.Disposition == "inline" {
		return e.Disposition
	}
	return "attachment"
}

//EmailResponse response from sending an email
type EmailResponse struct {
	MessageID string
}

//ValidationRequest struct describing params needed to validate an email
type ValidationRequest struct {
	Email string `json:"email"`
}

//ValidationResult container for ValidationResponse
type ValidationResult struct {
	ValidationResponse ValidationResponse `json:"result"`
}

//ValidationResponse holds the validation data we pass back
type ValidationResponse struct {
	Email     string           `json:"email"`
	Verdict   string           `json:"verdict"`
	Score     float64          `json:"score"`
	IPAddress string           `json:"ip_address"`
	Checks    ValidationChecks `json:"checks"`
}

//ValidationChecks some extra validation fields
type ValidationChecks struct {
	Domain     ValidationDomain     `json:"domain"`
	LocalPart  ValidationLocalPart  `json:"local_part"`
	Additional ValidationAdditional `json:"additional"`
}

//ValidationDomain extra valiation fields
type ValidationDomain struct {
	HasValidAddresssyntax        bool `json:"has_valid_address_syntax`
	HasMxOrARecord               bool `json:"has_mx_or_a_record"`
	IsSuspectedDisposableAddress bool `json:"is_suspected_disposable_address"`
}

//ValidationLocalPart extra valiation fields
type ValidationLocalPart struct {
	IsSuspectedRoleAddress bool `json:"is_suspected_role_address"`
}

//ValidationAdditional extra valiation fields
type ValidationAdditional struct {
	HasKnownBounces     bool `json:"has_known_bounces"`
	HasSuspectedBounces bool `json:"has_suspected_bounces"`
}

type sendGrid struct{}

//SendGrid is an interface that describes our sendgrid package's public methods
type SendGrid interface {
	GetValidation(ValidationRequest) (*ValidationResponse, error)
	SendEmail(request EmailRequest) (*EmailResponse, error)
	SendAttachmentEmail(request EmailAttachmentRequest) (*EmailResponse, error)
}

//NewSendGrid returns a new SendGrid interface
func NewSendGrid() SendGrid {
	return sendGrid{}
}

//GetValidation returns validation results for an email address
func (s sendGrid) GetValidation(vr ValidationRequest) (*ValidationResponse, error) {
	var vres ValidationResult

	apiKey, err := getApIKey("SENDGRID_API_VALIDATION_KEY")

	if err != nil {
		return nil, err
	}

	r := sg.GetRequest(apiKey, "/v3/validations/email", "")
	r.Method = "POST"

	body, err := json.Marshal(vr)

	if err != nil {
		return nil, err
	}

	r.Body = body

	resp, err := sg.MakeRequest(r)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(resp.Body), &vres)

	if err != nil {
		return nil, err
	}

	return &vres.ValidationResponse, err
}

func getApIKey(s string) (string, error) {
	ak := os.Getenv(s)

	if ak == "" {
		return "", fmt.Errorf("%s var is missing", s)
	}

	return ak, nil
}

func (s sendGrid) SendAttachmentEmail(request EmailAttachmentRequest) (*EmailResponse, error) {

	apiKey, err := getApIKey("SENDGRID_API_KEY")

	if err != nil {
		return nil, err
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
	if len(request.CCEmails) > 0 {
		for _, cc := range request.CCEmails {
			p.AddCCs(mail.NewEmail(cc.Name, cc.Email))
		}

	}
	email.AddPersonalizations(p)

	for _, attachment := range request.Attachment {
		a := mail.NewAttachment()
		a.SetContent(attachment.Attachment)
		a.SetType(attachment.ContentType)
		a.SetFilename(attachment.FileName)
		a.SetDisposition(attachment.GetDisposition())
		a.SetContentID(attachment.ContentID)

		email.AddAttachment(a)
	}

	client := sendgrid.NewSendClient(apiKey)

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
			MessageID: messageID,
		}, nil
	}
}

func (s sendGrid) SendEmail(request EmailRequest) (*EmailResponse, error) {

	if len(os.Getenv("SENDGRID_API_KEY")) == 0 {
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

	if len(request.CCEmails) > 0 {
		for _, cc := range request.CCEmails {
			p.AddCCs(mail.NewEmail(cc.Name, cc.Email))
		}

	}
	email.AddPersonalizations(p)

	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))

	response, err := client.Send(email)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var messageID string
	IDs := response.Headers["X-Message-Id"]
	if len(IDs) > 0 {
		messageID = IDs[0]
	}

	return &EmailResponse{
		MessageID: messageID,
	}, nil
}
