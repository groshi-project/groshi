package main

import (
	"fmt"
	"github.com/jieggii/groshi/groshi/config"
	"github.com/jieggii/groshi/groshi/database"
	"github.com/jieggii/groshi/groshi/logger"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func StartHTTPServer(host string, port int) {
	router := httprouter.New()

	//router.GET("/", ...)

	logger.Info.Printf("Running HTTP server on %v:%v.\n", host, port)
	err := http.ListenAndServe(
		fmt.Sprintf("%v:%v", host, port),
		router,
	)
	if err != nil {
		logger.Fatal.Fatalln(err)
	}
}

func main() {
	logger.Info.Println("Starting groshi server.")
	cfg := config.ReadFromEnv()
	if err := database.Connect(cfg.MongoHost, cfg.MongoPort, cfg.MongoDBName); err != nil {
		logger.Fatal.Fatalf("Could not connect to the mongodb database \"%v\" at %v:%v (%v).", cfg.MongoDBName, cfg.MongoHost, cfg.MongoPort, err)
	}
	StartHTTPServer(cfg.Host, cfg.Port)
}
