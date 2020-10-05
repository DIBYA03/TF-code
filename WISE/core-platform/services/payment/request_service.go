/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

package payment

import (
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/wiseco/core-platform/services/invoice"

	"github.com/jmoiron/sqlx"
	"github.com/microcosm-cc/bluemonday"
	"github.com/wiseco/core-platform/auth"
	"github.com/wiseco/core-platform/partner/service"
	"github.com/wiseco/core-platform/partner/service/sendgrid"
	"github.com/wiseco/core-platform/partner/service/stripe"
	"github.com/wiseco/core-platform/services"
	"github.com/wiseco/core-platform/services/banking/business"
	b "github.com/wiseco/core-platform/services/business"
	"github.com/wiseco/core-platform/services/contact"
	"github.com/wiseco/core-platform/services/data"
	"github.com/wiseco/core-platform/services/document"
	"github.com/wiseco/core-platform/services/pdf"
	"github.com/wiseco/core-platform/shared"
	"github.com/wiseco/go-lib/id"
	"mvdan.cc/xurls/v2"

	t "github.com/wiseco/core-platform/partner/service/twilio"
)

type requestDatastore struct {
	sourceReq services.SourceRequest
	*sqlx.DB
}

type RequestService interface {
	// Read
	GetByID(shared.PaymentRequestID, shared.BusinessID) (*Request, error)
	GetByIDInternal(shared.PaymentRequestID) (*Request, error)
	GetByRequestSourceIDInternal(id.ShopifyOrderID, shared.BusinessID) (*Request, error)
	ListByContactID(shared.BusinessID, string, int, int, string) ([]RequestPayment, error)
	List(shared.BusinessID, int, int, string) ([]RequestPayment, error)

	GetConnectionToken(PaymentConnectionRequest) (*PaymentConnectionResponse, error)

	CreatePaymentIntent(PaymentResponse) (*stripe.PaymentResponse, error)

	Request(*RequestInitiate) (*Request, error)
	UpdateRequestStatus(*RequestStatusUpdate) (*Request, error)
	UpdateRequestStatusByIntentID(string, PaymentRequestStatus) error
	UpdateInvoiceStatus(id.InvoiceID, PaymentRequestStatus) error

	CapturePayment(PaymentCaptureRequest) error

	Resend([]PaymentRequestResend, shared.BusinessID) error
}

func NewRequestService(r services.SourceRequest) RequestService {
	return &requestDatastore{r, data.DBWrite}
}

func (db *requestDatastore) GetByID(id shared.PaymentRequestID, businessID shared.BusinessID) (*Request, error) {
	// Check access
	err := auth.NewAuthService(db.sourceReq).CheckBusinessAccess(businessID)
	if err != nil {
		return nil, err
	}
	if os.Getenv("USE_INVOICE_SERVICE") == "true" {
		return newLegacyInvoiceService().GetInvoice(id)
	} else {
		a := Request{}

		err = db.Get(&a, "SELECT * FROM business_money_request WHERE id = $1 AND business_id = $2", id, businessID)
		if err != nil && err != sql.ErrNoRows {

			return nil, err
		}

		return &a, err
	}
}

func (db *requestDatastore) GetByIDInternal(id shared.PaymentRequestID) (*Request, error) {
	a := Request{}

	err := db.Get(&a, "SELECT * FROM business_money_request WHERE id = $1", id)
	if err != nil && err != sql.ErrNoRows {

		return nil, err
	}

	return &a, err
}

func (db *requestDatastore) GetByRequestSourceIDInternal(sourceID id.ShopifyOrderID, bID shared.BusinessID) (*Request, error) {
	a := Request{}

	err := db.Get(&a, "SELECT * FROM business_money_request WHERE request_source_id = $1 AND business_id = $2", sourceID, bID)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return &a, err
}

