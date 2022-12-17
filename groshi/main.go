package main

import (
	"fmt"
	"github.com/jieggii/groshi/groshi/auth/jwt"
	"github.com/jieggii/groshi/groshi/config"
	"github.com/jieggii/groshi/groshi/database"
	"github.com/jieggii/groshi/groshi/handlers"
	"github.com/jieggii/groshi/groshi/loggers"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func setupHandles(router *httprouter.Router) {
	router.Handle("POST", "/auth", handlers.Auth)
}

func startHTTPServer(host string, port int) error {
	router := httprouter.New()
	setupHandles(router)

	loggers.Info.Printf("Starting HTTP server on %v:%v.\n", host, port)
	return http.ListenAndServe(fmt.Sprintf("%v:%v", host, port), router)
}

func initializeApp(cfg *config.Config) error {
	jwt.SecretKey = cfg.JWTSecretKey
	if err := database.Connect(
		cfg.PostgresHost,
		cfg.PostgresPort,
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresDatabase,
	); err != nil {
		return fmt.Errorf("failed to connect to PostgreSQL database \"%v\" at %v:%v: %v", cfg.PostgresDatabase, cfg.PostgresHost, cfg.PostgresPort, err)
	}
	if err := database.Initialize(cfg.SuperuserUsername, cfg.SuperuserPassword); err != nil {
		return fmt.Errorf("failed to initialize PostgreSQL database \"%v\" at %v:%v: %v", cfg.PostgresDatabase, cfg.PostgresHost, cfg.PostgresPort, err)
	}
	return nil
}

func main() {
	loggers.Info.Println("Starting groshi server.")
	cfg := config.ReadFromEnv()
	if err := initializeApp(cfg); err != nil {
		loggers.Fatal.Fatalf("Error initializing groshi: %v.\n", err)
	}
	if err := startHTTPServer(cfg.Host, cfg.Port); err != nil {
		loggers.Fatal.Fatalf("Error starting HTTP server: %v.", err)
	}
}
