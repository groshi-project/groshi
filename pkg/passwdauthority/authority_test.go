package passwdauthority

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const bcryptCost = 1

func TestAuthority_HashPassword(t *testing.T) {
	a := New(bcryptCost)

	// [a.HashPassword] can hash an empty password, so check what you are passing!
	hash1, err := a.HashPassword("")
	assert.NoError(t, err)
	assert.NotEmpty(t, hash1)

	hash2, err := a.HashPassword("my-password")
	assert.NoError(t, err)
	assert.NotEmpty(t, hash2)
}

func TestAuthority_VerifyPassword(t *testing.T) {
	a := New(bcryptCost)

	rightPassword := "password-123"
	wrongPassword := "wrong lol"

	hash, err := a.HashPassword(rightPassword)
	assert.NoError(t, err)
	assert.NotEmpty(t, hash)

	ok := a.VerifyPassword(wrongPassword, hash)
	assert.False(t, ok)

	ok = a.VerifyPassword(rightPassword, hash)
	assert.True(t, ok)
}
