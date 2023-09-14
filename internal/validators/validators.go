package validators

import (
	"github.com/go-playground/validator/v10"
	"github.com/groshi-project/groshi/internal/currency/currency_rates"
	"github.com/groshi-project/groshi/internal/loggers"
	"regexp"
	"time"
)

// NonzeroTimeValidator ensures that date is not a zero date.
var NonzeroTimeValidator = func(fl validator.FieldLevel) bool {
	if fl.Field().Interface().(time.Time).IsZero() {
		return false
	}
	return true
}

// GetCurrencyValidator returns validator function for currencies.
func GetCurrencyValidator() validator.Func {
	currencies, err := currency_rates.FetchCurrencies()
	if err != nil {
		loggers.Error.Fatalf("could not fetch available currencies: %v", err)
	}

	return func(fl validator.FieldLevel) bool {
		for _, currency := range currencies {
			if fl.Field().Interface().(string) == currency {
				return true
			}
		}
		return false
	}
}

// GetRegexValidator returns validator function which checks if string matches regex pattern.
func GetRegexValidator(pattern *regexp.Regexp) validator.Func {
	return func(fl validator.FieldLevel) bool {
		return pattern.MatchString(fl.Field().Interface().(string))
	}
}
