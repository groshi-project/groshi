package auth

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
)

type PasswordAuthenticator interface {
	// HashPassword returns hash of a given password.
	HashPassword(password string) (string, error)

	// VerifyPassword returns true if a given password matches with a given hash.
	VerifyPassword(password string, hash string) (bool, error)
}

// DefaultPasswordAuthenticator represents password hashing and validation authority.
type DefaultPasswordAuthenticator struct {
	bcryptCost int
}

// NewPasswordAuthenticator creates a new instance of [DefaultPasswordAuthenticator] and returns pointer to it.
func NewPasswordAuthenticator(bcryptCost int) *DefaultPasswordAuthenticator {
	return &DefaultPasswordAuthenticator{bcryptCost: bcryptCost}
}

// HashPassword returns hash of a given password.
func (d *DefaultPasswordAuthenticator) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), d.bcryptCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// VerifyPassword returns true if a given password matches with a given hash.
func (d *DefaultPasswordAuthenticator) VerifyPassword(password string, hash string) (bool, error) {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
