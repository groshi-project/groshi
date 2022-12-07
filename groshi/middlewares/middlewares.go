package middlewares

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func UnmarshalMiddleware(handle httprouter.Handle, schema interface{}) httprouter.Handle {
	return func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		handle(writer, request, params)
	}
}

//
//func ValidateRequestMiddleware(handle httprouter.Handle) httprouter.Handle {
//	return func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
//		fmt.Println(params)
//		handle(writer, request, params)
//	}
//}

func AuthMiddleware() {

}
