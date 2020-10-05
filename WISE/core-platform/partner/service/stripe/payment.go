/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

package stripe

import (
	"fmt"
	"log"
	"os"

	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/paymentintent"
	"github.com/stripe/stripe-go/terminal/connectiontoken"
	"github.com/wiseco/core-platform/partner/service"
)

type PaymentRequest struct {
	Amount          float64
	Currency        service.Currency
	PaymentMethod   service.PaymentMethod
	Descriptor      string
	ReceiptEmail    string
	CaptureMethod   string
	PaymentMetadata PaymentMetadata
}

type PaymentMetadata struct {
	BusinessName      string
	BusinessOwnerName string
	AvailableBalance  float64
	IPAddress         string
	PaymentMethod     string
}

type PaymentResponse struct {
	IntentID     string
	Status       service.PaymentStatus
	ClientSecret string
}

type PaymentStatus string

const (
	PaymentStatusRequiresPaymentMethod = PaymentStatus("requires_payment_method")
	PaymentStatusRequiresConfirmation  = PaymentStatus("requires_confirmation")
	PaymentStatusRequiresAction        = PaymentStatus("requires_action")
	PaymentStatusRequiresProcessing    = PaymentStatus("processing")
	PaymentStatusRequiresCapture       = PaymentStatus("requires_capture")
	PaymentStatusRequiresCanceled      = PaymentStatus("canceled")
	PaymentStatusSucceeded             = PaymentStatus("succeeded")
)

var partnerPaymentStatusTo = map[PaymentStatus]service.PaymentStatus{
	PaymentStatusRequiresPaymentMethod: service.PaymentStatusRequiresPaymentMethod,
	PaymentStatusRequiresConfirmation:  service.PaymentStatusRequiresConfirmation,
	PaymentStatusRequiresAction:        service.PaymentStatusRequiresAction,
	PaymentStatusRequiresProcessing:    service.PaymentStatusRequiresProcessing,
	PaymentStatusRequiresCapture:       service.PaymentStatusRequiresCapture,
	PaymentStatusRequiresCanceled:      service.PaymentStatusRequiresCanceled,
	PaymentStatusSucceeded:             service.PaymentStatusSucceeded,
}

type stripeService struct {
	request *service.APIRequest
}

type StripeService interface {
	// Create payment intent object
	CreatePayment(PaymentRequest) (*PaymentResponse, error)

	// Fetches client secret in payment intent object
	GetClientSecret(string) (*string, error)

	GetPaymentIntent(string) (*stripe.PaymentIntent, error)

	GetConnectionToken() (*string, error)

	CapturePayment(string) error
}

func NewStripeService(request *service.APIRequest) StripeService {
	return &stripeService{
		request: request,
	}
}

func (p *stripeService) CreatePayment(request PaymentRequest) (*PaymentResponse, error) {
	stripe.Key = os.Getenv("STRIPE_KEY")

	params := &stripe.PaymentIntentParams{
		// Converting cents to dollars. Stripe uses smallest common currency unit
		Amount: stripe.Int64(int64(request.Amount * 100)),

		Currency: stripe.String(string(request.Currency)),
		PaymentMethodTypes: stripe.StringSlice([]string{
			string(request.PaymentMethod),
		}),
		StatementDescriptor: &request.Descriptor,
		ReceiptEmail:        &request.ReceiptEmail,
	}

	if request.CaptureMethod != "" {
		params.CaptureMethod = stripe.String(request.CaptureMethod)
	}

	params.AddMetadata("business_name", request.PaymentMetadata.BusinessName)
	params.AddMetadata("business_owner_name", request.PaymentMetadata.BusinessOwnerName)
	params.AddMetadata("available_balance", fmt.Sprintf("%f", request.PaymentMetadata.AvailableBalance))
	params.AddMetadata("ip_address", request.PaymentMetadata.IPAddress)
	params.AddMetadata("type", request.PaymentMetadata.PaymentMethod)

	intent, err := paymentintent.New(params)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	response := PaymentResponse{
		IntentID:     intent.ID,
		Status:       partnerPaymentStatusTo[PaymentStatus(string(intent.Status))],
		ClientSecret: intent.ClientSecret,
	}

	return &response, nil

}

func (p *stripeService) GetClientSecret(id string) (*string, error) {
	stripe.Key = os.Getenv("STRIPE_KEY")

	intent, err := paymentintent.Get(id, nil)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &intent.ClientSecret, nil

}

func (p *stripeService) GetPaymentIntent(paymentIntentID string) (*stripe.PaymentIntent, error) {
	stripe.Key = os.Getenv("STRIPE_KEY")

	intent, err := paymentintent.Get(paymentIntentID, nil)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return intent, nil
}

func (p *stripeService) GetConnectionToken() (*string, error) {
	stripe.Key = os.Getenv("STRIPE_KEY")

	params := &stripe.TerminalConnectionTokenParams{}
	ct, err := connectiontoken.New(params)
	if err != nil {
		return nil, err
	}

	return &ct.Secret, nil

}

func (p *stripeService) CapturePayment(intentID string) error {
	stripe.Key = os.Getenv("STRIPE_KEY")

	params := &stripe.PaymentIntentCaptureParams{}
	_, err := paymentintent.Capture(intentID, params)

	return err

}
