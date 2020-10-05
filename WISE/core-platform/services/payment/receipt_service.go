package payment

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/wiseco/go-lib/id"

	"github.com/jmoiron/sqlx"
	"github.com/wiseco/core-platform/partner/service/sendgrid"
	"github.com/wiseco/core-platform/services"
	"github.com/wiseco/core-platform/services/data"
	"github.com/wiseco/core-platform/services/document"
	"github.com/wiseco/core-platform/services/pdf"
	"github.com/wiseco/core-platform/shared"

	t "github.com/wiseco/core-platform/partner/service/twilio"
)

type receiptDatastore struct {
	sourceReq *services.SourceRequest
	*sqlx.DB
}

type ReceiptService interface {
	GetSignedURL(receiptID string, businessID shared.BusinessID) (*string, error)
	GetSignedURLByRequestID(requestID, businessID string) (*string, error)
	GetReceiptURLForInvoice(receiptID string, businessID shared.BusinessID) (*string, error)
	IsPOSInvoice(receiptID string, businessID shared.BusinessID) bool
}

func NewReceiptService(r services.SourceRequest) ReceiptService {
	return &receiptDatastore{&r, data.DBWrite}
}

func NewReceiptServiceInternal() ReceiptService {
	return &receiptDatastore{nil, data.DBWrite}
}

func (db *receiptDatastore) GetSignedURL(receiptID string, businessID shared.BusinessID) (*string, error) {
	var key string
	notFound := services.ErrorNotFound{}.New("")

	err := db.Get(&key, "SELECT storage_key from business_receipt WHERE id = $1 AND business_id = $2", receiptID, businessID)
	if err != nil && err == sql.ErrNoRows {
		return nil, notFound
	}

	if key == "" {
		return nil, notFound
	}

	if err != nil {
		log.Printf("Error getting document storage_key  error:%v", err)
		return nil, err
	}

	storer, err := document.NewStorerFromKey(key)

	if err != nil {
		log.Printf("Error creating storer  error:%v", err)
		return nil, err
	}
	url, err := storer.SignedUrl()
	if url == nil {
		log.Printf("no url url:%v err:%v", url, err)
		return nil, notFound
	}
	return url, err
}

func (db *receiptDatastore) IsPOSInvoice(receiptID string, businessID shared.BusinessID) bool {
	var receiptIDInDB string

	err := db.Get(&receiptIDInDB, `select br.id from business_receipt br join business_money_request bmr on br.request_id = bmr.id
	and bmr.request_type  = 'pos' where br.id = $1 and br.business_id = $2`, receiptID, businessID)
	if err != nil {
		return false
	}

	if receiptIDInDB == "" {
		return false
	}
	return true
}

func (db *receiptDatastore) GetReceiptURLForInvoice(receiptID string, businessID shared.BusinessID) (*string, error) {
	var key string
	notFound := services.ErrorNotFound{}.New("")

	err := db.Get(&key, "SELECT invoice_id_v2 from business_receipt WHERE id = $1 AND business_id = $2", receiptID, businessID)
	if err != nil {
		return nil, err
	}

	if key == "" {
		return nil, notFound
	}

	encodedInvoiceID := base64.RawURLEncoding.EncodeToString([]byte(fmt.Sprintf("%s%s", id.IDPrefixInvoice, key)))
	viewLink := fmt.Sprintf("%s/invoice-receipt?token=%s",
		os.Getenv("PAYMENTS_URL"), encodedInvoiceID)

	return &viewLink, err
}

func (db *receiptDatastore) GetSignedURLByRequestID(requestID, businessID string) (*string, error) {
	var key string
	notFound := services.ErrorNotFound{}.New("")

	err := db.Get(&key, "SELECT storage_key from business_receipt WHERE request_id = $1 AND business_id = $2", requestID, businessID)
	if err != nil && err == sql.ErrNoRows {
		return nil, notFound
	}

	if key == "" {
		return nil, notFound
	}

	if err != nil {
		log.Printf("Error getting document storage_key  error:%v", err)
		return nil, err
	}

	storer, err := document.NewStorerFromKey(key)

	if err != nil {
		log.Printf("Error creating storer  error:%v", err)
		return nil, err
	}
	url, err := storer.SignedUrl()
	if url == nil {
		log.Printf("no url url:%v err:%v", url, err)
		return nil, notFound
	}
	return url, err
}

