package payment

import (
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"strings"
	"time"

	"github.com/wiseco/core-platform/services/invoice"

	"github.com/wiseco/core-platform/services"
	"github.com/wiseco/core-platform/services/banking"
	"github.com/wiseco/core-platform/services/banking/business"
	b "github.com/wiseco/core-platform/services/business"
	"github.com/wiseco/core-platform/services/contact"
	"github.com/wiseco/core-platform/shared"
	"github.com/wiseco/go-lib/grpc"
	"github.com/wiseco/go-lib/id"
	"github.com/wiseco/go-lib/num"
	grpcShopify "github.com/wiseco/protobuf/golang/shopping/shopify"
)

type RequestJoin struct {
	PaymentID           string                   `db:"id"`
	PaymentStatus       string                   `db:"status"`
	SourcePaymentID     string                   `db:"source_payment_id"`
	RequestID           *shared.PaymentRequestID `db:"request_id"`
	ContactID           *string                  `db:"business_contact.id"`
	InvoiceID           *string                  `db:"business_invoice.id"`
	RequestStatus       PaymentRequestStatus     `db:"request_status"`
	UserID              shared.UserID            `db:"created_user_id"`
	BusinessID          shared.BusinessID        `db:"business_id"`
	Amount              float64                  `db:"amount"`
	FirstName           string                   `db:"consumer.first_name"`
	LastName            string                   `db:"consumer.last_name"`
	Email               string                   `db:"business.email"`
	ContactFirstName    *string                  `db:"business_contact.first_name"`
	ContactLastName     *string                  `db:"business_contact.last_name"`
	ContactEmail        *string                  `db:"business_contact.email"`
	ContactPhone        *string                  `db:"business_contact.phone_number"`
	ContactBusinessName *string                  `db:"business_contact.business_name"`
	BusinessName        string
	LegalName           *string              `db:"legal_name"`
	DBA                 services.StringArray `db:"dba"`
	BusinessPhone       string               `db:"business.phone"`
	Notes               *string              `db:"notes"`
	InvoiceNumber       *string              `db:"invoice_number"`
	RequestType         *PaymentRequestType  `db:"request_type"`
	RequestSource       *RequestSource       `db:"request_source"`
	InvoiceIdV2         *id.InvoiceID
	InvoiceViewLink     string
	InvoiceReceiptLink  string
}

