package passhash

import "golang.org/x/crypto/bcrypt"

const bcryptCost = 12 // todo: set convenient cost

// Hash returns hash of provided password.
func Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	return string(bytes), err
}

// Validate returns true if provided password matches with provided hash.
// Otherwise, false is returned.
func Validate(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
