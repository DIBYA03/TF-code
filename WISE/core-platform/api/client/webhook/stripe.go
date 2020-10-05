package webhook

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/webhook"
	"github.com/wiseco/core-platform/api"
	str "github.com/wiseco/core-platform/partner/service/stripe"
	"github.com/wiseco/core-platform/services"
	"github.com/wiseco/core-platform/services/payment"
)

//HandleWebhookRequest handles Stripe  requests
func HandleStripeRequest(r api.APIRequest) error {

	switch strings.ToUpper(r.HTTPMethod) {
	case http.MethodPost:
		return execute(r)
	default:
		return errors.New("Not supported")
	}

}

func execute(r api.APIRequest) error {

	endpointSecret := os.Getenv("STRIPE_WEBHOOK_SECRET")
	event, err := webhook.ConstructEvent([]byte(r.Body), r.Headers["Stripe-Signature"],
		endpointSecret)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error verifying webhook signature: %v\n", err)
		return err
	}

	switch event.Type {
	case "payment_intent.succeeded":
		return processPaymentIntentSucceeded(event.Data.Raw, r.SourceRequest())
	case "review.closed":
		return processReviewClosed(event.Data.Raw, r.SourceRequest())
	case "review.opened":
		return updatePaymentStatus(event.Data.Raw, payment.PaymentRequestStatusInProcess)
	default:
		return nil
	}
}

func processReviewClosed(data json.RawMessage, req services.SourceRequest) error {
	var review stripe.Review

	err := json.Unmarshal(data, &review)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing webhook JSON: %v\n", err)
		return err
	}

	if review.PaymentIntent != nil {
		switch review.Reason {
		case "approved":
			paymentIntent, err := str.NewStripeService(nil).GetPaymentIntent(review.PaymentIntent.ID)
			if err != nil {
				return err
			}

			return processPaymentIntent(*paymentIntent, req)
		default:
			return updatePaymentStatus(data, payment.PaymentRequestStatusPending)
		}
	}

	return nil
}

func processPaymentIntentSucceeded(data json.RawMessage, req services.SourceRequest) error {
	var paymentIntent stripe.PaymentIntent

	err := json.Unmarshal(data, &paymentIntent)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing webhook JSON: %v\n", err)
		return err
	}

	if paymentIntent.Charges != nil && paymentIntent.Charges.Data != nil {
		for _, charge := range paymentIntent.Charges.Data {

			if charge.Outcome != nil && charge.Outcome.Type == "authorized" {
				return processPaymentIntent(paymentIntent, req)
			}
		}
	}

	return nil
}

func processPaymentIntent(paymentIntent stripe.PaymentIntent, req services.SourceRequest) error {
	brand := ""
	last4 := ""
	walletType := ""
	var paidDate time.Time

	for _, data := range paymentIntent.Charges.Data {
		if data.Status == "succeeded" {
			if data.PaymentMethodDetails != nil {

				switch data.PaymentMethodDetails.Type {
				case "card_present":
					if data.PaymentMethodDetails.CardPresent != nil {
						brand = string(data.PaymentMethodDetails.CardPresent.Brand)
						last4 = data.PaymentMethodDetails.CardPresent.Last4
						if data.PaymentMethodDetails.Card != nil && data.PaymentMethodDetails.Card.Wallet != nil {
							walletType = string(data.PaymentMethodDetails.Card.Wallet.Type)
						}
					}
				case "card":
					if data.PaymentMethodDetails.Card != nil {
						brand = string(data.PaymentMethodDetails.Card.Brand)
						last4 = data.PaymentMethodDetails.Card.Last4
						if data.PaymentMethodDetails.Card != nil && data.PaymentMethodDetails.Card.Wallet != nil {
							walletType = string(data.PaymentMethodDetails.Card.Wallet.Type)
						}
					}

				default:
					e := fmt.Sprintf("Invalid payment type %s", data.PaymentMethodDetails.Type)
					return errors.New(e)
				}

				paidDate = time.Unix(data.Created, 0)

			}
		}

		if paymentIntent.Status == "succeeded" {
			resouce := payment.Payment{
				ID:          paymentIntent.ID,
				Status:      "succeeded",
				CardBrand:   &brand,
				CardLast4:   &last4,
				PaymentDate: &paidDate,
				WalletType:  &walletType,
			}

			err := payment.NewPaymentService(req).HandleWebhook(&resouce)
			if err != nil {
				return err
			}

		}

	}

	return nil
}

func updatePaymentStatus(data json.RawMessage, status payment.PaymentRequestStatus) error {
	var review stripe.Review

	err := json.Unmarshal(data, &review)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing webhook JSON: %v\n", err)
		return err
	}

	return payment.NewRequestService(services.NewSourceRequest()).UpdateRequestStatusByIntentID(review.PaymentIntent.ID, status)
}
