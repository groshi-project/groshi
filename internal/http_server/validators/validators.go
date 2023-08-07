package validators

//
//// This can be useful: https://gin-gonic.com/docs/examples/custom-validators/
//
//import (
//	"github.com/go-playground/validator/v10"
//	"regexp"
//)
//
//func regexBasedValidatorFactory(pattern *regexp.Regexp, valueRequired bool) validator.Func {
//	return func(fl validator.FieldLevel) bool {
//		value, ok := fl.Field().Interface().(string)
//		if !ok {
//			panic("todo")
//		}
//		if !valueRequired && value == "" {
//			return true
//		}
//		return pattern.MatchString(value)
//	}
//}
//
//var Username validator.Func = func(fl validator.FieldLevel) bool {
//	return regexBasedValidatorFactory(
//		regexp.MustCompile(".+"), true, // todo
//	)(fl)
//}
//
//var Password validator.Func = func(fl validator.FieldLevel) bool {
//	return regexBasedValidatorFactory(
//		regexp.MustCompile(".+"), true, // todo
//	)(fl)
//}
//
//var TransactionDescription validator.Func = func(fl validator.FieldLevel) bool {
//	return regexBasedValidatorFactory(
//		regexp.MustCompile(".+"), true, // todo
//	)(fl)
//}
//
//var Currency validator.Func = func(fl validator.FieldLevel) bool {
//	value, ok := fl.Field().Interface().(string)
//	if !ok {
//		panic("todo") // todo
//	}
//	for _, currency := range currencies {
//		if value == currency {
//			return true
//		}
//	}
//	return false
//}
