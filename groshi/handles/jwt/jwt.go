package jwt

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jieggii/groshi/groshi/handles/util"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"time"
)

const TTL = 10 * 365 * 24 * time.Hour // todo: set convenient

var SecretKey []byte

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// Validate validates claims provided by user.
// username check happens only if expectedUsername != ""
// todo: validate ValidAfter, etc.
func (claims *Claims) Validate(expectedUsername string) error {
	if expectedUsername != "" {
		if claims.Username != expectedUsername {
			return errors.New("no access to this user")
		}
	}
	return nil
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

func ParseJWT(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return SecretKey, nil
	})
	if err != nil {
		return nil, errors.New("could not parse token")
	}
	if !token.Valid {
		return nil, errors.New("invalid token") // todo
	}
	return claims, nil
}

type HandleWithJWTClaims func(
	http.ResponseWriter, *http.Request, httprouter.Params, *Claims,
)

func ValidateJWTMiddleware(handle HandleWithJWTClaims) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		token, err := parseJWTHeader(r.Header)
		fmt.Println(token)
		if err != nil {
			util.ReturnError(w, http.StatusBadRequest, "invalid headers todo")
			return
		}
		claims, err := ParseJWT(token)
		if err != nil {
			fmt.Println(err)
			util.ReturnError(w, http.StatusUnauthorized, "invalid JWT")
			return
		}
		handle(w, r, ps, claims)
	}
}
