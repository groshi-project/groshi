package main

import (
	"github.com/jieggii/groshi/groshi/config"
	"github.com/jieggii/groshi/groshi/database"
	"github.com/jieggii/groshi/groshi/logger"
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

func readEnv() {

}

func main() {
	logger.Info.Println("Starting groshi server.")
	_, _ = config.ReadFromEnv()
	host := "localhost"
	port := 27017
	dbName := "groshi"
	if err := database.Connect(host, port, dbName); err != nil {
		logger.Fatal.Fatalf("Could not connect to mongodb database \"%v\" at %v:%v (%v).", dbName, host, port, err)
	}
	StartHTTPServer(":8080")
}
