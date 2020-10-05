package shared

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// Date object representation (YYYY-MM-DD)
type Date time.Time

func (d Date) String() string {
	return time.Time(d).Format("2006-01-02")
}

func (d Date) Time() time.Time {
	return time.Time(d)
}

func (d *Date) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
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

// Expiration Date object representation (YYYY-MM)
type ExpDate time.Time

func (d ExpDate) String() string {
	return time.Time(d).Format("2006-01")
}

func (d ExpDate) Time() time.Time {
	return time.Time(d)
}

func (d *ExpDate) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

func (d *ExpDate) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	date, err := time.Parse("2006-01", s)
	if err != nil {
		return err
	}

	*d = ExpDate(date)
	return nil
}

// SQL value marshaller
func (d ExpDate) Value() (driver.Value, error) {
	return d.Time(), nil
}

func (d *ExpDate) Scan(src interface{}) error {
	*d = ExpDate(src.(time.Time))
	return nil
}

func (d ExpDate) IsZero() bool {
	return d.Time().IsZero()
}
