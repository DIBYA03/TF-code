package locale

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

type Date time.Time

func (d Date) Format() string {
	return time.Time(d).Format("2006-01-02")
}

func (d Date) FormatCardExpirationDate() string {
	return time.Time(d).Format("2006-01")
}

func ParseDate(s string) (Date, error) {
	t, err := time.Parse("2006-01-02", s)
	return Date(t), err
}

func ParseCardExpirationDate(s string) (Date, error) {
	t, err := time.Parse("2006-01", s)
	return Date(t), err
}

func (d Date) Time() time.Time {
	return time.Time(d)
}

func (d *Date) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Format())
}

func (d *Date) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	date, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}

	*d = Date(date)
	return nil
}

func (d Date) IsZero() bool {
	return d.Time().IsZero()
}

// SQL value marshaller
func (d Date) Value() (driver.Value, error) {
	return d.Time(), nil
}

func (d *Date) Scan(src interface{}) error {
	*d = Date(src.(time.Time))
	return nil
}

func NewDate(year, month, day int) Date {
	return Date(time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC))
}