func (db *requestDatastore) List(businessID shared.BusinessID, offset int, limit int, status string) ([]RequestPayment, error) {
	// Check access
	err := auth.NewAuthService(db.sourceReq).CheckBusinessAccess(businessID)
	if err != nil {
		return nil, err
	}

	rows := []RequestPayment{}

	statusCondition := ""
	if len(status) > 0 {
		statusCondition = "AND request_status = '" + status + "' "
	}

	if os.Getenv("USE_BANKING_SERVICE") == "true" {
		if os.Getenv("USE_INVOICE_SERVICE") == "true" {
			invSvc := newLegacyInvoiceService()
			payments, err := invSvc.GetManyInvoices(businessID, limit, offset, status, "")
			if err != nil {
				return nil, err
			}
			// map the payment date from payment

			// get list of invoice ids
			invoiceIds := []string{}
			for _, paymentModel := range payments {
				invoiceId := strings.Replace(paymentModel.ID.ToPrefixString(), shared.PaymentRequestPrefix, "", 1)
				invoiceIds = append(invoiceIds, fmt.Sprintf("'%s'", invoiceId))
			}
			// get payment info
			query := fmt.Sprintf(`SELECT invoice_id "business_money_request.id",  business_money_request_payment.payment_date from business_money_request_payment WHERE business_money_request_payment.invoice_id in (%s)`, strings.Join(invoiceIds, ","))
			log.Println(query)
			paymentInfos := []RequestPayment{}
			err = db.Select(&paymentInfos, query)
			if err != nil && err != sql.ErrNoRows {
				return nil, err
			}

			// make map of request_id:payment_date, as we might not have an entry in payment for invoice, but still invoice will be part of response
			mapRequestID := map[shared.PaymentRequestID]*time.Time{}
			for _, row := range paymentInfos {
				mapRequestID[row.Request.ID] = row.PaymentDate
			}

			// from the payment info
			for _, paymentInfo := range payments {
				row := RequestPayment{
					Request: paymentInfo,
				}
				if value, ok := mapRequestID[paymentInfo.ID]; ok {
					row.PaymentDate = value
				}
				rows = append(rows, row)
			}
		} else {
			query := `SELECT business_money_request.id "business_money_request.id",
    business_money_request.business_id "business_money_request.business_id",
    business_money_request.created_user_id "business_money_request.created_user_id",
    business_money_request.contact_id "business_money_request.contact_id",
    business_money_request.amount "business_money_request.amount",
    business_money_request.currency "business_money_request.currency",
    business_money_request.notes "business_money_request.notes",
    business_money_request.request_status "business_money_request.request_status",
    business_money_request.request_type "business_money_request.request_type",
    business_money_request.pos_id "business_money_request.pos_id",
    business_money_request.created "business_money_request.created",
    business_money_request.modified "business_money_request.modified",
    business_money_request_payment.payment_date
    FROM business_money_request
    LEFT JOIN business_money_request_payment ON business_money_request.id = business_money_request_payment.request_id
    WHERE business_money_request.business_id = $1 ` + statusCondition +
				`ORDER BY business_money_request.modified DESC LIMIT $2 OFFSET $3`

			err = db.Select(&rows, query, businessID, limit, offset)
			if err != nil && err != sql.ErrNoRows {
				return nil, err
			}
		}
		for _, r := range rows {
			bts, err := business.NewBankingTransferService()
			if err != nil {
				return nil, err
			}

			mt, err := bts.GetByPaymentRequestID(r.Request.ID)
			if err != nil {
				if err == sql.ErrNoRows {
					continue
				}

				return nil, err
			}

			if mt.PostedCreditTransactionID != nil {
				prID := shared.PostedTransactionID(*mt.PostedCreditTransactionID)

				r.TransactionID = &prID
			}
		}
	} else {
		query := `SELECT business_money_request.id "business_money_request.id",
    business_money_request.business_id "business_money_request.business_id",
    business_money_request.created_user_id "business_money_request.created_user_id",
    business_money_request.contact_id "business_money_request.contact_id",
    business_money_request.amount "business_money_request.amount",
    business_money_request.currency "business_money_request.currency",
    business_money_request.notes "business_money_request.notes",
    business_money_request.request_status "business_money_request.request_status",
    business_money_request.request_type "business_money_request.request_type",
    business_money_request.pos_id "business_money_request.pos_id",
    business_money_request.created "business_money_request.created",
    business_money_request.modified "business_money_request.modified",
    business_money_request_payment.payment_date,
	business_money_transfer.posted_credit_transaction_id
    FROM business_money_request
    LEFT JOIN business_money_request_payment ON business_money_request.id = business_money_request_payment.request_id
    LEFT JOIN business_money_transfer ON business_money_request.id = business_money_transfer.money_request_id
    WHERE business_money_request.business_id = $1 ` + statusCondition +
			`ORDER BY business_money_request.modified DESC LIMIT $2 OFFSET $3`

		err = db.Select(&rows, query, businessID, limit, offset)
		if err != nil && err != sql.ErrNoRows {
			return nil, err
		}
	}

	return rows, err
}

