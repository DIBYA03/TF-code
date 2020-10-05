/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

// Package for business contacts
package payment

import (
	"database/sql"
	"fmt"
	"log"
	"net/mail"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/wiseco/go-lib/id"

	"github.com/wiseco/core-platform/services/invoice"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/ttacon/libphonenumber"
	"github.com/wiseco/core-platform/auth"
	"github.com/wiseco/core-platform/partner/service/stripe"
	t "github.com/wiseco/core-platform/partner/service/twilio"
	"github.com/wiseco/core-platform/services"
	"github.com/wiseco/core-platform/services/banking/business"
	"github.com/wiseco/core-platform/services/data"
	"github.com/wiseco/core-platform/shared"
	grpcBanking "github.com/wiseco/protobuf/golang/banking"
	grpcSvc "github.com/wiseco/protobuf/golang/invoice"
)

type paymentDatastore struct {
	sourceReq *services.SourceRequest
	*sqlx.DB
}

type PaymentService interface {
	// Fetch payment details
	GetPaymentInfo(string) (*PaymentResponse, error)
	GetPaymentByRequestID(shared.PaymentRequestID) (*Payment, error)

	HandleWebhook(*Payment) error

	GenerateReceipt(r ReceiptGenerate) (*string, string, error)
	SendReceiptToBusiness(r ReceiptRequest, viewLink string) error
	SendReceiptToCustomer(r ReceiptRequest, viewLink string, receiptLink string) error

	SendCardReaderReceipt(CardReaderReceiptCreate) error
	GetReceiptInfo(string) (*CardReaderReceiptResponse, error)

	UpdatePaymentStatus(u *Payment) error
	UpdateRequestStatus(u *RequestUpdate) error

	// Invoice2.0 related methods
	CreatePaymentForInvoice(string) (*Payment, error)
	GetPaymentTokenFromInvoice(string) (*Payment, bool, error)
	GetPaymentReceiptFromInvoice(invoiceID id.InvoiceID) (*InvoicePaymentReceipt, error)
	GetPayments(requestIDs []shared.PaymentRequestID) ([]Payment, error)
}

func NewPaymentService(r services.SourceRequest) PaymentService {
	return &paymentDatastore{&r, data.DBWrite}
}

func NewPaymentServiceInternal() PaymentService {
	return &paymentDatastore{nil, data.DBWrite}
}

func (db *paymentDatastore) GetPaymentByRequestID(requestID shared.PaymentRequestID) (*Payment, error) {
	p := Payment{}

	err := db.Get(
		&p,
		`SELECT* FROM business_money_request_payment WHERE request_id = $1`, requestID)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &p, nil
}

