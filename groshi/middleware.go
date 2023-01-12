package main

import (
	"github.com/jieggii/groshi/groshi/auth/jwt"
	"github.com/jieggii/groshi/groshi/database"
	"github.com/jieggii/groshi/groshi/ghttp"
	"github.com/jieggii/groshi/groshi/handles/schema"
	"net/http"
)

type _JWTFieldHolder struct {
	Token string `json:"token"`
}

func middleware(auth bool, handle ghttp.Handle) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := ghttp.NewRequest(w, r)
		if req.RawRequest.Method != http.MethodPost {
			req.SendErrorResponse(schema.ClientSideError, "Invalid request method (use POST instead).", nil)
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
				req.SendErrorResponse(schema.ClientSideError, "Missing required field `token`.", nil)
				return
			}

			claims, err := jwt.ParseJWT(token)
			if err != nil {
				req.SendErrorResponse(schema.ClientSideError, "Invalid JWT.", nil)
				return
			}
			currentUser, err = database.FetchUserByUsername(claims.Username)
			if err != nil {
				req.SendErrorResponse(schema.ClientSideError, schema.UserNotFound, nil)
				return
			}
		}
		handle(req, currentUser)
	}
}