func (db *paymentDatastore) HandleWebhook(payment *Payment) error {

	r := RequestJoin{}

	if payment.Status == PaymentStatusSucceeded {
		paymentDbInfo := Payment{}
		isPOSPayment := db.isPOSPayment(payment.ID)
		if os.Getenv("USE_INVOICE_SERVICE") == "true" && !isPOSPayment {
			err := db.Get(&paymentDbInfo, "select * from business_money_request_payment where source_payment_id = $1", payment.ID)
			if err != nil {
				log.Println(err)
				return err
			}
			invSvc, err := invoice.NewInvoiceService()
			if err != nil {
				log.Println(err)
				return err
			}
			log.Println("------Processing webhook with invoice service------")
			log.Println(*paymentDbInfo.InvoiceID)
			inv, err := invSvc.GetInvoiceByID(*paymentDbInfo.InvoiceID)
			if err != nil {
				log.Println(err)
				return err
			}

			err = db.Get(
				&r, `
			SELECT
				business_money_request_payment.id,
				business_money_request_payment.status,				
				business_money_request_payment.request_id,
				business_money_request_payment.source_payment_id,
				business.legal_name,
				business.phone "business.phone",
				business.dba,
				business.email "business.email",
				business_contact.first_name "business_contact.first_name",
				business_contact.last_name "business_contact.last_name",
				business_contact.email "business_contact.email",
				business_contact.phone_number "business_contact.phone_number",
				business_contact.id "business_contact.id",
				business_contact.business_name "business_contact.business_name",
				consumer.first_name "consumer.first_name",
				consumer.last_name "consumer.last_name"
				FROM business_money_request_payment,business, business_contact business_contact,
				consumer inner join wise_user on consumer.id = wise_user.consumer_id
				WHERE business_money_request_payment.source_payment_id = $1
				and business.id = $2
				and business_contact.id = $3
				and wise_user.id = $4`,
				payment.ID,
				inv.BusinessID.UUIDString(),
				inv.ContactID.UUIDString(),
				inv.UserID.UUIDString())
			if userid, err := shared.ParseUserID(inv.UserID.String()); err != nil {
				return err
			} else {
				r.UserID = userid
			}
			if busId, err := shared.ParseBusinessID(inv.BusinessID.String()); err != nil {
				return err
			} else {
				r.BusinessID = busId
			}

			contactId := inv.ContactID.UUIDString()
			r.ContactID = &contactId

			if userid, err := shared.ParseUserID(inv.UserID.String()); err != nil {
				return err
			} else {
				r.UserID = userid
			}
			r.Notes = &inv.Notes
			reqType := PaymentRequestTypeInvoiceCard
			if inv.AllowCard {
				r.RequestType = &reqType
			}
			if amtFlt, ok := inv.Amount.Float64(); !ok {
				return err
			} else {
				r.Amount = amtFlt
			}
			r.RequestStatus = GetInvoiceStatus(inv.Status)
			r.InvoiceIdV2 = paymentDbInfo.InvoiceID
			r.RequestID = nil
			r.InvoiceViewLink = inv.InvoiceViewLink
			encodedInvoiceID := base64.RawURLEncoding.EncodeToString([]byte(r.InvoiceIdV2.String()))
			r.InvoiceReceiptLink = fmt.Sprintf("%s/invoice-receipt?token=%s",
				os.Getenv("PAYMENTS_URL"), encodedInvoiceID)
			if inv.RequestSource != "" && inv.RequestSource == string(RequestSourceShopify) {
				shopify := RequestSourceShopify
				r.RequestSource = &shopify
			}
			log.Println(fmt.Sprintf("setting invoiceidv2 as %s", paymentDbInfo.InvoiceID))
		} else {
			log.Println("--processing webhook with old invoice-----")
			err := db.Get(
				&r, `
			SELECT
				business_money_request.created_user_id,
				business_money_request.business_id,
				business_money_request.notes,
				business_money_request.request_type,
				business_money_request.request_source,
				business_invoice.invoice_number,
				business_invoice.id "business_invoice.id",
				business_money_request.amount,
				business_money_request.request_status,
				business_money_request_payment.id,
				business_money_request_payment.status,
				business.legal_name,
				business.phone "business.phone",
				business.dba,
				business_money_request_payment.request_id,
				business_money_request_payment.source_payment_id,
				business_contact.first_name "business_contact.first_name",
				business_contact.last_name "business_contact.last_name",
				business_contact.email "business_contact.email",
				business_contact.phone_number "business_contact.phone_number",
				business_contact.id "business_contact.id",
				business_contact.business_name "business_contact.business_name",
				business.email "business.email", 
				consumer.first_name "consumer.first_name",
				consumer.last_name "consumer.last_name"
			FROM business_money_request
			JOIN business_money_request_payment ON business_money_request.id = business_money_request_payment.request_id
			LEFT JOIN business_invoice ON business_money_request.id = business_invoice.request_id
			LEFT JOIN business_contact ON business_contact.id = business_money_request.contact_id
			JOIN wise_user ON business_money_request.created_user_id = wise_user.id
			JOIN business ON business.id = business_money_request.business_id
			JOIN consumer ON consumer.id = wise_user.consumer_id
			WHERE business_money_request_payment.source_payment_id = $1`,
				payment.ID,
			)
			if err != nil {
				if err == sql.ErrNoRows {
					return errors.New(fmt.Sprintf("Payment with id:%s not found", payment.ID))
				}

				return err
			}
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

		if payment.PaymentDate == nil {
			return errors.New("payment date is nil")
		}

		paidDate := (*payment.PaymentDate).Format("Jan _2, 2006")

		// Backward compatibility
		if r.RequestType == nil {
			reqType := PaymentRequestTypeInvoiceCard
			r.RequestType = &reqType
		}

		tz := os.Getenv("BATCH_TZ")
		if tz == "" {
			panic(errors.New("Local timezone missing"))
		}

		loc, err := time.LoadLocation(tz)
		if err != nil {
			panic(err)
		}

		nowUTC := time.Now().UTC()
		nowLocal := nowUTC.In(loc)

		// Subscription starts from July 1st
		wiseSubscriptionStartDate := time.Date(2020, time.July, 1, 0, 0, 0, 0, loc)
		var feeAmountDecimal num.Decimal
		invoiceAmount, transferAmount := r.Amount, r.Amount

		if wiseSubscriptionStartDate.Before(nowLocal) {
			feeAmount := (invoiceAmount * 3) / 100
			feeAmountRounded := math.Round(feeAmount*100) / 100

			transferAmount = invoiceAmount - feeAmountRounded

			fVal, err := num.NewDecimalFin(feeAmountRounded)
			if err != nil {
				return err
			}

			feeAmountDecimal = num.Decimal{V: fVal}
		} else {
			fVal, err := num.NewDecimalFin(0)
			if err != nil {
				return err
			}

			feeAmountDecimal = num.Decimal{V: fVal}
		}

		log.Println("Date check:", wiseSubscriptionStartDate, nowLocal)

		// Check to make sure payment is not already made
		if r.RequestStatus != PaymentRequestStatusComplete && r.RequestStatus != PaymentRequestStatusCanceled {
			// Email sender
			rNumber := shared.GetRandomString("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789", 8) + "-" +
				shared.GetRandomString("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789", 5)

			var content *string
			if *r.RequestType != PaymentRequestTypePOS {
				receiptGenerate := ReceiptGenerate{
					ContactFirstName:    r.ContactFirstName,
					ContactLastName:     r.ContactLastName,
					ContactBusinessName: r.ContactBusinessName,
					ContactEmail:        r.ContactEmail,
					PaymentDate:         payment.PaymentDate,
					PaymentBrand:        payment.CardBrand,
					PaymentNumber:       payment.CardLast4,
					ReceiptNumber:       rNumber,
					InvoiceNumber:       r.InvoiceNumber,
					BusinessName:        r.BusinessName,
					BusinessPhone:       r.BusinessPhone,
					Amount:              r.Amount,
					Notes:               r.Notes,
					UserID:              r.UserID,
					BusinessID:          r.BusinessID,
					ContactID:           r.ContactID,
					InvoiceID:           r.InvoiceID,
					RequestID:           r.RequestID,
					InvoiceIdV2:         r.InvoiceIdV2,
				}
				content, _, err = db.GenerateReceipt(receiptGenerate)
				if err != nil {
					log.Println("Receipt generation failed ", err)
				}

				receiptEmail := ReceiptRequest{
					ContactFirstName:    r.ContactFirstName,
					ContactLastName:     r.ContactLastName,
					ContactEmail:        r.ContactEmail,
					ContactPhone:        r.ContactPhone,
					ContactBusinessName: r.ContactBusinessName,
					BusinessName:        r.BusinessName,
					Amount:              r.Amount,
					Notes:               r.Notes,
					PaymentDate:         paidDate,
					ReceiptNumber:       rNumber,
					Content:             content,
				}

				db.SendReceiptToCustomer(receiptEmail, r.InvoiceViewLink, r.InvoiceReceiptLink)
			}

			// To pass through auth controls
			db.sourceReq.UserID = shared.UserID(clearingUserID)

			// Get wise clearing linked account number. Required for money transfer
			wiseLinkedAccount, err := business.NewLinkedAccountService(*db.sourceReq).GetByAccountIDInternal(shared.BusinessID(clearingBusinessID), clearingAccountID)
			if err != nil && err != sql.ErrNoRows {
				log.Println(err)
				return err
			}

			if err != nil {
				return errors.New("Wise account needs to be registered before sending money to business")
			}

			// Get requester business' account number
			accounts, err := business.NewBankAccountService(*db.sourceReq).GetByUserID(r.UserID, r.BusinessID)
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
			linkedAccount, err := business.NewLinkedAccountService(*db.sourceReq).GetByAccountNumber(
				shared.BusinessID(clearingBusinessID),
				business.AccountNumber(account.AccountNumber),
				account.RoutingNumber,
			)
			if err != nil && err != sql.ErrNoRows {
				log.Println(err)
				return err
			}

			// Check if account is already registered by wise clearing account
			if err != nil {
				// Step 1: Wise clearing account links requester's account
				laCreate := business.MerchantLinkedAccountCreate{
					UserID:            shared.UserID(clearingUserID),
					BusinessID:        shared.BusinessID(clearingBusinessID),
					AccountHolderName: r.BusinessName,
					AccountNumber:     business.AccountNumber(account.AccountNumber),
					AccountType:       banking.AccountType(account.AccountType),
					RoutingNumber:     account.RoutingNumber,
					Currency:          account.Currency,
					Permission:        banking.LinkedAccountPermissionSendAndRecieve,
				}
				la, err := business.NewLinkedAccountService(*db.sourceReq).LinkMerchantBankAccount(&laCreate)
				if err != nil {
					log.Println(err)
					return err
				}

				// Step 2: Move money
				transferInitiate := business.TransferInitiate{
					CreatedUserID:   shared.UserID(clearingUserID),
					BusinessID:      shared.BusinessID(clearingBusinessID),
					SourceAccountId: wiseLinkedAccount.Id,
					DestAccountId:   la.Id,
					Amount:          transferAmount,
					SourceType:      banking.TransferTypeAccount,
					DestType:        banking.TransferTypeAccount,
					Currency:        banking.CurrencyUSD,
					MoneyRequestID:  r.RequestID,
					Notes:           r.Notes,
				}
				if transferInitiate.MoneyRequestID == nil && r.InvoiceIdV2 != nil {
					transferInitiate.MoneyRequestID = getMoneyRequestId(*r.InvoiceIdV2)
				}

				s := db.sourceReq
				_, err = business.NewMoneyTransferService(*s).Transfer(&transferInitiate)
				if err != nil {
					log.Println(err)
					return err
				}
			} else {
				// Move money
				transferInitiate := business.TransferInitiate{
					CreatedUserID:   shared.UserID(clearingUserID),
					BusinessID:      shared.BusinessID(clearingBusinessID),
					SourceAccountId: wiseLinkedAccount.Id,
					DestAccountId:   linkedAccount.Id,
					Amount:          transferAmount,
					SourceType:      banking.TransferTypeAccount,
					DestType:        banking.TransferTypeAccount,
					Currency:        banking.CurrencyUSD,
					MoneyRequestID:  r.RequestID,
					Notes:           r.Notes,
				}
				if transferInitiate.MoneyRequestID == nil && r.InvoiceIdV2 != nil {
					transferInitiate.MoneyRequestID = getMoneyRequestId(*r.InvoiceIdV2)
				}

				s := db.sourceReq
				_, err = business.NewMoneyTransferService(*s).Transfer(&transferInitiate)
				if err != nil {
					log.Println(err)
					return err
				}
			}

			// Send email to business
			if *r.RequestType != PaymentRequestTypePOS {
				receiptEmail := ReceiptRequest{
					ContactFirstName:    r.ContactFirstName,
					ContactLastName:     r.ContactLastName,
					ContactEmail:        r.ContactEmail,
					ContactBusinessName: r.ContactBusinessName,
					BusinessEmail:       r.Email,
					BusinessName:        r.BusinessName,
					Amount:              r.Amount,
					Notes:               r.Notes,
					PaymentDate:         paidDate,
					ReceiptNumber:       rNumber,
					Content:             content,
				}
				db.SendReceiptToBusiness(receiptEmail, r.InvoiceViewLink)
			}

			// If invoiceIdv2 is not null , then send update request to invoice service
			if os.Getenv("USE_INVOICE_SERVICE") == "true" && r.InvoiceIdV2 != nil {
				log.Println("updating the invoice service")
				invSvc, err := invoice.NewInvoiceService()
				if err != nil {
					log.Println(err)
					return err
				}
				businessIdParsed, err := id.ParseBusinessID(r.BusinessID.ToPrefixString())
				if err != nil {
					return err
				}

				err = invSvc.UpdateInvoiceStatus(*r.InvoiceIdV2, businessIdParsed, GetInvoiceGrpcStatus(PaymentRequestStatusComplete))
				if err != nil {
					log.Println(err)
					return err
				}
			} else if r.RequestID != nil { // Update request money status in business_money_request table
				requestUpdate := RequestUpdate{
					ID:     shared.PaymentRequestID(*r.RequestID),
					Status: PaymentRequestStatusComplete,
				}

				if *r.RequestType == PaymentRequestTypeInvoiceCardAndBank {
					reqType := PaymentRequestTypeInvoiceCard
					requestUpdate.RequestType = &reqType
				}

				err = db.UpdateRequestStatus(&requestUpdate)
				if err != nil {
					return err
				}
			}

			receiptToken := shared.GetRandomString("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789", 16)

			// Update payment status in business_money_request_payment table and set token to null
			paymentUpdate := Payment{
				ID:           r.PaymentID,
				Status:       payment.Status,
				PaymentToken: nil,
				CardBrand:    payment.CardBrand,
				CardLast4:    payment.CardLast4,
				WalletType:   payment.WalletType,
				PaymentDate:  payment.PaymentDate,
				ReceiptToken: &receiptToken,
				RequestID:    r.RequestID,
				FeeAmount:    &feeAmountDecimal,
				InvoiceID:    r.InvoiceIdV2,
			}

			err = db.UpdatePaymentStatus(&paymentUpdate)
			if err != nil {
				return err
			}

			// check if pos receipt is set or not
			if *r.RequestType == PaymentRequestTypePOS {
				p, err := db.GetPaymentByRequestID(*r.RequestID)
				if err != nil {
					return err
				}

				if p.ReceiptMode == nil {
					return nil
				}

				return db.sendCardReaderReceipts(p)
			}
		}

	}

	return nil
}

func (db *paymentDatastore) createContact(bID shared.BusinessID, uID shared.UserID, wiserUserID shared.UserID, clearingBusinessID shared.BusinessID, firstName, lastName string) (*contact.Contact, error) {

	// To pass through auth controls
	db.sourceReq.UserID = uID

	b, err := b.NewBusinessService(*db.sourceReq).GetById(bID)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	} else if err != nil && err == sql.ErrNoRows {
		return nil, errors.New("business not found")
	}

	businessName := b.Name()
	contactCreate := contact.ContactCreate{
		UserID:       wiserUserID,
		BusinessID:   clearingBusinessID,
		Type:         contact.ContactTypeBusiness,
		BusinessName: &businessName,
		PhoneNumber:  *b.Phone,
		Email:        *b.Email,
	}

	// To pass through auth controls
	db.sourceReq.UserID = wiserUserID

	c, err := contact.NewContactService(*db.sourceReq).Create(&contactCreate)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func getMoneyRequestId(invoiceId id.InvoiceID) *shared.PaymentRequestID {
	monReqId, err := shared.ParsePaymentRequestID(fmt.Sprintf("%s%s", shared.PaymentRequestPrefix, invoiceId.UUIDString()))
	if err != nil {
		return nil
	}
	return &monReqId
}

func (db *paymentDatastore) UpdateRequestStatus(u *RequestUpdate) error {

	var columns []string

	columns = append(columns, "request_status = :request_status")

	if u.RequestType != nil {
		columns = append(columns, "request_type = :request_type")
	}

	_, err := db.NamedExec(fmt.Sprintf("UPDATE business_money_request SET %s WHERE id = '%s'", strings.Join(columns, ", "), u.ID), u)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (db *paymentDatastore) UpdatePaymentStatus(u *Payment) error {
	var columns []string

	columns = append(columns, "payment_token = :payment_token")

	if u.Status != "" {
		columns = append(columns, "status = :status")
	}

	if u.SourcePaymentID != nil {
		columns = append(columns, "source_payment_id = :source_payment_id")
	}

	if u.CardBrand != nil {
		columns = append(columns, "card_brand = :card_brand")
	}

	if u.CardLast4 != nil {
		columns = append(columns, "card_number = :card_number")
	}

	if u.ReceiptToken != nil {
		columns = append(columns, "receipt_token = :receipt_token")
	}

	if u.PaymentDate != nil {
		columns = append(columns, "payment_date = :payment_date")
	}

	if u.LinkedBankAccountID != nil {
		columns = append(columns, "linked_bank_account_id = :linked_bank_account_id")
	}

	if u.FeeAmount != nil {
		columns = append(columns, "fee_amount = :fee_amount")
	}

	if u.WalletType != nil {
		columns = append(columns, "wallet_type = :wallet_type")
	}

	_, err := db.NamedExec(fmt.Sprintf("UPDATE business_money_request_payment SET %s WHERE id = '%s'", strings.Join(columns, ", "), u.ID), u)
	if err != nil {
		log.Println(err)
		return err
	}
	if os.Getenv("USE_INVOICE_SERVICE") == "true" && u.InvoiceID != nil {
		invSvc, err := invoice.NewInvoiceService()
		if err != nil {
			log.Println(err)
			return err
		}
		invoiceModel, err := invSvc.GetInvoiceByID(*u.InvoiceID)
		if err != nil {
			log.Println(err)
			return err
		}
		if invoiceModel.RequestSource != "" && invoiceModel.RequestSource == string(RequestSourceShopify) {
			log.Println("updating the shopify request")
			err := db.updateShopifyRequest(GetInvoiceStatus(invoiceModel.Status), invoiceModel.BusinessID.String(),
				invoiceModel.Amount.FormatCurrency(), invoiceModel.RequestSourceID)
			if err != nil {
				log.Println(err)
				return err
			}
		}
	} else if u.RequestID != nil {
		req, err := NewRequestService(*db.sourceReq).GetByIDInternal(shared.PaymentRequestID(*u.RequestID))
		if err != nil {
			log.Println("Error finding request ID", err)
			return err
		}

		if req.RequestSource != nil && *req.RequestSource == RequestSourceShopify {
			amount := fmt.Sprintf("%f", req.Amount)
			err := db.updateShopifyRequest(*req.Status, req.BusinessID.ToPrefixString(), amount, *req.RequestSourceID)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (db *paymentDatastore) updateShopifyRequest(status PaymentRequestStatus, businessId string, amount string, sourceId string) error {

	var txnType grpcShopify.ShopifyOrderTransactionType
	if status == PaymentRequestStatusComplete {
		txnType = grpcShopify.ShopifyOrderTransactionType_SOTT_CAPTURE
	} else if status == PaymentRequestStatusCanceled {
		txnType = grpcShopify.ShopifyOrderTransactionType_SOTT_VOID
	} else {
		return nil
	}

	sn, err := grpc.GetConnectionStringForService(grpc.ServiceNameShopping)
	if err != nil {
		return err
	}

	client, err := grpc.NewInsecureClient(sn)
	if err != nil {
		return err
	}

	defer client.CloseAndCancel()
	shopifyServiceClient := grpcShopify.NewShopifyBusinessServiceClient(client.GetConn())

	req := grpcShopify.UpdateOrderPaymentRequest{
		BusinessId:      businessId,
		Amount:          amount,
		OrderId:         fmt.Sprintf("%s%s", id.IDPrefixShopifyOrder, sourceId),
		TransactionType: txnType,
	}
	log.Println(fmt.Sprintf("shopify request obj: %+v", req))
	order, err := shopifyServiceClient.UpdateOrderPaymentStatus(client.GetContext(), &req)
	if err != nil {
		return err
	}

	log.Println("Order status updated:", order.Id)

	return nil
}

func (db *paymentDatastore) isPOSPayment(sourcePaymentID string) bool {
	query := `SELECT business_money_request.request_type 
				FROM business_money_request JOIN business_money_request_payment 
				ON business_money_request.id = business_money_request_payment.request_id 
				WHERE business_money_request_payment.source_payment_id = $1`
	reqObj := RequestJoin{}
	err := db.Get(&reqObj, query, sourcePaymentID)
	if err != nil {
		return false
	}
	if reqObj.RequestType != nil && *reqObj.RequestType == PaymentRequestTypePOS {
		return true
	}
	return false
}
