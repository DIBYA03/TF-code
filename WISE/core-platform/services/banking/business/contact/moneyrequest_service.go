/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

// Package for business contacts
package contact

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/wiseco/core-platform/shared"
	"mvdan.cc/xurls/v2"

	"github.com/jmoiron/sqlx"
	"github.com/microcosm-cc/bluemonday"
	"github.com/pkg/errors"
	"github.com/wiseco/core-platform/auth"
	"github.com/wiseco/core-platform/partner/service"
	sendgrid "github.com/wiseco/core-platform/partner/service/sendgrid"
	stripe "github.com/wiseco/core-platform/partner/service/stripe"
	"github.com/wiseco/core-platform/services"
	"github.com/wiseco/core-platform/services/banking"
	"github.com/wiseco/core-platform/services/banking/business"
	b "github.com/wiseco/core-platform/services/business"
	"github.com/wiseco/core-platform/services/contact"
	"github.com/wiseco/core-platform/services/data"
	"github.com/wiseco/core-platform/services/document"
	"github.com/wiseco/core-platform/services/pdf"
)

type moneyRequestDatastore struct {
	sourceReq services.SourceRequest
	*sqlx.DB
}

type MoneyRequestService interface {
	// Read
	GetById(shared.PaymentRequestID, string, shared.BusinessID) (*business.MoneyRequest, error)
	GetByContactId(int, int, string, shared.BusinessID) ([]business.MoneyRequest, error)

	//Request
	Request(*business.RequestInitiate) (*business.MoneyRequest, error)
	UpdateRequestStatus(*banking.MoneyRequestUpdate) error
	UpdatePaymentStatus(*business.Payment) error

	// Webhook
	HandleWebhook(*business.Payment) error
}

func NewMoneyRequestService(r services.SourceRequest) MoneyRequestService {
	return &moneyRequestDatastore{r, data.DBWrite}
}

func (db *moneyRequestDatastore) GetById(id shared.PaymentRequestID, contactId string, businessID shared.BusinessID) (*business.MoneyRequest, error) {
	// Check access
	err := auth.NewAuthService(db.sourceReq).CheckBusinessAccess(businessID)
	if err != nil {
		return nil, err
	}

	a := business.MoneyRequest{}

	err = db.Get(&a, "SELECT * FROM business_money_request WHERE id = $1 AND contact_id = $2 AND business_id = $3", id, contactId, businessID)
	if err != nil && err != sql.ErrNoRows {

		return nil, err
	}

	return &a, err
}

func (db *moneyRequestDatastore) GetByContactId(offset int, limit int, id string, businessID shared.BusinessID) ([]business.MoneyRequest, error) {
	// Check access
	err := auth.NewAuthService(db.sourceReq).CheckBusinessAccess(businessID)
	if err != nil {
		return nil, err
	}

	rows := []business.MoneyRequest{}

	err = db.Select(&rows, "SELECT * FROM business_money_request WHERE contact_id = $1 AND business_id = $2", id, businessID)
	if err != nil && err != sql.ErrNoRows {

		return nil, err
	}

	return rows, err
}

