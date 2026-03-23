package lib

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

type JSONDate struct {
	time.Time
}

const dateLayout = "2006-01-02"

func (d JSONDate) MarshalJSON() ([]byte, error) {
	if d.Time.IsZero() {
		return []byte("null"), nil
	}
	return json.Marshal(d.Format(dateLayout))
}

func (d *JSONDate) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	parsed, err := time.Parse(dateLayout, s)
	if err != nil {
		return err
	}

	d.Time = parsed
	return nil
}

func (d *JSONDate) Scan(value any) error {
	switch v := value.(type) {
	case time.Time:
		d.Time = v
		return nil
	case []byte:
		t, err := time.Parse(dateLayout, string(v))
		if err != nil {
			return err
		}
		d.Time = t
		return nil
	case string:
		t, err := time.Parse(dateLayout, v)
		if err != nil {
			return err
		}
		d.Time = t
		return nil
	case nil:
		d.Time = time.Time{}
		return nil
	default:
		return fmt.Errorf("cannot scan %T into JSONDate", value)
	}
}

func (d JSONDate) Value() (driver.Value, error) {
	if d.Time.IsZero() {
		return nil, nil
	}
	return d.Format(dateLayout), nil
}

func (d *JSONDate) Parse(sdate string) error {
	parsed, err := time.Parse(dateLayout, sdate)
	if err != nil {
		return err
	}

	d.Time = parsed
	return nil
}

func JSONDateNil() JSONDate {
	var t time.Time
	d := JSONDate{}
	d.Time = t
	return d
}
