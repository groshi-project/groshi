package datatypes

import (
	"encoding/json"
)

// UnknownCurrencyError is returned by Currency JSON unmarshaler when unmarshalling unknown currency.
type UnknownCurrencyError struct{}

func (e *UnknownCurrencyError) Error() string {
	return ""
}

type Currency string

func (c *Currency) UnmarshalJSON(b []byte) error {
	var stringCurrency string
	if err := json.Unmarshal(b, &stringCurrency); err != nil {
		return err
	}

	if stringCurrency == "" {
		return nil // returning no error when currency is empty
		// because it should be checked after unmarshalling
	}

	for _, currency := range currencies {
		if stringCurrency == currency {
			*c = Currency(stringCurrency)
			return nil
		}
	}

	return new(UnknownCurrencyError)
}
