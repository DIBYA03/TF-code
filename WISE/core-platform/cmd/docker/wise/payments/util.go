package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/stripe/stripe-go"

	"github.com/wiseco/core-platform/services/contact"

	"github.com/wiseco/core-platform/services/invoice"

	"github.com/wiseco/core-platform/services"
	"github.com/wiseco/core-platform/services/banking"
	"github.com/wiseco/core-platform/services/banking/business"
	"github.com/wiseco/core-platform/services/payment"
	"github.com/wiseco/core-platform/shared"
	"github.com/wiseco/go-lib/id"
	grpcInvoice "github.com/wiseco/protobuf/golang/invoice"
)

func loadHTMLTemplate(t HTMLTemplate, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	var data Request
	var p *payment.PaymentResponse
	var err error

	params, ok := r.URL.Query()["token"]
	if !ok || len(params[0]) < 1 {
		renderError("requested url is invalid or has expired", w)
		return
	}

	var paymentToken = params[0]
	var invoiceV2Payment = false
	var invoiceID id.InvoiceID
	// for invoice 2.0 payment request
	//
	// get the paymentToken from invoice_id
	// if not found, then create any entry and use that payment token onwards
	if isInvoiceV2Payment(t, r) {
		invoiceV2Payment = true
		invoiceID, err := getInvoiceIDFromToken(paymentToken, w)
		if err != nil {
			renderError(err.Error(), w)
			return
		}
		fmt.Println(invoiceID)
		if t == HTMLTemplateInvoiceView {
			renderInvoiceView(invoiceID, w)
			return
		}

		if t == HTMLTemplateInvoiceSuccess {
			renderInvoiceReceipt(invoiceID, w)
			return
		}

		// fetch the payment token
		paymentTokenForInvoice, isFound, err := payment.NewPaymentServiceInternal().GetPaymentTokenFromInvoice(invoiceID.UUIDString())
		if err != nil {
			renderError(err.Error(), w)
			return
		}
		if !isFound {
			paymentTokenForInvoice, err := createPaymentToken(invoiceID)
			if err != nil {
				renderError(err.Error(), w)
				return
			}
			paymentToken = *paymentTokenForInvoice
			println(*paymentTokenForInvoice)
		} else if paymentTokenForInvoice.PaymentToken != nil {
			paymentToken = *paymentTokenForInvoice.PaymentToken
		} else {
			paymentToken = ""
		}
	}
	p, err = payment.NewPaymentServiceInternal().GetPaymentInfo(paymentToken)
	if err != nil {
		renderError(err.Error(), w)
		return
	}

	amount := shared.FormatFloatAmount(p.Amount)
	p.Notes = strings.Replace(p.Notes, "\n", ", ", -1)

	data = Request{
		BusinessName: p.BusinessName,
		Amount:       amount,
		Notes:        p.Notes,
		StripeKey:    p.StripeKey,
		PaymentToken: paymentToken,
		SegmentKey:   os.Getenv("SEGMENT_WEB_WRITE_KEY"),
		UserID:       p.UserID.ToPrefixString(),
	}
	if p.InvoiceID != nil && !p.InvoiceID.IsZero() {
		invoiceV2Payment = true
		invoiceID = *p.InvoiceID
	}
	if invoiceV2Payment {
		invoiceSvc, err := invoice.NewInvoiceService()
		if err != nil {
			renderError(err.Error(), w)
			return
		}

		invoiceModel, err := invoiceSvc.GetInvoiceByID(invoiceID)
		if err != nil {
			renderError(err.Error(), w)
		}

		data.ShowBankDetails = invoiceModel.ShowBankAccount
		data.InvoiceNumber = invoiceModel.Number
		data.DueDate = invoiceModel.Created
		data.InvoiceTitle = invoiceModel.Title
		data.AllowCardPayment = invoiceModel.AllowCard
		data.AccountNumber = invoiceModel.AccountNumber
		data.RoutingNumber = invoiceModel.RoutingNumber
		data.BusinessLogo = invoiceModel.BusinessLogo
		// This is for backward compatibility
		// if we get invoiceid from payment table, and still the template requested in index
		// then force it to render the new invoice payment page.
		if t == HTMLTemplateIndex {
			t = HTMLTemplateInvoiceRequestIndex
		}
		data.ClientSecret = getPaymentIntentClientSecret(*p, paymentToken)
	}
	data.Country = "US"
	data.Currency = string(payment.CurrencyUSD)
	data.AmountStripe = stripe.Int64(int64(p.Amount * 100))
	if err != nil {
		// Return error html page
		renderError(err.Error(), w)
		return
	}

	switch p.RequestType {
	case payment.PaymentRequestTypeInvoiceCardAndBank:
		data.AllowBankPayment = true
		data.AllowCardPayment = true
		fallthrough
	case payment.PaymentRequestTypeInvoiceBank:
		data.AllowBankPayment = true
	case payment.PaymentRequestTypeInvoiceCard:
		data.AllowCardPayment = true
	}

	switch p.Status {
	case payment.PaymentRequestStatusInProcess:
		if t != HTMLTemplateBankAccountSuccess {
			data.PaymentDate = p.PaymentDate.Format("Jan _2, 2006")

			receivedDate := p.PaymentDate.Add(time.Hour * time.Duration(72))
			data.PaymentReceivedDate = receivedDate.Format("Jan _2, 2006")

			tmpl := template.Must(template.ParseFiles(string(HTMLTemplateInvoicePaid)))
			tmpl.Execute(w, data)

			return
		}
	case payment.PaymentRequestStatusComplete:
		data.ErrorMessage = "requested url is invalid or has expired"

		tmpl := template.Must(template.ParseFiles(string(HTMLTemplateError)))
		tmpl.Execute(w, data)

		return
	case payment.PaymentRequestStatusCanceled:
		data.ErrorMessage = "invoice has been canceled"

		tmpl := template.Must(template.ParseFiles(string(HTMLTemplateError)))
		tmpl.Execute(w, data)

		return
	}

	switch t {
	case HTMLTemplateIndex:
		if p.Status == payment.PaymentRequestStatusPending {
			tmpl := template.Must(template.ParseFiles(string(HTMLTemplateIndex)))
			tmpl.Execute(w, data)
		}
	case HTMLTemplateInvoiceRequestIndex:
		if p.Status == payment.PaymentRequestStatusPending {
			funcMap := getInvoiceTemplateFuncMap(data.BusinessLogo)
			// In template parsing, Funcs is supposed to go before ParseFiles
			tmpl := template.Must(template.New(string(HTMLTemplateInvoiceRequestIndex)).Funcs(funcMap).ParseFiles(string(HTMLTemplateInvoiceRequestIndex)))
			tmpl.Execute(w, data)
		}
	case HTMLTemplateCardPayment:
		data.ClientSecret = getPaymentIntentClientSecret(*p, paymentToken)

		tmpl := template.Must(template.ParseFiles(string(HTMLTemplateCardPayment)))
		tmpl.Execute(w, data)
	case HTMLTemplateBankAccountPayment:
		data.PublicPlaidKey = os.Getenv("PLAID_PUBLIC_KEY")
		data.PlaidEnv = os.Getenv("PLAID_ENV")

		tmpl := template.Must(template.ParseFiles(string(HTMLTemplateBankAccountPayment)))
		tmpl.Execute(w, data)
	case HTMLTemplateBankAccountSuccess:
		tmpl := template.Must(template.ParseFiles(string(HTMLTemplateBankAccountSuccess)))
		tmpl.Execute(w, data)
	case HTMLTemplateCardSuccess:
		// Set token to null to prevent further payments
		sourceRequest := services.NewSourceRequest()
		paymentUpdate := payment.Payment{
			ID: p.PaymentID,
		}
		err = payment.NewPaymentService(sourceRequest).UpdatePaymentStatus(&paymentUpdate)
		if err != nil {
			log.Println("Error updating payment status ", err)

		}
		if os.Getenv("USE_INVOICE_SERVICE") == "true" {
			successResp := getInvoiceSuccessResponse(r, p)
			currentTime := time.Now()
			funcMap := getHTMLTemplateInvoiceSuccessFuncMap(&currentTime)
			tmpl := template.Must(template.New(string(HTMLTemplateInvoiceSuccess)).Funcs(funcMap).ParseFiles(string(HTMLTemplateInvoiceSuccess)))
			tmpl.Execute(w, successResp)

		} else {
			tmpl := template.Must(template.ParseFiles(string(HTMLTemplateCardSuccess)))
			tmpl.Execute(w, data)
		}
	case HTMLTemplateInvoicePaid:
		data.PaymentDate = p.PaymentDate.Format("Jan _2, 2006")
		receivedDate := p.PaymentDate.Add(time.Hour * time.Duration(72))
		data.PaymentReceivedDate = receivedDate.Format("Jan _2, 2006")

		tmpl := template.Must(template.ParseFiles("invoice-paid.html"))
		tmpl.Execute(w, data)
	case HTMLTemplatePlaid:
		accountID := r.FormValue("accountId")
		publicToken := r.FormValue("publicToken")

		err = bankTransferPayment(p, paymentToken, accountID, publicToken)
		if err != nil {
			log.Println("error in bank transfer ", err)
			http.Error(w, "error moving money", 500)
			return
		}

		w.Write([]byte("success"))
	default:
		data.ErrorMessage = "requested url is invalid or has expired"

		tmpl := template.Must(template.ParseFiles(string(HTMLTemplateError)))
		tmpl.Execute(w, data)
	}
}

