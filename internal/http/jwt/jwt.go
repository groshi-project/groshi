package jwt

import (
	"github.com/golang-jwt/jwt/v4"
	"time"
)

const TTL = 31 * 24 * time.Hour // todo: set convenient

var SecretKey []byte

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// Before validates claims provided by user.
// username check happens only if expectedUsername != ""
// todo: validate ValidAfter, etc.
//func (claims *Claims) Before(expectedUsername string) error {
//	if expectedUsername != "" {
//		if claims.Username != expectedUsername {
//			return errors.New("no access to this user")
//		}
//	}
//	return nil
//}

func GenerateJWT(username string) (string, error) {
	expirationTime := time.Now().Add(TTL)
	claims := &Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(SecretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ParseJWT(tokenString string) (*Claims, bool) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(
		tokenString, claims,
		func(token *jwt.Token) (interface{}, error) {
			return SecretKey, nil
		},
	)
	if err != nil {
		return nil, false
	}
	if !token.Valid {
		return nil, false
	}
	return claims, true
}
