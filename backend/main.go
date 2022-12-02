package main

import (
	"github.com/jieggii/groshi/backend/logger"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func StartHTTPServer(addr string) {
	router := httprouter.New()

	//router.GET("/", ...)

	err := http.ListenAndServe(addr, router)
	if err != nil {
		logger.Fatal.Fatalln(err)
	}
}

func main() {
	StartHTTPServer(":8080")
}