func getPaymentIntentClientSecret(p payment.PaymentResponse, paymentToken string) string {
	// Check if payment intent already exits
	if p.ClientSecret != nil {
		return *p.ClientSecret
	}

	// Create payment intent
	paymentIntent, err := payment.NewRequestService(services.NewSourceRequest()).CreatePaymentIntent(p)
	if err != nil {
		log.Println("error creating payment intent", err)
	}

	// update payments table
	paymentUpdate := payment.Payment{
		ID:              p.PaymentID,
		SourcePaymentID: &paymentIntent.IntentID,
		Status:          payment.PartnerTransferStatusToPaymentStatus[string(paymentIntent.Status)],
		PaymentToken:    &paymentToken,
	}

	err = payment.NewPaymentService(services.NewSourceRequest()).UpdatePaymentStatus(&paymentUpdate)
	if err != nil {
		log.Println("error updating payment")
	}

	return paymentIntent.ClientSecret

}

func bankTransferPayment(p *payment.PaymentResponse, paymentToken, accountID, token string) error {
	//1. Link account
	l := business.LinkedExternalAccountCreate{
		BusinessID:      p.BusinessID,
		ContactID:       p.ContactID,
		PublicToken:     token,
		SourceAccountId: accountID,
		UserID:          p.UserID,
	}

	sourceRequest := services.NewSourceRequest()
	sourceRequest.UserID = p.UserID

	la, err := business.NewLinkedAccountService(sourceRequest).LinkExternalBankAccount(&l)
	if la == nil {
		// show error message here
		log.Println("Error linking bank account ", err)
		return err
	}

	//2. Move money
	m := business.TransferInitiate{
		BusinessID:      p.BusinessID,
		ContactId:       p.ContactID,
		CreatedUserID:   p.UserID,
		Amount:          p.Amount,
		Currency:        banking.CurrencyUSD,
		SourceAccountId: la.Id,
		SourceType:      banking.TransferTypeAccount,
		DestAccountId:   p.RegisteredAccountID,
		DestType:        banking.TransferTypeAccount,
		MoneyRequestID:  p.MoneyRequestID,
		Notes:           &p.Notes,
	}
	mm, err := business.NewMoneyTransferService(sourceRequest).Transfer(&m)
	if err != nil {
		log.Println("Error moving money", err)
		return err
	}

	s := la.Id
	if !strings.HasPrefix(s, id.IDPrefixLinkedBankAccount.String()) {
		s = fmt.Sprintf("%s%s", id.IDPrefixLinkedBankAccount, s)
	}

	plaID, err := id.ParseLinkedBankAccountID(s)
	if err != nil {
		return err
	}

	laID := plaID.UUIDString()

	//3. Update payments table
	accNumber := string(la.AccountNumber[len(la.AccountNumber)-4:])
	paymentUpdate := payment.Payment{
		ID:                  p.PaymentID,
		Status:              payment.PartnerTransferStatusToPaymentStatus[mm.Status],
		SourcePaymentID:     &mm.Id,
		LinkedBankAccountID: &laID,
		PaymentDate:         &mm.Created,
		PaymentToken:        &paymentToken,
		CardBrand:           la.BankName,
		CardLast4:           &accNumber,
	}
	err = payment.NewPaymentService(sourceRequest).UpdatePaymentStatus(&paymentUpdate)
	if err != nil {
		log.Println("Error updating payment status ", err)
		return err
	}

	//4. Update money request table
	reqType := payment.PaymentRequestTypeInvoiceBank
	requestUpdate := payment.RequestUpdate{
		ID:          shared.PaymentRequestID(*p.MoneyRequestID),
		Status:      payment.MoveMoneyStatusToRequestStatus[mm.Status],
		RequestType: &reqType,
	}
	err = payment.NewPaymentService(sourceRequest).UpdateRequestStatus(&requestUpdate)
	if err != nil {
		log.Println("Error updating request status ", err)
		return err
	}

	// Update payment object
	p, err = payment.NewPaymentServiceInternal().GetPaymentInfo(paymentToken)
	if err != nil {
		log.Println("error fetching payment token", err)
	}

	return nil
}