func (db *requestDatastore) ListByContactID(businessID shared.BusinessID, contactID string, offset int, limit int, status string) ([]RequestPayment, error) {
	// Check access
	err := auth.NewAuthService(db.sourceReq).CheckBusinessAccess(businessID)
	if err != nil {
		return nil, err
	}

	rows := []RequestPayment{}

	statusCondition := ""
	if len(status) > 0 {
		statusCondition = "AND request_status = '" + status + "' "
	}

	if os.Getenv("USE_BANKING_SERVICE") == "true" {
		if os.Getenv("USE_INVOICE_SERVICE") == "true" {
			invSvc := newLegacyInvoiceService()
			payments, err := invSvc.GetManyInvoices(businessID, limit, offset, status, contactID)
			if err != nil {
				return nil, err
			}
			// map the payment date from payment

			// get list of invoice ids
			invoiceIds := []string{}
			for _, paymentModel := range payments {
				invoiceId := strings.Replace(paymentModel.ID.ToPrefixString(), shared.PaymentRequestPrefix, "", 1)
				invoiceIds = append(invoiceIds, fmt.Sprintf("'%s'", invoiceId))
			}
			// get payment info
			query := fmt.Sprintf(`SELECT invoice_id "business_money_request.id",  business_money_request_payment.payment_date from business_money_request_payment
						WHERE invoice_id in (%s)`, strings.Join(invoiceIds, ","))
			log.Println(query)
			paymentInfos := []RequestPayment{}
			err = db.Select(&paymentInfos, query)
			if err != nil && err != sql.ErrNoRows {
				return nil, err
			}

			// make map of request_id:payment_date, as we might not have an entry in payment for invoice, but still invoice will be part of response
			mapRequestID := map[shared.PaymentRequestID]*time.Time{}
			for _, row := range paymentInfos {
				mapRequestID[row.Request.ID] = row.PaymentDate
			}

			// from the payment info
			for _, paymentInfo := range payments {
				row := RequestPayment{
					Request: paymentInfo,
				}
				if value, ok := mapRequestID[paymentInfo.ID]; ok {
					row.PaymentDate = value
				}
				rows = append(rows, row)
			}
		} else {
			query := `SELECT business_money_request.id "business_money_request.id", 
	business_money_request.business_id "business_money_request.business_id", 
	business_money_request.created_user_id "business_money_request.created_user_id",
	business_money_request.contact_id "business_money_request.contact_id", 
	business_money_request.amount "business_money_request.amount", 
	business_money_request.currency "business_money_request.currency",
	business_money_request.notes "business_money_request.notes", 
	business_money_request.request_status "business_money_request.request_status", 
	business_money_request.request_type "business_money_request.request_type",
	business_money_request.pos_id "business_money_request.pos_id", 
	business_money_request.created "business_money_request.created",
	business_money_request.modified "business_money_request.modified",
	business_money_request_payment.payment_date
	FROM business_money_request
	LEFT JOIN business_money_request_payment ON business_money_request.id = business_money_request_payment.request_id
	WHERE business_money_request.contact_id = $1 AND business_money_request.business_id = $2 ` + statusCondition +
				`ORDER BY business_money_request.modified DESC LIMIT $3 OFFSET $4`

			err = db.Select(&rows, query, contactID, businessID, limit, offset)
			if err != nil && err != sql.ErrNoRows {
				return nil, err
			}
		}
		for _, r := range rows {
			bts, err := business.NewBankingTransferService()
			if err != nil {
				return nil, err
			}

			mt, err := bts.GetByPaymentRequestID(r.Request.ID)
			if err != nil {
				if err == sql.ErrNoRows {
					continue
				}

				return nil, err
			}

			if mt.PostedCreditTransactionID != nil {
				prID := shared.PostedTransactionID(*mt.PostedCreditTransactionID)

				r.TransactionID = &prID
			}
		}
	} else {
		query := `SELECT business_money_request.id "business_money_request.id", 
	business_money_request.business_id "business_money_request.business_id", 
	business_money_request.created_user_id "business_money_request.created_user_id",
	business_money_request.contact_id "business_money_request.contact_id", 
	business_money_request.amount "business_money_request.amount", 
	business_money_request.currency "business_money_request.currency",
	business_money_request.notes "business_money_request.notes", 
	business_money_request.request_status "business_money_request.request_status", 
	business_money_request.request_type "business_money_request.request_type",
	business_money_request.pos_id "business_money_request.pos_id", 
	business_money_request.created "business_money_request.created",
	business_money_request.modified "business_money_request.modified",
	business_money_request_payment.payment_date,
	business_money_transfer.posted_credit_transaction_id
	FROM business_money_request
	LEFT JOIN business_money_request_payment ON business_money_request.id = business_money_request_payment.request_id
	LEFT JOIN business_money_transfer ON business_money_request.id = business_money_transfer.money_request_id
	WHERE business_money_request.contact_id = $1 AND business_money_request.business_id = $2 ` + statusCondition +
			`ORDER BY business_money_request.modified DESC LIMIT $3 OFFSET $4`

		err = db.Select(&rows, query, contactID, businessID, limit, offset)
		if err != nil && err != sql.ErrNoRows {
			return nil, err
		}
	}

	return rows, err
}