func (db *moneyRequestDatastore) Request(request *business.RequestInitiate) (*business.MoneyRequest, error) {
	// Check access
	err := auth.NewAuthService(db.sourceReq).CheckBusinessAccess(request.BusinessID)
	if err != nil {
		return nil, err
	}

	maxAmount, err := strconv.ParseFloat(os.Getenv("WISE_CLEARING_MAX_MONEY_REQUEST_ALLOWED"), 64)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if request.Amount > maxAmount {
		e := fmt.Sprintf("Amount cannot exceed $%s", strconv.FormatFloat(maxAmount, 'f', 2, 64))
		return nil, errors.New(e)
	}

	if request.Amount == 0 {
		return nil, errors.New("Amount cannot be zero")
	}

	if string(request.Currency) == "" {
		return nil, errors.New("Currency is required")
	}

	if len(request.Notes) == 0 {
		return nil, errors.New("Notes is required")
	} else if len(strings.TrimSpace(request.Notes)) > 80 {
		return nil, errors.New("Notes cannot exceed more than 80 characters")
	}

	rxRelaxed := xurls.Relaxed()
	if len(rxRelaxed.FindString(request.Notes)) > 0 {
		return nil, errors.New("Notes cannot contain urls")
	}

	// sanitize request notes
	policy := bluemonday.StrictPolicy()
	notes := policy.Sanitize(
		request.Notes,
	)
	request.Notes = notes

	// Get business details
	b, err := b.NewBusinessService(db.sourceReq).GetById(request.BusinessID)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	} else if err != nil && err == sql.ErrNoRows {
		return nil, errors.New("business not found")
	}

	// Get contact details
	c, err := contact.NewContactService(db.sourceReq).GetById(request.ContactId, request.BusinessID)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	} else if err != nil && err == sql.ErrNoRows {
		return nil, errors.New("contact not found")
	}

	r := db.sourceReq.PartnerServiceRequest()

	// Create payment object using stripe
	p, err := stripe.NewStripeService(&r).CreatePayment(transformForPaymentService(b, request))
	if err != nil {
		return nil, err
	}

	moneyRequest := transformFromEmailService(request)

	// Default/mandatory fields
	columns := []string{
		"created_user_id", "business_id", "contact_id", "amount", "currency", "notes", "message_id",
	}
	// Default/mandatory values
	values := []string{
		":created_user_id", ":business_id", ":contact_id", ":amount", ":currency", ":notes", ":message_id",
	}

	sql := fmt.Sprintf("INSERT INTO business_money_request(%s) VALUES(%s) RETURNING *", strings.Join(columns, ", "), strings.Join(values, ", "))

	stmt, err := db.PrepareNamed(sql)
	if err != nil {
		return nil, err
	}

	m := &business.MoneyRequest{}

	err = stmt.Get(m, &moneyRequest)
	if err != nil {
		return nil, err
	}

	//-- Store payment object --//
	payment := transformFromPaymentService(p, m.ID)

	columns = []string{
		"request_id", "source_payment_id", "status", "expiration_date",
	}
	values = []string{
		":request_id", ":source_payment_id", ":status", ":expiration_date",
	}

	sql = fmt.Sprintf("INSERT INTO business_money_request_payment(%s) VALUES(%s) RETURNING *", strings.Join(columns, ", "), strings.Join(values, ", "))

	stmt, err = db.PrepareNamed(sql)
	if err != nil {
		return nil, err
	}

	requestPayment := &business.Payment{}

	err = stmt.Get(requestPayment, &payment)
	if err != nil {
		return nil, err
	}

	// Send request email
	response, err := db.sendInvoice(b, c, m, *requestPayment.PaymentToken)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// Update money request with send grid message id
	u := business.MoneyRequestIDUpdate{
		BusinessID: m.BusinessID,
	}
	u.ID = m.ID
	u.MessageId = response.MessageId
	u.ContactId = m.ContactId

	requestUpdate, err := db.UpdateMessageID(&u)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return requestUpdate, nil

}

type RequestJoin struct {
	PaymentID           string                     `db:"id"`
	PaymentStatus       string                     `db:"status"`
	SourcePaymentID     string                     `db:"source_payment_id"`
	RequestID           shared.PaymentRequestID    `db:"request_id"`
	ContactID           string                     `db:"business_contact.id"`
	InvoiceID           string                     `db:"business_invoice.id"`
	RequestStatus       banking.MoneyRequestStatus `db:"request_status"`
	UserID              shared.UserID              `db:"created_user_id"`
	BusinessID          shared.BusinessID          `db:"business_id"`
	Amount              float64                    `db:"amount"`
	FirstName           string                     `db:"consumer.first_name"`
	LastName            string                     `db:"consumer.last_name"`
	Email               string                     `db:"business.email"`
	ContactFirstName    string                     `db:"business_contact.first_name"`
	ContactLastName     string                     `db:"business_contact.last_name"`
	ContactEmail        string                     `db:"business_contact.email"`
	ContactBusinessName string                     `db:"business_contact.business_name"`
	BusinessName        string
	LegalName           *string              `db:"legal_name"`
	DBA                 services.StringArray `db:"dba"`
	BusinessPhone       string               `db:"business.phone"`
	Notes               string               `db:"notes"`
	InvoiceNumber       string               `db:"invoice_number"`
}

