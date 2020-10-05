package shared

import "strconv"

func StringValue(v *string) string {
	if v == nil {
		return ""
	}

	return *v
}

func IntValue(v *int) int {
	if v == nil {
		return 0
	}

	return *v
}

func Float64Value(v *float64) float64 {
	if v == nil {
		return 0
	}

	return *v
}

func DecimalValue(v *Decimal) Decimal {
	if v == nil {
		return Decimal{V: nil, Zero: true}
	}

	return *v
}

func FormatFloatAmount(v float64) string {
	invoiceAmt, err := NewDecimalFin(v)
	if err != nil {
		return strconv.FormatFloat(v, 'f', 2, 64)
	}

	amount := Decimal{V: invoiceAmt}

	return amount.FormatCurrencySep()
}
