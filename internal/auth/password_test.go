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

	t.Run("hash a regular password", func(t *testing.T) {
		hash, err := a.HashPassword("my-password")
		assert.NoError(t, err)
		assert.NotEmpty(t, hash)
	})

	t.Run("hash an empty password", func(t *testing.T) {
		hash, err := a.HashPassword("")
		assert.NoError(t, err)
		assert.NotEmpty(t, hash)
	})
}

func TestAuthority_VerifyPassword(t *testing.T) {
	const (
		correctPassword = "correct-password"
		wrongPassword   = "wrong-password"
	)

	a := newTestPasswordAuthenticator()
	correctPasswordHash, err := a.HashPassword(correctPassword)
	if err != nil {
		panic(err)
	}

	t.Run("verify correct password", func(t *testing.T) {
		ok, err := a.VerifyPassword(correctPassword, correctPasswordHash)
		assert.NoError(t, err)
		assert.True(t, ok)
	})

	t.Run("verify wrong password", func(t *testing.T) {
		ok, err := a.VerifyPassword(wrongPassword, correctPasswordHash)
		assert.NoError(t, err)
		assert.False(t, ok)
	})

	t.Run("verify empty password", func(t *testing.T) {
		ok, err := a.VerifyPassword("", correctPasswordHash)
		assert.NoError(t, err)
		assert.False(t, ok)
	})

	t.Run("verify correct password over an empty hash", func(t *testing.T) {
		ok, err := a.VerifyPassword(correctPasswordHash, "")
		assert.ErrorIs(t, err, bcrypt.ErrHashTooShort)
		assert.False(t, ok)
	})
}
