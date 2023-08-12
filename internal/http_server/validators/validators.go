package validators

import (
	"github.com/go-playground/validator/v10"
	"github.com/jieggii/groshi/internal/currency/currency_rates"
	"github.com/jieggii/groshi/internal/loggers"
	"regexp"
)

var UsernameValidator = getRegexValidator(regexp.MustCompile(".{2,40}"))
var PasswordValidator = getRegexValidator(regexp.MustCompile(".{2,250}"))

// GetCurrenciesValidator returns validator function for currencies.
func GetCurrenciesValidator() validator.Func {
	currencies, err := currency_rates.FetchCurrencies()
	if err != nil {
		loggers.Error.Fatalf("could not fetch available currencies: %v", err)
	}

	return func(fl validator.FieldLevel) bool {
		value, ok := fl.Field().Interface().(string)
		if !ok {
			loggers.Error.Printf("could not get value to be validated")
			return false
		}
		for _, currency := range currencies {
			if value == currency {
				return true
			}
		}
		return false
	}
}

// getRegexValidator returns validator function which checks if string matches regex pattern.
func getRegexValidator(pattern *regexp.Regexp) validator.Func {
	return func(fl validator.FieldLevel) bool {
		value, ok := fl.Field().Interface().(string)
		if !ok {
			loggers.Error.Printf("could not get value to be validated")
		}
		return pattern.MatchString(value)
	}
}
