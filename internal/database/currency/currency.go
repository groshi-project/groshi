package currency

import (
	"encoding/json"
)

// UnknownCurrencyError is returned by Currency JSON unmarshaler when unmarshalling unknown currency.
type UnknownCurrencyError struct{}

func (e *UnknownCurrencyError) Error() string {
	return "unknown currency"
}

type Currency string

func (c *Currency) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	if s == "" {
		return nil // returning no error when currency is empty
		// because it should be checked after unmarshalling
	}

	for _, currency := range currencies {
		if s == currency {
			*c = Currency(currency)
			return nil
		}
	}
	return &UnknownCurrencyError{}
}