func (db *paymentDatastore) GetPaymentInfo(token string) (*PaymentResponse, error) {
	// Fetch payment details
	p, err := db.getPaymentByToken(token)
	if err != nil {
		log.Println("error finding token", err)
		return nil, errors.New("Requested url is invalid or has expired")
	}

	// Check validity of token
	if p.PaymentToken == nil {
		log.Println("Token is empty")
		return nil, errors.New("Requested url is invalid or has expired")
	}

	expDate := p.ExpirationDate
	currDate := time.Now()

	// Check token expiration
	if expDate.Before(currDate) {
		log.Println("Token has expired")
		return nil, errors.New("Requested url is invalid or has expired")
	}

	useInvoiceService := false
	var invoiceReq *invoice.Invoice

	if p.InvoiceID != nil {
		useInvoiceService = true
		println("Fetching invoice from service")
		invoiceService, err := invoice.NewInvoiceService()
		if err != nil {
			return nil, err
		}
		invoiceReq, err = invoiceService.GetInvoiceByID(*p.InvoiceID)
		if err != nil {
			return nil, err
		}
	}
	paymentResponse := PaymentResponse{}

	if os.Getenv("USE_BANKING_SERVICE") == "true" {
		if !useInvoiceService {
			err = db.Get(
				&paymentResponse,
				`SELECT business_money_request.amount, business_money_request.notes, business_money_request.request_status, business_money_request.request_type,
		business_money_request.contact_id, business.owner_id, business.id "business.id",
		business.legal_name, business.dba, business_money_request_payment.id "business_money_request_payment.id",
		business_money_request_payment.request_id "business_money_request_payment.request_id",
		business_money_request_payment.payment_date "business_money_request_payment.payment_date"
		FROM business_money_request_payment
		JOIN business_money_request ON business_money_request_payment.request_id = business_money_request.id
		JOIN business ON business_money_request.business_id = business.id
		WHERE business_money_request_payment.payment_token = $1`, p.PaymentToken)
			if err != nil {
				return nil, err
			}
		} else {
			err = db.Get(
				&paymentResponse,
				`SELECT business.owner_id, business.id "business.id",
		business.legal_name, business.dba, business_money_request_payment.id "business_money_request_payment.id",
		business_money_request_payment.invoice_id "business_money_request_payment.invoice_id",
		business_money_request_payment.payment_date "business_money_request_payment.payment_date"
		FROM business_money_request_payment, business 
		WHERE business_money_request_payment.payment_token = $1 And business.id = $2`, p.PaymentToken, invoiceReq.BusinessID.UUIDString())
			if err != nil {
				return nil, err
			}
			if amount, ok := invoiceReq.Amount.Float64(); !ok {
				return nil, err
			} else {
				paymentResponse.Amount = amount
			}
			paymentResponse.Notes = invoiceReq.Notes
			paymentResponse.Status = GetInvoiceStatus(invoiceReq.Status)
			if invoiceReq.AllowCard {
				paymentResponse.RequestType = PaymentRequestTypeInvoiceCard
			}
			paymentResponse.MoneyRequestID = nil
			contactId := invoiceReq.ContactID.UUIDString()
			paymentResponse.ContactID = &contactId
			paymentResponse.InvoiceIPAddress = &invoiceReq.IPAddress
			paymentResponse.InvoiceID = p.InvoiceID
			paymentResponse.Status = GetInvoiceStatus(invoiceReq.Status)
			paymentResponse.InvoiceTitle = &invoiceReq.Title
		}

		blas, err := business.NewBankingLinkedAccountService()
		if err != nil {
			return nil, err
		}

		stfs := []grpcBanking.LinkedSubtype{
			grpcBanking.LinkedSubtype_LST_PRIMARY,
		}

		lbas, err := blas.List(paymentResponse.BusinessID, stfs, 1, 0)
		if err != nil {
			return nil, err
		}

		if len(lbas) != 1 {
			return nil, errors.New("Could not find linked account")
		}

		lba := lbas[0]

		paymentResponse.BusinessBankAccountID = lba.BusinessBankAccountId
		paymentResponse.RegisteredAccountID = lba.Id
	} else {
		err = db.Get(
			&paymentResponse,
			`SELECT business_money_request.amount, business_money_request.notes, business_money_request.request_status, business_money_request.request_type, 
		business_money_request.contact_id, business.owner_id, business.id "business.id",
		business.legal_name, business.dba, business_linked_bank_account.id "business_linked_bank_account.id", 
		business_linked_bank_account.business_bank_account_id, business_money_request_payment.id "business_money_request_payment.id",
		business_money_request_payment.request_id "business_money_request_payment.request_id", 
		business_money_request_payment.payment_date "business_money_request_payment.payment_date"
		FROM business_money_request_payment
		JOIN business_money_request ON business_money_request_payment.request_id = business_money_request.id
		JOIN business ON business_money_request.business_id = business.id 
		JOIN business_linked_bank_account ON business_money_request.business_id = business_linked_bank_account.business_id
		WHERE business_money_request_payment.payment_token = $1 AND business_linked_bank_account.business_bank_account_id IS NOT null`, p.PaymentToken)
		if err != nil {
			return nil, err
		}
	}

	var clientSecret *string
	if p.SourcePaymentID != nil && isPaymentIntent(*p.SourcePaymentID) {
		clientSecret, err = stripe.NewStripeService(nil).GetClientSecret(*p.SourcePaymentID)
		if err != nil {
			return nil, err
		}
	}

	paymentResponse.ClientSecret = clientSecret
	paymentResponse.StripeKey = os.Getenv("STRIPE_PUBLISH_KEY")
	paymentResponse.BusinessName = shared.GetBusinessName(paymentResponse.LegalName, paymentResponse.DBA)
	// Send back client secret
	return &paymentResponse, nil
}

