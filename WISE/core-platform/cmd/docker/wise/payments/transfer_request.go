package main

import (
	"database/sql"
	"encoding/base64"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/wiseco/core-platform/services"
	"github.com/wiseco/core-platform/services/banking"
	"github.com/wiseco/core-platform/services/banking/business"
	"github.com/wiseco/core-platform/services/banking/business/contact"
	con "github.com/wiseco/core-platform/services/contact"
	"github.com/wiseco/core-platform/services/payment"
	"github.com/wiseco/core-platform/shared"
)

type ErrorCode int

const (
	ErrorCodeBadRequest       = ErrorCode(400)
	ErrorCodeRegister         = ErrorCode(512)
	ErrorCodeMoveMoney        = ErrorCode(513)
	ErrorCodeInsufficientFund = ErrorCode(514)
	ErrorCodeGeneric          = ErrorCode(515)
	ErrorCodeAccountBlocked   = ErrorCode(516)
)

const (
	ErrorMessageInsufficientFund           = "Insufficient funds to initiate move money"
	ErrorMessageRegisterAccountFailure     = "There was an error registering account"
	ErrorMessageRegisterCardFailure        = "There was an error registering card"
	ErrorMessagePaymentFailure             = "There was an error while making payment"
	ErrorMessageInvalidExpirationDate      = "Invalid expiration date format"
	ErrorMessageParseExpirationDateFailure = "Error parsing expiration date"
	ErrorMessageInvalidUrl                 = "Requested url is invalid"
)

type InstantPayError struct {
	Error     error
	ErrorCode ErrorCode
}

type Country string

const (
	CountryUSA = Country("USA")
)

func loadTransferRequestHTMLTemplate(t HTMLTemplate, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	var data TransferRequest
	var p *payment.TransferRequestResponse
	var err error

	if t == HTMLTemplateBankPaymentSuccess || t == HTMLTemplateCardPaymentSuccess {
		params, ok := r.URL.Query()["token"]
		if !ok || len(params[0]) < 1 {

			data := Request{
				ErrorMessage: ErrorMessageInvalidUrl,
				SegmentKey:   os.Getenv("SEGMENT_WEB_WRITE_KEY"),
			}

			tmpl := template.Must(template.ParseFiles("error.html"))
			tmpl.Execute(w, data)

			return
		}

		// decode business name
		businessName, err := base64.StdEncoding.DecodeString(params[0])
		if err != nil {
			log.Println("error decoding business name", businessName)

			data := Request{
				ErrorMessage: ErrorMessageInvalidUrl,
				SegmentKey:   os.Getenv("SEGMENT_WEB_WRITE_KEY"),
			}

			tmpl := template.Must(template.ParseFiles("error.html"))
			tmpl.Execute(w, data)

			return
		}

		params, ok = r.URL.Query()["id"]
		if !ok || len(params[0]) < 1 {

			data := Request{
				ErrorMessage: ErrorMessageInvalidUrl,
				SegmentKey:   os.Getenv("SEGMENT_WEB_WRITE_KEY"),
			}

			tmpl := template.Must(template.ParseFiles("error.html"))
			tmpl.Execute(w, data)

			return
		}

		// decode user ID
		userID, err := base64.StdEncoding.DecodeString(params[0])
		if err != nil {
			log.Println("error decoding userID", userID)

			data := Request{
				ErrorMessage: ErrorMessageInvalidUrl,
				SegmentKey:   os.Getenv("SEGMENT_WEB_WRITE_KEY"),
			}

			tmpl := template.Must(template.ParseFiles("error.html"))
			tmpl.Execute(w, data)

			return
		}

		data = TransferRequest{
			BusinessName: string(businessName),
			UserID:       string(userID),
			SegmentKey:   os.Getenv("SEGMENT_WEB_WRITE_KEY"),
		}

	} else {
		params, ok := r.URL.Query()["token"]
		if !ok || len(params[0]) < 1 {
			data := Request{
				ErrorMessage: ErrorMessageInvalidUrl,
				SegmentKey:   os.Getenv("SEGMENT_WEB_WRITE_KEY"),
			}

			tmpl := template.Must(template.ParseFiles("error.html"))
			tmpl.Execute(w, data)

			return
		}

		sourceRequest := services.NewSourceRequest()

		p, err = payment.NewTransferService(sourceRequest).GetTransferRequestInfo(params[0])
		if err != nil && t != HTMLTemplateBankPaymentSuccess {
			// Return error html page
			log.Println("invalid token error", err)

			data := Request{
				ErrorMessage: err.Error(),
				SegmentKey:   os.Getenv("SEGMENT_WEB_WRITE_KEY"),
			}

			tmpl := template.Must(template.ParseFiles(string(HTMLTemplateError)))
			tmpl.Execute(w, data)

			return
		}

		amount := shared.FormatFloatAmount(p.Amount)

		p.Notes = strings.Replace(p.Notes, "\n", ", ", -1)

		data = TransferRequest{
			BusinessName: p.BusinessName,
			Amount:       amount,
			Notes:        p.Notes,
			PaymentToken: params[0],
			SegmentKey:   os.Getenv("SEGMENT_WEB_WRITE_KEY"),
			UserID:       p.UserID.ToPrefixString(),
		}
	}

	switch t {
	case HTMLTemplateGetStarted:
		tmpl := template.Must(template.ParseFiles(string(HTMLTemplateGetStarted)))
		tmpl.Execute(w, data)
	case HTMLTemplateCardPay:
		tmpl := template.Must(template.ParseFiles(string(HTMLTemplateCardPay)))
		tmpl.Execute(w, data)
	case HTMLTemplateBankPay:
		tmpl := template.Must(template.ParseFiles(string(HTMLTemplateBankPay)))
		tmpl.Execute(w, data)
	case HTMLTemplateCardPaymentSuccess:
		tmpl := template.Must(template.ParseFiles(string(HTMLTemplateCardPaymentSuccess)))
		tmpl.Execute(w, data)
	case HTMLTemplateBankPaymentSuccess:
		tmpl := template.Must(template.ParseFiles(string(HTMLTemplateBankPaymentSuccess)))
		tmpl.Execute(w, data)
	case HTMLTemplateCardNotSupported:
		tmpl := template.Must(template.ParseFiles(string(HTMLTemplateCardNotSupported)))
		tmpl.Execute(w, data)
	case HTMLTemplateLinkedCard:
		cardNumber := r.FormValue("cardNumber")
		expirationDate := r.FormValue("expirationDate")
		securityCode := r.FormValue("securityCode")

		allowToSave := r.FormValue("allowToSave")
		b, err := strconv.ParseBool(allowToSave)
		if err != nil {
			b = false
		}

		addressLine1 := r.FormValue("addressLine1")
		addressLine2 := r.FormValue("addressLine2")
		city := r.FormValue("city")
		state := r.FormValue("state")
		zip := r.FormValue("zip")

		linkedCardErr := transferMoneyToLinkedCard(p, b, cardNumber, expirationDate, securityCode, addressLine1, addressLine2, city, state, zip)
		if linkedCardErr != nil {
			log.Println("error in bank transfer ", linkedCardErr.Error, int(linkedCardErr.ErrorCode))
			http.Error(w, linkedCardErr.Error.Error(), int(linkedCardErr.ErrorCode))
			return
		}

		w.Write([]byte("success"))
	}

}