func (db *moneyRequestDatastore) HandleWebhook(payment *business.Payment) error {

	r := RequestJoin{}

	if payment.Status == string(business.PaymentStatusSucceeded) {

		err := db.Get(
			&r, `
			SELECT
				business_money_request.created_user_id, business_money_request.business_id, business_money_request.notes,
				business_invoice.invoice_number, business_invoice.id "business_invoice.id",
				business_money_request.amount, business_money_request.request_status,
				business_money_request_payment.id, business_money_request_payment.status, business.legal_name, business.dba, business.phone "business.phone",
				business_money_request_payment.request_id, business_money_request_payment.source_payment_id,
				business_contact.first_name "business_contact.first_name", business_contact.last_name "business_contact.last_name", business_contact.email "business_contact.email",
				business_contact.id "business_contact.id", business_contact.business_name "business_contact.business_name", business.email "business.email", 
				consumer.first_name "consumer.first_name", consumer.last_name "consumer.last_name" FROM business_money_request
				JOIN business_money_request_payment ON business_money_request.id = business_money_request_payment.request_id
				JOIN business_invoice ON business_money_request.id = business_invoice.request_id
				JOIN business_contact ON business_contact.id = business_money_request.contact_id
				JOIN wise_user ON business_money_request.created_user_id = wise_user.id
				JOIN business ON business.id = business_money_request.business_id
				JOIN consumer ON consumer.id = wise_user.consumer_id
				WHERE business_money_request_payment.source_payment_id = $1`,
			payment.Id,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				return errors.New(fmt.Sprintf("Payment with id:%s not found", payment.Id))
			}

			return err
		}

		r.BusinessName = shared.GetBusinessName(r.LegalName, r.DBA)

		// Wise clearing account
		clearingUserID := os.Getenv("WISE_CLEARING_USER_ID")
		if clearingUserID == "" {
			return errors.New("clearing user missing")
		}

		clearingBusinessID := os.Getenv("WISE_CLEARING_BUSINESS_ID")
		if clearingBusinessID == "" {
			return errors.New("clearing business missing")
		}

		clearingAccountID := os.Getenv("WISE_CLEARING_ACCOUNT_ID")
		if clearingAccountID == "" {
			return errors.New("clearing account missing")
		}

		paidDate := (*payment.PaymentDate).Format("Jan _2, 2006")

		// Check to make sure payment is not already made
		if r.RequestStatus != banking.MoneyRequestStatusComplete {
			// Email sender
			rNumber := shared.GetRandomString("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789", 8) + "-" +
				shared.GetRandomString("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789", 5)

			content, err := db.sendReceipt(payment, rNumber, r)
			if err != nil {
				log.Println("Receipt generation failed ", err)
			}

			sendReceiptToCustomer(r, rNumber, paidDate, content)

			// To pass through auth controls
			db.sourceReq.UserID = shared.UserID(clearingUserID)

			// Get wise clearing linked account number. Required for money transfer
			wiseLinkedAccount, err := business.NewLinkedAccountService(db.sourceReq).GetByAccountIDInternal(shared.BusinessID(clearingBusinessID), clearingAccountID)
			if err != nil && err != sql.ErrNoRows {
				log.Println(err)
				return err
			}

			if err != nil {
				return errors.New("Wise account needs to be registered before sending money to business")
			}

			// Get requester business' account number
			accounts, err := business.NewBankAccountService(db.sourceReq).GetByUserID(r.UserID, r.BusinessID)
			if err != nil {
				log.Println(err)
				return err
			}

			// Check for primary account
			var account business.BankAccount
			for _, acc := range *accounts {
				if acc.UsageType == business.UsageTypePrimary {
					account = acc
				}
			}

			if account.UsageType != business.UsageTypePrimary {
				return errors.New("Primary account not found")
			}

			// Check if requester is already registered by wise clearing account
			linkedAccount, err := business.NewLinkedAccountService(db.sourceReq).GetByAccountNumber(
				shared.BusinessID(clearingBusinessID),
				business.AccountNumber(account.AccountNumber),
				account.RoutingNumber,
			)
			if err != nil && err != sql.ErrNoRows {
				log.Println(err)
				return err
			}

			// Check if business is already registered by wise clearing account
			if err != nil {

				// Step 1: Create contact
				c, err := db.createContact(r.BusinessID, r.UserID, shared.UserID(clearingUserID), shared.BusinessID(clearingBusinessID), r.FirstName, r.LastName)
				if err != nil {
					log.Println(err)
					return err
				}

				// Step 2 : Wise clearing account links requester's account
				laCreate := business.ContactLinkedAccountCreate{
					UserID:              shared.UserID(clearingUserID),
					BusinessID:          shared.BusinessID(clearingBusinessID),
					ContactId:           c.ID,
					RegisteredAccountId: account.BankAccountId,
					RegisteredBankName:  string(account.BankName),
					AccountNumber:       business.AccountNumber(account.AccountNumber),
					AccountType:         banking.AccountType(account.AccountType),
					RoutingNumber:       account.RoutingNumber,
					Currency:            account.Currency,
					Permission:          banking.LinkedAccountPermissionSendAndRecieve,
				}
				la, err := NewLinkedAccountService(db.sourceReq).Create(&laCreate)
				if err != nil {
					log.Println(err)
					return err
				}

				// Step 3 : Move money
				transferInitiate := business.TransferInitiate{
					CreatedUserID:   shared.UserID(clearingUserID),
					BusinessID:      shared.BusinessID(clearingBusinessID),
					SourceAccountId: wiseLinkedAccount.Id,
					DestAccountId:   la.Id,
					Amount:          r.Amount,
					SourceType:      banking.TransferTypeAccount,
					DestType:        banking.TransferTypeAccount,
					ContactId:       la.ContactId,
					Currency:        banking.CurrencyUSD,
					MoneyRequestID:  &r.RequestID,
					Notes:           &r.Notes,
				}

				s := db.sourceReq

				_, err = NewMoneyTransferService(s).Transfer(&transferInitiate)
				if err != nil {
					log.Println(err)
					return err
				}

			} else {
				// Get contact Id
				if linkedAccount.ContactId == nil {
					return errors.New("Contact id missing for linked account")
				}

				// Move money
				transferInitiate := business.TransferInitiate{
					CreatedUserID:   shared.UserID(clearingUserID),
					BusinessID:      shared.BusinessID(clearingBusinessID),
					SourceAccountId: wiseLinkedAccount.Id,
					DestAccountId:   linkedAccount.Id,
					Amount:          r.Amount,
					SourceType:      banking.TransferTypeAccount,
					DestType:        banking.TransferTypeAccount,
					Currency:        banking.CurrencyUSD,
					ContactId:       linkedAccount.ContactId,
					MoneyRequestID:  &r.RequestID,
					Notes:           &r.Notes,
				}

				s := db.sourceReq

				_, err = NewMoneyTransferService(s).Transfer(&transferInitiate)
				if err != nil {
					log.Println(err)
					return err
				}

			}

			// Send email to business
			sendReceiptToBusiness(r, rNumber, paidDate, content)

			// Update request money status in business_money_request table
			requestUpdate := banking.MoneyRequestUpdate{
				ID:     r.RequestID,
				Status: banking.MoneyRequestStatusComplete,
			}

			err = db.UpdateRequestStatus(&requestUpdate)
			if err != nil {
				return err
			}

			receiptToken := shared.GetRandomString("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789", 16)

			// Update payment status in business_money_request_payment table and set token to null
			paymentUpdate := business.Payment{
				Id:           r.PaymentID,
				Status:       payment.Status,
				PaymentToken: nil,
				CardBrand:    payment.CardBrand,
				CardLast4:    payment.CardLast4,
				PaymentDate:  payment.PaymentDate,
				ReceiptToken: &receiptToken,
			}

			err = db.UpdatePaymentStatus(&paymentUpdate)
			if err != nil {
				return err
			}

		}

	}

	return nil

}

