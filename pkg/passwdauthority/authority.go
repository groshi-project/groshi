package passwdauthority

import "golang.org/x/crypto/bcrypt"

// PasswordAuthority represents password hashing and validation authority.
type PasswordAuthority struct {
	bcryptCost int
}

// New creates a new instance of [PasswordAuthority] and returns pointer to it.
func New(bcryptCost int) *PasswordAuthority {
	return &PasswordAuthority{bcryptCost: bcryptCost}
}

// HashPassword returns hash of a given password.
func (a *PasswordAuthority) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), a.bcryptCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// ValidatePassword returns true if a given password matches with a given hash.
func (a *PasswordAuthority) ValidatePassword(password string, hash string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return false
	}
	return true
}