func transferMoneyToLinkedAccount(transferRequest *payment.TransferRequestResponse, allowToSave bool, accountNumber, routingNumber string) *InstantPayError {
	//This is currently not being used, I'm leaving this method in here in case we want to turn this on in the future
	//1. Link account
	l := business.ContactLinkedAccountCreate{
		UserID:        transferRequest.UserID,
		BusinessID:    transferRequest.BusinessID,
		ContactId:     *transferRequest.ContactID,
		AccountType:   banking.AccountTypeChecking,
		AccountNumber: business.AccountNumber(accountNumber),
		RoutingNumber: routingNumber,
		Currency:      banking.CurrencyUSD,
		Permission:    banking.LinkedAccountPermissionRecieveOnly,
	}

	sourceRequest := services.NewSourceRequest()
	sourceRequest.UserID = transferRequest.UserID

	var ut business.UsageType
	if allowToSave {
		ut = business.UsageTypeContact
		l.UsageType = &ut
	} else {
		ut = business.UsageTypeContactInvisible
		l.UsageType = &ut
	}

	la, err := contact.NewLinkedAccountService(sourceRequest).Create(&l)
	if la == nil {
		e := InstantPayError{
			Error:     errors.Wrap(err, ErrorMessageRegisterAccountFailure),
			ErrorCode: ErrorCodeRegister,
		}
		log.Println("error fetching linked account", err)
		return &e
	}

	//2. Move money
	m := business.TransferInitiate{
		BusinessID:      transferRequest.BusinessID,
		ContactId:       transferRequest.ContactID,
		CreatedUserID:   transferRequest.UserID,
		Amount:          transferRequest.Amount,
		Currency:        banking.CurrencyUSD,
		SourceAccountId: transferRequest.RegisteredAccountID,
		SourceType:      banking.TransferTypeAccount,
		DestAccountId:   la.Id,
		DestType:        banking.TransferTypeAccount,
		Notes:           &transferRequest.Notes,
	}
	mt, err := contact.NewMoneyTransferService(sourceRequest).Transfer(&m)
	if err != nil {
		log.Println("string check", err.Error())
		if err.Error() == ErrorMessageInsufficientFund {
			e := InstantPayError{
				Error:     err,
				ErrorCode: ErrorCodeInsufficientFund,
			}
			return &e
		} else if strings.Contains(err.Error(), "Source account is blocked from") {
			e := InstantPayError{
				Error:     errors.New(ErrorMessagePaymentFailure),
				ErrorCode: ErrorCodeAccountBlocked,
			}
			return &e
		} else {
			e := InstantPayError{
				Error:     errors.Wrap(err, ErrorMessagePaymentFailure),
				ErrorCode: ErrorCodeMoveMoney,
			}
			return &e
		}
	}

	//3. Reset token
	u := payment.TransferRequestUpdate{
		ID:              transferRequest.TransferRequestID,
		MoneyTransferID: &mt.Id,
	}
	err = payment.NewTransferService(sourceRequest).UpdateTransferRequest(u)
	if err != nil {
		log.Println("error updating transfer request")
		e := InstantPayError{
			Error:     errors.Wrap(err, ErrorMessagePaymentFailure),
			ErrorCode: ErrorCodeGeneric,
		}
		return &e
	}

	// 4. Update usage type if required
	if allowToSave && la.UsageType != nil && *la.UsageType != business.UsageTypeContact {
		u := business.LinkedAccountUpdate{
			ID:        la.Id,
			UsageType: &ut,
		}
		sourceRequest := services.NewSourceRequest()
		sourceRequest.UserID = transferRequest.UserID

		err := contact.NewLinkedAccountService(sourceRequest).UpdateLinkedAccountUsageType(&u)
		if err != nil {
			log.Println("error updating account usage")
			e := InstantPayError{
				Error:     errors.Wrap(err, ErrorMessagePaymentFailure),
				ErrorCode: ErrorCodeGeneric,
			}
			return &e
		}
	}

	return nil
}