func (db *paymentDatastore) CreatePaymentForInvoice(invoiceId string) (*Payment, error) {
	payment := PartnerPaymentResponse{
		InvoiceID:       invoiceId,
		Status:          string(PaymentStatusPending),
		SourcePaymentID: "",
	}
	expTime := time.Now().UTC().Add(time.Hour * time.Duration(720))
	payment.ExpirationDate = &expTime

	columns := []string{
		"invoice_id", "status", "expiration_date", "source_payment_id",
	}
	values := []string{
		":invoice_id", ":status", ":expiration_date", ":source_payment_id",
	}
	sqlIns := fmt.Sprintf("INSERT INTO business_money_request_payment(%s) VALUES(%s) RETURNING *", strings.Join(columns, ", "), strings.Join(values, ", "))
	stmt, err := db.PrepareNamed(sqlIns)
	if err != nil {
		return nil, err
	}

	requestPayment := &Payment{}
	err = stmt.Get(requestPayment, &payment)
	if err != nil {
		return nil, err
	}
	return requestPayment, nil
}

func (db *paymentDatastore) GetPaymentTokenFromInvoice(invoiceId string) (*Payment, bool, error) {
	paymentResponse := &Payment{}

	err := db.Get(
		paymentResponse,
		`SELECT invoice_id, status, expiration_date, payment_token
		FROM business_money_request_payment
		WHERE business_money_request_payment.invoice_id = $1`, invoiceId)
	if err != nil && err != sql.ErrNoRows {
		return nil, false, err
	} else if err == sql.ErrNoRows {
		return nil, false, nil
	}
	return paymentResponse, true, nil
}

func isPaymentIntent(paymentID string) bool {
	r, err := regexp.Compile("pi_(.+)")
	if err != nil {
		log.Println(err)
		return false
	}

	str := r.FindString(paymentID)
	if len(str) == 0 {
		return false
	}

	return true
}

func (db *paymentDatastore) getPaymentByToken(token string) (*Payment, error) {
	p := Payment{}

	err := db.Get(&p, "SELECT * FROM business_money_request_payment WHERE payment_token = $1", token)
	if err != nil && err != sql.ErrNoRows {
		log.Println(err)
		return nil, err
	}

	return &p, err
}

func (db *paymentDatastore) getPaymentByInvoiceId(invoiceId string) (*Payment, error) {
	p := Payment{}

	err := db.Get(&p, "SELECT * FROM business_money_request_payment WHERE invoice_id = $1", invoiceId)
	if err != nil && err != sql.ErrNoRows {
		log.Println(err)
		return nil, err
	}

	return &p, err
}

func (db *paymentDatastore) getPaymentByReceiptToken(token string) (*Payment, error) {
	p := Payment{}

	err := db.Get(&p, "SELECT * FROM business_money_request_payment WHERE receipt_token = $1", token)
	if err != nil && err != sql.ErrNoRows {
		log.Println(err)
		return nil, err
	}

	return &p, err
}

