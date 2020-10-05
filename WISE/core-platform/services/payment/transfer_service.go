package payment

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/wiseco/core-platform/partner/service/sendgrid"
	t "github.com/wiseco/core-platform/partner/service/twilio"
	"github.com/wiseco/core-platform/services"
	"github.com/wiseco/core-platform/services/banking/business"
	"github.com/wiseco/core-platform/services/contact"
	"github.com/wiseco/core-platform/services/data"
	"github.com/wiseco/core-platform/shared"
	grpcBanking "github.com/wiseco/protobuf/golang/banking"
)

type transferDatastore struct {
	sourceReq services.SourceRequest
	*sqlx.DB
}

type TransferService interface {
	GetTransferRequestInfo(string) (*TransferRequestResponse, error)
	SendTransferRequest(TransferRequestCreate) (*TransferRequest, error)
	UpdateTransferRequest(TransferRequestUpdate) error
}

func NewTransferService(r services.SourceRequest) TransferService {
	return &transferDatastore{r, data.DBWrite}
}

func (db *transferDatastore) SendTransferRequest(c TransferRequestCreate) (*TransferRequest, error) {

	t := &TransferRequest{}

	log.Println("sent request mode is ", c.RequestMode)

	val, ok := RequestModeToRequestMode[c.RequestMode]

	log.Println("map value is  ", val, ok)

	if !ok {
		return nil, errors.New("invalid request mode")
	}

	if len(os.Getenv("PAYMENTS_URL")) == 0 {
		log.Println("Missing payments url")
		return nil, errors.New("missing payments url")
	}

	// Set token expiry to 7 days(168 hours)
	expTime := time.Now().UTC().Add(time.Hour * time.Duration(168))
	c.ExpirationDate = &expTime

	columns := []string{
		"created_user_id", "business_id", "contact_id", "request_mode", "expiration_date", "amount", "notes",
	}
	values := []string{
		":created_user_id", ":business_id", ":contact_id", ":request_mode", ":expiration_date", ":amount", ":notes",
	}

	sql := fmt.Sprintf("INSERT INTO money_transfer_request(%s) VALUES(%s) RETURNING *", strings.Join(columns, ", "), strings.Join(values, ", "))

	stmt, err := db.PrepareNamed(sql)
	if err != nil {
		return nil, err
	}

	err = stmt.Get(t, &c)
	if err != nil {
		return nil, err
	}

	b := BusinessContact{}

	err = db.Get(
		&b, `
		SELECT
			business.legal_name, business.dba, business_contact.first_name, business_contact.last_name,
			business_contact.business_name, business_contact.phone_number, business_contact.email,
			business_contact.contact_type FROM business
			JOIN business_contact ON business.id = business_contact.business_id
			WHERE business.id = $1 AND business_contact.id = $2`, c.BusinessID, c.ContactID,
	)
	if err != nil {
		log.Println(err)
		return nil, errors.New("unable to find business")
	}

	switch c.RequestMode {
	case ReqeuestModeSMS:
		err = sendSMS(t, b)
		if err != nil {
			return nil, err
		}
	case RequestModeEmail:
		err = sendEmail(t, b)
		if err != nil {
			return nil, err
		}
	case RequestModeSMSEmail:
		err = sendSMS(t, b)
		if err != nil {
			return nil, err
		}

		err = sendEmail(t, b)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("invalid request mode")
	}

	return t, nil
}