func (db *paymentDatastore) GenerateReceipt(r ReceiptGenerate) (*string, string, error) {
	amount := shared.FormatFloatAmount(r.Amount)

	var customerName string
	if r.ContactBusinessName != nil && len(*r.ContactBusinessName) > 0 {
		customerName = *r.ContactBusinessName
	} else if r.ContactFirstName != nil && r.ContactLastName != nil {
		customerName = (*r.ContactFirstName + " " + *r.ContactLastName)
	}

	paidDate := (r.PaymentDate).Format("Jan _2, 2006")

	// Generate PDF
	receipt := pdf.Receipt{
		BusinessName:  r.BusinessName,
		BusinessPhone: r.BusinessPhone,
		Amount:        amount,
		Notes:         *r.Notes,
		WisePhone:     os.Getenv("WISE_SUPPORT_PHONE"),
		ReceiptNo:     r.ReceiptNumber,
		CardBrand:     *r.PaymentBrand,
		CardLast4:     *r.PaymentNumber,
		PaidDate:      paidDate,
	}

	if r.ContactEmail != nil {
		receipt.ContactEmail = r.ContactEmail
	}

	if customerName != "" {
		receipt.ContactName = &customerName
	}

	if r.InvoiceNumber != nil {
		receipt.InvoiceNo = r.InvoiceNumber
	}

	content, err := pdf.NewReceiptService(receipt).GenerateReceipt()
	if err != nil {
		log.Println(err)
		return nil, "", err
	}

	// Upload document to aws
	store, err := document.NewAWSS3DocStorageFromContent(string(r.BusinessID), document.BusinessPrefix, "application/pdf", content)
	if err != nil {
		log.Println("Error uploading receipt to s3", err)
		return nil, "", err
	}

	key, err := store.Key()
	if err != nil {
		log.Println("Error getting receipt document key", err)
		return nil, "", err
	}

	// Store in database
	receiptCreate := ReceiptCreate{
		RequestID:     r.RequestID,
		BusinessID:    r.BusinessID,
		CreatedUserID: r.UserID,
		ReceiptNumber: r.ReceiptNumber,
		ContactID:     r.ContactID,
		InvoiceID:     r.InvoiceID,
		StorageKey:    *key,
		InvoiceIdV2:   r.InvoiceIdV2,
	}

	re, err := db.saveReceipt(receiptCreate)
	if err != nil {
		log.Println(err)
		return nil, "", err
	}

	return content, re.ID, nil
}

func (db *paymentDatastore) saveReceipt(receipt ReceiptCreate) (*Receipt, error) {
	columns := []string{
		"request_id", "invoice_id", "created_user_id", "business_id", "contact_id", "receipt_number", "storage_key", "invoice_id_v2",
	}
	values := []string{
		":request_id", ":invoice_id", ":created_user_id", ":business_id", ":contact_id", ":receipt_number", ":storage_key", ":invoice_id_v2",
	}

	sql := fmt.Sprintf("INSERT INTO business_receipt(%s) VALUES(%s) RETURNING *", strings.Join(columns, ", "), strings.Join(values, ", "))

	stmt, err := db.PrepareNamed(sql)
	if err != nil {
		return nil, err
	}

	r := &Receipt{}

	err = stmt.Get(r, &receipt)
	if err != nil {
		return nil, err
	}

	return r, nil
}