func (db *moneyRequestDatastore) sendInvoice(b *b.Business, c *contact.Contact, m *business.MoneyRequest, token string) (*sendgrid.EmailResponse, error) {
	// Generate invoice ID
	invoiceNo, err := db.getNextInvoiceSequence(b.ID, b.EmployerNumber)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	amount := strconv.FormatFloat(m.Amount, 'f', 2, 64)
	issueDate := m.Created.Format("Jan _2, 2006")

	// Generate PDF
	invoice := pdf.Invoice{
		BusinessName:  b.Name(),
		BusinessPhone: *b.Phone,
		Amount:        amount,
		ContactEmail:  c.Email,
		ContactName:   c.Name(),
		InvoiceNo:     *invoiceNo,
		Notes:         m.Notes,
		WisePhone:     os.Getenv("WISE_SUPPORT_PHONE"),
		IssueDate:     issueDate,
	}

	content, err := pdf.NewInvoiceService(invoice).GenerateInvoice()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// Upload document to aws
	store, err := document.NewAWSS3DocStorageFromContent(string(b.ID), document.BusinessPrefix, "application/pdf", content)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	key, err := store.Key()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// Store in database
	invoiceCreate := business.InvoiceCreate{
		RequestID:     m.ID,
		BusinessID:    m.BusinessID,
		CreatedUserID: m.CreatedUserID,
		InvoiceNumber: *invoiceNo,
		ContactId:     m.ContactId,
		StorageKey:    *key,
	}

	_, err = db.saveInvoice(invoiceCreate)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// Send email to contact
	response, err := sendInvoiceToCustomer(b, c, m, invoiceNo, token, issueDate, content)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// send email to business
	sendInvoiceToBusiness(b, c, m, invoiceNo, token, issueDate, content)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return response, nil

}

