package bbva

import (
	"encoding/json"
	"time"
)

type DateTimeLocal time.Time

func (d DateTimeLocal) String() string {
	return time.Time(d).Format("2006-01-02T15:04:05")
}

func (d DateTimeLocal) Time() time.Time {
	return time.Time(d)
}

func (d *DateTimeLocal) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

func (d *DateTimeLocal) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	date, err := time.Parse("2006-01-02T15:04:05", s)
	if err != nil {
		date, err = time.Parse(time.RFC3339, s)
		if err != nil {
			return err
		}
	}

	*d = DateTimeLocal(date)
	return nil
}

func (d DateTimeLocal) IsZero() bool {
	return d.Time().IsZero()
}
