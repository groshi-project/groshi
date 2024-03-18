package jwtauthority

import (
	"context"
	"github.com/go-chi/jwtauth/v5"
	"time"
)

// hashing algorithm used for JWT verification.
const algorithm = "HS256"

// JWTAuthority represents JWT authority.
type JWTAuthority struct {
	JWTAuth *jwtauth.JWTAuth

	jwtTTL time.Duration
}

// New creates a new instance of [JWTAuthority] and returns pointer to it.
func New(ttl time.Duration, secretKey string) *JWTAuthority {
	return &JWTAuthority{
		JWTAuth: jwtauth.New(algorithm, []byte(secretKey), nil),
		jwtTTL:  ttl,
	}
}

// GenerateToken generates a new jwt for a user with the given username.
// Returns token string and token's expiration date.
func (a *JWTAuthority) GenerateToken(username string) (string, time.Time, error) {
	// setup claims:
	claims := make(map[string]any)

	now := time.Now()
	jwtauth.SetIssuedAt(claims, now)

	expires := now.Add(a.jwtTTL)
	jwtauth.SetExpiry(claims, expires)

	claims["username"] = username

	// generate token:
	_, tokenString, err := a.JWTAuth.Encode(claims)
	if err != nil {
		return "", expires, err
	}

	return tokenString, expires, nil
}

// extractClaims extracts claims from provided context.
func (a *JWTAuthority) extractClaims(ctx context.Context) (map[string]any, error) {
	_, claims, err := jwtauth.FromContext(ctx)
	if err != nil {
		return nil, err
	}
	return claims, nil
}

// ExtractUsername extracts "username" claim from provided context.
func (a *JWTAuthority) ExtractUsername(ctx context.Context) (string, error) {
	claims, err := a.extractClaims(ctx)
	if err != nil {
		return "", err
	}
	return claims["username"].(string), nil
}
