package auth

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

const (
	secretKey    = "some-secret-key"
	longTokenTTL = time.Duration(10) * time.Minute
	zeroTokenTTL = time.Duration(0) * time.Nanosecond
)

func NewTestJWTAuthenticator(tokenTTL time.Duration) *DefaultJWTAuthenticator {
	return NewJWTAuthenticator(secretKey, tokenTTL)
}

func TestAuthority_CreateToken(t *testing.T) {
	jwtAuth := NewTestJWTAuthenticator(longTokenTTL)

	token, expires, err := jwtAuth.CreateToken("test-username")
	if assert.NoError(t, err) {
		assert.NotEmpty(t, token)
		assert.NotZero(t, expires)
	}
}

func TestAuthority_VerifyToken(t *testing.T) {
	const testUsername = "test-username"

	t.Run("verify valid token", func(t *testing.T) {
		jwtAuth := NewTestJWTAuthenticator(longTokenTTL)

		// create a new token for the test user:
		token, _, _ := jwtAuth.CreateToken(testUsername)

		claims, err := jwtAuth.VerifyToken(token)
		if assert.NoError(t, err) {
			if assert.NotEmpty(t, claims) {
				assert.Equal(t, testUsername, claims[JWTClaimUsername])
			}
		}
	})

	t.Run("verify expired token", func(t *testing.T) {
		jwtAuth := NewTestJWTAuthenticator(zeroTokenTTL)

		// create a new token for the test user:
		token, _, _ := jwtAuth.CreateToken(testUsername)

		claims, err := jwtAuth.VerifyToken(token)
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, jwt.ErrTokenExpired)
		}
		assert.Empty(t, claims)
	})

	// todo: verify token which is not yet valid somehow

	t.Run("verify malformed token", func(t *testing.T) {
		jwtAuth := NewTestJWTAuthenticator(longTokenTTL)

		claims, err := jwtAuth.VerifyToken("malformed-token")
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, jwt.ErrTokenMalformed)
		}
		assert.Empty(t, claims)
	})

}