func sendEmail(r *TransferRequest, b BusinessContact) error {
	businessName := shared.GetBusinessName(&b.LegalName, b.DBA)

	var contactName string
	switch b.ContactType {
	case contact.ContactTypePerson:
		contactName = *b.ContactFirstName + " " + *b.ContactLastName
	case contact.ContactTypeBusiness:
		contactName = *b.ContactBusinessName
	}

	url := os.Getenv("PAYMENTS_URL") + "/transfer?token=" + *r.PaymentToken

	subject := fmt.Sprintf(services.TransferRequestSubject, businessName)
	body := fmt.Sprintf(services.TransferRequestEmail, contactName, businessName, url)

	email := sendgrid.EmailRequest{
		SenderEmail:   os.Getenv("WISE_INVOICE_EMAIL"),
		SenderName:    os.Getenv("WISE_SUPPORT_NAME"),
		ReceiverName:  contactName,
		ReceiverEmail: b.ContactEmail,
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

func sendSMS(r *TransferRequest, b BusinessContact) error {
	businessName := shared.GetBusinessName(&b.LegalName, b.DBA)

	var contactName string
	switch b.ContactType {
	case contact.ContactTypePerson:
		contactName = *b.ContactFirstName + " " + *b.ContactLastName
	case contact.ContactTypeBusiness:
		contactName = *b.ContactBusinessName
	}

	url := os.Getenv("PAYMENTS_URL") + "/transfer?token=" + *r.PaymentToken
	body := fmt.Sprintf(services.TransferRequestSMS, contactName, businessName, url)
	request := t.SMSRequest{
		Body:  body,
		Phone: b.ContactPhone,
	}

	err := t.NewTwilioService().SendSMS(request)
	if err != nil {
		log.Println("error sending receipt sms ", err)
		return err
	}

	return nil
}

func (db *transferDatastore) GetTransferRequestInfo(token string) (*TransferRequestResponse, error) {
	// Fetch payment details
	p, err := db.getTransferRequestByToken(token)
	if err != nil {
		log.Println("error finding token", err)
		return nil, errors.New("This payment has been completed")
	}

	expDate := p.ExpirationDate
	currDate := time.Now()

	// Check token expiration
	if expDate.Before(currDate) {
		log.Println("Token has expired")
		return nil, errors.New("Requested url has expired")
	}

	response := TransferRequestResponse{}

	if os.Getenv("USE_BANKING_SERVICE") == "true" {
		err = db.Get(
			&response,
			`SELECT money_transfer_request.amount, money_transfer_request.notes, money_transfer_request.id "money_transfer_request.id",
		money_transfer_request.contact_id, business.owner_id, business.id "business.id",
		business.legal_name, business.dba, business_contact.first_name "business_contact.first_name",
		business_contact.last_name "business_contact.last_name", business_contact.business_name "business_contact.business_name",
		business_contact.contact_type "business_contact.contact_type"
		FROM money_transfer_request
		JOIN business ON money_transfer_request.business_id = business.id
		JOIN business_contact ON money_transfer_request.contact_id = business_contact.id
		WHERE money_transfer_request.payment_token = $1`, p.PaymentToken)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		blas, err := business.NewBankingLinkedAccountService()
		if err != nil {
			return nil, err
		}

		stfs := []grpcBanking.LinkedSubtype{
			grpcBanking.LinkedSubtype_LST_PRIMARY,
		}

		lbas, err := blas.List(response.BusinessID, stfs, 1, 0)
		if err != nil {
			return nil, err
		}

		if len(lbas) != 1 {
			return nil, errors.New("Could not find linked account")
		}

		lba := lbas[0]

		response.BusinessBankAccountID = lba.BusinessBankAccountId
		response.RegisteredAccountID = lba.Id
	} else {
		err = db.Get(
			&response,
			`SELECT money_transfer_request.amount, money_transfer_request.notes, money_transfer_request.id "money_transfer_request.id",
		money_transfer_request.contact_id, business.owner_id, business.id "business.id",
		business.legal_name, business.dba, business_linked_bank_account.id "business_linked_bank_account.id", 
		business_linked_bank_account.business_bank_account_id, business_contact.first_name "business_contact.first_name",
		business_contact.last_name "business_contact.last_name", business_contact.business_name "business_contact.business_name",
		business_contact.contact_type "business_contact.contact_type"
		FROM money_transfer_request
		JOIN business ON money_transfer_request.business_id = business.id 
		JOIN business_contact ON money_transfer_request.contact_id = business_contact.id 
		JOIN business_linked_bank_account ON money_transfer_request.business_id = business_linked_bank_account.business_id
		WHERE money_transfer_request.payment_token = $1 AND business_linked_bank_account.business_bank_account_id IS NOT null`, p.PaymentToken)
		if err != nil {
			log.Println(err)
			return nil, err
		}
	}

	response.BusinessName = shared.GetBusinessName(response.LegalName, response.DBA)

	// Check validity of token
	if p.PaymentToken == nil {
		log.Println("Token is empty")
		err = errors.New("This payment has been completed")
	}

	return &response, err
}

func (db *transferDatastore) getTransferRequestByToken(token string) (*TransferRequest, error) {
	p := TransferRequest{}

	err := db.Get(&p, "SELECT * FROM money_transfer_request WHERE payment_token = $1", token)
	if err != nil && err != sql.ErrNoRows {
		log.Println(err)
		return nil, err
	}

	return &p, err
}

func (db *transferDatastore) UpdateTransferRequest(u TransferRequestUpdate) error {
	var columns []string

	columns = append(columns, "payment_token = :payment_token")

	if u.MoneyTransferID != nil {
		columns = append(columns, "money_transfer_id = :money_transfer_id")
	}

	_, err := db.NamedExec(fmt.Sprintf("UPDATE money_transfer_request SET %s WHERE id = '%s'", strings.Join(columns, ", "), u.ID), u)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
