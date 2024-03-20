package jwtauthority

import (
	"github.com/golang-jwt/jwt/v4"
	"time"
)

// Authority represents JWT authority.
type Authority struct {
	// signing method used to sign token claims.
	signingMethod jwt.SigningMethod

	// secret key used to sign token claims.
	secretKey []byte

	// duration of a token validity.
	tokenTTL time.Duration
}

// New creates a new instance of [Authority] and returns pointer to it.
func New(signingMethod jwt.SigningMethod, secretKey string, tokenTTL time.Duration) *Authority {
	return &Authority{
		signingMethod: signingMethod,
		secretKey:     []byte(secretKey),
		tokenTTL:      tokenTTL,
	}
}

// CreateToken generates a new JWT and returns its string representation and expiration timestamp.
// todo: should `expires` be returned and is it necessary for a user?
func (a *Authority) CreateToken(username string) (string, time.Time, error) {
	issued := time.Now()
	expires := time.Now().Add(a.tokenTTL)
	token := jwt.NewWithClaims(a.signingMethod, jwt.MapClaims{
		"username": username,
		"exp":      expires.Unix(),
		"iat":      issued.Unix(),
	})

	tokenString, err := token.SignedString(a.secretKey)
	if err != nil {
		return "", expires, err
	}

	return tokenString, expires, nil
}

// VerifyToken verifies that JWT token is valid and not expired, returns claims it contains.
func (a *Authority) VerifyToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		// todo: verify signing method (https://pkg.go.dev/github.com/golang-jwt/jwt/v5#example-Parse-Hmac)
		return a.secretKey, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, jwt.ErrTokenMalformed
	}
	return claims, nil
}
