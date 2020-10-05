/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

// Package for business contacts
package contact

import (
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/wiseco/core-platform/services"
	"github.com/wiseco/core-platform/shared"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	stripe "github.com/wiseco/core-platform/partner/service/stripe"
	"github.com/wiseco/core-platform/services/banking/business"
	"github.com/wiseco/core-platform/services/data"
)

type paymentDatastore struct {
	token *string
	*sqlx.DB
}

type PaymentResponse struct {
	LegalName    *string              `db:"legal_name"`
	DBA          services.StringArray `db:"dba"`
	BusinessName string
	Amount       float64 `db:"amount"`
	Notes        string  `db:"notes"`
	ClientSecret string
	StripeKey    string
}

type PaymentService interface {
	// Fetch payment details
	GetPaymentInfo() (*PaymentResponse, error)
	GetPaymentByRequestID(shared.PaymentRequestID, shared.BusinessID) (*business.Payment, error)
}

func NewPaymentService(token *string) PaymentService {
	return &paymentDatastore{token, data.DBWrite}
}

func (db *paymentDatastore) GetPaymentByRequestID(requestID shared.PaymentRequestID, businessID shared.BusinessID) (*business.Payment, error) {
	p := business.Payment{}

	err := db.Get(
		&p,
		`SELECT* FROM business_money_request_payment WHERE request_id = $1 AND business_id`, requestID, businessID)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &p, nil

}

func (db *paymentDatastore) GetPaymentInfo() (*PaymentResponse, error) {
	// Fetch payment details
	p, err := db.getPaymentByToken(*db.token)
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

	paymentResponse := PaymentResponse{}

	err = db.Get(
		&paymentResponse,
		`SELECT business_money_request.amount, business_money_request.notes, business.legal_name, business.dba FROM business_money_request_payment
		JOIN business_money_request ON business_money_request_payment.request_id = business_money_request.id
		JOIN business ON business_money_request.business_id = business.id WHERE business_money_request_payment.token = $1`, p.PaymentToken)

	println("join values are ", paymentResponse.Amount, paymentResponse.BusinessName)

	clientSecret, err := stripe.NewStripeService(nil).GetClientSecret(*p.SourcePaymentID)
	if err != nil {
		return nil, err
	}

	paymentResponse.ClientSecret = *clientSecret
	paymentResponse.StripeKey = os.Getenv("STRIPE_PUBLISH_KEY")
	paymentResponse.BusinessName = shared.GetBusinessName(paymentResponse.LegalName, paymentResponse.DBA)

	// Send back client secret
	return &paymentResponse, nil

}

func (db *paymentDatastore) getPaymentByToken(token string) (*business.Payment, error) {
	p := business.Payment{}

	err := db.Get(&p, "SELECT * FROM business_money_request_payment WHERE token = $1", token)
	if err != nil && err != sql.ErrNoRows {
		log.Println(err)
		return nil, err
	}

	return &p, err
}
