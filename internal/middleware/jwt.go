package middleware

import (
	"context"
	"errors"
	"github.com/groshi-project/groshi/internal/auth"
	"github.com/groshi-project/groshi/internal/service/handler/httpresp"
	"github.com/groshi-project/groshi/internal/service/handler/model"
	"net/http"
	"strings"
)

const (
	UsernameContextVar  = "username" // todo: is it the right place for this const?
	authorizationHeader = "Authorization"
)

var (
	errEmptyOrMissingAuthHeader = errors.New("empty or missing authorization header")
	errInvalidAuthHeader        = errors.New("invalid authorization header")
)

// tokenFromHeader extracts token from a header value.
// For example, it will extract "some-token" from string "Bearer some-token".
func tokenFromHeader(headerValue string) (string, error) {
	if headerValue == "" {
		return "", errEmptyOrMissingAuthHeader
	}

	tokens := strings.SplitN(headerValue, " ", 3)
	if len(tokens) != 2 || tokens[0] != "Bearer" || tokens[1] == "" {
		return "", errInvalidAuthHeader
	}

	return tokens[1], nil
}

// NewJWT returns new JWT middleware which extracts and verifies JWT from authorization header.
// Additionally, sets [UsernameContextVar] context key to the authorized user's username.
func NewJWT(jwtAuthority auth.JWTAuthenticator) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// extract token from authorization header value:
			token, err := tokenFromHeader(r.Header.Get(authorizationHeader))
			if err != nil {
				switch {
				case errors.Is(err, errEmptyOrMissingAuthHeader):
					httpresp.Render(w, httpresp.New(http.StatusUnauthorized, model.NewError("empty or missing required authorization header")))
					return
				case errors.Is(err, errInvalidAuthHeader):
					httpresp.Render(w, httpresp.New(http.StatusBadRequest, model.NewError("invalid authorization header")))
					return
				default:
					panic(err)
				}
			}

			claims, err := jwtAuthority.VerifyToken(token)
			if err != nil {
				// todo: is it safe to display error?
				httpresp.Render(w, httpresp.New(http.StatusUnauthorized, model.NewError(err.Error())))
				return
			}

			ctx := context.WithValue(r.Context(), UsernameContextVar, claims["username"])
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
