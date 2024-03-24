package auth

import (
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func newTestPasswordAuthenticator() *DefaultPasswordAuthenticator {
	return NewPasswordAuthenticator(1)
}

func TestAuthority_HashPassword(t *testing.T) {
	a := newTestPasswordAuthenticator()

	// [a.HashPassword] can hash an empty password, so check what you are passing!
	hash1, err := a.HashPassword("")
	assert.NoError(t, err)
	assert.NotEmpty(t, hash1)

	hash2, err := a.HashPassword("my-password")
	assert.NoError(t, err)
	assert.NotEmpty(t, hash2)
}

func TestAuthority_VerifyPassword(t *testing.T) {
	a := newTestPasswordAuthenticator()

	rightPassword := "password-123"
	wrongPassword := "wrong lol"

	hash, err := a.HashPassword(rightPassword)
	assert.NoError(t, err)
	assert.NotEmpty(t, hash)

	// test with wrong password:
	ok, err := a.VerifyPassword(wrongPassword, hash)
	assert.NoError(t, err)
	assert.False(t, ok)

	// test with right password:
	ok, err = a.VerifyPassword(rightPassword, hash)
	assert.NoError(t, err)
	assert.True(t, ok)

	// test with empty password:
	ok, err = a.VerifyPassword("", hash)
	assert.NoError(t, err)
	assert.False(t, ok)

	// test with empty hash:
	ok, err = a.VerifyPassword(rightPassword, "")
	assert.Error(t, bcrypt.ErrHashTooShort)
	assert.False(t, ok)
}
