package validators

import (
	"github.com/go-playground/validator/v10"
	"github.com/groshi-project/groshi/internal/currency/rates"
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
// `isOptional` param indicates if the valid currency can be an empty string.
func GetCurrencyValidator(isOptional bool) validator.Func {
	currencies, err := currency_rates.GetCurrencies()
	if err != nil {
		loggers.Error.Fatalf("could not fetch available currencies: %v", err)
	}

	return func(fl validator.FieldLevel) bool {
		currency := fl.Field().Interface().(string)
		if currency == "" { // return true on zero currency if it is optional, otherwise false.
			return isOptional
		}

		for _, possibleCurrency := range currencies {
			if currency == possibleCurrency {
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
