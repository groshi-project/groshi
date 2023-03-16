package datatypes

import (
	"encoding/json"
	"time"
)

// InvalidISO8601DateError is returned by ISO8601Date JSON unmarshaler
// when unmarshalling invalid date in ISO-8601 format.
type InvalidISO8601DateError struct{}

func (e *InvalidISO8601DateError) Error() string {
	return ""
}

type ISO8601Date struct {
	time.Time
}

func (d *ISO8601Date) UnmarshalJSON(b []byte) error {
	if err := json.Unmarshal(b, &d.Time); err != nil {
		return new(InvalidISO8601DateError)
	}
	return nil
}