func (db *moneyRequestDatastore) sendReceipt(p *business.Payment, rNumber string, r RequestJoin) (*string, error) {
	amount := strconv.FormatFloat(r.Amount, 'f', 2, 64)

	customerName := r.ContactFirstName + " " + r.ContactLastName
	if len(r.ContactBusinessName) > 0 {
		customerName = r.ContactBusinessName
	}

	paidDate := (*p.PaymentDate).Format("Jan _2, 2006")

	// Generate PDF
	receipt := pdf.Receipt{
		BusinessName:  r.BusinessName,
		BusinessPhone: r.BusinessPhone,
		Amount:        amount,
		ContactEmail:  &r.ContactEmail,
		ContactName:   &customerName,
		InvoiceNo:     &r.InvoiceNumber,
		Notes:         r.Notes,
		WisePhone:     os.Getenv("WISE_SUPPORT_PHONE"),
		ReceiptNo:     rNumber,
		CardBrand:     *p.CardBrand,
		CardLast4:     *p.CardLast4,
		PaidDate:      paidDate,
	}

	content, err := pdf.NewReceiptService(receipt).GenerateReceipt()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// Upload document to aws
	store, err := document.NewAWSS3DocStorageFromContent(string(r.BusinessID), document.BusinessPrefix, "application/pdf", content)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	key, err := store.Key()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// Store in database
	receiptCreate := business.ReceiptCreate{
		RequestID:     r.RequestID,
		BusinessID:    r.BusinessID,
		CreatedUserID: r.UserID,
		ReceiptNumber: rNumber,
		ContactId:     r.ContactID,
		InvoiceId:     r.InvoiceID,
		StorageKey:    *key,
	}

	_, err = db.saveReceipt(receiptCreate)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return content, nil

}

