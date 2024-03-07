package jwt

import (
	"context"
	"github.com/go-chi/jwtauth/v5"
	"time"
)

const algorithm = "HS256"

// Authority represents JWT authority.
type Authority struct {
	JWTAuth *jwtauth.JWTAuth

	jwtTTL time.Duration
}

// NewAuthority creates a new instance of [Authority] and returns pointer to it.
func NewAuthority(ttl time.Duration, secretKey string) *Authority {
	return &Authority{
		JWTAuth: jwtauth.New(algorithm, []byte(secretKey), nil),
		jwtTTL:  ttl,
	}
}

// GenerateToken generates a new jwt for a user with the given username.
// Returns token string and token's expiration date.
func (a *Authority) GenerateToken(username string) (string, time.Time, error) {
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

func (a *Authority) ExtractClaims(ctx context.Context) (map[string]any, error) {
	_, claims, err := jwtauth.FromContext(ctx)
	if err != nil {
		return nil, err
	}
	return claims, nil
}
