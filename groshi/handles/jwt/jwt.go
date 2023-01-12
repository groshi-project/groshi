package jwt

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

const TTL = 31 * 24 * time.Hour // todo: set convenient

var SecretKey []byte

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// Validate validates claims provided by user.
// username check happens only if expectedUsername != ""
// todo: validate ValidAfter, etc.
//func (claims *Claims) Validate(expectedUsername string) error {
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

//type HandleWithJWTClaims func(
//	http.ResponseWriter, *http.Request, *Claims,
//)

//func ValidateJWTMiddleware(handle HandleWithJWTClaims) http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		jwtFieldHolder := _JWTFieldHolder{}
//		req, ok := util.NewSafelyParsedRequest(w, r, &jwtFieldHolder)
//		if !ok {
//			return
//		}
//		token := jwtFieldHolder.Token
//		if ok = req.WrapCondition(token != "", "Missing required field `token`."); !ok {
//			return
//		}
//		claims, err := ParseJWT(token)
//		if err != nil {
//			req.SendErrorResponse(schema.ClientSideError, "Invalid JWT.", nil)
//			return
//		}
//		handle(w, r, claims)
//	}
//}
