package handlers

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type _credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Auth(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	var credentials _credentials
	if ok := decodeBodyJSON(writer, request, &credentials); !ok {
		return
	}
	if len(credentials.Username) == 0 || len(credentials.Password) == 0 {
		returnError(writer, http.StatusBadRequest, "Invalid bruh")
		return
	}
}
