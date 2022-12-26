package jwt

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jieggii/groshi/groshi/handles/schema"
	"github.com/jieggii/groshi/groshi/handles/util"
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
	http.ResponseWriter, *http.Request, *Claims,
)

type _requestWithJWT struct {
	Token string `json:"token"`
}

func ValidateJWTMiddleware(handle HandleWithJWTClaims) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := _requestWithJWT{}
		if ok := util.DecodeBodyJSON(w, r, &req); !ok {
			return
		}
		claims, err := ParseJWT(req.Token)
		if err != nil {
			fmt.Println(err)
			util.ReturnErrorResponse(w, schema.ClientSideError, "Invalid token.", nil)
			return
		}
		handle(w, r, claims)
	}
}