func createPaymentToken(invoiceID id.InvoiceID) (*string, error) {
	//create bus_mon_req_pay and return the payment Token
	paymentResp, err := payment.NewPaymentServiceInternal().CreatePaymentForInvoice(invoiceID.UUIDString())
	if err != nil {
		return nil, err
	}
	return paymentResp.PaymentToken, nil
}

func getInvoiceIDFromToken(token string, w http.ResponseWriter) (id.InvoiceID, error) {
	var invoiceID id.InvoiceID
	invoiceIDByte, err := base64.RawURLEncoding.DecodeString(token)
	if err != nil {
		return invoiceID, err
	}
	invoiceIDStr := string(invoiceIDByte)
	invoiceIDParsed, err := id.ParseInvoiceID(invoiceIDStr)
	if err != nil {
		return invoiceID, err
	}
	return invoiceIDParsed, nil
}

func isInvoiceV2Payment(t HTMLTemplate, r *http.Request) bool {
	if t == HTMLTemplateInvoiceRequestIndex || t == HTMLTemplateInvoiceView ||
		t == HTMLTemplateInvoiceSuccess {
		return true
	}
	return false
}

func renderError(err string, w http.ResponseWriter) {
	data := Request{
		ErrorMessage: err,
		SegmentKey:   os.Getenv("SEGMENT_WEB_WRITE_KEY"),
	}

	tmpl := template.Must(template.ParseFiles("error.html"))
	tmpl.Execute(w, data)
}