func transferMoneyToLinkedCard(transferRequest *payment.TransferRequestResponse, allowToSave bool, cardNumber, expirationDate, securityCode,
	addressLine1, addressLine2, city, state, postalCode string) *InstantPayError {

	//1. Link account
	a := services.Address{
		StreetAddress: addressLine1,
		City:          city,
		State:         state,
		Country:       string(CountryUSA),
		PostalCode:    postalCode,
	}

	if len(addressLine2) > 0 {
		a.AddressLine2 = addressLine2
	}

	var contactName string
	switch transferRequest.ContactType {
	case con.ContactTypeBusiness:
		contactName = transferRequest.BusinessName
	case con.ContactTypePerson:
		contactName = *transferRequest.ContactFirstName + " " + *transferRequest.ContactLastName
	}

	expDate := strings.Split(expirationDate, "/")
	if len(expDate) != 2 {
		e := InstantPayError{
			Error:     errors.New(ErrorMessageInvalidExpirationDate),
			ErrorCode: ErrorCodeBadRequest,
		}
		log.Println("invalid expiration date format")
		return &e
	}

	dateFormat := expDate[0] + "01" + expDate[1]
	layout := "010206"
	t, err := time.Parse(layout, dateFormat)
	if err != nil {
		e := InstantPayError{
			Error:     errors.New(ErrorMessageParseExpirationDateFailure),
			ErrorCode: ErrorCodeBadRequest,
		}
		log.Println("error parsing expiration date")
		return &e
	}

	c := business.LinkedCardCreate{
		UserID:         transferRequest.UserID,
		BusinessID:     transferRequest.BusinessID,
		ContactId:      transferRequest.ContactID,
		CardNumber:     business.CardNumber(cardNumber),
		CVVCode:        securityCode,
		ExpirationDate: shared.ExpDate(t),
		CardHolderName: contactName,
		Alias:          contactName,
		Permission:     banking.LinkedAccountPermissionRecieveOnly,
		BillingAddress: &a,
	}
	sourceRequest := services.NewSourceRequest()
	sourceRequest.UserID = transferRequest.UserID

	var ut business.UsageType
	if allowToSave {
		ut = business.UsageTypeContact
		c.UsageType = &ut
	} else {
		ut = business.UsageTypeContactInvisible
		c.UsageType = &ut
	}

	var lc *business.LinkedCard
	hash := c.HashLinkedCard()

	lc, err = contact.NewLinkedCardService(sourceRequest).GetByLinkedCardHashAndContactID(transferRequest.BusinessID, *transferRequest.ContactID, *hash)
	if err != nil && err != sql.ErrNoRows {
		e := InstantPayError{
			Error:     errors.Wrap(err, ErrorMessageRegisterCardFailure),
			ErrorCode: ErrorCodeRegister,
		}
		log.Println("error fetching linked card contact id and hash", err)
		return &e
	}

	if lc == nil {
		//Lets check to see if this card has been registered by another contact
		lcs, err := contact.NewLinkedCardService(sourceRequest).GetByLinkedCardHash(transferRequest.BusinessID, *hash)
		if err != nil {
			e := InstantPayError{
				Error:     errors.Wrap(err, ErrorMessageRegisterCardFailure),
				ErrorCode: ErrorCodeRegister,
			}
			log.Println("error fetching linked card hash", err)
			return &e
		}

		//If we already have a linked card with this hash, lets not re-register it with bbva, but copy the info over
		if len(lcs) > 0 {
			lc, err = contact.NewLinkedCardService(sourceRequest).RegisterExistingCard(transferRequest.BusinessID, *transferRequest.ContactID, *hash)
			if err != nil {
				e := InstantPayError{
					Error:     errors.Wrap(err, ErrorMessageRegisterCardFailure),
					ErrorCode: ErrorCodeRegister,
				}
				log.Println("error registering existing linked card", err)
				return &e
			}
		} else {
			lc, err = contact.NewLinkedCardService(sourceRequest).Create(&c)
			if lc == nil {
				e := InstantPayError{
					Error:     errors.Wrap(err, ErrorMessageRegisterCardFailure),
					ErrorCode: ErrorCodeRegister,
				}
				log.Println("Error linking card ", err)
				return &e
			}
		}
	}

	//2. Move money
	m := business.TransferInitiate{
		BusinessID:      transferRequest.BusinessID,
		ContactId:       transferRequest.ContactID,
		CreatedUserID:   transferRequest.UserID,
		Amount:          transferRequest.Amount,
		Currency:        banking.CurrencyUSD,
		SourceAccountId: transferRequest.RegisteredAccountID,
		SourceType:      banking.TransferTypeAccount,
		DestAccountId:   lc.Id,
		DestType:        banking.TransferTypeCard,
		Notes:           &transferRequest.Notes,
	}
	mt, err := contact.NewMoneyTransferService(sourceRequest).Transfer(&m)
	if err != nil {
		log.Println("Error moving money", err.Error())
		if err.Error() == ErrorMessageInsufficientFund {
			e := InstantPayError{
				Error:     err,
				ErrorCode: ErrorCodeInsufficientFund,
			}
			return &e
		} else if strings.Contains(err.Error(), "Source account is blocked from") {
			e := InstantPayError{
				Error:     errors.New(ErrorMessagePaymentFailure),
				ErrorCode: ErrorCodeAccountBlocked,
			}
			return &e
		} else {
			e := InstantPayError{
				Error:     errors.Wrap(err, ErrorMessagePaymentFailure),
				ErrorCode: ErrorCodeMoveMoney,
			}
			return &e
		}
	}

	//3. Reset token
	u := payment.TransferRequestUpdate{
		ID:              transferRequest.TransferRequestID,
		MoneyTransferID: &mt.Id,
	}
	err = payment.NewTransferService(sourceRequest).UpdateTransferRequest(u)
	if err != nil {
		e := InstantPayError{
			Error:     errors.Wrap(err, ErrorMessagePaymentFailure),
			ErrorCode: ErrorCodeGeneric,
		}
		log.Println("error updating transfer request")
		return &e
	}

	// 4. Update usage type if required
	if allowToSave && lc.UsageType != nil && *lc.UsageType != business.UsageTypeContact {
		u := business.LinkedCardUpdate{
			ID:        lc.Id,
			UsageType: &ut,
		}
		sourceRequest := services.NewSourceRequest()
		sourceRequest.UserID = transferRequest.UserID

		err := contact.NewLinkedCardService(sourceRequest).UpdateLinkedCardUsageType(&u)
		if err != nil {
			e := InstantPayError{
				Error:     errors.Wrap(err, ErrorMessagePaymentFailure),
				ErrorCode: ErrorCodeGeneric,
			}
			log.Println("error updating card usage type", err)
			return &e
		}
	}

	return nil
}