func (db *requestDatastore) Request(request *RequestInitiate) (*Request, error) {
	// Check access
	if request.CardReaderID != nil {
		err := auth.NewAuthService(db.sourceReq).CheckBusinessCardReaderAccess(string(*request.CardReaderID))
		if err != nil {
			return nil, err
		}
	} else {
		err := auth.NewAuthService(db.sourceReq).CheckBusinessAccess(request.BusinessID)
		if err != nil {
			return nil, err
		}
	}

	maxCardReaderAmount, err := strconv.ParseFloat(os.Getenv("CARD_READER_MAX_MONEY_REQUEST_ALLOWED"), 64)
	if err != nil {
		return nil, err
	}

	maxCardOnlineAmount, err := strconv.ParseFloat(os.Getenv("CARD_ONLINE_MAX_MONEY_REQUEST_ALLOWED"), 64)
	if err != nil {
		return nil, err
	}

	if request.Amount == 0 {
		return nil, errors.New("Amount cannot be zero")
	}

	if string(request.Currency) == "" {
		return nil, errors.New("Currency is required")
	}

	if len(request.RequestType) == 0 {
		return nil, errors.New("Request type is required")
	}

	// Get business details
	b, err := b.NewBusinessService(db.sourceReq).GetById(request.BusinessID)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	} else if err != nil && err == sql.ErrNoRows {
		return nil, errors.New("business not found")
	}

	var p *stripe.PaymentResponse
	switch request.RequestType {
	case PaymentRequestTypePOS:
		if request.Amount > maxCardReaderAmount {
			e := fmt.Sprintf("Amount cannot exceed $%s", strconv.FormatFloat(maxCardReaderAmount, 'f', 2, 64))
			return nil, errors.New(e)
		}

		if request.CardReaderID == nil {
			return nil, errors.New("POS ID is required")
		}

		metadata, err := db.getPaymentIntentMetadata(nil, b.ID, b.Name())
		if err != nil {
			return nil, err
		}

		metadata.IPAddress = request.IPAddress

		r := db.sourceReq.PartnerServiceRequest()

		// Create payment object using stripe
		p, err = stripe.NewStripeService(&r).CreatePayment(transformForPaymentService(b.Name(), request.RequestType, request.Amount, request.Currency, metadata))
		if err != nil {
			return nil, err
		}
	case PaymentRequestTypeInvoiceCard:
		if request.Amount > maxCardOnlineAmount {
			e := fmt.Sprintf("Amount cannot exceed $%s", strconv.FormatFloat(maxCardOnlineAmount, 'f', 2, 64))
			return nil, errors.New(e)
		}

		fallthrough
	case PaymentRequestTypeInvoiceBank:
		fallthrough
	case PaymentRequestTypeInvoiceNone:
		fallthrough
	case PaymentRequestTypeInvoiceCardAndBank:
		if request.Amount > maxCardOnlineAmount {
			request.RequestType = PaymentRequestTypeInvoiceNone
		}

		if request.Notes == nil {
			return nil, errors.New("Notes is required")
		} else if len(strings.TrimSpace(*request.Notes)) > 80 {
			return nil, errors.New("Notes cannot exceed more than 80 characters")
		}

		rxRelaxed := xurls.Relaxed()
		if len(rxRelaxed.FindString(*request.Notes)) > 0 {
			return nil, errors.New("Notes cannot contain urls")
		}

		if request.ContactID == nil {
			return nil, errors.New("Contact id is required")
		}

		// sanitize request notes
		policy := bluemonday.StrictPolicy()
		notes := policy.Sanitize(
			*request.Notes,
		)
		*request.Notes = notes
		break
	default:
		return nil, errors.New("invalid payment request type")
	}

	//These are disabled per https://wise.atlassian.net/browse/P2-450
	if request.RequestType == PaymentRequestTypeInvoiceBank {
		request.RequestType = PaymentRequestTypeInvoiceNone
	}
	moneyRequest := transformFromEmailService(request)
	// Get account details
	var a *business.BankAccount
	accounts, err := business.NewBankAccountService(db.sourceReq).List(request.BusinessID, 10, 0)
	if err != nil {
		return nil, errors.New("Unable to fetch bank account details of the business")
	}

	if len(*accounts) > 0 {
		a = &(*accounts)[0]
	} else {
		return nil, errors.New("Unable to fetch bank account details of the business")
	}
	if os.Getenv("USE_INVOICE_SERVICE") == "true" && request.RequestType != PaymentRequestTypePOS {
		m := &Request{}
		invSvc := newLegacyInvoiceService()
		moneyRequest.ContactID = request.ContactID
		_, err := invSvc.CreateInvoice(moneyRequest)
		if err != nil {
			log.Println(err)
			return m, fmt.Errorf("Unable to request")
		}
		pendingStatus := PaymentRequestStatusPending
		m.RequestType = moneyRequest.RequestType
		m.Amount = moneyRequest.Amount
		m.BusinessID = moneyRequest.BusinessID
		m.ContactID = moneyRequest.ContactID
		m.CreatedUserID = moneyRequest.CreatedUserID
		m.Currency = moneyRequest.Currency
		m.Notes = moneyRequest.Notes
		m.Status = &pendingStatus
		return m, nil
	} else {
		// Default/mandatory fields
		columns := []string{
			"created_user_id", "business_id", "amount", "currency", "notes", "message_id", "request_type", "pos_id", "request_ip_address", "request_source", "request_source_id",
		}
		// Default/mandatory values
		values := []string{
			":created_user_id", ":business_id", ":amount", ":currency", ":notes", ":message_id", ":request_type", ":pos_id", ":request_ip_address", ":request_source", ":request_source_id",
		}

		// Get contact details
		var c *contact.Contact
		if request.ContactID != nil {
			c, err = contact.NewContactService(db.sourceReq).GetById(*request.ContactID, request.BusinessID)
			if err != nil && err != sql.ErrNoRows {
				return nil, err
			} else if err != nil && err == sql.ErrNoRows {
				return nil, errors.New("contact not found")
			}

			columns = append(columns, "contact_id")
			values = append(values, ":contact_id")

			moneyRequest.ContactID = request.ContactID
		}
		sql := fmt.Sprintf("INSERT INTO business_money_request(%s) VALUES(%s) RETURNING *", strings.Join(columns, ", "), strings.Join(values, ", "))

		stmt, err := db.PrepareNamed(sql)
		if err != nil {
			return nil, err
		}

		m := &Request{}

		err = stmt.Get(m, &moneyRequest)
		if err != nil {
			return nil, err
		}

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

		requestPayment := &Payment{}

		err = stmt.Get(requestPayment, &payment)
		if err != nil {
			return nil, err
		}

		switch request.RequestType {
		case PaymentRequestTypeInvoiceNone:
			fallthrough
		case PaymentRequestTypeInvoiceCard:
			fallthrough
		case PaymentRequestTypeInvoiceBank:
			fallthrough
		case PaymentRequestTypeInvoiceCardAndBank:
			// Send request email
			responseID, err := db.sendInvoice(b, c, m, a, *requestPayment.PaymentToken)
			if err != nil {
				return nil, err
			}

			// Update payment request with send grid message id
			if responseID != nil {
				u := RequestIDUpdate{
					BusinessID: m.BusinessID,
					ContactID:  m.ContactID,
					ID:         m.ID,
					MessageID:  *responseID,
				}
				requestUpdate, err := db.UpdateMessageID(&u)
				if err != nil {
					return nil, err
				}

				return requestUpdate, nil
			}

			return m, nil
		case PaymentRequestTypePOS:
			// send back payment intent secret
			m.PaymentIntentToken = payment.ClientSecret

			return m, nil
		default:
			return nil, fmt.Errorf("Unable to request")
		}
	}
}

func transformForPaymentService(businessName string, requestType PaymentRequestType, amount float64, currency Currency, metadata *RequestMetadata) stripe.PaymentRequest {
	descriptor := businessName
	if len(descriptor) > 22 {
		descriptor = descriptor[:22]
	}

	// strip off characters like <>\'"*
	reg := regexp.MustCompile("[<>\\'\"*]+")
	descriptor = reg.ReplaceAllString(descriptor, "")

	var paymentMethod string
	var captureMethod string
	switch requestType {
	case PaymentRequestTypePOS:
		paymentMethod = "card_present"
		captureMethod = "manual"
	default:
		paymentMethod = "card"
		captureMethod = ""
	}

	pm := stripe.PaymentMetadata{
		BusinessName:     metadata.BusinessName,
		AvailableBalance: metadata.AvailableBalance,
		PaymentMethod:    paymentMethod,
	}

	if metadata.IPAddress != nil {
		pm.IPAddress = *metadata.IPAddress
	}

	ownerName := metadata.FirstName
	if metadata.MiddleName != nil {
		ownerName = ownerName + " " + *metadata.MiddleName
	}
	ownerName = ownerName + " " + metadata.LastName

	pm.BusinessOwnerName = ownerName
	pm.BusinessName = metadata.BusinessName

	return stripe.PaymentRequest{
		PaymentMethod:   service.PaymentMethod(paymentMethod),
		Amount:          amount,
		Currency:        service.Currency(currency),
		Descriptor:      descriptor,
		ReceiptEmail:    os.Getenv("WISE_INVOICE_EMAIL"),
		CaptureMethod:   captureMethod,
		PaymentMetadata: pm,
	}
}

