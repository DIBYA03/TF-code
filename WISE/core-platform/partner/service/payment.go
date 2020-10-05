package service

type Currency string

func (c Currency) String() string {
	return string(c)
}

const (
	CurrencyUSD = Currency("usd")
)

type PaymentMethod string

func (c PaymentMethod) String() string {
	return string(c)
}

const (
	PaymentMethodCard = PaymentMethod("card")
)

type PaymentStatus string

const (
	PaymentStatusRequiresPaymentMethod = PaymentStatus("requiresPaymentMethod")
	PaymentStatusRequiresConfirmation  = PaymentStatus("requiresConfirmation")
	PaymentStatusRequiresAction        = PaymentStatus("requiresAction")
	PaymentStatusRequiresProcessing    = PaymentStatus("processing")
	PaymentStatusRequiresCapture       = PaymentStatus("requiresCapture")
	PaymentStatusRequiresCanceled      = PaymentStatus("canceled")
	PaymentStatusSucceeded             = PaymentStatus("succeeded")
)