func (db *paymentDatastore) SendCardReaderReceipt(r CardReaderReceiptCreate) error {
	// check access
	err := auth.NewAuthService(*db.sourceReq).CheckBusinessAccess(r.BusinessID)
	if err != nil {
		return err
	}

	// validate input
	if len(r.RequestID) == 0 {
		return errors.New("Request Id is required")
	}

	if r.ReceiptMode != nil {
		switch *r.ReceiptMode {
		case CardReaderReceiptModeEmail:
			if r.CustomerContact == nil {
				return errors.New("Email address is required")
			}

			e, err := mail.ParseAddress(*r.CustomerContact)
			if err != nil {
				return errors.New("Invalid email address")
			}

			*r.CustomerContact = e.Address
		case CardReaderReceiptModeSMS:
			if r.CustomerContact == nil {
				return errors.New("Phone number is required")
			}

			// Validate phone no default
			ph, err := libphonenumber.Parse(*r.CustomerContact, "")
			if err != nil {
				return errors.New("Invalid phone number")
			}

			*r.CustomerContact = libphonenumber.Format(ph, libphonenumber.E164)
		default:
			return errors.New("Invalid receipt mode")
		}
	} else {
		mode := CardReaderReceiptModeNone
		r.ReceiptMode = &mode
	}

	// insert in payment table
	cardReaderReceipt := CardReaderReceiptCreate{
		BusinessID:      r.BusinessID,
		RequestID:       r.RequestID,
		ReceiptMode:     r.ReceiptMode,
		CustomerContact: r.CustomerContact,
	}
	pr, err := db.saveReceiptMode(cardReaderReceipt)
	if err != nil {
		log.Println(err)
		return nil
	}

	// Check if stripe has notified via webhook or not
	if pr.CardBrand != nil && pr.CardLast4 != nil {
		return db.sendCardReaderReceipts(pr)
	}

	return nil
}