func transformFromEmailService(request *RequestInitiate) Request {

	r := Request{
		BusinessID: request.BusinessID,
	}

	r.CreatedUserID = request.CreatedUserID
	r.Currency = request.Currency
	r.Amount = request.Amount
	r.Notes = request.Notes
	r.MessageID = ""
	r.RequestType = &request.RequestType
	r.CardReaderID = request.CardReaderID
	r.IPAddress = request.IPAddress
	r.RequestSource = request.RequestSource
	r.RequestSourceID = request.RequestSourceID

	return r
}

func transformFromPaymentService(stripe *stripe.PaymentResponse, requestId shared.PaymentRequestID) PartnerPaymentResponse {

	p := PartnerPaymentResponse{
		RequestID: requestId,
	}

	if stripe != nil {
		p.SourcePaymentID = stripe.IntentID
		p.ClientSecret = stripe.ClientSecret
		p.Status = string(stripe.Status)
	} else {
		p.Status = string(PaymentStatusPending)
	}

	// Set token expiry to 30 days(720 hours)
	expTime := time.Now().UTC().Add(time.Hour * time.Duration(720))
	p.ExpirationDate = &expTime

	return p
}

func (db *requestDatastore) getNextInvoiceSequence(bID shared.BusinessID, employerNumber string) (*string, error) {
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

func (db *requestDatastore) UpdateRequestStatus(u *RequestStatusUpdate) (*Request, error) {
	if os.Getenv("USE_INVOICE_SERVICE") == "true" {
		err := newLegacyInvoiceService().UpdateInvoiceStatus(u.ID, u.Status)
		if err != nil {
			return nil, err
		}
	} else {
		r, err := db.GetByID(u.ID, u.BusinessID)
		if err != nil {
			return nil, err
		}

		var paymentStatus PaymentStatus
		switch u.Status {
		case PaymentRequestStatusCanceled:
			if r.Status != nil && *r.Status == PaymentRequestStatusInProcess {
				return nil, errors.New("invoices with payment in process cannot be canceled")
			}
			paymentStatus = PaymentStatusCanceled
		case PaymentRequestStatusComplete:
			paymentStatus = PaymentStatusSucceeded
			break
		default:
			return nil, errors.New("invalid request status")
		}

		_, err = db.NamedExec(
			fmt.Sprintf(
				"UPDATE business_money_request SET request_status = :request_status WHERE id = '%s'",
				u.ID,
			), u,
		)
		if err != nil {
			return nil, err
		}

		payment, err := NewPaymentService(db.sourceReq).GetPaymentByRequestID(u.ID)
		if err != nil {
			return nil, err
		}

		// Update payments also
		paymentUpdate := Payment{
			PaymentToken: nil,
			ID:           payment.ID,
			RequestID:    payment.RequestID,
			Status:       paymentStatus,
		}
		err = NewPaymentService(db.sourceReq).UpdatePaymentStatus(&paymentUpdate)
		if err != nil {
			return nil, err
		}
	}
	return db.GetByID(u.ID, u.BusinessID)
}

func (db *requestDatastore) UpdateRequestStatusByIntentID(intentID string, status PaymentRequestStatus) error {
	// check if we have invoice_id for business_money_request_payment, then make grpc call to update the status
	payment := Payment{}
	err := db.Get(&payment, "select * from business_money_request_payment where source_payment_id = $1", intentID)
	if err != nil {
		log.Println(err)
		return err
	}
	if payment.InvoiceID != nil && !payment.InvoiceID.IsZero() {
		// if the status is pending, then delete the entry from payment info, it will make it as initial state
		if status == PaymentRequestStatusPending {
			_, err := db.Exec(`DELETE FROM business_money_request_payment WHERE invoice_id=$1`, *payment.InvoiceID)
			if err != nil {
				log.Println(err)
				return err
			}
		}
		return db.UpdateInvoiceStatus(*payment.InvoiceID, status)
	} else { // else keep it as it is
		r := Request{}

		err := db.Get(&r, `SELECT business_money_request.* FROM business_money_request 
	JOIN business_money_request_payment 
	ON business_money_request.id = business_money_request_payment.request_id
	WHERE business_money_request_payment.source_payment_id = $1`, intentID)
		if err != nil {
			log.Println(err)
			return err
		}

		_, err = db.Exec(`UPDATE business_money_request SET request_status = $1 WHERE id = $2`, status, r.ID)
		return err
	}
}

func (db *requestDatastore) UpdateMessageID(u *RequestIDUpdate) (*Request, error) {
	_, err := db.NamedExec(
		fmt.Sprintf(
			"UPDATE business_money_request SET message_id = :message_id WHERE id = '%s'",
			u.ID,
		), u,
	)

	if err != nil {
		return nil, err
	}

	m, err := db.GetByID(u.ID, u.BusinessID)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (db *requestDatastore) saveInvoice(invoice InvoiceCreate) (*Invoice, error) {
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

	i := &Invoice{}

	err = stmt.Get(i, &invoice)
	if err != nil {
		return nil, err
	}

	return i, nil
}

func (db *requestDatastore) sendInvoice(b *b.Business, c *contact.Contact, m *Request, a *business.BankAccount, token string) (*string, error) {
	// Generate invoice ID
	invoiceNo, err := db.getNextInvoiceSequence(b.ID, b.EmployerNumber)
	if err != nil {
		return nil, err
	}

	issueDate := m.Created.Format("Jan _2, 2006")

	// Generate PDF
	invoice := pdf.Invoice{
		BusinessName:  b.Name(),
		BusinessPhone: *b.Phone,
		Amount:        shared.FormatFloatAmount(m.Amount),
		ContactEmail:  c.Email,
		ContactName:   c.Name(),
		InvoiceNo:     *invoiceNo,
		Notes:         *m.Notes,
		WisePhone:     os.Getenv("WISE_SUPPORT_PHONE"),
		IssueDate:     issueDate,
		AccountNumber: a.AccountNumber,
		RoutingNumber: a.RoutingNumber,
	}

	content, err := pdf.NewInvoiceService(invoice).GenerateInvoice()
	if err != nil {
		return nil, err
	}

	// Upload document to aws
	store, err := document.NewAWSS3DocStorageFromContent(string(b.ID), document.BusinessPrefix, "application/pdf", content)
	if err != nil {
		return nil, err
	}

	key, err := store.Key()
	if err != nil {
		return nil, err
	}

	// Store in database
	invoiceCreate := InvoiceCreate{
		RequestID:     &m.ID,
		BusinessID:    m.BusinessID,
		CreatedUserID: m.CreatedUserID,
		InvoiceNumber: *invoiceNo,
		ContactID:     c.ID,
		StorageKey:    *key,
	}

	_, err = db.saveInvoice(invoiceCreate)
	if err != nil {
		return nil, err
	}

	// In case of POS request type, invoice should not be sent to contact and business
	if *m.RequestType == PaymentRequestTypePOS {
		return nil, err

	} else {
		// Send invoice to contact
		responseID, err := sendInvoiceToCustomer(b, c, m, invoiceNo, token, issueDate, content, false)
		if err != nil {
			return nil, err
		}

		// send invoice to business
		sendInvoiceToBusiness(b, c, m, invoiceNo, token, issueDate, content)
		if err != nil {
			return nil, err
		}

		return responseID, nil
	}
}

func sendInvoiceSMSToCustomer(request t.SMSRequest) error {

	err := t.NewTwilioService().SendSMS(request)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func sendInvoiceEmailToCustomer(request sendgrid.EmailAttachmentRequest) (*string, error) {
	response, err := sendgrid.NewSendGridServiceWithout().SendAttachmentEmail(request)
	if err != nil {
		return nil, err
	}

	return &response.MessageId, nil
}

func sendInvoiceToCustomer(b *b.Business, c *contact.Contact, transfer *Request, invoiceNumber *string,
	token string, issueDate string, attachment *string, resend bool) (*string, error) {

	amount := shared.FormatFloatAmount(transfer.Amount)

	url := os.Getenv("PAYMENTS_URL") + "/request?token=" + token

	*transfer.Notes = strings.Replace(*transfer.Notes, "\n", "<br>", -1)

	var subject string
	if !resend {
		subject = fmt.Sprintf(services.CustomerInvoiceSubject, b.Name())
	} else {
		subject = fmt.Sprintf(services.CustomerResendInvoiceSubject, b.Name())
	}

	a := sendgrid.EmailAttachment{
		Attachment:  *attachment,
		ContentType: "application/pdf",
		FileName:    "Invoice - " + *invoiceNumber + ".pdf",
		ContentID:   "Invoice",
	}

	if c.Email == "" && c.PhoneNumber == "" {
		return nil, errors.New("Both phone and email are empty")
	}

	var body string
	var respID *string
	var err error
	if c.Email != "" {
		if *transfer.RequestType == PaymentRequestTypeInvoiceNone {
			body = fmt.Sprintf(services.CustomerInvoiceEmailWithoutPayURL, c.Name(), b.Name(), issueDate, *transfer.Notes, amount, amount, b.Name(), b.Name())
		} else {
			body = fmt.Sprintf(services.CustomerInvoiceEmail, c.Name(), b.Name(), issueDate, *transfer.Notes, amount, url, b.Name(), amount, amount, b.Name(), b.Name())
		}

		emailRequest := sendgrid.EmailAttachmentRequest{
			SenderEmail:   os.Getenv("WISE_INVOICE_EMAIL"),
			SenderName:    os.Getenv("WISE_SUPPORT_NAME"),
			ReceiverEmail: c.Email,
			ReceiverName:  c.Name(),
			Subject:       subject,
			Body:          body,
			Attachment:    []sendgrid.EmailAttachment{a},
		}

		respID, err = sendInvoiceEmailToCustomer(emailRequest)
		if err != nil {
			return nil, err
		}
	}

	// Send SMS only for Shopify requests
	if c.PhoneNumber != "" && transfer.isShopifyRequest() {
		if *transfer.RequestType == PaymentRequestTypeInvoiceNone {
			body = fmt.Sprintf(services.CustomerInvoiceSMSWithoutPayURL, c.Name(), b.Name(), issueDate, *transfer.Notes, amount, amount, b.Name(), b.Name())
		} else {
			body = fmt.Sprintf(services.CustomerInvoiceSMS, c.Name(), b.Name(), issueDate, *transfer.Notes, amount, url, amount, b.Name(), b.Name())
		}

		smsReq := t.SMSRequest{
			Body:  body,
			Phone: c.PhoneNumber,
		}

		err = sendInvoiceSMSToCustomer(smsReq)
		if err != nil {
			return nil, err
		}
	}

	return respID, err
}

func sendInvoiceToBusiness(b *b.Business, c *contact.Contact, transfer *Request, invoiceNumber *string, token string,
	issueDate string, attachment *string) (*sendgrid.EmailResponse, error) {

	amount := shared.FormatFloatAmount(transfer.Amount)

	*transfer.Notes = strings.Replace(*transfer.Notes, "\n", "<br>", -1)

	subject := fmt.Sprintf(services.BusinessInvoiceSubject, amount, c.Name())
	body := fmt.Sprintf(services.BusinessInvoiceEmail, b.Name(), c.Name(), issueDate, *transfer.Notes, amount, amount, b.Name())

	a := sendgrid.EmailAttachment{
		Attachment:  *attachment,
		ContentType: "application/pdf",
		FileName:    "Invoice - " + *invoiceNumber + ".pdf",
		ContentID:   "Invoice",
	}

	email := sendgrid.EmailAttachmentRequest{
		SenderEmail:   os.Getenv("WISE_INVOICE_EMAIL"),
		SenderName:    os.Getenv("WISE_SUPPORT_NAME"),
		BCCEmail:      os.Getenv("WISE_FRAUDOPS_INVOICE_EMAIL"),
		BCCName:       os.Getenv("WISE_FRAUDOPS_INVOICE_NAME"),
		ReceiverEmail: *b.Email,
		ReceiverName:  b.Name(),
		Subject:       subject,
		Body:          body,
		Attachment:    []sendgrid.EmailAttachment{a},
	}

	response, err := sendgrid.NewSendGridServiceWithout().SendAttachmentEmail(email)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (db *requestDatastore) GetConnectionToken(request PaymentConnectionRequest) (*PaymentConnectionResponse, error) {
	// Check access
	err := auth.NewAuthService(db.sourceReq).CheckBusinessAccess(request.BusinessID)
	if err != nil {
		return nil, err
	}

	r := db.sourceReq.PartnerServiceRequest()

	token, err := stripe.NewStripeService(&r).GetConnectionToken()
	if err != nil {
		return nil, err
	}

	return &PaymentConnectionResponse{
		ConnectionToken: *token,
	}, nil
}

func (db *requestDatastore) CapturePayment(c PaymentCaptureRequest) error {
	// Check access
	err := auth.NewAuthService(db.sourceReq).CheckBusinessAccess(c.BusinessID)
	if err != nil {
		return err
	}

	if len(c.RequestID) == 0 {
		return errors.New("Request Id is required")
	}

	if c.PurchaseAddress == nil {
		return errors.New("Purchase address is required")
	}

	// Get payment intent ID
	payment, err := db.getPaymentIntentID(c.RequestID)
	if err != nil {
		return err
	}

	// capture payment
	r := db.sourceReq.PartnerServiceRequest()

	err = stripe.NewStripeService(&r).CapturePayment(*payment.SourcePaymentID)
	if err != nil {
		return errors.New("Error capturing payment")
	}

	// store purchase address
	_, err = db.NamedExec(
		fmt.Sprintf(
			"UPDATE business_money_request_payment SET purchase_address = :purchase_address WHERE request_id = '%s'",
			c.RequestID,
		), c,
	)
	if err != nil {
		return err
	}

	return nil
}

func (db *requestDatastore) getPaymentIntentID(requestID shared.PaymentRequestID) (*Payment, error) {
	p := Payment{}

	err := db.Get(&p, "SELECT * FROM business_money_request_payment WHERE request_id = $1", requestID)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return &p, err
}

func (db *requestDatastore) CreatePaymentIntent(pr PaymentResponse) (*stripe.PaymentResponse, error) {
	r := db.sourceReq.PartnerServiceRequest()

	metadata, err := db.getPaymentIntentMetadata(pr.MoneyRequestID, pr.BusinessID, pr.BusinessName)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if metadata.IPAddress == nil {
		metadata.IPAddress = pr.InvoiceIPAddress
	}

	p, err := stripe.NewStripeService(&r).CreatePayment(transformForPaymentService(pr.BusinessName, PaymentRequestTypeInvoiceCard, pr.Amount, CurrencyUSD, metadata))
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return p, nil
}

func (db *requestDatastore) getPaymentByPaymentRequestID(pID shared.PaymentRequestID) (*Payment, error) {
	var p Payment
	err := db.Get(&p, "SELECT * FROM business_money_request_payment WHERE request_id = $1", pID)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (db *requestDatastore) getInvoiceByPaymentRequestID(pID shared.PaymentRequestID) (*Invoice, error) {
	var inv Invoice
	err := db.Get(&inv, "SELECT * FROM business_invoice WHERE request_id = $1", pID)
	if err != nil {
		return nil, err
	}

	return &inv, nil
}

func (db *requestDatastore) ExtendInvoiceExpiryDate(requestPaymentID string, expiryDate time.Time) error {
	_, err := db.Exec("UPDATE business_money_request_payment SET expiration_date = $1 WHERE id = $2", expiryDate, requestPaymentID)
	return err
}

func (db *requestDatastore) Resend(reqs []PaymentRequestResend, bID shared.BusinessID) error {
	// Check access
	err := auth.NewAuthService(db.sourceReq).CheckBusinessAccess(bID)
	if err != nil {
		return err
	}

	// Get business
	bus, err := b.NewBusinessService(db.sourceReq).GetById(bID)
	if err != nil {
		return err
	}
	if os.Getenv("USE_INVOICE_SERVICE") == "true" {
		return newLegacyInvoiceService().SendReminder(reqs)
	} else {
		for _, r := range reqs {
			// Get payment request
			req, err := db.GetByID(r.RequestID, bID)
			if err != nil {
				log.Println("Failed to get payment request: ", r.RequestID, err)
				continue
			}

			if req.Status != nil && *req.Status == PaymentRequestStatusComplete {
				log.Println("Payment is already complete")
				continue
			}

			if req.Status != nil && *req.Status == PaymentRequestStatusCanceled {
				log.Println("Invoice has been canceled")
				continue
			}

			// Check contact
			if req.ContactID == nil {
				log.Println("Contact not found")
				continue
			}

			con, err := contact.NewContactService(db.sourceReq).GetById(*req.ContactID, bID)
			if err != nil {
				log.Println("Failed to get contact: ", req.ContactID, err)
				continue
			}

			// Get Payment
			pr, err := db.getPaymentByPaymentRequestID(r.RequestID)
			if err != nil {
				log.Println("Payment not found for payment request: ", r.RequestID, err)
				continue
			}

			// Extend expiry date by a week
			expTime := time.Now().UTC().Add(time.Hour * time.Duration(720))
			err = db.ExtendInvoiceExpiryDate(pr.ID, expTime)
			if err != nil {
				log.Println("Error extending invoice expiry date", err)
				continue
			}

			// Get invoice
			inv, err := db.getInvoiceByPaymentRequestID(r.RequestID)
			if err != nil {
				log.Println("Invoice not found for payment request: ", r.RequestID, err)
				continue
			}

			st, err := document.NewStorerFromKey(inv.StorageKey)
			if err != nil {
				log.Println("Error loading invoice from storage: ", inv.ID, err)
				continue
			}

			by, err := st.Content()
			if err != nil {
				log.Println("Error loading invoice contents: ", inv.ID, err)
				continue
			}

			content := base64.StdEncoding.EncodeToString(by)
			respID, err := sendInvoiceToCustomer(bus, con, req, &inv.InvoiceNumber, shared.StringValue(pr.PaymentToken), req.Created.Format("Jan _2, 2006"), &content, true)
			if err != nil {
				log.Println("Failed to send invoice: ", inv.InvoiceNumber, err)
				continue
			}

			// Update payment request with send grid message id
			if respID != nil {
				u := RequestIDUpdate{
					ID:         r.RequestID,
					BusinessID: bID,
					ContactID:  req.ContactID,
					MessageID:  *respID,
				}
				_, err = db.UpdateMessageID(&u)
				if err != nil {
					log.Println("Failed update message id: ", respID, err)
					continue
				}

				if os.Getenv("DEBUG") != "" {
					log.Println("Invoice sent successfully with id: ", respID)
				}
			}
		}
		return nil
	}
}

func (db *requestDatastore) getPaymentIntentMetadata(requestID *shared.PaymentRequestID, bID shared.BusinessID, businessName string) (*RequestMetadata, error) {
	metadata := RequestMetadata{}

	var query string
	var err error

	if os.Getenv("USE_BANKING_SERVICE") == "true" {
		if requestID != nil {
			query = `select first_name, middle_name, last_name, request_ip_address
		FROM business
		JOIN wise_user ON business.owner_id = wise_user.id
		JOIN consumer ON consumer.id = wise_user.consumer_id
		JOIN business_money_request ON business_money_request.business_id = business.id
		WHERE business.id = $1 AND business_money_request.id = $2`

			err = db.Get(&metadata, query, bID, *requestID)
		} else {
			query = `select first_name, middle_name, last_name
		FROM business
		JOIN wise_user ON business.owner_id = wise_user.id
		JOIN consumer ON consumer.id = wise_user.consumer_id
		WHERE business.id = $1`

			err = db.Get(&metadata, query, bID)
		}
		if err != nil {
			return nil, err
		}

		bas, err := business.NewBankingAccountService()
		if err != nil {
			return nil, err
		}

		//TODO If this is still around when we allow multiple bank accounts per business this needs to be refactored
		bs, err := bas.GetByBusinessID(bID, 1, 0)
		if err != nil {
			return nil, err
		}

		if len(bs) != 1 {
			return nil, errors.New("Could not find bank account for business")
		}

		ba := bs[0]

		metadata.AvailableBalance = ba.AvailableBalance

	} else {
		if requestID != nil {
			query = `select available_balance, first_name, middle_name, last_name, request_ip_address
		FROM business
		JOIN business_bank_account ON business.id = business_bank_account.business_id
		JOIN wise_user ON business.owner_id = wise_user.id
		JOIN consumer ON consumer.id = wise_user.consumer_id
		JOIN business_money_request ON business_money_request.business_id = business.id
		WHERE business.id = $1 AND business_money_request.id = $2`

			err = db.Get(&metadata, query, bID, *requestID)
		} else {
			query = `select available_balance, first_name, middle_name, last_name
		FROM business
		JOIN business_bank_account ON business.id = business_bank_account.business_id
		JOIN wise_user ON business.owner_id = wise_user.id
		JOIN consumer ON consumer.id = wise_user.consumer_id
		WHERE business.id = $1`

			err = db.Get(&metadata, query, bID)
		}
	}

	if err != nil {
		log.Println(err)
		return nil, err
	}

	metadata.BusinessName = businessName
	return &metadata, nil
}

func (db *requestDatastore) UpdateInvoiceStatus(invoiceId id.InvoiceID, status PaymentRequestStatus) error {
	invSvc, err := invoice.NewInvoiceService()
	if err != nil {
		log.Println(err)
		return err
	}
	// check if we have invoice for this
	inv, err := invSvc.GetInvoiceByID(invoiceId)
	if err != nil {
		log.Println(err)
		return err
	}

	invSvc, err = invoice.NewInvoiceService()
	if err != nil {
		log.Println(err)
		return err
	}

	err = invSvc.UpdateInvoiceStatus(invoiceId, inv.BusinessID, GetInvoiceGrpcStatus(status))
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