func (db *moneyRequestDatastore) createContact(bID shared.BusinessID, uID shared.UserID, clearingUserID shared.UserID, clearingBusinessID shared.BusinessID, firstName, lastName string) (*contact.Contact, error) {

	// To pass through auth controls
	db.sourceReq.UserID = uID

	b, err := b.NewBusinessService(db.sourceReq).GetById(bID)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	} else if err != nil && err == sql.ErrNoRows {
		return nil, errors.New("business not found")
	}

	businessName := b.Name()
	contactCreate := contact.ContactCreate{
		UserID:       clearingUserID,
		BusinessID:   clearingBusinessID,
		Type:         contact.ContactTypeBusiness,
		BusinessName: &businessName,
		PhoneNumber:  *b.Phone,
		Email:        *b.Email,
	}

	// To pass through auth controls
	db.sourceReq.UserID = clearingUserID

	c, err := contact.NewContactService(db.sourceReq).Create(&contactCreate)
	if err != nil {
		return nil, err
	}

	return c, nil

}

func (db *moneyRequestDatastore) UpdateMessageID(u *business.MoneyRequestIDUpdate) (*business.MoneyRequest, error) {
	_, err := db.NamedExec(
		fmt.Sprintf(
			"UPDATE business_money_request SET message_id = :message_id WHERE id = '%s'",
			u.ID,
		), u,
	)

	if err != nil {
		return nil, errors.Cause(err)
	}

	m, err := db.GetById(u.ID, u.ContactId, u.BusinessID)
	if err != nil {
		return nil, errors.Cause(err)
	}

	return m, nil

}

func (db *moneyRequestDatastore) UpdateRequestStatus(u *banking.MoneyRequestUpdate) error {
	_, err := db.NamedExec(
		fmt.Sprintf(
			"UPDATE business_money_request SET request_status = :request_status WHERE id = '%s'",
			u.ID,
		), u,
	)

	if err != nil {
		return errors.Cause(err)
	}

	return nil

}

func (db *moneyRequestDatastore) UpdatePaymentStatus(u *business.Payment) error {
	_, err := db.NamedExec(
		fmt.Sprintf(
			"UPDATE business_money_request_payment SET status = :status, token = :token, card_brand = :card_brand, card_number = :card_number, payment_date = :payment_date WHERE id = '%s'",
			u.Id,
		), u,
	)

	if err != nil {
		return errors.Cause(err)
	}

	return nil

}

func sendInvoiceToCustomer(b *b.Business, c *contact.Contact, transfer *business.MoneyRequest, invoiceNumber *string,
	token string, issueDate string, attachment *string) (*sendgrid.EmailResponse, error) {

	amount := strconv.FormatFloat(transfer.Amount, 'f', 2, 64)

	url := os.Getenv("PAYMENTS_URL") + "/request?token=" + token

	transfer.Notes = strings.Replace(transfer.Notes, "\n", "<br>", -1)

	subject := fmt.Sprintf(services.CustomerInvoiceSubject, b.Name())
	body := fmt.Sprintf(services.CustomerInvoiceEmail, c.Name(), b.Name(), issueDate, transfer.Notes, amount, url, b.Name(), amount, amount, b.Name(), b.Name())

	a := sendgrid.EmailAttachment{
		Attachment:  *attachment,
		ContentType: "application/pdf",
		FileName:    "Invoice - " + *invoiceNumber + ".pdf",
		ContentID:   "Invoice",
	}

	email := sendgrid.EmailAttachmentRequest{
		SenderEmail:   os.Getenv("WISE_INVOICE_EMAIL"),
		SenderName:    os.Getenv("WISE_SUPPORT_NAME"),
		ReceiverEmail: c.Email,
		ReceiverName:  c.Name(),
		Subject:       subject,
		Body:          body,
		Attachment:    []sendgrid.EmailAttachment{a},
	}

	response, err := sendgrid.NewSendGridServiceWithout().SendAttachmentEmail(email)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return response, nil

}

