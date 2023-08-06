package password_hashing

import "golang.org/x/crypto/bcrypt"

const bcryptCost = 12 // todo: set convenient cost

// HashPassword returns hash of provided password.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	return string(bytes), err
}

// ValidatePassword returns true if provided password matches with provided hash.
// Otherwise, false is returned.
func ValidatePassword(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
