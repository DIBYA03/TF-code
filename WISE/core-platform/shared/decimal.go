package shared

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/ericlagergren/decimal"
	"github.com/wiseco/go-lib/num"
	"golang.org/x/text/message"
)

const (
	MaxIntegralDigits   = 131072 // max digits before the decimal point
	MaxFractionalDigits = 16383  // max digits after the decimal point
)

// LengthError is returned from Decimal.Value when either its integral (digits
// before the decimal point) or fractional (digits after the decimal point)
// parts are too long for PostgresSQL.
type LengthError struct {
	Part string // "integral" or "fractional"
	N    int    // length of invalid part
	max  int
}

func (e LengthError) Error() string {
	return fmt.Sprintf("%s (%d digits) is too long (%d max)", e.Part, e.N, e.max)
}

// Decimal is a PostgreSQL DECIMAL. Its zero value is valid for use with both
// Value and Scan.
type Decimal struct {
	V     *decimal.Big
	Round bool // round if the decimal exceeds the bounds for DECIMAL
	Zero  bool // return "0" if V == nil
}

// Value implements driver.Valuer.
func (d Decimal) Value() (driver.Value, error) {
	if d.V == nil {
		if d.Zero {
			return "0", nil
		}
		return nil, nil
	}

	v := d.V
	if v.IsNaN(0) {
		return "NaN", nil
	}

	if v.IsInf(0) {
		return nil, errors.New("Decimal.Value: DECIMAL does not accept Infinities")
	}

	dl := v.Precision()  // length of d
	sl := int(v.Scale()) // length of fractional part

	if il := dl - sl; il > MaxIntegralDigits {
		if !d.Round {
			return nil, &LengthError{Part: "integral", N: il, max: MaxIntegralDigits}
		}
		// Rounding down the integral part automatically chops off the fractional
		// part.
		return v.Round(MaxIntegralDigits).String(), nil
	}

	if sl > MaxFractionalDigits {
		if !d.Round {
			return nil, &LengthError{Part: "fractional", N: sl, max: MaxFractionalDigits}
		}
		v.Round(dl - (sl - MaxFractionalDigits))
	}

	return v.String(), nil
}

// Scan implements sql.Scanner.
func (d *Decimal) Scan(val interface{}) error {
	if d.V == nil {
		d.V = new(decimal.Big)
	}

	switch t := val.(type) {
	case string:
		if _, ok := d.V.SetString(t); !ok {
			if err := d.V.Context.Err(); err != nil {
				return err
			}
			return fmt.Errorf("Decimal.Scan: invalid syntax: %q", t)
		}

		d.Zero = d.V.Sign() == 0
		return nil
	case []byte:
		err := d.V.UnmarshalText(t)
		if err != nil {
			return err
		}

		d.Zero = d.V.Sign() == 0
		return nil
	default:
		return fmt.Errorf("Decimal.Scan: unknown value: %#v", val)
	}
}

func (d *Decimal) UnmarshalJSON(b []byte) error {
	if d.V == nil {
		d.V = new(decimal.Big)
	}

	err := json.Unmarshal(b, &d.V)
	if err != nil {
		return err
	}

	d.Zero = d.V.Sign() == 0
	return nil
}

func (d Decimal) MarshalJSON() ([]byte, error) {
	// Default to zero for nil internal value
	if d.V == nil {
		return json.Marshal("0")
	}

	return d.V.MarshalText()
}

var ContextDecimalFin = decimal.Context{
	Precision:     19,
	RoundingMode:  decimal.ToNearestEven,
	OperatingMode: decimal.GDA,
	Traps:         ^(decimal.Inexact | decimal.Rounded | decimal.Subnormal),
	MaxScale:      4,
	MinScale:      -4,
}

func NewDecimalFin(v float64) (*decimal.Big, error) {
	return NewDecimal(ContextDecimalFin, v)
}

func NewDecimal(c decimal.Context, v float64) (*decimal.Big, error) {
	val := decimal.WithContext(c).SetFloat64(v)
	if val == nil {
		return nil, errors.New("invalid fixed decimal value")
	}

	if val.IsNaN(0) {
		return nil, errors.New("invalid fixed decimal value")
	}

	return val, nil
}

func (d *Decimal) IsNil() bool {
	return d == nil || (d.V == nil && !d.Zero)
}

func (d *Decimal) Float64() (f float64, ok bool) {
	return d.V.Float64()
}

func (d *Decimal) FormatCurrency() string {
	if d == nil {
		return "<nil>"
	}

	return fmt.Sprintf("%.2f", d.V)
}

func (d *Decimal) FormatCurrencySep() string {
	if d == nil {
		return "<nil>"
	}

	f, _ := d.Float64()
	p := message.NewPrinter(message.MatchLanguage("en"))
	return p.Sprintf("%.2f", f)
}

func (d *Decimal) Abs() *Decimal {
	if d == nil || d.IsNil() {
		return nil
	}

	b := *d.V
	return &Decimal{V: b.Abs(&b), Zero: b.Sign() == 0}
}

func (d *Decimal) IsNegative() bool {
	return !d.V.IsNaN(0) && d.V.Sign() < 0
}

func (d *Decimal) IsPositive() bool {
	return !d.V.IsNaN(0) && d.V.Sign() > 0
}

func (d Decimal) NumDecimal() num.Decimal {
	return num.Decimal{d.V, d.Round, d.Zero}
}

func ParseDecimal(s string) (Decimal, error) {
	d := Decimal{V: new(decimal.Big)}

	// Default to Zero on empty string
	if s == "" {
		b, err := NewDecimalFin(0)
		return Decimal{V: b, Zero: true}, err
	}

	if _, ok := d.V.SetString(s); !ok {
		if err := d.V.Context.Err(); err != nil {
			return d, err
		}

		return d, fmt.Errorf("parse decimal invalid syntax: %s", s)
	}

	d.Zero = d.V.Sign() == 0
	return d, nil
}

func (d Decimal) Add(y Decimal) Decimal {
	ret := NewZero()

	ret.V.Add(d.V, y.V)

	ret.Zero = ret.V.Sign() == 0

	return ret
}

func NewZero() Decimal {
	return Decimal{V: &decimal.Big{}, Zero: true}
}
