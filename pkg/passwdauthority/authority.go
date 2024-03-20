package passwdauthority

import "golang.org/x/crypto/bcrypt"

// Authority represents password hashing and validation authority.
type Authority struct {
	bcryptCost int
}

// New creates a new instance of [Authority] and returns pointer to it.
func New(bcryptCost int) *Authority {
	return &Authority{bcryptCost: bcryptCost}
}

// HashPassword returns hash of a given password.
func (a *Authority) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), a.bcryptCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// VerifyPassword returns true if a given password matches with a given hash.
func (a *Authority) VerifyPassword(password string, hash string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return false
	}
	return true
}
