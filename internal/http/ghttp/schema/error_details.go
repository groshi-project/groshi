package schema

import (
	"fmt"
	"strings"
)

const (
	TransactionNotFoundErrorDetail           = "transaction not found"
	TransactionDoesNotBelongToYouErrorDetail = "this transaction does not belong to you"
)

func MissingRequiredFieldErrorDetail(field string) string {
	return fmt.Sprintf("missing required field %v", field)
}

func MissingRequiredFieldsErrorDetail(fields ...string) string {
	return fmt.Sprintf(
		"one of the required fields is missing: %v", strings.Join(fields, ", "),
	)
}

func AtLeastOneOfFieldsIsRequiredErrorDetail(fields ...string) string {
	return fmt.Sprintf(
		"at least one of these fields is required: %v", strings.Join(fields, ", "),
	)
}
