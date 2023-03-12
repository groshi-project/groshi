package middleware

import (
	"github.com/jieggii/groshi/internal/database"
	"github.com/jieggii/groshi/internal/http/ghttp"
	"github.com/jieggii/groshi/internal/http/ghttp/schema"
	"github.com/jieggii/groshi/internal/http/jwt"
	"net/http"
)

type _jwtFieldHolder struct {
	Token string `json:"token"`
}

// Middleware is the main middleware which must be used for all groshi handles.
// It validates request and ensures if the user is authorized if it is required.
func Middleware(authRequired bool, handle ghttp.Handle) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := ghttp.NewRequest(w, r)
		if req.RawRequest.Method != http.MethodPost {
			req.SendClientSideErrorResponse(
				schema.InvalidRequestErrorTag,
				"Invalid request method (POST must be used).",
			)
			return
		}

		var currentUser *database.User = nil

		if authRequired {
			jwtFieldHolder := _jwtFieldHolder{}
			if ok := req.Decode(&jwtFieldHolder); !ok {
				return
			}

			token := jwtFieldHolder.Token
			if token == "" {
				req.SendClientSideErrorResponse(
					schema.UnauthorizedErrorTag,
					"This method requires authorization, but required field `token` in the request body is missing.",
				)
				return
			}

			claims, ok := jwt.ParseJWT(token)
			if !ok {
				req.SendClientSideErrorResponse(
					schema.AccessDeniedErrorTag, "Invalid JWT.",
				)
				return
			}

			var err error
			currentUser, err = database.FetchUserByUsername(claims.Username)
			if err != nil {
				req.SendClientSideErrorResponse(
					schema.ObjectNotFoundErrorTag,
					"The user you are authorized under has not been found.",
				)
				return
			}
		}
		handle(req, currentUser)
	}
}
