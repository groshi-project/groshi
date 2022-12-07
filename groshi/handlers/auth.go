package handlers

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type credentials struct {
	Username string
	Password string
}

func Auth(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	//request.Body
}