func renderInvoiceView(invoiceID id.InvoiceID, w http.ResponseWriter) {
	invSvc, err := invoice.NewInvoiceService()
	if err != nil {
		renderError(err.Error(), w)
		return
	}
	inv, err := invSvc.GetGrpcInvoiceByID(invoiceID)

	funcMap := getInvoiceTemplateFuncMap(inv.BusinessDetail.Logo)
	if !inv.PaymentTypes.AllowCards.ToBool() && !inv.PaymentTypes.AllowBankTransfer.ToBool() && !inv.ShowBankAccount {
		inv.PaymentLink = ""
	}
	// In template parsing, Funcs is supposed to go before ParseFiles
	tmpl := template.Must(template.New(string(HTMLTemplateInvoiceView)).Funcs(funcMap).ParseFiles(string(HTMLTemplateInvoiceView)))
	if err != nil {
		renderError(err.Error(), w)
		return
	}
	tmpl.Execute(w, inv)
}

func renderInvoiceReceipt(invoiceID id.InvoiceID, w http.ResponseWriter) {
	invSvc, err := invoice.NewInvoiceService()
	if err != nil {
		renderError(err.Error(), w)
		return
	}
	inv, err := invSvc.GetGrpcInvoiceByID(invoiceID)
	if err != nil {
		renderError(err.Error(), w)
		return
	}
	// including payment in process as there might be error while transferring the money
	// but still the receipt is sent, so we need show the receipt
	if inv.Status != grpcInvoice.InvoiceRequestStatus_IRT_PAID &&
		inv.Status != grpcInvoice.InvoiceRequestStatus_IRT_PROCESSING_PAYMENT {
		renderError("invoice payment is not complete", w)
		return
	}

	resp := &InvoiceSuccessResponse{
		Amount:                inv.TotalAmount,
		BusinessName:          inv.BusinessDetail.Name,
		CustomerName:          inv.ContactDetails.Name,
		InvoiceTitle:          inv.Title,
		IsWalletInfoAvailable: false,
		IsCardInfoAvailable:   false,
	}

	receiptInfo, err := payment.NewPaymentServiceInternal().GetPaymentReceiptFromInvoice(invoiceID)
	if err != nil {
		renderError(err.Error(), w)
		return
	}

	if receiptInfo.WalletType != nil && *receiptInfo.WalletType != "" {
		fmt.Println(receiptInfo.WalletType)
		if val, ok := payment.PaymentWalletTypeMap[*receiptInfo.WalletType]; ok {
			resp.WalletType = val
			resp.IsWalletInfoAvailable = true
		}
	}
	if receiptInfo.ReceiptNumber != nil {
		resp.ReceiptNumber = *receiptInfo.ReceiptNumber
	}

	if receiptInfo.CardBrand != nil && *receiptInfo.CardBrand != "" &&
		receiptInfo.CardLast4 != nil && *receiptInfo.CardLast4 != "" {
		resp.CardType = *receiptInfo.CardBrand
		resp.CardNumber = *receiptInfo.CardLast4
		resp.IsCardInfoAvailable = true
	}
	funcMap := getHTMLTemplateInvoiceSuccessFuncMap(receiptInfo.PaymentDate)
	tmpl := template.Must(template.New(string(HTMLTemplateInvoiceSuccess)).Funcs(funcMap).ParseFiles(string(HTMLTemplateInvoiceSuccess)))
	if err != nil {
		renderError(err.Error(), w)
		return
	}
	tmpl.Execute(w, resp)
}

