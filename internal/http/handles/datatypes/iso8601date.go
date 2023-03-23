package datatypes

import (
	"encoding/json"
	"time"
)

type ISO8601Date struct {
	Time time.Time

	IsValid bool
}

func (d *ISO8601Date) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &d.Time)
	if err != nil {
		d.IsValid = false
	} else {
		d.IsValid = true
	}
	return nil
}
