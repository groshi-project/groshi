package middleware

import (
	"github.com/jieggii/groshi/internal/database"
	"github.com/jieggii/groshi/internal/http/ghttp"
	"github.com/jieggii/groshi/internal/http/jwt"
	"github.com/jieggii/groshi/internal/http/schema"
	"net/http"
)

type _JWTFieldHolder struct {
	Token string `json:"token"`
}

// Middleware todo...
func Middleware(authRequired bool, handle ghttp.Handle) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := ghttp.NewRequest(w, r)
		if req.RawRequest.Method != http.MethodPost {
			req.SendClientSideErrorResponse(
				schema.InvalidRequestErrorTag, "Invalid request method (POST must be used)",
			)
			return
		}

		var currentUser *database.User = nil

		if authRequired {
			jwtFieldHolder := _JWTFieldHolder{}
			if ok := req.DecodeSafe(&jwtFieldHolder); !ok {
				return
			}

			token := jwtFieldHolder.Token
			if token == "" {
				req.SendClientSideErrorResponse(
					schema.UnauthorizedErrorTag,
					"This method requires authorization, but required `token` field is missing.",
				)
				return
			}

			claims, err := jwt.ParseJWT(token)
			if err != nil {
				req.SendClientSideErrorResponse(
					schema.AccessDeniedErrorTag, "Invalid JWT",
				)
				return
			}

			currentUser, err = database.FetchUserByUsername(claims.Username)
			if err != nil {
				req.SendClientSideErrorResponse(
					schema.ObjectNotFoundErrorTag,
					"The user you authorized yourself to was not found",
				)
				return
			}
		}
		handle(req, currentUser)
	}
}
