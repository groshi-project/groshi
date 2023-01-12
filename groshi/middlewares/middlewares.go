package middlewares

import (
	"github.com/jieggii/groshi/groshi/ghttp"
	"github.com/jieggii/groshi/groshi/handles/jwt"
	"github.com/jieggii/groshi/groshi/handles/schema"
	"net/http"
)

type _JWTFieldHolder struct {
	Token string `json:"token"`
}

func ParseRequest(handle ghttp.Handle) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := ghttp.NewRequest(w, r)
		if r.Method != http.MethodPost {
			req.SendErrorResponse(
				schema.ClientSideError, "Invalid request method.", nil,
			)
			return
		}

		handle(req, nil)
	}
}

func ValidateJWT(handle ghttp.Handle) ghttp.Handle {
	return func(req *ghttp.Request, _ *jwt.Claims) {
		jwtFieldHolder := _JWTFieldHolder{}
		if ok := req.DecodeSafe(&jwtFieldHolder); !ok {
			return
		}

		token := jwtFieldHolder.Token
		if ok := req.WrapCondition(token != "", "Missing required field `token`."); !ok {
			return
		}
		claims, err := jwt.ParseJWT(token)

		if ok := req.WrapCondition(err != nil, "Invalid JWT."); !ok {
			return
		}
		handle(req, claims)
	}
}
