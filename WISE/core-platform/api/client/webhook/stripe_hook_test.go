package webhook

import (
	"encoding/json"
	"github.com/stripe/stripe-go"
	"github.com/wiseco/core-platform/test"
	"testing"
)

func TestStripeHook(t *testing.T) {
	data := `{
    "id": "pi_1GxwQ8Du1MErS7u8S1pCr6rL",
    "object": "payment_intent",
    "amount": 220,
    "amount_capturable": 0,
    "amount_received": 220,
    "application": null,
    "application_fee_amount": null,
    "canceled_at": null,
    "cancellation_reason": null,
    "capture_method": "automatic",
    "charges": {
      "object": "list",
      "data": [
        {
          "id": "ch_1GxwQXDu1MErS7u84Bdyp4Ui",
          "object": "charge",
          "amount": 220,
          "amount_refunded": 0,
          "application": null,
          "application_fee": null,
          "application_fee_amount": null,
          "balance_transaction": "txn_1GxwQYDu1MErS7u8PkVQUK2M",
          "billing_details": {
            "address": {
              "city": null,
              "country": null,
              "line1": null,
              "line2": null,
              "postal_code": null,
              "state": null
            },
            "email": null,
            "name": "arindam",
            "phone": null
          },
          "calculated_statement_descriptor": "BUSINESS TESTING",
          "captured": true,
          "created": 1593096705,
          "currency": "usd",
          "customer": null,
          "description": null,
          "destination": null,
          "dispute": null,
          "disputed": false,
          "failure_code": null,
          "failure_message": null,
          "fraud_details": {
          },
          "invoice": null,
          "livemode": false,
          "metadata": {
            "ip_address": "[::1]:52982",
            "type": "card",
            "business_name": "Business Testing",
            "business_owner_name": "Alberta  Charleson",
            "available_balance": "10763.170000"
          },
          "on_behalf_of": null,
          "order": null,
          "outcome": {
            "network_status": "approved_by_network",
            "reason": null,
            "risk_level": "normal",
            "risk_score": 34,
            "seller_message": "Payment complete.",
            "type": "authorized"
          },
          "paid": true,
          "payment_intent": "pi_1GxwQ8Du1MErS7u8S1pCr6rL",
          "payment_method": "pm_1GxwQXDu1MErS7u8tf82SmQL",
          "payment_method_details": {
            "card": {
              "brand": "visa",
              "checks": {
                "address_line1_check": null,
                "address_postal_code_check": null,
                "cvc_check": "pass"
              },
              "country": "US",
              "exp_month": 11,
              "exp_year": 2022,
              "fingerprint": "NLRoyt9IjUNk7mzd",
              "funding": "credit",
              "installments": null,
              "last4": "4242",
              "network": "visa",
              "three_d_secure": null,
              "wallet": null
            },
            "type": "card"
          },
          "receipt_email": "noreply@wise.us",
          "receipt_number": null,
          "receipt_url": "https://pay.stripe.com/receipts/acct_1EIlHuDu1MErS7u8/ch_1GxwQXDu1MErS7u84Bdyp4Ui/rcpt_HX0R4Wy8pdydrzCbyUYhhmmnpUOTRpE",
          "refunded": false,
          "refunds": {
            "object": "list",
            "data": [
            ],
            "has_more": false,
            "total_count": 0,
            "url": "/v1/charges/ch_1GxwQXDu1MErS7u84Bdyp4Ui/refunds"
          },
          "review": null,
          "shipping": null,
          "source": null,
          "source_transfer": null,
          "statement_descriptor": "Business Testing",
          "statement_descriptor_suffix": null,
          "status": "succeeded",
          "transfer_data": null,
          "transfer_group": null
        }
      ],
      "has_more": false,
      "total_count": 1,
      "url": "/v1/charges?payment_intent=pi_1GxwQ8Du1MErS7u8S1pCr6rL"
    },
    "client_secret": "pi_1GxwQ8Du1MErS7u8S1pCr6rL_secret_vGuiYcxF80qzxAKOvSlB0EmP4",
    "confirmation_method": "automatic",
    "created": 1593096680,
    "currency": "usd",
    "customer": null,
    "description": null,
    "invoice": null,
    "last_payment_error": null,
    "livemode": false,
    "metadata": {
      "ip_address": "[::1]:52982",
      "type": "card",
      "business_name": "Business Testing",
      "business_owner_name": "Alberta  Charleson",
      "available_balance": "10763.170000"
    },
    "next_action": null,
    "on_behalf_of": null,
    "payment_method": "pm_1GxwQXDu1MErS7u8tf82SmQL",
    "payment_method_options": {
      "card": {
        "installments": null,
        "network": null,
        "request_three_d_secure": "automatic"
      }
    },
    "payment_method_types": [
      "card"
    ],
    "receipt_email": "noreply@wise.us",
    "review": null,
    "setup_future_usage": null,
    "shipping": null,
    "source": null,
    "statement_descriptor": "Business Testing",
    "statement_descriptor_suffix": null,
    "status": "succeeded",
    "transfer_data": null,
    "transfer_group": null
  }`
	apReq := test.TestRequest("post")
	apReq.Body = data

	var paymentIntent stripe.PaymentIntent

	err := json.Unmarshal([]byte(data), &paymentIntent)
	err = processPaymentIntent(paymentIntent, apReq.SourceRequest())
	println(err)
}
