// Package validators contains useful functions to validate request params.
package validators

import (
	"errors"
	"fmt"
	"github.com/jieggii/groshi/internal/http/handles/datatypes"
)

// user username settings
const usernamePattern = ".+"

// user password settings
const minPasswordLen = 8
const maxPasswordLen = 128

// transaction description settings
const minDescriptionLen = 1
const maxDescriptionLen = 999

func ValidateUserPassword(password string) error {
	if len(password) < minPasswordLen || len(password) > maxPasswordLen {
		return fmt.Errorf("password must contain from %v to %v characters", minPasswordLen, maxPasswordLen)
	}
	return nil
}

func ValidateUserUsername(username string) error {
	return nil
}

func ValidateCurrency(currency datatypes.Currency) error {
	if !currency.IsKnown {
		return errors.New("unknown currency")
	}
	return nil
}

func ValidateTransactionDescription(description string) error {
	if len(description) < minDescriptionLen || len(description) > maxDescriptionLen {
		return fmt.Errorf("description must contain from %v to %v characters", minDescriptionLen, maxDescriptionLen)
	}
	return nil
}