// Send email to contact who made the payment
func (db *paymentDatastore) SendReceiptToCustomer(r ReceiptRequest, invoiceLink string, receiptLink string) error {

	amount := shared.FormatFloatAmount(r.Amount)

	var customerName string
	if r.ContactBusinessName != nil && len(*r.ContactBusinessName) > 0 {
		customerName = *r.ContactBusinessName
	} else if r.ContactFirstName != nil && r.ContactLastName != nil {
		customerName = (*r.ContactFirstName + " " + *r.ContactLastName)
	}

	notes := strings.Replace(*r.Notes, "\n", "<br>", -1)
	r.Notes = &notes

	subject := fmt.Sprintf(services.CustomerReceiptSubject, amount, r.BusinessName)
	body := fmt.Sprintf(services.CustomerReceiptEmail, customerName, r.BusinessName, r.PaymentDate, *r.Notes, amount, amount, r.BusinessName, r.BusinessName)

	if r.ContactEmail != nil && len(*r.ContactEmail) > 0 {
		// if viewLink is passed as non-empty, then don't attach the PDF
		if invoiceLink != "" && receiptLink != "" {
			body := fmt.Sprintf(services.CustomerReceiptEmailWithInvoiceViewLink, customerName, r.BusinessName,
				r.PaymentDate, *r.Notes, amount, invoiceLink, receiptLink, amount, r.BusinessName, r.BusinessName)

			email := sendgrid.EmailRequest{
				SenderEmail:  os.Getenv("WISE_INVOICE_EMAIL"),
				SenderName:   os.Getenv("WISE_SUPPORT_NAME"),
				ReceiverName: customerName,
				Subject:      subject,
				Body:         body,
			}

			if r.ContactEmail != nil {
				email.ReceiverEmail = *r.ContactEmail
			} else {
				email.ReceiverEmail = ""
			}
			_, err := sendgrid.NewSendGridServiceWithout().SendEmail(email)
			if err != nil {
				log.Println(err)
				return err
			}
		} else {
			a := sendgrid.EmailAttachment{
				Attachment:  *r.Content,
				ContentType: "application/pdf",
				FileName:    "Receipt - " + r.ReceiptNumber + ".pdf",
				ContentID:   "Receipt",
			}

			email := sendgrid.EmailAttachmentRequest{
				SenderEmail:  os.Getenv("WISE_INVOICE_EMAIL"),
				SenderName:   os.Getenv("WISE_SUPPORT_NAME"),
				ReceiverName: customerName,
				Subject:      subject,
				Body:         body,
				Attachment:   []sendgrid.EmailAttachment{a},
			}

			if r.ContactEmail != nil {
				email.ReceiverEmail = *r.ContactEmail
			} else {
				email.ReceiverEmail = ""
			}

			_, err := sendgrid.NewSendGridServiceWithout().SendAttachmentEmail(email)
			if err != nil {
				log.Println(err)
				return err
			}
		}

	}

	// Send SMS only for Shopify request
	if r.ContactPhone != nil && len(*r.ContactPhone) > 0 && r.isShopifyRequest() {
		body = strings.ReplaceAll(body, "<br />", "")

		smsReq := t.SMSRequest{
			Body:  body,
			Phone: *r.ContactPhone,
		}

		err := t.NewTwilioService().SendSMS(smsReq)
		if err != nil {
			log.Println(err)
			return err
		}
	}

	return nil
}