func sendInvoiceToBusiness(b *b.Business, c *contact.Contact, transfer *business.MoneyRequest, invoiceNumber *string, token string,
	issueDate string, attachment *string) (*sendgrid.EmailResponse, error) {

	amount := strconv.FormatFloat(transfer.Amount, 'f', 2, 64)
	transfer.Notes = strings.Replace(transfer.Notes, "\n", "<br>", -1)

	subject := fmt.Sprintf(services.BusinessInvoiceSubject, amount, c.Name())
	body := fmt.Sprintf(services.BusinessInvoiceEmail, b.Name(), c.Name(), issueDate, transfer.Notes, amount, amount, b.Name())

	a := sendgrid.EmailAttachment{
		Attachment:  *attachment,
		ContentType: "application/pdf",
		FileName:    "Invoice - " + *invoiceNumber + ".pdf",
		ContentID:   "Invoice",
	}

	email := sendgrid.EmailAttachmentRequest{
		SenderEmail:   os.Getenv("WISE_INVOICE_EMAIL"),
		SenderName:    os.Getenv("WISE_SUPPORT_NAME"),
		ReceiverEmail: *b.Email,
		ReceiverName:  b.Name(),
		Subject:       subject,
		Body:          body,
		Attachment:    []sendgrid.EmailAttachment{a},
	}

	response, err := sendgrid.NewSendGridServiceWithout().SendAttachmentEmail(email)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return response, nil

}

func transformFromEmailService(request *business.RequestInitiate) business.MoneyRequest {

	r := business.MoneyRequest{
		BusinessID: request.BusinessID,
	}

	r.CreatedUserID = request.CreatedUserID
	r.ContactId = request.ContactId
	r.Currency = request.Currency
	r.Amount = request.Amount
	r.Notes = request.Notes
	r.MessageId = ""

	return r

}

func transformForPaymentService(b *b.Business, transfer *business.RequestInitiate) stripe.PaymentRequest {

	descriptor := b.Name()
	if len(descriptor) > 22 {
		descriptor = descriptor[:22]
	}

	// strip off characters like <>\'"*
	reg := regexp.MustCompile("[<>\\'\"*]+")
	descriptor = reg.ReplaceAllString(descriptor, "")

	return stripe.PaymentRequest{
		PaymentMethod: "card",
		Amount:        transfer.Amount,
		Currency:      service.Currency(transfer.Currency),
		Descriptor:    descriptor,
		ReceiptEmail:  os.Getenv("WISE_INVOICE_EMAIL"),
	}

}

func transformFromPaymentService(stripe *stripe.PaymentResponse, requestId shared.PaymentRequestID) business.PaymentResponse {

	p := business.PaymentResponse{
		RequestID:       requestId,
		SourcePaymentID: stripe.IntentID,
		Status:          string(stripe.Status),
	}

	// Set token expiry to 30 days(720 hours)
	expTime := time.Now().UTC().Add(time.Hour * time.Duration(720))
	p.ExpirationDate = &expTime

	return p

}

