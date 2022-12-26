package handles

import (
	"github.com/jieggii/groshi/groshi/auth"
	"github.com/jieggii/groshi/groshi/database"
	"github.com/jieggii/groshi/groshi/handles/jwt"
	"github.com/jieggii/groshi/groshi/handles/schema"
	"github.com/jieggii/groshi/groshi/handles/util"
	"net/http"
	"time"
)

type _request struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type _response struct {
	Token string `json:"token"` // JWT
	TTL   int    `json:"TTL"`   // JWT TTL in seconds
}

func Auth(w http.ResponseWriter, r *http.Request) {
	var credentials _request
	if !util.DecodeBodyJSON(w, r, &credentials) {
		return
	}
	if !util.ValidateBody(w, len(credentials.Username) == 0 || len(credentials.Password) == 0) {
		return
	}

	user := new(database.User)
	err := database.Db.NewSelect().Model(user).Where("username = ?", credentials.Username).Scan(database.Ctx)
	if err != nil {
		util.ReturnErrorResponse(w, schema.ClientSideError, "User does not exist.", nil)
		return
	}

	if !auth.CheckPasswordHash(credentials.Password, user.Password) {
		util.ReturnErrorResponse(w, schema.ClientSideError, "Invalid password.", nil)
		return
	}

	token, err := jwt.GenerateJWT(credentials.Username)
	if err != nil {
		util.ReturnErrorResponse(w, schema.ServerSideError, "Could not generate JWT.", err)
		return
	}
	response := _response{Token: token, TTL: int(jwt.TTL / time.Second)}
	util.ReturnSuccessResponse(w, &response)
	return
}