func (db *paymentDatastore) SendCardReaderReceiptToCustomer(r ReceiptRequest) error {

	amount := shared.FormatFloatAmount(r.Amount)

	notes := strings.Replace(*r.Notes, "\n", "<br>", -1)
	r.Notes = &notes

	subject := fmt.Sprintf(services.CustomerCardReaderReceiptSubject, amount, r.BusinessName)
	body := fmt.Sprintf(services.CustomerCardReaderReceiptEmail, amount, r.BusinessName, r.PaymentDate, *r.Notes, amount, amount, r.BusinessName, r.BusinessName)

	a := sendgrid.EmailAttachment{
		Attachment:  *r.Content,
		ContentType: "application/pdf",
		FileName:    "Receipt - " + r.ReceiptNumber + ".pdf",
		ContentID:   "Receipt",
	}

	email := sendgrid.EmailAttachmentRequest{
		SenderEmail: os.Getenv("WISE_INVOICE_EMAIL"),
		SenderName:  os.Getenv("WISE_SUPPORT_NAME"),
		Subject:     subject,
		Body:        body,
		Attachment:  []sendgrid.EmailAttachment{a},
	}

	if r.ContactEmail != nil {
		email.ReceiverEmail = *r.ContactEmail
	} else {
		email.ReceiverEmail = ""
	}

	_, err := sendgrid.NewSendGridServiceWithout().SendAttachmentEmail(email)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// Send email to business
func (db *paymentDatastore) SendReceiptToBusiness(r ReceiptRequest, viewLink string) error {
	amount := shared.FormatFloatAmount(r.Amount)

	var customerName string
	if r.ContactBusinessName != nil && len(*r.ContactBusinessName) > 0 {
		customerName = *r.ContactBusinessName
	} else if r.ContactFirstName != nil && r.ContactLastName != nil {
		customerName = (*r.ContactFirstName + " " + *r.ContactLastName)
	}

	notes := strings.Replace(*r.Notes, "\n", "<br>", -1)
	r.Notes = &notes

	subject := fmt.Sprintf(services.BusinessReceiptSubject, amount, customerName)
	if viewLink != "" {
		body := fmt.Sprintf(services.BusinessReceiptEmailWithInvoiceViewLink, r.BusinessName, amount, customerName, r.PaymentDate, *r.Notes, amount, viewLink, amount, r.BusinessName)

		email := sendgrid.EmailRequest{
			SenderEmail:   os.Getenv("WISE_INVOICE_EMAIL"),
			SenderName:    os.Getenv("WISE_SUPPORT_NAME"),
			ReceiverEmail: r.BusinessEmail,
			ReceiverName:  r.BusinessName,
			Subject:       subject,
			Body:          body,
		}

		_, err := sendgrid.NewSendGridServiceWithout().SendEmail(email)
		if err != nil {
			log.Println(err)
			return err
		}
		return nil
	}
	body := fmt.Sprintf(services.BusinessReceiptEmail, r.BusinessName, amount, customerName, r.PaymentDate, *r.Notes, amount, amount, r.BusinessName)

	a := sendgrid.EmailAttachment{
		Attachment:  *r.Content,
		ContentType: "application/pdf",
		FileName:    "Receipt - " + r.ReceiptNumber + ".pdf",
		ContentID:   "Receipt",
	}

	email := sendgrid.EmailAttachmentRequest{
		SenderEmail:   os.Getenv("WISE_INVOICE_EMAIL"),
		SenderName:    os.Getenv("WISE_SUPPORT_NAME"),
		ReceiverEmail: r.BusinessEmail,
		ReceiverName:  r.BusinessName,
		Subject:       subject,
		Body:          body,
		Attachment:    []sendgrid.EmailAttachment{a},
	}

	_, err := sendgrid.NewSendGridServiceWithout().SendAttachmentEmail(email)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// Send email to business
func (db *paymentDatastore) SendCardReaderReceiptToBusiness(r ReceiptRequest) error {
	amount := shared.FormatFloatAmount(r.Amount)

	notes := strings.Replace(*r.Notes, "\n", "<br>", -1)
	r.Notes = &notes

	subject := fmt.Sprintf(services.BusinessCardReaderReceiptSubject, amount)
	body := fmt.Sprintf(services.BusinessCardReaderReceiptEmail, r.BusinessName, amount, r.PaymentDate, *r.Notes, amount, amount, r.BusinessName)

	a := sendgrid.EmailAttachment{
		Attachment:  *r.Content,
		ContentType: "application/pdf",
		FileName:    "Receipt - " + r.ReceiptNumber + ".pdf",
		ContentID:   "Receipt",
	}

	email := sendgrid.EmailAttachmentRequest{
		SenderEmail:   os.Getenv("WISE_INVOICE_EMAIL"),
		SenderName:    os.Getenv("WISE_SUPPORT_NAME"),
		ReceiverEmail: r.BusinessEmail,
		ReceiverName:  r.BusinessName,
		Subject:       subject,
		Body:          body,
		Attachment:    []sendgrid.EmailAttachment{a},
	}

	_, err := sendgrid.NewSendGridServiceWithout().SendAttachmentEmail(email)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
