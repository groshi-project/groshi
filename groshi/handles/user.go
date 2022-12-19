package handles

import (
	"github.com/jieggii/groshi/groshi/handles/jwt"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func UserCreate(w http.ResponseWriter, r *http.Request, ps httprouter.Params, claims *jwt.Claims) {
}

func UserRead(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}

func UserUpdate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}

func UserDelete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}