// Send email to contact who made the payment
func sendReceiptToCustomer(r RequestJoin, rNumber string, date string, attachment *string) error {

	amount := strconv.FormatFloat(r.Amount, 'f', 2, 64)

	customerName := r.ContactFirstName + " " + r.ContactLastName
	if len(r.ContactBusinessName) > 0 {
		customerName = r.ContactBusinessName
	}

	r.Notes = strings.Replace(r.Notes, "\n", "<br>", -1)

	subject := fmt.Sprintf(services.CustomerReceiptSubject, amount, r.BusinessName)
	body := fmt.Sprintf(services.CustomerReceiptEmail, customerName, r.BusinessName, date, r.Notes, amount, amount, r.BusinessName, r.BusinessName)

	a := sendgrid.EmailAttachment{
		Attachment:  *attachment,
		ContentType: "application/pdf",
		FileName:    "Receipt - " + rNumber + ".pdf",
		ContentID:   "Receipt",
	}

	email := sendgrid.EmailAttachmentRequest{
		SenderEmail:   os.Getenv("WISE_INVOICE_EMAIL"),
		SenderName:    os.Getenv("WISE_SUPPORT_NAME"),
		ReceiverEmail: r.ContactEmail,
		ReceiverName:  customerName,
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
func sendReceiptToBusiness(r RequestJoin, rNumber string, date string, attachment *string) error {

	amount := strconv.FormatFloat(r.Amount, 'f', 2, 64)

	customerName := r.ContactFirstName + " " + r.ContactLastName
	if len(r.ContactBusinessName) > 0 {
		customerName = r.ContactBusinessName
	}

	r.Notes = strings.Replace(r.Notes, "\n", "<br>", -1)

	subject := fmt.Sprintf(services.BusinessReceiptSubject, amount, customerName)
	body := fmt.Sprintf(services.BusinessReceiptEmail, r.BusinessName, amount, customerName, date, r.Notes, amount, amount, r.BusinessName)

	a := sendgrid.EmailAttachment{
		Attachment:  *attachment,
		ContentType: "application/pdf",
		FileName:    "Receipt - " + rNumber + ".pdf",
		ContentID:   "Receipt",
	}

	email := sendgrid.EmailAttachmentRequest{
		SenderEmail:   os.Getenv("WISE_INVOICE_EMAIL"),
		SenderName:    os.Getenv("WISE_SUPPORT_NAME"),
		ReceiverEmail: r.Email,
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

func (db *moneyRequestDatastore) getNextInvoiceSequence(bID shared.BusinessID, employerNumber string) (*string, error) {
	i := business.Invoice{}

	err := db.Get(&i, "SELECT * FROM business_invoice where business_id = $1 order by created desc limit 1", bID)
	if err != nil && err != sql.ErrNoRows {
		return nil, err

	}

	if err != nil {
		seq := employerNumber + "-" + lpad("1", "0", 5)
		return &seq, nil
	}

	//Split by -
	words := strings.Split(i.InvoiceNumber, "-")
	if len(words) > 1 {
		val, err := strconv.ParseInt(words[1], 10, 64)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		val = val + 1

		seq := employerNumber + "-" + lpad(strconv.FormatInt(val, 10), "0", 5)

		return &seq, nil

	}

	return nil, err
}

func lpad(s string, pad string, plength int) string {
	for i := len(s); i < plength; i++ {
		s = pad + s
	}
	return s
}

func (db *moneyRequestDatastore) saveInvoice(invoice business.InvoiceCreate) (*business.Invoice, error) {
	columns := []string{
		"request_id", "created_user_id", "business_id", "contact_id", "invoice_number", "storage_key",
	}
	values := []string{
		":request_id", ":created_user_id", ":business_id", ":contact_id", ":invoice_number", ":storage_key",
	}

	sql := fmt.Sprintf("INSERT INTO business_invoice(%s) VALUES(%s) RETURNING *", strings.Join(columns, ", "), strings.Join(values, ", "))

	stmt, err := db.PrepareNamed(sql)
	if err != nil {
		return nil, err
	}

	i := &business.Invoice{}

	err = stmt.Get(i, &invoice)
	if err != nil {
		return nil, err
	}

	return i, nil
}

func (db *moneyRequestDatastore) saveReceipt(receipt business.ReceiptCreate) (*business.Receipt, error) {
	columns := []string{
		"request_id", "invoice_id", "created_user_id", "business_id", "contact_id", "receipt_number", "storage_key",
	}
	values := []string{
		":request_id", ":invoice_id", ":created_user_id", ":business_id", ":contact_id", ":receipt_number", ":storage_key",
	}

	sql := fmt.Sprintf("INSERT INTO business_receipt(%s) VALUES(%s) RETURNING *", strings.Join(columns, ", "), strings.Join(values, ", "))

	stmt, err := db.PrepareNamed(sql)
	if err != nil {
		return nil, err
	}

	r := &business.Receipt{}

	err = stmt.Get(r, &receipt)
	if err != nil {
		return nil, err
	}

	return r, nil
}