func (db *paymentDatastore) sendCardReaderReceipts(p *Payment) error {

	req := RequestJoin{}

	err := db.Get(
		&req, `
		SELECT
			business_money_request.created_user_id, business_money_request.business_id, business_money_request.notes, business_money_request.request_type,
			business_invoice.invoice_number, business_invoice.id "business_invoice.id",
			business_money_request.amount, business_money_request.request_status, business_money_request.request_source,
			business_money_request_payment.id, business_money_request_payment.status, business.legal_name, business.phone "business.phone",
			business_money_request_payment.request_id, business_money_request_payment.source_payment_id,
			business_contact.first_name "business_contact.first_name", business_contact.last_name "business_contact.last_name", business_contact.email "business_contact.email",
			business_contact.id "business_contact.id", business_contact.business_name "business_contact.business_name", business.email "business.email", 
			consumer.first_name "consumer.first_name", consumer.last_name "consumer.last_name"  FROM business_money_request
			JOIN business_money_request_payment ON business_money_request.id = business_money_request_payment.request_id
			LEFT JOIN business_invoice ON business_money_request.id = business_invoice.request_id
			LEFT JOIN business_contact ON business_contact.id = business_money_request.contact_id
			JOIN wise_user ON business_money_request.created_user_id = wise_user.id
			JOIN business ON business.id = business_money_request.business_id
			JOIN consumer ON consumer.id = wise_user.consumer_id
			WHERE business_money_request_payment.request_id = $1`,
		p.RequestID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New(fmt.Sprintf("Payment with request id:%v not found", p.RequestID))
		}

		return err
	}

	req.BusinessName = shared.GetBusinessName(req.LegalName, req.DBA)

	req.ContactEmail = p.CustomerContact
	notes := "Card at " + req.BusinessName
	req.Notes = &notes

	// create receipt
	rNumber := shared.GetRandomString("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789", 8) + "-" +
		shared.GetRandomString("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789", 5)
	receiptGenerate := ReceiptGenerate{
		ContactFirstName: req.ContactFirstName,
		ContactLastName:  req.ContactLastName,
		ContactEmail:     req.ContactEmail,
		PaymentDate:      p.PaymentDate,
		PaymentBrand:     p.CardBrand,
		PaymentNumber:    p.CardLast4,
		ReceiptNumber:    rNumber,
		InvoiceNumber:    req.InvoiceNumber,
		BusinessName:     req.BusinessName,
		BusinessPhone:    req.BusinessPhone,
		Amount:           req.Amount,
		Notes:            req.Notes,
		UserID:           req.UserID,
		BusinessID:       req.BusinessID,
		ContactID:        req.ContactID,
		InvoiceID:        req.InvoiceID,
		RequestID:        req.RequestID,
		InvoiceIdV2:      req.InvoiceIdV2,
	}
	content, rID, err := db.GenerateReceipt(receiptGenerate)
	if err != nil {
		log.Println(err)
		return err
	}

	// update receipt id in pos receipt table
	err = db.updateReceiptID(p.ID, rID)
	if err != nil {
		log.Println(err)
		return err
	}

	if p.PaymentDate == nil {
		return errors.New("payment date is nil")
	}

	paidDate := (*p.PaymentDate).Format("Jan _2, 2006")

	// send receipt to business
	r := ReceiptRequest{
		ContactFirstName:    req.ContactFirstName,
		ContactLastName:     req.ContactLastName,
		ContactEmail:        req.ContactEmail,
		ContactBusinessName: req.ContactBusinessName,
		BusinessEmail:       req.Email,
		BusinessName:        req.BusinessName,
		Amount:              req.Amount,
		Notes:               req.Notes,
		PaymentDate:         paidDate,
		ReceiptNumber:       rNumber,
		Content:             content,
	}
	db.SendCardReaderReceiptToBusiness(r)

	if p.ReceiptMode == nil {
		return nil
	}

	switch *p.ReceiptMode {
	case CardReaderReceiptModeEmail:
		// send receipt to customer
		r := ReceiptRequest{
			ContactFirstName:    req.ContactFirstName,
			ContactLastName:     req.ContactLastName,
			ContactEmail:        req.ContactEmail,
			ContactBusinessName: req.ContactBusinessName,
			BusinessName:        req.BusinessName,
			Amount:              req.Amount,
			Notes:               req.Notes,
			PaymentDate:         paidDate,
			ReceiptNumber:       rNumber,
			Content:             content,
		}
		db.SendCardReaderReceiptToCustomer(r)
	case CardReaderReceiptModeSMS:
		// Send sms using twilio
		if len(os.Getenv("PAYMENTS_URL")) == 0 {
			log.Println("Missing receipt url")
			return errors.New("Missing receipt url")
		}

		url := os.Getenv("PAYMENTS_URL") + "/receipt?token=" + *p.ReceiptToken
		body := fmt.Sprintf(services.ReceiptSMS, req.BusinessName, url)
		r := t.SMSRequest{
			Body:  body,
			Phone: *p.CustomerContact,
		}

		err := t.NewTwilioService().SendSMS(r)
		if err != nil {
			log.Println("Error sending receipt sms ", err)
			return err
		}
	default:
		break
	}

	return nil
}

func (db *paymentDatastore) saveReceiptMode(receipt CardReaderReceiptCreate) (*Payment, error) {

	_, err := db.NamedExec(
		fmt.Sprintf(
			"UPDATE business_money_request_payment SET receipt_mode = :receipt_mode, customer_contact = :customer_contact WHERE request_id = '%s'",
			receipt.RequestID,
		), receipt,
	)

	if err != nil {
		return nil, err
	}

	p, err := db.GetPaymentByRequestID(receipt.RequestID)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return p, nil
}

