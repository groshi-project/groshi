package auth

import "golang.org/x/crypto/bcrypt"

// PasswordAuthenticator represents password hashing and validation authority.
type PasswordAuthenticator struct {
	bcryptCost int
}

// NewPasswordAuthenticator creates a new instance of [PasswordAuthenticator] and returns pointer to it.
func NewPasswordAuthenticator(bcryptCost int) *PasswordAuthenticator {
	return &PasswordAuthenticator{bcryptCost: bcryptCost}
}

// HashPassword returns hash of a given password.
func (a *PasswordAuthenticator) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), a.bcryptCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// VerifyPassword returns true if a given password matches with a given hash.
func (a *PasswordAuthenticator) VerifyPassword(password string, hash string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return false
	}
	return true
}
