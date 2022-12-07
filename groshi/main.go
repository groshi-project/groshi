package main

import (
	"fmt"
	"github.com/jieggii/groshi/groshi/config"
	"github.com/jieggii/groshi/groshi/database"
	"github.com/jieggii/groshi/groshi/handlers"
	"github.com/jieggii/groshi/groshi/jwt"
	"github.com/jieggii/groshi/groshi/logger"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func setupHandles(router *httprouter.Router) {
	router.Handle("POST", "/auth", handlers.Auth)
}

func startHTTPServer(host string, port int) {
	router := httprouter.New()
	setupHandles(router)

	logger.Info.Printf("Starting HTTP server on %v:%v.\n", host, port)
	err := http.ListenAndServe(fmt.Sprintf("%v:%v", host, port), router)
	if err != nil {
		logger.Fatal.Fatalln(err)
	}
}

func main() {
	logger.Info.Println("Starting groshi server.")
	cfg := config.ReadFromEnv()
	jwt.SecretKey = cfg.JWTSecretKey

	err := database.Initialize(
		cfg.PostgresHost,
		cfg.PostgresPort,
		cfg.PostgresUsername,
		cfg.PostgresPassword,
		cfg.PostgresDatabaseName,
	)
	if err != nil {
		logger.Fatal.Fatalf("Could not initialize PostgreSQL database \"%v\" at %v:%v (%v).", cfg.PostgresDatabaseName, cfg.PostgresHost, cfg.PostgresPort, err)
	}
	startHTTPServer(cfg.Host, cfg.Port)
}
