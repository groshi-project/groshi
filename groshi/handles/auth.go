package handles

import (
	"github.com/jieggii/groshi/groshi/auth"
	"github.com/jieggii/groshi/groshi/database"
	"github.com/jieggii/groshi/groshi/handles/jwt"
	"github.com/jieggii/groshi/groshi/handles/util"
	"github.com/julienschmidt/httprouter"
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

func Auth(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	var credentials _request
	if ok := util.DecodeBodyJSON(writer, request, &credentials); !ok {
		return
	}
	if len(credentials.Username) == 0 || len(credentials.Password) == 0 {
		util.ReturnError(writer, http.StatusBadRequest, "Invalid request body.")
		return
	}

	user := new(database.User)
	err := database.Db.NewSelect().Model(user).Where("username = ?", credentials.Username).Scan(database.Ctx)
	if err != nil {
		util.ReturnError(writer, http.StatusUnauthorized, "User does not exist.")
		return
	}

	if !auth.CheckPasswordHash(credentials.Password, user.Password) {
		util.ReturnError(writer, http.StatusUnauthorized, "Invalid password.")
		return
	}

	token, err := jwt.GenerateJWT(credentials.Username)
	if err != nil {
		util.ReturnError(writer, http.StatusInternalServerError, "Could not generate JWT.")
		return
	}
	response := _response{Token: token, TTL: int(jwt.TTL / time.Second)}
	util.ReturnJSON(writer, http.StatusOK, &response)

	return
}
