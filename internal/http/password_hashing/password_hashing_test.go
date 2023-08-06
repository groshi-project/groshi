package password_hashing

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHashPassword(t *testing.T) {
	password := "123456789"
	hash, err := HashPassword(password)

	assert.NoError(t, err)
	assert.NotNil(t, hash)
}

func TestValidatePassword(t *testing.T) {
	password1 := "my-super-secret-password-123"
	password2 := "test-password-123"

	hash1, err := HashPassword(password1)
	assert.NoError(t, err)

	var ok bool
	ok = ValidatePassword(password1, hash1)
	assert.Equal(t, true, ok)

	ok = ValidatePassword(password2, hash1)
	assert.Equal(t, false, ok)
}
