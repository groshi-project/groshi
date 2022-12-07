package jwt

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

const TTL = 5 * time.Minute // todo: set convenient

var SecretKey []byte

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

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

//func GetJWT(headers http.Header) {
//
//}

func ParseJWT(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return SecretKey, nil
	})
	if err != nil {
		return nil, err // todo: check error type and return less detailed error
	}
	if !token.Valid {
		return nil, errors.New("invalid token") // todo
	}
	return claims, nil
}