func getInvoiceSuccessResponse(r *http.Request, p *payment.PaymentResponse) *InvoiceSuccessResponse {
	resp := &InvoiceSuccessResponse{
		BusinessName:          p.BusinessName,
		Amount:                shared.FormatFloatAmount(p.Amount),
		InvoiceTitle:          *p.InvoiceTitle,
		WalletType:            "",
		ReceiptNumber:         "",
		IsWalletInfoAvailable: false,
	}
	resp.IsCardInfoAvailable = false
	if params, ok := r.URL.Query()["resp"]; ok {
		cardDetails, err := base64.StdEncoding.DecodeString(params[0])
		if err == nil {
			cardData := strings.Split(string(cardDetails), ",")
			if len(cardData) == 2 {
				resp.CardNumber = cardData[0]
				resp.CardType = cardData[1]
				resp.IsCardInfoAvailable = true
			}
		}
	}

	sourceRequest := services.NewSourceRequest()
	contactResp, err := contact.NewContactService(sourceRequest).GetByIDInternal(*p.ContactID)
	if err == nil {
		resp.CustomerName = contactResp.Name()
	}
	return resp
}

func getInvoiceTemplateFuncMap(businessLogo string) template.FuncMap {
	funcMap := template.FuncMap{
		"GetFormatedCreated": func(timestamp *timestamp.Timestamp) string {
			created, err := ptypes.Timestamp(timestamp)
			if err != nil {
				log.Println(fmt.Sprintf("Created time error: %v", err.Error()))
				return "NA"
			}
			return created.Format("Jan 02, 2006")
		},
		"ShowBusinessLogo": func() string {
			if strings.TrimSpace(businessLogo) == "" {
				return "false"
			}
			return "true"
		},
	}

	return funcMap
}

func getHTMLTemplateInvoiceSuccessFuncMap(paymentDate *time.Time) template.FuncMap {
	funcMap := template.FuncMap{
		"GetPaymentDate": func() string {
			if paymentDate == nil {
				log.Println("invoice payment date not available")
				return ""
			}
			return paymentDate.Format("Jan 02, 2006")
		},
		"IsPaymentDateAvailable": func() bool {
			return paymentDate != nil && !paymentDate.IsZero()
		},
	}

	return funcMap
}
