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

func initializeApp(cfg *config.Config) {
	jwt.SecretKey = cfg.JWTSecretKey
	if err := database.Connect(
		cfg.PostgresHost,
		cfg.PostgresPort,
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresDatabase,
	); err != nil {
		logger.Fatal.Fatalf("Failed to connect to PostgreSQL database \"%v\" at %v:%v: %v.", cfg.PostgresDatabase, cfg.PostgresHost, cfg.PostgresPort, err)
	}
	if err := database.Initialize(); err != nil {
		logger.Fatal.Fatalf("Failed to initialize PostgreSQL database \"%v\" at %v:%v: %v.", cfg.PostgresDatabase, cfg.PostgresHost, cfg.PostgresPort, err)
	}
}

func main() {
	logger.Info.Println("Starting groshi server.")
	cfg := config.ReadFromEnv()
	initializeApp(cfg)
	startHTTPServer(cfg.Host, cfg.Port)
}
