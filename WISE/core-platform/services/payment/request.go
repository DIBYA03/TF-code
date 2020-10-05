package payment

import "github.com/wiseco/core-platform/shared"

type PaymentRequestResend struct {
	RequestID shared.PaymentRequestID `json:"requestId"`
}
