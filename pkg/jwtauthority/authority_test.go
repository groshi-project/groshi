package jwtauthority

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var signingMethod = jwt.SigningMethodHS256

const (
	secretKey    = "some-secret-key"
	longTokenTTL = time.Duration(10) * time.Minute
	zeroTokenTTL = time.Duration(0) * time.Nanosecond
)

func NewTestAuthority(tokenTTL time.Duration) *Authority {
	return New(signingMethod, secretKey, tokenTTL)
}

func TestAuthority_CreateToken(t *testing.T) {
	a := NewTestAuthority(longTokenTTL)

	token, expires, err := a.CreateToken("jieggii")
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.NotZero(t, expires)
}

func TestAuthority_VerifyToken(t *testing.T) {
	a1 := NewTestAuthority(longTokenTTL)

	// test valid token:
	token1, _, _ := a1.CreateToken("jieggii")
	claims1, err := a1.VerifyToken(token1)
	assert.NoError(t, err)
	assert.NotEmpty(t, claims1)

	// test expired token:
	a2 := NewTestAuthority(zeroTokenTTL)
	token2, _, _ := a2.CreateToken("jieggii")
	claims2, err := a2.VerifyToken(token2)
	if assert.Error(t, err) {
		assert.ErrorIs(t, err, jwt.ErrTokenExpired)
	}
	assert.Empty(t, claims2)

	// test malformed token:
	claims3, err := a1.VerifyToken("blablabla")
	if assert.Error(t, err) {
		assert.ErrorIs(t, err, jwt.ErrTokenMalformed)
	}
	assert.Empty(t, claims3)
}
