package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

func handle_some(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	_, err := fmt.Fprint(w, "Hello, world!")
	if err != nil {
		return
	}
}

func main() {
	router := httprouter.New()
	router.GET("/", handle_some)
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", router))
}
