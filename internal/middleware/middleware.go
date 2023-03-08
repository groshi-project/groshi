package middleware

import (
	"github.com/jieggii/groshi/internal/auth/jwt"
	"github.com/jieggii/groshi/internal/database"
	"github.com/jieggii/groshi/internal/ghttp"
	"github.com/jieggii/groshi/internal/schema"
	"net/http"
)

type _JWTFieldHolder struct {
	Token string `json:"token"`
}

func Middleware(auth bool, handle ghttp.Handle) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := ghttp.NewRequest(w, r)
		if req.RawRequest.Method != http.MethodPost {
			req.SendClientSideErrorResponse(schema.InvalidRequestMethod)
			return
		}

		var currentUser *database.User = nil

		if auth {
			jwtFieldHolder := _JWTFieldHolder{}
			if ok := req.DecodeSafe(&jwtFieldHolder); !ok {
				return
			}

			token := jwtFieldHolder.Token
			if token == "" {
				req.SendClientSideErrorResponse(schema.MissingJWTField)
				return
			}

			claims, err := jwt.ParseJWT(token)
			if err != nil {
				req.SendClientSideErrorResponse(schema.InvalidJWT)
				return
			}
			currentUser, err = database.FetchUserByUsername(claims.Username)
			if err != nil {
				req.SendClientSideErrorResponse(schema.UserNotFound)
				return
			}
		}
		handle(req, currentUser)
	}
}