func (db *paymentDatastore) updateReceiptID(ID string, receiptID string) error {
	sql := `UPDATE business_money_request_payment SET receipt_id = $1 WHERE id = $2`
	_, err := db.Exec(sql, receiptID, ID)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (db *paymentDatastore) GetPaymentReceiptFromInvoice(invoiceID id.InvoiceID) (*InvoicePaymentReceipt, error) {
	paymentReceipt := &InvoicePaymentReceipt{}
	query := `SELECT card_brand , card_number, payment_date, wallet_type, receipt_number, bmrp.status status
				FROM business_money_request_payment bmrp left join business_receipt br  
				on bmrp.invoice_id  = br.invoice_id_v2 
				where bmrp.invoice_id  = $1;`
	err := db.Get(paymentReceipt, query, invoiceID)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return paymentReceipt, nil
}

func (db *paymentDatastore) GetReceiptInfo(receiptToken string) (*CardReaderReceiptResponse, error) {
	// Fetch payment details
	p, err := db.getPaymentByReceiptToken(receiptToken)
	if err != nil {
		log.Println("error finding receipt token", err)
		return nil, errors.New("Unable to display receipt")
	}

	cardReaderResponse := CardReaderReceiptResponse{}

	err = db.Get(
		&cardReaderResponse,
		`SELECT business_money_request.amount, business.legal_name, business_money_request_payment.card_brand, business_money_request_payment.card_number, business_money_request.business_id, 
		business_money_request_payment.receipt_id, business_money_request_payment.purchase_address, business_money_request_payment.payment_date, business.dba FROM business_money_request_payment
		JOIN business_money_request ON business_money_request_payment.request_id = business_money_request.id
		JOIN business ON business_money_request.business_id = business.id WHERE business_money_request_payment.receipt_token = $1`, p.ReceiptToken)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	println("join values are ", cardReaderResponse.Amount, cardReaderResponse.BusinessName)

	clientSecret, err := stripe.NewStripeService(nil).GetClientSecret(*p.SourcePaymentID)
	if err != nil {
		return nil, err
	}

	cardReaderResponse.ClientSecret = *clientSecret
	cardReaderResponse.StripeKey = os.Getenv("STRIPE_PUBLISH_KEY")
	cardReaderResponse.BusinessName = shared.GetBusinessName(cardReaderResponse.LegalName, cardReaderResponse.DBA)

	// Send back client secret
	return &cardReaderResponse, nil
}

func (db *paymentDatastore) GetPayments(requestIDs []shared.PaymentRequestID) ([]Payment, error) {
	query := `SELECT *  FROM business_money_request_payment  WHERE invoice_id in (?) `
	payments := []Payment{}
	q, args, err := sqlx.In(query, requestIDs)
	if err != nil {
		return nil, err
	}

	q = db.DB.Rebind(q)
	err = db.DB.Select(&payments, q, args...)
	if err != nil {
		return nil, err
	}
	return payments, nil
}

func GetInvoiceStatus(status grpcSvc.InvoiceRequestStatus) PaymentRequestStatus {
	switch status {
	case grpcSvc.InvoiceRequestStatus_IRT_OPEN:
		return PaymentRequestStatusPending
	case grpcSvc.InvoiceRequestStatus_IRT_PAID:
		return PaymentRequestStatusComplete
	case grpcSvc.InvoiceRequestStatus_IRT_CANCELLED:
		return PaymentRequestStatusCanceled
	case grpcSvc.InvoiceRequestStatus_IRT_PROCESSING_PAYMENT:
		return PaymentRequestStatusInProcess
	default:
		return PaymentRequestStatusPending // TODO : have unspecified here
	}
}

func GetInvoiceGrpcStatus(status PaymentRequestStatus) grpcSvc.InvoiceRequestStatus {
	switch status {
	case PaymentRequestStatusPending:
		return grpcSvc.InvoiceRequestStatus_IRT_OPEN
	case PaymentRequestStatusComplete:
		return grpcSvc.InvoiceRequestStatus_IRT_PAID
	case PaymentRequestStatusCanceled:
		return grpcSvc.InvoiceRequestStatus_IRT_CANCELLED
	case PaymentRequestStatusInProcess:
		return grpcSvc.InvoiceRequestStatus_IRT_PROCESSING_PAYMENT
	default:
		return grpcSvc.InvoiceRequestStatus_IRT_UNSPECIFIED
	}
}
